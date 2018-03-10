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
	Dx float32 `json:"dx"`
	Dy float32 `json:"dy"`
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
