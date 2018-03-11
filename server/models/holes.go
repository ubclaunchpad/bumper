package models

// Hole related constants
const (
	MinHoleRadius = 15
	MaxHoleRadius = 45
	MinHoleLife   = 25
	MaxHoleLife   = 75
)

// Hole contains the data for a hole's position and size
type Hole struct {
	Position Position `json:"position"`
	Radius   float64  `json:"radius"`
	Life     float64  `json:"life"`
}

// Set this hole to a new position and lifespan
func (h *Hole) startNewLife() {

}
