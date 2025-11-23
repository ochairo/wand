package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ochairo/wand/internal/domain/entities"
	errs "github.com/ochairo/wand/internal/domain/errors"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// ShellPluginService handles shell plugin installations
type ShellPluginService struct {
	fs            interfaces.FileSystem
	shellExecutor interfaces.ShellExecutor
	homeDir       string
	wandDir       string
}

// NewShellPluginService creates a new shell plugin service
func NewShellPluginService(
	fs interfaces.FileSystem,
	shellExecutor interfaces.ShellExecutor,
	homeDir string,
	wandDir string,
) *ShellPluginService {
	return &ShellPluginService{
		fs:            fs,
		shellExecutor: shellExecutor,
		homeDir:       homeDir,
		wandDir:       wandDir,
	}
}

const (
	wandShellMarkerStart = "# >>> wand shell plugins >>>"
	wandShellMarkerEnd   = "# <<< wand shell plugins <<<"
)

// Install installs a shell plugin from Git
func (s *ShellPluginService) Install(formula *entities.Formula, version string) error {
	if formula.ShellPlugin == nil {
		return errs.New(errs.ErrConfigInvalid, "Formula is not a shell plugin")
	}

	config := formula.ShellPlugin
	installPath := filepath.Join(s.wandDir, config.InstallPath)

	// Clone or update repository
	if s.fs.Exists(installPath) {
		// Already exists, pull latest
		if _, err := s.shellExecutor.ExecuteInDir(installPath, "git pull origin main || git pull origin master"); err != nil {
			return errs.Wrap(errs.ErrNetworkUnreachable, "Failed to update shell plugin", err)
		}

		// Checkout specific version if not "latest"
		if version != "latest" && version != "" {
			if _, err := s.shellExecutor.ExecuteInDir(installPath, fmt.Sprintf("git checkout %s", version)); err != nil {
				return errs.Wrap(errs.ErrVersionNotFound, fmt.Sprintf("Failed to checkout version %s", version), err)
			}
		}
	} else {
		// Clone fresh
		if _, err := s.shellExecutor.Execute("git", "clone", config.GitURL, installPath); err != nil {
			return errs.Wrap(errs.ErrNetworkUnreachable, "Failed to clone shell plugin repository", err)
		}

		// Checkout specific version if not "latest"
		if version != "latest" && version != "" {
			if _, err := s.shellExecutor.ExecuteInDir(installPath, fmt.Sprintf("git checkout %s", version)); err != nil {
				return errs.Wrap(errs.ErrVersionNotFound, fmt.Sprintf("Failed to checkout version %s", version), err)
			}
		}
	}

	// Configure shell
	if err := s.configureShell(formula); err != nil {
		return errs.Wrap(errs.ErrConfigInvalid, "Failed to configure shell", err)
	}

	return nil
}

// configureShell adds source lines to the shell config file
func (s *ShellPluginService) configureShell(formula *entities.Formula) error {
	config := formula.ShellPlugin
	shellConfigPath := formula.GetShellConfigPath(s.homeDir)

	if shellConfigPath == "" {
		return errs.New(errs.ErrConfigInvalid, "Unsupported shell: "+config.Shell)
	}

	// Ensure config file exists
	if !s.fs.Exists(shellConfigPath) {
		// Create empty config file
		if err := s.fs.WriteFile(shellConfigPath, []byte("# Created by wand\n"), 0644); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to create shell config file", err)
		}
	}

	// Read current config
	content, err := s.fs.ReadFile(shellConfigPath)
	if err != nil {
		return errs.Wrap(errs.ErrFileNotFound, "Failed to read shell config", err)
	}

	configStr := string(content)

	// Build source line(s)
	var sourceLines []string
	if config.CustomSetup && len(config.SourceLines) > 0 {
		// Use custom source lines (e.g., oh-my-zsh)
		sourceLines = config.SourceLines
	} else {
		// Standard plugin source line
		installPath := filepath.Join(s.wandDir, config.InstallPath)
		sourcePath := filepath.Join(installPath, config.SourceFile)
		sourceLines = []string{fmt.Sprintf("source %s", sourcePath)}
	}

	// Check if already configured (avoid duplicates)
	for _, line := range sourceLines {
		if strings.Contains(configStr, line) {
			// Already configured, skip
			return nil
		}
	}

	// Check if wand section exists
	hasWandSection := strings.Contains(configStr, wandShellMarkerStart)

	var newConfig string
	if hasWandSection {
		// Insert into existing wand section
		newConfig = s.insertIntoWandSection(configStr, sourceLines)
	} else {
		// Create new wand section
		newConfig = s.createWandSection(configStr, sourceLines)
	}

	// Write updated config
	if err := s.fs.WriteFile(shellConfigPath, []byte(newConfig), 0644); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to update shell config", err)
	}

	return nil
}

// insertIntoWandSection inserts source lines into existing wand section
func (s *ShellPluginService) insertIntoWandSection(config string, sourceLines []string) string {
	lines := strings.Split(config, "\n")
	result := make([]string, 0, len(lines)+len(sourceLines))

	inserted := false

	for _, line := range lines {
		if strings.Contains(line, wandShellMarkerStart) {
			result = append(result, line)
			continue
		}

		if strings.Contains(line, wandShellMarkerEnd) {
			// Insert new lines before the end marker
			if !inserted {
				result = append(result, sourceLines...)
				inserted = true
			}
			result = append(result, line)
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// createWandSection creates a new wand section with source lines
func (s *ShellPluginService) createWandSection(config string, sourceLines []string) string {
	// Add wand section at the end
	section := "\n" + wandShellMarkerStart + "\n"
	for _, line := range sourceLines {
		section += line + "\n"
	}
	section += wandShellMarkerEnd + "\n"

	return config + section
}

// Uninstall removes a shell plugin
func (s *ShellPluginService) Uninstall(formula *entities.Formula) error {
	if formula.ShellPlugin == nil {
		return errs.New(errs.ErrConfigInvalid, "Formula is not a shell plugin")
	}

	config := formula.ShellPlugin
	installPath := filepath.Join(s.wandDir, config.InstallPath)

	// Remove directory
	if s.fs.Exists(installPath) {
		if err := s.fs.RemoveAll(installPath); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to remove shell plugin directory", err)
		}
	}

	// Remove from shell config
	if err := s.unconfigureShell(formula); err != nil {
		return errs.Wrap(errs.ErrConfigInvalid, "Failed to unconfigure shell", err)
	}

	return nil
}

// unconfigureShell removes source lines from shell config
func (s *ShellPluginService) unconfigureShell(formula *entities.Formula) error {
	config := formula.ShellPlugin
	shellConfigPath := formula.GetShellConfigPath(s.homeDir)

	if shellConfigPath == "" || !s.fs.Exists(shellConfigPath) {
		return nil
	}

	// Read current config
	content, err := s.fs.ReadFile(shellConfigPath)
	if err != nil {
		return errs.Wrap(errs.ErrFileNotFound, "Failed to read shell config", err)
	}

	configStr := string(content)

	// Build source line(s) to remove
	var sourceLines []string
	if config.CustomSetup && len(config.SourceLines) > 0 {
		sourceLines = config.SourceLines
	} else {
		installPath := filepath.Join(s.wandDir, config.InstallPath)
		sourcePath := filepath.Join(installPath, config.SourceFile)
		sourceLines = []string{fmt.Sprintf("source %s", sourcePath)}
	}

	// Remove lines
	lines := strings.Split(configStr, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		shouldRemove := false
		for _, sourceLine := range sourceLines {
			if strings.TrimSpace(line) == strings.TrimSpace(sourceLine) {
				shouldRemove = true
				break
			}
		}

		if !shouldRemove {
			result = append(result, line)
		}
	}

	// Clean up empty wand section
	newConfig := strings.Join(result, "\n")
	newConfig = s.cleanupEmptyWandSection(newConfig)

	// Write updated config
	if err := s.fs.WriteFile(shellConfigPath, []byte(newConfig), 0644); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to update shell config", err)
	}

	return nil
}

// cleanupEmptyWandSection removes wand section if it's empty
func (s *ShellPluginService) cleanupEmptyWandSection(config string) string {
	lines := strings.Split(config, "\n")
	result := make([]string, 0, len(lines))

	inWandSection := false
	wandSectionLines := make([]string, 0)

	for _, line := range lines {
		if strings.Contains(line, wandShellMarkerStart) {
			inWandSection = true
			wandSectionLines = []string{line}
			continue
		}

		if strings.Contains(line, wandShellMarkerEnd) {
			inWandSection = false
			wandSectionLines = append(wandSectionLines, line)

			// Check if section has any content besides markers
			hasContent := false
			for _, sectionLine := range wandSectionLines {
				trimmed := strings.TrimSpace(sectionLine)
				if trimmed != "" && !strings.Contains(trimmed, ">>>") && !strings.Contains(trimmed, "<<<") {
					hasContent = true
					break
				}
			}

			// Only add section if it has content
			if hasContent {
				result = append(result, wandSectionLines...)
			}
			wandSectionLines = nil
			continue
		}

		if inWandSection {
			wandSectionLines = append(wandSectionLines, line)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
