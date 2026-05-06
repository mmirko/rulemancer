(deftemplate game-config
  (slot game-name)
  (slot description)
  (slot num-players))

(deftemplate assertable
  (slot name)
  (multislot relations))

(deftemplate results
  (slot name)
  (multislot relations))

(deftemplate queryable
  (slot name)
  (multislot relations))

(deftemplate condition
  (slot id)
  (slot type)
  (slot operand)
  (slot value1)
  (slot value2))

(deftemplate acl
  (slot chain)
  (slot position)
  (multislot conditions)
  (slot action))
  