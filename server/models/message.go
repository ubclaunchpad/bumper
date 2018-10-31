package models

// Message is the generic schema for client/server communication
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ConnectionMessage defines the initial connection message
type ConnectionMessage struct {
	ArenaWidth  float64 `json:"arenaWidth"`
	ArenaHeight float64 `json:"arenaHeight"`
	PlayerID    string  `json:"playerID"`
}

// UpdateMessage defines the schema for a state update message
type UpdateMessage struct {
	Holes   []*HoleMessage   `json:"holes"`
	Junk    []*JunkMessage   `json:"junk"`
	Players []*PlayerMessage `json:"players"`
}

// SpawnHandlerMessage defines a spawn message
type SpawnHandlerMessage struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// KeyHandlerMessage defines a player key press message
type KeyHandlerMessage struct {
	Key       int  `json:"key"`
	IsPressed bool `json:"isPressed"`
}
