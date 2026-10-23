#!/bin/bash

# Built directly from WAT, no wasm-opt, so that the opcodes stay exactly as written.

set -e

cd "$(dirname "$0")"
mkdir -p output
wat2wasm bulk-memory.wat -o output/bulk-memory.wasm
