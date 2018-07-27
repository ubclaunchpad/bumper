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
	IsAlive       bool     `json:"isAlive"`
	Life          float64  `json:"-"`
	StartingLife  float64  `json:"-"`
}

// CreateHole initializes and returns an instance of a Hole
func CreateHole(position Position) *Hole {
	life := math.Floor(rand.Float64()*((MaxHoleLife-MinHoleLife)+1)) + MinHoleLife
	radius := math.Floor(rand.Float64()*((MaxHoleRadius-MinHoleRadius)+1)) + MinHoleRadius
	h := Hole{
		Position:      position,
		Radius:        radius,
		GravityRadius: radius * gravityRadiusFactor,
		Life:          life,
		IsAlive:       false,
		StartingLife:  life,
	}
	return &h
}

// Update reduces this holes life and increases radius if max not reached
func (h *Hole) Update() {
	h.Life--

	if h.Life < h.StartingLife-HoleInfancy {
		h.IsAlive = true
	}
	if h.Radius < MaxHoleRadius*1.2 {
		h.Radius += 0.02
		h.GravityRadius += 0.03
	}
}

// IsDead checks the lifespan of the hole
func (h *Hole) IsDead() bool {
	return h.Life < 0
}

// ApplyGravity applys a vector towards this hole
func (h *Hole) ApplyGravity(o interface{}) {
	gravityVector := Velocity{0, 0}
	j, okay := o.(*Junk)
	if okay {
		gravityVector.Dx = h.Position.X - j.Position.X
		gravityVector.Dy = h.Position.Y - j.Position.Y

		inverseMagnitude := 1.0 / gravityVector.magnitude()
		gravityVector.normalize()

		//Velocity is affected by how close you are, the size of the hole, and a damping factor.
		j.Velocity.Dx += gravityVector.Dx * inverseMagnitude * h.Radius * JunkGravityDamping
		j.Velocity.Dy += gravityVector.Dy * inverseMagnitude * h.Radius * JunkGravityDamping
	}
	p, okay := o.(*Player)
	if okay {
		gravityVector.Dx = h.Position.X - p.Position.X
		gravityVector.Dy = h.Position.Y - p.Position.Y

		inverseMagnitude := 1.0 / gravityVector.magnitude()
		gravityVector.normalize()

		//Velocity is affected by how close you are, the size of the hole, and a damping factor.
		p.Velocity.Dx += gravityVector.Dx * inverseMagnitude * h.Radius * PlayerGravityDamping
		p.Velocity.Dy += gravityVector.Dy * inverseMagnitude * h.Radius * PlayerGravityDamping
	}
}
