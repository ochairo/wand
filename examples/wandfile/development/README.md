# Development Environment Wandfile Example

Complete development environment with language runtimes, tools, and utilities.

## Usage

```bash
wand install --wandfile ./Wandfile
```

## What's Included

### Language Runtimes
- **Node.js** - JavaScript runtime
- **Python** - Python interpreter
- **Go** - Go compiler
- **Rust** - Rust toolchain

### Development Tools
- **Git** - Version control
- **Docker** - Containerization
- **Neovim** - Text editor

### CLI Utilities
- **jq, ripgrep, fd, bat** - Essential CLI tools
- **lazygit** - Git TUI
- **gh** - GitHub CLI

### Build Tools
- **make** - Build automation
- **cmake** - Cross-platform build system

## Per-Project Override

Use `.wandrc` in your project to override versions:

```yaml
# .wandrc in your project
packages:
  - name: node
    version: "18.0.0"  # Project uses Node 18
```

## Team Usage

Share this Wandfile with your team to ensure everyone has the same tools:

```bash
git add Wandfile
git commit -m "Add development environment config"
git push
```

Team members can then run:

```bash
wand install --wandfile ./Wandfile
```
