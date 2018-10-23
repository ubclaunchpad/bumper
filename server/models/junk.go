package models

import (
	"math"
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
}

// CreateJunk initializes and returns an instance of a Junk
func CreateJunk(position Position) *Junk {
	return &Junk{
		Position:  position,
		Velocity:  Velocity{0, 0},
		Color:     "white",
		Debounce:  0,
		jDebounce: 0,
	}
}

// UpdatePosition Update Junk's position based on calculations of position/velocity
func (j *Junk) UpdatePosition(height float64, width float64) {
	if j.Position.X+j.Velocity.Dx > width-JunkRadius || j.Position.X+j.Velocity.Dx < JunkRadius {
		j.Velocity.Dx = -j.Velocity.Dx
	}
	if j.Position.Y+j.Velocity.Dy > height-JunkRadius || j.Position.Y+j.Velocity.Dy < JunkRadius {
		j.Velocity.Dy = -j.Velocity.Dy
	}

	j.Velocity.Dx *= JunkFriction
	j.Velocity.Dy *= JunkFriction

	j.Position.X += j.Velocity.Dx
	j.Position.Y += j.Velocity.Dy

	if j.Debounce > 0 {
		j.Debounce--
	} else {
		j.Debounce = 0
	}

	if j.jDebounce > 0 {
		j.jDebounce--
	} else {
		j.jDebounce = 0
	}
}

// HitBy Update Junks's velocity based on calculations of being hit by a player
func (j *Junk) HitBy(p *Player) {
	pVelocity := p.getVelocity()
	// We don't want this collision till the debounce is down.
	if j.Debounce != 0 {
		return
	}

	j.Color = p.getColor() //Assign junk to last recently hit player color
	j.LastPlayerHit = p

	if pVelocity.Dx < 0 {
		j.Velocity.Dx = math.Min(pVelocity.Dx*BumpFactor, -MinimumBump)
	} else {
		j.Velocity.Dx = math.Max(pVelocity.Dx*BumpFactor, MinimumBump)
	}

	if pVelocity.Dy < 0 {
		j.Velocity.Dy = math.Min(pVelocity.Dy*BumpFactor, -MinimumBump)
	} else {
		j.Velocity.Dy = math.Max(pVelocity.Dy*BumpFactor, MinimumBump)
	}

	p.hitJunk()
	j.Debounce = JunkDebounceTicks
}

// HitJunk Update Junks's velocity based on calculations of being hit by another Junk
func (j *Junk) HitJunk(jh *Junk) {
	// We don't want this collision till the debounce is down.
	if j.jDebounce != 0 {
		return
	}

	initalVelocity := j.Velocity

	//Calculate this junks's new velocity
	j.Velocity.Dx = (j.Velocity.Dx * -JunkVTransferFactor) + (jh.Velocity.Dx * JunkVTransferFactor)
	j.Velocity.Dy = (j.Velocity.Dy * -JunkVTransferFactor) + (jh.Velocity.Dy * JunkVTransferFactor)

	//Calculate other junk's new velocity
	jh.Velocity.Dx = (jh.Velocity.Dx * -JunkVTransferFactor) + (initalVelocity.Dx * JunkVTransferFactor)
	jh.Velocity.Dy = (jh.Velocity.Dy * -JunkVTransferFactor) + (initalVelocity.Dy * JunkVTransferFactor)

	j.jDebounce = JunkDebounceTicks
	jh.jDebounce = JunkDebounceTicks
}

// ApplyGravity applys a vector towards given position
func (j *Junk) ApplyGravity(h *Hole) {
	gravityVector := Velocity{0, 0}
	gravityVector.Dx = h.Position.X - j.Position.X
	gravityVector.Dy = h.Position.Y - j.Position.Y

	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	j.Velocity.Dx += gravityVector.Dx * inverseMagnitude * h.Radius * JunkGravityDamping
	j.Velocity.Dy += gravityVector.Dy * inverseMagnitude * h.Radius * JunkGravityDamping
}
