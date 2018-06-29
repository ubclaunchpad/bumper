import React from 'react';
import { drawPlayer, drawJunk, drawMapHole } from './GameObjects';

const EDGE_BUFFER = 5;
const MAP_SCALE = 12;
const OBJECT_SCALE = 4;
const HOLE_SCALE = 10;

export default class Minimap extends React.Component {
  constructor(props) {
    super(props);
    this.canvas = props.canvas;

    this.mapWidth = props.arena.width / MAP_SCALE;
    this.mapHeight = props.arena.height / MAP_SCALE;
    this.mapX = props.canvas.width - this.mapWidth - EDGE_BUFFER;
    this.mapY = props.canvas.height - this.mapHeight - EDGE_BUFFER;

    this.state = {
      junk: props.junk,
      holes: props.holes,
      players: props.players,
    };

    this.drawMap = this.drawMap.bind(this);
  }

  drawMap() {
    const ctx = this.canvas.getContext('2d');

    // Deep copy and then translate all game objects
    const junk = JSON.parse(JSON.stringify(this.state.junk));
    const holes = JSON.parse(JSON.stringify(this.state.holes));
    const players = JSON.parse(JSON.stringify(this.state.players));

    players.forEach((p) => {
      p.position.x = (p.position.x / MAP_SCALE) + this.mapX;
      p.position.y = (p.position.y / MAP_SCALE) + this.mapY;
    });
    holes.forEach((h) => {
      h.position.x = (h.position.x / MAP_SCALE) + this.mapX;
      h.position.y = (h.position.y / MAP_SCALE) + this.mapY;
    });
    junk.forEach((j) => {
      j.position.x = (j.position.x / MAP_SCALE) + this.mapX;
      j.position.y = (j.position.y / MAP_SCALE) + this.mapY;
    });

    // draw map bg
    ctx.beginPath();
    ctx.rect(this.mapX, this.mapY, this.mapWidth, this.mapHeight);
    ctx.fillStyle = 'rgba(5,225,255,0.3)';
    ctx.fill();
    ctx.closePath();

    // draw players and whatnot on map
    players.forEach((p) => {
      drawPlayer(p, this.canvas, OBJECT_SCALE);
    });
    junk.forEach((j) => {
      drawJunk(j, this.canvas, OBJECT_SCALE);
    });
    holes.forEach((h) => {
      drawMapHole(h, this.canvas, HOLE_SCALE);
    });
  }

  render() {
    this.drawMap();

    return (
      // <div style={styles.mapModal}>
      <div>
        Minimap
      </div>
    );
  }
}

// const styles = {
//   mapModal: {
//     display: 'flex',
//     flexDirection: 'column',
//     justifyContent: 'center',
//     position: 'fixed',
//     top: 0,
//     bottom: 0,
//     left: 0,
//     right: 0,
//     backgroundColor: 'rgba(5,225,255,0.3)',
//     borderRadius: 5,
//     height: window.innerHeight / MAP_SCALE,
//     width: window.innerWidth / MAP_SCALE,
//     padding: 50,
//     zIndex: 10,
//   },
//   buttonLayout: {
//     flexDirection: 'row',
//   },
// };
