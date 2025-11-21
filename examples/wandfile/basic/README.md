# Basic Wandfile Example

A minimal Wandfile with essential CLI tools.

## Usage

```bash
# Install all packages
wand install --wandfile ./Wandfile

# Or just copy to your project
cp Wandfile ~/myproject/
cd ~/myproject
wand install
```

## What's Included

- **jq** - JSON processor
- **ripgrep** - Fast text search
- **fd** - Fast file finder
- **bat** - Cat clone with syntax highlighting

## Customization

Edit the `Wandfile` and specify exact versions:

```yaml
packages:
  - name: jq
    version: "1.7.1"  # Pin to specific version
```

Or use latest:

```yaml
packages:
  - name: jq
    version: latest  # Always get latest
```
