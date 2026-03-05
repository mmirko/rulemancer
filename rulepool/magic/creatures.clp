; Validate: wrong player
(defrule cast-creature-wrong-player
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&main1|main2) (priority-player ?pp))
  ?action <- (cast-creature (player ?p&:(neq ?p ?pp)) (card-id ?cid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Not your priority."))))

; Validate: wrong phase
(defrule cast-creature-wrong-phase
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&~main1&~main2))
  ?action <- (cast-creature (player ?p) (card-id ?cid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Can only cast creatures in main phase."))))

; Validate: card not in hand
(defrule cast-creature-not-in-hand
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (cast-creature (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (zone ?z&~hand))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Card not in hand."))))

; Validate: not enough mana
(defrule cast-creature-not-enough-mana
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (cast-creature (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (mana-cost ?mc))
  ?ps <- (player-state (player-id ?p) (current-mana ?cm&:(< ?cm ?mc)))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Not enough mana."))))

; Validate: not a creature
(defrule cast-creature-not-a-creature
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (cast-creature (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (type ?t&~creature))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Card is not a creature."))))

; Valid cast creature
(defrule cast-creature-valid
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&main1|main2) (priority-player ?p))
  ?action <- (cast-creature (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (type creature) (zone hand) (owner ?p) (mana-cost ?mc))
  ?ps <- (player-state (player-id ?p) (current-mana ?cm&:(>= ?cm ?mc)))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Creature cast successfully.")))
  (modify ?c (zone battlefield))
  (modify ?ps (current-mana (- ?cm ?mc)))
  (assert (permanent (card-id ?cid) (controller ?p) (tapped no) (summoning-sick yes) (attacking no) (blocking no) (blocked-attacker none))))
