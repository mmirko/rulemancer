; Validate: not defending player
(defrule declare-blocker-not-defending-player
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase combat-declare-blockers) (active-player ?ap))
  ?action <- (declare-blocker (player ?p&:(eq ?p ?ap)) (blocker-id ?bid) (attacker-id ?aid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Only defending player can declare blockers."))))

; Validate: wrong phase
(defrule declare-blocker-wrong-phase
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&~combat-declare-blockers))
  ?action <- (declare-blocker (player ?p) (blocker-id ?bid) (attacker-id ?aid))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Not in declare blockers phase."))))

; Validate: blocker tapped
(defrule declare-blocker-creature-tapped
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (declare-blocker (player ?p) (blocker-id ?bid) (attacker-id ?aid))
  ?perm <- (permanent (card-id ?bid) (tapped yes))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Blocker is tapped."))))

; Validate: attacker not attacking
(defrule declare-blocker-not-attacking
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?action <- (declare-blocker (player ?p) (blocker-id ?bid) (attacker-id ?aid))
  ?perm <- (permanent (card-id ?aid) (attacking no))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid no) (reason "Target is not attacking."))))

; Valid declare blocker
(defrule declare-blocker-valid
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase combat-declare-blockers) (active-player ?ap))
  ?action <- (declare-blocker (player ?p&:(neq ?p ?ap)) (blocker-id ?bid) (attacker-id ?aid))
  ?blocker-card <- (card (card-id ?bid) (type creature) (zone battlefield))
  ?blocker <- (permanent (card-id ?bid) (controller ?p) (tapped no))
  ?attacker <- (permanent (card-id ?aid) (attacking yes))
  =>
  (retract ?action)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Blocker declared.")))
  (modify ?blocker (blocking yes) (blocked-attacker ?aid)))

; Defending player passes - move to damage
(defrule declare-blockers-done
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase combat-declare-blockers) (active-player ?ap) (priority-player ?dp&:(neq ?dp ?ap)))
  ?pass <- (pass-priority (player ?dp))
  =>
  (retract ?pass)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Moving to combat damage.")))
  (modify ?gs (phase combat-damage)))
