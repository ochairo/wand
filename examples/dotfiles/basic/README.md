# Basic Dotfiles Example

This example shows how to manage shell configuration files (`.zshrc`, `.zshenv`) using wand.

## Structure

```
basic/
├── README.md
├── shell-dotfiles.yaml    # Wand formula
└── configs/               # Configuration files
    ├── .zshrc
    ├── .zshenv
    ├── .zsh_aliases
    └── .zsh_functions
```

## Installation

```bash
# Install the dotfiles formula
wand install shell-dotfiles

# This will:
# 1. Clone the repository (if remote)
# 2. Create symlinks: ~/.zshrc -> configs/.zshrc
# 3. Source configurations automatically
```

## Formula Definition

See `shell-dotfiles.yaml` for the formula that manages these dotfiles.

## Usage

After installation, your shell will automatically load:
- `.zshrc` - Main shell configuration
- `.zshenv` - Environment variables
- `.zsh_aliases` - Command aliases
- `.zsh_functions` - Custom shell functions

## Customization

Edit the files in `configs/` and they'll be automatically reflected in your shell (symlinked).
