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

// validatePath validates the path to prevent path traversal attacks.
func validatePath(path string) error {
	// Check for path traversal (..)
	if strings.Contains(path, "..") {
		return fmt.Errorf("path contains '..' which is not allowed")
	}

	// Check for absolute paths
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute paths are not allowed")
	}

	// Check for invalid characters
	if strings.Contains(path, "\n") || strings.Contains(path, "\r") {
		return fmt.Errorf("path contains invalid characters")
	}

	return nil
}

// CreateArchive creates a tar.gz archive at outputPath containing the provided files.
func CreateArchive(outputPath string, files []File) error {
	if err := validatePath(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	f, err := os.Create(outputPath) // #nosec G304 -- path is validated above
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
		hdr := &tar.Header{
			Name: filepath.ToSlash(file.Name),
			Mode: 0o644,
			Size: int64(len(file.Data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("write tar header for %s: %w", file.Name, err)
		}
		if _, err := tw.Write(file.Data); err != nil {
			return fmt.Errorf("write tar data for %s: %w", file.Name, err)
		}
	}

	return nil
}

// ExtractArchive reads a tar.gz archive from path and returns its files.
func ExtractArchive(path string) ([]File, error) {
	if err := validatePath(path); err != nil {
		return nil, fmt.Errorf("invalid archive path: %w", err)
	}

	f, err := os.Open(path) // #nosec G304 -- path is validated above
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

		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("read file %s from archive: %w", hdr.Name, err)
		}

		out = append(out, File{
			Name: hdr.Name,
			Data: data,
		})
	}

	return out, nil
}
