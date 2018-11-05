package models

import (
	"sync"
)

// Junk related constants
const (
	JunkFriction          = 0.98
	JunkRadius            = 15
	JunkDebounceTicks     = 10
	JunkGravityDamping    = 0.025
	JunkMass              = 1
	JunkRestitutionFactor = 1
)

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	PhysicsBody   `json:"body"`
	Color         string  `json:"color"`
	LastPlayerHit *Player `json:"-"`
	Debounce      int     `json:"-"`
	jDebounce     int
	rwMutex       sync.RWMutex
}

// JunkMessage contains the data the client needs about a junk
type JunkMessage struct {
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
	Color    string   `json:"color"`
}

// CreateJunk initializes and returns an instance of a Junk
func CreateJunk(position Position) *Junk {
	lock := sync.RWMutex{}
	return &Junk{
		PhysicsBody: CreateBody(position, JunkRadius, JunkMass, JunkRestitutionFactor, &lock),
		Color:       "white",
		Debounce:    0,
		jDebounce:   0,
		rwMutex:     lock,
	}
}

// Getters
func (j *Junk) getDebounce() int {
	j.rwMutex.RLock()
	defer j.rwMutex.RUnlock()

	return j.Debounce
}

func (j *Junk) getJDebounce() int {
	j.rwMutex.RLock()
	defer j.rwMutex.RUnlock()

	return j.jDebounce
}

// Setters
func (j *Junk) setDebounce(debounce int) {
	j.rwMutex.Lock()
	defer j.rwMutex.Unlock()

	j.Debounce = debounce
}

func (j *Junk) setJDebounce(jDebounce int) {
	j.rwMutex.Lock()
	defer j.rwMutex.Unlock()

	j.jDebounce = jDebounce
}

func (j *Junk) setColor(color string) {
	j.rwMutex.Lock()
	defer j.rwMutex.Unlock()

	j.Color = color
}

func (j *Junk) setLastPlayerHit(player *Player) {
	j.rwMutex.Lock()
	defer j.rwMutex.Unlock()

	j.LastPlayerHit = player
}

// UpdatePosition Update Junk's position based on calculations of position/velocity
func (j *Junk) UpdatePosition(height float64, width float64) {
	// Check if the junk should reflect off the walls first
	if j.GetX()+j.GetVelocity().Dx > width-j.GetRadius() || j.GetX()+j.GetVelocity().Dx < j.GetRadius() {
		j.SetDx(-j.GetVelocity().Dx)
	}
	if j.GetY()+j.GetVelocity().Dy > height-j.GetRadius() || j.GetY()+j.GetVelocity().Dy < j.GetRadius() {
		j.SetDy(-j.GetVelocity().Dy)
	}

	j.ApplyFactor(JunkFriction)
	j.ApplyVelocity()

	if jDebounce := j.getDebounce(); jDebounce > 0 {
		j.setDebounce(jDebounce - 1)
	} else {
		j.setDebounce(0)
	}

	if jjDebounce := j.getJDebounce(); jjDebounce > 0 {
		j.setJDebounce(jjDebounce - 1)
	} else {
		j.setJDebounce(0)
	}
}

// HitBy causes a collision event between this junk and given player.
func (j *Junk) HitBy(p *Player) {
	// We don't want this collision till the debounce is done.
	if j.getDebounce() != 0 {
		return
	}

	j.setColor(p.getColor()) //Assign junk to last recently hit player color
	j.setLastPlayerHit(p)

	InelasticCollision(&j.PhysicsBody, &p.PhysicsBody)

	j.setDebounce(JunkDebounceTicks)
}

// HitJunk causes a collision event between this junk and given junk.
func (j *Junk) HitJunk(jh *Junk) {
	// We don't want this collision till the debounce is done.
	if j.getJDebounce() != 0 {
		return
	}
	InelasticCollision(&j.PhysicsBody, &jh.PhysicsBody)

	j.setJDebounce(JunkDebounceTicks)
	jh.setJDebounce(JunkDebounceTicks)
}

// MakeMessage returns a JunkMessage with this junk's data
func (j *Junk) MakeMessage() *JunkMessage {
	return &JunkMessage{
		Position: j.GetPosition(),
		Velocity: j.GetVelocity(),
		Color:    j.Color,
	}
}

// ApplyGravity Applys a vector towards given position
func (j *Junk) ApplyGravity(h *Hole) {
	jVelocity := j.GetVelocity()
	jPosition := j.GetPosition()
	hPosition := h.GetPosition()

	gravityVector := Velocity{0, 0}
	gravityVector.Dx = hPosition.X - jPosition.X
	gravityVector.Dy = hPosition.Y - jPosition.Y
	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	jVelocity.Dx += gravityVector.Dx * inverseMagnitude * h.GetRadius() * JunkGravityDamping
	jVelocity.Dy += gravityVector.Dy * inverseMagnitude * h.GetRadius() * JunkGravityDamping
	j.SetVelocity(jVelocity)
}
