const PLAYER_RADIUS = 25;

export default class Player {
  constructor(props) {
    this.canvas = props.canvas;
    this.position = props.position || { x: window.innerWidth / 2, y: window.innerHeight / 2 };
    this.velocity = { dx: 0, dy: 0 };
    this.theta = props.theta;

    this.mass = props.mass || 10;
    this.color = props.color || '#FF0000';
    this.name = props.name || 'Default name';
    this.alive = true;

    this.drawPlayer = this.drawPlayer.bind(this);
  }

  drawPlayer() {
    const ctx = this.canvas.getContext('2d');
    const { x, y } = this.position;
    ctx.beginPath();
    ctx.arc(x, y, PLAYER_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = '#00FFFF';
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
}
