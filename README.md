# Rulemancer

<p align="center">
  <img src="logo.png" alt="Rulemancer Logo"/>
</p>

A Go application that embeds the CLIPS expert system engine to power rules-based games. Define game logic using CLIPS (an expressive rule/fact inference engine), then interact with it via HTTP or CLI.

## Features

- **CLIPS Integration**: Leverage CLIPS for complex rule-based inference and fact management
- **Multi-Game Support**: Host multiple game types simultaneously with dynamic game loading
- **Room-Based Multiplayer**: Create isolated game rooms for concurrent sessions
- **HTTP API**: Comprehensive REST endpoints for system, game, and room management with TLS support
- **Flexible Configuration**: JSON-based configuration with support for multiple game definitions
- **CLI Tools**: Commands for testing, building, and serving games
- **Client Management**: Track and manage connected clients per room

## Quick Start

### Prerequisites

- Go 1.25+
- C compiler (for CLIPS 6.4 compilation)
- re2c with go-bindings (for the build subcommand)

### Installation

```bash
git clone https://github.com/mmirko/rulemancer.git
cd rulemancer
./install-clips.sh
make
```
The above commands will compile CLIPS and build the Rulemancer binary placed in the project root: `./rulemancer`

### 

### Usage

- `./rulemancer serve` - Start HTTP server (listens on :3000 with TLS)
- `./rulemancer test` - Run test suite
- `./rulemancer build` - Build the extras tools

## Project Structure

- **`core/`** - CLIPS 6.4 C source files and headers
- **`cmd/`** - CLI commands (serve, test, build, root)
- **`pkg/rulemancer/`** - Core engine, CLIPS bindings, HTTP handlers, and game management
- **`rulepool/`** - CLIPS rule files (`.clp`) for game definitions (e.g., Tic-Tac-Toe example)
- **`interface/`** - Client interface examples and utilities
- **`testpool/`** - Test rule files for development (unit tests for Tic-Tac-Toe game logic)

## Configuration

Edit `rulemancer.json`:

```json
{
  "debug": true,
  "debug_level": 10,
  "tls_cert_file": "server.crt",
  "tls_key_file": "server.key",
  "clipsless_mode": false,
  "games": ["rulepool"]
}
```

### Configuration Options

- **debug**: Enable debug logging
- **debug_level**: Verbosity level for debugging (0-10)
- **tls_cert_file**: Path to TLS certificate file
- **tls_key_file**: Path to TLS private key file
- **clipsless_mode**: Run without CLIPS for testing purposes
- **games**: Array of game directories to load

### Game Definition

Each game directory should contain CLIPS files with:

- **Game metadata** via `game-config` fact:
  ```clips
  (game-config 
    (game-name "TicTacToe")
    (description "Classic 3x3 grid game"))
  ```
- **Assertable facts**: Facts that can be asserted by clients
- **Queryable facts**: Facts that can be queried by clients  
- **Response facts**: Facts returned after assertions
- **Game rules**: CLIPS rules implementing game logic

for more details check the [Game Definition](README-GAME-DEFINITION.md) document.

See [rulepool/tictactoe.clp](rulepool/tictactoe.clp) for a complete example.


## Example: Tic-Tac-Toe

See [rulepool/tictactoe.clp](rulepool/tictactoe.clp) and [rulepool/tictactoemeta.clp](rulepool/tictactoemeta.clp) for a complete game implementation using CLIPS rules and facts.

### Creating a Room

```bash
curl -k -X POST https://localhost:3000/api/v1/room/create \
  -H "Content-Type: application/json" \
  -d '{"name": "My Game", "description": "Test room", "game_ref": "tictactoe"}'
```
See `rulepool/tictactoe.clp` for a complete game implementation using CLIPS rules and facts.

## License

See LICENSE file
