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

# --------Configuration----------

# Set up Python environment
info "Running Python setup script..."
if [ ! -x "python-setup.sh" ]; then
    chmod +x ./python-setup.sh
fi
./python-setup.sh

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
