class Player {
	var position = {null, null}; 
	var mass = null;
	var color = null;
	var name = null;
	var canvas = null;
	var velocity = {0.0, 0.0};
	var alive = true; 
	
	constructor(mass, color, name, canvas) { 
		this.mass = mass;
		this.color = color; 
		this.name = name;
		this.canvas = canvas;
		this.position = generatePlayerCoords(); // Need to change this function to return position array, or take this line out and keep as is
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
		
	drawPlayer() {
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