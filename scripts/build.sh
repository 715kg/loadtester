#!/bin/bash

echo "========================================"
echo "   Load Tester - Cross Platform Build"
echo "========================================"
echo

# Получаем версию из git или используем дефолтную
VERSION=$(git describe --tags --always 2>/dev/null || echo "v1.0.0")
echo "Building version: $VERSION"
echo

# Создаем директорию для сборок
mkdir -p dist/{windows,linux,macos}

# Получаем дополнительную информацию для сборки
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_USER=$(whoami 2>/dev/null || echo "unknown")
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "unknown")

# Устанавливаем расширенные флаги сборки
LDFLAGS="-ldflags -s -w -X main.version=$VERSION -X main.gitCommit=$GIT_COMMIT -X 'main.buildTime=$BUILD_TIME' -X main.buildUser=$BUILD_USER"
export CGO_ENABLED=0

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build $LDFLAGS -o dist/windows/loadtester-$VERSION-windows-amd64.exe ./cmd/loadtester
if [ $? -ne 0 ]; then
    echo "ERROR: Windows build failed"
    exit 1
fi

echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build $LDFLAGS -o dist/linux/loadtester-$VERSION-linux-amd64 ./cmd/loadtester
if [ $? -ne 0 ]; then
    echo "ERROR: Linux build failed"
    exit 1
fi

echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build $LDFLAGS -o dist/linux/loadtester-$VERSION-linux-arm64 ./cmd/loadtester
if [ $? -ne 0 ]; then
    echo "ERROR: Linux ARM64 build failed"
    exit 1
fi

echo "Building for macOS (amd64 - Intel)..."
GOOS=darwin GOARCH=amd64 go build $LDFLAGS -o dist/macos/loadtester-$VERSION-darwin-amd64 ./cmd/loadtester
if [ $? -ne 0 ]; then
    echo "ERROR: macOS Intel build failed"
    exit 1
fi

echo "Building for macOS (arm64 - Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build $LDFLAGS -o dist/macos/loadtester-$VERSION-darwin-arm64 ./cmd/loadtester
if [ $? -ne 0 ]; then
    echo "ERROR: macOS Apple Silicon build failed"
    exit 1
fi

echo
echo "========================================"
echo "          Build Complete!"
echo "========================================"
echo

echo "Built files:"
find dist -type f -exec basename {} \;

echo
echo "File sizes:"
find dist -type f -exec ls -lh {} \; | awk '{print $9 " - " $5}'

echo
echo "Build completed successfully!"
echo "Files are located in the 'dist' directory."

# Делаем Unix файлы исполняемыми
chmod +x dist/linux/*
chmod +x dist/macos/*

echo "Unix executables have been made executable."