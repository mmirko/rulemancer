; LEt's move to next turn
(defrule end-phase-cleanup
  ?gs <- (game-state (phase end) (active-player ?ap) (turn-number ?tn))
  =>
  (modify ?gs 
    (phase untap) 
    (turn-number (+ ?tn 1))
    (active-player (other-player ?ap))
    (priority-player (other-player ?ap))))

; Reset attacking/blocking status
(defrule end-phase-reset-combat
  ?gs <- (game-state (phase end))
  ?perm <- (permanent (attacking yes))
  =>
  (modify ?perm (attacking no) (blocking no) (blocked-attacker none)))

(defrule end-phase-reset-blocking
  ?gs <- (game-state (phase end))
  ?perm <- (permanent (blocking yes))
  =>
  (modify ?perm (blocking no) (blocked-attacker none)))

; Remove damages
(defrule end-phase-remove-damage
  ?gs <- (game-state (phase end))
  ?c <- (card (zone battlefield) (damage ?d&:(> ?d 0)))
  =>
  (modify ?c (damage 0)))