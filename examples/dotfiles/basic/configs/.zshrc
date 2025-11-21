#!/usr/bin/env zsh
# Basic shell configuration for macOS and Linux

# History configuration
HISTFILE=~/.zsh_history
HISTSIZE=10000
SAVEHIST=10000
setopt HIST_IGNORE_DUPS SHARE_HISTORY

# Load aliases and functions
[[ -f ~/.zsh_aliases ]] && source ~/.zsh_aliases
[[ -f ~/.zsh_functions ]] && source ~/.zsh_functions

# Basic PATH
export PATH="/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"

# Editor
export EDITOR="vim"
export VISUAL="vim"

# Simple prompt
PROMPT='%n@%m:%~$ '
