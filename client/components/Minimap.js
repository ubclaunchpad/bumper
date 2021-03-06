import React from 'react';
import { drawPlayer, drawJunk, drawMapHole } from './GameObjects';

const EDGE_BUFFER = 5;
const BORDER_WIDTH = 3;
const MAP_SCALE = 12;
const OBJECT_SCALE = 4;
const HOLE_SCALE = 10;

export default class Minimap extends React.Component {
  constructor(props) {
    super(props);

    this.junk = props.junk;
    this.holes = props.holes;
    this.players = props.players;

    this.drawMap = this.drawMap.bind(this);
  }

  async componentDidMount() {
    this.canvas = document.getElementById('mapctx');

    this.mapWidth = this.props.arena.width / MAP_SCALE;
    this.mapHeight = this.props.arena.height / MAP_SCALE;

    this.mapX = this.canvas.width - this.mapWidth - EDGE_BUFFER;
    this.mapY = this.canvas.height - this.mapHeight - EDGE_BUFFER;
  }

  componentDidUpdate() {
    this.junk = this.props.junk;
    this.holes = this.props.holes;
    this.players = this.props.players;
  }

  drawMapBorder() {
    const ctx = this.canvas.getContext('2d');

    ctx.beginPath();
    ctx.rect(this.mapX - BORDER_WIDTH, this.mapY - BORDER_WIDTH, BORDER_WIDTH, this.mapHeight + (BORDER_WIDTH * 2)); // Left
    ctx.rect(this.mapX + this.mapWidth, this.mapY - BORDER_WIDTH, BORDER_WIDTH, this.mapHeight + (BORDER_WIDTH * 2)); // Right
    ctx.rect(this.mapX - BORDER_WIDTH, this.mapY - BORDER_WIDTH, this.mapWidth + (BORDER_WIDTH * 2), BORDER_WIDTH); // Top
    ctx.rect(this.mapX - BORDER_WIDTH, this.mapY + this.mapHeight, this.mapWidth + (BORDER_WIDTH * 2), BORDER_WIDTH); // Bottom
    ctx.fillStyle = 'yellow';
    ctx.fill();
    ctx.closePath();
  }

  drawMap() {
    const ctx = this.canvas.getContext('2d');

    // Deep copy and then translate all game objects
    const holes = this.holes.map((h) => {
      if (!h.position) {
        return {
          position: { x: 0, y: 0 },
          radius: 0,
          isAlive: false,
        };
      }

      const newPosition = {
        x: (h.position.x / MAP_SCALE) + this.mapX,
        y: (h.position.y / MAP_SCALE) + this.mapY,
      };
      return {
        position: newPosition,
        radius: h.radius,
        isAlive: h.isAlive,
      };
    });

    const players = this.players.map((p) => {
      if (!p.position) {
        return {
          position: { x: 0, y: 0 },
          color: 'black',
          angle: 0,
        };
      }

      const newPosition = {
        x: (p.position.x / MAP_SCALE) + this.mapX,
        y: (p.position.y / MAP_SCALE) + this.mapY,
      };
      return {
        position: newPosition,
        color: p.color,
        angle: p.angle,
        name: p.name,
      };
    });

    const junk = this.junk.map((j) => {
      if (!j.position) {
        return {
          position: { x: 0, y: 0 },
          color: 'black',
        };
      }

      const newPosition = {
        x: (j.position.x / MAP_SCALE) + this.mapX,
        y: (j.position.y / MAP_SCALE) + this.mapY,
      };
      return {
        position: newPosition,
        color: j.color,
      };
    });

    // Draw map bg
    ctx.beginPath();
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    this.drawMapBorder();
    ctx.rect(this.mapX, this.mapY, this.mapWidth, this.mapHeight);
    ctx.fillStyle = 'rgba(5,225,255,0.3)';
    ctx.fill();
    ctx.closePath();

    // Draw all game objects on map
    junk.forEach((j) => { drawJunk(j, this.canvas, OBJECT_SCALE); });
    players.forEach((p) => { drawPlayer(p, this.canvas, OBJECT_SCALE); });
    holes.forEach((h) => { drawMapHole(h, this.canvas, HOLE_SCALE); });
  }

  render() {
    if (this.canvas) {
      this.drawMap();
    }

    return (
      <div style={styles.canvasContainer}>
        <canvas id="mapctx" style={styles.canvas} display="inline" width={window.innerWidth - 20} height={window.innerHeight - 20} margin={0} />
      </div>
    );
  }
}

const styles = {
  canvas: {
    textAlign: 'center',
  },
  canvasContainer: {
    textAlign: 'center',
    position: 'fixed',
    left: 0,
    right: 0,
    top: 0,
    bottom: 0,
  },
};
