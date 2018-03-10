package models

import "math"

// LeftKey is the left keyboard button
const LeftKey = 37

// RightKey is the left keyboard button
const RightKey = 39

// UpKey is the up keyboard button
const UpKey = 38

// DownKey is the down keyboard button
const DownKey = 40

// JunkBounceFactor is how much hitting a junk affects your velocity
const JunkBounceFactor = -0.25

// WallBounceFactor is how much hitting a wall affects your velocity
const WallBounceFactor = -1.5

// PlayerRadius is how big
const PlayerRadius = 25

// PlayerAcceleration is how much a key press affects the player's velocity
const PlayerAcceleration = 0.5

// PlayerFriction is how much damping they experience per tick
const PlayerFriction = 0.97

// MaxVelocity caps how fast a player can get going
const MaxVelocity = 15

// Player contains data and state about a player's object
type Player struct {
	ID       int         `json:"id"`
	Position Position    `json:"position"`
	Velocity Velocity    `json:"velocity"`
	Color    string      `json:"color"`
	Angle    float64     `json:"angle"`
	Controls KeysPressed `json:"controls"`
}

// KeysPressed contains a boolean about each key, true if it's down
type KeysPressed struct {
	Right bool `json:"right"`
	Left  bool `json:"left"`
	Up    bool `json:"up"`
	Down  bool `json:"down"`
}

//Update Player's position based on calculations of position/velocity
func (p *Player) updatePosition(a *Arena) {
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
	if p.Position.X+PlayerRadius > a.Width {
		p.Velocity.Dx *= WallBounceFactor
	} else if p.Position.X-PlayerRadius < 0 {
		p.Velocity.Dx *= WallBounceFactor
	}

	if p.Position.Y+PlayerRadius > a.Height {
		p.Velocity.Dy *= WallBounceFactor
	} else if p.Position.Y-PlayerRadius < 0 {
		p.Velocity.Dy *= WallBounceFactor
	}
}

//Update Player's position based on calculations of hitting junk
func (p *Player) hitJunk() {
	p.Velocity.Dx *= JunkBounceFactor
	p.Velocity.Dy *= JunkBounceFactor
}

//Handle a key press
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

//Handle a key release
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
