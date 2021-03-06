package arena

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Arena related constants
const (
	MinDistanceBetween = models.MaxHoleRadius
)

// MessageChannel is used by the server to emit messages to a client (injected global from Main)
var MessageChannel chan models.Message

// Arena container for play area information including all objects
type Arena struct {
	rwMutex sync.RWMutex
	Height  float64
	Width   float64
	Holes   []*models.Hole
	Junk    []*models.Junk
	Players map[string]*models.Player
}

// CreateArena constructor for arena initializes holes and junk
func CreateArena(height float64, width float64, holeCount int, junkCount int) *Arena {
	a := Arena{
		sync.RWMutex{},
		height,
		width,
		make([]*models.Hole, 0, holeCount),
		make([]*models.Junk, 0, junkCount),
		make(map[string]*models.Player),
	}

	for i := 0; i < holeCount; i++ {
		a.addHole()
	}

	for i := 0; i < junkCount; i++ {
		a.addJunk()
	}

	return &a
}

// GetHoles returns a list of holes
func (a *Arena) GetHoles() []*models.Hole {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	return a.Holes
}

// GetJunk returns a list of junk
func (a *Arena) GetJunk() []*models.Junk {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	return a.Junk
}

// GetPlayers returns a list of players
func (a *Arena) GetPlayers() []*models.Player {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	players := make([]*models.Player, 0, len(a.Players))
	for _, player := range a.Players {
		players = append(players, player)
	}
	return players
}

// UpdatePositions calculates the next state of each object
func (a *Arena) UpdatePositions() {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	for i, hole := range a.Holes {
		hole.Update()
		if hole.IsDead() {
			a.removeHole(i)
			a.addHole()
		}
	}
	for _, junk := range a.Junk {
		junk.UpdatePosition(a.Height, a.Width)
	}
	for _, player := range a.Players {
		// check whether player has been spawned
		if player.GetName() != "" {
			player.UpdatePosition(a.Height, a.Width)
		}
	}
}

// CollisionDetection loops through players and holes and determines if a collision has occurred
func (a *Arena) CollisionDetection() {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.playerCollisions()
	a.holeCollisions()
	a.junkCollisions()
}

// GetState assembles an UpdateMessage from the current state of the arena
func (a *Arena) GetState() *models.UpdateMessage {
	return &models.UpdateMessage{
		Holes:   a.GetHoles(),
		Junk:    a.GetJunk(),
		Players: a.GetPlayers(),
	}
}

// AddPlayer adds a new player to the arena
// player has no position or name until spawned
// TODO player has no color until spawned
func (a *Arena) AddPlayer(ws *websocket.Conn) (*models.Player, error) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	color, err := a.generateRandomColor()
	if err != nil {
		return nil, err
	}

	p := models.CreatePlayer("", color, ws)
	a.Players[p.GetID()] = p
	return p, nil
}

// GetPlayer gets the specified player
func (a *Arena) GetPlayer(id string) *models.Player {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()
	return a.Players[id]
}

// RemovePlayer removes the specified player from the arena
func (a *Arena) RemovePlayer(p *models.Player) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	delete(a.Players, p.GetID())
}

// SpawnPlayer spawns the player with a position on the map
// TODO choose color here as well
func (a *Arena) SpawnPlayer(id string, name string, country string) error {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	position := a.generateCoordinate(models.PlayerRadius)
	a.Players[id].Position = position
	a.Players[id].Name = name
	a.Players[id].Country = country
	return nil
}

// generateCoordinate creates a position coordinate
// coordinates are constrained within the Arena's width/height and spacing
// they are all valid
func (a *Arena) generateCoordinate(objectRadius float64) models.Position {
	maxWidth := a.Width - objectRadius
	maxHeight := a.Height - objectRadius

	dummy := models.Hole{
		Position: models.Position{},
		Radius:   MinDistanceBetween,
	}
	for {
		x := math.Floor(rand.Float64()*(maxWidth)) + objectRadius
		y := math.Floor(rand.Float64()*(maxHeight)) + objectRadius
		dummy.Position = models.Position{X: x, Y: y}
		if a.isPositionValid(&dummy) {
			return dummy.Position
		}
	}
}

func (a *Arena) isPositionValid(obj models.Object) bool {
	for _, hole := range a.Holes {
		if areCirclesColliding(hole, obj) {
			return false
		}
	}
	for _, junk := range a.Junk {
		if areCirclesColliding(junk, obj) {
			return false
		}
	}
	for _, player := range a.Players {
		if areCirclesColliding(player, obj) {
			return false
		}
	}

	return true
}

// detect collision between objects
// (x2-x1)^2 + (y1-y2)^2 <= (r1+r2)^2
func areCirclesColliding(obj models.Object, other models.Object) bool {
	p := obj.GetPosition()
	q := other.GetPosition()
	return math.Pow(p.X-q.X, 2)+math.Pow(p.Y-q.Y, 2) <= math.Pow(obj.GetRadius()+other.GetRadius(), 2)
}

/*
collisionPlayer checks for collisions between players to junk, holes, and other players
Duplicate calculations are kept track of using the memo map to store collisions detected
between player-to-player.
*/
func (a *Arena) playerCollisions() {
	memo := make(map[*models.Player]*models.Player)
	for _, player := range a.Players {
		for _, playerHit := range a.Players {
			if player == playerHit || memo[player] == playerHit {
				continue
			}
			if areCirclesColliding(player, playerHit) {
				memo[playerHit] = player
				player.HitPlayer(playerHit)
			}
		}
		for _, junk := range a.Junk {
			if areCirclesColliding(player, junk) {
				junk.HitBy(player)
			}
		}
	}
}

func (a *Arena) holeCollisions() {
	for _, hole := range a.Holes {
		if !hole.IsAlive {
			continue
		}

		gravityField := models.Hole{
			Position: hole.GetPosition(),
			Radius:   hole.GetGravityRadius(),
		}

		for name, player := range a.Players {
			if areCirclesColliding(player, hole) {
				playerScored := player.LastPlayerHit
				if playerScored != nil {
					playerScored.AddPoints(models.PointsPerPlayer)
					// go database.UpdatePlayerScore(playerScored)
				}

				deathMsg := models.Message{
					Type: "death",
					Data: name,
				}
				MessageChannel <- deathMsg
			} else if areCirclesColliding(player, gravityField) {
				player.ApplyGravity(hole)
			}
		}

		for i, junk := range a.Junk {
			if areCirclesColliding(junk, hole) {
				playerScored := junk.LastPlayerHit
				if playerScored != nil {
					playerScored.AddPoints(models.PointsPerJunk)
				}

				a.removeJunk(i)
				a.addJunk()
			} else if areCirclesColliding(junk, gravityField) {
				junk.ApplyGravity(hole)
			}
		}
	}
}

// Checks for junk on junk collisions
func (a *Arena) junkCollisions() {
	memo := make(map[*models.Junk]*models.Junk)
	for _, junk := range a.Junk {
		for _, junkHit := range a.Junk {
			if junk == junkHit || memo[junkHit] == junk {
				continue
			}
			if areCirclesColliding(junk, junkHit) {
				memo[junkHit] = junk
				junk.HitJunk(junkHit)
			}
		}
	}
}

// generate random hex value
func (a *Arena) generateRandomColor() (string, error) {
	letterSet := [13]string{"3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	colorSet := make(map[string]bool)
	for _, player := range a.Players {
		colorSet[player.Color] = true
	}

	var (
		color   string
		timeout int
	)
	for {
		var buffer bytes.Buffer
		buffer.WriteString("#")
		for i := 0; i < 6; i++ {
			c := letterSet[rand.Intn(12)]
			buffer.WriteString(c)
		}
		color = buffer.String()

		if !colorSet[color] {
			return color, nil
		}

		timeout++
		if timeout == 5 {
			return "", errors.New("Cannot generate unique random color")
		}
	}
}

// adds a junk in a random spot
func (a *Arena) addJunk() {
	position := a.generateCoordinate(models.JunkRadius)
	junk := models.CreateJunk(position)
	a.Junk = append(a.Junk, junk)
}

// remove junk without considering order
func (a *Arena) removeJunk(index int) bool {
	if len(a.Junk) < index+1 {
		return false
	}

	a.Junk[index] = a.Junk[len(a.Junk)-1]
	a.Junk[len(a.Junk)-1] = nil
	a.Junk = a.Junk[:len(a.Junk)-1]
	return true
}

// adds a hole in a random spot
func (a *Arena) addHole() {
	h := models.CreateHole(a.generateCoordinate(models.MinHoleRadius))
	a.Holes = append(a.Holes, h)
}

// remove hole without considering order
func (a *Arena) removeHole(index int) bool {
	if len(a.Holes) < index+1 {
		return false
	}

	a.Holes[index] = a.Holes[len(a.Holes)-1]
	a.Holes[len(a.Holes)-1] = nil
	a.Holes = a.Holes[:len(a.Holes)-1]
	return true
}
