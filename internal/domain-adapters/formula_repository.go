package domainadapters

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// FormulaRepository implements formula loading from YAML files
type FormulaRepository struct {
	fs          interfaces.FileSystem
	formulasDir string
}

// NewFormulaRepository creates a new FormulaRepository
func NewFormulaRepository(fs interfaces.FileSystem, formulasDir string) interfaces.FormulaRepository {
	return &FormulaRepository{
		fs:          fs,
		formulasDir: formulasDir,
	}
}

// GetFormula loads a formula by name
func (r *FormulaRepository) GetFormula(name string) (*entities.Formula, error) {
	formulaPath := filepath.Join(r.formulasDir, name+".yaml")

	if !r.fs.Exists(formulaPath) {
		return nil, fmt.Errorf("formula not found: %s", name)
	}

	data, err := r.fs.ReadFile(formulaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read formula file: %w", err)
	}

	var formula entities.Formula
	if err := yaml.Unmarshal(data, &formula); err != nil {
		return nil, fmt.Errorf("failed to parse formula: %w", err)
	}

	return &formula, nil
}

// ListFormulas lists all available formulas
func (r *FormulaRepository) ListFormulas() ([]*entities.Formula, error) {
	var formulas []*entities.Formula

	err := r.fs.Walk(r.formulasDir, func(path string, isDir bool, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if isDir || filepath.Ext(path) != ".yaml" {
			return nil
		}

		data, err := r.fs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var formula entities.Formula
		if err := yaml.Unmarshal(data, &formula); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		formulas = append(formulas, &formula)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return formulas, nil
}

// Sync updates the formulas directory from the remote repository
func (r *FormulaRepository) Sync() error {
	const repoURL = "https://github.com/ochairo/potions.git"

	// Check if formulas directory is a git repository
	gitDir := filepath.Join(r.formulasDir, ".git")
	isGitRepo := r.fs.Exists(gitDir)

	if isGitRepo {
		// Pull latest changes if already a git repo
		cmd := exec.Command("git", "-C", r.formulasDir, "pull", "origin", "main") //nolint:gosec // G204: hardcoded git command
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to pull formulas: %w", err)
		}
	} else {
		// Clone repository if not yet cloned
		cmd := exec.Command("git", "clone", repoURL, r.formulasDir) //nolint:gosec // G204: hardcoded git command
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone formulas repository: %w", err)
		}
	}

	return nil
}
