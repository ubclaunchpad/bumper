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

// GetID returns the ID of this junk
func (j Junk) GetID() string {
	return ""
}

// GetColor returns the color of this junk
func (j Junk) GetColor() string {
	return j.Color
}

// GetPosition returns the position of this jun
func (j Junk) GetPosition() Position {
	return j.Position
}

// GetVelocity returns the velocity of this junk
func (j Junk) GetVelocity() Velocity {
	return j.Velocity
}

// GetRadius returns the radius of this junk
func (j Junk) GetRadius() float64 {
	return JunkRadius
}

func (j *Junk) getDebounce() int {
	return j.Debounce
}

func (j *Junk) getJDebounce() int {
	return j.jDebounce
}

func (j *Junk) setPosition(pos Position) {
	j.Position = pos
}

func (j *Junk) setVelocity(v Velocity) {
	j.Velocity = v
}

func (j *Junk) setDebounce(debounce int) {
	j.Debounce = debounce
}

func (j *Junk) setJDebounce(jDebounce int) {
	j.jDebounce = jDebounce
}

func (j *Junk) setColor(color string) {
	j.Color = color
}

func (j *Junk) setLastPlayerHit(player *Player) {
	j.LastPlayerHit = player
}

// UpdatePosition Update Junk's position based on calculations of position/velocity
func (j *Junk) UpdatePosition(height float64, width float64) {
	positionVector := j.GetPosition()
	velocityVector := j.GetVelocity()
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
	pVelocity := p.GetVelocity()
	jVelocity := j.GetVelocity()
	// We don't want this collision till the debounce is down.
	if j.getDebounce() != 0 {
		return
	}

	j.setColor(p.GetColor()) //Assign junk to last recently hit player color
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

	jInitialVelocity := j.GetVelocity()
	jVelocity := jInitialVelocity
	jhVelocity := jh.GetVelocity()
	//Calculate this junks's new velocity
	jVelocity.Dx = (jVelocity.Dx * -JunkVTransferFactor) + (jhVelocity.Dx * JunkVTransferFactor)
	jVelocity.Dy = (jVelocity.Dy * -JunkVTransferFactor) + (jhVelocity.Dy * JunkVTransferFactor)

	//Calculate other junk's new velocity
	jhVelocity.Dx = (jhVelocity.Dx * -JunkVTransferFactor) + (jInitialVelocity.Dx * JunkVTransferFactor)
	jhVelocity.Dy = (jhVelocity.Dy * -JunkVTransferFactor) + (jInitialVelocity.Dy * JunkVTransferFactor)

	j.setVelocity(jVelocity)
	jh.setVelocity(jhVelocity)
	j.setJDebounce(JunkDebounceTicks)
	jh.setJDebounce(JunkDebounceTicks)
}

// ApplyGravity applys a vector towards given position
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
	j.setVelocity(jVelocity)
}
