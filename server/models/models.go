package models

import (
	"math"
)

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
