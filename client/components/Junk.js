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
}
