# Production Deployment

## System Requirements

- macOS 11+ or Linux (Ubuntu 20.04+, Fedora 33+, Debian 11+)
- x86_64 or ARM64
- 500MB+ disk space
- Bash 4.0+ or Zsh 5.0+
- `git`, `curl`, `openssl`

## Installation

Standard installation works for production:

```bash
curl -fsSL https://raw.githubusercontent.com/ochairo/wand/main/install.sh | sh
```

Then add to PATH:
```bash
export PATH="$HOME/.wand/shims:$PATH"
```

## Verify Installation

```bash
wand version
wand doctor          # Full system check
wand help            # Show all commands
```

## Basic Security

```bash
# Set proper permissions
chmod 700 ~/.wand
chmod 755 ~/.wand/shims
chmod 755 ~/.wand/packages

# Enable checksum verification (default)
export WAND_VERIFY_CHECKSUMS=true
```

## Configuration

Create `.wandrc` in your project:

```yaml
packages:
  - name: nano
    version: "8.7"
  - name: make
    version: "4.4"
```

Install all:
```bash
wand install
```

## Health Checks

```bash
# Check system health
wand doctor

# List installed packages
wand list

# Check for updates
wand outdated

# Update all packages
wand update
```

## Uninstall

```bash
rm /usr/local/bin/wand
rm -rf ~/.wand
```

For details on error handling, see [ERROR_CODES.md](ERROR_CODES.md).
For development setup, see [DEVELOPMENT.md](DEVELOPMENT.md).
