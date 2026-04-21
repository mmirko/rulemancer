(deftemplate game-state
  (slot phase)          ; setup | playing | hand-ended | game-ended
  (slot hand-number)
  (slot dealer)
  (slot leader)
  (slot current-player)
  (slot trick-number)
  (slot target-score))

(deftemplate deck-card
  (slot position)
  (slot suit)
  (slot rank))

(deftemplate dealing
  (slot hand))

(deftemplate dealt-marker
  (slot hand)
  (slot position))

(deftemplate card-in-hand
  (slot hand)
  (slot player)
  (slot suit)
  (slot rank))

(deftemplate initial-hand-card
  (slot hand)
  (slot player)
  (slot suit)
  (slot rank))

(deftemplate trick-state
  (slot hand)
  (slot trick)
  (slot cards-played)
  (slot lead-suit)
  (slot winning-player)
  (slot winning-rank))

(deftemplate trick-card
  (slot hand)
  (slot trick)
  (slot order)
  (slot player)
  (slot suit)
  (slot rank))

(deftemplate hand-tally
  (slot hand)
  (slot team)
  (slot thirds)
  (slot tricks))

(deftemplate team-score
  (slot team)
  (slot points))

(deftemplate match-winner
  (slot team))

(deftemplate player-team
  (slot player)
  (slot team))

(deftemplate accuse-used
  (slot hand)
  (slot player)
  (slot type)
  (slot key))

(deftemplate play-card
  (slot player)
  (slot suit)
  (slot rank))

(deftemplate declare-accuse
  (slot player)
  (slot type) ; napoli | tris
  (slot suit (default none)) ; needed for napoli
  (slot rank (default none))) ; needed for tris: ace | two | three

(deftemplate action-result
  (slot valid)
  (slot code)
  (slot reason))

(deftemplate trick-result
  (slot hand)
  (slot trick)
  (slot winner)
  (slot team)
  (slot points-third))

(deftemplate hand-result
  (slot hand)
  (slot team1-points)
  (slot team2-points)
  (slot team1-total)
  (slot team2-total))

(deftemplate accuse-result
  (slot hand)
  (slot player)
  (slot team)
  (slot type)
  (slot key)
  (slot points))
