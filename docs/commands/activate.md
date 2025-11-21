# wand activate

Switch to a different version of an installed package.

## Syntax

```bash
wand activate PACKAGE VERSION
```

## Description

Sets a specific version of a package as the active version. This changes which binary is executed when you run the command.

## Usage

### Activate different version

```bash
wand activate nano@8.1
```

### Verify active version

```bash
nano --version
```

## Examples

### Switch to different version

```bash
$ wand activate nano@8.1
✓ Activated nano@8.1

$ nano --version
  GNU nano, version 8.1
```

### Switch back to latest

```bash
$ wand activate nano@8.7
$ nano --version
  GNU nano, version 8.7.0
```

## Requirements

- Package must be installed with multiple versions
- Target version must exist

## Error Handling

- `PACKAGE_NOT_INSTALLED` - Package not found
- `VERSION_NOT_FOUND` - Version not installed

## See Also

- [install](./install.md) - Install a package version
- [list](./list.md) - Show installed versions
