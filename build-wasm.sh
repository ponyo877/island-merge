#!/bin/bash

# Build the WebAssembly binary
GOOS=js GOARCH=wasm go build -o web/wasm/game.wasm cmd/game/main.go

# Copy the wasm_exec.js support file from Go installation
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/

echo "WebAssembly build complete!"
echo "Files created:"
echo "  - web/wasm/game.wasm"
echo "  - web/wasm_exec.js"