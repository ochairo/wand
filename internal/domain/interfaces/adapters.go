// Package interfaces defines the domain interfaces.
package interfaces

import "io"

// Downloader defines the interface for downloading files
type Downloader interface {
	Download(url, destPath string) error
	DownloadWithProgress(url, destPath string, progress io.Writer) error
	VerifyChecksum(filePath, checksumURL string) error
}

// Extractor defines the interface for extracting archives
type Extractor interface {
	Extract(archivePath, destDir string) error
	ExtractFile(archivePath, fileName, destPath string) error
}

// FileSystem defines the interface for file system operations
type FileSystem interface {
	Exists(path string) bool
	IsDir(path string) bool
	MkdirAll(path string, perm uint32) error
	Remove(path string) error
	RemoveAll(path string) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm uint32) error
	Symlink(oldname, newname string) error
	ReadSymlink(name string) (string, error)
	Chmod(name string, mode uint32) error
	Walk(root string, walkFn func(path string, isDir bool, err error) error) error
}

// ShellExecutor defines the interface for executing shell commands
type ShellExecutor interface {
	Execute(command string, args ...string) (string, error)
	ExecuteInDir(dir, command string, args ...string) (string, error)
	ExecuteWithEnv(env map[string]string, command string, args ...string) (string, error)
}

// GitClient defines the interface for Git operations
type GitClient interface {
	Clone(repoURL, destDir string) error
	Pull(repoDir string) error
	Status(repoDir string) (string, error)
	Add(repoDir string, files ...string) error
	Commit(repoDir, message string) error
	Push(repoDir string) error
}
