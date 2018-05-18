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
	t.Run("Top wall test", func(t *testing.T) {
		testJunkWallCollision(t, Velocity{0, -2}, Position{testWidth / 2, 0 + JunkRadius + 1})
	})
	t.Run("Bottom wall test", func(t *testing.T) {
		testJunkWallCollision(t, Velocity{0, 2}, Position{testWidth / 2, testHeight - JunkRadius - 1})
	})
	t.Run("Left wall test", func(t *testing.T) {
		testJunkWallCollision(t, Velocity{-2, 0}, Position{0 + JunkRadius + 1, testHeight / 2})
	})
	t.Run("Right wall test", func(t *testing.T) {
		testJunkWallCollision(t, Velocity{2, 0}, Position{testWidth - JunkRadius - 1, testHeight / 2})
	})
}

func testJunkWallCollision(t *testing.T, junkVelocity Velocity, junkPosition Position) {

	dyDirection := 1.0
	dxDirection := 1.0

	if junkVelocity.Dy != 0 {
		dyDirection = -1
	}
	if junkVelocity.Dx != 0 {
		dxDirection = -1
	}

	// Create junk near wall moving towards it
	j := CreateJunk(junkPosition)
	j.Velocity = junkVelocity

	j.UpdatePosition(testHeight, testWidth)

	// Junk should have bounced off the wall
	if j.Position.X != junkPosition.X+junkVelocity.Dx*JunkFriction*dxDirection || j.Position.Y != junkPosition.Y+junkVelocity.Dy*JunkFriction*dyDirection {
		t.Error("Error: Junk bounced incorrectly")
	}

	// Junks velocity should have had one direction inverted
	if j.Velocity.Dx != junkVelocity.Dx*JunkFriction*dxDirection || j.Velocity.Dy != junkVelocity.Dy*JunkFriction*dyDirection {
		t.Error("Error: Junk velocity incorrectly affected, top wall test")
	}
}

func TestPlayerJunkCollisions(t *testing.T) {
	t.Run("Bump from top right", func(t *testing.T) { testPlayerBumpJunk(t, Velocity{-testVelocity.Dx, testVelocity.Dy}) })
	t.Run("Bump from bottom left", func(t *testing.T) { testPlayerBumpJunk(t, Velocity{testVelocity.Dx, -testVelocity.Dy}) })
	t.Run("Bump from top left", func(t *testing.T) { testPlayerBumpJunk(t, Velocity{testVelocity.Dx, testVelocity.Dy}) })
	t.Run("Bump from bottom right", func(t *testing.T) { testPlayerBumpJunk(t, Velocity{-testVelocity.Dx, -testVelocity.Dy}) })
}

func testPlayerBumpJunk(t *testing.T, intialPlayerVelocity Velocity) {

	// Create junk
	j := CreateJunk(centerPos)
	intialJunkVelocity := testVelocity
	j.Velocity = intialJunkVelocity

	// Create a Player
	p := new(Player)
	p.Color = "red"
	p.Velocity = intialPlayerVelocity

	// Hit Junk with Player
	j.HitBy(p)

	// Junk should take player's colour and ID
	if j.Color != p.Color || j.LastPlayerHit != p {
		t.Error("Error: Junk Collsion didn't transfer ownership")
	}

	// Junks velocity should have been affected:
	//   either the minimum bump factor or the damped players Velocity
	if intialPlayerVelocity.Dx < 0 {
		if j.Velocity.Dx != -MinimumBump && j.Velocity.Dx != intialPlayerVelocity.Dx*BumpFactor {
			t.Error("Error: Junk velocity incorrectly affected")
		}
	} else {
		if j.Velocity.Dx != MinimumBump && j.Velocity.Dx != intialPlayerVelocity.Dx*BumpFactor {
			t.Error("Error: Junk velocity incorrectly affected")
		}
	}

	if intialPlayerVelocity.Dy < 0 {
		if j.Velocity.Dy != -MinimumBump && j.Velocity.Dy != intialPlayerVelocity.Dy*BumpFactor {
			t.Error("Error: Junk velocity incorrectly affected")
		}
	} else {
		if j.Velocity.Dy != MinimumBump && j.Velocity.Dy != intialPlayerVelocity.Dy*BumpFactor {
			t.Error("Error: Junk velocity incorrectly affected")
		}
	}

	// Collision also affects Players velocity
	if p.Velocity.Dx != intialPlayerVelocity.Dx*JunkBounceFactor || p.Velocity.Dy != intialPlayerVelocity.Dy*JunkBounceFactor {
		t.Error("Error: Player velocity not affected")
	}

	// Second collision right away should have no effect becuase of the debounce period.
	lastVelocity := j.Velocity
	j.HitBy(p)
	if j.Velocity.Dx != lastVelocity.Dx || j.Velocity.Dy != lastVelocity.Dy {
		t.Error("Error: Junk/Player collision debouncing failed")
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
		{"Hole NE", Position{centerPos.X + 1, centerPos.Y - 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			h := CreateHole(tc.holePosition)
			j := CreateJunk(centerPos)

			vector := Velocity{h.Position.X - j.Position.X, h.Position.Y - j.Position.Y}
			j.ApplyGravity(&h)

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
	intialJunkVelocity := Velocity{testVelocity.Dx, testVelocity.Dy}
	j1.Velocity = intialJunkVelocity

	j2 := CreateJunk(centerPos)
	otherJunkVelocity := Velocity{-testVelocity.Dx, testVelocity.Dy}
	j2.Velocity = otherJunkVelocity

	// Hit Junk with Other Junk
	j1.HitJunk(&j2)

	// Both Junk's velocities should have been affected, not black boxed :(
	if j1.Velocity.Dx != (intialJunkVelocity.Dx*-JunkVTransferFactor)+(otherJunkVelocity.Dx*JunkVTransferFactor) ||
		j1.Velocity.Dy != (intialJunkVelocity.Dy*-JunkVTransferFactor)+(otherJunkVelocity.Dy*JunkVTransferFactor) {
		t.Error("Error: Junk 1's velocity incorrectly affected")
	}

	if j2.Velocity.Dx != (otherJunkVelocity.Dx*-JunkVTransferFactor)+(intialJunkVelocity.Dx*JunkVTransferFactor) ||
		j2.Velocity.Dy != (otherJunkVelocity.Dy*-JunkVTransferFactor)+(intialJunkVelocity.Dy*JunkVTransferFactor) {
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
