# Island Merge

A browser-based puzzle game where you connect islands by building bridges. Built with Go and WebAssembly.

## How to Play

- Click on sea tiles adjacent to islands or bridges to build connections
- Connect all islands to win!
- Try to complete the puzzle in the minimum number of moves

## Building and Running

### Prerequisites

- Go 1.16 or higher
- Python 3 (for the local web server)

### Build

```bash
./build-wasm.sh
```

### Run

```bash
./serve.sh
```

Then open your browser to http://localhost:8080

## Project Structure

- `cmd/game/` - Main entry point
- `pkg/core/` - Core game loop and world state
- `pkg/island/` - Game logic (board, tiles, Union-Find)
- `pkg/systems/` - Input and rendering systems
- `web/` - HTML and WebAssembly files

## Features

- MVP Implementation:
  - 5x5 grid with 3 islands
  - Click to build bridges on adjacent sea tiles
  - Union-Find algorithm for connectivity checking
  - Victory detection when all islands are connected
  - Move counter
  - Simple colored tile rendering