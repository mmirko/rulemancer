# Tressette Ruleset for Rulemancer

This folder contains a full 4-player partnership Tressette implementation for Rulemancer.

## Rules Modeled

This CLIPS implementation follows the standard partnership rules from:

- https://www.pagat.com/tressette/tressette.html
- https://en.wikipedia.org/wiki/Tressette
- https://it.wikipedia.org/wiki/Tressette

Implemented core rules:

- 4 players in fixed partnerships: `p1+p3` vs `p2+p4`
- 40-card Italian deck (denari, coppe, spade, bastoni)
- No trump suit
- Rank order (high to low): `3 2 A K Knight Knave 7 6 5 4`
- Must follow suit if possible
- Trick winner: highest rank in lead suit
- 10 tricks per hand
- Base card scoring:
  - Ace = 1 point
  - 3, 2, king, knight, knave = 1/3 point
  - 7, 6, 5, 4 = 0
- Last trick = +1 point
- Team hand points are computed by flooring thirds (`div thirds 3`)
- Match target score = 21 points
- Accuse support in first trick:
  - `napoli` (A,2,3 same suit) = +3 points
  - `tris` of aces/twos/threes = +3 points
  - four of a kind (aces/twos/threes) = +4 points

## Rulemancer Interface

### Assertables

- `play-card`
- `declare-accuse`

### Results

- `action-result`
- `trick-result`
- `hand-result`
- `accuse-result`
- `match-winner`

### Queryables

- `state`
- `hand`
- `table`
- `score`
- `winner`

## Notes

- The engine deals cards deterministically from the predefined deck order.
- Dealer rotates each hand.
- First lead is the player after the dealer.
- Historical special win conditions (cappotto, cappottone, stramazzo variants) are intentionally not modeled as separate multi-game outcomes.
