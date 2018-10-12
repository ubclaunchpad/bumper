package models

import (
	"log"
	"math"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

// Player related constants
const (
	LeftKey                = 37
	RightKey               = 39
	UpKey                  = 38
	DownKey                = 40
	JunkBounceFactor       = -0.25
	VelocityTransferFactor = 0.75
	WallBounceFactor       = -1.5
	PlayerRadius           = 25
	PlayerAcceleration     = 0.5
	PlayerFriction         = 0.97
	MaxVelocity            = 15
	PointsPerJunk          = 100
	PointsPerPlayer        = 500
	gravityDamping         = 0.075
	PlayerDebounceTicks    = 15
	PointsDebounceTicks    = 100
)

// KeysPressed contains a boolean about each key, true if it's down
type KeysPressed struct {
	Right bool `json:"right"`
	Left  bool `json:"left"`
	Up    bool `json:"up"`
	Down  bool `json:"down"`
}

// Player contains data and state about a player's object
type Player struct {
	Name           string      `json:"name"`
	ID             string      `json:"id"`
	Country        string      `json:"country"`
	Position       Position    `json:"position"`
	Velocity       Velocity    `json:"-"`
	Color          string      `json:"color"`
	Angle          float64     `json:"angle"`
	Controls       KeysPressed `json:"-"`
	Points         int         `json:"points"`
	LastPlayerHit  *Player     `json:"-"`
	pointsDebounce int
	pDebounce      int
	rwMutex        sync.RWMutex
	ws             *websocket.Conn
}

// CreatePlayer constructs an instance of player with
// given position, color, and WebSocket connection
func CreatePlayer(id string, name string, pos Position, color string, ws *websocket.Conn) *Player {
	return &Player{
		Name:           name,
		ID:             id,
		Position:       pos,
		Velocity:       Velocity{},
		Color:          color,
		Angle:          math.Pi,
		Controls:       KeysPressed{},
		pDebounce:      0,
		pointsDebounce: 0,
		rwMutex:        sync.RWMutex{},
		ws:             ws,
	}
}

// GetName returns name of player
func (p *Player) GetName() string {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.Name
}

func (p *Player) getVelocity() Velocity {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.Velocity
}

func (p *Player) getControls() KeysPressed {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.Controls
}

func (p *Player) getPosition() Position {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.Position
}

func (p *Player) getPDebounce() int {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.pDebounce
}

func (p *Player) getPointsDebounce() int {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.pointsDebounce
}

func (p *Player) getLastPlayerHit() *Player {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.LastPlayerHit
}

func (p *Player) setVelocity(v Velocity) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.Velocity = v
}

func (p *Player) setControls(k KeysPressed) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.Controls = k
}

func (p *Player) setAngle(a float64) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.Angle = a
}

func (p *Player) setPosition(pos Position) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.Position = pos
}

func (p *Player) setPDebounce(pDebounce int) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.pDebounce = pDebounce
}

func (p *Player) setPointsDebounce(pointsDebounce int) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.pointsDebounce = pointsDebounce
}

func (p *Player) setLastPlayerHit(playerHit *Player) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.LastPlayerHit = playerHit
}

// GenUniqueID generates unique id...
func GenUniqueID() string {
	id := xid.New()
	return id.String()
}

// AddPoints adds numPoints to player p
func (p *Player) AddPoints(numPoints int) {
	p.Points = p.Points + numPoints
}

// SendJSON sends JSON data through the player's websocket connection
func (p *Player) SendJSON(m *Message) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.ws.WriteJSON(m)
}

// Close ends the WebSocket connection with the player
func (p *Player) Close() {
	err := p.ws.Close()
	if err != nil {
		log.Printf("Failed to close connection:\n%v", err)
	}
}

// UpdatePosition based on calculations of position/velocity
func (p *Player) UpdatePosition(height float64, width float64) {

	controlsVector := Velocity{0, 0}

	if p.getControls().Left {
		p.setAngle(math.Mod(p.Angle+0.1, 360))
	}

	if p.getControls().Right {
		p.setAngle(math.Mod(p.Angle-0.1, 360))
	}

	if p.getControls().Up {
		controlsVector.Dy = (0.5 * PlayerRadius * math.Cos(p.Angle))
		controlsVector.Dx = (0.5 * PlayerRadius * math.Sin(p.Angle))
	}

	controlsVector.normalize()
	controlsVector.Dx *= PlayerAcceleration
	controlsVector.Dy *= PlayerAcceleration

	positionVector := p.getPosition()
	velocityVector := p.getVelocity()
	velocityVector.Dx = (velocityVector.Dx * PlayerFriction) + controlsVector.Dx
	velocityVector.Dy = (velocityVector.Dy * PlayerFriction) + controlsVector.Dy

	// Ensure it never gets going too fast
	if velocityVector.magnitude() > MaxVelocity {
		velocityVector.normalize()
		velocityVector.Dx *= MaxVelocity
		velocityVector.Dy *= MaxVelocity
	}

	// Apply player's velocity vector
	p.setVelocity(velocityVector)

	// Calculate next position
	positionVector.X = positionVector.X + velocityVector.Dx
	positionVector.Y = positionVector.Y + velocityVector.Dy

	// Set position
	p.setPosition(positionVector)

	p.checkWalls(height, width)

	if pDebounce := p.getPDebounce(); pDebounce > 0 {
		p.setPDebounce(pDebounce - 1)
	} else {
		p.setPDebounce(0)
	}

	if pointsDebounce := p.getPointsDebounce(); pointsDebounce > 0 {
		p.setPointsDebounce(pointsDebounce - 1)
	} else {
		p.setLastPlayerHit(nil)
		p.setPointsDebounce(0)
	}
}

func (p *Player) hitJunk() {
	p.Velocity.Dx *= JunkBounceFactor
	p.Velocity.Dy *= JunkBounceFactor
}

// HitPlayer calculates collision, update Player's velocity based on calculation of hitting another player
func (p *Player) HitPlayer(ph *Player) {
	if p.pDebounce != 0 {
		return
	}

	initalVelocity := p.Velocity

	//Calculate player's new velocity
	p.Velocity.Dx = (p.Velocity.Dx * -VelocityTransferFactor) + (ph.Velocity.Dx * VelocityTransferFactor)
	p.Velocity.Dy = (p.Velocity.Dy * -VelocityTransferFactor) + (ph.Velocity.Dy * VelocityTransferFactor)

	//Calculate the player you hits new velocity
	ph.Velocity.Dx = (ph.Velocity.Dx * -VelocityTransferFactor) + (initalVelocity.Dx * VelocityTransferFactor)
	ph.Velocity.Dy = (ph.Velocity.Dy * -VelocityTransferFactor) + (initalVelocity.Dy * VelocityTransferFactor)

	ph.LastPlayerHit = p
	p.LastPlayerHit = ph
	p.pointsDebounce = PointsDebounceTicks
	ph.pointsDebounce = PointsDebounceTicks
	p.pDebounce = PlayerDebounceTicks
}

// ApplyGravity applys a vector towards given position
func (p *Player) ApplyGravity(h *Hole) {
	gravityVector := Velocity{0, 0}
	gravityVector.Dx = h.Position.X - p.Position.X
	gravityVector.Dy = h.Position.Y - p.Position.Y

	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	p.Velocity.Dx += gravityVector.Dx * inverseMagnitude * h.Radius * gravityDamping
	p.Velocity.Dy += gravityVector.Dy * inverseMagnitude * h.Radius * gravityDamping
}

// checkWalls check if the player is attempting to exit the walls, reverse they're direction
func (p *Player) checkWalls(height float64, width float64) {
	if p.Position.X+PlayerRadius > width {
		p.Position.X = width - PlayerRadius - 1
		p.Velocity.Dx *= WallBounceFactor
	} else if p.Position.X-PlayerRadius < 0 {
		p.Velocity.Dx *= WallBounceFactor
		p.Position.X = PlayerRadius + 1
	}

	if p.Position.Y+PlayerRadius > height {
		p.Velocity.Dy *= WallBounceFactor
		p.Position.Y = height - PlayerRadius - 1
	} else if p.Position.Y-PlayerRadius < 0 {
		p.Velocity.Dy *= WallBounceFactor
		p.Position.Y = PlayerRadius + 1
	}
}

// KeyDownHandler sets this players given key as pressed down
func (p *Player) KeyDownHandler(key int) {
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

// KeyUpHandler sets this players given key as released
func (p *Player) KeyUpHandler(key int) {
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
