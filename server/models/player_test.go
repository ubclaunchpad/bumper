package models

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestAddPoints(t *testing.T) {
	p := new(Player)

	p.AddPoints(10)
	if p.Points != 10 {
		t.Error("Error adding points")
	}
}

func TestCreatePlayer(t *testing.T) {
	testPosition := Position{5, 5}
	testName := "testy"
	testColor := "blue"
	ws := new(*websocket.Conn)

	//Test initialization of player
	p := CreatePlayer(testName, testPosition, testColor, *ws)

	//Test name assignment of player
	if p.Name != testName {
		t.Error("Error assigning name")
	}
}

func TestUpdatePosition(t *testing.T) {
	//Mock player and info
	testHeight1 := 10.0
	testWidth1 := 20.0
	testAngle1 := 0.0
	testPosition1 := Position{5, 5}
	p := new(Player)
	p.Position = testPosition1
	p.Angle = testAngle1
	p.Controls.Left = false
	p.Controls.Right = false
	p.Controls.Up = false

	//Test left control
	p.Controls.Left = true
	p.UpdatePosition(testHeight1, testWidth1)
	if p.Angle != (testAngle1 + 0.1) {
		t.Error("Error in left key player control")
	}
	p.Controls.Left = false
	p.Controls.Right = true
	p.UpdatePosition(testHeight1, testWidth1)
	if p.Angle != (testAngle1) {
		t.Error("Error in right key player control or symmetry")
	}
	p.Controls.Right = false
	p.Controls.Up = true
	p.UpdatePosition(testHeight1, testWidth1)
	if p.Angle != (testAngle1) {
		t.Error("Error in up key player control")
	}
	p.UpdatePosition(testHeight1, testWidth1)

}

func TestHitjunk(t *testing.T) {
	p := new(Player)
	testVelocity := Velocity{4, 4}
	p.Velocity = testVelocity
	p.hitJunk()
	if p.Velocity.Dx != testVelocity.Dx*JunkBounceFactor {
		t.Error("Error calculating player Dx hitting junk")
	}
	if p.Velocity.Dy != testVelocity.Dy*JunkBounceFactor {
		t.Error("Error calculating player Dy hitting junk")
	}
}

func TestKeyDownHandler(t *testing.T) {
	p := new(Player)
	key := UpKey
	p.KeyDownHandler(key)
	if p.Controls.Up != true {
		t.Error("Error key-down-handling up key control")
	}
	key = RightKey
	p.KeyDownHandler(key)
	if p.Controls.Right != true {
		t.Error("Error key-down-handling right key control")
	}
	key = DownKey
	p.KeyDownHandler(key)
	if p.Controls.Down != true {
		t.Error("Error key-down-handling down key control")
	}
	key = LeftKey
	p.KeyDownHandler(key)
	if p.Controls.Left != true {
		t.Error("Error key-down-handling left key control")
	}
}

func TestKeyUpHandler(t *testing.T) {
	p := new(Player)
	key := UpKey
	p.KeyUpHandler(key)
	if p.Controls.Up != false {
		t.Error("Error key-up-handling up key control")
	}
	key = RightKey
	p.KeyUpHandler(key)
	if p.Controls.Right != false {
		t.Error("Error key-up-handling right key control")
	}
	key = DownKey
	p.KeyUpHandler(key)
	if p.Controls.Down != false {
		t.Error("Error key-up-handling down key control")
	}
	key = LeftKey
	p.KeyUpHandler(key)
	if p.Controls.Left != false {
		t.Error("Error key-up-handling left key control")
	}
}
