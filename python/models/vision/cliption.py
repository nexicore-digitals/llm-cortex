import sys
import argparse
import json
import time
from typing import Dict
from PIL import Image
import torch
from transformers import BlipProcessor, BlipForConditionalGeneration

class CLIPtionPlugin:
    def __init__(self, model_path: str, device: str = "cpu", dtype=torch.float32):
        """
        Initializes the CLIPtion plugin with a BLIP image captioning model.
        """
        self.device = device
        self.dtype = dtype
        self.model_path = model_path

        print(f"[CLIPtion] Loading processor and model from {model_path}...")
        self.processor = BlipProcessor.from_pretrained(model_path)
        self.model = BlipForConditionalGeneration.from_pretrained(
            model_path,
            dtype=self.dtype
        ).to(self.device) #type: ignore
        print("[CLIPtion] Ready.")

    def invoke(self, image_path: str, prompt: str = "", max_length: int = 50) -> Dict[str, str | float]:
        """
        Generates a caption for the given image.
        """
        start = time.time()

        with Image.open(image_path) as img:
            raw_image = img.convert("RGB")

        inputs = self.processor(raw_image, prompt, return_tensors="pt").to(self.device, self.dtype) #type: ignore
        out = self.model.generate(**inputs, max_length=max_length)
        caption = self.processor.decode(out[0], skip_special_tokens=True)

        latency = time.time() - start

        return {
            "caption": caption,
            "latency": latency,
            "prompt": prompt,
            "image": image_path
        }

def main():
    """
    Main function to run the CLIPtion model from the command line.
    """
    parser = argparse.ArgumentParser(description="Run CLIPtion inference on an image.")
    parser.add_argument("--model-path", required=True, help="Path to the local CLIPtion model directory.")
    parser.add_argument("--image-path", required=True, help="Path to the input image.")
    parser.add_argument("--prompt", default="", help="Optional prompt for the model.")
    parser.add_argument("--max-length", type=int, default=50, help="Maximum number of tokens to generate.")
    args = parser.parse_args()

    try:
        device = "cuda" if torch.cuda.is_available() else "cpu"
        plugin = CLIPtionPlugin(model_path=args.model_path, device=device, dtype=torch.float16 if device == "cuda" else torch.float32)
        result = plugin.invoke(image_path=args.image_path, prompt=args.prompt, max_length=args.max_length)
        print(json.dumps(result, indent=2))
    except FileNotFoundError:
        error_msg = {"error": f"Image or model not found. Searched for image at '{args.image_path}' and model at '{args.model_path}'."}
        print(json.dumps(error_msg), file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        error_msg = {"error": f"An unexpected error occurred in CLIPtion: {e}"}
        print(json.dumps(error_msg), file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()