package game

import (
	"math"
	"math/rand"

	"github.com/ubclaunchpad/bumper/server/models"
)

// Game related constants
const (
	JunkCount          = 10
	HoleCount          = 10
	MinDistanceBetween = models.MaxHoleRadius
)

// Arena container for play area information including all objects
type Arena struct {
	Height  float64 // Height of play area in pixels
	Width   float64 // Width of play area in pixels
	Holes   []models.Hole
	Junk    []models.Junk
	Players []models.Player
}

// CreateArena constructor for arena initializes holes and junk
func CreateArena(height float64, width float64) *Arena {
	a := Arena{height, width, nil, nil, nil}

	// create holes
	for i := 0; i < HoleCount; i++ {
		position := a.generateCoord(models.MinHoleRadius)
		hole := models.CreateHole(position)
		a.Holes = append(a.Holes, hole)
	}

	// create junk
	for i := 0; i < JunkCount; i++ {
		position := a.generateCoord(models.JunkRadius)
		junk := models.Junk{
			Position: position,
			Velocity: models.Velocity{Dx: 0, Dy: 0},
			Color:    "white",
			ID:       0}
		a.Junk = append(a.Junk, junk)
	}

	return &a
}

// UpdatePositions calculates the next state of each object
func (a *Arena) UpdatePositions() {
	// for _, hole := range a.Holes {

	// }
	// for _, junk := range a.Junk {

	// }
	// for _, player := range a.Players {

	// }
}

// CollisionDetection loops through players and holes and determines if a collision has occurred
func (a *Arena) CollisionDetection() {
	a.collisionPlayer()
	a.collisionHole()
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
	memo := make(map[int]int)          //memo keeps track of collisions calculated between B and A
	for _, player := range a.Players { //Loop through all players in the arena
		for _, playerHit := range a.Players { //Check player collisions
			if player == playerHit || memo[playerHit.ID] == playerHit.ID { //Skip duplicate calculation
				continue
			}
			if areCirclesColliding(player.Position, models.PlayerRadius, playerHit.Position, models.PlayerRadius) {
				memo[playerHit.ID] = player.ID //Keep track of already calculated collisions
				player.HitPlayer()
				playerHit.HitPlayer()
			}
		}
		for _, junk := range a.Junk { //Check if player hits a junk
			if areCirclesColliding(player.Position, models.PlayerRadius, junk.Position, models.JunkRadius) {
				junk.HitBy(&player)       //Junk calls player.hitJunk function to calculate player state
				junk.ID = player.ID       //Assign junk to last recently hit player id
				junk.Color = player.Color //Assign junk to last recently hit player color
			}
		}
	}
}

func (a *Arena) collisionHole() {
	for _, hole := range a.Holes { //Loop through all holes in the arena
		for _, player := range a.Players { //Check if hole collides with a player
			if areCirclesColliding(player.Position, models.PlayerRadius, hole.Position, hole.Radius) {
				//Player falls into hole
				//TODO: implement events for player falling into hole, removing the player
			}
		}
		for _, junk := range a.Junk { //Check if hole collides with junk
			if areCirclesColliding(junk.Position, models.JunkRadius, hole.Position, hole.Radius) {
				for _, playerPt := range a.Players { //Loop through player ID's and add points
					if playerPt.ID == junk.ID {
						playerPt.Points += models.PointsPerJunk
					}
				}
				//TODO: implement deleting junk
			}
		}
	}
}
