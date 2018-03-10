package models

// LeftKey is the left keyboard button
const LeftKey = 37

// RightKey is the left keyboard button
const RightKey = 39

// UpKey is the up keyboard button
const UpKey = 38

// DownKey is the down keyboard button
const DownKey = 40

// Player contains data and state about a player's object
type Player struct {
	ID       int      `json:"id"`
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
	Color    string   `json:"color"`
	Angle    float32  `json:"angle"`
}

//Update Player's position based on calculations of position/velocity
func (p *Player) updatePosition() {

}

//Update Player's position based on calculations of hitting junk
func (p *Player) hitJunk() {

}

//Handle the effect of a key press
func (p *Player) keyHandler(key int) {

}
