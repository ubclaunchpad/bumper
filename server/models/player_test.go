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
	p := CreatePlayer(testName, testPosition, testColor, *ws)

	//Name Assignment
	if p.Name != testName {
		t.Error("Error assigning name")
	}

	//To do: test other initializations
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
