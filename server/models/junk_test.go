package models

import (
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
