# wand cache

Manage package cache.

## Syntax

```bash
wand cache SUBCOMMAND [ARGS]
```

## Description

Manages the local package cache. Can clean, clear, and inspect cache.

## Subcommands

### clean

Remove cached downloads for specific package:

```bash
wand cache clean [PACKAGE]
```

### clear

Clear entire cache:

```bash
wand cache clear
```

### size

Show cache size:

```bash
wand cache size
```

## Usage

### Clean cache for package

```bash
wand cache clean nano
```

### Clear all cached files

```bash
wand cache clear
```

### Check cache size

```bash
wand cache size
```

## Examples

### Remove specific package cache

```bash
$ wand cache clean nano
✓ Cleaned cache for nano (234 MB freed)
```

### Clear entire cache

```bash
$ wand cache clear
Clear entire cache? (y/n) y
✓ Cache cleared (1.2 GB freed)
```

### Check cache usage

```bash
$ wand cache size
Cache size: 1.2 GB
Max size: 5 GB (24% used)

Top packages:
  nano: 456 MB
  make: 234 MB
  zsh: 156 MB
```

## When to Clean Cache

- Disk space running low
- Package downloads corrupted
- Freeing up space before large installation
- Cache seems stale or outdated

## Cache Configuration

```bash
# Disable caching
wand config set cache_enabled false

# Set max cache size
wand config set cache_max_size 2GB

# Set cache time-to-live
wand config set cache_ttl 86400  # 24 hours
```

## See Also

- [config](./config.md) - Manage configuration
- [clean](./clean.md) - General cleanup command
