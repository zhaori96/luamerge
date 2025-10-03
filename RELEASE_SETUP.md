# Release Setup Guide

This guide explains how to set up automated releases for luamerge on GitHub.

## Prerequisites

1. A GitHub repository for this project
2. Git installed and configured
3. GitHub account with access to create releases

## Setup Steps

### 1. Update Repository References

Before pushing, replace `YOUR_USERNAME` in the following files with your actual GitHub username:

- `install.sh` (line 11)
- `.goreleaser.yml` (line 58-59)
- `README.md` (lines 21, 27, 34, 39)

Example:
```bash
# Replace YOUR_USERNAME with your actual GitHub username
sed -i 's/YOUR_USERNAME/yourusername/g' install.sh
sed -i 's/YOUR_USERNAME/yourusername/g' .goreleaser.yml
sed -i 's/YOUR_USERNAME/yourusername/g' README.md
```

### 2. Push to GitHub

```bash
git add .
git commit -m "Add release automation and installation script"
git push origin main
```

### 3. Create a Release

To create a new release, push a tag:

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will automatically:
- Build binaries for Linux, macOS, and Windows
- Support for amd64, arm64, and arm architectures
- Create a GitHub release with all binaries
- Generate checksums
- Create archives (tar.gz for Linux/macOS, zip for Windows)

### 4. Verify Release

After pushing the tag:
1. Go to https://github.com/YOUR_USERNAME/luamerge/actions
2. Wait for the "Release" workflow to complete
3. Check https://github.com/YOUR_USERNAME/luamerge/releases
4. Your new release should be available with all binaries

## Installation Commands

Once you have a release, users can install using:

### Install Latest Version
```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```

### Install Specific Version
```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | VERSION=v1.0.0 bash
```

### Install to Custom Directory
```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | INSTALL_DIR=/usr/local/bin bash
```

### Combine Options
```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | VERSION=v1.2.0 INSTALL_DIR=/custom/path bash
```

## Release Versioning

Follow semantic versioning:
- `v1.0.0` - Major release
- `v1.1.0` - Minor release (new features)
- `v1.0.1` - Patch release (bug fixes)

## Manual Testing of Install Script

Before the first release, you can test the install script locally:

```bash
# Build a test binary
go build -o luamerge ./cmd/cli

# Test manual installation
mkdir -p ~/.local/bin
cp luamerge ~/.local/bin/
chmod +x ~/.local/bin/luamerge

# Verify installation
luamerge --help
```

## Files Created

- `install.sh` - Installation script for users
- `.goreleaser.yml` - GoReleaser configuration
- `.github/workflows/release.yml` - GitHub Actions workflow
- `RELEASE_SETUP.md` - This guide (can be deleted after setup)

## Troubleshooting

### GitHub Actions fails
- Check that you have enabled Actions in your repository settings
- Verify that the GITHUB_TOKEN has write permissions

### Install script fails
- Ensure you've replaced YOUR_USERNAME in all files
- Verify that a release exists on GitHub
- Check that binaries were uploaded correctly

### Binary naming issues
If the install script can't find binaries, check that GoReleaser is creating them with the expected names:
- Linux: `luamerge_linux_amd64`
- macOS: `luamerge_darwin_amd64`
- Windows: `luamerge_windows_amd64.exe`

## Notes

- The install script installs to `~/.local/bin` by default
- Users can override with: `INSTALL_DIR=/custom/path ./install.sh`
- The script automatically detects OS and architecture
- Supports Linux, macOS, and Windows (via Git Bash/WSL)
