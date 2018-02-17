export default class Junk {
  constructor(props) {
    this.mass = props.mass || 10;
	this.pointVal = props.pointVal || 50;
	this.position = props.position || {x: (window.innerWidth / 2), y: (window.innerHeight) / 2)};
	this.velocity = {dx: 0, dy: 0};
	this.lastBumped = null;
	this.alive = true;
	
	this.drawJunk = this.drawJunk.bind(this);
  }
  
  drawJunk() {

  }
}