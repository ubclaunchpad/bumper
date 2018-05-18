package models

import (
	"fmt"
	"testing"
)

func TestCreateHole(t *testing.T) {
	p := Position{
		X: 5,
		Y: 10,
	}
	h := CreateHole(p)
	if h.Life < MinHoleLife || h.Life > MaxHoleLife {
		t.Error("hole life span is created too large or too small")
	}
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
	testCases := []struct {
		radius            float64
		radiusWant        float64
		lifeWant          float64
		gravityRadiusWant float64
	}{
		{20, 20.02, 199, 5.03},
		{MaxHoleRadius * 1.2, MaxHoleRadius * 1.2, 199, 5},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test updateHole with radius %v", tc.radius), func(t *testing.T) {
			h := Hole{
				Position:      Position{X: 5, Y: 10},
				Radius:        tc.radius,
				Life:          200,
				GravityRadius: 5,
				IsAlive:       false,
				StartingLife:  200,
			}
			h.Update()
			if h.Radius != tc.radiusWant {
				t.Errorf("got %g; want %g", h.Radius, tc.radiusWant)
			}
			if h.GravityRadius != tc.gravityRadiusWant {
				t.Errorf("got %g; want %g", h.GravityRadius, tc.gravityRadiusWant)
			}
			if h.Life != tc.lifeWant {
				t.Errorf("got %g; want %g", h.Life, tc.lifeWant)
			}
		})
	}

}
