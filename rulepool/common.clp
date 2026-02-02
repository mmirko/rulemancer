(deftemplate game-config
  (slot game-name)
  (slot description))

(deftemplate assertable
  (slot name)
  (multislot relations))

(deftemplate results
  (slot name)
  (multislot relations))

(deftemplate queryable
  (slot name)
  (multislot relations))