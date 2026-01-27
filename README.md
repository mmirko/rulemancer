# Rulemancer

![Logo](./logo.png)

A Go application that embeds the CLIPS expert system engine to power rules-based games. Define game logic using CLIPS (an expressive rule/fact inference engine), then interact with it via HTTP or CLI.

## Features

- **CLIPS Integration**: Leverage CLIPS for complex rule-based inference and fact management
- **HTTP API**: Serve games via REST endpoints with optional TLS support
- **Fact/Rule Management**: Assert facts, query game state, and receive results from inference
- **CLI Tools**: Commands for testing, rebuilding, and serving games
- **Configuration**: JSON-based configuration for game setup, facts, and query rules

## Quick Start

### Prerequisites

- Go 1.25+
- C compiler (for CLIPS 6.4 compilation)

### Installation

```bash
./install-clips.sh
go build
```

### Commands

- `./rulemancer serve` - Start HTTP server
- `./rulemancer test` - Run test suite

## Project Structure

- **`core/`** - CLIPS C source files
- **`cmd/`** - CLI commands (serve, test, root)
- **`pkg/game/`** - Specific game logic and handlers
- **`pkg/rulemancer/`** - Core engine bindings
- **`rulepool/`** - CLIPS rule files (`.clp`)
- **`examples/`** - Usage examples

## Configuration

Create the rule pool directory:

```bash
mkdir rulepool
```

Place your CLIPS rule files (e.g., `tictactoe.clp`) in the `rulepool/` directory.

Edit `rulemancer.json`:

```json
{
  "debug": true,
  "tls_cert_file": "server.crt",
  "tls_key_file": "server.key",
  "rule_pool": "rulepool",
  "assertables": ["move"],
  "results": { "move": ["last-move"] },
  "querables": ["cell", "winner"]
}
```

- **assertables**: Fact types the client can assert
- **querables**: Fact types the client can query
- **results**: Return values from assertions

## Example: Tic-Tac-Toe

See `rulepool/tictactoe.clp` for a complete game implementation using CLIPS rules and facts.

## License

See LICENSE file
