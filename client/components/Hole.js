const MAX_RADIUS = 45;

export default class Hole {
  constructor(props) {
    this.canvas = props.canvas;
    this.position = props.position;
    this.radius = props.radius;
    this.lifespan = props.lifespan;
    this.getPositionAndRadius = this.getPositionAndRadius.bind(this);
    this.drawHole = this.drawHole.bind(this);

    this.lifeTimer = setInterval(
      () => {
        this.lifespan -= 1;
        if (this.radius < MAX_RADIUS) this.radius += 0.25;
      },
      250,
    );
  }

  startNewLife(newCoords, newRadius, newLifespan) {
    this.position = newCoords;
    this.radius = newRadius;
    this.lifespan = newLifespan;
  }

  getPositionAndRadius() {
    return { position: this.position, radius: this.radius };
  }

  drawHole() {
    if (this.lifespan > 0) {
      const ctx = this.canvas.getContext('2d');
      ctx.beginPath();
      ctx.arc(this.position.x, this.position.y, this.radius, 0, Math.PI * 2);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
      return true;
    }

    return false;
  }
}
