(deffunction next-player (?player)
  (if (eq ?player p1) then p2
   else
    (if (eq ?player p2) then p3
     else
      (if (eq ?player p3) then p4 else p1))))

(deffunction team-of (?player)
  (if (or (eq ?player p1) (eq ?player p3)) then team1 else team2))

(deffunction partner-of (?player)
  (if (eq ?player p1) then p3
   else
    (if (eq ?player p2) then p4
     else
      (if (eq ?player p3) then p1 else p2))))

(deffunction rank-strength (?rank)
  (if (eq ?rank three) then 10
   else
    (if (eq ?rank two) then 9
     else
      (if (eq ?rank ace) then 8
       else
        (if (eq ?rank king) then 7
         else
          (if (eq ?rank knight) then 6
           else
            (if (eq ?rank knave) then 5
             else
              (if (eq ?rank seven) then 4
               else
                (if (eq ?rank six) then 3
                 else
                  (if (eq ?rank five) then 2 else 1)))))))))))

(deffunction card-point-thirds (?rank)
  (if (eq ?rank ace) then 3
   else
    (if (or (eq ?rank three) (eq ?rank two) (eq ?rank king) (eq ?rank knight) (eq ?rank knave))
      then 1
      else 0)))

(deffunction hand-card-count (?hand ?player)
  (length$
    (find-all-facts
      ((?c card-in-hand))
      (and (eq ?c:hand ?hand) (eq ?c:player ?player)))))

(deffunction has-suit-in-hand (?hand ?player ?suit)
  (> (length$
       (find-all-facts
         ((?c card-in-hand))
         (and (eq ?c:hand ?hand) (eq ?c:player ?player) (eq ?c:suit ?suit))))
     0))

(deffunction has-initial-card (?hand ?player ?suit ?rank)
  (> (length$
       (find-all-facts
         ((?c initial-hand-card))
         (and (eq ?c:hand ?hand) (eq ?c:player ?player) (eq ?c:suit ?suit) (eq ?c:rank ?rank))))
     0))

(deffunction count-initial-rank (?hand ?player ?rank)
  (length$
    (find-all-facts
      ((?c initial-hand-card))
      (and (eq ?c:hand ?hand) (eq ?c:player ?player) (eq ?c:rank ?rank)))))

(deffunction player-for-deal-position (?dealer ?position)
  (bind ?offset (+ 1 (mod (- ?position 1) 4)))
  (if (= ?offset 1) then (return (next-player ?dealer)))
  (if (= ?offset 2) then (return (next-player (next-player ?dealer))))
  (if (= ?offset 3) then (return (next-player (next-player (next-player ?dealer)))))
  (return ?dealer))
