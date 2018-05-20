package models

import (
	"math"
	"testing"
)

const (
	testHeight = 400
	testWidth  = 800
)

var (
	testVelocity = Velocity{1, 1}
	centerPos    = Position{testWidth / 2, testHeight / 2}
)

func TestUpdateJunkPosition(t *testing.T) {

	// Create still junk in middle
	initialPosition := centerPos
	j := CreateJunk(initialPosition)
	j.UpdatePosition(testHeight, testWidth)

	// Junk with no velocity shouldn't move
	if j.Position.X != initialPosition.X || j.Position.Y != initialPosition.Y {
		t.Error("Error: Still Junk moved")
	}

	// Apply vector
	j.Velocity = testVelocity
	j.UpdatePosition(testHeight, testWidth)

	// Junk should have moved in that direction, but not more than the velocity
	if j.Position.X != initialPosition.X+testVelocity.Dx*JunkFriction || j.Position.Y != initialPosition.Y+testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk moved incorrectly")
	}

	// Junks velocity should have had friction applied.
	if j.Velocity.Dx != testVelocity.Dx*JunkFriction || j.Velocity.Dy != testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk friction not applied")
	}

	// Update Position Again
	lastPosition := j.Position
	j.UpdatePosition(testHeight, testWidth)

	// Junk should have moved in that direction again
	if j.Position.X != lastPosition.X+testVelocity.Dx*JunkFriction*JunkFriction || j.Position.Y > lastPosition.Y+testVelocity.Dy*JunkFriction*JunkFriction {
		t.Error("Error: Junk moved incorrectly")
	}

	// Junks velocity should have had friction applied again
	if j.Velocity.Dx != testVelocity.Dx*JunkFriction*JunkFriction || j.Velocity.Dy != testVelocity.Dy*JunkFriction*JunkFriction {
		t.Error("Error: Junk friction not applied")
	}
}

func TestJunkWallConstraints(t *testing.T) {
	testCases := []struct {
		description  string
		junkVelocity Velocity
		junkPosition Position
	}{
		{"Top wall test", Velocity{0, -2}, Position{testWidth / 2, 0 + JunkRadius + 1}},
		{"Bottom wall test", Velocity{0, 2}, Position{testWidth / 2, testHeight - JunkRadius - 1}},
		{"Left wall test", Velocity{-2, 0}, Position{0 + JunkRadius + 1, testHeight / 2}},
		{"Right wall test", Velocity{2, 0}, Position{testWidth - JunkRadius - 1, testHeight / 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dyDirection := 1.0
			dxDirection := 1.0

			if tc.junkVelocity.Dy != 0 {
				dyDirection = -1
			}
			if tc.junkVelocity.Dx != 0 {
				dxDirection = -1
			}

			// Create junk near wall moving towards it
			j := CreateJunk(tc.junkPosition)
			j.Velocity = tc.junkVelocity

			j.UpdatePosition(testHeight, testWidth)

			// Junk should have bounced off the wall
			if j.Position.X != tc.junkPosition.X+tc.junkVelocity.Dx*JunkFriction*dxDirection || j.Position.Y != tc.junkPosition.Y+tc.junkVelocity.Dy*JunkFriction*dyDirection {
				t.Error("Error: Junk bounced incorrectly")
			}

			// Junks velocity should have had one direction inverted
			if j.Velocity.Dx != tc.junkVelocity.Dx*JunkFriction*dxDirection || j.Velocity.Dy != tc.junkVelocity.Dy*JunkFriction*dyDirection {
				t.Error("Error: Junk velocity incorrectly affected, top wall test")
			}
		})
	}
}

func TestPlayerBumpJunk(t *testing.T) {
	testCases := []struct {
		description           string
		initialPlayerVelocity Velocity
	}{
		{"Bump from top right", Velocity{-testVelocity.Dx, testVelocity.Dy}},
		{"Bump from bottom left", Velocity{testVelocity.Dx, -testVelocity.Dy}},
		{"Bump from top left", Velocity{testVelocity.Dx, testVelocity.Dy}},
		{"Bump from bottom right", Velocity{-testVelocity.Dx, -testVelocity.Dy}},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create junk
			j := CreateJunk(centerPos)
			initialJunkVelocity := testVelocity
			j.Velocity = initialJunkVelocity

			// Create a Player
			p := new(Player)
			p.Color = "red"
			p.Velocity = tc.initialPlayerVelocity

			// Hit Junk with Player
			j.HitBy(p)

			// Junk should take player's colour and ID
			if j.Color != p.Color || j.LastPlayerHit != p {
				t.Error("Error: Junk Collsion didn't transfer ownership")
			}

			// Junks velocity should have been affected:
			//   either the minimum bump factor or the damped players Velocity
			if tc.initialPlayerVelocity.Dx < 0 {
				if j.Velocity.Dx != -MinimumBump && j.Velocity.Dx != tc.initialPlayerVelocity.Dx*BumpFactor {
					t.Error("Error: Junk velocity incorrectly affected")
				}
			} else {
				if j.Velocity.Dx != MinimumBump && j.Velocity.Dx != tc.initialPlayerVelocity.Dx*BumpFactor {
					t.Error("Error: Junk velocity incorrectly affected")
				}
			}

			if tc.initialPlayerVelocity.Dy < 0 {
				if j.Velocity.Dy != -MinimumBump && j.Velocity.Dy != tc.initialPlayerVelocity.Dy*BumpFactor {
					t.Error("Error: Junk velocity incorrectly affected")
				}
			} else {
				if j.Velocity.Dy != MinimumBump && j.Velocity.Dy != tc.initialPlayerVelocity.Dy*BumpFactor {
					t.Error("Error: Junk velocity incorrectly affected")
				}
			}

			// Collision also affects Players velocity
			if p.Velocity.Dx != tc.initialPlayerVelocity.Dx*JunkBounceFactor || p.Velocity.Dy != tc.initialPlayerVelocity.Dy*JunkBounceFactor {
				t.Error("Error: Player velocity not affected")
			}

			// Second collision right away should have no effect becuase of the debounce period.
			lastVelocity := j.Velocity
			j.HitBy(p)
			if j.Velocity.Dx != lastVelocity.Dx || j.Velocity.Dy != lastVelocity.Dy {
				t.Error("Error: Junk/Player collision debouncing failed")
			}
		})
	}
}

func TestJunkGravity(t *testing.T) {
	testCases := []struct {
		description  string
		holePosition Position
	}{
		{"Hole NW", Position{centerPos.X - 1, centerPos.Y + 1}},
		{"Hole NE", Position{centerPos.X + 1, centerPos.Y + 1}},
		{"Hole SW", Position{centerPos.X - 1, centerPos.Y - 1}},
		{"Hole SE", Position{centerPos.X + 1, centerPos.Y - 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			h := CreateHole(tc.holePosition)
			j := CreateJunk(centerPos)

			vector := Velocity{h.Position.X - j.Position.X, h.Position.Y - j.Position.Y}
			j.ApplyGravity(h)

			if !checkDirection(vector, j.Velocity) {
				t.Error("Error: Gravity wasn't applied in the correct direction")
			}
		})
	}
}

// Test Junk bumping off other junk
func TestJunkBumpJunk(t *testing.T) {

	// Create 2 junk
	j1 := CreateJunk(centerPos)
	initialJunkVelocity := Velocity{testVelocity.Dx, testVelocity.Dy}
	j1.Velocity = initialJunkVelocity

	j2 := CreateJunk(centerPos)
	otherJunkVelocity := Velocity{-testVelocity.Dx, testVelocity.Dy}
	j2.Velocity = otherJunkVelocity

	// Hit Junk with Other Junk
	j1.HitJunk(&j2)

	// Both Junk's velocities should have been affected, not black boxed :(
	if j1.Velocity.Dx != (initialJunkVelocity.Dx*-JunkVTransferFactor)+(otherJunkVelocity.Dx*JunkVTransferFactor) ||
		j1.Velocity.Dy != (initialJunkVelocity.Dy*-JunkVTransferFactor)+(otherJunkVelocity.Dy*JunkVTransferFactor) {
		t.Error("Error: Junk 1's velocity incorrectly affected")
	}

	if j2.Velocity.Dx != (otherJunkVelocity.Dx*-JunkVTransferFactor)+(initialJunkVelocity.Dx*JunkVTransferFactor) ||
		j2.Velocity.Dy != (otherJunkVelocity.Dy*-JunkVTransferFactor)+(initialJunkVelocity.Dy*JunkVTransferFactor) {
		t.Error("Error: Junk 2's velocity incorrectly affected")
	}

	// Second collision right away should have no effect becuase of the debounce period.
	lastVelocity := j1.Velocity
	j1.HitJunk(&j2)
	if j1.Velocity.Dx != lastVelocity.Dx || j1.Velocity.Dy != lastVelocity.Dy {
		t.Error("Error: Junk/Junk collision debouncing failed")
	}
}

// Helper function, compares two velocities to see if they're in the same direction
func checkDirection(v1 Velocity, v2 Velocity) bool {
	v1.normalize()
	v2.normalize()

	// Round floats to 3 decimal places for comparison
	if math.Ceil(v1.Dx*1000)/1000 == math.Ceil(v2.Dx*1000)/1000 && math.Ceil(v1.Dy*1000)/1000 == math.Ceil(v2.Dy*1000)/1000 {
		return true
	}
	return false
}
