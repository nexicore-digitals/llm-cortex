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
SWAP_SIZE=""

# --- Argument Parsing ---
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --size) SWAP_SIZE="$2"; shift ;;
        *) error "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# --- Main Script ---

# Check for root privileges
if [ "$EUID" -ne 0 ]; then
  error "Please run as root or with sudo."
  exit 1
fi

# Check if a size argument was provided
if [ -z "$SWAP_SIZE" ]; then
  error "Usage: sudo $0 --size <size>"
  info "Example: sudo $0 --size 8G"
  exit 1
fi

# Check if swap is already active to avoid errors
if swapon --show | grep -q "$SWAP_FILE"; then
  info "Swap file $SWAP_FILE is already active. No action taken."
  exit 0
fi

# Check if swap file already exists
if [ -f "$SWAP_FILE" ]; then
  error "Swap file $SWAP_FILE already exists but is not active."
  info "Please remove it manually or use the teardown_swap.sh script before creating a new one."
  exit 1
fi

info "Creating a ${SWAP_SIZE} swap file at ${SWAP_FILE}..."
fallocate -l "$SWAP_SIZE" "$SWAP_FILE"

info "Setting correct permissions (600) on the swap file..."
chmod 600 "$SWAP_FILE"

info "Formatting the file as swap space..."
mkswap "$SWAP_FILE"

info "Activating the swap file..."
swapon "$SWAP_FILE"

info "Making the swap file persistent across reboots..."
# Check if the entry already exists in /etc/fstab to avoid duplicates
if ! grep -q "$SWAP_FILE" /etc/fstab; then
  echo "$SWAP_FILE none swap sw 0 0" >> /etc/fstab
  info "Added swap entry to /etc/fstab."
else
  info "Swap entry already exists in /etc/fstab."
fi

success "Swap setup complete. Current swap status:"
swapon --show