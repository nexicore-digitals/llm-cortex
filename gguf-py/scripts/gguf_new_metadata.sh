#!/usr/bin/env bash

# --- Configuration Section: EDIT THESE VALUES ---
# 1. SET THE NEW MODEL NAME
NEW_NAME="Mistral-7B-Cortex-V1.0"

# 2. SET THE NEW MODEL DESCRIPTION
NEW_DESCRIPTION="Mistral 7B v0.1 model, Q4_K_M. Optimized for use with the LLM Cortex Go inference engine."
# ------------------------------------------------

# --- Core Script Logic: DO NOT EDIT BELOW THIS LINE ---

# Define the absolute path to the Python script for reliable execution
GGUF_SCRIPT="/home/owen/projects/llm-cortex/gguf-py/scripts/gguf_new_metadata.py"

# Define the input file path (relative to where the script is run)
INPUT_FILE="../../models/mistral-7b/mistral-7b-v0.1.Q4_K_M.gguf"

# Define the output file name, using the NEW_NAME variable for clarity
# This uses sed to replace spaces with hyphens for a clean filename
OUTPUT_FILE="../../models/$(echo "$NEW_NAME" | sed 's/ /-/g').gguf"

# The chat template remains static
CHAT_TEMPLATE="{% for message in messages %}{% if message['role'] == 'user' %}[INST] {{ message['content'] }} [/INST]{% elif message['role'] == 'assistant' %}{{ message['content'] }}{% endif %}{% endfor %}"

# --------------------

echo "--- GGUF Metadata Update ---"
echo "  Input File:  ${INPUT_FILE}"
echo "  Output File: ${OUTPUT_FILE}"
echo "  New Name:    ${NEW_NAME}"
echo "  New Desc:    ${NEW_DESCRIPTION}"
echo "----------------------------"

# --- Execute the Python Script ---
python3 "${GGUF_SCRIPT}" \
    "${INPUT_FILE}" \
    "${OUTPUT_FILE}" \
    --chat-template "${CHAT_TEMPLATE}" \
    --general-name "${NEW_NAME}" \
    --general-description "${NEW_DESCRIPTION}" \
    --force

echo "GGUF file successfully created/updated at: ${OUTPUT_FILE}"