package models

import (
	"math"
	"math/rand"
)

// Hole related constants
const (
	MinHoleRadius       = 15
	MaxHoleRadius       = 45
	gravityRadiusFactor = 5
	HzToSeconds         = 60
	MinHoleLife         = 25 * HzToSeconds
	MaxHoleLife         = 75 * HzToSeconds
	HoleInfancy         = 2 * HzToSeconds
)

// Hole contains the data for a hole's position and size
type Hole struct {
	Position      Position `json:"position"`
	Radius        float64  `json:"radius"`
	GravityRadius float64  `json:"-"`
	IsCollidable  bool     `json:"isCollidable"`
	IsAlive       bool     `json:"isAlive"`
	Life          float64  `json:"-"`
	StartingLife  float64  `json:"-"`
}

// CreateHole initializes and returns an instance of a Hole
func CreateHole(position Position) Hole {
	life := math.Floor(rand.Float64()*((MaxHoleLife-MinHoleLife)+1)) + MinHoleLife
	radius := math.Floor(rand.Float64()*((MaxHoleRadius-MinHoleRadius)+1)) + MinHoleRadius
	return Hole{
		Position:      position,
		Radius:        radius,
		GravityRadius: radius * gravityRadiusFactor,
		Life:          life,
		IsCollidable:  false,
		IsAlive:       true,
		StartingLife:  life,
	}
}

// Update reduces this holes life and increases radius if max not reached
func (h *Hole) Update() {
	h.Life--

	if h.Life < h.StartingLife-HoleInfancy {
		h.IsCollidable = true
	}
	if h.Life < 0 {
		h.IsAlive = false
	}
	if h.Radius < MaxHoleRadius*1.2 {
		h.Radius += 0.02
		h.GravityRadius += 0.03
	}
}
