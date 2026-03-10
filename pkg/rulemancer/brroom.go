/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"errors"
	"time"
)

type BrRoom struct {
	id            string
	bridge        *Bridge
	clipsInstance *ClipsInstance
	lastActive    int64
}

func (e *Engine) newBrRoom(name, bridgeRef string) (*BrRoom, error) {

	bridge, err := e.searchBridge(bridgeRef)
	if err != nil {
		return nil, err
	}

	rulesLocation := bridge.rulesLocation

	// Ensure unique ID generation and locking on the rooms map
	var cli *ClipsInstance
	if !e.ClipsLessMode {
		cli = e.NewClipsInstance()
		if err := cli.InitClips(); err != nil {
			return nil, err
		}
		if err := cli.loadGame(rulesLocation); err != nil {
			cli.Dispose()
			return nil, err
		}
	}
	e.brRoomsMutex.Lock()
	defer e.brRoomsMutex.Unlock()
	brRoom := &BrRoom{
		id:            name,
		bridge:        bridge,
		clipsInstance: cli,
		lastActive:    time.Now().Unix(),
	}
	e.numRooms++
	e.brRooms[brRoom.id] = brRoom

	bridge.BrRoomsMutex.Lock()
	defer bridge.BrRoomsMutex.Unlock()
	bridge.runningBrooms[brRoom.id] = brRoom

	return brRoom, nil
}

func (e *Engine) searchBrRoom(id string) (*BrRoom, error) {
	e.brRoomsMutex.RLock()
	defer e.brRoomsMutex.RUnlock()
	if room, exists := e.brRooms[id]; exists {
		return room, nil
	}
	return nil, errors.New("bridge room not found")
}

func (e *Engine) removeBrRoom(id string) (*BrRoom, error) {
	e.brRoomsMutex.Lock()
	defer e.brRoomsMutex.Unlock()
	if room, exists := e.brRooms[id]; exists {
		if !e.ClipsLessMode {
			room.clipsInstance.Dispose()
		}
		delete(e.brRooms, id)
		e.numBrRooms--
		return room, nil
	}
	return nil, errors.New("bridge room not found")
}

func (e *Engine) listBrRooms() []string {
	e.brRoomsMutex.RLock()
	defer e.brRoomsMutex.RUnlock()
	rooms := make([]string, 0, len(e.brRooms))
	for _, room := range e.brRooms {
		rooms = append(rooms, room.id)
	}
	return rooms
}
