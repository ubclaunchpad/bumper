class Player {
	var velocity = { 0.0, 0.0 };
	var position;
	var mass;
	var color;
	var name;
	var alive = true; // NEED TO CHANGE TO FALSE IF PLAYER FALLS INTO THE AREA OF A HOLE, ETC.
	
	// Note: Only mass and position are required parameters to the Player contructor
  constructor(mass, color, name) { 
    this.mass = mass;
	this.color = color; // May be undefined for now
	this.name = name; // May be undefined for now
	
	// Assuming a spawnPlayer function exists with these parameters,
	// and returns the position the player is placed
	this.position = spawnPlayer(color, name);
	
	move();
  }
	
   move() { // Assuming velocity instance variable is updated if a collision occurs
	updateVelocity();
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
	  
	// Can either input the arrow key that is currently being pressed
	// (NOT TAKING INTO ACCOUNT 2 BEING PRESSED AT THE SAME TIME YET),
	// or a calculated velocity value (i.e. from after a collision)
	updateVelocity(value) {
		switch(value) {
			case up:
				// CALCULATION : TAKE X AND Y COMPONENTS OF VELOCITY
				// AND MAKE THEM ALL IN THE NEGATIVE (b/c switched with canvas) Y DIRECTION
				break;
			case down:
				// CALCULATION : TAKE X AND Y COMPONENTS OF VELOCITY
				// AND MAKE THEM ALL IN THE POSITIVE (b/c switched with canvas) Y DIRECTION
				break;
			case left:
				// CALCULATION : TAKE X AND Y COMPONENTS OF VELOCITY
				// AND MAKE THEM ALL IN THE NEGATIVE X DIRECTION
				break;
			case right:
				// CALCULATION : TAKE X AND Y COMPONENTS OF VELOCITY
				// AND MAKE THEM ALL IN THE POSITIVE X DIRECTION
				break;
			case noKeyPressed:
				velocity = { 0.0, 0.0 };
				break;
			default: this.velocity = velocity
	}
  
}