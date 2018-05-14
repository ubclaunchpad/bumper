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

// Game related constants
const (
	JunkCount          = 30
	HoleCount          = 20
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
func CreateArena(height float64, width float64) *Arena {
	a := Arena{sync.RWMutex{}, height, width, nil, nil, nil}
	a.Players = make(map[string]*models.Player)

	for i := 0; i < HoleCount; i++ {
		position := a.generateCoordinate(models.MinHoleRadius)
		hole := models.CreateHole(position)
		a.Holes = append(a.Holes, &hole)
	}

	for i := 0; i < JunkCount; i++ {
		position := a.generateCoordinate(models.JunkRadius)
		junk := models.Junk{
			Position: position,
			Velocity: models.Velocity{},
			Color:    "white",
		}
		a.Junk = append(a.Junk, &junk)
	}

	return &a
}

// UpdatePositions calculates the next state of each object
func (a *Arena) UpdatePositions() {
	for _, hole := range a.Holes {
		hole.Update()
		if hole.Life < 0 {
			hole.StartNewLife(a.generateCoordinate(models.MaxHoleRadius))
		}
	}
	for _, junk := range a.Junk {
		junk.UpdatePosition(a.Height, a.Width)
	}
	for _, player := range a.Players {
		player.UpdatePosition(a.Height, a.Width)
	}
}

// CollisionDetection loops through players and holes and determines if a collision has occurred
func (a *Arena) CollisionDetection() {
	a.playerCollisions()
	a.holeCollisions()
	a.junkCollisions()
}

// AddPlayer adds a new player to the arena
func (a *Arena) AddPlayer(n string, ws *websocket.Conn) error {
	color, err := a.generateRandomColor()
	if err != nil {
		return err
	}

	position := a.generateCoordinate(models.PlayerRadius)
	a.Players[n] = models.CreatePlayer(n, position, color, ws)
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

		// TODO: Add a timeout here
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
				// TODO: send a you're dead signal - err := client.WriteJSON(&msg)
				// Also should award some points to the bumper... Not as straight forward as the junk
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

				// remove that junk from the junk
				a.Junk = append(a.Junk[:i], a.Junk[i+1:]...)
				//create a new junk to hold the count steady
				a.generateJunk()
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
func (a *Arena) generateJunk() {
	position := a.generateCoordinate(models.JunkRadius)
	junk := models.Junk{
		Position: position,
		Velocity: models.Velocity{Dx: 0, Dy: 0},
		Color:    "white",
	}
	a.Junk = append(a.Junk, &junk)
}
