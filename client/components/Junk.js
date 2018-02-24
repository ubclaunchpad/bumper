const JUNK_SIZE = 15;
const JUNK_FRICTION = 0.99;
const JUNK_MINBUMP = 0.5;

export default class Junk {
  constructor(props) {
    this.canvas = props.canvas;
    this.position = props.position;
    this.velocity = { dx: 0, dy: 0 };
    this.lastBumped = null;
    this.color = 'white';

    this.mass = props.mass || 10;
    this.pointVal = props.pointVal || 50;
    this.alive = true;
    this.drawJunk = this.drawJunk.bind(this);
  }

  drawJunk() {
    const ctx = this.canvas.getContext('2d');
    ctx.beginPath();
    ctx.rect(this.position.x, this.position.y, JUNK_SIZE, JUNK_SIZE);
    ctx.fillStyle = this.color;
    ctx.fill();
    ctx.closePath();
  }

  hitBy(player) {
    this.color = player.color;
    if (player.velocity.dx < 0) {
      this.velocity.dx = Math.min(player.velocity.dx * 1.05, -JUNK_MINBUMP);
    } else {
      this.velocity.dx = Math.max(player.velocity.dx * 1.05, JUNK_MINBUMP);
    }
    if (player.velocity.dy < 0) {
      this.velocity.dy = Math.min(player.velocity.dy * 1.05, -JUNK_MINBUMP);
    } else {
      this.velocity.dy = Math.max(player.velocity.dy * 1.05, JUNK_MINBUMP);
    }
  }

  updatePosition() {
    const r = JUNK_SIZE / 2;
    if (this.position.x + this.velocity.dx > this.canvas.width - r || this.position.x + this.velocity.dx < r) {
      this.velocity.dx = -this.velocity.dx;
    }
    if (this.position.y + this.velocity.dy > this.canvas.height - r || this.position.y + this.velocity.dy < r) {
      this.velocity.dy = -this.velocity.dy;
    }

    this.velocity.dy *= JUNK_FRICTION;
    this.velocity.dx *= JUNK_FRICTION;

    this.position.x += this.velocity.dx;
    this.position.y += this.velocity.dy;
  }
}
