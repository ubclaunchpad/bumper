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
)

var (
	centerPosPlayerTest = Position{testWidthPlayerTest / 2, testHeightPlayerTest / 2}
)

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
	p := CreatePlayer(testNamePlayerTest, centerPosPlayerTest, testColorPlayerTest, *ws)

	//Test name assignment of player
	if p.Name != testNamePlayerTest {
		t.Error("Error assigning name")
	}
}

func TestUpdatePosition(t *testing.T) {
	//Mock player and info
	p := new(Player)
	p.Angle = 0
	p.Position = centerPosPlayerTest
	rangeAngle := 2
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
			p.Velocity = tc.playerVelocity
			p.Position = tc.playerPosition
			p.Angle = tc.playerAngle

			p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			//Test max velocity
			if p.Velocity.magnitude() > MaxVelocity {
				t.Error("Error calculating max velocity")
			}
			//Test directional controls
			//Left
			p.Controls.Left = true
			for i := 0; i < rangeAngle; i++ {
				prevAngle := p.Angle
				p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
				if (p.Angle - prevAngle) != 0.1 {
					t.Error("Error calculating left control")
				}
			}
			p.Controls.Left = false
			p.Controls.Right = true
			//Right
			for i := 0; i < rangeAngle; i++ {
				prevAngle := p.Angle
				p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
				if (p.Angle - prevAngle) != -0.1 {
					t.Error("Error calculating right control")
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
			prevMagnitude := p.Velocity.magnitude()
			p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			if p.Velocity.magnitude() > prevMagnitude {
				t.Error("Error calculating friction")
			}
		})

		//Test friction and accelerate
		t.Run(tc.description, func(t *testing.T) {
			p := new(Player)
			p.Velocity = tc.playerVelocity
			p.Position = tc.playerPosition
			p.Angle = tc.playerAngle
			//Test Friction
			prevMagnitude := p.Velocity.magnitude()
			p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			if p.Velocity.magnitude() > prevMagnitude {
				t.Error("Error calculating friction")
			}

			// p.Controls.Up = true
			// //Test Accelerate
			// prevMagnitude = p.Velocity.magnitude()
			// p.UpdatePosition(testHeightPlayerTest, testWidthPlayerTest)
			// if p.Velocity.magnitude() < prevMagnitude {
			// 	t.Error("Error calculating acceleration")
			// }
		})
	}
}
func TestHitjunk(t *testing.T) {
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
			p := new(Player)
			p.Velocity = tc.playerVelocity
			p.hitJunk()
			if p.Velocity.Dx != tc.playerVelocity.Dx*JunkBounceFactor {
				t.Error("Error calculating player Dx hitting junk")
			}
			if p.Velocity.Dy != tc.playerVelocity.Dy*JunkBounceFactor {
				t.Error("Error calculating player Dy hitting junk")
			}
			if p.Velocity.Dx < tc.playerVelocity.Dx*JunkBounceFactor {
				t.Error("Error x axis bounce factor greater than 1")
			}
			if p.Velocity.Dy < tc.playerVelocity.Dy*JunkBounceFactor {
				t.Error("Error y axis bounce factor greater than 1")
			}
		})
	}
}

func TestKeyDownHandler(t *testing.T) {
	p := new(Player)
	key := UpKey
	p.KeyDownHandler(key)
	if !p.Controls.Up {
		t.Error("Error key-down-handling up key control")
	}
	key = RightKey
	p.KeyDownHandler(key)
	if !p.Controls.Right {
		t.Error("Error key-down-handling right key control")
	}
	key = DownKey
	p.KeyDownHandler(key)
	if !p.Controls.Down {
		t.Error("Error key-down-handling down key control")
	}
	key = LeftKey
	p.KeyDownHandler(key)
	if !p.Controls.Left {
		t.Error("Error key-down-handling left key control")
	}
}

func TestKeyUpHandler(t *testing.T) {
	p := new(Player)
	key := UpKey
	p.KeyUpHandler(key)
	if p.Controls.Up {
		t.Error("Error key-up-handling up key control")
	}
	key = RightKey
	p.KeyUpHandler(key)
	if p.Controls.Right {
		t.Error("Error key-up-handling right key control")
	}
	key = DownKey
	p.KeyUpHandler(key)
	if p.Controls.Down {
		t.Error("Error key-up-handling down key control")
	}
	key = LeftKey
	p.KeyUpHandler(key)
	if p.Controls.Left {
		t.Error("Error key-up-handling left key control")
	}
}
