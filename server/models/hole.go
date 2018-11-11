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

// GetID returns the hole's ID
func (h Hole) GetID() string {
	return ""
}

// GetColor returns this hole's color
func (h Hole) GetColor() string {
	return ""
}

// GetPosition returns this hole's position
func (h Hole) GetPosition() Position {
	return h.Position
}

// GetVelocity returns this hole's velocity
func (h Hole) GetVelocity() Velocity {
	return Velocity{}
}

// GetRadius returns this hole's radius
func (h Hole) GetRadius() float64 {
	return h.Radius
}

func (h *Hole) getLife() float64 {
	return h.Life
}

// GetGravityRadius returns this hole's gravitational radius
func (h *Hole) GetGravityRadius() float64 {
	return h.GravityRadius
}

func (h *Hole) getStartingLife() float64 {
	return h.StartingLife
}

func (h *Hole) setIsAlive(isAlive bool) {
	h.IsAlive = isAlive
}

func (h *Hole) setLife(life float64) {
	h.Life = life
}

func (h *Hole) setRadius(radius float64) {
	h.Radius = radius
}

func (h *Hole) setGravityRadius(gravityRadius float64) {
	h.GravityRadius = gravityRadius
}

// Update reduces this holes life and increases radius if max not reached
func (h *Hole) Update() {
	hLife := h.getLife()
	hLife--
	h.setLife(hLife)

	if hLife < h.getStartingLife()-HoleInfancy {
		h.setIsAlive(true)
	}
	if hRadius := h.GetRadius(); hRadius < MaxHoleRadius*1.2 {
		h.setRadius(hRadius + 0.02)
		h.setGravityRadius(h.GetGravityRadius() + 0.03)
	}
}

// IsDead checks the lifespan of the hole
func (h *Hole) IsDead() bool {
	return h.getLife() < 0
}
