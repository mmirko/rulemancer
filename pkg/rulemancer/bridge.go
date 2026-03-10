package rulemancer

import (
	"errors"
	"log"
	"os"
	"sync"
)

type Bridge struct {
	name           string
	id             string
	rulesLocation  string
	runningBrRooms map[string]*BrRoom
	BrRoomsMutex   sync.RWMutex
}

func (g *Bridge) Info() map[string]any {
	return map[string]any{
		"id":             g.id,
		"name":           g.name,
		"rulesLocation":  g.rulesLocation,
		"runningBrRooms": g.runningBrRooms,
	}
}

func (e *Engine) loadBridges() {
	// Load bridges from the configured bridges map
	for name, rulesLocation := range e.Bridges {
		if err := e.newBridge(name, rulesLocation); err != nil {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/loadBridges]")+" ", 0)
			l.Printf("error loading bridge %s from %s: %v", name, rulesLocation, err)
		} else {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/loadBridges]")+" ", 0)
				l.Printf("successfully loaded bridge %s from %s", name, rulesLocation)
			}
		}
	}
}

func (e *Engine) newBridge(name, rulesLocation string) error {

	e.bridgesMutex.Lock()
	defer e.bridgesMutex.Unlock()
	bridge := &Bridge{
		rulesLocation:  rulesLocation,
		name:           name,
		id:             e.generateBridgeUniqueID(),
		runningBrRooms: make(map[string]*BrRoom),
		BrRoomsMutex:   sync.RWMutex{},
	}
	e.numBridges++
	e.bridges[bridge.id] = bridge

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/newBridge]")+" ", 0)
		l.Printf("Loaded bridge %s with ID %s", bridge.rulesLocation, bridge.id)
		l.Print(bridge)
	}
	return nil
}

func (e *Engine) generateBridgeUniqueID() string {
	for {
		newId := randStringBytes(16)
		if _, exists := e.bridges[newId]; !exists {
			return newId
		}
	}
}

// Search for a bridge by its ID
func (e *Engine) searchBridge(id string) (*Bridge, error) {
	e.bridgesMutex.RLock()
	defer e.bridgesMutex.RUnlock()

	// Search by ID
	if bridge, exists := e.bridges[id]; exists {
		return bridge, nil
	}

	// Search by name if no bridge found by ID
	for _, bridge := range e.bridges {
		if bridge.name == id {
			return bridge, nil
		}
	}

	return nil, errors.New("bridge not found")
}

func (e *Engine) listBridges() []string {
	e.bridgesMutex.RLock()
	defer e.bridgesMutex.RUnlock()
	bridges := make([]string, 0, len(e.bridges))
	for _, bridge := range e.bridges {
		bridges = append(bridges, bridge.id)
	}
	return bridges
}
