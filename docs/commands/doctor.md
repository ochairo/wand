# wand doctor

Check system health and diagnose issues.

## Syntax

```bash
wand doctor
```

## Description

Runs comprehensive health checks on Wand installation and system configuration. Identifies problems and suggests fixes.

## Usage

### Run diagnostics

```bash
wand doctor
```

### Verbose diagnostics

```bash
wand doctor --verbose
```

## Checks Performed

- ✓ Wand installation integrity
- ✓ ~/.wand directory structure
- ✓ ~/.wand/shims in PATH
- ✓ Network connectivity to GitHub
- ✓ Available disk space
- ✓ Installed packages integrity
- ✓ Shell completion configuration
- ✓ Permission issues

## Output Example

```
Wand Health Check

✓ Wand version: 1.0.0
✓ Home directory: /Users/user/.wand
✓ Shims directory: /Users/user/.wand/shims
✓ Shims in PATH: Yes
✓ Network: Connected (GitHub API accessible)
✓ Disk space: 2.5 GB available
✓ Installed packages: 5 packages, all valid
✓ Permissions: All correct

System is healthy!
```

## Examples

### Basic check

```bash
$ wand doctor
Wand Health Check
✓ System is healthy
```

### Detailed diagnostics

```bash
$ wand doctor --verbose
[Detailed output with all checks and their results]
```

### When issues are found

```bash
$ wand doctor
✗ Shims not in PATH
  Add to your shell profile:
  export PATH="$HOME/.wand/shims:$PATH"

✗ Low disk space (150MB available)
  Run: wand clean
```

## Troubleshooting

If issues are found, follow the suggestions printed by `wand doctor`. Common fixes:

```bash
# Fix PATH issues
echo 'export PATH="$HOME/.wand/shims:$PATH"' >> ~/.bashrc

# Fix permissions
chmod 700 ~/.wand

# Clean up and fix registry
wand clean
wand verify
```

## See Also

- [ERROR_CODES.md](../ERROR_CODES.md) - Error reference
- [TROUBLESHOOTING.md](../TROUBLESHOOTING.md) - Common issues
