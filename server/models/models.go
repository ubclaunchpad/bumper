package models

import (
	"math"
)

// PhysicsBody houses the physical properties of any object.
type PhysicsBody struct {
	Position    Position `json:"position"`
	Velocity    Velocity
	Radius      float64 `json:"radius"`
	Mass        float64
	Restitution float64
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

func (v Velocity) dot(v2 Velocity) float64 {
	return v.Dx*v2.Dx + v.Dy*v2.Dy
}

func (v Velocity) sub(v2 Velocity) Velocity {
	return Velocity{
		Dx: v.Dx - v2.Dx,
		Dy: v.Dy - v2.Dy,
	}
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
func CreateBody(p Position, rad float64, m float64, res float64) PhysicsBody {
	return PhysicsBody{
		Position:    p,
		Velocity:    Velocity{},
		Radius:      rad,
		Mass:        m,
		Restitution: res,
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

// VelocityMagnitude returns the magnitude of this body's velocity
func (b *PhysicsBody) VelocityMagnitude() float64 {
	return math.Hypot(b.Velocity.Dx, b.Velocity.Dy)
}

// NormalizeVelocity normalizes this body's velocity
func (b *PhysicsBody) NormalizeVelocity() {
	mag := b.VelocityMagnitude()
	if mag > 0 {
		b.Velocity.Dx /= mag
		b.Velocity.Dy /= mag
	}
}

// ApplyFactor applies given factor to the velocity of this body
func (b *PhysicsBody) ApplyFactor(factor float64) {
	b.Velocity.Dx *= factor
	b.Velocity.Dy *= factor
}

// ApplyXFactor applies given factor to the Dx velocity of this body
func (b *PhysicsBody) ApplyXFactor(factor float64) {
	b.Velocity.Dx *= factor
}

// ApplyYFactor applies given factor to the Dy velocity of this body
func (b *PhysicsBody) ApplyYFactor(factor float64) {
	b.Velocity.Dy *= factor
}

// ApplyVector applies given vector to the velocity of this body
func (b *PhysicsBody) ApplyVector(vector Velocity) {
	b.Velocity.Dx += vector.Dx
	b.Velocity.Dy += vector.Dy
}

// ApplyVelocity applies this body's velocity to its position
func (b *PhysicsBody) ApplyVelocity() {
	b.Position.X += b.Velocity.Dx
	b.Position.Y += b.Velocity.Dy
}

// GetPosition returns the body's position
func (b *PhysicsBody) GetPosition() Position {
	return b.Position
}

// GetX returns the body's X position
func (b *PhysicsBody) GetX() float64 {
	return b.Position.X
}

// GetY returns the body's Y position
func (b *PhysicsBody) GetY() float64 {
	return b.Position.Y
}

// GetVelocity returns the body's velocity
func (b *PhysicsBody) GetVelocity() Velocity {
	return b.Velocity
}

// GetDx sets the Dx of this body's velocity
func (b *PhysicsBody) GetDx() float64 {
	return b.Velocity.Dx
}

// GetDy returns the Dy of this body's velocity
func (b *PhysicsBody) GetDy() float64 {
	return b.Velocity.Dy
}

// GetMass returns the body's Mass
func (b *PhysicsBody) GetMass() float64 {
	return b.Mass
}

// GetRadius returns the body's Radius
func (b *PhysicsBody) GetRadius() float64 {
	return b.Radius
}

// GetRestitution returns the body's Restitution Factor
func (b *PhysicsBody) GetRestitution() float64 {
	return b.Restitution
}

// SetPosition sets the body's position
func (b *PhysicsBody) SetPosition(x float64, y float64) {
	b.Position.X = x
	b.Position.Y = y
}

// SetX sets the body's X position
func (b *PhysicsBody) SetX(x float64) {
	b.Position.X = x
}

// SetY sets the body's Y position
func (b *PhysicsBody) SetY(y float64) {
	b.Position.Y = y
}

// SetVelocity sets the body's velocity
func (b *PhysicsBody) SetVelocity(dX float64, dY float64) {
	b.Velocity.Dx = dX
	b.Velocity.Dy = dY
}

// SetDx sets the Dx of this body's velocity
func (b *PhysicsBody) SetDx(dX float64) {
	b.Velocity.Dx = dX
}

// SetDy sets the Dy of this body's velocity
func (b *PhysicsBody) SetDy(dY float64) {
	b.Velocity.Dy = dY
}

// SetMass sets the body's Mass
func (b *PhysicsBody) SetMass(mass float64) {
	b.Mass = mass
}

// SetRadius sets the body's Radius
func (b *PhysicsBody) SetRadius(radius float64) {
	b.Radius = radius
}

// SetRestitution sets the body's Restitution Factor
func (b *PhysicsBody) SetRestitution(restitution float64) {
	b.Restitution = restitution
}

// InelasticCollision update
func InelasticCollision(b1 *PhysicsBody, b2 *PhysicsBody) {
	// func InelasticCollision(b1 struct{ PhysicsBody }, b2 struct{ PhysicsBody }) {

	// Math: https://en.wikipedia.org/wiki/Elastic_collision
	x1Minusx2 := Velocity{
		Dx: b1.GetX() - b2.GetX(),
		Dy: b1.GetY() - b2.GetY(),
	}
	x2Minusx1 := Velocity{
		Dx: b2.GetX() - b1.GetX(),
		Dy: b2.GetY() - b1.GetY(),
	}
	// Maybe this function should be made critical as one block anyways?
	// b1.SetDx(b1.GetDx() * b1.GetRestitution() - (2*b2.GetMass()/(b1.GetMass()+b2.GetMass()))*(b1.)
	b1.Velocity.Dx = b1.Velocity.Dx*b1.Restitution - (2*b2.Mass/(b1.Mass+b2.Mass))*(b1.Velocity.sub(b2.Velocity).dot(x1Minusx2))/(math.Pow(x1Minusx2.magnitude(), 2))*(x1Minusx2.Dx)
	b1.Velocity.Dy = b1.Velocity.Dy*b1.Restitution - (2*b2.Mass/(b1.Mass+b2.Mass))*(b1.Velocity.sub(b2.Velocity).dot(x1Minusx2))/(math.Pow(x1Minusx2.magnitude(), 2))*(x1Minusx2.Dy)

	b2.Velocity.Dx = b2.Velocity.Dx*b2.Restitution - (2*b1.Mass/(b1.Mass+b2.Mass))*(b2.Velocity.sub(b1.Velocity).dot(x2Minusx1))/(math.Pow(x2Minusx1.magnitude(), 2))*(x2Minusx1.Dx)
	b2.Velocity.Dy = b2.Velocity.Dy*b2.Restitution - (2*b1.Mass/(b1.Mass+b2.Mass))*(b2.Velocity.sub(b1.Velocity).dot(x2Minusx1))/(math.Pow(x2Minusx1.magnitude(), 2))*(x2Minusx1.Dy)

}
