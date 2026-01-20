/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import "sync"

type Engine struct {
	rooms      map[string]*Room
	roomsMutex sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		rooms:      make(map[string]*Room),
		roomsMutex: sync.RWMutex{},
	}
}

func (e *Engine) SpawnEngine(c *Config, rulePool string) error {
	// Implement the logic to spawn and run the CLIPS engine
	// using the provided configuration and rule pool directory
	return nil
}
