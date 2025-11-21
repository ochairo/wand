# Team Wandfile Example

Standardized tooling configuration for development teams.

## Purpose

This Wandfile ensures all team members have:
- Same tool versions for consistency
- Required development dependencies
- Standardized CLI utilities
- Matching build environments

## Usage

### For New Team Members

```bash
# Clone the repo
git clone https://github.com/yourteam/project.git
cd project

# Install all team tools
wand install --wandfile ./Wandfile

# Verify installation
wand doctor
```

### For Existing Team Members

```bash
# Update to latest team standards
git pull
wand install --wandfile ./Wandfile

# Check for outdated packages
wand outdated
```

## Version Strategy

- **Language Runtimes**: Pinned to specific versions (consistency)
- **CLI Tools**: Latest versions (get improvements)
- **Build Tools**: Pinned to tested versions (stability)

## Project-Specific Overrides

Individual projects can override versions with `.wandrc`:

```yaml
# frontend/.wandrc
packages:
  - name: node
    version: "18.0.0"  # Frontend uses Node 18
```

```yaml
# backend/.wandrc
packages:
  - name: node
    version: "20.10.0"  # Backend uses Node 20
```

## CI/CD Integration

Use in GitHub Actions:

```yaml
- name: Install team tools
  run: |
    curl -sSL https://install.wand.sh | sh
    wand install --wandfile ./Wandfile
```

## Maintenance

Update this file when:
- Adopting new tools
- Upgrading language versions
- Deprecating old tools

Communicate changes in team standup or Slack.
