# Wand API Usage Example

This example demonstrates how to use wand programmatically as a Go library.

## Overview

The wand public API allows you to integrate package management into your own tools and applications. This example shows common operations like listing formulas, searching packages, and managing installations.

## Building

```bash
go build -o wand-example examples/api-usage/main.go
```

## Running

```bash
./wand-example
```

## Features Demonstrated

### 1. Formula Discovery
- List all available formulas
- Search formulas by keyword
- Get formula details

### 2. Package Management
- List installed packages
- Check package versions
- Get global version settings

### 3. Version Information
- Check available versions for packages
- Compare installed vs. available versions

## API Reference

The example uses `github.com/ochairo/wand/pkg/client` which provides:

```go
// Create a client
client, err := client.New("") // Uses ~/.wand by default

// List formulas
formulas, err := client.ListFormulas()

// Search for specific packages
results, err := client.SearchFormulas("json")

// List installed packages
packages, err := client.ListPackages()

// Get global version
version, err := client.GetGlobalVersion("neovim")
```

## Use Cases

- **CI/CD Integration**: Automate package installation in build pipelines
- **Development Tools**: Build custom package managers on top of wand
- **System Automation**: Script complex environment setups
- **Monitoring**: Track installed package versions across systems

## See Also

- [Wand Public API Documentation](../../pkg/client/)
- [CLI Usage](../../docs/GETTING_STARTED.md)
- [Wandfile Examples](../wandfile/)
