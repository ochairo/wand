// Package domainadapters provides adapter implementations for domain interfaces.
package domainadapters

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// DotfileRepository implements dotfile configuration persistence
type DotfileRepository struct {
	fs      interfaces.FileSystem
	wandDir string
}

// NewDotfileRepository creates a new DotfileRepository
func NewDotfileRepository(fs interfaces.FileSystem, wandDir string) interfaces.DotfileRepository {
	return &DotfileRepository{
		fs:      fs,
		wandDir: wandDir,
	}
}

// Load loads the dotfile configuration from disk
func (r *DotfileRepository) Load() (*entities.DotfileConfig, error) {
	configPath := filepath.Join(r.wandDir, "dotfiles.json")

	if !r.fs.Exists(configPath) {
		return nil, fmt.Errorf("dotfile config not found")
	}

	data, err := r.fs.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dotfile config: %w", err)
	}

	var config entities.DotfileConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse dotfile config: %w", err)
	}

	return &config, nil
}

// Save saves the dotfile configuration to disk
func (r *DotfileRepository) Save(config *entities.DotfileConfig) error {
	configPath := filepath.Join(r.wandDir, "dotfiles.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize dotfile config: %w", err)
	}

	if err := r.fs.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write dotfile config: %w", err)
	}

	return nil
}

// Exists checks if the dotfile configuration file exists
func (r *DotfileRepository) Exists() bool {
	configPath := filepath.Join(r.wandDir, "dotfiles.json")
	return r.fs.Exists(configPath)
}
