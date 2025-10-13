#!/usr/bin/env bash

set -e # Exit immediately if a command exits with a non-zero status

VENV_DIR="python_venv"
FLAG_FILE="$VENV_DIR/.setup_complete"

info() {
  echo -e "\e[1;33m$1\e[0m"
}

success() {
  echo -e "\e[1;32m$1\e[0m"
}

error() {
  echo -e "\e[1;31m$1\e[0m"
}


info "Setting up Python environment..."

# Check for python3
if ! command -v python3 &> /dev/null; then
    error "python3 could not be found. Please install Python 3 and run this script again."
    exit 1
fi

# On Debian/Ubuntu, the 'venv' module is in a separate package.
if [ -f /etc/debian_version ]; then
    # Check if python3-venv is installed
    if ! dpkg -s python3-venv >/dev/null 2>&1; then
        info "'python3-venv' package not found. It is required to create virtual environments. Installing..."
        sudo apt update && sudo apt install -y python3-venv
    fi
fi

if [ ! -f "$VENV_DIR/bin/activate" ]; then
    info "Virtual environment is incomplete or does not exist."
    if [ -d "$VENV_DIR" ]; then
        info "Removing incomplete venv directory: $VENV_DIR"
        rm -rf "$VENV_DIR"
    fi
    info "Creating Python virtual environment in '$VENV_DIR'..."
    python3 -m venv "$VENV_DIR"
fi

if [ -f "python/requirements.txt" ]; then
    # Check if the flag file exists and if requirements.txt is not newer than the flag
    if [ -f "$FLAG_FILE" ] && [ ! "python/requirements.txt" -nt "$FLAG_FILE" ]; then
        success "Python dependencies are already installed and up-to-date."
    else
        info "Installing/updating Python dependencies from python/requirements.txt..."
        # Activate venv, install requirements, and then create the flag file on success
        source "$VENV_DIR/bin/activate"
        pip install -r python/requirements.txt
        deactivate
        touch "$FLAG_FILE"
        success "Python dependencies installed successfully."
    fi
else
    info "python/requirements.txt not found, skipping dependency installation."
fi

success "Python setup complete."