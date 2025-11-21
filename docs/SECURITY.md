# Security Policy

Wand v1.0.0 implements security-first design with cryptographic verification, secure defaults, and clean architecture to prevent supply chain attacks.

## Security Features

### Package Integrity
- SHA256 checksums verify all downloads automatically
- Failed verification prevents installation
- Both binary and formula checksums validated

### Secure Defaults
- HTTPS enforced (no HTTP fallback)
- Certificate validation enabled
- Checksums required by default
- Strict file permissions (700 for ~/.wand)

### Architecture
- Clean Architecture prevents external dependency exposure
- Dependency isolation prevents supply chain attacks
- Minimal external dependencies (core: standard library only)

## Best Practices for Users

```bash
# Keep Wand updated
wand check-updates
wand update

# Verify your installation
wand verify
wand doctor

# Use HTTPS only
export WAND_HTTPS_ONLY=true

# Secure file permissions
chmod 700 ~/.wand
chmod 600 ~/.wand/.registry
chmod 600 .wandrc
```

## For System Administrators

```bash
# Centralized deployment (Ansible, Puppet, Chef recommended)
# Use configuration management for standardization

# Audit logging
export WAND_LOG_LEVEL=debug
export WAND_LOG_FORMAT=json

# User isolation (per-user installations, no shared state)
# Recommended for multi-user systems

# Access control
# Create restricted service account if needed
sudo useradd -r -s /bin/false wand-service
```

## Reporting Security Issues

**DO NOT report security vulnerabilities in public GitHub issues.**

Email: `security@ochairo.com` with:
- Affected version(s)
- Steps to reproduce
- Impact assessment
- Suggested fix (if available)

Response time: 24 hours acknowledgment, fix released as soon as possible.

## Dependency Security

Wand uses minimal dependencies:
- Standard library: `crypto/sha256`, `net/http`, `encoding/json`, `os/exec`
- Internal: `github.com/ochairo/go-utils` (maintained by team)
- External frameworks: `cobra` (CLI), `gopkg.in/yaml.v3` (YAML parsing)

All dependencies:
- Pinned to specific versions in `go.mod`
- Regularly audited: `nancy sleuth`
- Source-reviewed before inclusion
- Public repositories only

## Compliance

- OWASP Top 10 adherence
- CERT secure coding guidelines
- CWE/SANS 25 best practices

See [DEPLOYMENT.md](DEPLOYMENT.md) for production hardening.
See [ERROR_CODES.md](ERROR_CODES.md) for error handling and debugging.

---

**Security Policy Version**: 1.0.0
**Last Updated**: 2024
For questions: contact the Wand team
