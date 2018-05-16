package models

import (
	"math"
	"testing"
)

const (
	testHeight     = 400
	testWidth      = 800
	testVelocityDx = 1
	testVelocityDy = 1
)

func TestUpdateJunkPosition(t *testing.T) {

	// Create still junk in middle
	j := new(Junk)
	j.Velocity = Velocity{0, 0}
	intialPosition := Position{testWidth / 2, testHeight / 2}
	j.Position = intialPosition
	j.UpdatePosition(testHeight, testWidth)

	// Junk with no velocity shouldn't move
	if j.Position.X != intialPosition.X || j.Position.Y != intialPosition.Y {
		t.Error("Error: Still Junk moved")
	}

	// Apply vector
	testVelocity := Velocity{testVelocityDx, testVelocityDy}
	j.Velocity = testVelocity
	j.UpdatePosition(testHeight, testWidth)

	// Junk should have moved in that direction, but not more than the velocity
	if j.Position.X != intialPosition.X+testVelocityDx*JunkFriction || j.Position.Y != intialPosition.Y+testVelocityDy*JunkFriction {
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

	// Create junk near top wall moving towards it
	j := new(Junk)
	testVelocity := Velocity{0, -2}
	intialPosition := Position{testWidth / 2, 0 + JunkRadius + 1}
	j.Velocity = testVelocity
	j.Position = intialPosition

	j.UpdatePosition(testHeight, testWidth)

	// Junk should have bounced off the top wall
	if j.Position.X != intialPosition.X+testVelocity.Dx*JunkFriction || j.Position.Y != intialPosition.Y-testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk moved incorrectly, top wall test")
	}

	// Junks velocity should have had Dy inverted
	if j.Velocity.Dx != testVelocity.Dx*JunkFriction || j.Velocity.Dy != -testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk velocity incorrectly affected, top wall test")
	}

	// Test Bottom Wall
	testVelocity = Velocity{0, 2}
	intialPosition = Position{testWidth / 2, testHeight - JunkRadius - 1}
	j.Velocity = testVelocity
	j.Position = intialPosition

	j.UpdatePosition(testHeight, testWidth)

	// Junk should have bounced off the Bottom wall
	if j.Position.X != intialPosition.X+testVelocity.Dx*JunkFriction || j.Position.Y != intialPosition.Y-testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk moved incorrectly, bottom wall test")
	}

	// Junks velocity should have had Dy inverted
	if j.Velocity.Dx != testVelocity.Dx*JunkFriction || j.Velocity.Dy != -testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk velocity incorrectly affected, bottom wall test")
	}

	// Test Left Wall
	testVelocity = Velocity{-2, 0}
	intialPosition = Position{0 + JunkRadius + 1, testHeight / 2}
	j.Velocity = testVelocity
	j.Position = intialPosition

	j.UpdatePosition(testHeight, testWidth)

	// Junk should have bounced off the Left wall
	if j.Position.X != intialPosition.X-testVelocity.Dx*JunkFriction || j.Position.Y != intialPosition.Y+testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk moved incorrectly, left wall test")
	}

	// Junks velocity should have had Dx inverted
	if j.Velocity.Dx != -testVelocity.Dx*JunkFriction || j.Velocity.Dy != testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk velocity incorrectly affected, left wall test")
	}

	// Test Right Wall
	testVelocity = Velocity{2, 0}
	intialPosition = Position{testWidth - JunkRadius - 1, testHeight / 2}
	j.Velocity = testVelocity
	j.Position = intialPosition

	j.UpdatePosition(testHeight, testWidth)

	// Junk should have bounced off the Right wall
	if j.Position.X != intialPosition.X-testVelocity.Dx*JunkFriction || j.Position.Y != intialPosition.Y+testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk moved incorrectly, right wall test")
	}

	// Junks velocity should have had Dx inverted
	if j.Velocity.Dx != -testVelocity.Dx*JunkFriction || j.Velocity.Dy != testVelocity.Dy*JunkFriction {
		t.Error("Error: Junk velocity incorrectly affected, right wall test")
	}
}

func TestPlayerJunkCollisions(t *testing.T) {
	testPlayerBumpJunk(t, 0)
	testPlayerBumpJunk(t, 1)
	testPlayerBumpJunk(t, 2)
	testPlayerBumpJunk(t, 3)
}

func testPlayerBumpJunk(t *testing.T, direction int) {

	// Create junk
	j := new(Junk)
	intialJunkVelocity := Velocity{testVelocityDx, testVelocityDy}
	j.Velocity = intialJunkVelocity

	// Create a Player
	p := new(Player)
	p.Color = "red"
	intialPlayerVelocity := Velocity{0, 0}
	switch direction {
	case 0:
		intialPlayerVelocity = Velocity{-testVelocityDx, testVelocityDy}
	case 1:
		intialPlayerVelocity = Velocity{testVelocityDx, -testVelocityDy}
	case 2:
		intialPlayerVelocity = Velocity{testVelocityDx, testVelocityDy}
	case 3:
		intialPlayerVelocity = Velocity{-testVelocityDx, -testVelocityDy}
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

// Test that gravity vectors are applied to the junk's velocity in the direction of the hole
func TestJunkGravity(t *testing.T) {

	// Create Junk
	j := new(Junk)
	intialPosition := Position{testWidth / 2, testHeight / 2}
	// intialJunkVelocity := Velocity{testVelocityDx, testVelocityDy}
	j.Velocity = Velocity{0, 0}
	j.Position = intialPosition

	// Create Holes slightly off in each direction to the junk
	hNW := CreateHole(Position{intialPosition.X - 1, intialPosition.Y + 1})
	hNE := CreateHole(Position{intialPosition.X + 1, intialPosition.Y + 1})
	hSW := CreateHole(Position{intialPosition.X - 1, intialPosition.Y - 1})
	hSE := CreateHole(Position{intialPosition.X + 1, intialPosition.Y - 1})

	vector := Velocity{hNW.Position.X - j.Position.X, hNW.Position.Y - j.Position.Y}
	j.ApplyGravity(&hNW)

	if !checkDirection(vector, j.Velocity) {
		t.Error("Error: Gravity wasn't applied in the correct direction")
	}

	j.Velocity = Velocity{0, 0}
	vector = Velocity{hNE.Position.X - j.Position.X, hNE.Position.Y - j.Position.Y}
	j.ApplyGravity(&hNE)

	if !checkDirection(vector, j.Velocity) {
		t.Error("Error: Gravity wasn't applied in the correct direction")
	}

	j.Velocity = Velocity{0, 0}
	vector = Velocity{hSW.Position.X - j.Position.X, hSW.Position.Y - j.Position.Y}
	j.ApplyGravity(&hSW)

	if !checkDirection(vector, j.Velocity) {
		t.Error("Error: Gravity wasn't applied in the correct direction")
	}

	j.Velocity = Velocity{0, 0}
	vector = Velocity{hSE.Position.X - j.Position.X, hSE.Position.Y - j.Position.Y}
	j.ApplyGravity(&hSE)

	if !checkDirection(vector, j.Velocity) {
		t.Error("Error: Gravity wasn't applied in the correct direction")
	}
}

// Test Junk bumping off other junk
func TestJunkBumpJunk(t *testing.T) {

	// Create 2 junk
	j1 := new(Junk)
	intialJunkVelocity := Velocity{testVelocityDx, testVelocityDy}
	j1.Velocity = intialJunkVelocity

	j2 := new(Junk)
	otherJunkVelocity := Velocity{-testVelocityDx, testVelocityDy}
	j2.Velocity = otherJunkVelocity

	// Hit Junk with Other Junk
	j1.HitJunk(j2)

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
	j1.HitJunk(j2)
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
