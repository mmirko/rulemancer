(deftemplate persona
   (slot nome)
   (slot eta))

(defrule saluta-adulto
   (persona (nome ?n) (eta ?e&:(>= ?e 18)))
   =>
   (printout t "Ciao " ?n ", sei maggiorenne." crlf))
