<div align="center">

# 🪄 Wand

A package manager with shim-based version control for CLI tools, GUI apps, and dotfiles.

[![Version](https://img.shields.io/github/v/release/ochairo/wand?label=version)](https://github.com/ochairo/wand/releases)
[![Tests](https://img.shields.io/github/actions/workflow/status/ochairo/wand/test.yml?branch=main&label=tests&logo=github)](https://github.com/ochairo/wand/actions/workflows/test.yml)
[![Security](https://img.shields.io/github/actions/workflow/status/ochairo/wand/security-scan.yml?branch=main&label=security&logo=github)](https://github.com/ochairo/wand/actions/workflows/security-scan.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ochairo/wand)](https://goreportcard.com/report/github.com/ochairo/wand)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/ochairo/wand/badge)](https://scorecard.dev/viewer/?uri=github.com/ochairo/wand)
[![License](https://img.shields.io/github/license/ochairo/wand)](https://github.com/ochairo/wand/blob/main/LICENSE)

[Features](#-features) • [Formulas](#-supported-formulas) • [Quick Start](#-quick-start) • [Third-Party Integration](#-third-party-integration) • [Documentation](#-documentation)

</div>

## ✨ Features

- **CLI Tools with Version Management**: Install multiple versions, switch per-project
- **GUI Applications**: Single-version installs for desktop apps
- **Dotfiles Management**: Git-based dotfile repository with symlink mapping
- **Declarative Configuration**: Wandfile for system-wide tool configuration
- **Per-Project Overrides**: `.wandrc` files for project-specific version pinning

## 🔮 Supported Formulas

- [View all supported formulas →](./formulas/)

## 🚀 Quick Start

```bash
# Installation
curl -sSL https://raw.githubusercontent.com/ochairo/wand/main/scripts/install.sh | bash
```

```bash
# List available versions
wand list jq --remote

# Install latest version
wand install jq

# Install specific version
wand install jq@1.7.1

# Check installed versions
wand list jq

# Switch versions
wand switch jq@1.6.0

# Show active version details
wand info jq
```

## 🔌 Third-Party Integration

Wand provides a public API for building custom integrations like TUIs, web dashboards, and IDE extensions.

> **Building something cool with wand?** We'd love to see it! Share your project in [GitHub Discussions](https://github.com/ochairo/wand/discussions).

### Use Cases

- **Terminal UIs (TUI)**: Build interactive package browsers with [Bubble Tea](https://github.com/charmbracelet/bubbletea) or [tview](https://github.com/rivo/tview)
- **Web Dashboards**: Create web-based package management interfaces
- **IDE Extensions**: Integrate wand into VS Code, IntelliJ, or other editors
- **CI/CD Tools**: Automate package installations in build pipelines
- **Custom Workflows**: Build domain-specific package management tools

### Quick Example

```go
import "github.com/ochairo/wand/pkg/client"

// Initialize client
c, _ := client.New("")

// Install a package
c.Install("jq", "1.7.1")

// List installed packages
packages, _ := c.ListPackages()

// Get package details
formula, _ := c.GetFormula("jq")
```

See [pkg/README.md](./pkg/README.md) for complete API documentation and [examples/](./examples/).

## 🏛️ Documentation

- [User Guide](./docs/GETTING_STARTED.md) - Getting started and usage instructions
- [Public API](./pkg/README.md) - Programmatic API for third-party integrations
- [Contributing](./docs/CONTRIBUTING.md) - Development setup and adding packages

<br><br>

<div align="center">

[Report Bug](https://github.com/ochairo/wand/issues) • [Request Feature](https://github.com/ochairo/wand/issues) • [Documentation](./docs/)

**Made with ❤︎ by [ochairo](https://github.com/ochairo)**

</div>
