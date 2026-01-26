package rulemancer

func (e *Engine) RebuildEngine() error {
	// The rebuild engine reads the rule pool directory and the assertables,results and querables
	// from the configuration and uses re2c to write the template files under the pkg/game/ directory.

	// This way, the engine can be rebuilt with the new possible queries and assertables, adapting to
	// different game logics. This way the game logic is decoupled from the engine itself and completely
	// contained in the CLIPS rules pool.

	// TODO
	return nil
}
