(defrule clear-previous-results-on-play
  ?pc <- (play-card)
  ?ar <- (action-result)
  =>
  (retract ?ar))

(defrule clear-previous-results-on-accuse
  ?da <- (declare-accuse)
  ?ar <- (action-result)
  =>
  (retract ?ar))

(defrule deal-each-card
  ?d <- (dealing (hand ?h))
  ?gs <- (game-state (phase setup) (hand-number ?h) (dealer ?dealer))
  ?dc <- (deck-card (position ?p) (suit ?s) (rank ?r))
  (not (dealt-marker (hand ?h) (position ?p)))
  =>
  (bind ?owner (player-for-deal-position ?dealer ?p))
  (assert (dealt-marker (hand ?h) (position ?p)))
  (assert (card-in-hand (hand ?h) (player ?owner) (suit ?s) (rank ?r)))
  (assert (initial-hand-card (hand ?h) (player ?owner) (suit ?s) (rank ?r))))

(defrule start-first-trick-when-deal-complete
  ?d <- (dealing (hand ?h))
  ?gs <- (game-state (phase setup) (hand-number ?h) (dealer ?dealer))
  (test (= (length$ (find-all-facts ((?m dealt-marker)) (eq ?m:hand ?h))) 40))
  =>
  (retract ?d)
  (bind ?leader (next-player ?dealer))
  (modify ?gs (phase playing) (leader ?leader) (current-player ?leader) (trick-number 1))
  (assert (trick-state
            (hand ?h)
            (trick 1)
            (cards-played 0)
            (lead-suit none)
            (winning-player none)
            (winning-rank four)))
  (assert (hand-tally (hand ?h) (team team1) (thirds 0) (tricks 0)))
  (assert (hand-tally (hand ?h) (team team2) (thirds 0) (tricks 0)))
  (assert (action-result (valid yes) (code hand-started) (reason "Hand dealt. Play begins."))))

(defrule reject-play-when-game-not-playing
  (declare (salience 100))
  ?pc <- (play-card (player ?p))
  (game-state (phase ?phase&:(neq ?phase playing)))
  =>
  (retract ?pc)
  (assert (action-result (valid no) (code wrong-phase) (reason "Cannot play a card outside the playing phase."))))

(defrule reject-play-out-of-turn
  (declare (salience 90))
  ?pc <- (play-card (player ?p))
  (game-state (phase playing) (current-player ?current&:(neq ?current ?p)))
  =>
  (retract ?pc)
  (assert (action-result (valid no) (code wrong-turn) (reason "It is not this player's turn."))))

(defrule reject-play-card-not-owned
  (declare (salience 80))
  ?pc <- (play-card (player ?p) (suit ?s) (rank ?r))
  (game-state (phase playing) (hand-number ?h))
  (not (card-in-hand (hand ?h) (player ?p) (suit ?s) (rank ?r)))
  =>
  (retract ?pc)
  (assert (action-result (valid no) (code card-missing) (reason "Player does not hold this card."))))

(defrule reject-play-must-follow-suit
  (declare (salience 70))
  ?pc <- (play-card (player ?p) (suit ?sPlayed))
  (game-state (phase playing) (hand-number ?h))
  (trick-state (hand ?h) (trick ?t) (cards-played ?n&:(> ?n 0)) (lead-suit ?lead&:(neq ?lead none)))
  (test (and (neq ?sPlayed ?lead) (has-suit-in-hand ?h ?p ?lead)))
  =>
  (retract ?pc)
  (assert (action-result (valid no) (code follow-suit) (reason "Player must follow the lead suit when possible."))))

(defrule accept-first-card-of-trick
  ?pc <- (play-card (player ?p) (suit ?s) (rank ?r))
  ?gs <- (game-state (phase playing) (hand-number ?h) (current-player ?p) (trick-number ?t))
  ?ts <- (trick-state (hand ?h) (trick ?t) (cards-played 0) (lead-suit none) (winning-player none) (winning-rank ?wr))
  ?ch <- (card-in-hand (hand ?h) (player ?p) (suit ?s) (rank ?r))
  =>
  (retract ?pc)
  (retract ?ch)
  (modify ?ts
    (cards-played 1)
    (lead-suit ?s)
    (winning-player ?p)
    (winning-rank ?r))
  (modify ?gs (current-player (next-player ?p)))
  (assert (trick-card (hand ?h) (trick ?t) (order 1) (player ?p) (suit ?s) (rank ?r)))
  (assert (action-result (valid yes) (code card-played) (reason "Card accepted."))))

(defrule accept-following-card-beats-current
  ?pc <- (play-card (player ?p) (suit ?s) (rank ?r))
  ?gs <- (game-state (phase playing) (hand-number ?h) (current-player ?p))
  ?ts <- (trick-state (hand ?h) (trick ?t) (cards-played ?n&:(> ?n 0)&:(< ?n 4)) (lead-suit ?lead) (winning-player ?wp) (winning-rank ?wr))
  ?ch <- (card-in-hand (hand ?h) (player ?p) (suit ?s) (rank ?r))
  (test (and (eq ?s ?lead) (> (rank-strength ?r) (rank-strength ?wr))))
  =>
  (retract ?pc)
  (retract ?ch)
  (bind ?nextCount (+ ?n 1))
  (modify ?ts
    (cards-played ?nextCount)
    (winning-player ?p)
    (winning-rank ?r))
  (modify ?gs (current-player (next-player ?p)))
  (assert (trick-card (hand ?h) (trick ?t) (order ?nextCount) (player ?p) (suit ?s) (rank ?r)))
  (assert (action-result (valid yes) (code card-played) (reason "Card accepted."))))

(defrule accept-following-card-does-not-beat
  ?pc <- (play-card (player ?p) (suit ?s) (rank ?r))
  ?gs <- (game-state (phase playing) (hand-number ?h) (current-player ?p))
  ?ts <- (trick-state (hand ?h) (trick ?t) (cards-played ?n&:(> ?n 0)&:(< ?n 4)) (lead-suit ?lead) (winning-player ?wp) (winning-rank ?wr))
  ?ch <- (card-in-hand (hand ?h) (player ?p) (suit ?s) (rank ?r))
  (test (or (neq ?s ?lead) (<= (rank-strength ?r) (rank-strength ?wr))))
  =>
  (retract ?pc)
  (retract ?ch)
  (bind ?nextCount (+ ?n 1))
  (modify ?ts (cards-played ?nextCount))
  (modify ?gs (current-player (next-player ?p)))
  (assert (trick-card (hand ?h) (trick ?t) (order ?nextCount) (player ?p) (suit ?s) (rank ?r)))
  (assert (action-result (valid yes) (code card-played) (reason "Card accepted."))))

(defrule close-trick
  ?gs <- (game-state (phase playing) (hand-number ?h) (trick-number ?t))
  ?ts <- (trick-state (hand ?h) (trick ?t) (cards-played 4) (winning-player ?winner))
  ?tc1 <- (trick-card (hand ?h) (trick ?t) (order 1) (rank ?r1))
  ?tc2 <- (trick-card (hand ?h) (trick ?t) (order 2) (rank ?r2))
  ?tc3 <- (trick-card (hand ?h) (trick ?t) (order 3) (rank ?r3))
  ?tc4 <- (trick-card (hand ?h) (trick ?t) (order 4) (rank ?r4))
  ?ht <- (hand-tally (hand ?h) (team ?team) (thirds ?thirds) (tricks ?tricks))
  (test (eq ?team (team-of ?winner)))
  =>
  (bind ?trickThirds (+ (card-point-thirds ?r1)
                        (card-point-thirds ?r2)
                        (card-point-thirds ?r3)
                        (card-point-thirds ?r4)))
  (bind ?bonus (if (= ?t 10) then 3 else 0))
  (bind ?totalThirds (+ ?thirds ?trickThirds ?bonus))
  (bind ?nextTricks (+ ?tricks 1))

  (retract ?tc1)
  (retract ?tc2)
  (retract ?tc3)
  (retract ?tc4)
  (retract ?ts)

  (modify ?ht (thirds ?totalThirds) (tricks ?nextTricks))
  (assert (trick-result
            (hand ?h)
            (trick ?t)
            (winner ?winner)
            (team ?team)
            (points-third (+ ?trickThirds ?bonus))))

  (if (< ?t 10)
    then
      (bind ?nextTrick (+ ?t 1))
      (modify ?gs
        (leader ?winner)
        (current-player ?winner)
        (trick-number ?nextTrick))
      (assert (trick-state
                (hand ?h)
                (trick ?nextTrick)
                (cards-played 0)
                (lead-suit none)
                (winning-player none)
                (winning-rank four)))
    else
      (modify ?gs (phase hand-ended) (leader ?winner) (current-player ?winner))))

(defrule score-hand
  ?gs <- (game-state (phase hand-ended) (hand-number ?h) (target-score ?target) (dealer ?dealer))
  ?t1 <- (hand-tally (hand ?h) (team team1) (thirds ?th1))
  ?t2 <- (hand-tally (hand ?h) (team team2) (thirds ?th2))
  ?s1 <- (team-score (team team1) (points ?p1))
  ?s2 <- (team-score (team team2) (points ?p2))
  =>
  (bind ?add1 (div ?th1 3))
  (bind ?add2 (div ?th2 3))
  (bind ?new1 (+ ?p1 ?add1))
  (bind ?new2 (+ ?p2 ?add2))

  (modify ?s1 (points ?new1))
  (modify ?s2 (points ?new2))

  (assert (hand-result
            (hand ?h)
            (team1-points ?add1)
            (team2-points ?add2)
            (team1-total ?new1)
            (team2-total ?new2)))

  (if (or (>= ?new1 ?target) (>= ?new2 ?target))
    then
      (modify ?gs (phase game-ended))
      (if (>= ?new1 ?target)
        then (assert (match-winner (team team1)))
        else (assert (match-winner (team team2))))
      (assert (action-result (valid yes) (code game-ended) (reason "Match finished.")))
    else
      (bind ?newHand (+ ?h 1))
      (bind ?newDealer (next-player ?dealer))
      (modify ?gs
        (phase setup)
        (hand-number ?newHand)
        (dealer ?newDealer)
        (leader none)
        (current-player none)
        (trick-number 0))
      (assert (dealing (hand ?newHand)))
      (assert (action-result (valid yes) (code next-hand) (reason "Hand scored. Next hand dealt.")))))

(defrule cleanup-old-hand-data-after-scoring
  ?gs <- (game-state (phase setup) (hand-number ?current))
  ?c <- (card-in-hand (hand ?old&:(< ?old ?current)))
  =>
  (retract ?c))

(defrule cleanup-old-initial-hand-after-scoring
  ?gs <- (game-state (phase setup) (hand-number ?current))
  ?c <- (initial-hand-card (hand ?old&:(< ?old ?current)))
  =>
  (retract ?c))

(defrule cleanup-old-markers-after-scoring
  ?gs <- (game-state (phase setup) (hand-number ?current))
  ?m <- (dealt-marker (hand ?old&:(< ?old ?current)))
  =>
  (retract ?m))

(defrule cleanup-old-tallies-after-scoring
  ?gs <- (game-state (phase setup) (hand-number ?current))
  ?h <- (hand-tally (hand ?old&:(< ?old ?current)))
  =>
  (retract ?h))

(defrule cleanup-old-accuse-used-after-scoring
  ?gs <- (game-state (phase setup) (hand-number ?current))
  ?a <- (accuse-used (hand ?old&:(< ?old ?current)))
  =>
  (retract ?a))

(defrule reject-accuse-when-not-playing
  (declare (salience 100))
  ?da <- (declare-accuse)
  (game-state (phase ?phase&:(neq ?phase playing)))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code wrong-phase) (reason "Accuse can be declared only while playing."))))

(defrule reject-accuse-after-first-trick
  (declare (salience 90))
  ?da <- (declare-accuse)
  (game-state (phase playing) (trick-number ?t&:(> ?t 1)))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code accuse-window) (reason "Accuse is allowed only during the first trick."))))

(defrule reject-accuse-invalid-type
  (declare (salience 80))
  ?da <- (declare-accuse (type ?type&:(and (neq ?type napoli) (neq ?type tris))))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code accuse-type) (reason "Unsupported accuse type."))))

(defrule reject-accuse-invalid-tris-rank
  (declare (salience 75))
  ?da <- (declare-accuse (type tris) (rank ?rank&:(and (neq ?rank ace) (neq ?rank two) (neq ?rank three))))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code accuse-rank) (reason "Tris accuse supports only ace, two or three."))))

(defrule reject-accuse-duplicate
  (declare (salience 70))
  ?da <- (declare-accuse (player ?p) (type ?type) (suit ?s) (rank ?r))
  (game-state (hand-number ?h))
  (accuse-used (hand ?h) (player ?p) (type ?type) (key ?k&:(or (eq ?k ?s) (eq ?k ?r))))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code accuse-duplicate) (reason "This accuse has already been declared by the same player."))))

(defrule declare-napoli-valid
  ?da <- (declare-accuse (player ?p) (type napoli) (suit ?s))
  ?gs <- (game-state (phase playing) (hand-number ?h))
  ?score <- (team-score (team ?team) (points ?pts))
  (test (eq ?team (team-of ?p)))
  (test (and (has-initial-card ?h ?p ?s ace)
             (has-initial-card ?h ?p ?s two)
             (has-initial-card ?h ?p ?s three)))
  =>
  (retract ?da)
  (modify ?score (points (+ ?pts 3)))
  (assert (accuse-used (hand ?h) (player ?p) (type napoli) (key ?s)))
  (assert (accuse-result (hand ?h) (player ?p) (team ?team) (type napoli) (key ?s) (points 3)))
  (assert (action-result (valid yes) (code accuse-ok) (reason "Napoli declared: +3 points."))))

(defrule reject-napoli-missing-cards
  ?da <- (declare-accuse (player ?p) (type napoli) (suit ?s))
  (game-state (phase playing) (hand-number ?h))
  (test (not (and (has-initial-card ?h ?p ?s ace)
                  (has-initial-card ?h ?p ?s two)
                  (has-initial-card ?h ?p ?s three))))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code accuse-invalid) (reason "Player does not hold a valid Napoli in this hand."))))

(defrule declare-tris-valid
  ?da <- (declare-accuse (player ?p) (type tris) (rank ?r))
  ?gs <- (game-state (phase playing) (hand-number ?h))
  ?score <- (team-score (team ?team) (points ?pts))
  (test (eq ?team (team-of ?p)))
  (test (>= (count-initial-rank ?h ?p ?r) 3))
  =>
  (retract ?da)
  (bind ?count (count-initial-rank ?h ?p ?r))
  (bind ?bonus (if (= ?count 4) then 4 else 3))
  (modify ?score (points (+ ?pts ?bonus)))
  (assert (accuse-used (hand ?h) (player ?p) (type tris) (key ?r)))
  (assert (accuse-result (hand ?h) (player ?p) (team ?team) (type tris) (key ?r) (points ?bonus)))
  (assert (action-result (valid yes) (code accuse-ok) (reason "Tris declared."))))

(defrule reject-tris-missing-cards
  ?da <- (declare-accuse (player ?p) (type tris) (rank ?r))
  (game-state (phase playing) (hand-number ?h))
  (test (< (count-initial-rank ?h ?p ?r) 3))
  =>
  (retract ?da)
  (assert (action-result (valid no) (code accuse-invalid) (reason "Player does not hold a valid tris in this hand."))))

(defrule reject-play-after-match-ended
  (declare (salience 120))
  ?pc <- (play-card)
  (match-winner)
  =>
  (retract ?pc)
  (assert (action-result (valid no) (code game-ended) (reason "Match already ended."))))

(defrule reject-accuse-after-match-ended
  (declare (salience 120))
  ?da <- (declare-accuse)
  (match-winner)
  =>
  (retract ?da)
  (assert (action-result (valid no) (code game-ended) (reason "Match already ended."))))
