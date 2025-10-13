# llm-cortex

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Language: Go](https://img.shields.io/badge/Language-Go-blue.svg)
![Language: Shell](https://img.shields.io/badge/Language-Shell-lightgrey.svg)
![Language: Python](https://img.shields.io/badge/Language-Python-3776AB.svg)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)

A modular shell for orchestrating local models. It uses Go for GGUF-based LLMs and Python for other modalities like vision and audio, with swap-aware memory management, plugin-centric routing, and model lifecycle control.
The Go application acts as the central orchestrator, calling standalone Python scripts for specific tasks like image captioning.

## Features

- 🌐 **Multi-language Orchestration:** Uses Go for high-performance LLM management and calls standalone Python scripts for vision/audio tasks.
- 🧠 Model registry with plugin intent mapping
- 💾 Swap provisioning and memory-aware loading
- 🔁 Sequential and fallback model orchestration
- 🛠️ Scripts for swap monitoring and CLI wrapping

## Getting Started

All you need to do is run the setup script. This will install all necessary dependencies, build the required binaries, and configure your environment.

```bash
chmod +x ./setup.sh
./setup.sh
```

## Usage

1. **Add GGUF models:** Place your GGUF-formatted models into the `models/` directory.
2. **Provision swap (Optional):** If you plan on running large models, ensure you have enough swap space. A helper script is provided in `scripts/`.
3. **Start the shell:** Run the main Go application to start the orchestrator.

    ```bash
    go run router/cortex_shell.go
    ```

## Directory Structure

llm-cortex/
├── bin/                  # All binaries live here
├── models/               # GGUF models
├── python/
│   ├── models/
│   │   └── vision/       # Python scripts for vision models (e.g., blip.py)
│   └── requirements.txt  # Python dependencies
├── router/               # Go orchestration logic
├── scripts/              # Bash helpers
├── README.md
├── LICENSE
├── setup.sh              # Project setup and dependency installer
├── python-setup.sh       # Python environment setup script
└── llama-cli.sh          # Wrapper script for llama-cli

## models/

This folder is for local GGUF models used by the orchestration shell.

**Do not commit model files.** They are large and system-specific.

To use:

1. Download GGUF models from Hugging Face
2. Place them here following the structure in `router/model_registry.go`

## Contributing

Contributions are welcome! Please open issues or pull requests for enhancements or bug fixes.

## License

MIT
See [LICENSE](LICENSE) for details.
