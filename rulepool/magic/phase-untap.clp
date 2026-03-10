; Rule 502.3 - Automatically untap all permanents
(defrule untap-phase-untap-permanents
  ?gs <- (game-state (phase untap) (active-player ?ap))
  ?perm <- (permanent (controller ?ap) (tapped yes))
  =>
  (modify ?perm (tapped no)))

; Remove summoning sickness from creatures
(defrule untap-phase-remove-summoning-sickness
  ?gs <- (game-state (phase untap) (active-player ?ap))
  ?perm <- (permanent (controller ?ap) (summoning-sick yes))
  =>
  (modify ?perm (summoning-sick no)))

; Refill mana to max
(defrule untap-phase-refill-mana
  ?gs <- (game-state (phase untap) (active-player ?ap))
  ?ps <- (player-state (player-id ?ap) (max-mana ?mm))
  =>
  (modify ?ps (current-mana ?mm) (lands-played 0)))

; Move to upkeep phase
(defrule untap-to-upkeep
  ?gs <- (game-state (phase untap) (active-player ?ap))
  =>
  (modify ?gs (phase upkeep)))