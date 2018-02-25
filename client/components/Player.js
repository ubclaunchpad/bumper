import { magnitude, normalize } from '../utils/utils';

const PLAYER_RADIUS = 25;
const MAX_VELOCITY = 15;
const PLAYER_ACCELERATION = 0.5;
const PLAYER_FRICTION = 0.97;
const WALL_BOUNCE_FACTOR = 1.5;
const JUNK_BOUNCE_FACTOR = 0.25;

export default class Player {
  constructor(props) {
    this.canvas = props.canvas;
    this.position = props.position || { x: window.innerWidth / 2, y: window.innerHeight / 2 };
    this.velocity = { dx: 0, dy: 0 };
    this.theta = props.theta;
    this.mass = props.mass || 10;

    this.alive = true;
    this.name = props.name || 'Default name';

    let c = '';
    while (c.length < 6) {
      c += (Math.random()).toString(16).substr(-6).substr(-1);
    }
    this.color = `#${c}`;

    this.rightPressed = false;
    this.leftPressed = false;
    this.upPressed = false;
    this.downPressed = false;

    this.drawPlayer = this.drawPlayer.bind(this);
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);
    window.addEventListener('keydown', this.keyDownHandler);
    window.addEventListener('keyup', this.keyUpHandler);
  }

  drawPlayer() {
    const ctx = this.canvas.getContext('2d');
    const { x, y } = this.position;
    ctx.beginPath();
    ctx.arc(x, y, PLAYER_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = this.color;
    ctx.fill();
    ctx.closePath();

    ctx.beginPath();
    ctx.moveTo(x + (PLAYER_RADIUS * Math.sin(this.theta)), y + (PLAYER_RADIUS * Math.cos(this.theta)));
    ctx.lineTo(x - (PLAYER_RADIUS * Math.sin(this.theta)), y - (PLAYER_RADIUS * Math.cos(this.theta)));
    ctx.strokeStyle = '#000000';
    ctx.strokeWidth = 5;
    ctx.stroke();

    const backCenterX = x - ((PLAYER_RADIUS * Math.sin(this.theta)) / 2);
    const backCenterY = y - ((PLAYER_RADIUS * Math.cos(this.theta)) / 2);
    const backLength = (2.5 * ((PLAYER_RADIUS / 2) / Math.tan(45)));
    ctx.beginPath();
    ctx.moveTo(backCenterX - (backLength * Math.cos(this.theta)), backCenterY + (backLength * Math.sin(this.theta)));
    ctx.lineTo(backCenterX + (backLength * Math.cos(this.theta)), backCenterY - (backLength * Math.sin(this.theta)));
    ctx.strokeStyle = '#0000000';
    ctx.strokeWidth = 5;
    ctx.stroke();
  }

  updatePosition() {
    const controlsVector = { dx: 0, dy: 0 };

    if (this.leftPressed) {
      this.theta = (this.theta + 0.1) % 360;
    }

    if (this.rightPressed) {
      this.theta = (this.theta - 0.1) % 360;
    }

    if (this.downPressed) {
      controlsVector.dy = -(0.5 * (PLAYER_RADIUS * Math.cos(this.theta)));
      controlsVector.dx = -(0.5 * (PLAYER_RADIUS * Math.sin(this.theta)));
    }

    if (this.upPressed) {
      controlsVector.dy = (0.5 * (PLAYER_RADIUS * Math.cos(this.theta)));
      controlsVector.dx = (0.5 * (PLAYER_RADIUS * Math.sin(this.theta)));
    }

    // Normalize controls vector and apply speed
    normalize(controlsVector);
    controlsVector.dx *= PLAYER_ACCELERATION;
    controlsVector.dy *= PLAYER_ACCELERATION;

    // Apply some friction damping
    this.velocity.dx = this.velocity.dx * PLAYER_FRICTION;
    this.velocity.dy = this.velocity.dy * PLAYER_FRICTION;

    this.velocity.dx += controlsVector.dx;
    this.velocity.dy += controlsVector.dy;

    // console.log("vdx: " + this.velocity.dx + "vdy: " + this.velocity.dy);

    // Ensure it never gets going too fast
    if (magnitude(this.velocity) > MAX_VELOCITY) {
      normalize(this.velocity);
      this.velocity.dx = this.velocity.dx * MAX_VELOCITY;
      this.velocity.dy = this.velocity.dy * MAX_VELOCITY;
    }

    // Apply player's velocity vector
    this.position.x += this.velocity.dx;
    this.position.y += this.velocity.dy;

    // Check wall collisions
    if (this.position.x + PLAYER_RADIUS > this.canvas.width) {
      this.velocity.dx = -this.velocity.dx * WALL_BOUNCE_FACTOR;
    } else if (this.position.x - PLAYER_RADIUS < 0) {
      this.velocity.dx = -this.velocity.dx * WALL_BOUNCE_FACTOR;
    }

    if (this.position.y + PLAYER_RADIUS > this.canvas.height) {
      this.velocity.dy = -this.velocity.dy * WALL_BOUNCE_FACTOR;
    } else if (this.position.y - PLAYER_RADIUS < 0) {
      this.velocity.dy = -this.velocity.dy * WALL_BOUNCE_FACTOR;
    }
  }

  hitJunk() {
    this.velocity.dx *= -JUNK_BOUNCE_FACTOR;
    this.velocity.dy *= -JUNK_BOUNCE_FACTOR;
  }

  keyDownHandler(e) {
    if (e.keyCode === 39) {
      this.rightPressed = true;
    } else if (e.keyCode === 37) {
      this.leftPressed = true;
    } else if (e.keyCode === 38) {
      this.upPressed = true;
    } else if (e.keyCode === 40) {
      this.downPressed = true;
    }
  }

  keyUpHandler(e) {
    if (e.keyCode === 39) {
      this.rightPressed = false;
    } else if (e.keyCode === 37) {
      this.leftPressed = false;
    } else if (e.keyCode === 38) {
      this.upPressed = false;
    } else if (e.keyCode === 40) {
      this.downPressed = false;
    }
  }
}