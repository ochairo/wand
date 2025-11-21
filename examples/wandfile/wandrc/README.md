# Project-Specific .wandrc Example

Project-level package version overrides.

## Usage

Place a `.wandrc` file in your project root to override system-wide versions defined in `Wandfile`.

## Example Scenarios

### Frontend Project

```yaml
# frontend/.wandrc
packages:
  - name: node
    version: "18.17.0"  # This project uses Node 18
```

### Python Data Science Project

```yaml
# ml-project/.wandrc
packages:
  - name: python
    version: "3.10.0"  # Compatible with TensorFlow
```

### Legacy Project

```yaml
# legacy-app/.wandrc
packages:
  - name: node
    version: "16.20.0"  # Old project, not upgraded yet
  - name: npm
    version: "8.19.0"
```

## How It Works

1. System-wide `Wandfile` defines default versions
2. Project `.wandrc` overrides specific packages
3. Wand automatically switches versions when you `cd` into the project

```bash
cd ~/projects/frontend    # Auto-switches to Node 18
node --version            # v18.17.0

cd ~/projects/backend     # Auto-switches to Node 20
node --version            # v20.10.0
```

## Benefits

- **Per-project versions** without conflicts
- **Team consistency** via committed `.wandrc`
- **Automatic switching** based on directory
- **No manual version management**

## Related

- [Wandfile examples](../)
- [Getting Started Guide](../../docs/GETTING_STARTED.md)
