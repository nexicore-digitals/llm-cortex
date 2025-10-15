#!/usr/bin/env bash

set -e # Exit immediately if a command exits with a non-zero status

# --- Helper Functions ---
info() {
  echo -e "\e[1;33m$1\e[0m"
}

error() {
  echo -e "\e[1;31m$1\e[0m"
}

success() {
  echo -e "\e[1;32m$1\e[0m"
}

# --- Configuration ---
SWAP_FILE="/swapfile"

# --- Main Script ---

# Check for root privileges
if [ "$EUID" -ne 0 ]; then
  error "Please run as root or with sudo."
  exit 1
fi

# Check if the swap file is active and deactivate it
if swapon --show | grep -q "$SWAP_FILE"; then
  info "Deactivating swap file: $SWAP_FILE..."
  swapoff "$SWAP_FILE"
else
  info "Swap file $SWAP_FILE is not active."
fi

# Remove the swap entry from /etc/fstab if it exists
if grep -q "$SWAP_FILE" /etc/fstab; then
  info "Removing swap entry from /etc/fstab..."
  sed -i "\|$SWAP_FILE|d" /etc/fstab
fi

# Remove the swap file itself if it exists
if [ -f "$SWAP_FILE" ]; then
  info "Deleting swap file: $SWAP_FILE..."
  rm "$SWAP_FILE"
fi

success "Swap teardown complete. No swap should be active."
swapon --show