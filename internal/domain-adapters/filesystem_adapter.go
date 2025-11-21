package domainadapters

import (
	"os"
	"path/filepath"

	"github.com/ochairo/wand/internal/domain/interfaces"
)

// FileSystemAdapter implements the FileSystem interface using os package
type FileSystemAdapter struct{}

// NewFileSystemAdapter creates a new FileSystemAdapter
func NewFileSystemAdapter() interfaces.FileSystem {
	return &FileSystemAdapter{}
}

// Exists checks if a file or directory exists
func (fs *FileSystemAdapter) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir checks if path is a directory
func (fs *FileSystemAdapter) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// MkdirAll creates a directory and all parent directories
func (fs *FileSystemAdapter) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// Remove removes a file or empty directory
func (fs *FileSystemAdapter) Remove(path string) error {
	return os.Remove(path)
}

// RemoveAll removes a path and all its contents
func (fs *FileSystemAdapter) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// ReadFile reads the entire file
func (fs *FileSystemAdapter) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path) //nolint:gosec // G304: path is from application config //nolint:gosec // G304: path is from application config
}

// WriteFile writes data to a file
func (fs *FileSystemAdapter) WriteFile(path string, data []byte, perm uint32) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

// Symlink creates a symbolic link
func (fs *FileSystemAdapter) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

// ReadSymlink returns the target of a symbolic link
func (fs *FileSystemAdapter) ReadSymlink(name string) (string, error) {
	return os.Readlink(name)
}

// Chmod changes the file mode
func (fs *FileSystemAdapter) Chmod(name string, mode uint32) error {
	return os.Chmod(name, os.FileMode(mode))
}

// Walk walks the file tree
func (fs *FileSystemAdapter) Walk(root string, walkFn func(path string, isDir bool, err error) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return walkFn(path, false, err)
		}
		return walkFn(path, info.IsDir(), nil)
	})
}
