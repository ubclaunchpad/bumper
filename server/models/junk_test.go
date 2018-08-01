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
	if j.Body.Position.X != initialPosition.X || j.Body.Position.Y != initialPosition.Y {
		t.Error("Error: Still Junk moved")
	}

	// Apply vector
	j.Body.Velocity = testVelocity
	j.UpdatePosition(testHeight, testWidth)

	// Junk should have moved in that direction, but not more than the velocity
	if j.Body.Position.X != initialPosition.X+testVelocity.Dx*JunkFriction || j.Body.Position.Y != initialPosition.Y+testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk moved incorrectly")
	}

	// Junks velocity should have had friction applied.
	if j.Body.Velocity.Dx != testVelocity.Dx*JunkFriction || j.Body.Velocity.Dy != testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk friction not applied")
	}

	// Update Position Again
	lastPosition := j.Body.Position
	j.UpdatePosition(testHeight, testWidth)

	// Junk should have moved in that direction again
	if j.Body.Position.X != lastPosition.X+testVelocity.Dx*JunkFriction*JunkFriction || j.Body.Position.Y > lastPosition.Y+testVelocity.Dy*JunkFriction*JunkFriction {
		t.Error("Error: Junk moved incorrectly")
	}

	// Junks velocity should have had friction applied again
	if j.Body.Velocity.Dx != testVelocity.Dx*JunkFriction*JunkFriction || j.Body.Velocity.Dy != testVelocity.Dy*JunkFriction*JunkFriction {
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
			j.Body.Velocity = tc.junkVelocity

			j.UpdatePosition(testHeight, testWidth)

			// Junk should have bounced off the wall
			if j.Body.Position.X != tc.junkPosition.X+tc.junkVelocity.Dx*JunkFriction*dxDirection || j.Body.Position.Y != tc.junkPosition.Y+tc.junkVelocity.Dy*JunkFriction*dyDirection {
				t.Error("Error: Junk bounced incorrectly")
			}

			// Junks velocity should have had one direction inverted
			if j.Body.Velocity.Dx != tc.junkVelocity.Dx*JunkFriction*dxDirection || j.Body.Velocity.Dy != tc.junkVelocity.Dy*JunkFriction*dyDirection {
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
			j.Body.Velocity = initialJunkVelocity

			// Create a Player
			p := new(Player)
			p.Color = "red"
			p.Body.Velocity = tc.initialPlayerVelocity

			// Hit Junk with Player
			j.HitBy(p)

			// Junk should take player's colour and ID
			if j.Color != p.Color || j.LastPlayerHit != p {
				t.Error("Error: Junk Collsion didn't transfer ownership")
			}

			direction := Velocity{0, 0}
			direction.Dx = j.Body.Position.X - p.Body.Position.X
			direction.Dy = j.Body.Position.Y - p.Body.Position.Y
			direction.normalize()

			minimumVelocity := Velocity{MinimumBump, MinimumBump}

			// Junks velocity should have been affected in the correct direction and at least minimum amount
			if !checkDirection(direction, j.Body.Velocity) || j.Body.Velocity.magnitude() < minimumVelocity.magnitude() {
				t.Error("Error: Junk not bumped in correct direction or hard enough")
			}

			// Collision also affects Players velocity
			if p.Body.Velocity.Dx != tc.initialPlayerVelocity.Dx*1 || p.Body.Velocity.Dy != tc.initialPlayerVelocity.Dy*1 {
				t.Error("Error: Player velocity not affected")
			}

			// Second collision right away should have no effect because of the debounce period.
			lastVelocity := j.Body.Velocity
			j.HitBy(p)
			if j.Body.Velocity.Dx != lastVelocity.Dx || j.Body.Velocity.Dy != lastVelocity.Dy {
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

			vector := Velocity{h.Body.Position.X - j.Body.Position.X, h.Body.Position.Y - j.Body.Position.Y}
			h.ApplyGravity(&j.Body, JunkGravityDamping)

			if !checkDirection(vector, j.Body.Velocity) {
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
	j1.Body.Velocity = initialJunkVelocity

	j2 := CreateJunk(centerPos)
	otherJunkVelocity := Velocity{-testVelocity.Dx, testVelocity.Dy}
	j2.Body.Velocity = otherJunkVelocity

	// Hit Junk with Other Junk
	j1.HitJunk(j2)

	// Both Junk's velocities should have been affected, not black boxed :(
	if j1.Body.Velocity.Dx != (initialJunkVelocity.Dx*-JunkVTransferFactor)+(otherJunkVelocity.Dx*JunkVTransferFactor) ||
		j1.Body.Velocity.Dy != (initialJunkVelocity.Dy*-JunkVTransferFactor)+(otherJunkVelocity.Dy*JunkVTransferFactor) {
		t.Error("Error: Junk 1's velocity incorrectly affected")
	}

	if j2.Body.Velocity.Dx != (otherJunkVelocity.Dx*-JunkVTransferFactor)+(initialJunkVelocity.Dx*JunkVTransferFactor) ||
		j2.Body.Velocity.Dy != (otherJunkVelocity.Dy*-JunkVTransferFactor)+(initialJunkVelocity.Dy*JunkVTransferFactor) {
		t.Error("Error: Junk 2's velocity incorrectly affected")
	}

	// Second collision right away should have no effect because of the debounce period.
	lastVelocity := j1.Body.Velocity
	j1.HitJunk(j2)
	if j1.Body.Velocity.Dx != lastVelocity.Dx || j1.Body.Velocity.Dy != lastVelocity.Dy {
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
