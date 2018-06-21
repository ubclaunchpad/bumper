import { drawPlayer } from './Player';

const JUNK_SIZE = 2;

const EDGE_BUFFER = 5;
const MAP_SCALE = 12;

export default class Minimap {
  constructor(props) {
    this.canvas = props.canvas;

    this.mapWidth = props.arena.width / MAP_SCALE;
    this.mapHeight = props.arena.height / MAP_SCALE;
    this.mapX = props.canvas.width - this.mapWidth - EDGE_BUFFER;
    this.mapY = props.canvas.height - this.mapHeight - EDGE_BUFFER;

    this.junk = null;
    this.holes = null;
    this.players = null;

    this.drawMap = this.drawMap.bind(this);
    this.update = this.update.bind(this);
  }

  update(data) {
    this.junk = JSON.parse(JSON.stringify(data.junk));
    this.holes = JSON.parse(JSON.stringify(data.holes));
    this.players = JSON.parse(JSON.stringify(data.players));

    this.players.forEach((p) => {
      p.position.x = (p.position.x / MAP_SCALE) + this.mapX;
      p.position.y = (p.position.y / MAP_SCALE) + this.mapY;
    });
    this.holes.forEach((h) => {
      h.position.x = (h.position.x / MAP_SCALE) + this.mapX;
      h.position.y = (h.position.y / MAP_SCALE) + this.mapY;
    });
    this.junk.forEach((j) => {
      j.position.x = (j.position.x / MAP_SCALE) + this.mapX;
      j.position.y = (j.position.y / MAP_SCALE) + this.mapY;
    });
  }

  drawMap() {
    const ctx = this.canvas.getContext('2d');

    // draw map bg
    ctx.beginPath();
    ctx.rect(this.mapX, this.mapY, this.mapWidth, this.mapHeight);
    ctx.fillStyle = 'rgba(5,225,255,0.3)';
    ctx.fill();
    ctx.closePath();

    // draw players and whatnot on map
    if (this.players) {
      this.players.forEach((p) => {
        drawPlayer(p, this.canvas, 4);
      });
    }
    if (this.junk) {
      this.junk.forEach((j) => {
        ctx.beginPath();
        ctx.rect(j.position.x - (JUNK_SIZE / 2), j.position.y - (JUNK_SIZE / 2), JUNK_SIZE, JUNK_SIZE);
        ctx.fillStyle = j.color;
        ctx.fill();
        ctx.closePath();
      });
    }
    if (this.holes) {
      this.holes.forEach((h) => {
        ctx.beginPath();
        ctx.arc(h.position.x, h.position.y, h.radius / 8, 0, 2*Math.PI);
        ctx.fillStyle = 'rgb(255,225,225)';
        ctx.fill();
        ctx.stroke();
      });
    }
  }
}
