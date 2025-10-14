from transformers import Blip2Processor, Blip2ForConditionalGeneration
from PIL import Image
import torch
import time
import argparse
import json
import sys
from typing import Dict


class BlipPlugin:
    def __init__(self, model_path: str, device: str = "cpu", dtype=torch.float32, use_fast: bool = True, legacy: bool = True):
        self.device = device
        self.dtype = dtype
        self.model_path = model_path
        self.use_fast = use_fast
        self.legacy = legacy

        print("[BLIP] Loading processor and model...")
        self.processor = Blip2Processor.from_pretrained(model_path, use_fast=self.use_fast, legacy=self.legacy)
        self.model = Blip2ForConditionalGeneration.from_pretrained(
            model_path,
            device_map=self.device if self.device != "cpu" else "auto",
            dtype=self.dtype
        )
        print("[BLIP] Ready.", flush=True)
    def invoke(self, image_path: str, prompt: str = "", max_length: int = 50) -> Dict[str, str | float]:
        """
        Runs BLIP-2 FLAN-T5-XL on the given image and prompt.

        Returns:
            dict: {
                "caption": str,
                "latency": float,
                "prompt": str,
                "image": str
            }
        """
        with Image.open(image_path) as img:
            image = img.convert("RGB")
        inputs = self.processor(images=image, text=prompt, return_tensors="pt").to(self.device, self.dtype) # type: ignore

        start = time.time()
        out = self.model.generate(**inputs, max_length=max_length)
        latency = time.time() - start

        caption = self.processor.decode(out[0], skip_special_tokens=True)
        return {
            "caption": caption,
            "latency": latency,
            "prompt": prompt,
            "image": image_path
        }

def main():
    """
    Main function to run the BLIP model from the command line.
    """
    parser = argparse.ArgumentParser(description="Run BLIP-2 inference on an image.")
    parser.add_argument("--model-path", type=str, required=True, help="Path to the local BLIP model directory.")
    parser.add_argument("--interactive", action="store_true", help="Run in interactive mode.")
    # Non-interactive mode arguments
    parser.add_argument("--image-path", type=str, help="Path to the input image (for non-interactive mode).")
    parser.add_argument("--device", type=str, help="Device to use for inference, e.g., 'cpu' or 'cuda'.")
    parser.add_argument("--prompt", type=str, default="", help="Optional prompt for the model (for non-interactive mode).")
    parser.add_argument("--use-fast", action="store_true", help="Use fast image processor if available (for non-interactive mode).")
    parser.add_argument("--max-length", type=int, default=75, help="Maximum number of tokens to generate (for non-interactive mode).")
    parser.add_argument("--no-legacy", action="store_false", dest="legacy", help="Use new tokenizer behavior instead of legacy (for non-interactive mode).")
    args = parser.parse_args()

    try:
        device = "cuda" if torch.cuda.is_available() else "cpu"
        plugin = BlipPlugin(model_path=args.model_path, device=device)

        if args.interactive:
            for line in sys.stdin:
                try:
                    input_data = json.loads(line)
                    if input_data.get("command") == "exit":
                        sys.exit(0)

                    # The processor is loaded once, so we can't change these per request,
                    # but we can honor them on the first interactive request if needed.
                    # For now, we just ensure the invoke call gets the right params.
                    result = plugin.invoke(
                        image_path=input_data.get("image_path"),
                        prompt=input_data.get("prompt", ""),
                        max_length=input_data.get("max_length", 75)
                    )
                    print(json.dumps(result), flush=True)
                    print("END_OF_JSON", flush=True)
                except json.JSONDecodeError:
                    # Ignore invalid JSON lines
                    pass
                except Exception as e:
                    print(json.dumps({"error": str(e)}), flush=True)
                    print("END_OF_JSON", flush=True)
        else:
            result = plugin.invoke(image_path=args.image_path, prompt=args.prompt, max_length=args.max_length)
            print(json.dumps(result, indent=2))

    except FileNotFoundError:
        error_msg = {"error": f"Image or model not found. Searched for image at '{args.image_path}' and model at '{args.model_path}'."}
        print(json.dumps(error_msg), file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        error_msg = {"error": f"An unexpected error occurred: {e}"}
        print(json.dumps(error_msg), file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
