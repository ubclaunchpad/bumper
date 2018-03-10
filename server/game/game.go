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

// createArena constructor for arena to initialize a width/height
// TODO add initial junk and holes
// returns an Arena struct with width and height
func createArena(height float64, width float64) *Arena {
	a := Arena{height, width, nil, nil, nil}

	// create holes
	for i := 0; i < holeCount; i++ {
		position := a.generateCoord(minHoleRadius)
		hole := models.Hole{*position, minHoleRadius}
		a.Holes = append(a.Holes, hole)
	}

	// create junk
	for i := 0; i < junkCount; i++ {
		position := a.generateCoord(junkRadius)
		junk := models.Junk{*position, models.Velocity{0, 0}, 0}
		a.Junk = append(a.Junk, junk)
	}

	return &Arena{height, width, nil, nil, nil}
}

// TODO does this need to be a function on arena?
// generateCoord creates a position coordinate
// coordinates are constrained within the Arena's width/height and spacing
func (a *Arena) generateCoord(objectRadius float64) *models.Position {
	maxWidth := a.Width - objectRadius
	maxHeight := a.Height - objectRadius

	x := math.Floor(rand.Float64()*(maxWidth)) + objectRadius
	y := math.Floor(rand.Float64()*(maxHeight)) + objectRadius

	return &models.Position{x, y}
}
