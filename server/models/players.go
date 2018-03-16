package models

import (
	"math"
)

// Player related constants
const (
	LeftKey            = 37
	RightKey           = 39
	UpKey              = 38
	DownKey            = 40
	JunkBounceFactor   = -0.25
	WallBounceFactor   = -1.5
	PlayerRadius       = 25
	PlayerAcceleration = 0.5
	PlayerFriction     = 0.97
	MaxVelocity        = 15
	PointsPerJunk      = 100
)

// Player contains data and state about a player's object
type Player struct {
	ID       int         `json:"id"`
	Position Position    `json:"position"`
	Velocity Velocity    `json:"velocity"`
	Color    string      `json:"color"`
	Angle    float64     `json:"angle"`
	Controls KeysPressed `json:"controls"`
	Points   int         `json:"points"`
}

// KeysPressed contains a boolean about each key, true if it's down
type KeysPressed struct {
	Right bool `json:"right"`
	Left  bool `json:"left"`
	Up    bool `json:"up"`
	Down  bool `json:"down"`
}

//Update Player's position based on calculations of position/velocity
func (p *Player) updatePosition(height float64, width float64) {
	controlsVector := Velocity{0, 0}

	if p.Controls.Left {
		p.Angle = math.Mod(p.Angle+0.1, 360)
	}

	if p.Controls.Right {
		p.Angle = math.Mod(p.Angle-0.1, 360)
	}

	if p.Controls.Up {
		controlsVector.Dy = (0.5 * PlayerRadius * math.Cos(p.Angle))
		controlsVector.Dx = (0.5 * PlayerRadius * math.Sin(p.Angle))
	}

	controlsVector.normalize()
	controlsVector.Dx *= PlayerAcceleration
	controlsVector.Dy *= PlayerAcceleration

	p.Velocity.Dx *= PlayerFriction
	p.Velocity.Dy *= PlayerFriction

	p.Velocity.Dx += controlsVector.Dx
	p.Velocity.Dy += controlsVector.Dy

	// Ensure it never gets going too fast
	if p.Velocity.magnitude() > MaxVelocity {
		p.Velocity.normalize()
		p.Velocity.Dx *= MaxVelocity
		p.Velocity.Dy *= MaxVelocity
	}

	// Apply player's velocity vector
	p.Position.X += p.Velocity.Dx
	p.Position.Y += p.Velocity.Dy

	// Check wall collisions
	if p.Position.X+PlayerRadius > width {
		p.Velocity.Dx *= WallBounceFactor
	} else if p.Position.X-PlayerRadius < 0 {
		p.Velocity.Dx *= WallBounceFactor
	}

	if p.Position.Y+PlayerRadius > height {
		p.Velocity.Dy *= WallBounceFactor
	} else if p.Position.Y-PlayerRadius < 0 {
		p.Velocity.Dy *= WallBounceFactor
	}
}

func (p *Player) hitJunk() {
	p.Velocity.Dx *= JunkBounceFactor
	p.Velocity.Dy *= JunkBounceFactor
}

// HitPlayer calculates collision, update Player's position based on calculation of hitting another player
func (p *Player) HitPlayer() {
	p.Velocity.Dx *= JunkBounceFactor
	p.Velocity.Dy *= JunkBounceFactor
}

func (p *Player) keyDownHandler(key int) {
	if key == RightKey {
		p.Controls.Right = true
	} else if key == LeftKey {
		p.Controls.Left = true
	} else if key == UpKey {
		p.Controls.Up = true
	} else if key == DownKey {
		p.Controls.Down = true
	}
}

func (p *Player) keyUpHandler(key int) {
	if key == RightKey {
		p.Controls.Right = false
	} else if key == LeftKey {
		p.Controls.Left = false
	} else if key == UpKey {
		p.Controls.Up = false
	} else if key == DownKey {
		p.Controls.Down = false
	}
}
