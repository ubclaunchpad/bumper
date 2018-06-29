import React from 'react';
import { drawPlayer, drawJunk, drawMapHole } from './GameObjects';

const EDGE_BUFFER = 5;
const MAP_SCALE = 12;
const OBJECT_SCALE = 4;

export default class Minimap extends React.Component {
  constructor(props) {
    super(props);
    this.canvas = props.canvas;

    this.mapWidth = props.arena.width / MAP_SCALE;
    this.mapHeight = props.arena.height / MAP_SCALE;
    this.mapX = props.canvas.width - this.mapWidth - EDGE_BUFFER;
    this.mapY = props.canvas.height - this.mapHeight - EDGE_BUFFER;

    this.state = {
      junk: null,
      holes: null,
      players: null,
    };

    this.drawMap = this.drawMap.bind(this);
    this.update = this.update.bind(this);
  }

  update(data) {
    this.state.junk = JSON.parse(JSON.stringify(data.junk));
    this.state.holes = JSON.parse(JSON.stringify(data.holes));
    this.state.players = JSON.parse(JSON.stringify(data.players));

    this.state.players.forEach((p) => {
      p.position.x = (p.position.x / MAP_SCALE) + this.mapX;
      p.position.y = (p.position.y / MAP_SCALE) + this.mapY;
    });
    this.state.holes.forEach((h) => {
      h.position.x = (h.position.x / MAP_SCALE) + this.mapX;
      h.position.y = (h.position.y / MAP_SCALE) + this.mapY;
    });
    this.state.junk.forEach((j) => {
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
        drawPlayer(p, this.canvas, OBJECT_SCALE);
      });
    }
    if (this.junk) {
      this.junk.forEach((j) => {
        drawJunk(j, this.canvas, OBJECT_SCALE);
      });
    }
    if (this.holes) {
      this.holes.forEach((h) => {
        drawMapHole(h, this.canvas);
      });
    }
  }

  // renders() {
  //   const minutes = this.props.finalTime.getMinutes();
  //   const seconds = this.props.finalTime.getSeconds();
  //   const timeString = `${minutes}:${seconds < 10 ? `0${seconds}` : seconds}`;
  //   return (
  //     <div style={styles.backdrop}>
  //       <div style={styles.modal}>
  //         <b>GAME OVER</b>
  //         <div>
  //           <div><b>Time alive:</b> <span>{timeString}</span></div>
  //           <div><b>Points earned:</b> <span>{this.props.finalPoints}</span></div>
  //           <div><b>Final ranking:</b> <span>{this.props.finalRanking}</span></div>
  //         </div>
  //         <button style={styles.restartButton} onClick={() => this.props.onRestart()}>
  //           Restart game
  //         </button>
  //       </div>
  //     </div>
  //   );
  // }
  render() {
    return (
      <div style={styles.mapModal}>
        Minimap
      </div>
    );
  }
}

const styles = {
  mapModal: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    position: 'fixed',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(5,225,255,0.3)',
    borderRadius: 5,
    height: window.innerHeight / MAP_SCALE,
    width: window.innerWidth / MAP_SCALE,
    padding: 50,
    zIndex: 10,
  },
  buttonLayout: {
    flexDirection: 'row',
  },
};
