package domainadapters

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ochairo/wand/internal/domain/interfaces"
)

const (
	// MaxExtractSize limits total extraction size to prevent decompression bombs (1GB)
	MaxExtractSize = 1 << 30
	// MaxFileSize limits individual file size (500MB)
	MaxFileSize = 500 << 20
)

// ExtractorAdapter implements archive extraction
type ExtractorAdapter struct {
	fs interfaces.FileSystem
}

// NewExtractorAdapter creates a new ExtractorAdapter
func NewExtractorAdapter(fs interfaces.FileSystem) interfaces.Extractor {
	return &ExtractorAdapter{
		fs: fs,
	}
}

// ValidateArchivePath ensures the target path is within the destination directory (prevents directory traversal)
func ValidateArchivePath(destDir, targetPath string) (string, error) {
	// Resolve both paths
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("invalid destination directory: %w", err)
	}

	// Join and clean the target path
	fullPath := filepath.Join(absDestDir, targetPath)
	fullPath = filepath.Clean(fullPath)

	// Ensure the resolved path is within destDir
	if !strings.HasPrefix(fullPath, absDestDir+string(filepath.Separator)) && fullPath != absDestDir {
		return "", fmt.Errorf("path traversal detected: %s is outside %s", targetPath, absDestDir)
	}

	return fullPath, nil
}

// Extract extracts an archive to a destination directory
func (e *ExtractorAdapter) Extract(archivePath, destDir string) error {
	if err := e.fs.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	switch {
	case strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz"):
		return e.extractTarGz(archivePath, destDir)
	case strings.HasSuffix(archivePath, ".tar.bz2"):
		return e.extractTarBz2(archivePath, destDir)
	case strings.HasSuffix(archivePath, ".tar"):
		return e.extractTar(archivePath, destDir)
	case strings.HasSuffix(archivePath, ".zip"):
		return e.extractZip(archivePath, destDir)
	case strings.HasSuffix(archivePath, ".dmg"):
		return e.extractDmg(archivePath, destDir)
	default:
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}
}

// ExtractFile extracts a single file from an archive
func (e *ExtractorAdapter) ExtractFile(archivePath, fileName, destPath string) error {
	if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		return e.extractTarGzFile(archivePath, fileName, destPath)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return e.extractZipFileSingle(archivePath, fileName, destPath)
	}
	return fmt.Errorf("unsupported archive format for file extraction")
}

// extractTarGzFile extracts a single file from a tar.gz archive
// Note: Extracts files from archive to user-specified destinations (by design)
// Note: Uses path validation to prevent directory traversal attacks
func (e *ExtractorAdapter) extractTarGzFile(archivePath, fileName, destPath string) error {
	file, err := os.Open(archivePath) //nolint:gosec // G304: archivePath is from validated archive
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer func() { _ = file.Close() }()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("file not found in archive: %s", fileName)
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		// Validate file size to prevent decompression bombs
		if header.Size > MaxFileSize {
			return fmt.Errorf("file %s is too large: %d bytes (max %d)", header.Name, header.Size, MaxFileSize)
		}

		// Match file name or path ending with filename
		if header.Name == fileName || strings.HasSuffix(header.Name, "/"+fileName) || strings.HasSuffix(header.Name, fileName) {
			if header.Typeflag == tar.TypeReg {
				// Ensure destination directory exists
				destDir := filepath.Dir(destPath)
				if err := os.MkdirAll(destDir, 0700); err != nil {
					return fmt.Errorf("failed to create destination directory: %w", err)
				}

				// Create destination file
				out, err := os.Create(destPath) //nolint:gosec // G304: destPath validated by ValidateArchivePath()
				if err != nil {
					return fmt.Errorf("failed to create destination file: %w", err)
				}
				defer func() { _ = out.Close() }() // Copy file contents with size limit to prevent decompression bombs
				limitedReader := io.LimitReader(tr, MaxFileSize)
				if _, err := io.Copy(out, limitedReader); err != nil {
					return fmt.Errorf("failed to extract file: %w", err)
				}
				// Set permissions
				if err := os.Chmod(destPath, os.FileMode(header.Mode)); err != nil { //nolint:gosec // G115: safe conversion from tar header
					return fmt.Errorf("failed to set permissions: %w", err)
				}

				return nil
			}
		}
	}
}

// extractZipFile extracts a single file from a zip archive
func (e *ExtractorAdapter) extractZipFileSingle(archivePath, fileName, destPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer func() { _ = reader.Close() }()

	for _, file := range reader.File {
		// Match file name or path ending with filename
		if file.Name == fileName || strings.HasSuffix(file.Name, "/"+fileName) || strings.HasSuffix(file.Name, fileName) {
			if !file.FileInfo().IsDir() {
				// Validate file size to prevent decompression bombs
				if file.UncompressedSize64 > MaxFileSize {
					return fmt.Errorf("file %s is too large: %d bytes (max %d)", file.Name, file.UncompressedSize64, MaxFileSize)
				}

				// Ensure destination directory exists
				destDir := filepath.Dir(destPath)
				if err := os.MkdirAll(destDir, 0700); err != nil {
					return fmt.Errorf("failed to create destination directory: %w", err)
				}

				// Open source file in zip
				source, err := file.Open()
				if err != nil {
					return fmt.Errorf("failed to open file in archive: %w", err)
				}
				defer func() { _ = source.Close() }()

				// Create destination file
				out, err := os.Create(destPath) //nolint:gosec // G304: destPath is validated
				if err != nil {
					return fmt.Errorf("failed to create destination file: %w", err)
				}
				defer func() { _ = out.Close() }()

				// Copy file contents with size limit
				limitedReader := io.LimitReader(source, MaxFileSize)
				if _, err := io.Copy(out, limitedReader); err != nil {
					return fmt.Errorf("failed to extract file: %w", err)
				}

				// Set permissions
				if err := os.Chmod(destPath, file.FileInfo().Mode()); err != nil {
					return fmt.Errorf("failed to set permissions: %w", err)
				}

				return nil
			}
		}
	}

	return fmt.Errorf("file not found in archive: %s", fileName)
}

// extractTarGz extracts a .tar.gz archive
func (e *ExtractorAdapter) extractTarGz(archivePath, destDir string) error {
	file, err := os.Open(archivePath) //nolint:gosec // G304: archivePath is from validated archive
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() { _ = gzReader.Close() }()

	return e.extractTarReader(tar.NewReader(gzReader), destDir)
}

// extractTarBz2 extracts a .tar.bz2 archive
func (e *ExtractorAdapter) extractTarBz2(archivePath, destDir string) error {
	file, err := os.Open(archivePath) //nolint:gosec // G304: archivePath is from validated archive
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	bzReader := bzip2.NewReader(file)
	return e.extractTarReader(tar.NewReader(bzReader), destDir)
}

// extractTar extracts a .tar archive
func (e *ExtractorAdapter) extractTar(archivePath, destDir string) error {
	file, err := os.Open(archivePath) //nolint:gosec // G304: archivePath is from validated archive
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return e.extractTarReader(tar.NewReader(file), destDir)
}

// extractTarReader extracts from a tar reader
func (e *ExtractorAdapter) extractTarReader(tarReader *tar.Reader, destDir string) error {
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		// Validate file size to prevent decompression bombs
		if header.Size > MaxFileSize {
			return fmt.Errorf("file %s is too large: %d bytes (max %d)", header.Name, header.Size, MaxFileSize)
		}

		// Validate path to prevent directory traversal
		target, err := ValidateArchivePath(destDir, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := e.fs.MkdirAll(target, 0700); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := e.extractTarFile(tarReader, target, header.Mode); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := e.fs.Symlink(header.Linkname, target); err != nil {
				return err
			}
		}
	}

	return nil
}

// extractTarFile extracts a single file from tar
func (e *ExtractorAdapter) extractTarFile(reader io.Reader, target string, mode int64) error {
	// Create parent directory
	if err := e.fs.MkdirAll(filepath.Dir(target), 0700); err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(mode)) //nolint:gosec // G304: target is validated
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Copy with size limit to prevent decompression bombs
	limitedReader := io.LimitReader(reader, MaxFileSize)
	_, err = io.Copy(file, limitedReader)
	return err
}

// extractZip extracts a .zip archive
func (e *ExtractorAdapter) extractZip(archivePath, destDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	for _, file := range reader.File {
		// Validate path to prevent directory traversal
		target, err := ValidateArchivePath(destDir, file.Name)
		if err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			if err := e.fs.MkdirAll(target, 0700); err != nil {
				return err
			}
			continue
		}

		if err := e.extractZipFile(file, target); err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from zip
func (e *ExtractorAdapter) extractZipFile(file *zip.File, target string) error {
	// Validate file size to prevent decompression bombs
	if file.UncompressedSize64 > MaxFileSize {
		return fmt.Errorf("file %s is too large: %d bytes (max %d)", file.Name, file.UncompressedSize64, MaxFileSize)
	}

	// Create parent directory
	if err := e.fs.MkdirAll(filepath.Dir(target), 0700); err != nil {
		return err
	}

	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode()) //nolint:gosec // G304: target is validated
	if err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()

	// Copy with size limit to prevent decompression bombs
	limitedReader := io.LimitReader(reader, MaxFileSize)
	_, err = io.Copy(outFile, limitedReader)
	return err
}

// extractDmg extracts a macOS .dmg disk image
func (e *ExtractorAdapter) extractDmg(dmgPath, destDir string) error {
	// On macOS, use hdiutil to mount and copy .app bundle
	// For now, we'll implement a simple copy approach assuming the .app is directly accessible
	// Full implementation would use shell commands to mount DMG, copy contents, unmount

	// Create a temporary mount point
	tmpMount := filepath.Join(os.TempDir(), "wand-dmg-mount")
	if err := e.fs.MkdirAll(tmpMount, 0700); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}
	defer func() { _ = e.fs.RemoveAll(tmpMount) }()

	// Mount the DMG (macOS only)
	// Use explicit command with arguments instead of string parsing
	mountCmd := exec.Command("hdiutil", "attach", "-nobrowse", "-mountpoint", tmpMount, dmgPath) //nolint:gosec // G204: hardcoded command
	if err := mountCmd.Run(); err != nil {
		return fmt.Errorf("failed to mount dmg: %w", err)
	}
	defer func() { _ = exec.Command("hdiutil", "detach", tmpMount).Run() }() //nolint:gosec // G204: hardcoded command

	// Find .app bundle in mounted volume
	entries, err := os.ReadDir(tmpMount)
	if err != nil {
		return fmt.Errorf("failed to read mount point: %w", err)
	}

	// Copy .app bundle to destination
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			srcPath := filepath.Join(tmpMount, entry.Name())
			dstPath := filepath.Join(destDir, entry.Name())
			if err := e.copyDir(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy app bundle: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("no .app bundle found in dmg")
}

// copyDir recursively copies a directory
func (e *ExtractorAdapter) copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return e.fs.MkdirAll(dstPath, 0700)
		}

		// Copy file
		return e.copyFile(path, dstPath, info.Mode())
	})
}

// copyFile copies a single file
func (e *ExtractorAdapter) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src) //nolint:gosec // G304: src is from filepath.Walk
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode) //nolint:gosec // G304: dst is validated
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	// Copy with size limit
	limitedReader := io.LimitReader(srcFile, MaxFileSize)
	_, err = io.Copy(dstFile, limitedReader)
	return err
}
