package arena

import (
	"fmt"
	"testing"
)

const (
	testHeight    = 2400
	testWidth     = 2800
	testHoleCount = 20
	testJunkCount = 30
)

// sets up a fresh instance of an arena for testing
func createTestArena() *Arena {
	return CreateArena(testHeight, testWidth, testHoleCount, testJunkCount)
}

func TestCreateArena(t *testing.T) {
	a := createTestArena()

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
	a := createTestArena()

	numPlayers := 3
	for i := 0; i < numPlayers; i++ {
		t.Run("TestAddPlayer", func(t *testing.T) {

			err := a.AddPlayer(fmt.Sprintf("player%d", i), nil)
			if err != nil {
				t.Errorf("Failed to add player: %v", err)
			}

			if len(a.Players) != i+1 {
				t.Errorf("Player map error")
			}
		})
	}
}
