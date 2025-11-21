# wand list

List all installed packages and their versions.

## Syntax

```bash
wand list [PACKAGE]
```

## Description

Shows all installed packages and their versions. Can show details for specific package.

## Usage

### List all packages

```bash
wand list
```

### Show specific package versions

```bash
wand list nano
```

### List in JSON format

```bash
wand list --format json
```

## Flags

- `--format string` - Output format: `table`, `json`, `csv` (default: `table`)
- `--sort string` - Sort by: `name`, `version`, `date` (default: `name`)
- `--verbose` - Show additional details

## Output Formats

### Table (default)

```
PACKAGE     VERSION    ACTIVE    INSTALLED
nano        8.7.0      ✓         2025-01-15
make        4.4        ✓         2025-01-10
zsh         5.9                  2025-01-05
```

### JSON

```json
{
  "packages": [
    {
      "name": "nano",
      "versions": ["8.7.0", "8.6.0"],
      "active": "8.7.0",
      "installed": "2025-01-15T10:30:00Z"
    }
  ]
}
```

## Examples

### Simple list

```bash
$ wand list
nano      8.7.0
make      4.4
zsh       5.9
```

### Show specific package

```bash
$ wand list nano
nano:
  8.7.0 (active)
  8.6.0
  8.5.0
```

### Export as JSON

```bash
$ wand list --format json > packages.json
```

### Sort by installation date

```bash
$ wand list --sort date
```

## See Also

- [search](./search.md) - Find available packages
- [info](./info.md) - Show package details
- [outdated](./outdated.md) - Show packages with updates
