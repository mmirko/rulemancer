package rulemancer

type Room struct {
	name        string
	description string
	id          string
}

func (e *Engine) createRoom(name, description string) *Room {
	e.roomsMutex.Lock()
	defer e.roomsMutex.Unlock()
	return &Room{
		name:        name,
		description: description,
		id:          e.generateUniqueID(),
	}
}

func (e *Engine) generateUniqueID() string {
	for {
		newId := RandStringBytes(16)
		if _, exists := e.rooms[newId]; !exists {
			return newId
		}
	}
}
