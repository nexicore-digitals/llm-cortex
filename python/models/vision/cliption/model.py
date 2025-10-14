from typing import List, Tuple

import torch
import torch.nn as nn
import torch.nn.functional as F
from transformers import CLIPModel, CLIPProcessor


# ----------------------------
# Decoder / Captioner Classes
# ----------------------------
class DecoderBlock(nn.Module):
    def __init__(self, embed_dim: int, num_heads: int):
        super().__init__()
        self.norm1 = nn.LayerNorm(embed_dim)
        self.self_attn = nn.MultiheadAttention(embed_dim, num_heads, batch_first=True)
        self.norm2 = nn.LayerNorm(embed_dim)
        self.cross_attn = nn.MultiheadAttention(embed_dim, num_heads, batch_first=True)
        self.norm3 = nn.LayerNorm(embed_dim)

        # Matches checkpoint MLP structure
        self.mlp = nn.ModuleDict({
            "0": nn.Linear(embed_dim, embed_dim * 4),
            "1": nn.GELU(),
            "3": nn.Linear(embed_dim * 4, embed_dim),
        })

    def forward(self, x, memory, self_attn_mask=None):
        residual = x
        x = self.norm1(x)
        attn_output, _ = self.self_attn(query=x, key=x, value=x, attn_mask=self_attn_mask, is_causal=True)
        x = residual + attn_output

        residual = x
        x = self.norm2(x)
        attn_output, _ = self.cross_attn(query=x, key=memory, value=memory)
        x = residual + attn_output

        residual = x
        x = self.norm3(x)
        x = residual + self.mlp["3"](self.mlp["1"](self.mlp["0"](x)))
        return x


class Captioner(nn.Module):
    def __init__(self, config, vision_embed_dim: int, vocab_size: int):
        super().__init__()
        self.hidden_dim = config.hidden_dim
        self.max_length = config.max_length
        self.vocab_size = vocab_size

        self.projection = nn.Linear(vision_embed_dim, self.hidden_dim)
        self.memory_pos_embedding = nn.Parameter(torch.zeros(1, 257, self.hidden_dim))
        self.layers = nn.ModuleList(
            [DecoderBlock(config.hidden_dim, config.num_heads) for _ in range(config.num_blocks)]
        )

        causal_mask = nn.Transformer.generate_square_subsequent_mask(self.max_length)
        self.register_buffer("causal_mask", causal_mask, persistent=False)


# ----------------------------
# CLIPtion Model
# ----------------------------
class CLIPtionModel(nn.Module):
    def __init__(self, clip_model: CLIPModel, processor: CLIPProcessor, config):
        super().__init__()
        self.clip = clip_model
        self.processor = processor
        self.device = next(clip_model.parameters()).device
        self.tokenizer = processor.tokenizer # type: ignore

        self.captioner = Captioner(config, vision_embed_dim=1024, vocab_size=self.tokenizer.vocab_size)
        self.text_projection = nn.Linear(768, 768, bias=False)
        self.output_projection = nn.Linear(self.captioner.hidden_dim, self.tokenizer.vocab_size, bias=False)
        self.output_projection.weight = nn.Parameter(
            self.clip.text_model.embeddings.token_embedding.weight.clone()
        )

    # ----------------------------
    # Generation methods
    # ----------------------------
    def generate(self, images: torch.Tensor, seed: int = 42, temperature: float = 0.7,
                 best_of: int = 1, ramble: bool = False) -> List[str]:
        device = self.device
        image_features, image_embeds = self._images_to_embeds(images, device)

        captions = []
        for idx in range(image_features.size(0)):
            features = image_features[idx:idx + 1]
            tokens = self._batch_generate(features, temperature, best_of, seed + idx, ramble)
            text = self.tokenizer.decode(tokens[0], skip_special_tokens=True, clean_up_tokenization_spaces=True)
            captions.append(text)
        return captions

    def generate_beam(self, images: torch.Tensor, beam_width: int = 4, ramble: bool = False) -> List[str]:
        device = self.device
        image_features, image_embeds = self._images_to_embeds(images, device)

        captions = []
        for idx in range(image_features.size(0)):
            features = image_features[idx].unsqueeze(0)
            candidates = self._beam_search(features, image_embeds[idx:idx + 1], device, beam_width, ramble)
            candidates.sort(key=lambda x: x[0], reverse=True)
            captions.append(candidates[0][1])
        return captions

    # ----------------------------
    # Helper methods
    # ----------------------------
    def _images_to_embeds(self, images: torch.Tensor, device: torch.device) -> Tuple[torch.Tensor, torch.Tensor]:
        """
        Extract per-patch features (for decoder) and global embeddings (for normalization).
        """
        # Vision backbone features
        vision_outputs = self.clip.vision_model(images)
        features = vision_outputs.last_hidden_state.to(device=device, dtype=images.dtype)

        # Global image embeddings from full CLIPModel
        embeds = self.clip.get_image_features(images) # type: ignore
        embeds = embeds.to(device=device, dtype=images.dtype)
        embeds /= embeds.norm(dim=-1, keepdim=True)

        return features, embeds

    def _batch_generate(self, image_features: torch.Tensor, temperature: float, batch_size: int,
                        seed: int = None, ramble: bool = False) -> torch.Tensor: # type: ignore
        tokenizer = self.tokenizer
        output_proj = self.output_projection
        token_embedding_ = self.clip.text_model.embeddings.token_embedding
        pos_embedding_ = self.clip.text_model.embeddings.position_embedding

        memory = self.captioner.projection(image_features)
        memory = memory + self.captioner.memory_pos_embedding
        memory = memory.repeat(batch_size, 1, 1)

        sequences = torch.full((batch_size, self.captioner.max_length), tokenizer.eos_token_id,
                               dtype=torch.long, device=image_features.device)
        sequences[:, 0] = tokenizer.bos_token_id

        generator = torch.Generator(device=image_features.device)
        if seed is not None:
            generator.manual_seed(seed)

        for t in range(1, self.captioner.max_length - 1):
            token_embeds = token_embedding_(sequences[:, :t])
            pos_embeds = pos_embedding_(torch.arange(t, device=sequences.device))
            x = token_embeds + pos_embeds

            mask = self.captioner.causal_mask[:t, :t] # type: ignore
            for layer in self.captioner.layers:
                x = layer(x, memory, self_attn_mask=mask)

            logits = output_proj(x[:, -1:]) / temperature

            prev_is_eos = sequences[:, t - 1] == tokenizer.eos_token_id
            vocab_mask = torch.zeros_like(logits)
            vocab_mask[prev_is_eos, :, :] = float("-inf")
            vocab_mask[prev_is_eos, :, tokenizer.eos_token_id] = 0
            if ramble:
                vocab_mask[~prev_is_eos, :, tokenizer.eos_token_id] = float("-inf")
            logits = logits + vocab_mask

            probs = F.softmax(logits, dim=-1)
            next_tokens = torch.multinomial(probs.squeeze(1), 1, generator=generator)
            sequences[:, t] = next_tokens.squeeze(-1)

            if not ramble and (next_tokens == tokenizer.eos_token_id).all():
                break

        return sequences

    def _beam_search(self, image_features, image_embed, device, beam_width=4, ramble=False):
        return [(0.0, self.tokenizer.decode(self._batch_generate(image_features,1, batch_size=1)[0], skip_special_tokens=True,))] # type: ignore