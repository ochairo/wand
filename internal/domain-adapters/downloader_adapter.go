package domainadapters

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ochairo/wand/internal/domain/interfaces"
)

const (
	maxRetries    = 3
	retryDelay    = 2 * time.Second
	retryDelayMax = 10 * time.Second
)

// DownloaderAdapter implements file downloading
type DownloaderAdapter struct{}

// NewDownloaderAdapter creates a new DownloaderAdapter
func NewDownloaderAdapter() interfaces.Downloader {
	return &DownloaderAdapter{}
}

// Download downloads a file from a URL to a destination path
func (d *DownloaderAdapter) Download(url, destPath string) error {
	return d.downloadWithRetry(url, destPath, nil)
}

// DownloadWithProgress downloads a file with progress reporting
func (d *DownloaderAdapter) DownloadWithProgress(url, destPath string, progress io.Writer) error {
	return d.downloadWithRetry(url, destPath, progress)
}

// downloadWithRetry implements download with exponential backoff retry logic
func (d *DownloaderAdapter) downloadWithRetry(url, destPath string, progress io.Writer) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2s, 4s, 8s (capped at 10s)
			// Safe cast: attempt is always >= 1 at this point, uint conversion safe
			shift := attempt - 1
			delay := retryDelay * time.Duration(1<<uint(shift)) //nolint:gosec // G115: shift >= 0 guaranteed by guard
			if delay > retryDelayMax {
				delay = retryDelayMax
			}
			time.Sleep(delay)
		}

		err := d.doDownload(url, destPath, progress)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on 404 or other client errors
		if isClientError(err) {
			return err
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// doDownload performs a single download attempt
// Note: Uses http.Get with variable URL as this is by design - tool downloads from user-specified URLs
// Note: Uses os.Create with variable paths as this is by design - tool creates files at user-specified locations
func (d *DownloaderAdapter) doDownload(url, destPath string, progress io.Writer) error {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return &httpError{
			statusCode: resp.StatusCode,
			status:     resp.Status,
		}
	}

	outFile, err := os.Create(destPath) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer func() { _ = outFile.Close() }()

	var writer io.Writer = outFile
	if progress != nil {
		writer = io.MultiWriter(outFile, progress)
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// httpError represents an HTTP error response
type httpError struct {
	statusCode int
	status     string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("download failed with status %d: %s", e.statusCode, e.status)
}

// isClientError checks if error is a 4xx client error (should not retry)
func isClientError(err error) bool {
	var httpErr *httpError
	if errors.As(err, &httpErr) {
		return httpErr.statusCode >= 400 && httpErr.statusCode < 500
	}
	return false
}

// VerifyChecksum verifies the SHA256 checksum of a file
// Note: Uses http.Get with variable URL as this is by design - checksum URLs from package formulas
// Note: Uses os.Open with variable paths as this is by design - tool opens downloaded files
func (d *DownloaderAdapter) VerifyChecksum(filePath, checksumURL string) error {
	// Download checksum file
	resp, err := http.Get(checksumURL) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to download checksum: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("checksum download failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	// Read expected checksum
	checksumData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read checksum: %w", err)
	}

	// Parse checksum (format: "<hash>  <filename>" or just "<hash>")
	expectedChecksum := strings.TrimSpace(string(checksumData))
	if idx := strings.Index(expectedChecksum, " "); idx != -1 {
		expectedChecksum = expectedChecksum[:idx]
	}
	expectedChecksum = strings.ToLower(expectedChecksum)

	// Calculate actual checksum
	file, err := os.Open(filePath) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to open file for checksum: %w", err)
	}
	defer func() { _ = file.Close() }()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))

	// Compare checksums
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}
