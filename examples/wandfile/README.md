# Wandfile Examples

This directory contains example wandfiles demonstrating the format and capabilities.

## Basic Example

See [`wandfile.yaml`](./wandfile.yaml) for a complete development environment setup.

## Format

```yaml
formulas:
  - package-name              # Install latest version
  - package-name@1.2.3        # Install specific version
  - shell-plugin-name         # Install shell plugins (auto-configured)

symlinks:
  .target: source/path        # Creates ~/.target -> ~/.dotfiles/source/path
```

## Features

- **Simple formula syntax**: Just package names with optional `@version`
- **Shell plugin support**: Plugins like `zsh-syntax-highlighting` are automatically configured
- **Symlink management**: Define dotfile symlinks from `~/.dotfiles/` to `~/`
- **Version pinning**: Lock specific versions with `@version` syntax

## Usage

```bash
# Validate
wand wandfile validate examples/wandfile.yaml

# Show summary
wand wandfile show examples/wandfile.yaml

# Install (careful - this will modify your system!)
wand wandfile install examples/wandfile.yaml
```

## Shell Plugins

Shell plugins are automatically configured in your shell config (`.zshrc`, `.bashrc`, etc.):

```yaml
formulas:
  - zsh-syntax-highlighting
  - zsh-autosuggestions
  - oh-my-zsh
```

These will:
1. Clone the Git repository to `~/.wand/plugins/`
2. Add source lines to your shell config
3. Use wand markers to track managed plugins
4. Prevent duplicate source lines

## Symlinks

Symlinks assume you have your dotfiles in `~/.dotfiles/`:

```yaml
symlinks:
  .zshrc: shell/zshrc                    # ~/.zshrc -> ~/.dotfiles/shell/zshrc
  .config/nvim/init.lua: nvim/init.lua   # ~/.config/nvim/init.lua -> ~/.dotfiles/nvim/init.lua
```

The system will:
- Create parent directories automatically
- Backup existing files with `.wand-backup` extension
- Skip if symlink already points to the correct target
