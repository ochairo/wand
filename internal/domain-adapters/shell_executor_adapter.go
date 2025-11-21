package domainadapters

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ochairo/wand/internal/domain/interfaces"
)

// ShellExecutorAdapter implements shell command execution
type ShellExecutorAdapter struct{}

// NewShellExecutorAdapter creates a new ShellExecutorAdapter
func NewShellExecutorAdapter() interfaces.ShellExecutor {
	return &ShellExecutorAdapter{}
}

// Execute runs a command with arguments
func (s *ShellExecutorAdapter) Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// ExecuteInDir runs a command in a specific directory
func (s *ShellExecutorAdapter) ExecuteInDir(dir, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// ExecuteWithEnv runs a command with custom environment variables
func (s *ShellExecutorAdapter) ExecuteWithEnv(env map[string]string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Set environment variables
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}
