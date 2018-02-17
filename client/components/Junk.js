export default class Junk {
	
  constructor(mass, pointVal, canvas) {
    this.mass = mass;
	this.pointVal = pointVal;
	this.position = generateJunkCoordinates(); // Need to change this function to return position array, or take this line out and keep as is
	this.velocity = [0.0, 0.0];
	this.lastBumped = null;
	this.alive = true; 
  }
  
  drawJunk() {
	// Check if new coords are within a hole
	// If yes, set alive to false
	if(!alive) {
		// Don't redraw - game over code
	}
	else {
		// Draw player with current coords
	}
  }
  
}