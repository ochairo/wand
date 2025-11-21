#!/bin/bash

# Wand Shell Integration Setup Script
# This script helps you set up shell completions and other integrations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect shell
detect_shell() {
    if [ -n "$BASH_VERSION" ]; then
        echo "bash"
    elif [ -n "$ZSH_VERSION" ]; then
        echo "zsh"
    else
        echo ""
    fi
}

# Install bash completion
install_bash_completion() {
    echo "Setting up Bash completion..."
    
    local completion_dir="/etc/bash_completion.d"
    
    if [ ! -d "$completion_dir" ]; then
        completion_dir="$HOME/.bash_completion.d"
        mkdir -p "$completion_dir"
    fi
    
    if [ -f "scripts/completion.bash" ]; then
        cp scripts/completion.bash "$completion_dir/wand"
        echo -e "${GREEN}✓${NC} Bash completion installed to $completion_dir/wand"
        echo "  You may need to reload your shell: source ~/.bashrc"
    else
        echo -e "${RED}✗${NC} completion.bash not found"
        return 1
    fi
}

# Install zsh completion
install_zsh_completion() {
    echo "Setting up Zsh completion..."
    
    local completion_dir="${fpath[1]}"
    
    if [ -z "$completion_dir" ] || [ ! -d "$completion_dir" ]; then
        completion_dir="$HOME/.zsh/completions"
        mkdir -p "$completion_dir"
    fi
    
    if [ -f "scripts/completion.zsh" ]; then
        cp scripts/completion.zsh "$completion_dir/_wand"
        echo -e "${GREEN}✓${NC} Zsh completion installed to $completion_dir/_wand"
        echo "  You may need to reload your shell: source ~/.zshrc"
    else
        echo -e "${RED}✗${NC} completion.zsh not found"
        return 1
    fi
}

# Main setup
main() {
    echo "Wand Shell Integration Setup"
    echo "============================="
    echo ""
    
    local current_shell
    current_shell=$(detect_shell)
    
    if [ -z "$current_shell" ]; then
        echo -e "${YELLOW}?${NC} Could not detect shell. Please specify:"
        echo "  ./scripts/setup-shell.sh bash"
        echo "  ./scripts/setup-shell.sh zsh"
        exit 1
    fi
    
    echo "Current shell: $current_shell"
    echo ""
    
    # Allow override
    if [ -n "$1" ]; then
        current_shell="$1"
    fi
    
    case "$current_shell" in
        bash)
            install_bash_completion
            ;;
        zsh)
            install_zsh_completion
            ;;
        *)
            echo -e "${RED}✗${NC} Unsupported shell: $current_shell"
            exit 1
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}Setup complete!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Reload your shell configuration"
    echo "2. Test completion by typing: wand <TAB>"
}

main "$@"
