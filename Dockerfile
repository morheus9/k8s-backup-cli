# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" \
    -o /kubectl-backup ./cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates \
    && addgroup -g 1000 -S backup \
    && adduser -u 1000 -S backup -G backup

USER backup:backup
COPY --from=builder --chown=backup:backup /kubectl-backup /usr/local/bin/kubectl-backup

ENTRYPOINT ["kubectl-backup"]