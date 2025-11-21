package domainadapters

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// WandfileRepository implements wandfile file operations
type WandfileRepository struct {
	fs interfaces.FileSystem
}

// NewWandfileRepository creates a new WandfileRepository
func NewWandfileRepository(fs interfaces.FileSystem) interfaces.WandfileRepository {
	return &WandfileRepository{
		fs: fs,
	}
}

// Load loads a wandfile from the specified path
func (r *WandfileRepository) Load(path string) (*entities.Wandfile, error) {
	data, err := r.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read wandfile: %w", err)
	}

	var wandfile entities.Wandfile
	if err := yaml.Unmarshal(data, &wandfile); err != nil {
		return nil, fmt.Errorf("failed to parse wandfile: %w", err)
	}

	return &wandfile, nil
}

// Save saves a wandfile to the specified path
func (r *WandfileRepository) Save(path string, wandfile *entities.Wandfile) error {
	data, err := yaml.Marshal(wandfile)
	if err != nil {
		return fmt.Errorf("failed to serialize wandfile: %w", err)
	}

	if err := r.fs.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write wandfile: %w", err)
	}

	return nil
}

// Exists checks if a wandfile exists at the specified path
func (r *WandfileRepository) Exists(path string) bool {
	return r.fs.Exists(path)
}
