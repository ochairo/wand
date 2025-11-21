# Getting Started with Wand

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/ochairo/wand/main/install.sh | sh
```

See [Installation Guide](INSTALLATION.md) for detailed setup.

## Your First Package

```bash
# Search for a package
wand search nano

# Install it
wand install nano

# Verify
nano --version

# Uninstall
wand uninstall nano
```

## Common Commands

```bash
wand list              # Show installed packages
wand search nano       # Find packages
wand install nano      # Install latest
wand install nano@8.2  # Install specific version
wand update nano       # Update to latest
wand uninstall nano    # Remove package
wand doctor            # Check system health
```

## What's a Formula?

A **formula** is a package specification (YAML file) that tells Wand how to download and install a tool. The `ochairo/potions` repository contains all available formulas.

## What's a Shim?

A **shim** is a shell script in `~/.wand/shims/` that launches the actual installed binary. When you install `nano`, Wand creates a `nano` shim so you can just type `nano` to run it.

## Project Configuration

Create a `wandfile.yaml` to pin versions for your project:

```yaml
packages:
  - name: node
    version: "20.10.0"
  - name: make
    version: "4.4"
```

Then install all at once:

```bash
wand install
```

## Documentation

- **[Installation](INSTALLATION.md)** - Setup details
- **[Command Reference](COMMAND_REFERENCE.md)** - All commands
- **[Error Codes](ERROR_CODES.md)** - Error reference and solutions
- **[Formula Guide](FORMULA_GUIDE.md)** - Creating formulas
- **[Deployment](DEPLOYMENT.md)** - Production setup

## Help

- 💬 [GitHub Discussions](https://github.com/ochairo/wand/discussions)
- 🐛 [GitHub Issues](https://github.com/ochairo/wand/issues)
