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

// CreateArena constructor for arena to initialize a width/height
// returns an Arena struct with width and height
func CreateArena(height int, width int) *Arena {
	return &Arena{height, width, nil, nil, nil}
}

// func (a *Arena) hello() {

// }
