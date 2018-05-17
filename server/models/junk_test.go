package models

import (
	"fmt"
	"math"
	"testing"
)

const (
	testHeight = 400
	testWidth  = 800
)

var testVelocity = Velocity{1, 1}
var centerPos = Position{testWidth / 2, testHeight / 2}

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
	for i := 0; i < 4; i++ {
		t.Run(fmt.Sprintf("Wall test %d", i), func(t *testing.T) { testJunkWallCollision(t, i) })
	}
}

func testJunkWallCollision(t *testing.T, wall int) {

	testVelocity := Velocity{0, 0}
	initialPosition := centerPos
	dyDirection := 1.0
	dxDirection := 1.0

	switch wall {
	case 0: // Top wall
		testVelocity = Velocity{0, -2}
		initialPosition = Position{testWidth / 2, 0 + JunkRadius + 1}
		dyDirection = -1
	case 1: // Bottom wall
		testVelocity = Velocity{0, 2}
		initialPosition = Position{testWidth / 2, testHeight - JunkRadius - 1}
		dyDirection = -1
	case 2: // Left wall
		testVelocity = Velocity{-2, 0}
		initialPosition = Position{0 + JunkRadius + 1, testHeight / 2}
		dxDirection = -1
	case 3: // Right wall
		testVelocity = Velocity{2, 0}
		initialPosition = Position{testWidth - JunkRadius - 1, testHeight / 2}
		dxDirection = -1
	default:
		t.Error("Error: Invalid Wall specified")
	}

	// Create junk near wall moving towards it
	j := CreateJunk(initialPosition)
	j.Velocity = testVelocity

	j.UpdatePosition(testHeight, testWidth)

	// Junk should have bounced off the wall
	if j.Position.X != initialPosition.X+testVelocity.Dx*JunkFriction*dxDirection || j.Position.Y != initialPosition.Y+testVelocity.Dy*JunkFriction*dyDirection {
		t.Error("Error: Junk bounced incorrectly")
	}

	// Junks velocity should have had one direction inverted
	if j.Velocity.Dx != testVelocity.Dx*JunkFriction*dxDirection || j.Velocity.Dy != testVelocity.Dy*JunkFriction*dyDirection {
		t.Error("Error: Junk velocity incorrectly affected, top wall test")
	}
}

func TestPlayerJunkCollisions(t *testing.T) {
	for i := 0; i < 4; i++ {
		t.Run(fmt.Sprintf("Direction test %d", i), func(t *testing.T) { testPlayerBumpJunk(t, i) })
	}
}

func testPlayerBumpJunk(t *testing.T, direction int) {

	// Create junk
	j := CreateJunk(centerPos)
	intialJunkVelocity := testVelocity
	j.Velocity = intialJunkVelocity

	// Create a Player
	p := new(Player)
	p.Color = "red"
	intialPlayerVelocity := Velocity{0, 0}
	switch direction {
	case 0:
		intialPlayerVelocity = Velocity{-testVelocity.Dx, testVelocity.Dy}
	case 1:
		intialPlayerVelocity = Velocity{testVelocity.Dx, -testVelocity.Dy}
	case 2:
		intialPlayerVelocity = Velocity{testVelocity.Dx, testVelocity.Dy}
	case 3:
		intialPlayerVelocity = Velocity{-testVelocity.Dx, -testVelocity.Dy}
	default:
		t.Error("Error: Invalid Direction specified")
	}
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
	for i := 0; i < 4; i++ {
		t.Run(fmt.Sprintf("Gravity test %d", i), func(t *testing.T) { testJunkGravity(t, i) })
	}
}

// Test that gravity vectors are applied to the junk's velocity in the direction of the hole
func testJunkGravity(t *testing.T, direction int) {

	h := CreateHole(Position{0, 0})
	initialPosition := centerPos

	switch direction { // Create Hole slightly off in a direction to the junk
	case 0:
		h = CreateHole(Position{initialPosition.X - 1, initialPosition.Y + 1})
	case 1:
		h = CreateHole(Position{initialPosition.X + 1, initialPosition.Y + 1})
	case 2:
		h = CreateHole(Position{initialPosition.X - 1, initialPosition.Y - 1})
	case 3:
		h = CreateHole(Position{initialPosition.X + 1, initialPosition.Y - 1})
	default:
		t.Error("Error: Invalid Direction specified")
	}

	j := CreateJunk(initialPosition)

	vector := Velocity{h.Position.X - j.Position.X, h.Position.Y - j.Position.Y}
	j.ApplyGravity(&h)

	if !checkDirection(vector, j.Velocity) {
		t.Error("Error: Gravity wasn't applied in the correct direction")
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
