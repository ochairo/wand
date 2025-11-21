#compdef wand

# Zsh completion for wand

_wand_commands() {
    local -a commands
    commands=(
        'install:Install a package'
        'uninstall:Remove an installed package'
        'list:List all installed packages'
        'search:Search for available formulas'
        'info:Show package information'
        'activate:Switch to a different version'
        'update:Update packages to latest version'
        'outdated:Show packages with available updates'
        'doctor:Check system health'
        'config:Manage configuration'
        'cache:Manage package cache'
        'validate:Validate wandfile or formula'
        'version:Show version information'
        'help:Show help'
    )
    _describe 'command' commands
}

_wand_packages() {
    local -a packages
    packages=($(wand list --format json 2>/dev/null | grep -o '"name":"[^"]*"' | cut -d'"' -f4))
    _values 'installed packages' "${packages[@]}"
}

_wand_formulas() {
    local -a formulas
    formulas=($(wand search --format json 2>/dev/null | grep -o '"name":"[^"]*"' | cut -d'"' -f4))
    _values 'available formulas' "${formulas[@]}"
}

_wand_config_commands() {
    local -a config_commands
    config_commands=(
        'get:Get a configuration value'
        'set:Set a configuration value'
        'list:List all configuration values'
        'reset:Reset to defaults'
    )
    _describe 'config command' config_commands
}

_wand_cache_commands() {
    local -a cache_commands
    cache_commands=(
        'clean:Clean cache for a package'
        'clear:Clear entire cache'
        'size:Show cache size'
    )
    _describe 'cache command' cache_commands
}

_wand() {
    local curcontext="$curcontext" state line
    local -a args

    args=(
        '1: :->command'
        '*: :->args'
    )

    _arguments "$args[@]"

    case $state in
        command)
            _wand_commands
            ;;
        args)
            case "${words[2]}" in
                install)
                    _wand_formulas
                    ;;
                uninstall|info|activate|update)
                    _wand_packages
                    ;;
                config)
                    if [[ $CURRENT -eq 3 ]]; then
                        _wand_config_commands
                    fi
                    ;;
                cache)
                    if [[ $CURRENT -eq 3 ]]; then
                        _wand_cache_commands
                    fi
                    ;;
            esac
            ;;
    esac
}

_wand "$@"
