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

func (a Velocity) dot(b Velocity) float64 {
	return a.Dx*b.Dx + a.Dy*b.Dy
}

func (a Velocity) sub(b Velocity) Velocity {
	return Velocity{
		Dx: a.Dx - b.Dx,
		Dy: a.Dy - b.Dy,
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

// ApplyVelocity applies this body's velocity to its position
func (b *PhysicsBody) ApplyVelocity() {
	b.Position.X += b.Velocity.Dx
	b.Position.Y += b.Velocity.Dy
}

// InelasticCollision update
func InelasticCollision(b1 *PhysicsBody, b2 *PhysicsBody) {

	// print := true
	// d := math.Sqrt(math.Pow((b1.Position.X-b2.Position.X), 2) + math.Pow((b1.Position.Y-b2.Position.Y), 2))
	// if d == 0 {
	// 	return
	// }

	// SinYx := (b1.Position.Y - b2.Position.Y) / (b1.Position.X - b2.Position.X)
	// TanYv := (b2.Velocity.Dy - b1.Velocity.Dy) / (b2.Velocity.Dx - b1.Velocity.Dx)
	// if SinYx > 1 || SinYx < -1 {
	// 	SinYx = 1 / SinYx
	// 	TanYv = 1 / TanYv
	// }

	// Yx := math.Asin(SinYx)
	// Yv := math.Atan(TanYv)
	// alpha := math.Asin((d * math.Sin(Yx-Yv)) / (b2.Radius + b1.Radius))
	// a := math.Tan(Yv + alpha)
	// massRatio := float64(b1.Mass / b2.Mass)

	// Deltab1Vx := 2 * (b2.Velocity.Dx - b1.Velocity.Dx + a*(b2.Velocity.Dy-b1.Velocity.Dy)) / ((math.Pow(a, 2) + 1) * (massRatio + 1))
	// Deltab1Vx := 2 * (b2.Velocity.Dx - b1.Velocity.Dx + a*(b2.Velocity.Dy-b1.Velocity.Dy)) / ((2*a + 1) * (massRatio + 1))

	// if print {
	// 	fmt.Println("d", d)
	// 	fmt.Println("Yx", Yx)
	// 	fmt.Println("SinYx", SinYx)
	// 	fmt.Println("r1 + r2", b1.Radius+b2.Radius)
	// 	fmt.Println("Yv", Yv)
	// 	fmt.Println("Yx - Yv", Yx-Yv)
	// 	fmt.Println("Alpha", alpha)
	// 	fmt.Println("a", a)
	// 	fmt.Println("massRatio", massRatio)
	// 	fmt.Println("Delta", Deltab1Vx)
	// 	fmt.Println("B1 Dx", Deltab1Vx*b1.Restitution)
	// 	fmt.Println("B1 Dy", a*Deltab1Vx*b1.Restitution)
	// 	fmt.Println("B2 Dx", -massRatio*Deltab1Vx*b2.Restitution)
	// 	fmt.Println("B2 Dy", -a*massRatio*Deltab1Vx*b2.Restitution)
	// }

	// b1.Velocity.Dx += Deltab1Vx * b1.Restitution
	// b1.Velocity.Dy += a * Deltab1Vx * b1.Restitution
	// b2.Velocity.Dx -= massRatio * Deltab1Vx * b2.Restitution
	// b2.Velocity.Dy -= a * massRatio * Deltab1Vx * b2.Restitution

	// if b1.Velocity.Dx < 0 {
	// 	b1.Velocity.Dx = -b1.Velocity.Dx
	// }

	// if b1.Position.X > b2.Position.X {
	// 	temp := b1
	// 	b1 = b2
	// 	b2 = temp
	// }

	// v1 := b1.Velocity.magnitude()
	// v2 := b2.Velocity.magnitude()
	// theta1 := float64(0)
	// theta2 := float64(0)
	// // theta1cos := float64(0)
	// // theta2cos := float64(0)
	// if v1 != 0 {
	// 	theta1 = math.Asin(b1.Velocity.Dy / v1)
	// 	if b1.Velocity.Dx < 0 {
	// 		if theta1 > 0 {
	// 			theta1 = math.Pi - theta1
	// 		} else if theta1 > 0 {
	// 			theta1 = math.Pi + theta1
	// 		}
	// 	}
	// 	// if theta1 < 0 {
	// 	// 	theta1 = math.Pi + theta1
	// 	// }
	// 	// theta1cos = math.Acos(b1.Velocity.Dx / v1)
	// }
	// if v2 != 0 {
	// 	theta2 = math.Asin(b2.Velocity.Dy / v2)
	// 	if b2.Velocity.Dx < 0 {
	// 		if theta2 > 0 {
	// 			theta2 = math.Pi - theta2
	// 		} else if theta2 > 0 {
	// 			theta2 = math.Pi + theta2
	// 		}
	// 	}
	// 	// if theta2 < 0 {
	// 	// 	theta2 = math.Pi + theta2
	// 	// }
	// 	// theta2cos = math.Acos(b2.Velocity.Dx / v2)
	// }
	// phi := 0.0
	// if b1.Position.Y > b2.Position.Y {
	// 	phi = math.Asin((b1.Position.Y - b2.Position.Y) / (b1.Radius + b2.Radius))
	// } else {
	// 	phi = math.Asin((b2.Position.Y - b1.Position.Y) / (b1.Radius + b2.Radius))
	// }
	// // phi := alpha

	// b1factor := (v1*math.Cos(theta1-phi)*(b1.Mass-b2.Mass) + 2*b2.Mass*v2*math.Cos(theta2-phi)) / (b1.Mass + b2.Mass)
	// // fmt.Println("phi", phi)
	// // fmt.Println("alpha", alpha)
	// // fmt.Println("theta1", theta1)
	// // fmt.Println("theta1cos", theta1cos)
	// // fmt.Println("theta2", theta2)
	// // fmt.Println("theta2cos", theta2cos)
	// // fmt.Println("b1factor", b1factor)
	// b1.Velocity.Dx = b1factor*math.Cos(phi) + v1*math.Sin(theta1-phi)*math.Sin(phi)
	// b1.Velocity.Dy = b1factor*math.Sin(phi) + v1*math.Sin(theta1-phi)*math.Cos(phi)

	// b2factor := (v2*math.Cos(theta2-phi)*(b2.Mass-b1.Mass) + 2*b1.Mass*v1*math.Cos(theta1-phi)) / (b2.Mass + b1.Mass)
	// b2.Velocity.Dx = b2factor*math.Cos(phi) + v2*math.Sin(theta2-phi)*math.Sin(phi)
	// b2.Velocity.Dy = b2factor*math.Sin(phi) + v2*math.Sin(theta2-phi)*math.Cos(phi)

	// // b1.Velocity.Dx = (b1.Mass*b1.Velocity.Dx + b2.Mass*b2.Velocity.Dx) / (b1.Mass + b2.Mass)
	// // b1.Velocity.Dy = (b1.Mass*b1.Velocity.Dy + b2.Mass*b2.Velocity.Dy) / (b1.Mass + b2.Mass)

	// // b2.Velocity.Dx = (b1.Mass*b1.Velocity.Dx + b2.Mass*b2.Velocity.Dx) / (b1.Mass + b2.Mass)
	// // b2.Velocity.Dy = (b1.Mass*b1.Velocity.Dy + b2.Mass*b2.Velocity.Dy) / (b1.Mass + b2.Mass)

	x1Minusx2 := Velocity{
		Dx: b1.Position.X - b2.Position.X,
		Dy: b1.Position.Y - b2.Position.Y,
	}
	x2Minusx1 := Velocity{
		Dx: b2.Position.X - b1.Position.X,
		Dy: b2.Position.Y - b1.Position.Y,
	}
	b1.Velocity.Dx = b1.Velocity.Dx - (2*b2.Mass/(b1.Mass+b2.Mass))*(b1.Velocity.sub(b2.Velocity).dot(x1Minusx2))/(math.Pow(x1Minusx2.magnitude(), 2))*(x1Minusx2.Dx)
	b1.Velocity.Dy = b1.Velocity.Dy - (2*b2.Mass/(b1.Mass+b2.Mass))*(b1.Velocity.sub(b2.Velocity).dot(x1Minusx2))/(math.Pow(x1Minusx2.magnitude(), 2))*(x1Minusx2.Dy)

	b2.Velocity.Dx = b2.Velocity.Dx - (2*b1.Mass/(b1.Mass+b2.Mass))*(b2.Velocity.sub(b1.Velocity).dot(x2Minusx1))/(math.Pow(x2Minusx1.magnitude(), 2))*(x2Minusx1.Dx)
	b2.Velocity.Dy = b2.Velocity.Dy - (2*b1.Mass/(b1.Mass+b2.Mass))*(b2.Velocity.sub(b1.Velocity).dot(x2Minusx1))/(math.Pow(x2Minusx1.magnitude(), 2))*(x2Minusx1.Dy)

}
