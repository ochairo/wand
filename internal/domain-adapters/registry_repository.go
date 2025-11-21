package domainadapters

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// RegistryRepository implements registry persistence using JSON files
type RegistryRepository struct {
	fs          interfaces.FileSystem
	registryDir string
}

// NewRegistryRepository creates a new RegistryRepository
func NewRegistryRepository(fs interfaces.FileSystem, wandDir string) interfaces.RegistryRepository {
	return &RegistryRepository{
		fs:          fs,
		registryDir: wandDir,
	}
}

// Load loads the registry from disk
func (r *RegistryRepository) Load() (*entities.Registry, error) {
	registryPath := filepath.Join(r.registryDir, "registry.json")

	// If registry doesn't exist, return empty registry
	if !r.fs.Exists(registryPath) {
		return &entities.Registry{
			Packages:       make(map[string]*entities.PackageEntry),
			GlobalVersions: make(map[string]string),
		}, nil
	}

	data, err := r.fs.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry file: %w", err)
	}

	var registry entities.Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	return &registry, nil
}

// Save saves the registry to disk
func (r *RegistryRepository) Save(registry *entities.Registry) error {
	registryPath := filepath.Join(r.registryDir, "registry.json")

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize registry: %w", err)
	}

	if err := r.fs.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}

	return nil
}

// Exists checks if the registry file exists
func (r *RegistryRepository) Exists() bool {
	registryPath := filepath.Join(r.registryDir, "registry.json")
	return r.fs.Exists(registryPath)
}
