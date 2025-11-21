// main generates formula files from templates.
package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type FormulaTemplate struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Homepage    string            `yaml:"homepage"`
	License     string            `yaml:"license"`
	Tags        []string          `yaml:"tags"`
	Releases    map[string]string `yaml:"releases"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: wand-formula-gen create <package-name>")
			os.Exit(1)
		}
		createFormula(os.Args[2])
	case "validate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: wand-formula-gen validate <formula-file>")
			os.Exit(1)
		}
		validateFormula(os.Args[2])
	case "list-templates":
		listTemplates()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func createFormula(packageName string) {
	formula := FormulaTemplate{
		Name:        packageName,
		Version:     "1.0.0",
		Description: fmt.Sprintf("Brief description of %s", packageName),
		Homepage:    fmt.Sprintf("https://github.com/example/%s", packageName),
		License:     "MIT",
		Tags:        []string{"cli"},
		Releases: map[string]string{
			"darwin-x86_64": fmt.Sprintf("https://github.com/ochairo/potions/releases/download/%s-1.0.0/%s-1.0.0-darwin-x86_64.tar.gz", packageName, packageName),
			"darwin-arm64":  fmt.Sprintf("https://github.com/ochairo/potions/releases/download/%s-1.0.0/%s-1.0.0-darwin-arm64.tar.gz", packageName, packageName),
			"linux-amd64":   fmt.Sprintf("https://github.com/ochairo/potions/releases/download/%s-1.0.0/%s-1.0.0-linux-amd64.tar.gz", packageName, packageName),
			"linux-arm64":   fmt.Sprintf("https://github.com/ochairo/potions/releases/download/%s-1.0.0/%s-1.0.0-linux-arm64.tar.gz", packageName, packageName),
		},
	}

	data, err := yaml.Marshal(&formula)
	if err != nil {
		log.Fatalf("Error marshaling formula: %v", err)
	}

	filename := packageName + ".yaml"
	err = os.WriteFile(filename, data, 0644) //nolint:gosec
	if err != nil {
		log.Fatalf("Error writing formula file: %v", err)
	}

	fmt.Printf("✓ Created formula: %s\n", filename)
	fmt.Println("\nNext steps:")
	fmt.Printf("1. Edit %s with proper details\n", filename)
	fmt.Println("2. Ensure release URLs point to actual binaries")
	fmt.Println("3. Run: wand-formula-gen validate " + filename)
	fmt.Println("4. Test with: wand install " + packageName)
}

func validateFormula(filePath string) {
	data, err := os.ReadFile(filePath) //nolint:gosec
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var formula FormulaTemplate
	err = yaml.Unmarshal(data, &formula)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	errors := []string{}

	if formula.Name == "" {
		errors = append(errors, "name is required")
	}
	if formula.Version == "" {
		errors = append(errors, "version is required")
	}
	if formula.Description == "" {
		errors = append(errors, "description is required")
	}
	if formula.Homepage == "" {
		errors = append(errors, "homepage is required")
	}
	if formula.License == "" {
		errors = append(errors, "license is required")
	}
	if len(formula.Releases) == 0 {
		errors = append(errors, "at least one release is required")
	}

	if len(errors) > 0 {
		fmt.Println("✗ Validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("✓ Formula is valid!")
	fmt.Printf("  Name: %s\n", formula.Name)
	fmt.Printf("  Version: %s\n", formula.Version)
	fmt.Printf("  Releases: %d platforms\n", len(formula.Releases))
}

func listTemplates() {
	fmt.Println("Available templates:")
	fmt.Println()
	fmt.Println("CLI Tool (default)")
	fmt.Println("  Tags: cli, tool, development")
	fmt.Println()
	fmt.Println("Editor")
	fmt.Println("  Tags: editor, cli, text")
	fmt.Println()
	fmt.Println("Browser")
	fmt.Println("  Tags: gui, browser, internet")
	fmt.Println()
	fmt.Println("Language Runtime")
	fmt.Println("  Tags: language, runtime, development")
}

func printUsage() {
	fmt.Println("Wand Formula Generator")
	fmt.Println()
	fmt.Println("Usage: wand-formula-gen <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create <name>      Create a new formula template")
	fmt.Println("  validate <file>    Validate a formula YAML file")
	fmt.Println("  list-templates     Show available templates")
	fmt.Println()
}
