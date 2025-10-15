# LLM-Cortex Model Zoo

This document provides an overview of the machine learning models used within the LLM-Cortex framework. For the application to function correctly, these models should be downloaded and placed in the appropriate subdirectories within this `models/` folder.

## Vision Models

These models are orchestrated by the Go application and executed via persistent Python processes for high performance.

---

### 1. Blip2

- **Model Name**: Blip2-Flan-T5-XL
- **Family**: Vision-Language Model
- **Details**: A powerful model that excels at visual question answering and image captioning. It combines a frozen image encoder (ViT) and a frozen LLM (Flan-T5-XL) with a lightweight trainable Querying Transformer (Q-Former).
- **Use Cases**: Visual Question Answering (VQA), Image Captioning.
- **License**: BSD 3-Clause "New" or "Revised" License
- **Main Download Link**: Salesforce/blip2-flan-t5-xl
- **Local Path**: `models/blip2-flan-t5-xl/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [blip2-flan-t5-xl](https://huggingface.co/Salesforce/blip2-flan-t5-xl)
2. Download the model files.
3. Place the downloaded files in the `models/blip2-flan-t5-xl/` directory.

---

### 2. CLIP

- **Model Name**: CLIP ViT-B/32
- **Family**: Vision-Language Model
- **Details**: Contrastive Language-Image Pre-Training (CLIP) is a model trained on a wide variety of (image, text) pairs. It can be instructed in natural language to predict the most relevant text snippet for a given image, without directly optimizing for the task.
- **Use Cases**: Zero-shot image classification.
- **License**: MIT
- **Main Download Link**: openai/clip-vit-base-patch32
- **Local Path**: `models/clip-vit-b-32/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [clip-vit-base-patch32](https://huggingface.co/laion/CLIP-ViT-B-32-laion2B-s34B-b79K)
2. Download the model files.
3. Place the downloaded files in the `models/clip-vit-b-32/` directory

---

### 3. CLIPtion

- **Model Name**: CLIPtion
- **Family**: Image Captioning
- **Details**: A custom image captioning model that uses a CLIP model as its vision backbone. This implementation is based on a specific architecture that requires a pre-trained `CLIPtion_20241219_fp16.safetensors` file.
- **Use Cases**: High-quality image captioning.
- **License**: Custom (dependent on the training data and base models). The CLIP base model is under the MIT license.
- **Download Link**: pharmapsychotic/CLIPtion
- **Local Path**: `models/CLIPtion/`

**Setup:**

1. Download the `CLIPtion_20241219_fp16.safetensors` file from [CLIPtion](https://huggingface.co/pharmapsychotic/CLIPtion).
2. Create the `models/CLIPtion` directory.
3. Place your `CLIPtion_20241219_fp16.safetensors` file inside it.

---

### 4. SDXL-Turbo-Ryzen-AI

- **Model Name**: SDXL-Turbo-Ryzen-AI
- **Family**: Text-to-Image Generation
- **Details**: A high-performance text-to-image generation model optimized for AMD Ryzen processors.
- **Use Cases**: Generating high-quality images from textual descriptions.
- **License**: CreativeML Open RAIL-M
- **Main Download Link**: stabilityai/sdxl-turbo-ryzen-ai
- **Local Path**: `models/sdxl-turbo-ryzen-ai/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [SDXL-Turbo-Ryzen-AI](https://huggingface.co/stabilityai/sdxl-turbo-ryzen-ai)
2. Download the model files.
3. Place the downloaded files in the `models/sdxl-turbo-ryzen-ai/` directory.

---

## Audio Models

These models are also orchestrated by the Go application and executed via persistent Python processes.

---

### 1. Whisper-Large-V3

- **Model Name**: Whisper Large V3
- **Family**: Automatic Speech Recognition (ASR)
- **Details**: Whisper is a general-purpose speech recognition model trained on a large dataset of diverse audio. The Large V3 variant offers high accuracy for transcription and translation tasks.
- **Use Cases**: Transcription, Translation.
- **License**: MIT
- **Main Download Link**: openai/whisper-large-v3
- **Local Path**: `models/whisper-large-v3/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [whisper-large-v3](https://huggingface.co/openai/whisper-large-v3)
2. Download the model files.
3. Place the downloaded files in the `models/whisper-large-v3/` directory.

---

### 2. Whisper-Large-V3-Turbo

- **Model Name**: Whisper Large V3 Turbo
- **Family**: Automatic Speech Recognition (ASR)
- **Details**: A more efficient variant of the Whisper Large V3 model, optimized for faster inference while maintaining high accuracy.
- **Use Cases**: Real-time transcription, Translation.
- **License**: MIT
- **Main Download Link**: openai/whisper-large-v3-turbo
- **Local Path**: `models/whisper-large-v3-turbo/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [whisper-large-v3-turbo](https://huggingface.co/openai/whisper-large-v3-turbo)
2. Download the model files.
3. Place the downloaded files in the `models/whisper-large-v3-turbo/` directory.

---

### 3. Whisper-Small

- **Model Name**: Whisper Small
- **Family**: Automatic Speech Recognition (ASR)
- **Details**: A smaller variant of the Whisper model, designed for lower resource environments while still providing good transcription accuracy.
- **Use Cases**: Transcription, Translation.
- **License**: MIT
- **Main Download Link**: openai/whisper-small
- **Local Path**: `models/whisper-small/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [whisper-small](https://huggingface.co/openai/whisper-small)
2. Download the model files.
3. Place the downloaded files in the `models/whisper-small/` directory.

---

### 4. Whisper-Tiny

- **Model Name**: Whisper Tiny
- **Family**: Automatic Speech Recognition (ASR)
- **Details**: The smallest variant of the Whisper model, optimized for extremely low resource environments.
- **Use Cases**: Transcription, Translation.
- **License**: MIT
- **Main Download Link**: openai/whisper-tiny
- **Local Path**: `models/whisper-tiny/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [whisper-tiny](https://huggingface.co/openai/whisper-tiny)
2. Download the model files.
3. Place the downloaded files in the `models/whisper-tiny/` directory.

---

### 5. XTTS

- **Model Name**: XTTS
- **Family**: Text-to-Speech (TTS)
- **Details**: A state-of-the-art text-to-speech model designed for high-quality, natural-sounding speech synthesis.
- **Use Cases**: Voiceovers, Assistive technology, Audiobooks.
- **License**: MIT
- **Main Download Link**: openai/xtts
- **Local Path**: `models/xtts/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [xtts](https://huggingface.co/openai/xtts)
2. Download the model files.
3. Place the downloaded files in the `models/xtts/` directory.

---

## Language Models (GGUF)

The framework is also designed to orchestrate GGUF-based Large Language Models (LLMs) using the `llama.cpp` engine, which is built during the setup process.

### All-MiniLM-L6-v2

- **Model Name**: All-MiniLM-L6-v2
- **Family**: Language Model
- **Details**: A small, efficient transformer model optimized for sentence embeddings. It is suitable for tasks like semantic search and clustering.
- **Use Cases**: Semantic search, sentence embeddings.
- **License**: MIT
- **Main Download Link**: sentence-transformers/all-MiniLM-L6-v2
- **Local Path**: `models/all-MiniLM-L6-v2/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [all-MiniLM-L6-v2](https://huggingface.co/leliuga/all-MiniLM-L6-v2-GGUF)
2. Download the model files.
3. Place the downloaded files in the `models/all-MiniLM-L6-v2/` directory.

---

### Qwen 2.5 Coder

- **Model Name**: Qwen2.5-Coder-7B-Instruct-Q6_K_L
- **Family**: Qwen
- **Details**: A transformer-based LLM from Alibaba Cloud, optimized for coding tasks. GGUF is a quantized format that allows the model to run efficiently on a CPU.
- **Use Cases**: Code generation, conversational AI, instruction following.
- **License**: Tongyi Qianwen LICENSE AGREEMENT
- **Download Link**: Qwen/Qwen2.5-Coder-7B-Instruct-GGUF
- **Local Path**: `models/qwen/Qwen2.5-Coder-7B-Instruct-Q6_K_L.gguf`

**Setup Instructions:**

1. Go to the Hugging Face model page: [Qwen2.5-Coder-7B-Instruct-GGUF](https://huggingface.co/bartowski/Qwen2.5-Coder-7B-Instruct-GGUF)
2. Download the model files.
3. Place the downloaded files in the `models/qwen/` directory.

### Qwen 2.5 Vision-Language

- **Model Name**: Qwen2.5-VL-7B-Instruct-Q6_K
- **Family**: Qwen
- **Details**: A transformer-based LLM from Alibaba Cloud, optimized for vision-language tasks. GGUF is a quantized format that allows the model to run efficiently on a CPU.
- **Use Cases**: Vision-language tasks, image captioning, visual question answering.
- **License**: Tongyi Qianwen LICENSE AGREEMENT
- **Download Link**: Qwen/Qwen2.5-VL-7B-Instruct-GGUF
- **Local Path**: `models/qwen/Qwen2.5-VL-7B-Instruct-Q6_K.gguf`

**Setup Instructions:**

1. Go to the Hugging Face model page: [Qwen2.5-VL-7B-Instruct-GGUF](https://huggingface.co/unsloth/Qwen2.5-VL-7B-Instruct-GGUF)
2. Download the model files.
3. Place the downloaded files in the `models/qwen/` directory.

---

### StarCoder

- **Model Name**: StarCoder
- **Family**: Code Generation
- **Details**: A state-of-the-art model for code generation tasks, fine-tuned on a diverse range of programming languages and tasks.
- **Use Cases**: Code completion, code translation, code summarization.
- **License**: Apache 2.0
- **Main Download Link**: bigcode/starcoder2-15b
- **Local Path**: `models/starcoder/`

**Setup Instructions:**

1. Go to the Hugging Face model page: [StarCoder](https://huggingface.co/bartowski/starcoder2-15b-instruct-v0.1-GGUF)
2. Download the model files.
3. Place the downloaded files in the `models/starcoder/` directory.

---

You can use any other GGUF-compatible model by downloading it and placing it in the `models/` directory.

This file provides a clear and centralized reference for all models in your project.
