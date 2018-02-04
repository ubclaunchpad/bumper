class Hole {
  var position;
  var radius;
  var lifespan;
  
  // Note: Lifespan must be in milliseconds (for use in setTimeout)
  constructor(position, radius, lifespan) { 
	this.radius = radius;
	this.lifespan = lifespan;
	
	// Assuming a spawnHole function exists with this parameter
	// and returns the position the hole is placed
	this.position = spawnHole(radius); 
	
	var timerFlag = 0;
	var timer = setTimeout(timer,lifespan); // timerFlag = 1 after lifespan milliseconds
	drawHole();
  }
  
  // Continuously redraws the hole on each frame
  drawHole() {
	  while(timerFlag === 0) {
		  // CODE TO CLEAR CANVAS AND DRAW HOLE HERE
		  requestAnimationFrame(draw);
	  }
  }
  
  timer() {
	  timerFlag = 1;
  }
  
  
  // Getter method to be used to detect if a player/junk is within
  // the bounds of a hole
  getPositionAndRadius() {
	  return {position, radius};
  }
}