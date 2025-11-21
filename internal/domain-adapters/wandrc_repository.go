package domainadapters

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// WandRCRepository implements .wandrc file operations
type WandRCRepository struct {
	fs interfaces.FileSystem
}

// NewWandRCRepository creates a new WandRCRepository
func NewWandRCRepository(fs interfaces.FileSystem) interfaces.WandRCRepository {
	return &WandRCRepository{
		fs: fs,
	}
}

// Load loads a .wandrc file from a directory
func (r *WandRCRepository) Load(dir string) (*entities.WandRC, error) {
	wandrcPath := filepath.Join(dir, ".wandrc")

	data, err := r.fs.ReadFile(wandrcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .wandrc: %w", err)
	}

	var wandrc entities.WandRC
	if err := yaml.Unmarshal(data, &wandrc); err != nil {
		return nil, fmt.Errorf("failed to parse .wandrc: %w", err)
	}

	return &wandrc, nil
}

// Save saves a .wandrc file to a directory
func (r *WandRCRepository) Save(dir string, wandrc *entities.WandRC) error {
	wandrcPath := filepath.Join(dir, ".wandrc")

	data, err := yaml.Marshal(wandrc)
	if err != nil {
		return fmt.Errorf("failed to serialize .wandrc: %w", err)
	}

	if err := r.fs.WriteFile(wandrcPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write .wandrc: %w", err)
	}

	return nil
}

// Exists checks if a .wandrc file exists in a directory
func (r *WandRCRepository) Exists(dir string) bool {
	wandrcPath := filepath.Join(dir, ".wandrc")
	return r.fs.Exists(wandrcPath)
}

// FindInPath searches for .wandrc in the directory and parent directories
func (r *WandRCRepository) FindInPath(startDir string) (*entities.WandRC, string, error) {
	dir := startDir

	for {
		if r.Exists(dir) {
			wandrc, err := r.Load(dir)
			if err != nil {
				return nil, "", err
			}
			wandrcPath := filepath.Join(dir, ".wandrc")
			return wandrc, wandrcPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir || parent == "/" || parent == "." {
			// Reached root
			break
		}
		dir = parent
	}

	return nil, "", fmt.Errorf(".wandrc not found in path")
}
