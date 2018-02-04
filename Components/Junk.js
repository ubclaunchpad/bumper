class Junk {
	var velocity = { 0.0, 0.0 };
	var lastBumped = undefined;
	var mass;
	var position;
	var pointVal;
	var alive = true; // NEED TO CHANGE TO FALSE IF JUNK FALLS INTO THE AREA OF A HOLE, ETC.
	
  constructor(mass, pointVal) {
    this.mass = mass;
	this.pointVal = pointVal;
	
	// Assuming a spawnJunk function exists with this parameter
	// and returns the position the hole is placed
	this.position = spawnJunk();
	
	move();
  }

  move() { // Assuming velocity instance variable is updated if a collision occurs
	 var stationary = (this.velocity === { 0.0, 0.0 });
	 if(!stationary) { // Move junk based on velocity
		// Calculate new position
		// ADD CODE FOR CALCULATING NEW POSITION BASED ON THIS.POSITION AND VELOCITY
		// PUT NEW POSITION IN THIS.POSITION
	  }
	  // Draw new position
	  // ADD CODE FOR DRAWING NEW POSITION
	  if(alive) {
		requestAnimationFrame(move); // Continue looping move() if alive
	  }
  }
  
}