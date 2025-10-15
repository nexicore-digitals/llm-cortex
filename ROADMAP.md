# LLM-Cortex Project Roadmap

This document outlines the planned features and development milestones for LLM-Cortex. Our goal is to evolve this framework into a comprehensive, production-ready system for orchestrating a diverse range of local machine learning models.

---

## Phase 1: Core Architecture & Foundation (Current -> Q3 2024)

This phase focuses on solidifying the core architecture to create a stable and extensible foundation.

- **[Done]** **Process & Session Management**: Implement a robust `spawn` package for managing persistent Python processes.
- **[Done]** **Parallel Execution**: Create a `scheduler` package to run multiple models concurrently.
- **[In Progress]** **Centralized Engine**:
  - Fully implement the `/core/engine` package to manage the application lifecycle.
  - Abstract all model loading and orchestration logic away from `main.go` into the engine.
- **[In Progress]** **Centralized Configuration**:
  - Fully integrate the `/core/config` package.
  - Load all configurations (Python path, server port, model paths) from a `config.yaml` file and environment variables, removing all hardcoded values.
- **Formalize Model Plugin System**:
  - Define a standard Go interface for all models (e.g., `type ModelPlugin interface { Load(); Invoke(); Unload(); }`).
  - Refactor the existing vision models (`BLIP`, `CLIP`) to adhere to this new interface.

---

## Phase 2: Model Expansion & API Development (Q4 2024)

With a solid foundation, this phase focuses on expanding the variety of supported models and building a clean API.

- **Integrate Audio Models**:
  - Implement Go wrappers and Python scripts for **Whisper** (ASR) to support transcription.
  - Implement Go wrappers and Python scripts for **XTTS** (TTS) to support text-to-speech synthesis.
- **Integrate Text-to-Image Models**:
  - Add support for text-to-image models like **SDXL-Turbo**.
- **Official GGUF Integration**:
  - Create a dedicated Go wrapper for `llama.cpp` that uses the `spawn` package, similar to the vision models.
  - This will enable first-class support for chat, completion, and embedding tasks with GGUF models.
- **RESTful API v1**:
  - Design and implement a clean, resource-oriented RESTful API.
  - Create dedicated endpoints for each modality (e.g., `/api/v1/vision/caption`, `/api/v1/audio/transcribe`, `/api/v1/llm/chat`).

---

## Phase 3: Usability & Production Readiness (Q1 2025)

This phase focuses on making the framework easier to use, deploy, and monitor.

- **Streaming Support**:
  - Implement real-time streaming for ASR (Whisper) and LLM (GGUF) endpoints using WebSockets or Server-Sent Events (SSE).
- **Web UI**:
  - Develop a simple, functional web interface to demonstrate the capabilities of the different models served by the backend.
- **Dockerization**:
  - Create a `Dockerfile` and `docker-compose.yml` to containerize the entire application, including the Go backend, Python environment, and dependencies. This will drastically simplify setup and deployment.
- **Improved Logging & Monitoring**:
  - Integrate a structured logging library (e.g., `zerolog`).
  - Expose application metrics (e.g., model latency, memory usage) via a `/metrics` endpoint for Prometheus scraping.

---

## Future Ideas (Beyond Phase 3)

- **Model Hot-Swapping**: Allow loading and unloading models at runtime without restarting the application.
- **GPU Resource Management**: Intelligently assign models to specific GPUs and manage VRAM.
- **Request Batching**: Implement dynamic request batching for supported models to increase throughput under heavy load.
- **Distributed Orchestration**: Extend the engine to manage models running on multiple nodes.

---

*This roadmap is a living document and is subject to change based on project priorities and community feedback.*
