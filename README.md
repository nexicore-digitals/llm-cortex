# llm-cortex

A modular shell for orchestrating local GGUF-based LLMs with swap-aware memory management, plugin-centric routing, and model lifecycle control.

## Features

- 🧠 Model registry with plugin intent mapping
- 💾 Swap provisioning and memory-aware loading
- 🔁 Sequential and fallback model orchestration
- 🛠️ Scripts for swap monitoring and CLI wrapping

## Getting Started

1. Provision swap (see `scripts/swap-monitor.sh`)
2. Build `llama.cpp` and copy `llama-cli` here
3. Place GGUF models in `models/`
4. Run `go run router/cortex_shell.go`

## Directory Structure

```llm-cortex

llm-cortex/
├── bin/                  # All binaries live here
│   ├── llama-cli
│   ├── llama-server
│   └── llama-bench
├── models/               # GGUF models
├── router/               # Go orchestration logic
│   ├── cortex_shell.go
│   ├── model_registry.go
│   └── memory_manager.go
├── scripts/              # Bash helpers
│   └── swap-monitor.sh
├── README.md
└── LICENSE

```

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
