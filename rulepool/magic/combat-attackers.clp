; Validate: not active player
(defrule declare-attacker-not-active-player
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase combat-declare-attackers) (active-player ?ap))
  ?action <- (declare-attacker (player ?p&:(neq ?p ?ap)) (card-id ?cid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Only active player can declare attackers."))))

; Validate: wrong phase
(defrule declare-attacker-wrong-phase
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&~combat-declare-attackers))
  ?action <- (declare-attacker (player ?p) (card-id ?cid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Not in declare attackers phase."))))

; Validate: creature tapped
(defrule declare-attacker-creature-tapped
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (declare-attacker (player ?p) (card-id ?cid))
  ?perm <- (permanent (card-id ?cid) (tapped yes))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Creature is tapped."))))

; Validate: summoning sickness
(defrule declare-attacker-summoning-sick
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (declare-attacker (player ?p) (card-id ?cid))
  ?perm <- (permanent (card-id ?cid) (summoning-sick yes))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Creature has summoning sickness."))))

; Valid declare attacker
(defrule declare-attacker-valid
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase combat-declare-attackers) (active-player ?p))
  ?action <- (declare-attacker (player ?p) (card-id ?cid))
  ?c <- (card (card-id ?cid) (type creature) (zone battlefield))
  ?perm <- (permanent (card-id ?cid) (controller ?p) (tapped no) (summoning-sick no))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Attacker declared.")))
  (modify ?perm (attacking yes) (tapped yes)))

; Active player passes - move to declare blockers
(defrule declare-attackers-done
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase combat-declare-attackers) (active-player ?ap))
  ?pass <- (pass-priority (player ?ap))
  =>
  (retract ?pass)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Moving to declare blockers.")))
  (modify ?gs (phase combat-declare-blockers) (priority-player (other-player ?ap))))
