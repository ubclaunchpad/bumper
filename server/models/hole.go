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

// Hole contains the data for a hole object
type Hole struct {
	PhysicsBody   `json:"body"`
	GravityRadius float64 `json:"-"`
	IsAlive       bool    `json:"isAlive"`
	Life          float64 `json:"-"`
	StartingLife  float64 `json:"-"`
	rwMutex       sync.RWMutex
}

// HoleMessage contains the data the client needs about a hole
type HoleMessage struct {
	Position Position `json:"position"`
	IsAlive  bool     `json:"isAlive"`
	Radius   float64  `json:"radius"`
}

// CreateHole initializes and returns an instance of a Hole
func CreateHole(position Position) *Hole {
	life := math.Floor(rand.Float64()*((MaxHoleLife-MinHoleLife)+1)) + MinHoleLife
	radius := math.Floor(rand.Float64()*((MaxHoleRadius-MinHoleRadius)+1)) + MinHoleRadius
	lock := sync.RWMutex{}
	h := Hole{
		PhysicsBody:   CreateBody(position, radius, 0, 0, &lock),
		GravityRadius: radius * gravityRadiusFactor,
		Life:          life,
		IsAlive:       false,
		StartingLife:  life,
		rwMutex:       lock,
	}
	return &h
}

// Getters
func (h *Hole) getLife() float64 {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	return h.Life
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
	if h.GetRadius() < MaxHoleRadius*1.2 {
		h.SetRadius(h.GetRadius() + 0.02)
		h.GravityRadius += 0.03
	}
}

// IsDead checks the lifespan of the hole
func (h *Hole) IsDead() bool {
	return h.getLife() < 0
}

// ApplyGravity modifies given velocity based on given position and damping factor relative to this hole.
func (h *Hole) ApplyGravity(b1 *PhysicsBody, DampingFactor float64) {
	gravityVector := Velocity{0, 0}

	gravityVector.Dx = h.GetX() - b1.GetX()
	gravityVector.Dy = h.GetY() - b1.GetY()

	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	gravityVector.ApplyFactor(inverseMagnitude * h.GetRadius() * DampingFactor)
	b1.ApplyVector(gravityVector)
}

// MakeMessage returns a HoleMessage with this hole's data
func (h *Hole) MakeMessage() *HoleMessage {
	return &HoleMessage{
		Position: h.GetPosition(),
		IsAlive:  h.IsAlive,
		Radius:   h.GetRadius(),
	}
}
