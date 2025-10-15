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

class CLIPtionPlugin:
    def __init__(self, model_path: str, device: str = "cpu", dtype=torch.float16, use_fast: bool = True):
        print(f"[CLIPtion] Loading processor and model from {model_path}...")
        self.device = device
        self.dtype = dtype
        self.use_fast = use_fast
        self.clip_model = CLIPModel.from_pretrained("openai/clip-vit-large-patch14", dtype=self.dtype).to(self.device) # type: ignore
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-large-patch14", use_fast=self.use_fast)

        safetensor_file = os.path.join(model_path, "CLIPtion_20241219_fp16.safetensors")
        if not os.path.exists(safetensor_file):
            raise FileNotFoundError(f"Could not find {safetensor_file}")

        state_dict = {}
        with safe_open(safetensor_file, framework="pt", device="cpu") as f:
            for key in f.keys():
                state_dict[key] = f.get_tensor(key)

        tp_dict = {"weight": state_dict.pop("text_projection.weight")}
        config = SimpleNamespace(hidden_dim=768, num_heads=8, num_blocks=6, max_length=77)

        self.model = CLIPtionModel(self.clip_model, self.processor, config)
        self.model.captioner.load_state_dict(state_dict)
        self.model.text_projection.load_state_dict(tp_dict)
        self.model.eval().to(self.device, self.dtype)
        print("[CLIPtion] Ready.", flush=True)

    def invoke(self, image_path: str, beam_search: bool = False, beam_width: int = 5, best_of: int = 5, temperature: float = 1.0) -> dict:
        start_time = time.time()
        with Image.open(image_path) as img:
            image = img.convert("RGB")

        inputs = self.processor(images=image, return_tensors="pt").to(self.device, self.dtype) # type: ignore
        image_tensor = inputs["pixel_values"]

        if beam_search:
            captions = self.model.generate_beam(image_tensor, beam_width=beam_width)
        else:
            captions = self.model.generate(image_tensor, best_of=best_of, temperature=temperature)

        latency = time.time() - start_time
        return {"caption": captions[0], "latency": latency, "image": image_path}

# ----------------------------
# CLI
# ----------------------------
def main():
    parser = argparse.ArgumentParser(description="Run CLIPtion captioning on an image")
    parser.add_argument("--interactive", action="store_true", help="Run in interactive mode.")
    parser.add_argument("--model-path", required=True, help="Path to CLIPtion folder with safetensors")
    parser.add_argument("--device", default="auto", help="'cpu' or 'cuda' (auto-detects GPU)")
    # Non-interactive arguments
    parser.add_argument("--image-path", help="Path to input image")
    parser.add_argument("--beam-search", action="store_true", help="Use beam search decoding")
    parser.add_argument("--use-fast", action="store_true", help="Use fast image processor if available.")
    parser.add_argument("--beam-width", type=int, default=5, help="Beam width for beam search.")
    parser.add_argument("--best-of", type=int, default=5, help="Number of candidates for sampling.")
    parser.add_argument("--temperature", type=float, default=1.0, help="Temperature for sampling.")
    args = parser.parse_args()

    device = args.device if args.device != "auto" else ("cuda" if torch.cuda.is_available() else "cpu")
    dtype = torch.float32 if device == "cpu" else torch.float16

    try:
        plugin = CLIPtionPlugin(args.model_path, device=device, dtype=dtype, use_fast=args.use_fast)
        if args.interactive:
            for line in sys.stdin:
                try:
                    input_data = json.loads(line)
                    if input_data.get("command") == "exit":
                        sys.exit(0)

                    result = plugin.invoke(
                        image_path=input_data.get("image_path"),
                        beam_search=input_data.get("beam_search", False),
                        beam_width=input_data.get("beam_width", 5),
                        best_of=input_data.get("best_of", 5),
                        temperature=input_data.get("temperature", 1.0)
                    )
                    print(json.dumps(result), flush=True)
                    print("END_OF_JSON", flush=True)
                except Exception as e:
                    print(json.dumps({"error": str(e)}), flush=True)
                    print("END_OF_JSON", flush=True)
        else:
            if not args.image_path:
                print(json.dumps({"error": "image-path is required for non-interactive mode"}), file=sys.stderr)
                sys.exit(1)
            result = plugin.invoke(
                image_path=args.image_path,
                beam_search=args.beam_search,
                beam_width=args.beam_width,
                best_of=args.best_of,
                temperature=args.temperature
            )
            print(json.dumps(result, indent=2))

    except Exception as e:
        print(json.dumps({"error": f"Unexpected error: {e}"}))
        sys.exit(1)


if __name__ == "__main__":
    main()
