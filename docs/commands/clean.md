# wand clean

Clean up and remove unused files.

## Syntax

```bash
wand clean [OPTIONS]
```

## Description

Cleans up temporary files, old versions, and cache to free up disk space.

## Flags

- `--all` - Remove all temporary files and old versions
- `--old-versions` - Remove non-active package versions
- `--cache` - Clear package download cache
- `--tmp` - Remove temporary extraction files
- `--dry-run` - Show what would be removed
- `--force` - Skip confirmation

## Usage

### Basic cleanup

```bash
wand clean
```

### Remove old versions

```bash
wand clean --old-versions
```

### Clear cache

```bash
wand clean --cache
```

### Preview cleanup

```bash
wand clean --dry-run
```

## Examples

### Clean up everything

```bash
$ wand clean --all
Will remove:
  - 234 MB of cache
  - 156 MB of old versions
  - 89 MB of temporary files

Total: 479 MB to be freed

Continue? (y/n) y
✓ Cleanup complete. Freed 479 MB
```

### Show what would be cleaned

```bash
$ wand clean --all --dry-run
Would remove:
  - nano@8.6.0 (156 MB)
  - make@4.3 (78 MB)
  - cache/*.tar.gz (234 MB)
  - tmp/* (89 MB)
```

### Remove old versions only

```bash
$ wand clean --old-versions
✓ Removed old versions (156 MB freed)
```

## What Gets Cleaned

- Old package versions (keeping active version)
- Downloaded cache files
- Temporary extraction files
- Build artifacts

## Safety

Cleanup is safe:
- Only removes files we know about
- Keeps active versions
- Asks for confirmation
- Use `--dry-run` to preview first

## See Also

- [cache](./cache.md) - Manage cache specifically
- [uninstall](./uninstall.md) - Remove packages
