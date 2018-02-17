class Junk {
	var velocity = { 0.0, 0.0 };
	var lastBumped = undefined;
	var mass;
	var position;
	var pointVal;
	var alive = true; 
	
  constructor(mass, pointVal, canvas) {
    this.mass = mass;
	this.pointVal = pointVal;
	this.position = generateJunkCoordinates(); // Need to change this function to return position array, or take this line out and keep as is
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