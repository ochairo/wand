# Wand Error Codes Reference

This document provides a comprehensive guide to all error codes in Wand, including their meanings, common causes, and recommended solutions.

## Error Code Format

All Wand errors follow a consistent format with the following structure:

```
[ERROR_CODE] Message: Details
```

Example:
```
[PACKAGE_NOT_FOUND] Package not found: package: "nano"
[DOWNLOAD_FAILED] Failed to download nano@8.7: network error
```

## Error Categories

### Input Validation Errors

#### `INVALID_PACKAGE_NAME`
**When**: Package name contains invalid characters or format

**Common Causes**:
- Package name contains special characters (except hyphens and underscores)
- Package name is empty or whitespace-only
- Package name exceeds length limits

**Solutions**:
```bash
# ✓ Valid package names
wand install nano
wand install my-package
wand install package_name

# ✗ Invalid package names
wand install "my package"  # spaces not allowed
wand install "@package"    # @ not allowed
```

#### `INVALID_VERSION`
**When**: Version string doesn't follow semantic versioning

**Common Causes**:
- Version is not valid semver (e.g., `1.0` should be `1.0.0`)
- Version contains invalid characters
- Version doesn't exist for package

**Solutions**:
```bash
# ✓ Valid versions
wand install nano@8.7.0
wand install nano@8.7    # auto-normalized to 8.7.0
wand install nano@latest

# ✗ Invalid versions
wand install nano@8.7.alpha  # prerelease not supported
wand install nano@v8.7       # prefix not allowed
```

#### `INVALID_PATH`
**When**: File path doesn't meet validation requirements

**Common Causes**:
- Path contains invalid characters for the platform
- Path is not absolute when absolute is required
- Path exceeds maximum length

**Solutions**:
- Use absolute paths: `/Users/name/.wand` instead of `~/.wand`
- Avoid special characters in paths
- Ensure path length is under 260 characters (Windows limit)

#### `INVALID_URL`
**When**: Download URL is malformed or unreachable

**Common Causes**:
- URL has invalid scheme (not http/https)
- URL contains invalid characters
- URL has been removed from server

**Solutions**:
- Check formula definition in `formulas/` directory
- Verify GitHub repository still exists
- Contact maintainers if URL is outdated

---

### Package Management Errors

#### `PACKAGE_NOT_FOUND`
**When**: Package doesn't exist or formula isn't loaded

**Common Causes**:
- Package name is misspelled
- Package hasn't been added to Wand formulas
- Formula repository unreachable

**Solutions**:
```bash
# List available packages
wand list

# Check formula directory
ls formulas/

# Install specific package
wand install nano
```

#### `PACKAGE_NOT_INSTALLED`
**When**: Trying to operate on a package that isn't installed

**Common Causes**:
- Package hasn't been installed yet
- Package was uninstalled previously
- Using wrong package name

**Solutions**:
```bash
# Check installed packages
wand list --installed

# Install package
wand install nano

# Reinstall if needed
wand install nano --force
```

#### `PACKAGE_INSTALLED`
**When**: Package version already installed

**Common Causes**:
- Attempting to install a version that's already installed
- Re-running installation without uninstalling first

**Solutions**:
```bash
# Update to latest version
wand install nano@latest --force

# Uninstall first if needed
wand uninstall nano
wand install nano@8.7
```

#### `VERSION_NOT_FOUND`
**When**: Specific version doesn't exist for package

**Common Causes**:
- Version was never released
- Version is misspelled
- Release was deleted from GitHub

**Solutions**:
```bash
# List available versions
wand versions nano

# Install latest if unsure
wand install nano@latest

# Check GitHub releases
wand info nano
```

#### `VERSION_INSTALLED`
**When**: Version is already globally selected

**Common Causes**:
- Setting version that's already active
- Updating to same version number

**Solutions**:
```bash
# Check current versions
wand status

# Switch to different version if needed
wand install nano@8.6
```

---

### Installation Errors

#### `DOWNLOAD_FAILED`
**When**: Package download from GitHub fails

**Common Causes**:
- Network connectivity issues
- GitHub API rate limit exceeded (>60 requests/hour)
- Download URL no longer valid
- Insufficient disk space

**Solutions**:
```bash
# Retry with backoff
wand install nano

# Check network
ping github.com

# Free up disk space if needed
wand clean

# Check GitHub status
# https://www.githubstatus.com/
```

#### `CHECKSUM_MISMATCH`
**When**: Downloaded package integrity check fails

**Common Causes**:
- Downloaded file was corrupted during transfer
- Release was re-published with different content
- Network interference corrupted download
- SHA256 file doesn't match binary

**Solutions**:
```bash
# Retry installation
wand install nano

# Clean cache and retry
wand clean
wand install nano

# Report issue if persistent
# https://github.com/ochairo/wand/issues
```

#### `EXTRACTION_FAILED`
**When**: Package archive cannot be extracted

**Common Causes**:
- Archive is corrupted
- Unsupported archive format
- Insufficient permissions
- Insufficient disk space

**Solutions**:
```bash
# Ensure sufficient disk space
df -h

# Retry installation (downloads fresh)
wand install nano

# Check file permissions
ls -la ~/.wand/
```

#### `INSTALLATION_FAILED`
**When**: Installation process failed at any stage

**Common Causes**:
- Build from source failed
- Permission denied creating directories
- Post-install hooks failed
- Registry update failed

**Solutions**:
```bash
# Check logs
wand install nano -v  # verbose mode

# Ensure required build tools
command -v make gcc    # check if available

# Run diagnostics
wand doctor

# Check permissions
sudo chown -R $USER ~/.wand
```

#### `BINARY_NOT_FOUND`
**When**: Expected binary not found after extraction

**Common Causes**:
- Archive has unexpected structure
- Binary name doesn't match formula definition
- Extraction removed or renamed binary
- Platform-specific binary not included

**Solutions**:
```bash
# Check formula definition
cat formulas/nano.yaml

# Inspect archive contents
tar -tzf /path/to/archive.tar.gz

# Update formula if needed
wand generate-formula nano
```

---

### Shim Errors

#### `SHIM_CREATION_FAILED`
**When**: Cannot create command wrapper (shim)

**Common Causes**:
- Permissions denied on shims directory
- Shim already exists and is protected
- File system doesn't support scripts
- Disk is full

**Solutions**:
```bash
# Check permissions
ls -la ~/.wand/shims/

# Fix if needed
chmod 755 ~/.wand/shims/

# Ensure disk space
df -h

# Recreate shims
wand refresh-shims
```

#### `SHIM_EXECUTION_FAILED`
**When**: Shim script fails to execute

**Common Causes**:
- Binary path changed or was moved
- Permissions removed on binary
- Corrupted shim file
- Shell environment issue

**Solutions**:
```bash
# Verify binary exists
ls ~/.wand/packages/nano/*/bin/nano

# Fix permissions
chmod 755 ~/.wand/packages/nano/*/bin/*

# Recreate shim
wand refresh-shims

# Test execution
~/.wand/shims/nano --version
```

---

### File System Errors

#### `FILE_NOT_FOUND`
**When**: Required file doesn't exist

**Common Causes**:
- Installation incomplete
- File was deleted after installation
- Wrong file path in formula
- Symlink target was moved

**Solutions**:
```bash
# Reinstall package
wand install nano

# Check file locations
find ~/.wand -name nano

# Verify formula definition
cat formulas/nano.yaml
```

#### `DIR_NOT_FOUND`
**When**: Required directory doesn't exist

**Common Causes**:
- Installation directory missing
- Directory deleted after installation
- Wand directory not initialized

**Solutions**:
```bash
# Reinitialize Wand
wand setup

# Recreate directories
mkdir -p ~/.wand/packages
mkdir -p ~/.wand/apps
mkdir -p ~/.wand/shims
```

#### `PERMISSION_DENIED`
**When**: Insufficient permissions for operation

**Common Causes**:
- Files owned by different user
- Directory permissions too restrictive
- SELinux or AppArmor restrictions
- Running in restricted environment

**Solutions**:
```bash
# Fix ownership
sudo chown -R $USER ~/.wand

# Fix permissions
chmod -R u+rwx ~/.wand

# Check security contexts
ls -Z ~/.wand/  # SELinux
aa-status      # AppArmor
```

#### `DISK_SPACE_LOW`
**When**: Insufficient disk space for operation

**Common Causes**:
- Large package download
- Multiple versions installed
- Temporary files not cleaned up
- Disk nearly full

**Solutions**:
```bash
# Check available space
df -h

# Clean up old versions
wand clean

# Remove unused packages
wand uninstall old-package

# Clear downloads
rm -rf ~/.wand/tmp/*
```

---

### Network Errors

#### `NETWORK_UNREACHABLE`
**When**: Cannot connect to remote server

**Common Causes**:
- No internet connectivity
- GitHub API server down
- Firewall blocking connection
- DNS resolution failure
- Proxy configuration needed

**Solutions**:
```bash
# Test connectivity
ping 8.8.8.8              # DNS
ping github.com           # GitHub
curl https://api.github.com  # GitHub API

# Check DNS
nslookup github.com
dig github.com

# Test proxy
export HTTP_PROXY=http://proxy.example.com:8080
wand install nano
```

#### `HTTP_ERROR`
**When**: HTTP request returns error status

**Common Causes**:
- 404: File not found (URL outdated)
- 403: Access forbidden
- 429: Rate limited (too many requests)
- 500: Server error

**Solutions**:
```bash
# For rate limiting, wait and retry
# GitHub allows 60 requests/hour per IP
sleep 60
wand install nano

# For 404 errors, check formula is current
wand info nano

# For access issues, check GitHub status
# https://www.githubstatus.com/
```

#### `TIMEOUT`
**When**: Request takes too long and times out

**Common Causes**:
- Slow internet connection
- Network congestion
- Large file download
- Remote server unresponsive

**Solutions**:
```bash
# Retry with patience
wand install nano

# Check if only specific networks are slow
ping -c 5 github.com

# Increase timeout if needed (not currently configurable)
# Consider setting in firewall rules

# Download at off-peak hours
```

---

### Configuration Errors

#### `CONFIG_MISSING`
**When**: Configuration file doesn't exist when required

**Common Causes**:
- `.wandrc` not found in search path
- `.wand/config.yml` not initialized
- Configuration deleted or moved
- First-time setup not completed

**Solutions**:
```bash
# Initialize Wand
wand setup

# Create project .wandrc
cat > .wandrc <<EOF
cli:
  - name: nano
    version: 8.7
EOF

# Check configuration
wand config show
```

#### `CONFIG_INVALID`
**When**: Configuration file has invalid syntax

**Common Causes**:
- YAML syntax error
- Invalid package format
- Duplicate package names
- Invalid version specification

**Solutions**:
```bash
# Validate configuration
wand validate-config

# Show current config
wand config show

# Edit carefully
vim .wandrc

# Check YAML syntax
cat .wandrc | python3 -c "import sys, yaml; yaml.safe_load(sys.stdin)"
```

#### `REGISTRY_CORRUPTED`
**When**: Package registry cannot be read or is invalid

**Common Causes**:
- Registry file corrupted
- Manual edit introduced invalid JSON
- Disk write error
- File permissions changed

**Solutions**:
```bash
# Backup registry
cp ~/.wand/.registry ~/.wand/.registry.backup

# Rebuild from installed packages
wand refresh-registry

# If registry is corrupted, restore from backup
# and reinstall packages

# Check registry integrity
cat ~/.wand/.registry | python3 -c "import sys, json; json.load(sys.stdin)"
```

---

### System Errors

#### `SYSTEM_NOT_SUPPORTED`
**When**: Feature not available on current operating system

**Common Causes**:
- Using macOS feature on Linux
- Using Linux-specific installation method on macOS
- Trying to use GUI app on headless server
- Shell-specific features not available

**Solutions**:
```bash
# Check supported platforms
wand info nano

# Install platform-specific version
wand install nano@8.7

# Use alternative if available
# For example, use different desktop manager

# On headless systems, skip GUI apps
# Edit wandfile to exclude GUI packages
```

#### `ARCH_NOT_SUPPORTED`
**When**: Package not available for current architecture

**Common Causes**:
- Package only available for x86_64
- ARM build not released
- Trying to install Intel binary on Apple Silicon (or vice versa)
- Custom architecture not officially supported

**Solutions**:
```bash
# Check supported architectures
wand info nano

# Check current architecture
uname -m
arch

# On Apple Silicon, ensure universal binaries
# Some packages may need Rosetta 2
softwareupdate --install-rosetta

# Request support for your architecture
# https://github.com/ochairo/wand/issues
```

---

## Recovery Strategies

### General Recovery Steps

1. **Identify the Error**: Look at error code and message
2. **Collect Information**: Run `wand doctor` for diagnostics
3. **Consult This Guide**: Find error code section
4. **Try Solutions**: Follow recommended solutions in order
5. **Report if Persistent**: Open GitHub issue with full output

### Running Diagnostics

```bash
# Run full diagnostic
wand doctor

# Check system requirements
wand check-requirements

# Verify installation
wand verify-installation

# Test connectivity
wand test-connectivity
```

### Clean Installation

```bash
# Full cleanup
wand clean --all

# Remove all packages
wand uninstall --all

# Reinitialize
wand setup

# Reinstall packages
wand install nano
```

### Getting Help

- **Documentation**: https://github.com/ochairo/wand/blob/main/docs/
- **Issues**: https://github.com/ochairo/wand/issues
- **Discussions**: https://github.com/ochairo/wand/discussions
- **Email**: support@ochairo.com

When reporting issues, include:
1. Full error message with code
2. Output of `wand doctor`
3. Affected package and version
4. Your OS and architecture
5. Steps to reproduce

---

## Common Error Scenarios

### Scenario: Installation Fails for All Packages

**Symptoms**:
- Every package shows `DOWNLOAD_FAILED`
- Network error appears consistently

**Root Cause**: Network connectivity issue

**Solutions**:
```bash
# 1. Check network
ping github.com

# 2. Check GitHub API
curl -I https://api.github.com

# 3. Check firewall/proxy
# Configure if needed

# 4. Retry after network stable
wand install nano
```

### Scenario: Binary Not Found After Installation

**Symptoms**:
- Installation succeeds
- Running command: `command not found`
- Shim exists but fails

**Root Cause**: Binary wasn't extracted correctly or permissions wrong

**Solutions**:
```bash
# 1. Check binary exists
find ~/.wand/packages -name nano

# 2. Fix permissions
chmod +x ~/.wand/packages/*/bin/*

# 3. Refresh shims
wand refresh-shims

# 4. Test
nano --version
```

### Scenario: Permission Denied Errors

**Symptoms**:
- `PERMISSION_DENIED` errors during operations
- Cannot write to directories
- Shim creation fails

**Root Cause**: Wrong file ownership or restrictive permissions

**Solutions**:
```bash
# 1. Fix ownership
sudo chown -R $USER ~/.wand

# 2. Fix permissions
chmod -R u+rwx ~/.wand

# 3. Verify
ls -la ~/.wand

# 4. Retry
wand install nano
```

---

## Error Code Quick Reference

| Code | Category | Severity | Recoverable |
|------|----------|----------|-------------|
| INVALID_PACKAGE_NAME | Input | Low | Yes |
| INVALID_VERSION | Input | Low | Yes |
| INVALID_PATH | Input | Low | Yes |
| INVALID_URL | Input | Medium | Yes |
| PACKAGE_NOT_FOUND | Package | Medium | Yes |
| PACKAGE_NOT_INSTALLED | Package | Low | Yes |
| PACKAGE_INSTALLED | Package | Low | Yes |
| VERSION_NOT_FOUND | Package | Medium | Yes |
| VERSION_INSTALLED | Package | Low | Yes |
| DOWNLOAD_FAILED | Installation | High | Yes |
| CHECKSUM_MISMATCH | Installation | High | Yes |
| EXTRACTION_FAILED | Installation | High | Yes |
| INSTALLATION_FAILED | Installation | High | Yes |
| BINARY_NOT_FOUND | Installation | High | Yes |
| SHIM_CREATION_FAILED | Shim | Medium | Yes |
| SHIM_EXECUTION_FAILED | Shim | High | Yes |
| FILE_NOT_FOUND | FileSystem | Medium | Yes |
| DIR_NOT_FOUND | FileSystem | Medium | Yes |
| PERMISSION_DENIED | FileSystem | High | Yes |
| DISK_SPACE_LOW | FileSystem | High | Yes |
| NETWORK_UNREACHABLE | Network | High | Yes |
| HTTP_ERROR | Network | High | Yes |
| TIMEOUT | Network | Medium | Yes |
| CONFIG_MISSING | Config | Medium | Yes |
| CONFIG_INVALID | Config | High | Yes |
| REGISTRY_CORRUPTED | Config | Critical | Yes |
| SYSTEM_NOT_SUPPORTED | System | Low | No |
| ARCH_NOT_SUPPORTED | System | Low | No |

---

## Error Handling in Scripts

When using Wand programmatically, check error codes:

```bash
#!/bin/bash

# Capture full error
if ! output=$(wand install nano 2>&1); then
    if [[ "$output" == *"PACKAGE_NOT_FOUND"* ]]; then
        echo "Package not found - check available packages"
        wand list
    elif [[ "$output" == *"DOWNLOAD_FAILED"* ]]; then
        echo "Download failed - check network"
        ping github.com
    else
        echo "Unknown error: $output"
        exit 1
    fi
fi

echo "Installation successful"
```

Implement retry logic for transient failures:

```bash
retry_install() {
    local package=$1
    local max_attempts=3
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if wand install "$package"; then
            return 0
        fi

        echo "Attempt $attempt failed, retrying in 10 seconds..."
        sleep 10
        attempt=$((attempt + 1))
    done

    return 1
}

retry_install "nano"
```
