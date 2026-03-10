; Player loses when life <= 0
(defrule player-lost-no-life
  ?gs <- (game-state (phase ?phase&~game-over))
  ?ps <- (player-state (player-id ?p) (life ?life&:(<= ?life 0)))
  =>
  (modify ?gs (phase game-over))
  (assert (winner (player (other-player ?p)))))