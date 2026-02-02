package rulemancer

func (e *Engine) BuildEngineExtras() error {
	// The rebuild engine reads the rules games directories and the assertables,results and querables from there.
	// Then uses re2c to write the various artifacts needed to interact with the engine.

	// Load all games from the rules directory

	e.loadGames()

	// TODO
	return nil
}
