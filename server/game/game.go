package game

import (
	"github.com/ubclaunchpad/bumper/server/models"
)

// Arena container for play area information including all objects
type Arena struct {
	Height  int // Height of play area in pixels
	Width   int // Width of play area in pixels
	Holes   []models.Hole
	Junk    []models.Junk
	Players []models.Player
}

// createArena constructor for arena to initialize a width/height
// returns an Arena struct with width and height
func createArena(height int, width int) *Arena {
	a := new(Arena)
	a.Height = height
	a.Width = width
	return a
}
