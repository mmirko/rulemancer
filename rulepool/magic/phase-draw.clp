; Automatically move to main1
(defrule draw-to-main1
  ?gs <- (game-state (phase draw))
  =>
  (modify ?gs (phase main1)))
