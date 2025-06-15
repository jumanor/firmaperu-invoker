#!/bin/bash

# Obtener información del repo
VERSION=$(git describe --tags 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT="$(git rev-parse --abbrev-ref HEAD)-$(git rev-parse HEAD)"

# Función para construir en Linux
build_linux() {
    echo "Building for Linux..."
    GOOS=linux GOARCH=amd64 go build -ldflags "-s -w \
        -X main.Version=$VERSION \
        -X main.BuildTime=$BUILD_TIME \
        -X main.GitCommit=$GIT_COMMIT" \
        -o main
    echo "Built: ./main"
}

# Función para construir en Windows
build_windows() {
    echo "Building for Windows..."
    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w \
        -X main.Version=$VERSION \
        -X main.BuildTime=$BUILD_TIME \
        -X main.GitCommit=$GIT_COMMIT" \
        -o main.exe
    echo "Built: ./main.exe"
}

# Manejar argumentos
case "$1" in
    linux)
        build_linux
        ;;
    windows)
        build_windows
        ;;
    *)
        build_linux
        build_windows
        ;;
esac