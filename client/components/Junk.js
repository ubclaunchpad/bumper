export default class Junk {
  constructor(props) {
    this.mass = props.mass || 10;
	this.pointVal = props.pointVal || 50;
	this.position = props.position || { x: (window.innerWidth / 2), y: (window.innerHeight) / 2) };
	this.velocity = { dx: 0, dy: 0 };
	this.lastBumped = null;
	this.alive = true;
	
	this.drawJunk = this.drawJunk.bind(this);
  }
  
  drawJunk() {
	const ctx = this.canvas.getContext('2d');    
    for (const p of this.state.junkCoords) {
      ctx.beginPath();
      ctx.rect(p.x, p.y, JUNK_SIZE, JUNK_SIZE);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
  }
}
