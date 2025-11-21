# wand uninstall

Remove an installed package.

## Syntax

```bash
wand uninstall PACKAGE [VERSION]
```

## Description

Uninstalls a specific package version or all versions of a package.

## Usage

### Uninstall all versions

```bash
wand uninstall nano
```

### Uninstall specific version

```bash
wand uninstall nano@8.1
```

## Flags

- `--force` - Skip confirmation prompt
- `--keep-config` - Keep configuration files
- `--verbose` - Show detailed removal process

## Examples

### Remove all versions

```bash
$ wand uninstall nano
Uninstall nano and all its versions? (y/n) y
✓ Uninstalling nano@8.7.0...
✓ Removing shims...
✓ Uninstallation complete
```

### Remove specific version

```bash
$ wand uninstall nano@8.1
✓ Uninstalling nano@8.1...
```

### Force uninstall without confirmation

```bash
$ wand uninstall nano --force
✓ Uninstalling nano@8.7.0...
```

## What Gets Removed

- Package binaries and files
- Command shims
- Registry entries
- Temporary files

Configuration files can be kept with `--keep-config`.

## Error Handling

- `PACKAGE_NOT_INSTALLED` - Package not found. Check with `wand list`
- `PERMISSION_DENIED` - Cannot remove files. May need `sudo`

## See Also

- [install](./install.md) - Install a package
- [list](./list.md) - Show installed packages
