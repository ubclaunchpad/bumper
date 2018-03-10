package models

// Position x y position
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// Velocity dx dy velocity
type Velocity struct {
	Dx float32 `json:"dx"`
	Dy float32 `json:"dy"`
}

// Player contains data and state about a player's object
type Player struct {
	ID       int      `json:"id"`
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
	Color    string   `json:"color"`
}

// Hole contains the data for a hole's position and size
type Hole struct {
	Position Position `json:"position"`
	Radius   float32  `json:"radius"`
}

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
	Player   Player   `json:"player"`
}
