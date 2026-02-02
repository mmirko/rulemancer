# Rulemancer Game Setup Guide

This guide explains how to set up CLIPS rule files to enable a game to work with Rulemancer.

## Overview

Rulemancer uses CLIPS (C Language Integrated Production System) to define game logic. To integrate a new game, you need to create a `.clp` file that defines:

1. **Game Configuration** - Basic metadata about your game
2. **Game Interface** - How the game interacts with external systems through assertables, results, and queryables

## Required Schema

Your game file must respect the schema defined in `rulepool/common.clp`. This schema defines three core templates:

### 1. `game-config`
Provides basic metadata about your game.

**Structure:**
```clips
(deftemplate game-config
  (slot game-name)
  (slot description))
```

### 2. `assertable`
Defines facts that can be asserted into the CLIPS environment from external sources (e.g., player moves, game actions).

**Structure:**
```clips
(deftemplate assertable
  (slot name)
  (multislot relations))
```

### 3. `results`
Defines facts that are generated as results of game logic and should be returned to external systems.

**Structure:**
```clips
(deftemplate results
  (slot name)
  (multislot relations))
```

### 4. `queryable`
Defines facts that can be queried from the CLIPS environment (e.g., game state, winner information).

**Structure:**
```clips
(deftemplate queryable
  (slot name)
  (multislot relations))
```

## Step-by-Step Setup

### Step 1: Create Your Game File

Create a new `.clp` file in the `rulepool/` directory, e.g., `rulepool/yourgamemeta.clp`.

### Step 2: Define Game Configuration

Use `deffacts` to declare your game configuration:

```clips
(deffacts yourgame-config
  (game-config
    (game-name YourGame)
    (description "A description of your game and its rules.")))
```

**Fields:**
- `game-name`: Identifier for your game (no spaces, use CamelCase)
- `description`: Human-readable description of the game

### Step 3: Define Game Interface

Use `deffacts` to declare how your game interfaces with the outside world:

```clips
(deffacts yourgame-interface
  ; Define what can be asserted
  (assertable
    (name action-type-1)
    (relations relation-name-1))
  
  ; Define what results are produced
  (results 
    (name result-type-1)
    (relations result-relation-1))
  
  ; Define what can be queried
  (queryable
    (name query-type-1)
    (relations relation-name-1 relation-name-2))
)
```

#### Assertables

Assertables define **inputs** to your game - actions or information that can be pushed into the CLIPS environment:

- `name`: The type of fact that can be asserted
- `relations`: The relation name(s) used in the actual CLIPS facts

**Example:**
```clips
(assertable
  (name move)
  (relations move))
```

This means external systems can assert facts like `(move player row col)` into CLIPS.

#### Results

Results define **outputs** from your game logic - facts that should be returned after processing:

- `name`: The type of result fact
- `relations`: The relation name(s) that will be matched in results

**Example:**
```clips
(results 
  (name move)
  (relations last-move))
```

This means the system will look for facts like `(last-move player row col)` to return as results.

#### Queryables

Queryables define what information can be **queried** from the current game state:

- `name`: The type of query
- `relations`: The relation name(s) that can be queried

**Example:**
```clips
(queryable
  (name winner)
  (relations winner cell))
```

This means external systems can query for facts matching `(winner player)` or `(cell row col value)`.

## Complete Example: Tic-Tac-Toe

Here's the complete metadata file for Tic-Tac-Toe (`rulepool/tictactoemeta.clp`):

```clips
(deffacts tictactoe-config
  (game-config
    (game-name TicTacToe)
    (description "A simple Tic Tac Toe game between two players.")))

(deffacts tictactoe-interface
  (assertable
    (name move)
    (relations move))
  (results 
    (name move)
    (relations last-move))
  (queryable
    (name winner)
    (relations winner cell))
  (queryable
    (name cell)
    (relations cell)))
```

### What This Means:

1. **Game Config**: Identifies the game as "TicTacToe"

2. **Assertable `move`**: External systems can assert move facts like:
   ```clips
   (move X 1 1)  ; Player X moves to position (1,1)
   ```

3. **Results `move`**: After processing, the system returns facts like:
   ```clips
   (last-move X 1 1)  ; The last move was X at (1,1)
   ```

4. **Queryable `winner`**: Can query for winner status:
   ```clips
   (winner X)  ; X has won
   ```

5. **Queryable `cell`**: Can query the board state:
   ```clips
   (cell 1 1 X)  ; Cell at (1,1) contains X
   (cell 1 2 O)  ; Cell at (1,2) contains O
   ```

## Best Practices

1. **Naming Conventions**:
   - Use descriptive names for `game-name` (CamelCase, no spaces)
   - Keep relation names short but meaningful
   - Be consistent with naming across your game files

2. **Separation of Concerns**:
   - Keep metadata (interface definitions) in a separate `*meta.clp` file
   - Keep game logic (rules, templates, functions) in separate files

3. **Documentation**:
   - Always provide a clear `description` in your game config
   - Comment your code to explain complex rules

4. **Testing**:
   - Test that assertables work by asserting facts and checking the results
   - Test that queryables return expected game state
   - Verify results are properly generated

## Next Steps

After creating your metadata file:

1. **Create Game Logic Files**: Write CLIPS rules that implement your game logic
2. **Define Game Templates**: Create templates for your game-specific facts
3. **Implement Rules**: Write rules that respond to assertions and update game state
4. **Test**: Use the Rulemancer test framework to verify your game works correctly

## Additional Resources

- CLIPS Documentation: Learn more about CLIPS syntax and features
- Example Games: Check `rulepool/` directory for complete game implementations
- Rulemancer Documentation: See main README.md for system architecture

---

For questions or issues, please refer to the main project documentation or open an issue on the project repository.
