# wand search

Search for available formulas in the repository.

## Syntax

```bash
wand search [QUERY]
```

## Description

Searches available formulas by name or description. Shows matching packages with information.

## Usage

### Search by name

```bash
wand search nano
```

### Search by description

```bash
wand search "text editor"
```

### List all available formulas

```bash
wand search
```

## Flags

- `--limit int` - Max results to show (default: 20)
- `--format string` - Output format: `table`, `json` (default: `table`)
- `--verbose` - Show full descriptions

## Output

```
NAME             VERSION    DESCRIPTION
nano             8.7.0      Nano's ANOther editor
make             4.4        GNU Make build tool
zsh              5.9        Zsh shell
```

## Examples

### Search for text editor

```bash
$ wand search editor
nano             8.7.0      Nano's ANOther editor
vim              9.0        Vi improved
```

### Search specific package

```bash
$ wand search make
make             4.4        GNU Make build tool
```

### Show all available packages

```bash
$ wand search
```

### Get more details

```bash
$ wand search nano --verbose
nano 8.7.0
  Description: Nano's ANOther editor, an enhanced free Pico clone
  Repository: https://github.com/ochairo/formulas
  Categories: editors, tools
```

### Export results

```bash
$ wand search --format json | jq '.'
```

## See Also

- [info](./info.md) - Show detailed package info
- [install](./install.md) - Install a package
