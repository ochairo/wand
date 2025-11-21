package validation

import (
	"fmt"
	"net/url"
	"strings"
)

// URLValidator validates URL strings.
type URLValidator struct {
	allowedHosts []string
}

// NewURLValidator creates a new URL validator.
func NewURLValidator() *URLValidator {
	return &URLValidator{
		allowedHosts: []string{
			"github.com",
			"raw.githubusercontent.com",
		},
	}
}

// Validate checks if a URL is valid and safe.
func (v *URLValidator) Validate(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Scheme != "https" {
		return fmt.Errorf("only HTTPS URLs are allowed")
	}

	// If no host restrictions, only check HTTPS
	if len(v.allowedHosts) == 0 {
		return nil
	}

	allowed := false
	for _, host := range v.allowedHosts {
		if parsed.Host == host || strings.HasSuffix(parsed.Host, "."+host) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("URL host not in allowed list: %s", parsed.Host)
	}

	return nil
}
