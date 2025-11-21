# Neovim Configuration Example

This directory contains a minimal Neovim configuration example for reference.

## Structure

```
nvim/
├── init.lua                    # Main entry point
├── lua/
│   └── user/
│       ├── settings.lua        # Editor settings
│       ├── keymaps.lua         # Key mappings
│       └── plugins/
│           ├── init.lua        # Plugin manager (lazy.nvim)
│           ├── lsp.lua         # LSP configuration
│           └── telescope.lua   # Fuzzy finder
```

## Features

- **LSP Support**: Language server protocol for code intelligence
- **Fuzzy Finder**: Telescope for file/text search
- **Plugin Management**: lazy.nvim for fast plugin loading
- **Custom Keymaps**: Intuitive key bindings

## Installation via Wand

This configuration is automatically linked when you install `dev-dotfiles`:

```bash
wand install dev-dotfiles
```

## Standalone Installation

```bash
# Link this directory to ~/.config/nvim
ln -s /path/to/this/nvim ~/.config/nvim

# Launch nvim - plugins will auto-install
nvim
```

## Reference

For a complete, production-ready configuration, see:
<https://github.com/ochairo/dotfiles/tree/main/src/configs/nvim>
