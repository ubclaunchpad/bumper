package models

import (
	"math"
)

// PhysicsBody houses the physical properties of any object.
type PhysicsBody struct {
	Position    Position `json:"position"`
	Velocity    Velocity
	Mass        float64
	Restitution float64
}

// Vector represents a point in 2D space
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Position x y position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Velocity dx dy velocity
type Velocity struct {
	Dx float64 `json:"dx"`
	Dy float64 `json:"dy"`
}

func (v *Velocity) magnitude() float64 {
	return math.Hypot(v.Dx, v.Dy)
}

func (v *Velocity) normalize() {
	mag := v.magnitude()
	if mag > 0 {
		v.Dx /= mag
		v.Dy /= mag
	}
}

// CreateBody constructs a physics body at rest with
// given position, mass, and resitution factor.
func CreateBody(p Position, m float64, r float64) PhysicsBody {
	return PhysicsBody{
		Position:    p,
		Velocity:    Velocity{},
		Mass:        m,
		Restitution: r,
	}
}

// ApplyFactor applies given factor to this velocity
func (v *Velocity) ApplyFactor(factor float64) {
	v.Dx *= factor
	v.Dy *= factor
}

// ApplyVector applies given vector to this velocity
func (v *Velocity) ApplyVector(vector Velocity) {
	v.Dx += vector.Dx
	v.Dy += vector.Dy
}

// ApplyVelocity applies this body's velocity to its position
func (b *PhysicsBody) ApplyVelocity() {
	b.Position.X += b.Velocity.Dx
	b.Position.Y += b.Velocity.Dy
}

// func (p PhysicsBody) GetMass() float64 {
// 	return p.Mass
// }

// func (p PhysicsBody) GetPosition() Position {
// 	return p.Pos
// }
