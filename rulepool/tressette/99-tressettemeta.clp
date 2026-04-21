(deffacts tressette-config
  (game-config
    (game-name tressette)
    (description "4-player partnership Tressette with full trick-taking, scoring, and accuse support.")
    (num-players 4)))

(deffacts tressette-interface
  (assertable
    (name play-card)
    (relations play-card))
  (assertable
    (name declare-accuse)
    (relations declare-accuse))

  (results
    (name play-card)
    (relations action-result trick-result hand-result match-winner))
  (results
    (name declare-accuse)
    (relations action-result accuse-result team-score))

  (queryable
    (name state)
    (relations game-state trick-state team-score action-result trick-result hand-result accuse-result match-winner))
  (queryable
    (name hand)
    (relations card-in-hand))
  (queryable
    (name table)
    (relations trick-card))
  (queryable
    (name score)
    (relations team-score hand-tally))
  (queryable
    (name winner)
    (relations match-winner)))
