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
	PlayerMass              = 1.5
	PlayerRestitutionFactor = 0.95
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
	Body           PhysicsBody `json:"body"`
	Color          string      `json:"color"`
	Angle          float64     `json:"angle"`
	Controls       KeysPressed `json:"-"`
	Points         int         `json:"points"`
	LastPlayerHit  *Player     `json:"-"`
	pointsDebounce int
	pDebounce      int
	mutex          sync.Mutex
	ws             *websocket.Conn
}

// CreatePlayer constructs an instance of player with
// given position, color, and WebSocket connection
func CreatePlayer(id string, name string, pos Position, color string, ws *websocket.Conn) *Player {
	return &Player{
		Name:           name,
		ID:             id,
		Body:           CreateBody(pos, PlayerRadius, PlayerMass, PlayerRestitutionFactor),
		Color:          color,
		Angle:          math.Pi,
		Controls:       KeysPressed{},
		pDebounce:      0,
		pointsDebounce: 0,
		mutex:          sync.Mutex{},
		ws:             ws,
	}
}

// GenUniqueID generates a unique id string
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
	p.mutex.Lock()
	defer p.mutex.Unlock()

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
	controlsVector.ApplyFactor(PlayerAcceleration)

	p.Body.Velocity.ApplyFactor(PlayerFriction)
	p.Body.Velocity.ApplyVector(controlsVector)

	// Ensure it never gets going too fast
	if p.Body.Velocity.magnitude() > MaxVelocity {
		p.Body.Velocity.normalize()
		p.Body.Velocity.ApplyFactor(MaxVelocity)
	}

	// Apply player's velocity vector
	p.Body.ApplyVelocity()

	p.checkWalls(height, width)

	if p.pDebounce > 0 {
		p.pDebounce--
	} else {
		p.pDebounce = 0
	}

	if p.pointsDebounce > 0 {
		p.pointsDebounce--
	} else {
		p.LastPlayerHit = nil
		p.pointsDebounce = 0
	}
}

// HitPlayer calculates collision, update Player's position based on calculation of hitting another player
func (p *Player) HitPlayer(ph *Player) {
	if p.pDebounce != 0 {
		return
	}

	InelasticCollision(&p.Body, &ph.Body)

	ph.LastPlayerHit = p
	p.LastPlayerHit = ph
	p.pointsDebounce = PointsDebounceTicks
	ph.pointsDebounce = PointsDebounceTicks
	p.pDebounce = PlayerDebounceTicks
}

// checkWalls check if the player is attempting to exit the walls, reverse they're direction
func (p *Player) checkWalls(height float64, width float64) {
	if p.Body.Position.X+PlayerRadius > width {
		p.Body.Position.X = width - PlayerRadius - 1
		p.Body.Velocity.Dx *= WallBounceFactor
	} else if p.Body.Position.X-PlayerRadius < 0 {
		p.Body.Velocity.Dx *= WallBounceFactor
		p.Body.Position.X = PlayerRadius + 1
	}

	if p.Body.Position.Y+PlayerRadius > height {
		p.Body.Velocity.Dy *= WallBounceFactor
		p.Body.Position.Y = height - PlayerRadius - 1
	} else if p.Body.Position.Y-PlayerRadius < 0 {
		p.Body.Velocity.Dy *= WallBounceFactor
		p.Body.Position.Y = PlayerRadius + 1
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
