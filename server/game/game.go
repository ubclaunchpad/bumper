package game

import (
	"math"
	"math/rand"

	"github.com/ubclaunchpad/bumper/server/models"
)

const (
	playerRadius       = 25
	junkCount          = 10
	holeCount          = 10
	junkRadius         = 8
	minHoleRadius      = 15
	maxHoleRadius      = 30
	minDistanceBetween = maxHoleRadius
	minHoleLife        = 25
	maxHoleLife        = 75
)

// Arena container for play area information including all objects
type Arena struct {
	Height  float64 // Height of play area in pixels
	Width   float64 // Width of play area in pixels
	Holes   []models.Hole
	Junk    []models.Junk
	Players []models.Player
}

// CreateArena constructor for arena
// initializes holes and junk
func CreateArena(height float64, width float64) *Arena {
	a := Arena{height, width, nil, nil, nil}

	// create holes
	for i := 0; i < holeCount; i++ {
		position := a.generateCoord(minHoleRadius)
		initialRadius := math.Floor(rand.Float64()*((maxHoleRadius-minHoleRadius)+1)) + minHoleRadius
		lifespan := math.Floor(rand.Float64()*((maxHoleLife-minHoleLife)+1)) + minHoleLife
		hole := models.Hole{
			Position: position,
			Radius:   initialRadius,
			Life:     lifespan,
		}
		a.Holes = append(a.Holes, hole)
	}

	// create junk
	for i := 0; i < junkCount; i++ {
		position := a.generateCoord(junkRadius)
		junk := models.Junk{
			Position: position,
			Velocity: models.Velocity{Dx: 0, Dy: 0},
			Color:    "white",
			ID:       0}
		a.Junk = append(a.Junk, junk)
	}

	return &a
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
	}
}

func (a *Arena) isPositionValid(position models.Position) bool {
	for _, hole := range a.Holes {
		if areCirclesColliding(hole.Position, hole.Radius, position, minDistanceBetween) {
			return false
		}
	}
	for _, junk := range a.Junk {
		if areCirclesColliding(junk.Position, junkRadius, position, minDistanceBetween) {
			return false
		}
	}
	for _, player := range a.Players {
		if areCirclesColliding(player.Position, playerRadius, position, minDistanceBetween) {
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

// Hello test
func (a *Arena) Hello() {
}
