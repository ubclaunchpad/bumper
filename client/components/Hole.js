export default class Hole {
  // Note: Lifespan must be in milliseconds (for use in setTimeout)
  constructor(props) {
    this.canvas = props.canvas;
    this.radius = props.radius;
    this.lifespan = props.lifespan;
    this.position = {};

    this.getPositionAndRadius = this.getPositionAndRadius.bind(this);
    this.drawHole = this.drawHole.bind(this);
  }

  getPositionAndRadius() {
    return { position: this.position, radius: this.radius };
  }

  drawHole() {

  }
}
