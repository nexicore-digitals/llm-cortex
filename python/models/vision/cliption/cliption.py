import os
import sys
import json
import time
import argparse
from types import SimpleNamespace

import torch
from safetensors import safe_open
from PIL import Image
from transformers import CLIPModel, CLIPProcessor
from model import CLIPtionModel


# ----------------------------
# Loader
# ----------------------------
def load_cliption(model_path: str, device: str = "cpu", dtype=torch.float16, use_fast: bool = True) -> CLIPtionModel:
    print(f"[CLIPtion] Loading processor and model from {model_path}...")
    clip_model = CLIPModel.from_pretrained("openai/clip-vit-large-patch14", dtype=dtype).to(device)
    processor = CLIPProcessor.from_pretrained("openai/clip-vit-large-patch14", use_fast=use_fast)

    safetensor_file = os.path.join(model_path, "CLIPtion_20241219_fp16.safetensors")
    if not os.path.exists(safetensor_file):
        raise FileNotFoundError(f"Could not find {safetensor_file}")

    state_dict = {}
    with safe_open(safetensor_file, framework="pt", device="cpu") as f:
        for key in f.keys():
            state_dict[key] = f.get_tensor(key)

    tp_dict = {"weight": state_dict.pop("text_projection.weight")}
    config = SimpleNamespace(hidden_dim=768, num_heads=8, num_blocks=6, max_length=77)

    model = CLIPtionModel(clip_model, processor, config)
    model.captioner.load_state_dict(state_dict)
    model.text_projection.load_state_dict(tp_dict)
    model.eval().to(device, dtype=dtype)
    print("[CLIPtion] Ready.")
    return model


# ----------------------------
# CLI
# ----------------------------
def main():
    parser = argparse.ArgumentParser(description="Run CLIPtion captioning on an image")
    parser.add_argument("--model-path", required=True, help="Path to CLIPtion folder with safetensors")
    parser.add_argument("--image-path", required=True, help="Path to input image")
    parser.add_argument("--beam-search", action="store_true", help="Use beam search decoding")
    parser.add_argument("--use-fast", action="store_true", help="Use fast image processor if available.")
    parser.add_argument("--beam-width", type=int, default=5, help="Beam width for beam search.")
    parser.add_argument("--best-of", type=int, default=5, help="Number of candidates for sampling.")
    parser.add_argument("--temperature", type=float, default=1.0, help="Temperature for sampling.")
    parser.add_argument("--device", default="auto", help="'cpu' or 'cuda' (auto-detects GPU)")
    args = parser.parse_args()

    device = args.device if args.device != "auto" else ("cuda" if torch.cuda.is_available() else "cpu")

    if device == "cpu":
        dtype = torch.float32
    else:
        dtype = torch.float16

    try:
        start_time = time.time()
        model = load_cliption(args.model_path, device=device, dtype=dtype, use_fast=args.use_fast)

        # Load and preprocess image
        with Image.open(args.image_path) as img:
            image = img.convert("RGB")

        inputs = model.processor(images=image, return_tensors="pt").to(device, dtype)
        image_tensor = inputs["pixel_values"]  # [1, 3, 224, 224]

        # Generate caption
        if args.beam_search:
            captions = model.generate_beam(image_tensor, beam_width=args.beam_width)
        else:
            captions = model.generate(image_tensor, best_of=args.best_of, temperature=args.temperature)

        latency = time.time() - start_time
        print(json.dumps({"caption": captions[0], "latency": latency, "image": args.image_path}, indent=2))

    except Exception as e:
        print(json.dumps({"error": f"Unexpected error: {e}"}))
        sys.exit(1)


if __name__ == "__main__":
    main()
