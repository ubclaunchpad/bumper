package models

import (
	"math"
	"sync"
)

// Junk related constants
const (
	JunkFriction        = 0.99
	MinimumBump         = 0.5
	BumpFactor          = 1.05
	JunkRadius          = 11
	JunkDebounceTicks   = 15
	JunkVTransferFactor = 0.5
	JunkGravityDamping  = 0.025
)

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	Position      Position `json:"position"`
	Velocity      Velocity `json:"-"`
	Color         string   `json:"color"`
	LastPlayerHit *Player  `json:"-"`
	Debounce      int      `json:"-"`
	jDebounce     int
	rwMutex       sync.RWMutex
}

// CreateJunk initializes and returns an instance of a Junk
func CreateJunk(position Position) *Junk {
	return &Junk{
		Position:  position,
		Velocity:  Velocity{0, 0},
		Color:     "white",
		Debounce:  0,
		jDebounce: 0,
		rwMutex:   sync.RWMutex{},
	}
}

// Getters

func (j *Junk) getPosition() Position {
	j.rwMutex.RLock()
	defer j.rwMutex.RUnlock()

	return j.Position
}

func (j *Junk) getVelocity() Velocity {
	j.rwMutex.RLock()
	defer j.rwMutex.RUnlock()

	return j.Velocity
}

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

func (j *Junk) setPosition(pos Position) {
	j.rwMutex.Lock()
	defer j.rwMutex.Unlock()

	j.Position = pos
}

func (j *Junk) setVelocity(v Velocity) {
	j.rwMutex.Lock()
	defer j.rwMutex.Unlock()

	j.Velocity = v
}

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
	positionVector := j.getPosition()
	velocityVector := j.getVelocity()
	if positionVector.X+velocityVector.Dx > width-JunkRadius || positionVector.X+velocityVector.Dx < JunkRadius {
		velocityVector.Dx = -velocityVector.Dx
	}
	if positionVector.Y+velocityVector.Dy > height-JunkRadius || positionVector.Y+velocityVector.Dy < JunkRadius {
		velocityVector.Dy = -velocityVector.Dy
	}

	velocityVector.Dx *= JunkFriction
	velocityVector.Dy *= JunkFriction

	positionVector.X += velocityVector.Dx
	positionVector.Y += velocityVector.Dy

	j.setPosition(positionVector)
	j.setVelocity(velocityVector)

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

// HitBy Update Junks's velocity based on calculations of being hit by a player
func (j *Junk) HitBy(p *Player) {
	pVelocity := p.getVelocity()
	jVelocity := j.getVelocity()
	// We don't want this collision till the debounce is down.
	if j.getDebounce() != 0 {
		return
	}

	j.setColor(p.getColor()) //Assign junk to last recently hit player color
	j.setLastPlayerHit(p)

	if pVelocity.Dx < 0 {
		jVelocity.Dx = math.Min(pVelocity.Dx*BumpFactor, -MinimumBump)
	} else {
		jVelocity.Dx = math.Max(pVelocity.Dx*BumpFactor, MinimumBump)
	}

	if pVelocity.Dy < 0 {
		jVelocity.Dy = math.Min(pVelocity.Dy*BumpFactor, -MinimumBump)
	} else {
		jVelocity.Dy = math.Max(pVelocity.Dy*BumpFactor, MinimumBump)
	}

	j.setVelocity(jVelocity)
	p.hitJunk()
	j.setDebounce(JunkDebounceTicks)
}

// HitJunk Update Junks's velocity based on calculations of being hit by another Junk
func (j *Junk) HitJunk(jh *Junk) {
	// We don't want this collision till the debounce is down.
	if j.getJDebounce() != 0 {
		return
	}

	jInitialVelocity := j.getVelocity()
	jVelocity := jInitialVelocity
	jhVelocity := jh.getVelocity()
	//Calculate this junks's new velocity
	jVelocity.Dx = (jVelocity.Dx * -JunkVTransferFactor) + (jhVelocity.Dx * JunkVTransferFactor)
	jVelocity.Dy = (jVelocity.Dy * -JunkVTransferFactor) + (jhVelocity.Dy * JunkVTransferFactor)

	//Calculate other junk's new velocity
	jhVelocity.Dx = (jhVelocity.Dx * -JunkVTransferFactor) + (jInitialVelocity.Dx * JunkVTransferFactor)
	jhVelocity.Dy = (jhVelocity.Dy * -JunkVTransferFactor) + (jInitialVelocity.Dy * JunkVTransferFactor)

	j.setJDebounce(JunkDebounceTicks)
	jh.setJDebounce(JunkDebounceTicks)
}

// ApplyGravity applys a vector towards given position
func (j *Junk) ApplyGravity(h *Hole) {
	jVelocity := j.getVelocity()
	jPosition := j.getPosition()

	gravityVector := Velocity{0, 0}
	gravityVector.Dx = h.Position.X - jPosition.X
	gravityVector.Dy = h.Position.Y - jPosition.Y
	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	jVelocity.Dx += gravityVector.Dx * inverseMagnitude * h.Radius * JunkGravityDamping
	jVelocity.Dy += gravityVector.Dy * inverseMagnitude * h.Radius * JunkGravityDamping
	j.setVelocity(jVelocity)
}
