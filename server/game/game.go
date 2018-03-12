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

func (a *Arena) collisionPlayer() {

	//memo keeps track of collisions calculated between B and A
	//duplicate calculations between A and B will be skipped
	memo := make(map[int]int)
	//Check player collisions
	//Player A collides with Player B
	for _, player := range a.Players {
		//Player B collides with Player A
		for _, playerHit := range a.Players {

			//Player checks for collision on it's self
			//if true, skip the collision calculation
			//Check if the calculation was already done between B and A
			//skip the duplicate calculation if A and B are colliding
			//memo[playerB.ID] evaluates to 0 if it does not exist in the map
			if player == playerHit || memo[playerHit.ID] == playerHit.ID {
				continue
			}

			if areCirclesColliding(player.Position, models.PlayerRadius, playerHit.Position, models.PlayerRadius) {
				//Keep track of already calculated collisions
				memo[playerHit.ID] = player.ID
				player.HitPlayer()
				playerHit.HitPlayer()
			}

		}

		//Check if player hits a junk
		for _, junk := range a.Junk {
			if areCirclesColliding(player.Position, models.PlayerRadius, junk.Position, models.JunkRadius) {
				//Junk calls player.hitJunk function to calculate player state
				junk.HitBy(&player)
				//Assign junk to last recently hit player color/id
				junk.ID = player.ID
				junk.Color = player.Color
			}
		}

	}
}

func (a *Arena) collisionHole() {

	//Loop through all holes in the arena
	for _, hole := range a.Holes {
		//Check if hole collides with a player
		for _, player := range a.Players {
			if areCirclesColliding(player.Position, models.PlayerRadius, hole.Position, hole.Radius) {
				//Player falls into hole
				//TODO: implement events for player falling into hole, removing the player
			}
		}
		//Check if hole collides with junk
		for _, junk := range a.Junk {
			if areCirclesColliding(junk.Position, models.JunkRadius, hole.Position, hole.Radius) {
				//Junk falls into hole

				//Loop through player ID's and add points
				for _, playerPt := range a.Players {
					if playerPt.ID == junk.ID {
						playerPt.Points += models.PointsPerJunk
					}
				}
				//TODO: implement removing junk
			}
		}
	}
}
