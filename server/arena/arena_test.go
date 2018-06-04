package arena

import (
	"fmt"
	"testing"

	"github.com/ubclaunchpad/bumper/server/models"
)

const (
	testHeight    = 2400
	testWidth     = 2800
	testHoleCount = 20
	testJunkCount = 30
)

var (
	testVelocity    = models.Velocity{Dx: 1, Dy: 1}
	centerPosition  = models.Position{X: testWidth / 2, Y: testHeight / 2}
	quarterPosition = models.Position{X: testWidth / 4, Y: testHeight / 4}
)

func CreateArenaWithPlayer(p models.Position) *Arena {
	a := CreateArena(testHeight, testWidth, 0, 0)
	a.AddPlayer("test", nil)
	testPlayer := a.Players["test"]
	testPlayer.Position = p
	testPlayer.Velocity = testVelocity
	return a
}

func TestCreateArena(t *testing.T) {
	a := CreateArena(testHeight, testWidth, testHoleCount, testJunkCount)

	holeCount := len(a.Holes)
	if holeCount != testHoleCount {
		t.Errorf("Arena did not spawn enough holes. Got %d/%d Holes", holeCount, testHoleCount)
	}

	junkCount := len(a.Junk)
	if junkCount != testJunkCount {
		t.Errorf("Arena did not spawn enough junk. Got %d/%d Junk", junkCount, testJunkCount)
	}
}

func TestAddPlayer(t *testing.T) {
	a := CreateArena(testHeight, testWidth, testHoleCount, testJunkCount)

	numPlayers := 3
	for i := 0; i < numPlayers; i++ {
		t.Run("TestAddPlayer", func(t *testing.T) {

			name := fmt.Sprintf("player%d", i)
			err := a.AddPlayer(name, nil)
			if err != nil {
				t.Errorf("Failed to add player: %v", err)
			}

			if len(a.Players) != i+1 {
				t.Errorf("Player map error")
			}
		})
	}
}

func TestAddRemoveObject(t *testing.T) {
	a := CreateArena(testHeight, testWidth, 0, 0)

	testCount := 10
	testCases := []struct {
		description string
		remove      func(int) bool
		add         func()
	}{
		{"Holes", a.removeHole, a.addHole},
		{"Junk", a.removeJunk, a.addJunk},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			for i := 0; i < testCount; i++ {
				tc.add()
			}

			for i := 0; i < testCount; i++ {
				ok := tc.remove(0)
				if !ok {
					t.Errorf("%s removal error at count %d", tc.description, i)
				}
			}

			ok := tc.remove(0)
			if ok {
				t.Errorf("Removal from empty slice of %s returned true", tc.description)
			}
		})
	}
}

func TestPlayerToPlayerCollisions(t *testing.T) {
	testCases := []struct {
		otherPlayer      string
		testPosition     models.Position
		expectedPosition models.Position
	}{
		{"colliding", quarterPosition, models.Position{X: 700.375, Y: 600.375}},
		{"non-colliding", centerPosition, quarterPosition},
	}

	for _, tc := range testCases {
		t.Run("Player to Player collision", func(t *testing.T) {
			a := CreateArenaWithPlayer(quarterPosition)

			a.AddPlayer(tc.otherPlayer, nil)
			a.Players[tc.otherPlayer].Position = tc.testPosition

			a.playerCollisions()
			if a.Players["test"].Position != tc.expectedPosition {
				t.Errorf("%s detection failed. Got player at %v. Expected player at %v", tc.otherPlayer, a.Players["test"].Position, tc.expectedPosition)
			}
		})
	}
}

func TestPlayerToJunkCollisions(t *testing.T) {
	a := CreateArenaWithPlayer(quarterPosition)

	testCases := []struct {
		description    string
		testPosition   models.Position
		expectedPlayer *models.Player
	}{
		{"non-colliding", centerPosition, nil},
		{"colliding", quarterPosition, a.Players["test"]},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			a.addJunk()
			a.Junk[i].Position = tc.testPosition

			a.playerCollisions()
			if a.Junk[i].LastPlayerHit != tc.expectedPlayer {
				t.Errorf("%s detection failed. Test Player at %v. Junk at %v. Junk Last Player Hit %v", tc.description, a.Players["test"].Position, a.Junk[i].Position, a.Junk[i].LastPlayerHit)
			}
		})
	}
}

func TestJunkToJunkCollisions(t *testing.T) {
	testCases := []struct {
		description      string
		testPosition     models.Position
		expectedVelocity models.Velocity
	}{
		{"non-colliding", centerPosition, testVelocity},
		{"colliding", quarterPosition, models.Velocity{Dx: -0.5, Dy: -0.5}},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			a := CreateArena(testHeight, testWidth, 0, 0)
			a.addJunk()
			a.Junk[0].Position = quarterPosition
			a.Junk[0].Velocity = testVelocity

			a.addJunk()
			a.Junk[1].Position = tc.testPosition

			a.junkCollisions()
			if a.Junk[0].Velocity != tc.expectedVelocity {
				t.Errorf("%s detection failed. Expected %v. Got %v", tc.description, tc.expectedVelocity, a.Junk[0].Velocity)
			}
		})
	}
}

// TODO: Complete once Game package refactoring has happened
func TestHoleToPlayerCollisions(t *testing.T) {

}

// TODO: Complete once Game package refactoring has happened
func TestHoleToJunkCollisions(t *testing.T) {

}
