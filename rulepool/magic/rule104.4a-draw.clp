; Check for draw (both players at 0) rule 104.4a
(defrule game-draw
  ?gs <- (game-state (phase ?phase&~game-over))
  ?ps1 <- (player-state (player-id p1) (life ?life1&:(<= ?life1 0)))
  ?ps2 <- (player-state (player-id p2) (life ?life2&:(<= ?life2 0)))
  =>
  (modify ?gs (phase game-over))
  (assert (winner (player draw))))