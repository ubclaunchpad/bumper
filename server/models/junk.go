package models

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
	Body          PhysicsBody `json:"body"`
	Color         string      `json:"color"`
	LastPlayerHit *Player     `json:"-"`
	Debounce      int         `json:"-"`
	jDebounce     int
}

// CreateJunk initializes and returns an instance of a Junk
func CreateJunk(position Position) *Junk {
	return &Junk{
		Body:      CreateBody(position, JunkRadius, JunkMass, JunkRestitutionFactor),
		Color:     "white",
		Debounce:  0,
		jDebounce: 0,
	}
}

// UpdatePosition Update Junk's position based on calculations of position/velocity
func (j *Junk) UpdatePosition(height float64, width float64) {
	if j.Body.Position.X+j.Body.Velocity.Dx > width-j.Body.Radius || j.Body.Position.X+j.Body.Velocity.Dx < j.Body.Radius {
		j.Body.Velocity.Dx = -j.Body.Velocity.Dx
	}
	if j.Body.Position.Y+j.Body.Velocity.Dy > height-j.Body.Radius || j.Body.Position.Y+j.Body.Velocity.Dy < j.Body.Radius {
		j.Body.Velocity.Dy = -j.Body.Velocity.Dy
	}

	j.Body.Velocity.ApplyFactor(JunkFriction)
	j.Body.ApplyVelocity()

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

// HitBy causes a collision event between this junk and given player.
func (j *Junk) HitBy(p *Player) {
	// We don't want this collision till the debounce is down.
	if j.Debounce != 0 {
		return
	}

	j.Color = p.Color //Assign junk to last recently hit player color
	j.LastPlayerHit = p
	InelasticCollision(&j.Body, &p.Body)

	j.Debounce = JunkDebounceTicks
}

// HitJunk causes a collision event between this junk and given junk.
func (j *Junk) HitJunk(jh *Junk) {
	// We don't want this collision till the debounce is down.
	if j.jDebounce != 0 {
		return
	}
	InelasticCollision(&j.Body, &jh.Body)

	j.jDebounce = JunkDebounceTicks
	jh.jDebounce = JunkDebounceTicks
}
