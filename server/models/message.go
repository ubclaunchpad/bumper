package models

// Message is the generic schema for client/server communication
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ConnectionMessage defines the schema for the initial connection message
type ConnectionMessage struct {
	ArenaWidth  float64 `json:"arenaWidth"`
	ArenaHeight float64 `json:"arenaHeight"`
	PlayerID    string  `json:"playerID"`
}

// SpawnHandlerMessage defines the schema for a spawn message
type SpawnHandlerMessage struct {
	Name string `json:"name"`
}

// KeyHandlerMessage defines the schema for a player's key presses
type KeyHandlerMessage struct {
	Key       int  `json:"key"`
	IsPressed bool `json:"isPressed"`
}
