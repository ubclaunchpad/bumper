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

func (p *Player) getColor() string {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.Color
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

func (p *Player) getAngle() float64 {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.Angle
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

func (p *Player) setPoints(points int) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.Points = points
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
	p.setPoints(p.Points + numPoints)
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
		p.setAngle(math.Mod(p.getAngle()+0.1, 360))
	}

	if p.getControls().Right {
		p.setAngle(math.Mod(p.getAngle()-0.1, 360))
	}

	if p.getControls().Up {
		controlsVector.Dy = (0.5 * PlayerRadius * math.Cos(p.getAngle()))
		controlsVector.Dx = (0.5 * PlayerRadius * math.Sin(p.getAngle()))
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
	velocityVector := p.getVelocity()
	velocityVector.Dx *= JunkBounceFactor
	velocityVector.Dy *= JunkBounceFactor
	p.setVelocity(velocityVector)
}

// HitPlayer calculates collision, update Player's velocity based on calculation of hitting another player
func (p *Player) HitPlayer(ph *Player) {
	if p.getPDebounce() != 0 {
		return
	}

	pInitialVelocity := p.getVelocity()
	pVelocity := pInitialVelocity
	phVelocity := ph.getVelocity()

	//Calculate player's new velocity
	pVelocity.Dx = (pVelocity.Dx * -VelocityTransferFactor) + (phVelocity.Dx * VelocityTransferFactor)
	pVelocity.Dy = (pVelocity.Dy * -VelocityTransferFactor) + (phVelocity.Dy * VelocityTransferFactor)

	//Calculate hit player's new velocity
	phVelocity.Dx = (phVelocity.Dx * -VelocityTransferFactor) + (pInitialVelocity.Dx * VelocityTransferFactor)
	phVelocity.Dy = (phVelocity.Dy * -VelocityTransferFactor) + (pInitialVelocity.Dy * VelocityTransferFactor)

	p.setVelocity(pVelocity)
	ph.setVelocity(phVelocity)
	ph.setLastPlayerHit(p)
	p.setLastPlayerHit(ph)
	p.setPointsDebounce(PointsDebounceTicks)
	ph.setPointsDebounce(PointsDebounceTicks)
	p.setPDebounce(PlayerDebounceTicks)
}

// ApplyGravity applys a vector towards given position
func (p *Player) ApplyGravity(h *Hole) {
	gravityVector := Velocity{0, 0}
	pVelocity := p.getVelocity()
	pPosition := p.getPosition()
	hPosition := h.getPosition()

	gravityVector.Dx = hPosition.X - pPosition.X
	gravityVector.Dy = hPosition.Y - pPosition.Y

	inverseMagnitude := 1.0 / gravityVector.magnitude()
	gravityVector.normalize()

	//Velocity is affected by how close you are, the size of the hole, and a damping factor.
	pVelocity.Dx += gravityVector.Dx * inverseMagnitude * h.getRadius() * gravityDamping
	pVelocity.Dy += gravityVector.Dy * inverseMagnitude * h.getRadius() * gravityDamping

	p.setVelocity(pVelocity)
}

// checkWalls if the player is attempting to exit the walls, reverse their direction
func (p *Player) checkWalls(height float64, width float64) {
	positionVector := p.getPosition()
	velocityVector := p.getVelocity()
	if positionVector.X+PlayerRadius > width {
		positionVector.X = width - PlayerRadius - 1
		velocityVector.Dx *= WallBounceFactor
	} else if positionVector.X-PlayerRadius < 0 {
		positionVector.X = PlayerRadius + 1
		velocityVector.Dx *= WallBounceFactor
	}

	if positionVector.Y+PlayerRadius > height {
		positionVector.Y = height - PlayerRadius - 1
		velocityVector.Dy *= WallBounceFactor
	} else if positionVector.Y-PlayerRadius < 0 {
		positionVector.Y = PlayerRadius + 1
		velocityVector.Dy *= WallBounceFactor
	}

	p.setPosition(positionVector)
	p.setVelocity(velocityVector)
}

// KeyDownHandler sets this players given key as pressed down
func (p *Player) KeyDownHandler(key int) {
	pControls := p.getControls()

	if key == RightKey {
		pControls.Right = true
	} else if key == LeftKey {
		pControls.Left = true
	} else if key == UpKey {
		pControls.Up = true
	} else if key == DownKey {
		pControls.Down = true
	}

	p.setControls(pControls)
}

// KeyUpHandler sets this players given key as released
func (p *Player) KeyUpHandler(key int) {
	pControls := p.getControls()

	if key == RightKey {
		pControls.Right = false
	} else if key == LeftKey {
		pControls.Left = false
	} else if key == UpKey {
		pControls.Up = false
	} else if key == DownKey {
		pControls.Down = false
	}

	p.setControls(pControls)
}
