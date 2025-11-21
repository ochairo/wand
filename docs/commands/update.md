# wand update

Update a package to the latest version.

## Syntax

```bash
wand update [PACKAGE]...
```

## Description

Updates packages to their latest available versions. Can update specific packages or all installed packages.

## Usage

### Update specific package

```bash
wand update nano
```

### Update multiple packages

```bash
wand update nano make zsh
```

### Update all packages

```bash
wand update
```

### Preview updates without applying

```bash
wand update --dry-run
```

## Flags

- `--dry-run` - Show what would be updated without doing it
- `--force` - Force update even if already latest
- `--verbose` - Show detailed update process
- `--skip-confirmation` - Don't ask before updating

## Examples

### Update single package

```bash
$ wand update nano
✓ nano: 8.6.0 → 8.7.0
✓ Installation complete
```

### Check what would be updated

```bash
$ wand update --dry-run
Would update:
  nano: 8.6.0 → 8.7.0
  make: 4.3 → 4.4
```

### Update all packages

```bash
$ wand update
✓ nano: 8.6.0 → 8.7.0
✓ make: 4.3 → 4.4
✓ 2 packages updated
```

## Notes

- Old versions are kept (can be managed separately)
- Active version automatically switches to latest
- Updates respect version pinning in wandfile

## See Also

- [outdated](./outdated.md) - Show available updates
- [install](./install.md) - Install specific version
- [activate](./activate.md) - Switch versions
