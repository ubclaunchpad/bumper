import React from 'react';
import Player from './components/Player';
import Hole from './components/Hole';
import Junk from './components/Junk';
import { generateRandomColor } from './utils/color';

const PLAYER_RADIUS = 25;
const JUNK_COUNT = 10;
const JUNK_SIZE = 15;
const HOLE_COUNT = 10;
const MAX_DISTANCE_BETWEEN = 50;
const POINTS_PER_JUNK = 100;

const MIN_HOLE_RADIUS = 15;
const MAX_HOLE_RADIUS = 30;
const MIN_HOLE_LIFE = 25;
const MAX_HOLE_LIFE = 75;

const width = window.innerWidth;
const height = window.innerHeight;
const address = process.env.NODE_ENV === 'production'
  ? 'ws://ec2-18-218-38-25.us-east-2.compute.amazonaws.com:9090/connect'
  : 'ws://localhost:9090/connect';

// detect collision
// (x2-x1)^2 + (y1-y2)^2 <= (r1+r2)^2
function areCirclesColliding(p, r1, q, r2) {
  return (((p.x - q.x) ** 2) + ((p.y - q.y) ** 2)) <= ((r1 + r2) ** 2);
}

export default class App extends React.Component {
  constructor(props) {
    super(props);
    if (window.WebSocket) {
      this.socket = new WebSocket(address);
      this.socket.onopen = () => this.initialClientMessage();
      this.socket.onmessage = event => this.handleServerMessage(JSON.parse(event.data));
    } else {
      console.log('websocket not available');
    }

    this.state = {
      allCoords: [], // might need to change this
      junk: [],
      holes: [],
      player: null,
      points: 0,
    };

    this.resizeCanvas = this.resizeCanvas.bind(this);
    this.tick = this.tick.bind(this);
    this.initialClientMessage = this.initialClientMessage.bind(this);
  }

  async componentDidMount() {
    this.canvas = document.getElementById('ctx');
    this.generateJunk();
    this.generatePlayer();
    this.generateHoles();

    window.addEventListener('resize', this.resizeCanvas);
    this.tick();
  }


  handleServerMessage(msg) {
    if (msg.type === 'initial') {
      // add id to player
      // start update interval
      this.state.player.id = msg.id;
      this.setState({ player: this.state.player });

      this.timerID = setInterval(
        () => this.updateClientMessage(),
        17, // 60 Hz
      );
    } else if (msg.type === 'update') {
      this.setState({
        players: msg.players,
      });
    }
  }

  initialClientMessage() {
    if (this.socket.readyState !== 1) return;

    this.socket.send(JSON.stringify({
      type: 'initial',
      id: 1,
      message: 'hello',
    }));
  }

  updateClientMessage() {
    if (this.socket.readyState !== 1 || !this.state.player) return;

    this.socket.send(JSON.stringify({
      type: 'update',
      id: this.state.player.id,
      position: this.state.player.position,
    }));
  }

  generateJunk() {
    const newCoords = this.generateCoords(JUNK_COUNT);
    newCoords.forEach((coord) => {
      const props = {
        position: { x: coord.x, y: coord.y },
        canvas: this.canvas,
      };
      this.state.junk.push(new Junk(props));
    });
    this.setState(this.state);
  }

  generateHoles() {
    const newCoords = this.generateCoords(HOLE_COUNT);
    const newHoles = [];
    newCoords.forEach((coord) => {
      const props = {
        position: { x: coord.x, y: coord.y },
        radius: Math.floor(Math.random() * ((MAX_HOLE_RADIUS - MIN_HOLE_RADIUS) + 1)) + MIN_HOLE_RADIUS,
        lifespan: Math.floor(Math.random() * ((MAX_HOLE_LIFE - MIN_HOLE_LIFE) + 1)) + MIN_HOLE_LIFE,
        canvas: this.canvas,
      };
      const hole = new Hole(props);
      newHoles.push(hole);
    });
    this.setState({
      holes: newHoles,
    });
  }

  // TODO check for collisions
  generatePlayerCoords() {
    const maxWidth = (2 * width) / 3;
    const minWidth = width / 3;
    const maxHeight = (2 * height) / 3;
    const minHeight = height / 3;
    const x = Math.floor(Math.random() * ((maxWidth - minWidth) + 1)) + minWidth;
    const y = Math.floor(Math.random() * ((maxHeight - minHeight) + 1)) + minHeight;
    return {x, y};
  }

  // should appear somewhere in the centre
  generatePlayer() {
    const coords = this.generatePlayerCoords();
    const props = {
      x: coords.x,
      y: coords.y,
      canvas: this.canvas,
      theta: 0,
    };
    this.state.allCoords.push({ x: coords.x, y: coords.y });
    const player = new Player(props);
    this.setState({
      player,
      allCoords: this.state.allCoords,
    });
  }

  generateCoords(num) {
    // make sure object radius isn't outside of canvas
    const maxWidth = width - MAX_DISTANCE_BETWEEN;
    const minWidth = MAX_DISTANCE_BETWEEN;
    const maxHeight = height - MAX_DISTANCE_BETWEEN;
    const minHeight = MAX_DISTANCE_BETWEEN;

    let count = num;
    const coords = [];
    while (count > 0) {
      const x = Math.floor(Math.random() * ((maxWidth - minWidth) + 1)) + minWidth;
      const y = Math.floor(Math.random() * ((maxHeight - minHeight) + 1)) + minHeight;
      const isColliding = this.state.allCoords.some((p) => {
        return areCirclesColliding(p.x, p.y, MAX_DISTANCE_BETWEEN, x, y, MAX_DISTANCE_BETWEEN);
      });

      if (!isColliding) {
        const newAllCoords = this.state.allCoords.push({ x, y });
        this.setState({ allCoords: newAllCoords });
        coords.push({ x, y });
        count -= 1;
      }
    }
    return coords;
  }

  generateNewHoleCoords() {
    // make sure object radius isn't outside of canvas
    const maxWidth = width - MAX_DISTANCE_BETWEEN;
    const minWidth = MAX_DISTANCE_BETWEEN;
    const maxHeight = height - MAX_DISTANCE_BETWEEN;
    const minHeight = MAX_DISTANCE_BETWEEN;

    const coords = { x: 0, y: 0 };
    while (true) {
      coords.x = Math.floor(Math.random() * ((maxWidth - minWidth) + 1)) + minWidth;
      coords.y = Math.floor(Math.random() * ((maxHeight - minHeight) + 1)) + minHeight;

      let isColliding = false;
      this.state.holes.forEach((hole) => {
        const { position } = hole;
        // Check every other
        if (areCirclesColliding(position, MAX_HOLE_RADIUS, coords, MAX_HOLE_RADIUS)) {
          isColliding = true;
        }
      });
      this.state.junk.forEach((junk) => {
        const { position } = junk;
        // Check every junk so we don't swallow them up
        if (areCirclesColliding(position, JUNK_SIZE, coords, MAX_HOLE_RADIUS)) {
          isColliding = true;
        }
      });
      // Check player to junk collisions
      if (this.state.player && !isColliding) {
        if (areCirclesColliding(this.state.player.position, PLAYER_RADIUS * 3, coords, MAX_HOLE_RADIUS)) {
          isColliding = true;
        }
      }

      // Dangerous infite loop?
      if (!isColliding) {
        break;
      }
    }
    return coords;
  }

  drawHoles() {
    this.state.holes.forEach((hole) => {
      if (hole.drawHole() === false) {
        const newCoords = this.generateNewHoleCoords();
        const newRadius = Math.floor(Math.random() * ((MAX_HOLE_RADIUS - MIN_HOLE_RADIUS) + 1)) + MIN_HOLE_RADIUS;
        const newLifespan = Math.floor(Math.random() * ((MAX_HOLE_LIFE - MIN_HOLE_LIFE) + 1)) + MIN_HOLE_LIFE;
        hole.startNewLife(newCoords, newRadius, newLifespan);
        // this.state.holes = this.state.holes.filter(h => h !== hole);
        // this.setState(this.state);
      }
    });
  }

  drawJunk() {
    this.state.junk.forEach(j => j.drawJunk());
  }

  drawPlayers() {
    if (this.state.player) {
      this.state.player.drawPlayer();
    }

    if (!this.state.players) return;

    this.state.players.forEach((p) => {
      if (p.id === this.state.player.id) return;

      const ctx = this.canvas.getContext('2d');
      const { x, y } = p.position;
      ctx.beginPath();
      ctx.arc(x, y, PLAYER_RADIUS, 0, Math.PI * 2);
      ctx.fillStyle = generateRandomColor();
      ctx.fill();
      ctx.closePath();
    });
  }

  resizeCanvas() {
    const ctx = document.getElementById('ctx');
    ctx.width = window.innerWidth - 20;
    ctx.height = window.innerHeight - 20;
    ctx.textAlign = 'center';
    this.updateCanvas();
  }

  tick() {
    this.updateCanvas();
    // check for hole and player collisions
    // TODO check rest of the possible collisions
    this.checkForCollisions();
    // eslint-disable-next-line
    requestAnimationFrame(this.tick);
  }

  drawPlayerPoints() {
    if (!this.state.player) return;

    const ctx = this.canvas.getContext('2d');
    ctx.beginPath();
    const rectHeight = 40;
    const rectWidth = 150;
    const rectX = window.innerWidth - 150;
    const rectY = 0;
    ctx.rect(rectX, rectY, rectWidth, rectHeight);
    ctx.fillStyle = this.state.player.color;
    ctx.fill();
    ctx.font = '16px Lucida Sans Unicode';
    ctx.textAlign = 'center'; 
    ctx.textBaseline = 'middle';
    ctx.fillStyle = '#FFFFFF';
    ctx.fillText(`Points: ${this.state.player.points}`, rectX + (rectWidth / 2) - 10, rectY + (rectHeight / 2) + 2);
    ctx.closePath();
  }

  checkForCollisions() {
    // Check hole to player/junk collisions
    this.state.holes.forEach((hole) => {
      const { position, radius } = hole;
      // Check the player
      if (this.state.player) {
        if (areCirclesColliding(this.state.player.position, PLAYER_RADIUS, position, radius)) {
          this.setState({
            player: null,
          });
        }
      }

      // Check each junk
      this.state.junk.forEach((junk) => {
        if (areCirclesColliding(junk.position, JUNK_SIZE, position, radius)) {
          // Add points for the last bumper player here
          if (junk.lastHitBy !== null) {
            this.state.player.points += POINTS_PER_JUNK;
          }

          this.state.junk = this.state.junk.filter(j => j !== junk);
          this.setState(this.state);
        }
      });
    });

    // Check player to junk collisions
    this.state.junk.forEach((junk) => {
      const { position } = junk;
      if (this.state.player) {
        if (areCirclesColliding(this.state.player.position, PLAYER_RADIUS, position, JUNK_SIZE)) {
          junk.hitBy(this.state.player);
        }
      }
    });
  }

  updateCanvas() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, width, height);
    this.drawJunk();
    this.drawHoles();
    this.drawPlayers();
    this.drawPlayerPoints();
    this.calculateNextState();
  }

  calculateNextState() {
    // TODO check all players
    if (!this.state.player) {
      return;
    }

    this.state.player.updatePosition();
    this.state.junk.forEach(j => j.updatePosition());
  }

  render() {
    return (
      <div style={styles.canvasContainer}>
        <canvas id="ctx" style={styles.canvas} display="inline" width={window.innerWidth - 20} height={window.innerHeight - 20} margin={0} />
      </div>
    );
  }
}

const styles = {
  container: {
    display: 'flex',
  },
  canvas: {
    background: '#000000',
    textAlign: 'center',
  },
  canvasContainer: {
    textAlign: 'center',
  },
};

