#!/usr/bin/env bash

set -euo pipefail

REPO="arya237/Gyscope"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="gyscope"
PATH_LINE='export PATH="$HOME/.local/bin:$PATH"'

echo "Installing Gyscope..."
echo

OS="$(uname -s)"
ARCH="$(uname -m)"

if [[ "$OS" != "Linux" ]]; then
    echo "Error: Gyscope currently supports Linux only."
    exit 1
fi

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected architecture: $ARCH"

API_URL="https://api.github.com/repos/${REPO}/releases/latest"

echo "Finding latest release..."

ASSET_URL="$(
    curl -fsSL \
        -H "Accept: application/vnd.github+json" \
        "$API_URL" |
    grep -o '"browser_download_url": "[^"]*"' |
    grep "gyscope-linux-${ARCH}" |
    head -n 1 |
    cut -d '"' -f 4
)"

if [[ -z "$ASSET_URL" ]]; then
    echo "Error: Could not find a Gyscope binary for ${ARCH}."
    exit 1
fi

mkdir -p "$INSTALL_DIR"

TEMP_FILE="$(mktemp)"

trap 'rm -f "$TEMP_FILE"' EXIT

echo "Downloading Gyscope..."

curl -fL "$ASSET_URL" -o "$TEMP_FILE"

echo "Installing Gyscope..."

install -m 755 "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"

echo
echo "✓ Gyscope installed successfully."
echo
echo "Installed to:"
echo "  ${INSTALL_DIR}/${BINARY_NAME}"
echo

if [[ ":${PATH}:" != *":${INSTALL_DIR}:"* ]]; then
    echo "${INSTALL_DIR} is not in your PATH."

    if [[ -f "${HOME}/.bashrc" ]]; then
        if ! grep -Fqx "$PATH_LINE" "${HOME}/.bashrc"; then
            echo "Adding ${INSTALL_DIR} to ~/.bashrc..."
            printf '\n%s\n' "$PATH_LINE" >> "${HOME}/.bashrc"
            echo "✓ PATH configuration added to ~/.bashrc."
        else
            echo "PATH configuration already exists in ~/.bashrc."
        fi
    else
        echo "Warning: ~/.bashrc was not found."
        echo
        echo "Add this line to your shell configuration:"
        echo
        echo "  ${PATH_LINE}"
    fi

    echo
    echo "Restart your shell or run:"
    echo
    echo "  source ~/.bashrc"
else
    echo "${INSTALL_DIR} is already in your PATH."
    echo
    echo "Run:"
    echo
    echo "  gyscope"
fi