package models

// Position x y position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Velocity dx dy velocity
type Velocity struct {
	Dx float64 `json:"dx"`
	Dy float64 `json:"dy"`
}

// Hole contains the data for a hole's position and size
type Hole struct {
	Position Position `json:"position"`
	Radius   float64  `json:"radius"`
}

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
	ID int      `json:"id"` // ID of the player that last recently hit this junk
}
