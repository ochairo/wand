# Wand Public API

This directory contains the public API for wand, enabling third-party integrations, TUIs, web UIs, and other programmatic interactions.

## Packages

### `pkg/types`
Public types representing wand entities:
- `Package` - Installed package information
- `Version` - Semantic version representation
- `Formula` - Package definition/metadata
- `Registry` - Local package registry state
- `PackageEntry` - Grouped package versions

### `pkg/api`
Public interfaces for wand operations:
- `Installer` - Package installation/uninstallation
- `RegistryManager` - Registry state management
- `FormulaProvider` - Formula lookup
- `VersionResolver` - Version resolution

### `pkg/client`
High-level client for wand operations:
- `Client` - Main entry point for all wand functionality

## Usage Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/ochairo/wand/pkg/client"
)

func main() {
    // Create a new wand client (uses ~/.wand by default)
    c, err := client.New("")
    if err != nil {
        log.Fatal(err)
    }

    // Install a package
    pkg, err := c.Install("jq", "1.7.1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Installed %s\n", pkg.Identifier())

    // List available versions
    versions, err := c.ListAvailableVersions("jq")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Available versions: %d\n", len(versions))

    // List installed packages
    packages, err := c.ListPackages()
    if err != nil {
        log.Fatal(err)
    }

    for _, entry := range packages {
        fmt.Printf("%s (%s):\n", entry.Name, entry.Type)
        for version := range entry.Versions {
            fmt.Printf("  - %s\n", version)
        }
    }

    // Set global version
    err = c.SetGlobalVersion("jq", "1.7.1")
    if err != nil {
        log.Fatal(err)
    }

    // Get registry state
    registry, err := c.GetRegistry()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Registry last updated: %s\n", registry.UpdatedAt)
}
```

## Building a TUI

Example TUI using the public API:

```go
package main

import (
    "github.com/ochairo/wand/pkg/client"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    wand     *client.Client
    packages []string
    cursor   int
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q":
            return m, tea.Quit
        case "i":
            // Install selected package
            packageName := m.packages[m.cursor]
            m.wand.Install(packageName, "latest")
        }
    }
    return m, nil
}

func (m model) View() string {
    s := "Wand TUI\n\n"
    // Render package list...
    return s
}

func main() {
    c, _ := client.New("")

    p := tea.NewProgram(model{
        wand: c,
    })
    p.Run()
}
```

## Web API

Example HTTP server exposing wand functionality:

```go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/ochairo/wand/pkg/client"
)

var wandClient *client.Client

func main() {
    wandClient, _ = client.New("")

    http.HandleFunc("/packages", listPackages)
    http.HandleFunc("/install", installPackage)
    http.ListenAndServe(":8080", nil)
}

func listPackages(w http.ResponseWriter, r *http.Request) {
    packages, err := wandClient.ListPackages()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(packages)
}

func installPackage(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name    string `json:"name"`
        Version string `json:"version"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    pkg, err := wandClient.Install(req.Name, req.Version)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(pkg)
}
```

## API Stability

The public API (`pkg/` directory) follows semantic versioning:
- **Minor version bumps**: New methods/fields (backwards compatible)
- **Major version bumps**: Breaking changes

The `internal/` directory may change without notice.
