package models

import (
	"math"
	"math/rand"
)

// Hole related constants
const (
	MinHoleRadius = 15
	MaxHoleRadius = 45
	HzToSeconds   = 60
	MinHoleLife   = 25 * HzToSeconds
	MaxHoleLife   = 75 * HzToSeconds
	HoleInfancy   = 2 * HzToSeconds
)

// Hole contains the data for a hole's position and size
type Hole struct {
	Position      Position `json:"position"`
	Radius        float64  `json:"radius"`
	GravityRadius float64  `json:"gravrad"`
	Alive         bool     `json:"islive"`
	Life          float64  `json:"life"`
	StaringLife   float64  `json:"born"`
}

// CreateHole initializes and returns an instance of a Hole
func CreateHole(position Position) Hole {
	life := math.Floor(rand.Float64()*((MaxHoleLife-MinHoleLife)+1)) + MinHoleLife
	return Hole{
		Position:      position,
		Radius:        math.Floor(rand.Float64()*((MaxHoleRadius-MinHoleRadius)+1)) + MinHoleRadius,
		GravityRadius: MaxHoleRadius * 4,
		Life:          life,
		Alive:         false,
		StaringLife:   life,
	}
}

// Update reduces this holes life or star a new one for it
func (h *Hole) Update() {
	h.Life--

	if h.Life < h.StaringLife-HoleInfancy {
		h.Alive = true
	}
	if h.Radius < MaxHoleRadius*1.2 {
		h.Radius += 0.02
	}
}

// StartNewLife sets this hole to a new position and lifespan
func (h *Hole) StartNewLife(position Position) {
	h.Alive = false
	h.Life = math.Floor(rand.Float64()*((MaxHoleLife-MinHoleLife)+1)) + MinHoleLife
	h.StaringLife = h.Life
	h.Radius = math.Floor(rand.Float64()*((MaxHoleRadius-MinHoleRadius)+1)) + MinHoleRadius
	h.Position = position
}
