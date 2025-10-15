#!/usr/bin/env bash

set -e # Exit immediately if a command exits with a non zero status

info() {
  echo -e "\e[1;33m$1\e[0m"
}

error() {
  echo -e "\e[1;31m$1\e[0m"
}

success() {
  echo -e "\e[1;32m$1\e[0m"
}

# --- Swap File Logic ---
handle_swap() {
    # First, check if a swapfile already exists
    if [ -f "/swapfile" ]; then
        if swapon --show | grep -q "/swapfile"; then
            info "Swap file /swapfile already exists and is active."
            swapon --show
            return
        else
            info "Swap file /swapfile already exists but is not currently active."
            read -p "Would you like to activate it? (y/N): " -r choice
            if [[ "$choice" =~ ^[Yy]$ ]]; then
                info "Activating existing swap file. This requires sudo privileges."
                sudo swapon /swapfile
                success "Swap file activated."
                swapon --show
            fi
            # Whether activated or not, our job here is done.
            return
        fi
    fi

    # Check total RAM in GB
    local total_ram_gb
    total_ram_gb=$(free -g | awk '/^Mem:/{print $2}')

    # Check if the root partition is on an SSD (1 means rotational, 0 means non-rotational/SSD)
    local is_ssd=0
    local root_device
    root_device=$(df / | awk 'NR==2 {print $1}')
    # Follow symlinks to get the real device, e.g., /dev/dm-0 -> /dev/nvme0n1p2
    local real_device
    real_device=$(ls -l "$root_device" | awk '{print $NF}')
    if [[ "$real_device" == *"nvme"* || "$real_device" == *"sd"* ]]; then
        # Extract the base device name (e.g., nvme0n1, sda)
        local base_device
        base_device=$(echo "$real_device" | sed -E 's/p[0-9]+$//' | sed -E 's/[0-9]+$//' | xargs basename)
        if [ -f "/sys/block/$base_device/queue/rotational" ]; then
            if [ "$(cat "/sys/block/$base_device/queue/rotational")" -eq 1 ]; then
                is_ssd=1 # It's a rotational HDD
            fi
        fi
    fi

    if [ "$is_ssd" -eq 0 ]; then
        info "System appears to be running on an SSD."
        if [ "$total_ram_gb" -lt 16 ]; then
            info "You have ${total_ram_gb}GB of RAM. To prevent out-of-memory errors when running large models, creating a swap file (4GB-16GB recommended) can act as a safety net."
        else
            info "You have ${total_ram_gb}GB of RAM. Creating a swap file is optional, but can provide a safety net for exceptionally large models."
        fi

        read -p "Would you like to create a swap file? (y/N): " -r choice
        if [[ "$choice" =~ ^[Yy]$ ]]; then
            read -p "Enter the desired swap file size (e.g., 8G, 16G recommended): " -r swap_size
            if [ -z "$swap_size" ]; then
                error "No size entered. Skipping swap creation."
            else
                info "Attempting to create a ${swap_size} swap file. This requires sudo privileges."
                sudo ./scripts/setup_swap.sh --size "$swap_size"
            fi
        else
            info "Skipping swap file creation."
        fi
    else
        info "System appears to be running on a rotational HDD. Swap file is not recommended for performance."
    fi
}

# --------Configuration----------

# Set up Python environment
info "Running Python setup script..."
if [ ! -x "python-setup.sh" ]; then
    chmod +x ./python-setup.sh
fi
./python-setup.sh

info "Checking permissions for swap scripts..."
if [ -f "scripts/setup_swap.sh" ] && [ ! -x "scripts/setup_swap.sh" ]; then
    chmod +x scripts/setup_swap.sh
fi
if [ -f "scripts/teardown_swap.sh" ] && [ ! -x "scripts/teardown_swap.sh" ]; then
    chmod +x scripts/teardown_swap.sh
fi

info "Checking system memory and storage for swap file recommendation..."
handle_swap

if [ -f "bin/llama-cli" ]; then
    if [ ! -x "bin/llama-cli" ]; then
        echo "adding execution permission for the llama-cli"
        chmod 755 bin/llama-cli
    fi
  success "dependencies satisfied, skipping setup!"
  exit 0
fi

info "configuring and setting up project dependencies..."

# Check for existing llama.cpp build
info "Checking for existing llama.cpp build"

# Ensure Cmake is already installed 
if ! command -v cmake &> /dev/null; then
    info "cmake is not installed. Installing..."
    sudo apt update
    sudo apt install -y cmake
fi

# Check for git
if ! command -v git &> /dev/null; then
    error "git is not installed. Please install git and rerun this script."
    exit 1
fi

# Clone llama.cpp if not already present
if [ ! -d "llama.cpp" ]; then
    info "Cloning llama.cpp repository..."
    git clone https://github.com/ggml-org/llama.cpp
else
    info "llama.cpp repo already exists, pulling latest changes..."
    cd llama.cpp
    git pull
    cd ..
fi

cd llama.cpp

# Ensure libcurl is installed
if ! dpkg -s libcurl4-openssl-dev >/dev/null 2>&1; then
    info "libcurl not found. Installing..."
    sudo apt install -y libcurl4-openssl-dev
fi

info "Building llama.cpp with CUDA support..."

# Build the llama-cpp
# Check if CUDA is available
if command -v nvcc &> /dev/null; then
    info "CUDA detected! Building with GPU support..."
    cmake -B build -DGGML_CUDA=ON
else
    info "CUDA not detected. Building for CPU only..."
    cmake -B build
fi

cmake --build build --config Release

info "Moving binaries to root bin directory..."

# Copy the bin folder to the root folder
mkdir -p ../bin
mv build/bin/* ../bin/

# Navigate back to the root dir and remove unnecessary deps
cd ../
info "cleaning up..."

rm -rf ./llama.cpp


info "Adding execution permissions to binaries in bin/"
find "bin" -type f -exec chmod 755 {} +

info "Setting the LD_LIBRARY_PATH environment variable"
echo 'export LD_LIBRARY_PATH=./bin' >> ~/.bashrc
source ~/.bashrc

success "setup complete, happy coding!"
