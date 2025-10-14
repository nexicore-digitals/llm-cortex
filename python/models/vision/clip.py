import sys
import argparse
import json
import time
from typing import Dict, List
from PIL import Image
import torch
from transformers import CLIPProcessor, CLIPModel

class CLIPPlugin:
    def __init__(self, model_path: str, device: str = "cpu", dtype=torch.float32, use_fast: bool = True):
        self.device = device
        self.dtype = dtype
        self.model_path = model_path
        self.use_fast = use_fast

        print(f"[CLIP] Loading processor and model from {model_path}...")
        self.model = CLIPModel.from_pretrained(model_path, dtype=self.dtype).to(self.device) #type: ignore
        self.processor = CLIPProcessor.from_pretrained(model_path, use_fast=self.use_fast)
        print("[CLIP] Ready.", flush=True)

    def invoke(self, image_path: str, texts: List[str]) -> Dict:
        """
        Performs zero-shot classification on an image against a list of text labels.
        """
        start = time.time()
        with Image.open(image_path) as img:
            image = img.convert("RGB")

        inputs = self.processor(text=texts, images=image, return_tensors="pt", padding=True).to(self.device, self.dtype) # type: ignore

        with torch.no_grad():
            outputs = self.model(**inputs)

        # Cosine similarity as logits
        logits_per_image = outputs.logits_per_image
        # Softmax to get probabilities
        probs = logits_per_image.softmax(dim=1)

        latency = time.time() - start

        # Create a dictionary of text labels and their probabilities
        results = {text: prob.item() for text, prob in zip(texts, probs[0])}

        return {
            "results": results,
            "latency": latency,
            "image": image_path
        }

def main():
    """
    Main function to run the CLIP model from the command line.
    """
    parser = argparse.ArgumentParser(description="Run CLIP zero-shot classification on an image.")
    parser.add_argument("--model-path", type=str, required=True, help="Path to the local CLIP model directory.")
    parser.add_argument("--interactive", action="store_true", help="Run in interactive mode.")
    parser.add_argument("--device", type=str, default="auto", help="Device to run the model on, e.g., 'cpu' or 'cuda'.")
    # Non-interactive mode arguments
    parser.add_argument("--image-path", type=str, help="Path to the input image (for non-interactive mode).")
    parser.add_argument("--texts", nargs='+', help="A list of text labels to classify against (for non-interactive mode).")
    parser.add_argument("--use-fast", action="store_true", help="Use fast image processor if available (for non-interactive mode).")
    args = parser.parse_args()

    try:
        device = args.device if args.device != "auto" else ("cuda" if torch.cuda.is_available() else "cpu")
        dtype = torch.float16 if device == "cuda" else torch.float32
        plugin = CLIPPlugin(model_path=args.model_path, device=device, dtype=dtype)

        if args.interactive:
            for line in sys.stdin:
                try:
                    input_data = json.loads(line)
                    if input_data.get("command") == "exit":
                        sys.exit(0)

                    result = plugin.invoke(
                        image_path=input_data.get("image_path"),
                        texts=input_data.get("texts", [])
                    )
                    print(json.dumps(result), flush=True)
                    print("END_OF_JSON", flush=True)
                except Exception as e:
                    print(json.dumps({"error": str(e)}), flush=True)
                    print("END_OF_JSON", flush=True)
        else:
            plugin.use_fast = args.use_fast
            result = plugin.invoke(image_path=args.image_path, texts=args.texts)
            print(json.dumps(result, indent=2))

    except FileNotFoundError:
        error_msg = {"error": f"Image or model not found. Searched for image at '{args.image_path}' and model at '{args.model_path}'."}
        print(json.dumps(error_msg), file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        error_msg = {"error": f"An unexpected error occurred in CLIP: {e}"}
        print(json.dumps(error_msg), file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()