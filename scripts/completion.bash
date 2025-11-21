#!/bin/bash

# Bash completion for wand
# Place this file in /etc/bash_completion.d/wand or source it from ~/.bashrc

_wand_completions() {
    local cur prev words cword
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    words=("${COMP_WORDS[@]}")
    cword=$COMP_CWORD

    # Main commands
    local commands="install uninstall list search info activate update outdated doctor config cache validate version help"

    # If we're at the first argument, suggest commands
    if [[ $cword -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
        return 0
    fi

    # Get list of installed packages
    local packages=$(wand list --format json 2>/dev/null | grep -o '"name":"[^"]*"' | cut -d'"' -f4 | tr '\n' ' ')

    # Handle subcommands
    case "${COMP_WORDS[1]}" in
        install)
            if [[ $cword -eq 2 ]]; then
                # Suggest available formulas
                local formulas=$(wand search --format json 2>/dev/null | grep -o '"name":"[^"]*"' | cut -d'"' -f4 | tr '\n' ' ')
                COMPREPLY=( $(compgen -W "${formulas}" -- "${cur}") )
            fi
            ;;
        uninstall|info|activate|update)
            if [[ $cword -eq 2 ]]; then
                # Suggest installed packages
                COMPREPLY=( $(compgen -W "${packages}" -- "${cur}") )
            fi
            ;;
        config)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "get set list reset" -- "${cur}") )
            fi
            ;;
        cache)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "clean clear size" -- "${cur}") )
            fi
            ;;
    esac

    return 0
}

complete -F _wand_completions wand
