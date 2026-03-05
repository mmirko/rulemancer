; Pass priority - active player passes in main phases
(defrule handle-pass-priority-active-player
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase ?phase&main1|main2) (active-player ?ap) (priority-player ?ap))
  ?pass <- (pass-priority (player ?ap))
  ?ps <- (player-state (player-id ?ap))
  =>
  (retract ?pass)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Priority passed.")))
  ; Give priority to other player
  (modify ?gs (priority-player (other-player ?ap))))

; When both players pass, advance phase
(defrule both-players-pass-main1-to-combat
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase main1) (active-player ?ap) (priority-player ?opp&:(eq ?opp (other-player ?ap))))
  ?pass <- (pass-priority (player ?opp))
  =>
  (retract ?pass)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Moving to combat.")))
  (modify ?gs (phase combat-declare-attackers) (priority-player ?ap)))

(defrule both-players-pass-main2-to-end
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?gs <- (game-state (phase main2) (active-player ?ap) (priority-player ?opp&:(eq ?opp (other-player ?ap))))
  ?pass <- (pass-priority (player ?opp))
  =>
  (retract ?pass)
  (retract ?ar)
  (assert (action-result (valid yes) (reason "Moving to end phase.")))
  (modify ?gs (phase end)))