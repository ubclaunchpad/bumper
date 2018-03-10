package models

// Position x y position
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// Velocity dx dy velocity
type Velocity struct {
	Dx float32
	Dy float32
}

// Hole contains the data for a hole's position and size
type Hole struct {
	Position Position
	Radius   float32
}

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	Position Position
	Velocity Velocity
	Player   Player
}
