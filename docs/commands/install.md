# wand install

Install a package or packages from wandfile.

## Syntax

```bash
wand install [PACKAGE]... [VERSION]
```

## Description

Installs one or more packages. Can install from command-line arguments or from a `wandfile.yaml`.

## Usage

### Install latest version

```bash
wand install nano
```

### Install specific version

```bash
wand install nano@8.2
```

### Install multiple packages

```bash
wand install nano make zsh
```

### Install from wandfile

```bash
wand install
```

Will read packages from `.wandrc` or `wandfile.yaml` in current directory.

### Install with verbose output

```bash
wand install nano --verbose
```

## Flags

- `--force` - Force reinstall even if already installed
- `--pre` - Include pre-release versions
- `--verbose` - Show detailed installation progress
- `--dry-run` - Show what would be installed without doing it

## Examples

### Basic installation

```bash
$ wand install nano
✓ Installing nano@8.7.0...
✓ Verifying checksum...
✓ Creating shims...
✓ Installation complete
```

### Install specific version

```bash
$ wand install nano@8.2
✓ Installing nano@8.2...
```

### Install from wandfile

```bash
$ cat wandfile.yaml
cli:
  - name: nano
    version: 8.7
  - name: make
    version: 4.4

$ wand install
✓ Installing nano@8.7...
✓ Installing make@4.4...
```

### Preview before installing

```bash
$ wand install nano --dry-run
Would install: nano@8.7.0
```

## Error Handling

Common errors and solutions:

- `PACKAGE_NOT_FOUND` - Package not in formula repository. Check spelling or run `wand search`
- `DOWNLOAD_FAILED` - Network issue. Check connectivity and retry
- `CHECKSUM_MISMATCH` - Download corrupted. Will retry automatically
- `DISK_SPACE_LOW` - Not enough space. Run `wand clean` or free up space

## See Also

- [uninstall](./uninstall.md) - Remove a package
- [update](./update.md) - Update to latest version
- [info](./info.md) - Show package details
