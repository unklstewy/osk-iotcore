# OSK Release Process

This document describes the release process for the OSK (On-Screen Keyboard) project.

## Overview

The release process is automated using GitHub Actions and includes:
- Cross-compilation for linux/amd64 and linux/arm64
- CGO support for Wayland functionality
- Fallback builds without CGO for limited environments
- Asset packaging with layouts and themes
- Checksum generation for integrity verification

## Files

### `scripts/build_release.sh`
The main build script that handles cross-compilation and packaging.

**Features:**
- Cross-compiles for linux/amd64 and linux/arm64
- Uses CGO when available (with proper cross-compiler setup)
- Falls back to non-CGO builds for limited environments
- Packages binaries with assets (layouts, themes, keycaps)
- Generates SHA256 checksums
- Includes build metadata (version, git commit, build date)

**Usage:**
```bash
# Basic usage
./scripts/build_release.sh

# With custom version
VERSION=v1.0.0 ./scripts/build_release.sh
```

### `.github/workflows/release.yml`
GitHub Actions workflow that triggers on release creation.

**Features:**
- Runs on Ubuntu latest
- Sets up Go 1.23.0
- Installs cross-compilation dependencies
- Builds and packages releases
- Uploads assets to GitHub Release

## Build Targets

### linux/amd64
- **CGO**: Enabled (requires wayland-client libraries)
- **Features**: Full Wayland support
- **Binary**: Dynamically linked
- **Package**: `osk-lin-amd64.tar.gz`

### linux/arm64
- **CGO**: Enabled when cross-compiler available, otherwise disabled
- **Features**: Full Wayland support (with CGO) or limited (without CGO)
- **Binary**: Dynamically linked (with CGO) or statically linked (without CGO)
- **Package**: `osk-lin-arm64.tar.gz`

## Cross-Compilation Requirements

### For Full CGO Support on ARM64:
```bash
# Install cross-compiler
sudo apt-get install gcc-aarch64-linux-gnu

# Install ARM64 development libraries (if needed)
sudo apt-get install libwayland-dev:arm64
```

### Build Tags
- `cgo`: Enables CGO build with full Wayland support
- `!cgo`: Disables CGO, provides limited functionality but better portability

## Package Contents

Each release package (`*.tar.gz`) contains:
- Binary executable (`osk-lin-amd64` or `osk-lin-arm64`)
- Asset files:
  - `layouts/` - Keyboard layout definitions
  - `themes/` - UI theme configurations
  - `keycaps/` - Keycap assets (if any)

## Checksums

Each package includes a SHA256 checksum file:
- `osk-lin-amd64.sha256`
- `osk-lin-arm64.sha256`

**Verification:**
```bash
# Verify package integrity
sha256sum -c osk-lin-amd64.sha256
```

## Release Process

### Manual Release
1. Create and push a git tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. Create a GitHub Release from the tag

3. The GitHub Action will automatically:
   - Build binaries for both architectures
   - Package with assets
   - Generate checksums
   - Upload to the release

### Local Testing
```bash
# Test the build script locally
./scripts/build_release.sh

# Check outputs
ls -la dist/
ls -la dist/tarballs/

# Test binaries
./dist/osk-lin-amd64 --version
file dist/osk-lin-*
```

## Troubleshooting

### CGO Build Fails
- Ensure wayland development libraries are installed
- Check cross-compiler availability
- Review build logs for missing dependencies

### Package Size Issues
- Binaries are stripped (`-s -w` flags)
- Assets are compressed in tarballs
- Consider excluding unnecessary files from assets

### GitHub Actions Fails
- Check workflow logs
- Verify Go version compatibility
- Ensure all dependencies are available in CI environment

## Dependencies

### Build Dependencies
- Go 1.23.0 or later
- gcc (for amd64 CGO builds)
- gcc-aarch64-linux-gnu (for arm64 CGO builds)
- pkg-config
- wayland-client development libraries

### Runtime Dependencies
- Wayland compositor (for CGO builds)
- glibc (for dynamically linked builds)

## Version Information

The build script embeds version information into binaries:
- `Version`: Release version (from VERSION env var or default)
- `BuildDate`: UTC build timestamp
- `GitCommit`: Git commit hash
- `GitBranch`: Git branch name

Access via `--version` flag or programmatically in Go code.
