; Blocked attacker deals damage to blocker
(defrule combat-damage-blocked-attacker-to-blocker
  ?gs <- (game-state (phase combat-damage))
  ?attacker-perm <- (permanent (card-id ?aid) (attacking yes))
  ?attacker-card <- (card (card-id ?aid) (power ?ap) (zone battlefield))
  ?blocker-perm <- (permanent (card-id ?bid) (blocking yes) (blocked-attacker ?aid))
  ?blocker-card <- (card (card-id ?bid) (damage ?bd) (zone battlefield))
  =>
  (modify ?blocker-card (damage (+ ?bd ?ap))))

; Blocker deals damage to attacker
(defrule combat-damage-blocker-to-attacker
  ?gs <- (game-state (phase combat-damage))
  ?blocker-perm <- (permanent (card-id ?bid) (blocking yes) (blocked-attacker ?aid))
  ?blocker-card <- (card (card-id ?bid) (power ?bp) (zone battlefield))
  ?attacker-perm <- (permanent (card-id ?aid) (attacking yes))
  ?attacker-card <- (card (card-id ?aid) (damage ?ad) (zone battlefield))
  =>
  (modify ?attacker-card (damage (+ ?ad ?bp))))

; Unblocked attacker deals damage to defending player
(defrule combat-damage-unblocked-attacker-to-player
  ?gs <- (game-state (phase combat-damage) (active-player ?ap))
  ?attacker-perm <- (permanent (card-id ?aid) (attacking yes))
  ?attacker-card <- (card (card-id ?aid) (power ?power) (zone battlefield))
  (not (permanent (blocking yes) (blocked-attacker ?aid)))
  ?defender <- (player-state (player-id ?dp&:(eq ?dp (other-player ?ap))) (life ?life))
  =>
  (modify ?defender (life (- ?life ?power))))

; Check for dead creatures (damage >= toughness)
(defrule creature-dies-from-damage
  ?gs <- (game-state (phase combat-damage))
  ?c <- (card (card-id ?cid) (type creature) (zone battlefield) (toughness ?t) (damage ?d&:(>= ?d ?t)))
  ?perm <- (permanent (card-id ?cid))
  =>
  (retract ?perm)
  (modify ?c (zone graveyard) (damage 0)))

; After damage, move to main2
(defrule combat-damage-to-main2
  ?gs <- (game-state (phase combat-damage) (active-player ?ap))
  =>
  (modify ?gs (phase main2) (priority-player ?ap)))
