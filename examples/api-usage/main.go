// Example program demonstrating the wand public API.
//
// This shows how third-party tools can integrate with wand programmatically.
//
// Build: go build -o wand-example examples/api-usage/main.go
// Run:   ./wand-example
package main

import (
	"fmt"
	"log"

	"github.com/ochairo/wand/pkg/client"
)

func main() {
	fmt.Println("=== Wand API Example ===")
	fmt.Println()

	// Create a new wand client
	// Empty string uses default ~/.wand directory
	c, err := client.New("")
	if err != nil {
		log.Fatalf("Failed to create wand client: %v", err)
	}

	// Example 1: List all formulas
	fmt.Println("1. Available formulas:")
	formulas, err := c.ListFormulas()
	if err != nil {
		log.Printf("  Error listing formulas: %v", err)
	} else {
		count := len(formulas)
		if count > 5 {
			count = 5
		}
		for i := 0; i < count; i++ {
			formula := formulas[i]
			fmt.Printf("  - %s: %s\n", formula.Name, formula.Description)
		}
		fmt.Printf("  ... and %d more\n", len(formulas)-count)
	}
	fmt.Println()

	// Example 2: Search for formulas
	fmt.Println("2. Search formulas (query: 'json'):")
	results, err := c.SearchFormulas("json")
	if err != nil {
		log.Printf("  Error searching: %v", err)
	} else {
		for _, formula := range results {
			fmt.Printf("  - %s (%s)\n", formula.Name, formula.Type)
		}
	}
	fmt.Println()

	// Example 3: List installed packages
	fmt.Println("3. Installed packages:")
	packages, err := c.ListPackages()
	if err != nil {
		log.Printf("  Error listing packages: %v", err)
	} else {
		if len(packages) == 0 {
			fmt.Println("  No packages installed yet")
		} else {
			for _, entry := range packages {
				fmt.Printf("  %s (%s):\n", entry.Name, entry.Type)
				for version := range entry.Versions {
					globalVersion, _ := c.GetGlobalVersion(entry.Name)
					marker := " "
					if version == globalVersion {
						marker = "*"
					}
					fmt.Printf("    %s %s\n", marker, version)
				}
			}
		}
	}
	fmt.Println()

	// Example 4: Check available versions for a package
	packageName := "jq"
	fmt.Printf("4. Available versions of %s:\n", packageName)
	versions, err := c.ListAvailableVersions(packageName)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		count := len(versions)
		if count > 5 {
			count = 5
		}
		for i := 0; i < count; i++ {
			fmt.Printf("  - %s\n", versions[i].String())
		}
		if len(versions) > count {
			fmt.Printf("  ... and %d more\n", len(versions)-count)
		}
	}
	fmt.Println()

	// Example 5: Get registry state
	fmt.Println("5. Registry state:")
	registry, err := c.GetRegistry()
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		fmt.Printf("  Total packages: %d\n", len(registry.Packages))
		fmt.Printf("  Global versions set: %d\n", len(registry.GlobalVersions))
		fmt.Printf("  Last updated: %s\n", registry.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// Example 6: Get formula details
	fmt.Printf("6. Formula details for %s:\n", packageName)
	formula, err := c.GetFormula(packageName)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		fmt.Printf("  Name: %s\n", formula.Name)
		fmt.Printf("  Type: %s\n", formula.Type)
		fmt.Printf("  Description: %s\n", formula.Description)
		fmt.Printf("  Homepage: %s\n", formula.Homepage)
		fmt.Printf("  License: %s\n", formula.License)
		if len(formula.Tags) > 0 {
			fmt.Printf("  Tags: %v\n", formula.Tags)
		}
	}

	fmt.Println("\n=== API Example Complete ===")
	fmt.Println("\nTo install a package:")
	fmt.Println("  pkg, err := c.Install(\"jq\", \"1.7.1\")")
	fmt.Println("\nTo uninstall:")
	fmt.Println("  err := c.Uninstall(\"jq\", \"1.7.1\")")
	fmt.Println("\nTo set global version:")
	fmt.Println("  err := c.SetGlobalVersion(\"jq\", \"1.7.1\")")
}
