// Package errors defines error types and error codes for the domain.
package errors

import "fmt"

// ErrorCode represents error classification
type ErrorCode string

const (
	// ErrInvalidPackageName indicates a package name failed validation.
	ErrInvalidPackageName ErrorCode = "INVALID_PACKAGE_NAME"
	// ErrInvalidVersion indicates a version string failed validation.
	ErrInvalidVersion ErrorCode = "INVALID_VERSION"
	// ErrInvalidPath indicates a file path failed validation.
	ErrInvalidPath ErrorCode = "INVALID_PATH"
	// ErrInvalidURL indicates a URL failed validation.
	ErrInvalidURL ErrorCode = "INVALID_URL"
	// ErrPackageNotFound indicates a package was not found in the repository.
	ErrPackageNotFound ErrorCode = "PACKAGE_NOT_FOUND"
	// ErrPackageNotInstalled indicates a package is not currently installed.
	ErrPackageNotInstalled ErrorCode = "PACKAGE_NOT_INSTALLED"
	// ErrPackageInstalled indicates a package is already installed.
	ErrPackageInstalled ErrorCode = "PACKAGE_INSTALLED"
	// ErrVersionNotFound indicates a specific version was not found.
	ErrVersionNotFound ErrorCode = "VERSION_NOT_FOUND"
	// ErrVersionInstalled indicates a version is already installed.
	ErrVersionInstalled ErrorCode = "VERSION_INSTALLED"
	// ErrDownloadFailed indicates a download operation failed.
	ErrDownloadFailed ErrorCode = "DOWNLOAD_FAILED"
	// ErrChecksumMismatch indicates a checksum validation failed.
	ErrChecksumMismatch ErrorCode = "CHECKSUM_MISMATCH"
	// ErrExtractionFailed indicates archive extraction failed.
	ErrExtractionFailed ErrorCode = "EXTRACTION_FAILED"
	// ErrInstallationFailed indicates the installation process failed.
	ErrInstallationFailed ErrorCode = "INSTALLATION_FAILED"
	// ErrBinaryNotFound indicates a required binary was not found.
	ErrBinaryNotFound ErrorCode = "BINARY_NOT_FOUND"
	// ErrShimCreationFailed indicates shim creation failed.
	ErrShimCreationFailed ErrorCode = "SHIM_CREATION_FAILED"
	// ErrShimExecutionFailed indicates shim execution failed.
	ErrShimExecutionFailed ErrorCode = "SHIM_EXECUTION_FAILED"
	// ErrFileNotFound indicates a file was not found.
	ErrFileNotFound ErrorCode = "FILE_NOT_FOUND"
	// ErrDirNotFound indicates a directory was not found.
	ErrDirNotFound ErrorCode = "DIR_NOT_FOUND"
	// ErrPermissionDenied indicates insufficient permissions for an operation.
	ErrPermissionDenied ErrorCode = "PERMISSION_DENIED"
	// ErrDiskSpaceLow indicates insufficient disk space.
	ErrDiskSpaceLow ErrorCode = "DISK_SPACE_LOW"
	// ErrNetworkUnreachable indicates the network is unreachable.
	ErrNetworkUnreachable ErrorCode = "NETWORK_UNREACHABLE"
	// ErrHTTPError indicates an HTTP operation failed.
	ErrHTTPError ErrorCode = "HTTP_ERROR"
	// ErrTimeout indicates an operation timed out.
	ErrTimeout ErrorCode = "TIMEOUT"
	// ErrConfigMissing indicates a required configuration is missing.
	ErrConfigMissing ErrorCode = "CONFIG_MISSING"
	// ErrConfigInvalid indicates a configuration is invalid.
	ErrConfigInvalid ErrorCode = "CONFIG_INVALID"
	// ErrRegistryCorrupted indicates the registry is corrupted.
	ErrRegistryCorrupted ErrorCode = "REGISTRY_CORRUPTED"
	// ErrSystemNotSupported indicates an unsupported operating system.
	ErrSystemNotSupported ErrorCode = "SYSTEM_NOT_SUPPORTED"
	// ErrArchNotSupported indicates an unsupported CPU architecture.
	ErrArchNotSupported ErrorCode = "ARCH_NOT_SUPPORTED"
)

// WandError represents a Wand-specific error
type WandError struct {
	Code    ErrorCode
	Message string
	Details string
	Wrapped error
}

// Error implements error interface
func (e *WandError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns wrapped error
func (e *WandError) Unwrap() error {
	return e.Wrapped
}

// New creates new WandError
func New(code ErrorCode, message string) *WandError {
	return &WandError{
		Code:    code,
		Message: message,
	}
}

// NewWithDetails creates WandError with details
func NewWithDetails(code ErrorCode, message, details string) *WandError {
	return &WandError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Wrap wraps error with context
func Wrap(code ErrorCode, message string, err error) *WandError {
	return &WandError{
		Code:    code,
		Message: message,
		Wrapped: err,
	}
}

// PackageNotFound creates package not found error
func PackageNotFound(name string) *WandError {
	return NewWithDetails(ErrPackageNotFound, "Package not found", fmt.Sprintf("package: %q", name))
}

// VersionNotFound creates version not found error
func VersionNotFound(pkg, version string) *WandError {
	return NewWithDetails(ErrVersionNotFound, "Version not found", fmt.Sprintf("package: %q, version: %q", pkg, version))
}

// ChecksumMismatch creates checksum error
func ChecksumMismatch(expected, actual string) *WandError {
	return NewWithDetails(ErrChecksumMismatch, "Checksum mismatch", fmt.Sprintf("expected: %s, got: %s", expected, actual))
}

// InsufficientDiskSpace creates disk space error
func InsufficientDiskSpace(required, available int64) *WandError {
	return NewWithDetails(ErrDiskSpaceLow, "Insufficient disk space", fmt.Sprintf("required: %d bytes, available: %d bytes", required, available))
}
