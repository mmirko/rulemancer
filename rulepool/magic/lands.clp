; Validate: wrong player
(defrule play-land-wrong-player
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&main1|main2) (priority-player ?pp))
  ?action <- (play-land (player ?p&:(neq ?p ?pp)) (card-id ?cid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Not your priority."))))

; Validate: wrong phase
(defrule play-land-wrong-phase
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&~main1&~main2))
  ?action <- (play-land (player ?p) (card-id ?cid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Can only play lands in main phase."))))

; Validate: card not in hand
(defrule play-land-not-in-hand
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (play-land (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (zone ?z&~hand))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Card not in hand."))))

; Validate: already played land this turn
(defrule play-land-already-played
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (play-land (player ?p) (card-id ?cid))
  ?ps <- (player-state (player-id ?p) (lands-played ?lp&:(>= ?lp 1)))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Already played a land this turn."))))

; Validate: card is not a land
(defrule play-land-not-a-land
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (play-land (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (type ?t&~land))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Card is not a land."))))

; Valid play land
(defrule play-land-valid
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&main1|main2) (priority-player ?p))
  ?action <- (play-land (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (type land) (zone hand) (owner ?p))
  ?ps <- (player-state (player-id ?p) (lands-played ?lp&:(< ?lp 1)) (max-mana ?mm))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Land played successfully.")))
  (modify ?c (zone battlefield))
  (modify ?ps (lands-played (+ ?lp 1)) (max-mana (+ ?mm 1)) (current-mana (+ ?mm 1)))
  (assert (permanent (card-id ?cid) (controller ?p) (tapped no) (summoning-sick no) (attacking no) (blocking no) (blocked-attacker none))))
