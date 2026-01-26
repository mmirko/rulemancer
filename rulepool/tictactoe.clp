(deftemplate move
  (slot x)
  (slot y)
  (slot player)) ; x | o

(deftemplate last-move
  (slot valid) ; yes | no | none
  (slot reason)) ; description of reason for invalid move

(deftemplate cell
  (slot x)
  (slot y)
  (slot value)) ; x | o

(deftemplate turn
  (slot player)) ; x o o

(deftemplate state
  (slot phase)) ; playing | ended

(deffacts start
  (turn (player x))
  (state (phase playing))
  (last-move (valid none)))

(deffunction switch-player (?current)
  (if (eq ?current x) then o else x))

(defrule invalid-move-not-playing
  ?s <- (state (phase ended))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  =>
  (retract ?m)
  (assert (last-move (valid no) (reason "Game has ended."))))

(defrule invalid-move-cell-occupied
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  ?c <- (cell (x ?x) (y ?y) (value ?v))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-move (valid no) (reason "Cell is already occupied."))))

(defrule invalid-move-wrong-turn
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  ?t <- (turn (player ?current&:(not (eq ?current ?p))))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-move (valid no) (reason "It's not your turn."))))

(defrule valid-move
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  ?t <- (turn (player ?p))
  =>
  (retract ?m)
  (retract ?t)
  (retract ?l)
  (assert (turn (player (switch-player ?p))))
  (assert (last-move (valid yes) (reason "Move accepted.")))
  (assert (cell (x ?x) (y ?y) (value ?p))))

(defrule last-move-show
  ?l <- (last-move (valid ?v) (reason ?r))
  =>
  (printout t "Last move valid: " ?v ", Reason: " ?r crlf))