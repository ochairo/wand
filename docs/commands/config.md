# wand config

Manage Wand configuration.

## Syntax

```bash
wand config SUBCOMMAND [ARGS]
```

## Description

Manages Wand configuration settings. Can view, set, reset, and validate configuration.

## Subcommands

### get

Get a configuration value:

```bash
wand config get SETTING
```

### set

Set a configuration value:

```bash
wand config set SETTING VALUE
```

### list

List all configuration values:

```bash
wand config list
```

### reset

Reset to default values:

```bash
wand config reset [SETTING]
```

## Usage

### Show all settings

```bash
wand config list
```

### Get specific setting

```bash
wand config get cache_ttl
```

### Change setting

```bash
wand config set cache_ttl 3600
```

### Reset to defaults

```bash
wand config reset
```

## Configuration Options

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `cache_enabled` | bool | true | Enable package caching |
| `cache_ttl` | int | 3600 | Cache time-to-live (seconds) |
| `cache_max_size` | string | 1GB | Maximum cache size |
| `home_dir` | string | ~/.wand | Wand home directory |
| `github_token` | string | - | GitHub API token (optional) |
| `formula_repo` | string | official | Formula repository source |
| `log_level` | string | info | Logging level |
| `verify_checksums` | bool | true | Verify package checksums |

## Examples

### View configuration

```bash
$ wand config list
cache_enabled: true
cache_ttl: 3600
log_level: info
```

### Disable caching

```bash
$ wand config set cache_enabled false
✓ Configuration updated
```

### Increase cache size

```bash
$ wand config set cache_max_size 5GB
```

### Set GitHub token for higher API limits

```bash
$ wand config set github_token ghp_xxx...
```

### Reset single setting

```bash
$ wand config reset cache_ttl
```

## Files

Configuration can be set in multiple places:

1. `~/.wand/config.yml` - User configuration
2. `.wandrc` - Project configuration (overrides user)
3. Environment variables (override all)

## See Also

- [INSTALLATION.md](../INSTALLATION.md) - Initial setup
- [GETTING_STARTED.md](../GETTING_STARTED.md) - Quick start
