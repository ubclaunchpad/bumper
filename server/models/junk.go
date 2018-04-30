package models

import (
	"math"

	"github.com/gorilla/websocket"
)

// Junk related constants
const (
	JunkFriction       = 0.99
	MinimumBump        = 0.5
	BumpFactor         = 1.05
	JunkRadius         = 8
	JunkDebounceTicks  = 10
	junkGravityDamping = 0.025
)

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	Position Position        `json:"position"`
	Velocity Velocity        `json:"velocity"` // Don't need to send me to the clients
	Color    string          `json:"color"`
	ID       *websocket.Conn `json:"id"`
	Debounce float32         `json:"debounce"` // Don't need to send me to the clients
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
}

// HitBy Update Junks's velocity based on calculations of being hit by a player
func (j *Junk) HitBy(p *Player, ws *websocket.Conn) {
	if j.Debounce == 0 {
		j.Color = p.Color //Assign junk to last recently hit player color
		j.ID = ws         //Assign junk to last recently hit player id (websocket)

		if p.Velocity.Dx < 0 {
			j.Velocity.Dx = math.Min(p.Velocity.Dx*BumpFactor, -MinimumBump)
		} else {
			j.Velocity.Dx = math.Max(p.Velocity.Dx*BumpFactor, MinimumBump)
		}

		if p.Velocity.Dy < 0 {
			j.Velocity.Dy = math.Min(p.Velocity.Dy*BumpFactor, -MinimumBump)
		} else {
			j.Velocity.Dy = math.Max(p.Velocity.Dy*BumpFactor, MinimumBump)
		}

		p.hitJunk()
		j.Debounce = JunkDebounceTicks
	} else {
		//We don't want this collision till the debounce is down.
	}
}

// ApplyGravity applys a vector towards given position
func (j *Junk) ApplyGravity(h *Hole) {
	gravityVector := Velocity{0, 0}

	if j.Position.X < h.Position.X {
		gravityVector.Dx = h.Position.X - j.Position.X
	} else {
		gravityVector.Dx = h.Position.X - j.Position.X
	}

	if j.Position.Y < h.Position.Y {
		gravityVector.Dy = h.Position.Y - j.Position.Y
	} else {
		gravityVector.Dy = h.Position.Y - j.Position.Y
	}

	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	j.Velocity.Dx += gravityVector.Dx * inverseMagnitude * h.Radius * junkGravityDamping
	j.Velocity.Dy += gravityVector.Dy * inverseMagnitude * h.Radius * junkGravityDamping
}
