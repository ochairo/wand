# Wandfile Examples

Wandfiles define system-wide package installations and configurations declaratively.

## Overview

A Wandfile allows you to:
- Declare all packages and versions for your system
- Share configurations across machines
- Reproducibly set up development environments
- Version control your tool stack

## Structure

```
wandfile/
├── README.md
├── basic/
│   └── Wandfile              # Simple package list
├── development/
│   └── Wandfile              # Full dev environment
├── team/
│   └── Wandfile              # Team-wide standards
└── with-dotfiles/
    └── Wandfile              # Packages + dotfiles together
```

## Usage

```bash
# Install all packages from a Wandfile
wand install --wandfile ./Wandfile

# Validate a Wandfile
wand validate ./Wandfile

# Show what would be installed
wand install --wandfile ./Wandfile --dry-run
```

## Examples

### Basic Wandfile
See [basic/](./basic/) for a minimal example with essential CLI tools.

### Development Environment
See [development/](./development/) for a complete development setup.

### Team Configuration
See [team/](./team/) for standardized team tooling.

### With Dotfiles
See [with-dotfiles/](./with-dotfiles/) for combining packages and dotfiles in one Wandfile.

## Wandfile Format

```yaml
version: "1.0"

packages:
  - name: jq
    version: "1.7.1"
  - name: ripgrep
    version: latest

dotfiles:
  - name: my-dotfiles
    repository: https://github.com/user/dotfiles.git
    configs:
      - nvim
      - zsh
```

## Related Documentation

- [Wandfile Specification](../../docs/WANDFILE.md)
- [Getting Started Guide](../../docs/GETTING_STARTED.md)
