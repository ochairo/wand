# wand version

Show Wand version and build information.

## Syntax

```bash
wand version
```

## Description

Displays the currently installed Wand version, build information, and compilation details.

## Usage

### Show version

```bash
wand version
```

## Output Format

```
Wand v1.0.0
Build: darwin-arm64
Commit: abc1234
Built: 2025-11-20T10:30:00Z
```

## Examples

### Check version

```bash
$ wand version
Wand v1.0.0
```

### Full version info

```bash
$ wand version
Wand v1.0.0
Build Information:
  Platform: darwin-arm64
  Commit: abc1234567890def
  Built: 2025-11-20 at 10:30 UTC
  Go Version: go1.21
```

## Using in Scripts

```bash
#!/bin/bash

# Get version
version=$(wand version | head -1 | cut -d' ' -f2)

if [[ "$version" > "1.0.0" ]]; then
    echo "Wand is up to date"
else
    echo "Wand update available"
fi
```

## Checking for Updates

```bash
# Check for newer version
wand check-updates

# Update if available
wand update-self
```

## See Also

- [RELEASE_NOTES_v1.0.0.md](../RELEASE_NOTES_v1.0.0.md) - Release notes
- [INSTALLATION.md](../INSTALLATION.md) - Installation guide
