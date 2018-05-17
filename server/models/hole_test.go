package models

import (
	"testing"
)

// func TestAddPoints(t *testing.T) {
// 	p := new(Player)

// 	p.AddPoints(10)
// 	if p.Points != 10 {
// 		t.Error("Error adding points")
// 	}
// }

func TestCreateHole(t *testing.T) {
	p := Position{
		X: 5,
		Y: 10,
	}
	h := CreateHole(p)
	if h.IsAlive {
		t.Error("isAlive is incorrectly set")
	}
	if h.Position.X != 5 {
		t.Error("X position is not set correctly")
	}
	if h.Position.Y != 10 {
		t.Error("Y position is not set correctly")
	}
	if h.GravityRadius != h.Radius*gravityRadiusFactor {
		t.Error("Gravity radius is calculated incorrectly")
	}
}

func TestUpdateHole(t *testing.T) {
	p := Position{
		X: 5,
		Y: 10,
	}
	h := CreateHole(p)
	h.StartingLife = 200
	h.Life = 200
	h.Radius = 20
	h.GravityRadius = 5

	h.Update()
	if h.Life != 199 {
		t.Error("Life is incorrectly updated")
	}

	if h.Radius != 20.02 {
		t.Error("Radius is incorrectly updated")
	}
	if h.GravityRadius != 5.03 {
		t.Error("Radius is incorrectly updated")
	}

}

func TestUpdateMaxSizeHole(t *testing.T) {
	p := Position{
		X: 5,
		Y: 10,
	}
	h := CreateHole(p)
	h.Radius = MaxHoleRadius * 1.2
	h.GravityRadius = 5
	h.Update()
	if h.Radius != (MaxHoleRadius * 1.2) {
		t.Error("Radius increased over the max size")
	}
	if h.GravityRadius != 5 {
		t.Error("Radius is incorrectly increased")
	}
}

func TestStartNewLife(t *testing.T) {

}
