/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import "errors"

type Room struct {
	name          string
	description   string
	id            string
	clipsInstance *ClipsInstance
}

func (e *Engine) newRoom(name, description string) (*Room, error) {
	// Ensure unique ID generation and locking on the rooms map
	var cli *ClipsInstance
	if !e.ClipsLessMode {
		cli = e.NewClipsInstance()
		if err := cli.InitClips(); err != nil {
			return nil, err
		}
		if err := cli.LoadKnowledgeBase(); err != nil {
			cli.Dispose()
			return nil, err
		}
	}
	e.roomsMutex.Lock()
	defer e.roomsMutex.Unlock()
	room := &Room{
		name:          name,
		description:   description,
		id:            e.generateUniqueID(),
		clipsInstance: cli,
	}
	e.rooms[room.id] = room
	return room, nil
}

func (e *Engine) generateUniqueID() string {
	for {
		newId := RandStringBytes(16)
		if _, exists := e.rooms[newId]; !exists {
			return newId
		}
	}
}

func (e *Engine) searchRoom(id string) (*Room, error) {
	e.roomsMutex.RLock()
	defer e.roomsMutex.RUnlock()
	if room, exists := e.rooms[id]; exists {
		return room, nil
	}
	return nil, errors.New("room not found")
}

func (e *Engine) removeRoom(id string) (*Room, error) {
	e.roomsMutex.Lock()
	defer e.roomsMutex.Unlock()
	if room, exists := e.rooms[id]; exists {
		if !e.ClipsLessMode {
			room.clipsInstance.Dispose()
		}
		delete(e.rooms, id)
		return room, nil
	}
	return nil, errors.New("room not found")
}

func (e *Engine) listRooms() []string {
	e.roomsMutex.RLock()
	defer e.roomsMutex.RUnlock()
	rooms := make([]string, 0, len(e.rooms))
	for _, room := range e.rooms {
		rooms = append(rooms, room.id)
	}
	return rooms
}
