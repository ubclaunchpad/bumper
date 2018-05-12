package models

import "testing"

func TestAddPoints(t *testing.T) {
	p := new(Player)

	p.AddPoints(10)
	if p.Points != 10 {
		t.Error("Error adding points")
	}
}
