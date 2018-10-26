package models

// Junk related constants
const (
	JunkFriction          = 0.98
	JunkRadius            = 15
	JunkDebounceTicks     = 10
	JunkGravityDamping    = 0.025
	JunkMass              = 1
	JunkRestitutionFactor = 1
)

// Junk a position and velocity struct describing it's state and player struct to identify rewarding points
type Junk struct {
	PhysicsBody   `json:"body"`
	Color         string  `json:"color"`
	LastPlayerHit *Player `json:"-"`
	Debounce      int     `json:"-"`
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
	// Check if the junk should reflect off the walls first
	if j.GetX()+j.GetVelocity().Dx > width-j.GetRadius() || j.GetX()+j.GetVelocity().Dx < j.GetRadius() {
		j.SetDx(-j.GetVelocity().Dx)
	}
	if j.GetY()+j.GetVelocity().Dy > height-j.GetRadius() || j.GetY()+j.GetVelocity().Dy < j.GetRadius() {
		j.SetDy(-j.GetVelocity().Dy)
	}

	j.ApplyFactor(JunkFriction)
	j.ApplyVelocity()

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
	InelasticCollision(j, p)

	j.Debounce = JunkDebounceTicks
}

// HitJunk causes a collision event between this junk and given junk.
func (j *Junk) HitJunk(jh *Junk) {
	// We don't want this collision till the debounce is down.
	if j.jDebounce != 0 {
		return
	}
	InelasticCollision(j, jh)

	j.jDebounce = JunkDebounceTicks
	jh.jDebounce = JunkDebounceTicks
}
