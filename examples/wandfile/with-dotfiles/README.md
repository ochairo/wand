# Wandfile with Dotfiles Example

Complete example showing package installation AND dotfile management in a single Wandfile.

## Usage

```bash
wand install --wandfile ./Wandfile
```

This will:
1. Install all specified packages
2. Clone and link dotfile configurations
3. Set up the complete environment

## What's Included

### Packages
- Essential CLI tools (jq, ripgrep, fd, bat)
- Development tools (git, neovim)
- Shell enhancement (zsh, starship)

### Dotfiles
- Shell configuration (zsh, aliases, functions)
- Git configuration (aliases, settings)
- Neovim configuration (LSP, plugins)
- Starship prompt configuration

## Benefits

- **Single source of truth** for packages AND configurations
- **Reproducible environments** across machines
- **Team standardization** - share one Wandfile
- **Version control** both tools and configs

## Related Examples

- [Basic Wandfile](../basic/) - Packages only
- [Development Wandfile](../development/) - Full dev environment
- [Dotfiles Examples](../../dotfiles/) - Dotfile-only examples
