package models

import (
	"math"
	"testing"

	"github.com/gorilla/websocket"
)

const (
	testHeightPlayerTest = 400
	testWidthPlayerTest  = 800
	testNamePlayerTest   = "testy"
	testColorPlayerTest  = "blue"
	rangeAngle           = 360
	angleIncrement       = 0.1
	roundingError5SigFig = 0.00001
)

var (
	testID              = "testid"
	centerPosPlayerTest = Position{testWidthPlayerTest / 2, testHeightPlayerTest / 2}
)

func isWithinTolerance(test float64, target float64, tolerance float64) bool {
	if (test < (target + tolerance)) && (test > (target - tolerance)) {
		return true
	}
	return false
}

func TestAddPoints(t *testing.T) {
	p := new(Player)
	p.AddPoints(10)
	if p.Points != 10 {
		t.Error("Error adding points")
	}
}

func TestCreatePlayer(t *testing.T) {
	ws := new(*websocket.Conn)

	//Test initialization of player
	p := CreatePlayer(testID, testNamePlayerTest, centerPosPlayerTest, testColorPlayerTest, *ws)

	//Test name assignment of player
	if p.Name != testNamePlayerTest {
		t.Error("Error assigning name")
	}
}

func TestUpdatePosition(t *testing.T) {
	//Mock player and info
	p := new(Player)
	p.Angle = 0
	p.SetPosition(centerPosPlayerTest)
	testCases := []struct {
		description    string
		playerVelocity Velocity
		playerPosition Position
		playerAngle    float64
	}{
		{"Not moving", Velocity{0, 0}, centerPosPlayerTest, 0},
		{"Max velocity", Velocity{math.Sqrt(MaxVelocity), math.Sqrt(MaxVelocity)}, centerPosPlayerTest, 0},
		{"Moving N", Velocity{5, 0}, centerPosPlayerTest, 0},
		{"Moving E", Velocity{0, 5}, centerPosPlayerTest, 0},
		{"Moving S", Velocity{-5, 0}, centerPosPlayerTest, 0},
		{"Moving W", Velocity{0, -5}, centerPosPlayerTest, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			p := new(Player)
			p.SetVelocity(tc.playerVelocity)
			p.SetPosition(tc.playerPosition)
			p.Angle = tc.playerAngle

			p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			//Test max velocity
			if p.VelocityMagnitude() > MaxVelocity {
				t.Error("Error calculating max velocity")
			}
			//Test directional controls
			//Left
			p.Controls.Left = true
			for i := 0; i < rangeAngle; i++ {
				prevAngle := p.Angle
				p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
				angleDifference := p.Angle - prevAngle
				if !isWithinTolerance(math.Abs(angleDifference), angleIncrement, roundingError5SigFig) {
					t.Error("Error calculating left control, angle increment is", angleDifference)
				}
			}
			p.Controls.Left = false
			p.Controls.Right = true
			//Right
			for i := 0; i < rangeAngle; i++ {
				prevAngle := p.Angle
				p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
				angleDifference := p.Angle - prevAngle
				if !isWithinTolerance(math.Abs(angleDifference), angleIncrement, roundingError5SigFig) {
					t.Error("Error calculating right control, angle increment is", angleDifference)
				}
			}
			p.Controls.Left = true
			//Both
			for i := 0; i < rangeAngle; i++ {
				prevAngle := p.Angle
				p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
				if p.Angle != prevAngle {
					t.Error("Error calculating both controls")
				}
			}
			p.Controls.Left = false
			p.Controls.Right = false
			//Up
			//Friction
			prevMagnitude := p.VelocityMagnitude()
			p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			if p.VelocityMagnitude() > prevMagnitude {
				t.Error("Error calculating friction")
			}
		})

		//Test friction and accelerate
		t.Run(tc.description, func(t *testing.T) {
			p := new(Player)
			p.SetVelocity(tc.playerVelocity)
			p.SetPosition(tc.playerPosition)
			p.Angle = tc.playerAngle
			//Test Friction
			prevMagnitude := p.VelocityMagnitude()
			p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			if p.VelocityMagnitude() > prevMagnitude {
				t.Error("Error calculating friction", p.VelocityMagnitude(), "expected to be less than", prevMagnitude)
			}

			// p.Controls.Up = true
			// //Test Accelerate
			// prevMagnitude = p.Velocity.magnitude()
			// p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			// if p.Velocity.magnitude() < prevMagnitude {
			// 	t.Error("Error calculating acceleration", p.Velocity.magnitude(), "expected to be greater than", prevMagnitude)
			// }
		})
	}
}
func TestHitJunk(t *testing.T) {
	testCases := []struct {
		description    string
		playerVelocity Velocity
	}{
		{"Moving NW", Velocity{5, -5}},
		{"Moving NE", Velocity{5, 5}},
		{"Moving SW", Velocity{-5, -5}},
		{"Moving SE", Velocity{-5, 5}},
		{"Stationary", Velocity{0, 0}},
		{"Moving N", Velocity{5, 0}},
		{"Moving E", Velocity{0, 5}},
		{"Moving S", Velocity{-5, 0}},
		{"Moving W", Velocity{0, -5}},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			p := CreatePlayer(testID, testNamePlayerTest, centerPos, testColorPlayerTest, *new(*websocket.Conn))
			p.SetVelocity(tc.playerVelocity)
			j := CreateJunk(centerPosBelowRight)
			initialJunkVelocity := Velocity{0, 0}
			j.SetVelocity(initialJunkVelocity)

			j.HitBy(p)

			if p.GetDx() == tc.playerVelocity.Dx {
				t.Error("Error calculating player Dx hitting junk", p.GetDx(), tc.playerVelocity.Dx)
			}
			if p.GetDy() == tc.playerVelocity.Dy {
				t.Error("Error calculating player Dy hitting junk")
			}
		})
	}
}

func TestKeyHandler(t *testing.T) {
	p := new(Player)
	testCases := []struct {
		description  string
		key          int
		expectedKeys KeysPressed
		fp           func(int)
	}{
		{"Up Key KeyDownHandler", UpKey, KeysPressed{Up: true, Right: false, Left: false, Down: false}, p.KeyDownHandler},
		{"Right Key KeyDownHandler", RightKey, KeysPressed{Up: true, Right: true, Left: false, Down: false}, p.KeyDownHandler},
		{"Left Key KeyDownHandler", LeftKey, KeysPressed{Up: true, Right: true, Left: true, Down: false}, p.KeyDownHandler},
		{"Down Key KeyDownHandler", DownKey, KeysPressed{Up: true, Right: true, Left: true, Down: true}, p.KeyDownHandler},
		{"Up Key KeyUpHandler", UpKey, KeysPressed{Up: false, Right: true, Left: true, Down: true}, p.KeyUpHandler},
		{"Right Key KeyUpHandler", RightKey, KeysPressed{Up: false, Right: false, Left: true, Down: true}, p.KeyUpHandler},
		{"Left Key KeyUpHandler", LeftKey, KeysPressed{Up: false, Right: false, Left: false, Down: true}, p.KeyUpHandler},
		{"Down Key KeyUpHandler", DownKey, KeysPressed{Up: false, Right: false, Left: false, Down: false}, p.KeyUpHandler},
	}

	//test keydownhandler
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.fp(tc.key)
			if p.Controls != tc.expectedKeys {
				t.Errorf("Error: %s. Keys Pressed: %v", tc.description, p.Controls)
			}
		})
	}
}

func TestPlayerBumpPlayer(t *testing.T) {
	testCases := []struct {
		description    string
		playerVelocity Velocity
	}{
		{"Moving NW", Velocity{5, -5}},
		{"Moving NE", Velocity{5, 5}},
		{"Moving SW", Velocity{-5, -5}},
		{"Moving SE", Velocity{-5, 5}},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			p1 := CreatePlayer("id1", "player1", centerPos, testColorPlayerTest, nil)
			p2 := CreatePlayer("id2", "player2", centerPosBelowRight, testColorPlayerTest, nil)

			p1.SetVelocity(tc.playerVelocity)

			p1.HitPlayer(p2)

			if checkDirection(p1.GetVelocity(), tc.playerVelocity) || p1.GetDx() == tc.playerVelocity.Dx ||
				p1.GetDy() == tc.playerVelocity.Dy {
				t.Error("Player 1's velocity was unaffected")
			}

			if p2.GetDx() == 0 || p2.GetDy() == 0 {
				t.Error("Player 2's velocity was unaffected")
			}
		})
	}
}

func TestPlayerKillPlayer(t *testing.T) {
	testCases := []struct {
		description       string
		bumperPosition    Position
		bumpeePosition    Position
		bumperVelocity    Velocity
		holePosition      Position
		shouldAwardPoints bool
	}{
		{"Award Points", Position{40, 50}, Position{45, 50}, Velocity{5, 0}, Position{50, 50}, true},
		{"No Points", Position{5, 50}, Position{10, 50}, Velocity{5, 0}, Position{400, 50}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			p1 := CreatePlayer("id1", "player1", tc.bumperPosition, testColorPlayerTest, nil)
			p2 := CreatePlayer("id2", "player2", tc.bumpeePosition, testColorPlayerTest, nil)
			p1.SetVelocity(tc.bumperVelocity)
			h := CreateHole(tc.holePosition)

			p1.HitPlayer(p2)

			// Advance the game 200 time ticks
			for i := 0; i < 200; i++ {
				p1.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
				p2.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)

				if areCirclesColliding(p2.GetPosition(), p2.GetRadius(), h.GetPosition(), h.GetRadius()) {
					playerScored := p2.LastPlayerHit
					if playerScored != nil {
						playerScored.AddPoints(PointsPerPlayer)
					}
				}
			}

			// If this flag is set the hole was close enough to eat the bumpee within the debounce time
			// if not set the hole was too far away
			if tc.shouldAwardPoints {
				if p1.Points == 0 {
					t.Error("Bumpee was not killed correctly")
				}
			} else {
				if p1.Points != 0 {
					t.Error("Bumpee was killed incorrectly")
				}
			}
		})
	}
}

// Helper function, detect collision between objects
func areCirclesColliding(p Position, r1 float64, q Position, r2 float64) bool {
	return math.Pow(p.X-q.X, 2)+math.Pow(p.Y-q.Y, 2) <= math.Pow(r1+r2, 2)
}
