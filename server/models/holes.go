package models

// MaxRadius is the largest a hole will grow before disappearing
const MaxRadius = 45

// Hole contains the data for a hole's position and size
type Hole struct {
	Position Position `json:"position"`
	Radius   float64  `json:"radius"`
	Life     float64  `json:"life"`
}

// Set this hole to a new position and lifespan
func (h *Hole) startNewLife() {

}
