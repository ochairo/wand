# wand outdated

Show packages with available updates.

## Syntax

```bash
wand outdated [PACKAGE]
```

## Description

Shows which installed packages have newer versions available for updating.

## Usage

### Check all packages

```bash
wand outdated
```

### Check specific package

```bash
wand outdated nano
```

## Flags

- `--format string` - Output format: `table`, `json` (default: `table`)
- `--prerelease` - Include pre-release versions

## Output Format

```
PACKAGE     CURRENT    LATEST     UPDATE AVAILABLE
nano        8.6.0      8.7.0      Yes
make        4.3        4.4        Yes
zsh         5.9        5.9        No
```

## Examples

### Show all outdated packages

```bash
$ wand outdated
nano        8.6.0 → 8.7.0
make        4.3 → 4.4
```

### Check specific package

```bash
$ wand outdated nano
nano is outdated:
  Current: 8.6.0
  Latest: 8.7.0

  Run 'wand update nano' to update
```

### Export as JSON

```bash
$ wand outdated --format json
```

## Next Steps

To update a package, use:

```bash
wand update nano
```

## See Also

- [update](./update.md) - Update packages
- [list](./list.md) - Show installed packages
