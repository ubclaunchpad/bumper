package models

import (
	"math"
	"math/rand"
	"sync"
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
	rwMutex       sync.RWMutex
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
		rwMutex:       sync.RWMutex{},
	}
	return &h
}

// Getters
func (h *Hole) getPosition() Position {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	return h.Position
}

func (h *Hole) getLife() float64 {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	return h.Life
}

func (h *Hole) getRadius() float64 {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	return h.Radius
}

func (h *Hole) getGravityRadius() float64 {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	return h.GravityRadius
}

func (h *Hole) getStartingLife() float64 {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	return h.StartingLife
}

// Setters

func (h *Hole) setIsAlive(isAlive bool) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	h.IsAlive = isAlive
}

func (h *Hole) setLife(life float64) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	h.Life = life
}

func (h *Hole) setRadius(radius float64) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	h.Radius = radius
}

func (h *Hole) setGravityRadius(gravityRadius float64) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

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
	if hRadius := h.getRadius(); hRadius < MaxHoleRadius*1.2 {
		h.setRadius(hRadius + 0.02)
		h.setGravityRadius(h.getGravityRadius() + 0.03)
	}
}

// IsDead checks the lifespan of the hole
func (h *Hole) IsDead() bool {
	return h.getLife() < 0
}
