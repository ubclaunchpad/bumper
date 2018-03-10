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
	ID       int
	Position Position
	Velocity Velocity
	Color    string
	Angle    float32 //Direction of player (degrees)
}

func (p *Player) updatePosition() {

}

func (p *Player) hitJunk() {

}

func (p *Player) keyHandler(key int) {

}
