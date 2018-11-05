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
	LeftKey                 = 37
	RightKey                = 39
	UpKey                   = 38
	DownKey                 = 40
	WallBounceFactor        = -1.5
	PlayerRadius            = 25
	PlayerAcceleration      = 0.5
	PlayerFriction          = 0.97
	MaxVelocity             = 15
	PointsPerJunk           = 100
	PointsPerPlayer         = 500
	PlayerGravityDamping    = 0.075
	PlayerDebounceTicks     = 15
	PointsDebounceTicks     = 100
	PlayerMass              = 3
	PlayerRestitutionFactor = 0.8
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
	Name           string `json:"name"`
	ID             string `json:"id"`
	PhysicsBody    `json:"body"`
	Country        string      `json:"country"`
	Color          string      `json:"color"`
	Angle          float64     `json:"angle"`
	Controls       KeysPressed `json:"-"`
	Points         int         `json:"points"`
	LastPlayerHit  *Player     `json:"-"`
	pointsDebounce int
	pDebounce      int
	rwMutex        *sync.RWMutex
	ws             *websocket.Conn
}

// PlayerMessage contains the data the client needs about a player
type PlayerMessage struct {
	Name     string   `json:"name"`
	ID       string   `json:"id"`
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
	Country  string   `json:"country"`
	Color    string   `json:"color"`
	Angle    float64  `json:"angle"`
	Points   int      `json:"points"`
}

// CreatePlayer constructs an instance of player with
// given position, color, and WebSocket connection
func CreatePlayer(id string, name string, pos Position, color string, ws *websocket.Conn) *Player {
	lock := sync.RWMutex{}
	return &Player{
		Name:           name,
		ID:             id,
		PhysicsBody:    CreateBody(pos, PlayerRadius, PlayerMass, PlayerRestitutionFactor, &lock),
		Color:          color,
		Angle:          math.Pi,
		Controls:       KeysPressed{},
		pDebounce:      0,
		pointsDebounce: 0,
		rwMutex:        &lock,
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

// GenUniqueID generates a unique id string
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
	controlsVector.ApplyFactor(PlayerAcceleration)

	p.ApplyFactor(PlayerFriction)
	p.ApplyVector(controlsVector)

	// Ensure it never gets going too fast
	if p.VelocityMagnitude() > MaxVelocity {
		p.NormalizeVelocity()
		p.ApplyFactor(MaxVelocity)
	}

	// Apply player's velocity vector
	p.ApplyVelocity()

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

// HitPlayer calculates collision, update Player's position based on calculation of hitting another player
func (p *Player) HitPlayer(ph *Player) {
	if p.getPDebounce() != 0 {
		return
	}

	InelasticCollision(&p.PhysicsBody, &ph.PhysicsBody)

	ph.LastPlayerHit = p
	p.LastPlayerHit = ph
	p.pointsDebounce = PointsDebounceTicks
	ph.pointsDebounce = PointsDebounceTicks
	p.pDebounce = PlayerDebounceTicks
}

// checkWalls if the player is attempting to exit the walls, reverse their direction
func (p *Player) checkWalls(height float64, width float64) {
	positionVector := p.GetPosition()
	velocityVector := p.GetVelocity()
	playerRadius := p.GetRadius()

	if positionVector.X+playerRadius > width {
		positionVector.X = width - playerRadius - 1
		velocityVector.Dx *= WallBounceFactor
	} else if positionVector.X-playerRadius < 0 {
		positionVector.X = playerRadius + 1
		velocityVector.Dx *= WallBounceFactor
	}

	if positionVector.Y+playerRadius > height {
		positionVector.Y = height - playerRadius - 1
		velocityVector.Dy *= WallBounceFactor
	} else if positionVector.Y-playerRadius < 0 {
		positionVector.Y = playerRadius + 1
		velocityVector.Dy *= WallBounceFactor
	}

	p.SetPosition(positionVector)
	p.SetVelocity(velocityVector)
}

// checkWalls2 check if the player is attempting to exit the walls, reverse they're direction
func (p *Player) checkWalls2(height float64, width float64) {
	if p.GetX()+p.GetRadius() > width {
		p.ApplyXFactor(WallBounceFactor)
		p.SetX(width - p.GetRadius() - 1)
	} else if p.GetX()-p.GetRadius() < 0 {
		p.ApplyXFactor(WallBounceFactor)
		p.SetX(p.GetRadius() + 1)
	}

	if p.GetY()+p.GetRadius() > height {
		p.ApplyYFactor(WallBounceFactor)
		p.SetY(height - p.GetRadius() - 1)
	} else if p.GetY()-p.GetRadius() < 0 {
		p.ApplyYFactor(WallBounceFactor)
		p.SetY(p.GetRadius() + 1)
	}
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

// MakeMessage returns a PlayerMessage with this player's data
func (p *Player) MakeMessage() *PlayerMessage {
	return &PlayerMessage{
		Name:     p.Name,
		ID:       p.ID,
		Position: p.GetPosition(),
		Velocity: p.GetVelocity(),
		Country:  p.Country,
		Color:    p.Color,
		Angle:    p.Angle,
		Points:   p.Points,
	}
}
