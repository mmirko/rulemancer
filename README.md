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

- Git
- Go 1.25+
- C compiler (for CLIPS 6.4 compilation)
- re2c with go-bindings (for the build subcommand)

### Clone the Repository

```bash
git clone https://github.com/mmirko/rulemancer.git
cd rulemancer
```

### Project Structure

Inside the project directory, you'll find the following key folders:

- **`core/`** - CLIPS 6.4 C source files and headers (not included in the repo, install using `install-clips.sh`)
- **`cmd/`** - CLI commands (serve, test, build, root)
- **`pkg/rulemancer/`** - Core engine, CLIPS bindings, HTTP handlers, and game management
- **`rulepool/`** - CLIPS rule files (`.clp`) for game definitions (e.g., Tic-Tac-Toe example)
- **`interface/`** - Client interface examples and utilities (builded via `rulemancer build`)
- **`testpool/`** - Test rule files for development (unit tests for Tic-Tac-Toe game logic)

### Installation

To install dependencies, compile CLIPS, and build Rulemancer, run from the project root:

```bash
./install-clips.sh
make
```

The above commands will compile CLIPS and build the Rulemancer binary placed in the project root: `./rulemancer`

### 

### Basic Usage

- `./rulemancer test` - Run test suite
- `./rulemancer build` - Build the extras tools
- `./rulemancer serve` - Start HTTPS server (listens on :3000 with TLS)

Once the server is running, the API can be accessed by clients at `https://localhost:3000/api/v1/`

To ease the interaction, the `rulemancer build` command generates shell client interfaces in the `interface/` folder for all available games.

The [Rooms and Games](README-ROOMS-AND-GAMES.md) document provides detailed information on how to create and manage game rooms and interact with games via the API. The API endpoints are documented in the [API Endpoints](README-API.md) document.

The `rulemancer.json` configuration file can be edited to customize server settings.

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

## Example: Tic-Tac-Toe

See [rulepool/tictactoe.clp](rulepool/tictactoe.clp) and [rulepool/tictactoemeta.clp](rulepool/tictactoemeta.clp) for a complete game implementation using CLIPS rules and facts.

## License

See LICENSE file
