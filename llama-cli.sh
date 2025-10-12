#!/usr/bin/env bash

./bin/llama-cli -m "$1" --batch-size 4096 --threads 8 --no-mmap --jinja