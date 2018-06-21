import { drawPlayer } from './Player';

// const PLAYER_RADIUS = 25;
const EDGE_BUFFER = 5;
const MAP_SCALE = 12;

export default class Minimap {
  constructor(props) {
    this.canvas = props.canvas;
    this.arena = props.arena;

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
  }

  drawMap() {
    const ctx = this.canvas.getContext('2d');

    const mapWidth = this.arena.width / MAP_SCALE;
    const mapHeight = this.arena.height / MAP_SCALE;
    const mapX = this.canvas.width - mapWidth - EDGE_BUFFER;
    const mapY = this.canvas.height - mapHeight - EDGE_BUFFER;

    // draw map bg
    ctx.beginPath();
    ctx.rect(mapX, mapY, mapWidth, mapHeight);
    ctx.fillStyle = 'rgba(5,225,255,0.3)';
    ctx.fill();
    ctx.closePath();

    console.log("mapX: " + mapX + " MapY: " + mapY);
    console.log("mapWidth: " + mapWidth + " MapHeight: " + mapHeight);

    if (this.players) {
      this.players.forEach((p) => {
        const mapPlayer = JSON.parse(JSON.stringify(p));
        mapPlayer.position.x = (p.position.x / MAP_SCALE) + mapX;
        mapPlayer.position.y = (p.position.y / MAP_SCALE) + mapY;

        mapPlayer.name = 'maperplayer';

        // console.log(mapPlayer.position.x);
        drawPlayer(mapPlayer, this.canvas, 2);
      });
    }
  }
}
