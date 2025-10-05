#!/bin/bash
MODEL="./models/deepseek-7b/deepseek-llm-7b-chat.Q4_K_M.gguf"
PROMPT="Hello what is your name?"
THREADS=8
TOKENS=128

echo Deepseek thinking...
go run . --model "$MODEL" --prompt "$PROMPT" --threads "$THREADS" --n_predict "$TOKENS" 2>&1 \
| awk '/^Assistant:/ {sub(/^Assistant:/, "", $0); reply=$0} END {print reply}'
