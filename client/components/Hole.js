export default class Hole {
  // Note: Lifespan must be in milliseconds (for use in setTimeout)
  constructor(position, radius, lifespan, canvas) { 
	this.radius = radius;
	this.lifespan = lifespan;
	this.position = generateHoleCoordinates(); // Need to change this function to return position array, or take this line out and keep as is
	
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