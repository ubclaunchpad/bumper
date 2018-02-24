const JUNK_SIZE = 15;
const HOLE_COUNT = 10;

export default class Junk {
  constructor(props) {
    this.canvas = props.canvas;
    this.mass = props.mass || 10;
    this.pointVal = props.pointVal || 50;
    this.position = props.position;
    this.velocity = { dx: 1, dy: 1 };
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

  updatePosition() {
    if (this.position.x + this.velocity.dx > this.canvas.width - JUNK_SIZE || this.position.x + this.velocity.dx < JUNK_SIZE) {
      this.velocity.dx = -this.velocity.dx;
    }
    if (this.position.y + this.velocity.dy > this.canvas.height - JUNK_SIZE || this.position.y + this.velocity.dy < JUNK_SIZE) {
      this.velocity.dy = -this.velocity.dy;
    }

    this.position.x += this.velocity.dx;
    this.position.y += this.velocity.dy;
  }
}
