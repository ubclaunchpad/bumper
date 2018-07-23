package arena

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/database"
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

// UpdatePositions calculates the next state of each object
func (a *Arena) UpdatePositions() {
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
		if player.Name != "" {
			player.UpdatePosition(a.Height, a.Width)
		}
	}
}

// CollisionDetection loops through players and holes and determines if a collision has occurred
func (a *Arena) CollisionDetection() {
	a.playerCollisions()
	a.holeCollisions()
	a.junkCollisions()
}

// GetState assembles an UpdateMessage from the current state of the arena
func (a *Arena) GetState() *models.UpdateMessage {
	players := make([]*models.Player, 0, len(a.Players))
	for _, player := range a.Players {
		players = append(players, player)
	}

	return &models.UpdateMessage{
		Holes:   a.Holes,
		Junk:    a.Junk,
		Players: players,
	}
}

// AddPlayer adds a new player to the arena
// player has no position or name until spawned
// TODO player has no color until spawned
func (a *Arena) AddPlayer(id string, ws *websocket.Conn) error {
	color, err := a.generateRandomColor()
	if err != nil {
		return err
	}
	a.Players[id] = models.CreatePlayer(id, "", models.Position{}, color, ws)
	return nil
}

// SpawnPlayer spawns the player with a position on the map
// TODO choose color here as well
func (a *Arena) SpawnPlayer(id string, name string, country string) error {
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
	for {
		x := math.Floor(rand.Float64()*(maxWidth)) + objectRadius
		y := math.Floor(rand.Float64()*(maxHeight)) + objectRadius
		position := models.Position{X: x, Y: y}
		if a.isPositionValid(position) {
			return position
		}

		// TODO: Add a timeout here; return error here
	}
}

func (a *Arena) isPositionValid(position models.Position) bool {
	for _, hole := range a.Holes {
		if areCirclesColliding(hole.Position, hole.GravityRadius, position, MinDistanceBetween) {
			return false
		}
	}
	for _, junk := range a.Junk {
		if areCirclesColliding(junk.Position, models.JunkRadius, position, MinDistanceBetween) {
			return false
		}
	}
	for _, player := range a.Players {
		if areCirclesColliding(player.Position, models.PlayerRadius, position, MinDistanceBetween) {
			return false
		}
	}

	return true
}

// detect collision between objects
// (x2-x1)^2 + (y1-y2)^2 <= (r1+r2)^2
func areCirclesColliding(p models.Position, r1 float64, q models.Position, r2 float64) bool {
	return math.Pow(p.X-q.X, 2)+math.Pow(p.Y-q.Y, 2) <= math.Pow(r1+r2, 2)
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
			if player == playerHit || memo[playerHit] == player {
				continue
			}
			if areCirclesColliding(player.Position, models.PlayerRadius, playerHit.Position, models.PlayerRadius) {
				memo[playerHit] = player
				player.HitPlayer(playerHit, a.Height, a.Width)
			}
		}
		for _, junk := range a.Junk {
			if areCirclesColliding(player.Position, models.PlayerRadius, junk.Position, models.JunkRadius) {
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

		for name, player := range a.Players {
			if areCirclesColliding(player.Position, models.PlayerRadius, hole.Position, hole.Radius) {
				// TODO: Should award some points to the bumper... Not as straight forward as the junk
				go database.UpdatePlayerScore(player)
				deathMsg := models.Message{
					Type: "death",
					Data: name,
				}
				MessageChannel <- deathMsg
			} else if areCirclesColliding(player.Position, models.PlayerRadius, hole.Position, hole.GravityRadius) {
				player.ApplyGravity(hole)
			}
		}

		for i, junk := range a.Junk {
			if areCirclesColliding(junk.Position, models.JunkRadius, hole.Position, hole.Radius) {
				playerScored := junk.LastPlayerHit
				if playerScored != nil {
					playerScored.AddPoints(models.PointsPerJunk)
				}

				a.removeJunk(i)
				a.addJunk()
			} else if areCirclesColliding(junk.Position, models.JunkRadius, hole.Position, hole.GravityRadius) {
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
			if areCirclesColliding(junk.Position, models.JunkRadius, junkHit.Position, models.JunkRadius) {
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
