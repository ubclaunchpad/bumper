package game

import (
	"math"
	"math/rand"

	"../models"
)

const junkCount = 10
const holeCount = 10
const junkRadius = 8
const minHoleRadius = 15
const maxHoleRadius = 30
const minDistanceBetween = maxHoleRadius
const minHoleLife = 25
const maxHoleLife = 75

// Arena container for play area information including all objects
type Arena struct {
	Height  float64 // Height of play area in pixels
	Width   float64 // Width of play area in pixels
	Holes   []models.Hole
	Junk    []models.Junk
	Players []models.Player
}

// createArena constructor for arena
// initializes holes and junk
func createArena(height float64, width float64) *Arena {
	a := Arena{height, width, nil, nil, nil}

	// create holes
	for i := 0; i < holeCount; i++ {
		foundValidPos := false
		for !foundValidPos {
			position := a.generateCoord(minHoleRadius)
			if a.isPositionValid(position) {
				foundValidPos = true
				hole := models.Hole{position, minHoleRadius}
				a.Holes = append(a.Holes, hole)
			}
		}
	}

	// create junk
	for i := 0; i < junkCount; i++ {
		foundValidPos := false
		for !foundValidPos {
			position := a.generateCoord(junkRadius)
			if a.isPositionValid(position) {
				foundValidPos = true
				junk := models.Junk{position, models.Velocity{0, 0}, 0}
				a.Junk = append(a.Junk, junk)
			}
		}
	}

	return &a
}

// generateCoord creates a position coordinate
// coordinates are constrained within the Arena's width/height and spacing
func (a *Arena) generateCoord(objectRadius float64) models.Position {
	maxWidth := a.Width - objectRadius
	maxHeight := a.Height - objectRadius

	x := math.Floor(rand.Float64()*(maxWidth)) + objectRadius
	y := math.Floor(rand.Float64()*(maxHeight)) + objectRadius

	return models.Position{x, y}
}

func (a *Arena) isPositionValid(position models.Position) bool {
	for _, hole := range a.Holes {
		if hole.Position == position {
			return false
		}
	}
	for _, junk := range a.Junk {
		if junk.Position == position {
			return false
		}
	}
	for _, player := range a.Players {
		if player.Position == position {
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
