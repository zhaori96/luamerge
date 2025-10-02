# Installation Examples

This document provides various examples of how users can install `luamerge` using the installation script.

## 📥 Basic Installation

### Install Latest Version

The simplest way to install the latest version:

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```

**What it does:**
- Downloads the latest release from GitHub
- Installs to `~/.local/bin/luamerge`
- Makes the binary executable

### Install Specific Version

Install a specific version instead of the latest:

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | VERSION=v1.0.0 bash
```

**Examples:**
```bash
# Install version 1.0.0
VERSION=v1.0.0 curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash

# Install version 2.3.1
VERSION=v2.3.1 curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash

# Install a pre-release version
VERSION=v2.0.0-beta curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```

## 📂 Custom Installation Directory

### Install to System-Wide Location

Install to `/usr/local/bin` (requires sudo):

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | sudo INSTALL_DIR=/usr/local/bin bash
```

### Install to User's Local Bin

Default location (recommended):

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | INSTALL_DIR=$HOME/.local/bin bash
```

### Install to Custom Path

```bash
# Install to ~/bin
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | INSTALL_DIR=$HOME/bin bash

# Install to a project-specific directory
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | INSTALL_DIR=/opt/myproject/bin bash
```

## 🔄 Combining Options

### Specific Version + Custom Directory

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | \
  VERSION=v1.2.0 INSTALL_DIR=/usr/local/bin bash
```

### Download Script First, Then Install

Useful for reviewing the script before execution:

```bash
# Download the script
wget https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh

# Review it
cat install.sh

# Run with default options
chmod +x install.sh
./install.sh

# Or with custom options
VERSION=v1.0.0 ./install.sh

# Or with multiple options
VERSION=v1.0.0 INSTALL_DIR=$HOME/bin ./install.sh
```

## 🐧 Platform-Specific Examples

### Linux (Ubuntu/Debian)

```bash
# Install latest to user directory
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash

# Install to system-wide location
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | \
  sudo INSTALL_DIR=/usr/local/bin bash
```

### macOS

```bash
# Install latest to user directory
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash

# Install using Homebrew location
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | \
  INSTALL_DIR=/usr/local/bin bash
```

### Windows (Git Bash / WSL)

```bash
# In Git Bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash

# In WSL (Ubuntu)
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```

## 🔧 Advanced Usage

### Install Multiple Versions Side-by-Side

```bash
# Install v1.0.0 as luamerge-1.0.0
VERSION=v1.0.0 INSTALL_DIR=$HOME/.local/bin bash -c '
  curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
  mv $HOME/.local/bin/luamerge $HOME/.local/bin/luamerge-1.0.0
'

# Install v2.0.0 as luamerge-2.0.0
VERSION=v2.0.0 INSTALL_DIR=$HOME/.local/bin bash -c '
  curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
  mv $HOME/.local/bin/luamerge $HOME/.local/bin/luamerge-2.0.0
'

# Install latest as luamerge (default)
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```

### Automated Installation in CI/CD

```yaml
# GitHub Actions example
- name: Install luamerge
  run: |
    curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | \
      VERSION=v1.0.0 INSTALL_DIR=/usr/local/bin bash
    luamerge --help
```

```yaml
# GitLab CI example
install_luamerge:
  script:
    - curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
    - export PATH="$HOME/.local/bin:$PATH"
    - luamerge --help
```

### Docker Installation

```dockerfile
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl && \
    curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | \
    INSTALL_DIR=/usr/local/bin bash && \
    apt-get clean

CMD ["luamerge", "--help"]
```

## 🔍 Verification

After installation, verify it worked:

```bash
# Check installation
which luamerge

# Check version
luamerge --version

# Get help
luamerge --help
```

## 🔄 Upgrading

### Upgrade to Latest Version

Simply run the install script again:

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```

### Upgrade to Specific Version

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | VERSION=v2.0.0 bash
```

### Downgrade to Older Version

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | VERSION=v1.0.0 bash
```

## ❌ Uninstallation

To uninstall, simply remove the binary:

```bash
# If installed to ~/.local/bin (default)
rm ~/.local/bin/luamerge

# If installed to /usr/local/bin
sudo rm /usr/local/bin/luamerge

# If installed to custom location
rm /custom/path/luamerge
```

## 🆘 Troubleshooting

### "Command not found" after installation

The install directory is not in your PATH. Add it:

```bash
# For ~/.local/bin (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### "Version not found"

Check available versions:

```bash
# List all releases
curl -s https://api.github.com/repos/YOUR_USERNAME/luamerge/releases | grep tag_name
```

### Permission denied

Use sudo for system-wide directories:

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | \
  sudo INSTALL_DIR=/usr/local/bin bash
```

### Download failed

Check your internet connection and GitHub status:

```bash
# Test GitHub connectivity
curl -I https://github.com

# Try with verbose output
curl -v -sSL https://raw.githubusercontent.com/YOUR_USERNAME/luamerge/main/install.sh | bash
```
