package models

import (
	"fmt"
	"math"
)

// Junk related constants
const (
	JunkFriction          = 0.98
	MinimumBump           = 0.6
	BumpFactor            = 1.05
	JunkRadius            = 18 // 11
	JunkDebounceTicks     = 15
	JunkVTransferFactor   = 0.5
	JunkJunkBounceFactor  = 0.01
	JunkGravityDamping    = 0.025
	JunkMass              = 1
	JunkRestitutionFactor = 1
)

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	Position      Position `json:"position"`
	Velocity      Velocity `json:"-"`
	Color         string   `json:"color"`
	LastPlayerHit *Player  `json:"-"`
	Debounce      int      `json:"-"`
	jDebounce     int
}

// CreateJunk initializes and returns an instance of a Junk
func CreateJunk(position Position) *Junk {
	return &Junk{
		Position:  position,
		Velocity:  Velocity{0, 0},
		Color:     "white",
		Debounce:  0,
		jDebounce: 0,
	}
}

// UpdatePosition Update Junk's position based on calculations of position/velocity
func (j *Junk) UpdatePosition(height float64, width float64) {
	if j.Position.X+j.Velocity.Dx > width-JunkRadius || j.Position.X+j.Velocity.Dx < JunkRadius {
		j.Velocity.Dx = -j.Velocity.Dx
	}
	if j.Position.Y+j.Velocity.Dy > height-JunkRadius || j.Position.Y+j.Velocity.Dy < JunkRadius {
		j.Velocity.Dy = -j.Velocity.Dy
	}

	j.Velocity.Dx *= JunkFriction
	j.Velocity.Dy *= JunkFriction

	j.Position.X += j.Velocity.Dx
	j.Position.Y += j.Velocity.Dy

	if j.Debounce > 0 {
		j.Debounce--
	} else {
		j.Debounce = 0
	}

	if j.jDebounce > 0 {
		j.jDebounce--
	} else {
		j.jDebounce = 0
	}
}

// HitBy Update Junks's velocity based on calculations of being hit by a player
func (j *Junk) HitBy(p *Player) {
	// We don't want this collision till the debounce is down.
	if j.Debounce != 0 {
		return
	}

		// J is 2 P is 1

		d := math.Sqrt(math.Pow((j.Position.X-p.Body.Position.X), 2) + math.Pow((j.Position.Y-p.Body.Position.Y), 2))
		fmt.Println("d", d)
		SinYx := (j.Position.Y - p.Body.Position.Y) / (j.Position.X - p.Body.Position.X)
		// Yx := math.Asin((j.Position.Y - p.Position.Y) / (j.Position.X - p.Position.X))
		if SinYx > 1 || SinYx < -1 {
			fmt.Println("x1", p.Body.Position.X, "y1", p.Body.Position.Y)
			fmt.Println("x2", j.Position.X, "y2", j.Position.Y)
			InvSinYx := (p.Body.Position.Y - j.Position.Y) / (p.Body.Position.X - j.Position.X)
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
		r1plusr2 := JunkRadius + PlayerRadius
		fmt.Println("Yx", Yx)
		fmt.Println("SinYx", SinYx)
		fmt.Println("r1 + r2", r1plusr2)

		Yv := math.Atan((p.Body.Velocity.Dy - j.Velocity.Dy) / (p.Body.Velocity.Dx - j.Velocity.Dx))
		fmt.Println("Yv", Yv)
		alpha := math.Asin((d * math.Sin(Yx-Yv)) / (PlayerRadius + JunkRadius))
		fmt.Println("Alpha", alpha)
		a := math.Tan(Yv + alpha)
		fmt.Println("a", a)
		massRatio := float64(JunkMass / PlayerMass)
		fmt.Println("massRatio", massRatio)

		DeltaJVx := 2 * (p.Body.Velocity.Dx - j.Velocity.Dx + a*(p.Body.Velocity.Dy-j.Velocity.Dy)) / ((math.Pow(a, 2) + 1) * (massRatio + 1))
		fmt.Println("Delta", DeltaJVx)

		j.Velocity.Dx += DeltaJVx * JunkRestitutionFactor
		j.Velocity.Dy += a * DeltaJVx * JunkRestitutionFactor
		p.Body.Velocity.Dx -= massRatio * DeltaJVx * PlayerRestitutionFactor
		p.Body.Velocity.Dy -= a * massRatio * DeltaJVx * PlayerRestitutionFactor

		// direction := Velocity{0, 0}
		// direction.Dx = j.Position.X - p.Position.X
		// direction.Dy = j.Position.Y - p.Position.Y
		// direction.normalize()

		// j.Velocity.Dx += direction.Dx * math.Max(math.Abs(p.Velocity.Dx)*BumpFactor, MinimumBump)
		// j.Velocity.Dy += direction.Dy * math.Max(math.Abs(p.Velocity.Dy)*BumpFactor, MinimumBump)

		// p.hitJunk()
		j.Debounce = JunkDebounceTicks
	}

	p.hitJunk()
	j.Debounce = JunkDebounceTicks
}

// HitJunk Update Junks's velocity based on calculations of being hit by another Junk
func (j *Junk) HitJunk(jh *Junk) {
	// We don't want this collision till the debounce is down.
	if j.jDebounce != 0 {
		return
	}

		j.Velocity.Dx *= JunkJunkBounceFactor
		j.Velocity.Dy *= JunkJunkBounceFactor
		jh.Velocity.Dy *= JunkJunkBounceFactor
		jh.Velocity.Dy *= JunkJunkBounceFactor

		direction := Velocity{0, 0}
		if j.Position.X > jh.Position.X {
			direction.Dx = j.Position.X - jh.Position.X
			direction.Dy = j.Position.Y - jh.Position.Y
			direction.normalize()
			// fmt.Println("1")
			j.Velocity.Dx -= direction.Dx * jh.Velocity.Dx * JunkVTransferFactor
			j.Velocity.Dy -= direction.Dy * jh.Velocity.Dy * JunkVTransferFactor
			jh.Velocity.Dx += direction.Dx * initalVelocity.Dx * JunkVTransferFactor
			jh.Velocity.Dy += direction.Dy * initalVelocity.Dy * JunkVTransferFactor
		} else {
			direction.Dx = jh.Position.X - j.Position.X
			direction.Dy = jh.Position.Y - j.Position.Y
			direction.normalize()
			// fmt.Println("2")
			j.Velocity.Dx -= direction.Dx * jh.Velocity.Dx * JunkVTransferFactor
			j.Velocity.Dy -= direction.Dy * jh.Velocity.Dy * JunkVTransferFactor
			jh.Velocity.Dx += direction.Dx * initalVelocity.Dx * JunkVTransferFactor
			jh.Velocity.Dy += direction.Dy * initalVelocity.Dy * JunkVTransferFactor
		}

		// //Calculate this junks's new velocity
		// j.Velocity.Dx += direction.Dx * jh.Velocity.Dx * JunkVTransferFactor
		// // fmt.Println(direction.Dx * jh.Velocity.Dx * JunkVTransferFactor)
		// j.Velocity.Dy += direction.Dy * jh.Velocity.Dy * JunkVTransferFactor
		// jh.Velocity.Dx += direction.Dx * initalVelocity.Dx * JunkVTransferFactor
		// jh.Velocity.Dy += direction.Dy * initalVelocity.Dy * JunkVTransferFactor

		// j.Velocity.Dx = (j.Velocity.Dx * -JunkVTransferFactor) + (jh.Velocity.Dx * JunkVTransferFactor)
		// j.Velocity.Dy = (j.Velocity.Dy * -JunkVTransferFactor) + (jh.Velocity.Dy * JunkVTransferFactor)
		// //Calculate other junk's new velocity
		// jh.Velocity.Dx = (jh.Velocity.Dx * -JunkVTransferFactor) + (initalVelocity.Dx * JunkVTransferFactor)
		// jh.Velocity.Dy = (jh.Velocity.Dy * -JunkVTransferFactor) + (initalVelocity.Dy * JunkVTransferFactor)

	//Calculate other junk's new velocity
	jh.Velocity.Dx = (jh.Velocity.Dx * -JunkVTransferFactor) + (initalVelocity.Dx * JunkVTransferFactor)
	jh.Velocity.Dy = (jh.Velocity.Dy * -JunkVTransferFactor) + (initalVelocity.Dy * JunkVTransferFactor)

	j.jDebounce = JunkDebounceTicks
	jh.jDebounce = JunkDebounceTicks
}
