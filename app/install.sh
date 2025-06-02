#!/usr/bin/env bash
set -euo pipefail 

REPO_BASE_URL="https://codeberg.org/Zusty/GCT/releases/download/latest"
INSTALL_DIR="${HOME}/.local/bin"
BIN_NAME="gct" 

info() {
    echo "[INFO] $1"
}
error() {
    echo "[ERROR] $1" >&2
    exit 1
}
warn() {
     echo "[WARN] $1" >&2
}
require_util() {
    command -v "$1" >/dev/null 2>&1 || error "'$1' command is required but not found. Please install it."
}

require_util "curl"
require_util "uname"
require_util "chmod"
require_util "mkdir"
require_util "tar"
require_util "xz"

os=""
arch=""
case "$(uname -s)" in
    Linux*)  os="linux";;
    Darwin*) os="darwin";;
    MINGW*|CYGWIN*|MSYS*)
             error "Windows detected via MingW/MSYS/Cygwin. Please use the 'install.ps1' script directly in PowerShell." ;;
    *)       error "Unsupported Operating System: $(uname -s)" ;;
esac

case "$(uname -m)" in
    x86_64|amd64) arch="amd64" ;; 
    arm64|aarch64) arch="arm64" ;; 
    *)          error "Unsupported Architecture: $(uname -m)" ;;
esac

TARGET_BIN="gct-${os}-${arch}.tar.xz"
DOWNLOAD_URL="${REPO_BASE_URL}/${TARGET_BIN}"
INSTALL_PATH="${INSTALL_DIR}/${BIN_NAME}"
TEMP_ARCHIVE="/tmp/${TARGET_BIN}"
TEMP_EXTRACT_DIR="/tmp/gct_extract_${os}_${arch}"

info "Installing/Updating GCT for ${os} (${arch})..."
info "Target: ${INSTALL_PATH}"

if [ ! -d "$INSTALL_DIR" ]; then
    info "Creating installation directory: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR" || error "Failed to create directory: $INSTALL_DIR"
fi

info "Downloading GCT from: ${DOWNLOAD_URL}"
if curl --fail --location --progress-bar --output "$TEMP_ARCHIVE" "$DOWNLOAD_URL"; then
    info "Download successful to ${TEMP_ARCHIVE}."
else
    rm -f "$TEMP_ARCHIVE"
    error "Download failed. Please check the URL and your connection: ${DOWNLOAD_URL}"
fi

if [ -f "$INSTALL_PATH" ]; then
    info "Removing existing binary at $INSTALL_PATH..."
    rm "$INSTALL_PATH" || warn "Failed to remove existing binary, proceeding with caution."
fi

info "Extracting archive..."
mkdir -p "$TEMP_EXTRACT_DIR" || error "Failed to create temporary extraction directory: $TEMP_EXTRACT_DIR"

if tar -xf "$TEMP_ARCHIVE" -C "$TEMP_EXTRACT_DIR"; then
    info "Extraction successful to $TEMP_EXTRACT_DIR."
else
    rm -rf "$TEMP_ARCHIVE" "$TEMP_EXTRACT_DIR"
    error "Extraction failed."
fi

EXTRACTED_BINARY="${TEMP_EXTRACT_DIR}/gct-${os}-${arch}"
if [ ! -f "$EXTRACTED_BINARY" ]; then
    EXTRACTED_BINARY="${TEMP_EXTRACT_DIR}/${BIN_NAME}"
    if [ ! -f "$EXTRACTED_BINARY" ]; then
        EXTRACTED_BINARY=$(find "$TEMP_EXTRACT_DIR" -type f -executable -name "gct*" 2>/dev/null | head -n 1)
        if [ -z "$EXTRACTED_BINARY" ]; then
            rm -rf "$TEMP_ARCHIVE" "$TEMP_EXTRACT_DIR"
            error "Could not find executable in extracted contents."
        fi
    fi
fi

info "Found binary at: $EXTRACTED_BINARY"
info "Moving extracted binary to $INSTALL_DIR..."
mv "$EXTRACTED_BINARY" "$INSTALL_DIR/" || error "Failed to move extracted binary to $INSTALL_DIR."

rm -rf "$TEMP_ARCHIVE" "$TEMP_EXTRACT_DIR"

info "Making binary executable..."
chmod +x "$INSTALL_PATH" || error "Failed to set execute permission on: $INSTALL_PATH"

info "Checking if '${INSTALL_DIR}' is in PATH..."
if [ -n "$PATH" ] && [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    info "'${INSTALL_DIR}' is not found in your current PATH."
    info "Attempting to add it to your shell profile..."

    PROFILE_FILE=""
    if [ -n "$ZSH_VERSION" ]; then
        PROFILE_FILE="${ZDOTDIR:-$HOME}/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        PROFILE_FILE="$HOME/.bashrc"
    elif [ -f "$HOME/.profile" ]; then
        PROFILE_FILE="$HOME/.profile"
    else
        if [ -f "$HOME/.bash_profile" ]; then
            PROFILE_FILE="$HOME/.bash_profile"
        elif [ -f "$HOME/.zprofile" ]; then
             PROFILE_FILE="$HOME/.zprofile"
        fi
    fi

    if [ -n "$PROFILE_FILE" ] && [ -f "$PROFILE_FILE" ]; then
        info "Detected profile file: $PROFILE_FILE"
        EXPORT_LINE="export PATH=\"\$PATH:${INSTALL_DIR}\"" 
        COMMENT_LINE="# Add GCT installation directory to PATH"
        if ! grep -qF -- "$EXPORT_LINE" "$PROFILE_FILE"; then
            info "Adding PATH update to $PROFILE_FILE..."
            [[ $(tail -c1 "$PROFILE_FILE" | wc -l) -eq 0 ]] && echo >> "$PROFILE_FILE"
            echo "" >> "$PROFILE_FILE" 
            echo "$COMMENT_LINE" >> "$PROFILE_FILE"
            echo "$EXPORT_LINE" >> "$PROFILE_FILE"
            info "Successfully updated profile. Please run 'source ${PROFILE_FILE}' or restart your shell."
        else
            info "PATH update line already seems to exist in $PROFILE_FILE."
        fi
    else
        warn "Could not automatically detect or access a suitable shell profile file (.zshrc, .bashrc, .profile)."
        warn "Please add the following line to your shell configuration file manually:"
        warn "  export PATH=\"\$PATH:${INSTALL_DIR}\""
    fi
else
    info "'${INSTALL_DIR}' seems to be already in your PATH."
fi

echo ""
info "GCT ($(basename $TARGET_BIN)) installed/updated successfully to: ${INSTALL_DIR}"
info "Run 'gct --version' in a *new* shell/terminal tab to verify."
