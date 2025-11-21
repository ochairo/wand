# Installation

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/ochairo/wand/main/install.sh | sh
```

The installer will detect your OS and architecture, download the binary, create `~/.wand/shims`, and show setup instructions.

## Setup PATH

Add to your shell config (`~/.zshrc`, `~/.bashrc`, or `~/.config/fish/config.fish`):

```bash
export PATH="$HOME/.wand/shims:$PATH"
```

Reload your shell:
```bash
source ~/.zshrc
```

## Verify

```bash
wand --version
wand help
```

## Uninstall

```bash
rm /usr/local/bin/wand
rm -rf ~/.wand
```

Then remove the PATH line from your shell config.

## Manual Installation

Visit [Releases](https://github.com/ochairo/wand/releases), download for your platform, extract, and move to `/usr/local/bin`.

## Shell Completion (Optional)

```bash
# Zsh
mkdir -p ~/.zsh/completions
cp scripts/completion.zsh ~/.zsh/completions/_wand
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
compinit
```

See [Getting Started](GETTING_STARTED.md) next.
