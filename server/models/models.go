package models

import (
	"math"
	"sync"
)

// PhysicsBody houses the physical properties of any object.
type PhysicsBody struct {
	Position    Position
	Velocity    Velocity
	Radius      float64
	Mass        float64
	Restitution float64
	rwMutex     *sync.RWMutex
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
func CreateBody(p Position, rad float64, m float64, res float64, lock *sync.RWMutex) PhysicsBody {
	return PhysicsBody{
		Position:    p,
		Velocity:    Velocity{},
		Radius:      rad,
		Mass:        m,
		Restitution: res,
		rwMutex:     lock,
	}
}

// ApplyFactor applies given factor to this velocity
func (v *Velocity) ApplyFactor(factor float64) {
	v.Dx *= factor
	v.Dy *= factor
}

// VelocityMagnitude returns the magnitude of this body's velocity
func (b *PhysicsBody) VelocityMagnitude() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return math.Hypot(b.Velocity.Dx, b.Velocity.Dy)
}

// NormalizeVelocity normalizes this body's velocity
func (b *PhysicsBody) NormalizeVelocity() {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	mag := b.VelocityMagnitude()
	if mag > 0 {
		b.Velocity.Dx /= mag
		b.Velocity.Dy /= mag
	}
}

// ApplyFactor applies given factor to the velocity of this body
func (b *PhysicsBody) ApplyFactor(factor float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dx *= factor
	b.Velocity.Dy *= factor
}

// ApplyXFactor applies given factor to the Dx velocity of this body
func (b *PhysicsBody) ApplyXFactor(factor float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dx *= factor
}

// ApplyYFactor applies given factor to the Dy velocity of this body
func (b *PhysicsBody) ApplyYFactor(factor float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dy *= factor
}

// ApplyVector applies given vector to the velocity of this body
func (b *PhysicsBody) ApplyVector(vector Velocity) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dx += vector.Dx
	b.Velocity.Dy += vector.Dy
}

// ApplyVelocity applies this body's velocity
func (b *PhysicsBody) ApplyVelocity() {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Position.X += b.Velocity.Dx
	b.Position.Y += b.Velocity.Dy
}

// GetPosition returns this object's position
func (b *PhysicsBody) GetPosition() Position {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Position
}

// GetX return's this object's X position
func (b *PhysicsBody) GetX() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Position.X
}

// GetY return's this object's Y position
func (b *PhysicsBody) GetY() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Position.Y
}

// GetVelocity return's this object's velocity
func (b *PhysicsBody) GetVelocity() Velocity {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Velocity
}

// GetDx return's this object's velocity's Dx component
func (b *PhysicsBody) GetDx() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Velocity.Dx
}

// GetDy return's this object's velocity's Dy component
func (b *PhysicsBody) GetDy() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Velocity.Dy
}

// GetMass return's this object's mass
func (b *PhysicsBody) GetMass() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Mass
}

// GetRadius return's this object's radius
func (b *PhysicsBody) GetRadius() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Radius
}

// GetRestitution return's this object's restitution factor
func (b *PhysicsBody) GetRestitution() float64 {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	return b.Restitution
}

// SetPosition set's this body's position
func (b *PhysicsBody) SetPosition(p Position) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Position.X = p.X
	b.Position.Y = p.Y
}

// SetX sets the body's X position
func (b *PhysicsBody) SetX(x float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Position.X = x
}

// SetY sets the body's Y position
func (b *PhysicsBody) SetY(y float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Position.Y = y
}

// SetVelocity sets the body's velocity
func (b *PhysicsBody) SetVelocity(v Velocity) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dx = v.Dx
	b.Velocity.Dy = v.Dy
}

// SetDx sets the Dx of this body's velocity
func (b *PhysicsBody) SetDx(dX float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dx = dX
}

// SetDy sets the Dy of this body's velocity
func (b *PhysicsBody) SetDy(dY float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Velocity.Dy = dY
}

// SetMass sets the body's Mass
func (b *PhysicsBody) SetMass(mass float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Mass = mass
}

// SetRadius sets the body's Radius
func (b *PhysicsBody) SetRadius(radius float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Radius = radius
}

// SetRestitution sets the body's Restitution Factor
func (b *PhysicsBody) SetRestitution(restitution float64) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.Restitution = restitution
}

// InelasticCollision update
func InelasticCollision(b1 *PhysicsBody, b2 *PhysicsBody) {
	// XXX - This feels like it could lead to deadlocks:
	b1.rwMutex.Lock()
	defer b1.rwMutex.Unlock()
	b2.rwMutex.Lock()
	defer b2.rwMutex.Unlock()

	// Math: https://en.wikipedia.org/wiki/Elastic_collision
	x1Minusx2 := Velocity{
		Dx: b1.Position.X - b2.Position.X,
		Dy: b1.Position.Y - b2.Position.Y,
	}
	x2Minusx1 := Velocity{
		Dx: b2.Position.X - b1.Position.X,
		Dy: b2.Position.Y - b1.Position.Y,
	}

	b1.Velocity.Dx = b1.Velocity.Dx*b1.Restitution - (2*b2.Mass/(b1.Mass+b2.Mass))*(b1.Velocity.sub(b2.Velocity).dot(x1Minusx2))/(math.Pow(x1Minusx2.magnitude(), 2))*(x1Minusx2.Dx)
	b1.Velocity.Dy = b1.Velocity.Dy*b1.Restitution - (2*b2.Mass/(b1.Mass+b2.Mass))*(b1.Velocity.sub(b2.Velocity).dot(x1Minusx2))/(math.Pow(x1Minusx2.magnitude(), 2))*(x1Minusx2.Dy)

	b2.Velocity.Dx = b2.Velocity.Dx*b2.Restitution - (2*b1.Mass/(b1.Mass+b2.Mass))*(b2.Velocity.sub(b1.Velocity).dot(x2Minusx1))/(math.Pow(x2Minusx1.magnitude(), 2))*(x2Minusx1.Dx)
	b2.Velocity.Dy = b2.Velocity.Dy*b2.Restitution - (2*b1.Mass/(b1.Mass+b2.Mass))*(b2.Velocity.sub(b1.Velocity).dot(x2Minusx1))/(math.Pow(x2Minusx1.magnitude(), 2))*(x2Minusx1.Dy)

}
