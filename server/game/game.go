package game

import (
	"math"
	"math/rand"

	"github.com/ubclaunchpad/bumper/server/models"
)

// Game related constants
const (
	junkCount          = 10
	holeCount          = 10
	minDistanceBetween = models.MaxHoleRadius
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
	for i := 0; i < holeCount; i++ {
		position := a.generateCoord(models.MinHoleRadius)
		initialRadius := math.Floor(rand.Float64()*((models.MaxHoleRadius-models.MinHoleRadius)+1)) + models.MinHoleRadius
		lifespan := math.Floor(rand.Float64()*((models.MaxHoleLife-models.MinHoleLife)+1)) + models.MinHoleLife
		hole := models.Hole{
			Position: position,
			Radius:   initialRadius,
			Life:     lifespan,
		}
		a.Holes = append(a.Holes, hole)
	}

	// create junk
	for i := 0; i < junkCount; i++ {
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
	// for _, player := range a.Players {
	// 	// TODO: Player to Player
	// 	// TODO: Player to Junk
	// }
	// for _, hole := range a.Holes {
	// 	// TODO: Hole to Player
	// 	// TODO: Hole to Junk
	// }
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
		if areCirclesColliding(hole.Position, hole.Radius, position, minDistanceBetween) {
			return false
		}
	}
	for _, junk := range a.Junk {
		if areCirclesColliding(junk.Position, models.JunkRadius, position, minDistanceBetween) {
			return false
		}
	}
	for _, player := range a.Players {
		if areCirclesColliding(player.Position, models.PlayerRadius, position, minDistanceBetween) {
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

func (a *Arena) collisionPlayerToPlayer() {
	//Check player collisions
	//Player A collides with Player B
	for _, playerA := range a.Players {
		//Player B collides with Player A
		for _, playerB := range a.Players {

			//Player checks for collision on it's self
			//if true, skip the collision calculation
			if playerA == playerB {
				continue
			}

			//TODO: Add logic to only calculate player-player collisions once
			if areCirclesColliding(playerA.Position, models.PlayerRadius, playerB.Position, models.PlayerRadius) {
				// playerA.hitPlayer()
				// playerB.hitPlayer()
			}

		}
	}
}
