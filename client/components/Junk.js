export default class Junk {
  constructor(props) {
    this.mass = props.mass;
	this.pointVal = props.pointVal;
	this.position = props.position;
	this.velocity = [0, 0];
	this.lastBumped = null;
	this.alive = true;
	
	this.drawJunk = this.drawJunk.bind(this);
  }
  
  drawJunk() {

  }
}