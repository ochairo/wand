package errors

import (
	stderrors "errors"
	"testing"
)

func TestWandErrorFormat(t *testing.T) {
	tests := []struct {
		name     string
		err      *WandError
		expected string
	}{
		{
			name: "error with message",
			err: &WandError{
				Code:    ErrPackageNotFound,
				Message: "package not found",
			},
			expected: "[PACKAGE_NOT_FOUND] package not found",
		},
		{
			name: "error with details",
			err: &WandError{
				Code:    ErrInvalidPackageName,
				Message: "invalid package name",
				Details: "package name contains invalid characters",
			},
			expected: "[INVALID_PACKAGE_NAME] invalid package name: package name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("got %q, want %q", tt.err.Error(), tt.expected)
			}
		})
	}
}

func TestErrorHelpers(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func() *WandError
		wantCode ErrorCode
	}{
		{
			name:     "PackageNotFound",
			errFunc:  func() *WandError { return PackageNotFound("test-pkg") },
			wantCode: ErrPackageNotFound,
		},
		{
			name:     "VersionNotFound",
			errFunc:  func() *WandError { return VersionNotFound("test-pkg", "1.0.0") },
			wantCode: ErrVersionNotFound,
		},
		{
			name:     "ChecksumMismatch",
			errFunc:  func() *WandError { return ChecksumMismatch("abc", "def") },
			wantCode: ErrChecksumMismatch,
		},
		{
			name:     "InsufficientDiskSpace",
			errFunc:  func() *WandError { return InsufficientDiskSpace(1024, 2048) },
			wantCode: ErrDiskSpaceLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err.Code != tt.wantCode {
				t.Errorf("got code %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	underlying := stderrors.New("underlying error")
	wrapped := &WandError{
		Code:    ErrNetworkUnreachable,
		Message: "network request failed",
		Wrapped: underlying,
	}

	if !stderrors.Is(wrapped, underlying) {
		t.Error("wrapped error should be findable with errors.Is")
	}

	if stderrors.Unwrap(wrapped) != underlying { //nolint:errorlint
		t.Error("errors.Unwrap should return underlying error")
	}
}

func TestErrorConstructors(t *testing.T) {
	// Test New
	err := New(ErrPackageNotFound, "test message")
	if err.Code != ErrPackageNotFound || err.Message != "test message" {
		t.Errorf("New() failed: got %v, %s", err.Code, err.Message)
	}

	// Test NewWithDetails
	err = NewWithDetails(ErrInvalidVersion, "test message", "test details")
	if err.Code != ErrInvalidVersion || err.Details != "test details" {
		t.Errorf("NewWithDetails() failed: got %v, %s, %s", err.Code, err.Message, err.Details)
	}

	// Test Wrap
	underlying := stderrors.New("underlying")
	err = Wrap(ErrDownloadFailed, "download failed", underlying)
	if err.Code != ErrDownloadFailed || err.Wrapped != underlying { //nolint:errorlint
		t.Errorf("Wrap() failed: got %v, wrapped=%v", err.Code, err.Wrapped)
	}
}

func TestErrorCodes(t *testing.T) {
	// Verify all error codes are defined and non-empty
	codes := []ErrorCode{
		ErrInvalidPackageName,
		ErrInvalidVersion,
		ErrInvalidPath,
		ErrInvalidURL,
		ErrPackageNotFound,
		ErrPackageNotInstalled,
		ErrPackageInstalled,
		ErrVersionNotFound,
		ErrVersionInstalled,
		ErrDownloadFailed,
		ErrChecksumMismatch,
		ErrExtractionFailed,
		ErrInstallationFailed,
		ErrBinaryNotFound,
		ErrShimCreationFailed,
		ErrShimExecutionFailed,
		ErrFileNotFound,
		ErrDirNotFound,
		ErrPermissionDenied,
		ErrDiskSpaceLow,
		ErrNetworkUnreachable,
		ErrHTTPError,
		ErrTimeout,
		ErrConfigMissing,
		ErrConfigInvalid,
		ErrRegistryCorrupted,
		ErrSystemNotSupported,
		ErrArchNotSupported,
	}

	for _, code := range codes {
		if code == "" {
			t.Error("error code should not be empty")
		}
	}
}
