; Going from upkeep to draw phase, it is a sort of placeholder for upkeep triggers
(defrule upkeep-to-draw
  ?gs <- (game-state (phase upkeep))
  =>
  (modify ?gs (phase draw)))