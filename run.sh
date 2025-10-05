#!/bin/bash
MODEL="./models/deepseek-7b/deepseek-llm-7b-chat.Q4_K_M.gguf"
PROMPT="Hello what is your name?"
THREADS=8
TOKENS=128

go run . --model "$MODEL" --prompt "$PROMPT" --threads "$THREADS" --n_predict "$TOKENS"
