#!/bin/bash

set -e

# Ensure paths are relative to the script location
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$DIR/.."

# Build information
VERSION="${VERSION:-v0.1.0}"
BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
GIT_BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')"

echo "=== OSK Release Build ==="
echo "Version: $VERSION"
echo "Build Date: $BUILD_DATE"
echo "Git Commit: $GIT_COMMIT"
echo "Git Branch: $GIT_BRANCH"
echo "======================"

# Create dist directory
mkdir -p $ROOT_DIR/dist

# Build ldflags with version information
LDFLAGS="-s -w -X 'main.Version=$VERSION' -X 'main.BuildDate=$BUILD_DATE' -X 'main.GitCommit=$GIT_COMMIT'"

echo "Build flags: $LDFLAGS"

# Function to check if cross-compilation environment is set up
check_cross_compile_env() {
    local arch=$1
    local cc_var=$2
    
    if command -v $cc_var > /dev/null 2>&1; then
        echo "✓ Cross-compiler found: $cc_var"
        return 0
    else
        echo "✗ Cross-compiler not found: $cc_var"
        return 1
    fi
}

# Build for linux/amd64
echo "Building for linux/amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags="cgo" -ldflags="-s -w" -o $ROOT_DIR/dist/osk-lin-amd64 $ROOT_DIR/cmd/oskway
echo "✓ AMD64 build completed"

# Build for linux/arm64
echo "Building for linux/arm64..."
if check_cross_compile_env "arm64" "aarch64-linux-gnu-gcc"; then
    echo "Building ARM64 with CGO enabled..."
    # Set up environment for ARM64 cross-compilation
    export CC=aarch64-linux-gnu-gcc
    export CXX=aarch64-linux-gnu-g++
    export CGO_ENABLED=1
    export GOOS=linux
    export GOARCH=arm64
    
    # Set PKG_CONFIG_PATH for cross-compilation if needed
    # export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig
    
    go build -tags="cgo" -ldflags="-s -w" -o $ROOT_DIR/dist/osk-lin-arm64 $ROOT_DIR/cmd/oskway || {
        echo "ARM64 CGO build failed, trying without CGO..."
        CGO_ENABLED=0 go build -tags="!cgo" -ldflags="-s -w" -o $ROOT_DIR/dist/osk-lin-arm64 $ROOT_DIR/cmd/oskway
        echo "⚠ ARM64 built without CGO (Wayland functionality may be limited)"
    }
else
    echo "Cross-compiler not available. Building ARM64 without CGO..."
    echo "For full functionality, install: sudo apt-get install gcc-aarch64-linux-gnu"
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags="!cgo" -ldflags="-s -w" -o $ROOT_DIR/dist/osk-lin-arm64 $ROOT_DIR/cmd/oskway
    echo "⚠ ARM64 built without CGO (Wayland functionality may be limited)"
fi
echo "✓ ARM64 build completed"

# Creating tarballs with assets
echo "Packaging binaries with assets..."
mkdir -p $ROOT_DIR/dist/tarballs

tar -czf $ROOT_DIR/dist/tarballs/osk-lin-amd64.tar.gz -C $ROOT_DIR/dist osk-lin-amd64 -C $ROOT_DIR/assets .
tar -czf $ROOT_DIR/dist/tarballs/osk-lin-arm64.tar.gz -C $ROOT_DIR/dist osk-lin-arm64 -C $ROOT_DIR/assets .

# Generate checksums
echo "Generating checksums..."
sha256sum $ROOT_DIR/dist/tarballs/osk-lin-amd64.tar.gz > $ROOT_DIR/dist/tarballs/osk-lin-amd64.sha256
sha256sum $ROOT_DIR/dist/tarballs/osk-lin-arm64.tar.gz > $ROOT_DIR/dist/tarballs/osk-lin-arm64.sha256

# Note: To hook this script to a GitHub Release action, you can call this script
# from your GitHub Actions workflow file like '.github/workflows/release.yml'
echo "Build and packaging complete."

