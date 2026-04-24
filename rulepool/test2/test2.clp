(deftemplate cast-number
  (multislot number) ; 1-9
  (slot player)) ; x | o

(deftemplate last-cast-number
  (multislot number) ; 1-9
  (slot player)) ; x | o

(deftemplate number
  (multislot value)) ; 1-9

(deffacts start
  (last-cast-number (number none) (player none))
  (number (value 0)))

(defrule valid-move
  ?l <- (last-cast-number (number ?n) (player ?r))
  ?m <- (cast-number (number ?a) (player ?p))
  ?number <- (number (value ?aa))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-cast-number (number ?a) (player ?p)))
  (modify ?number (value ?a)))
