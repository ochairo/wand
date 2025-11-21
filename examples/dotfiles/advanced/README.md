# Advanced Dotfiles Example

This example demonstrates managing configurations for multiple development tools.

## Structure

```
advanced/
├── README.md
├── dev-dotfiles.yaml      # Wand formula
└── configs/
    ├── nvim/              # Neovim configuration
    ├── starship/          # Starship prompt
    ├── zsh/               # Zsh configuration
    ├── git/               # Git configuration
    └── wezterm/           # WezTerm terminal
```

## Installation

```bash
# Install the complete dev environment dotfiles
wand install dev-dotfiles

# Or install specific configurations
wand install dev-dotfiles --config nvim
wand install dev-dotfiles --config starship
```

## Included Configurations

### Neovim
- LSP configurations for multiple languages
- Plugin management with lazy.nvim
- Custom keymaps and settings

### Starship
- Custom prompt with Git integration
- Performance-optimized configuration
- Module customization

### Zsh
- Oh My Zsh integration
- Custom aliases and functions
- Fast syntax highlighting
- Auto-suggestions

### Git
- Global gitconfig
- Git aliases
- Diff and merge tools

### WezTerm
- Modern terminal emulator config
- Custom key bindings
- Theme configuration

## Formula Definition

See `dev-dotfiles.yaml` for the complete formula definition.

## Real-World Reference

This example is inspired by: <https://github.com/ochairo/dotfiles>
