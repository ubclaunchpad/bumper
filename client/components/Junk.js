const JUNK_SIZE = 15;
const HOLE_COUNT = 10;

export default class Junk {
  constructor(props) {
	this.canvas = props.canvas;
  	this.mass = props.mass || 10;
	this.pointVal = props.pointVal || 50;
	this.position = props.position;
	this.velocity = { dx: 0, dy: 0 };
	this.lastBumped = null;
	this.alive = true;
	this.walls = {left: false, right: false, top: false, bottom: false};
	
	this.drawJunk = this.drawJunk.bind(this);
  }
  
  drawJunk() {
    const ctx = this.canvas.getContext('2d');   
    ctx.beginPath();
    ctx.rect(this.position.x, this.position.y, JUNK_SIZE, JUNK_SIZE);
    ctx.fillStyle = 'white';
    ctx.fill();
    ctx.closePath();
  }
  
  updatePosition(screen) {
	  // To-do: Add code for updating junk's position
	  
	  const junkRadius = JUNK_SIZE / 2;
	  var topWall = this.position.y - junkRadius == 0;
	  var bottomWall = this.position.y + junkRadius + 20 == height;
	  var leftWall = this.position.x - junkRadius == 0;
	  var rightWall = this.position.x + junkRadius + 20 == width;
	  
	  // Wall detection for the junk object: (should test once junk objects are moving)
	  if(topWall) { this.walls.top = true; }
	  		else { this.walls.top = false; }
	  if(bottomWall) { this.walls.bottom = true; }
	 		else { this.walls.bottom = false; }
	  if(leftWall) { this.walls.left = true; }
	  		else { this.walls.left = false; }
	  if(rightWall) { this.walls.right = true; }
	  		else { this.walls.right = false; }
  }
  
}