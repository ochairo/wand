# Dotfiles Management Examples

This directory contains examples of using wand to manage dotfiles configurations.

## Overview

Wand can manage dotfiles by:
- Installing configuration files from Git repositories
- Creating symlinks to dotfile configurations
- Managing multiple dotfile sources for different tools
- Version-controlling tool configurations

## Example Structure

```
~/.config/
├── nvim/          -> dotfiles/nvim
├── starship/      -> dotfiles/starship
├── zellij/        -> dotfiles/zellij
└── wezterm/       -> dotfiles/wezterm
```

## Examples

### 1. Basic Dotfiles Installation

See [basic/](./basic/) for a simple example of managing shell configurations.

### 2. Advanced Multi-Tool Setup

See [advanced/](./advanced/) for managing multiple tool configurations (nvim, zsh, starship, etc.).

## Formula Examples

Dotfile formulas define:
- Repository URL
- Target installation directory
- Symlink mappings
- Post-install scripts

See the example formulas in each subdirectory for reference.

## Related Documentation

- [Dotfile formulas in wand](../../formulas/)
- [Real-world dotfiles example](https://github.com/ochairo/dotfiles)
