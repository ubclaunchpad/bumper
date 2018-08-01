package models

import (
	"fmt"
	"math"
)

// PhysicsBody houses the physical properties of any object.
type PhysicsBody struct {
	Position    Position `json:"position"`
	Velocity    Velocity
	Radius      float64
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

// ApplyVelocity applies this body's velocity to its position
func (b *PhysicsBody) ApplyVelocity() {
	b.Position.X += b.Velocity.Dx
	b.Position.Y += b.Velocity.Dy
}

// InelasticCollision update
func InelasticCollision(b1 *PhysicsBody, b2 *PhysicsBody) {

	// J is 2 P is 1
	d := math.Sqrt(math.Pow((b1.Position.X-b2.Position.X), 2) + math.Pow((b1.Position.Y-b2.Position.Y), 2))
	fmt.Println("d", d)
	SinYx := (b1.Position.Y - b2.Position.Y) / (b1.Position.X - b2.Position.X)
	// Yx := math.Asin((b1.Position.Y - b2.Position.Y) / (b1.Position.X - b2.Position.X))
	if SinYx > 1 || SinYx < -1 {
		fmt.Println("x1", b2.Position.X, "y1", b2.Position.Y)
		fmt.Println("x2", b1.Position.X, "y2", b1.Position.Y)
		InvSinYx := (b2.Position.Y - b1.Position.Y) / (b2.Position.X - b1.Position.X)
		fmt.Println("SinYx", SinYx)
		fmt.Println("InvSinYx", InvSinYx)
		SinYx = 1 / SinYx
		// if SinYx < -1 {
		// 	SinYx = -1
		// }
		// if SinYx > 1 {
		// 	SinYx = 1
		// }
	}
	Yx := math.Asin(SinYx)
	r1plusr2 := b1.Radius + b2.Radius
	fmt.Println("Yx", Yx)
	fmt.Println("SinYx", SinYx)
	fmt.Println("r1 + r2", r1plusr2)

	Yv := math.Atan((b2.Velocity.Dy - b1.Velocity.Dy) / (b2.Velocity.Dx - b1.Velocity.Dx))
	fmt.Println("Yv", Yv)
	alpha := math.Asin((d * math.Sin(Yx-Yv)) / (b2.Radius + b1.Radius))
	fmt.Println("Alpha", alpha)
	a := math.Tan(Yv + alpha)
	fmt.Println("a", a)
	massRatio := float64(b1.Mass / b2.Mass)
	fmt.Println("massRatio", massRatio)

	Deltab1Vx := 2 * (b2.Velocity.Dx - b1.Velocity.Dx + a*(b2.Velocity.Dy-b1.Velocity.Dy)) / ((math.Pow(a, 2) + 1) * (massRatio + 1))
	fmt.Println("Delta", Deltab1Vx)

	b1.Velocity.Dx += Deltab1Vx * b1.Restitution
	b1.Velocity.Dy += a * Deltab1Vx * b1.Restitution
	b2.Velocity.Dx -= massRatio * Deltab1Vx * b2.Restitution
	b2.Velocity.Dy -= a * massRatio * Deltab1Vx * b2.Restitution
}
