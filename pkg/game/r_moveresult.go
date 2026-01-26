package game

type MoveResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
