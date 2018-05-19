package models

import (
	"fmt"
	"math"
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
		numUpdates        int
	}{
		{20, 20.02, 199, 5.03, 1},
		{MaxHoleRadius * 1.2, MaxHoleRadius * 1.2, 199, 5, 1},
		{20, 20.08, 196, 5.12, 4},
		{MaxHoleRadius * 1.2, MaxHoleRadius * 1.2, 195, 5, 5},
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
			for i := 0; i < tc.numUpdates; i++ {
				h.Update()
			}
			if h.Radius != tc.radiusWant {
				t.Errorf("got %g; want %g", h.Radius, tc.radiusWant)
			}
			if diff := h.GravityRadius - tc.gravityRadiusWant; math.Abs(diff) > 1e-9 {
				t.Errorf("got %g; want %g", h.GravityRadius, tc.gravityRadiusWant)
			}

			if h.Life != tc.lifeWant {
				t.Errorf("got %g; want %g", h.Life, tc.lifeWant)
			}
			// if !h.IsAlive {
			// 	t.Errorf("hole isAlive is false")
			// }
		})
	}

}

func TestHoleLifeCycle(t *testing.T) {
	testCases := []struct {
		life       float64
		numUpdates int
		wantIsDead bool // means that hole dies and starts a new life if false
	}{
		{MinHoleLife, 1, false},
		{MinHoleLife, MinHoleLife - 1, false},
		{MinHoleLife, MinHoleLife + 1, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test hole lifecycle with number of lives %v and number of updates %v", tc.life, tc.numUpdates), func(t *testing.T) {
			p := Position{X: 5, Y: 10}
			h := Hole{
				Position:      p,
				Radius:        20,
				Life:          tc.life,
				GravityRadius: 5,
				IsAlive:       true,
				IsCollidable:  false,
				StartingLife:  tc.life,
			}

			for i := 0; i < tc.numUpdates; i++ {
				h.Update()
			}
			if h.IsHoleDead() != tc.wantIsDead {
				t.Errorf("End of hole lifecycle is incorrectly reached or not reached")
			}

		})
	}
}
