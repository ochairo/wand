---
name: Formula Submission
about: Submit a new formula to Wand
title: '[FORMULA] Add <package-name>'
labels: formula
assignees: ''
---

## Package Information
- **Name:**
- **Version:**
- **Homepage:**
- **License:**
- **Description:**

## Checklist
- [ ] Formula file follows naming convention (`<package>.yaml`)
- [ ] All required fields are present (name, version, description, homepage, license, releases)
- [ ] Package name is lowercase with hyphens only
- [ ] Version follows semver format
- [ ] All URLs use HTTPS
- [ ] Release URLs point to `github.com/ochairo/potions/releases`
- [ ] Binaries uploaded to potions repository
- [ ] Tested installation on at least one platform
- [ ] Formula validation passes locally
- [ ] Read [Formula Submission Guidelines](../docs/FORMULA_SUBMISSION.md)

## Platforms
Which platforms are supported?
- [ ] darwin-x86_64 (macOS Intel)
- [ ] darwin-arm64 (macOS Apple Silicon)
- [ ] linux-amd64 (Linux x86-64)
- [ ] linux-arm64 (Linux ARM64)

## Testing
Describe how you tested the formula:
```bash
# Example:
wand install <package>
<package> --version
wand uninstall <package>
```

## Additional Notes
Any additional context or notes about this formula.
