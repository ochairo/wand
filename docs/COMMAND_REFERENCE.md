# Wand Command Reference

A complete reference to all Wand commands. For detailed information about each command, see the individual command documentation in the `docs/commands/` directory.

## Command Index

### Installation & Management

| Command | Description |
|---------|-------------|
| [install](./commands/install.md) | Install a package or packages from wandfile |
| [uninstall](./commands/uninstall.md) | Remove an installed package |
| [update](./commands/update.md) | Update packages to their latest versions |
| [activate](./commands/activate.md) | Switch to a different version of an installed package |
| [clean](./commands/clean.md) | Clean up and remove unused files |

### Discovery & Information

| Command | Description |
|---------|-------------|
| [list](./commands/list.md) | List all installed packages and their versions |
| [search](./commands/search.md) | Search for available formulas in the repository |
| [info](./commands/info.md) | Show detailed information about a package |
| [outdated](./commands/outdated.md) | Show packages with available updates |

### System & Configuration

| Command | Description |
|---------|-------------|
| [doctor](./commands/doctor.md) | Check system health and diagnose issues |
| [config](./commands/config.md) | Manage Wand configuration |
| [cache](./commands/cache.md) | Manage package cache |
| [validate](./commands/validate.md) | Validate wandfile or formula YAML |

### Utility

| Command | Description |
|---------|-------------|
| [version](./commands/version.md) | Show Wand version and build information |
| [help](./commands/help.md) | Show help information |

## Global Flags

All commands support the following global flags:

| Flag | Description |
|------|-------------|
| `--help, -h` | Show help for command |
| `--verbose` | Enable detailed output |
| `--config string` | Path to configuration file |
| `--dry-run` | Preview action without executing |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `WAND_HOME` | Override home directory (default: ~/.wand) |
| `WAND_CONFIG` | Override config file path |
| `WAND_CACHE_DIR` | Override cache directory |
| `WAND_LOG_LEVEL` | Set logging level (debug, info, warn, error) |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Command not found |
| `3` | Invalid arguments |
| `4` | Package not found |
| `5` | Installation failed |
| `6` | Network error |

## Quick Examples

```bash
# Install a package
wand install nano

# List installed packages
wand list

# Search for a package
wand search "text editor"

# Update all packages
wand update

# Check system health
wand doctor

# Show help
wand help COMMAND
```

## Related Documentation

- [GETTING_STARTED.md](./GETTING_STARTED.md) - Quick start guide
- [ERROR_CODES.md](./ERROR_CODES.md) - Error code reference and solutions
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Production setup
