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

func (v *Velocity) magnitude() float64 {
	return math.Sqrt(v.Dx*v.Dx + v.Dy*v.Dy)
}

func (v *Velocity) normalize() *Velocity {
	var mag = v.magnitude()
	if mag > 0 {
		return &Velocity{
			Dx: v.Dx / mag,
			Dy: v.Dy / mag,
		}
	}
	return v
}
