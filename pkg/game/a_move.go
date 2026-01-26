package game

type Move struct {
	Player string `json:"player"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func move2ClipsAssert(move Move) string {
	return `(move (player ` + move.Player + `) (x ` +
		string(rune('0'+move.X)) + `) (y ` +
		string(rune('0'+move.Y)) + `))`
}

func clipsAssert2MoveResult(clipsOutput string) (MoveResult, error) {
	var result MoveResult
	return result, nil
}
