# Contributing to Wand

Thank you for your interest in contributing! See [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md) for our community standards.

## How to Contribute

### Report Bugs

1. Check [GitHub Issues](https://github.com/ochairo/wand/issues)
2. Use the bug report template
3. Include: version, OS/arch, steps to reproduce, expected vs actual behavior

### Suggest Features

1. Check [GitHub Issues](https://github.com/ochairo/wand/issues) and [Discussions](https://github.com/ochairo/wand/discussions)
2. Use the feature request template
3. Explain the use case and benefits

### Contribute Code

See [DEVELOPMENT.md](DEVELOPMENT.md) for architecture, setup, and coding standards.

**Quick Start**:
```bash
git clone https://github.com/ochairo/wand.git
cd wand
go mod download
make build
make test
```

**Pull Request Process**:
1. Fork, create feature branch: `git checkout -b feature/name`
2. Follow code standards in [DEVELOPMENT.md](DEVELOPMENT.md)
3. Write tests for new features
4. Run: `go test ./...`, `go fmt ./...`
5. Commit with [Conventional Commits](https://www.conventionalcommits.org/) format
6. Push and create PR

### Contribute Formulas

See [DEVELOPMENT.md](DEVELOPMENT.md) under "Contributing Formulas" for creating package definitions.

## Commit Format

```bash
<type>(<scope>): <subject>

# Examples:
feat(installer): add custom mirror support
fix(shim): resolve version lookup in nested dirs
docs: clarify installation steps
```

**Types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

**Scopes**: `cli`, `installer`, `shim`, `formula`, `registry`

## Resources

- [Development Guide](DEVELOPMENT.md) - Architecture, setup, formulas
- [Error Codes](ERROR_CODES.md) - Error reference for debugging
- [Security Policy](SECURITY.md) - Security practices
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
