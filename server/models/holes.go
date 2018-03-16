package models

import (
	"math"
	"math/rand"
)

// Hole related constants
const (
	MinHoleRadius = 15
	MaxHoleRadius = 45
	MinHoleLife   = 25
	MaxHoleLife   = 75
)

// Hole contains the data for a hole's position and size
type Hole struct {
	Position Position `json:"position"`
	Radius   float64  `json:"radius"`
	Life     float64  `json:"life"`
}

// CreateHole initializes and returns an instance of a Hole
func CreateHole(position Position) Hole {
	return Hole{
		Position: position,
		Radius:   math.Floor(rand.Float64()*((MaxHoleRadius-MinHoleRadius)+1)) + MinHoleRadius,
		Life:     math.Floor(rand.Float64()*((MaxHoleLife-MinHoleLife)+1)) + MinHoleLife,
	}
}

// Set this hole to a new position and lifespan
func (h *Hole) startNewLife() {

}
