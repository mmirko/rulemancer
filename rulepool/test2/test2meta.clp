(deffacts tictactoe-config
  (game-config
    (game-name test2)
    (description "A simple game where two players alternate placing numbers")
    (num-players 2)))

(deffacts tictactoe-interface
  (assertable
    (name cast-number)
    (relations cast-number))
  (results 
    (name cast-number)
    (relations last-cast-number))
  (queryable
    (name number)
    (relations number)))