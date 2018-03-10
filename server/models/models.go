package models

import (
	"math"
)

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

// Player contains data and state about a player's object
type Player struct {
	ID       int
	Theta    int
	Position Position
	Velocity Velocity
	Color    string
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

func (v *Velocity) magnitude() float64 {
	return math.Sqrt(float64((v.Dx * v.Dx) + (v.Dy * v.Dy)))
}

func (v *Velocity) normalize() *Velocity {
	var mag = v.magnitude()
	if mag > 0 {
		return &Velocity{
			Dx: v.Dx / float32(mag),
			Dy: v.Dy / float32(mag),
		}
	}

	return v
}
