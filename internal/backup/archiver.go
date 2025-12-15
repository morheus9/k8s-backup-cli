package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// File represents a named file to be written into or read from an archive.
type File struct {
	Name string
	Data []byte
}

// validatePath ensures the path is within allowed boundaries
func validatePath(path string) (string, error) {
	// Clean the path to resolve any ../ or ./ sequences
	cleanPath := filepath.Clean(path)

	// Ensure the path doesn't start with .. or /
	if strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return cleanPath, nil
}

// CreateArchive creates a tar.gz archive at outputPath containing the provided files.
func CreateArchive(outputPath string, files []File) error {
	// Validate output path
	outputPath = filepath.Clean(outputPath)

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	gw := gzip.NewWriter(f)
	defer func() {
		_ = gw.Close()
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		_ = tw.Close()
	}()

	for _, file := range files {
		// Validate each file name to prevent path traversal
		cleanName, err := validatePath(file.Name)
		if err != nil {
			return fmt.Errorf("invalid filename %s: %w", file.Name, err)
		}

		hdr := &tar.Header{
			Name: filepath.ToSlash(cleanName),
			Mode: 0o644,
			Size: int64(len(file.Data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("write tar header for %s: %w", cleanName, err)
		}
		if _, err := tw.Write(file.Data); err != nil {
			return fmt.Errorf("write tar data for %s: %w", cleanName, err)
		}
	}

	return nil
}

// ExtractArchive reads a tar.gz archive from path and returns its files.
func ExtractArchive(path string) ([]File, error) {
	// Validate input path
	path = filepath.Clean(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("create gzip reader: %w", err)
	}
	defer func() {
		_ = gr.Close()
	}()

	tr := tar.NewReader(gr)
	var out []File

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar header: %w", err)
		}

		// Validate file names from archive to prevent path traversal
		cleanName, err := validatePath(hdr.Name)
		if err != nil {
			return nil, fmt.Errorf("invalid filename in archive: %s: %w", hdr.Name, err)
		}

		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("read file %s from archive: %w", cleanName, err)
		}

		out = append(out, File{
			Name: cleanName,
			Data: data,
		})
	}

	return out, nil
}
