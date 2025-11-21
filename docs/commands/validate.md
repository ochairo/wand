# wand validate

Validate wandfile or formula YAML.

## Syntax

```bash
wand validate [FILE]
```

## Description

Validates YAML files for correct syntax and required fields. Supports both wandfiles and formula definitions.

## Usage

### Validate wandfile in current directory

```bash
wand validate
```

### Validate specific file

```bash
wand validate my-formula.yaml
```

### Verbose validation

```bash
wand validate --verbose
```

## Validations

For wandfiles:
- Valid YAML syntax
- Required fields present
- Valid package names
- Valid version specifications

For formulas:
- Required formula fields
- Valid platform configurations
- URL format validation
- Version format validation

## Output

### Valid file

```bash
$ wand validate
✓ wandfile.yaml is valid
```

### Invalid file

```bash
$ wand validate
✗ wandfile.yaml has errors:
  - Missing required field: 'name' (line 5)
  - Invalid package version: '1.0.invalid' (line 12)
  - Duplicate package: 'nano' (line 15)
```

## Examples

### Validate project wandfile

```bash
$ cat wandfile.yaml
cli:
  - name: nano
    version: 8.7

$ wand validate
✓ Valid wandfile
```

### Validate formula

```bash
$ wand validate formulas/nano.yaml
✓ Formula is valid and ready to use
```

### Check for issues

```bash
$ wand validate --verbose
Validating wandfile.yaml...
✓ YAML syntax: OK
✓ Required fields: OK
✓ Package names: OK (2 packages)
✓ Versions: OK
```

## Wandfile Format

```yaml
name: My Project

cli:
  - name: nano
    version: "8.7"  # Required
    pin: false      # Optional

gui:
  - microsoft-edge  # Name only

dotfiles:          # Optional
  repo: https://github.com/username/dotfiles
  symlinks:
    .bashrc: bash/bashrc
```

## Error Messages

Common validation errors:

- `Invalid package name` - Package name contains invalid characters
- `Invalid version format` - Version doesn't follow semver
- `Duplicate package` - Same package listed multiple times
- `Missing required field` - Required YAML field is missing

## See Also

- [FORMULA_GUIDE.md](../FORMULA_GUIDE.md) - Creating formulas
- [COMMAND_REFERENCE.md](../COMMAND_REFERENCE.md) - Command help
