package game

import (
	"bytes"
	"math"
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Game related constants
const (
	JunkCount          = 10
	HoleCount          = 10
	MinDistanceBetween = models.MaxHoleRadius
)

var lastID = 1
var MessageChannel chan models.Message

// Arena container for play area information including all objects
type Arena struct {
	Height  float64                            `json:"height"`
	Width   float64                            `json:"width"`
	Holes   []*models.Hole                     `json:"holes"`
	Junk    []*models.Junk                     `json:"junk"`
	Players map[*websocket.Conn]*models.Player `json:"players"`
}

// CreateArena constructor for arena initializes holes and junk
func CreateArena(height float64, width float64) *Arena {
	a := Arena{height, width, nil, nil, nil}
	a.Players = make(map[*websocket.Conn]*models.Player)

	for i := 0; i < HoleCount; i++ {
		position := a.generateCoord(models.MinHoleRadius)
		hole := models.CreateHole(position)
		a.Holes = append(a.Holes, &hole)
	}

	for i := 0; i < JunkCount; i++ {
		position := a.generateCoord(models.JunkRadius)
		junk := models.Junk{
			Position: position,
			Velocity: models.Velocity{Dx: 0, Dy: 0},
			Color:    "white",
			ID:       nil}
		a.Junk = append(a.Junk, &junk)
	}

	return &a
}

// UpdatePositions calculates the next state of each object
func (a *Arena) UpdatePositions() {
	// for _, hole := range a.Holes {

	// }
	for _, junk := range a.Junk {
		junk.UpdatePosition(a.Height, a.Width)
	}
	for _, player := range a.Players {
		player.UpdatePosition(a.Height, a.Width)
	}
}

// CollisionDetection loops through players and holes and determines if a collision has occurred
func (a *Arena) CollisionDetection() {
	a.collisionPlayer()
	a.collisionHole()
}

// AddPlayer adds a new player to the arena
func (a *Arena) AddPlayer(ws *websocket.Conn) {
	player := models.Player{
		ID:       generateID(),
		Position: a.generateCoord(models.PlayerRadius),
		Velocity: models.Velocity{0, 0},
		Color:    a.generateRandomColor(),
		Angle:    math.Pi,
		Controls: models.KeysPressed{false, false, false, false},
	}
	a.Players[ws] = &player
}

// generateCoord creates a position coordinate
// coordinates are constrained within the Arena's width/height and spacing
// they are all valid
func (a *Arena) generateCoord(objectRadius float64) models.Position {
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
		if areCirclesColliding(hole.Position, hole.Radius, position, MinDistanceBetween) {
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
	return (math.Pow((p.X-q.X), 2) + math.Pow((p.Y-q.Y), 2)) <= math.Pow((r1+r2), 2)
}

/*
collisionPlayer checks for collisions between players to junk, holes, and other players
Duplicate calculations are kept track of using the memo map to store collisions detected
between player-to-player.
*/
func (a *Arena) collisionPlayer() {
	memo := make(map[*models.Player]*models.Player)
	for ws, player := range a.Players {
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
				junk.HitBy(player, ws)
			}
		}
	}
}

func (a *Arena) collisionHole() {
	for _, hole := range a.Holes {
		for client, player := range a.Players {
			if areCirclesColliding(player.Position, models.PlayerRadius, hole.Position, hole.Radius) {
				// TODO: award some points to the bumper... Not as straight forward as the junk
				// send update to client with player missing

				deathMsg := models.Message{
					Type: "death",
					Data: client,
				}
				MessageChannel <- deathMsg
				delete(a.Players, client)
			}
		}
		for i, junk := range a.Junk {
			if areCirclesColliding(junk.Position, models.JunkRadius, hole.Position, hole.Radius) {
				if a.Players[junk.ID] != nil {
					playerScored := a.Players[junk.ID]
					playerScored.AddPoints(models.PointsPerJunk)
				}

				// remove that junk from the junk
				a.Junk = append(a.Junk[:i], a.Junk[i+1:]...)
				//create a new junk to hold the count steady
				a.generateJunk()
			}
		}
	}
}

// generate random hex value
func (a *Arena) generateRandomColor() string {
	found := false
	color := ""

	for !found {
		letters := [13]string{"3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
		var buffer bytes.Buffer
		buffer.WriteString("#")
		for i := 0; i < 6; i++ {
			c := letters[rand.Intn(12)]
			buffer.WriteString(c)
		}
		color = buffer.String()
		for _, p := range a.Players {
			if p.Color == color {
				break
			}
		}
		found = true
	}
	return color
}

// TODO generate player id check whether any current players have this id
func generateID() int {
	id := lastID
	lastID++
	return id
}

// adds a junk in a random spot
func (a *Arena) generateJunk() {
	position := a.generateCoord(models.JunkRadius)
	junk := models.Junk{
		Position: position,
		Velocity: models.Velocity{Dx: 0, Dy: 0},
		Color:    "white",
		ID:       nil}
	a.Junk = append(a.Junk, &junk)
}
