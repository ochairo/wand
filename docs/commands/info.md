# wand info

Show detailed information about a package.

## Syntax

```bash
wand info PACKAGE [VERSION]
```

## Description

Displays comprehensive information about a formula including versions, platforms, dependencies, and installation details.

## Usage

### Show package info

```bash
wand info nano
```

### Show info for specific version

```bash
wand info nano@8.7
```

### Export as JSON

```bash
wand info nano --format json
```

## Flags

- `--format string` - Output format: `table`, `json` (default: `table`)
- `--checksums` - Show download checksums
- `--platforms` - Show supported platforms
- `--dependencies` - Show package dependencies

## Output

```
Package: nano
Latest Version: 8.7.0
Description: Nano's ANOther editor, an enhanced free Pico clone

Available Versions:
  8.7.0 (latest)
  8.6.0
  8.5.0

Supported Platforms:
  macOS (x86_64, arm64)
  Linux (x86_64, arm64)

Homepage: https://www.nano-editor.org
Repository: https://github.com/ochairo/formulas
License: GPLv3
```

## Examples

### Get package details

```bash
$ wand info nano
Package: nano
Version: 8.7.0
Description: Nano editor
Platforms: macOS, Linux
License: GPLv3
```

### Show checksums

```bash
$ wand info nano --checksums
8.7.0 darwin-amd64:
  sha256: abc123def456...
8.7.0 linux-amd64:
  sha256: xyz789uvw012...
```

### Export as JSON

```bash
$ wand info nano --format json
```

## See Also

- [search](./search.md) - Find packages
- [install](./install.md) - Install a package
