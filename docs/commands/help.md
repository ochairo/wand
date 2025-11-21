# wand help

Show help information.

## Syntax

```bash
wand help [COMMAND]
```

## Description

Displays help information for Wand or for a specific command.

## Usage

### General help

```bash
wand help
```

### Command help

```bash
wand help install
```

### List all commands

```bash
wand help commands
```

## Examples

### Get command help

```bash
$ wand help install
Usage: wand install [PACKAGE] [VERSION]

Install a package or packages from wandfile.

Examples:
  wand install nano
  wand install nano@8.7
  wand install
```

### List available commands

```bash
$ wand help commands
Available Commands:
  activate      Switch to a different package version
  cache         Manage package cache
  config        Manage Wand configuration
  doctor        Check system health
  help          Show help information
  info          Show package information
  install       Install a package
  list          List installed packages
  outdated      Show packages with available updates
  search        Search for packages
  uninstall     Remove a package
  update        Update packages
  validate      Validate wandfile or formula
  version       Show Wand version
```

### Show general help

```bash
$ wand help
Wand - Universal Package Manager

Usage:
  wand [command]

Available Commands:
  See 'wand help commands' for full list

Global Flags:
  --help, -h            Show help
  --version, -v         Show version
  --verbose             Enable verbose output
  --config string       Config file path

Examples:
  wand install nano
  wand list
  wand doctor

For more help:
  wand help [command]  - Show command help
  wand doctor          - Diagnose problems
```

## Global Flags

All commands support:

- `--help, -h` - Show help for command
- `--verbose` - Enable detailed output
- `--config string` - Path to config file

## Tips

- Use `wand help COMMAND` for detailed help on any command
- Use `wand COMMAND --help` as shorthand
- Use `wand doctor` to diagnose problems
- Check documentation at: https://github.com/ochairo/wand/blob/main/docs/

## See Also

- [COMMAND_REFERENCE.md](../COMMAND_REFERENCE.md) - Full command reference
- [GETTING_STARTED.md](../GETTING_STARTED.md) - Quick start guide
