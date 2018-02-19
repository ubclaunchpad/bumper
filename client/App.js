import React from 'react';
import Player from './components/Player';
import Hole from './components/Hole';

const JUNK_COUNT = 10;
const JUNK_SIZE = 15;
const HOLE_COUNT = 10;
const MAX_DISTANCE_BETWEEN = 50;

const width = window.innerWidth;
const height = window.innerHeight;
const address = 'ws://localhost:9090/connect';

export default class App extends React.Component {
  constructor(props) {
    super(props);
    if (window.WebSocket) {
      this.socket = new WebSocket(address);
      this.socket.onmessage = event => console.log(event.data);
    } else {
      console.log('websocket not available');
    }

    this.state = {
      rightPressed: false,
      leftPressed: false,
      upPressed: false,
      downPressed: false,
      allCoords: [], // might need to change this
      junkCoords: [],
      holes: [],
    };

    this.resizeCanvas = this.resizeCanvas.bind(this);
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);
    this.tick = this.tick.bind(this);
  }

  componentDidMount() {
    this.canvas = document.getElementById('ctx');
    this.generatePlayerCoordinates();
    this.generateJunkCoordinates();
    this.generateHoles();

    window.addEventListener('resize', this.resizeCanvas);
    window.addEventListener('keydown', this.keyDownHandler);
    window.addEventListener('keyup', this.keyUpHandler);
    this.tick();
  }

  generateJunkCoordinates() {
    const newCoords = this.generateCoords(JUNK_COUNT);
    this.setState({ junkCoords: newCoords });
  }

  generateHoles() {
    const newCoords = this.generateCoords(HOLE_COUNT);
    const newHoles = [];
    newCoords.forEach((coord) => {
      const props = {
        position: { x: coord.x, y: coord.y },
        canvas: this.canvas,
      };
      const hole = new Hole(props);
      newHoles.push(hole);
    });
    this.setState({
      holes: newHoles,
    });
  }

  // should appear somewhere in the centre
  generatePlayerCoordinates() {
    const maxWidth = (2 * width) / 3;
    const minWidth = width / 3;
    const maxHeight = (2 * height) / 3;
    const minHeight = height / 3;
    const x = Math.floor(Math.random() * ((maxWidth - minWidth) + 1)) + minWidth;
    const y = Math.floor(Math.random() * ((maxHeight - minHeight) + 1)) + minHeight;
    const props = {
      x,
      y,
      canvas: this.canvas,
      theta: 0,
    };
    this.state.allCoords.push({ x, y });
    this.setState({
      player: new Player(props),
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
      let placed = true;

      for (const p of this.state.allCoords) { //es-lint-disable no-restricted-syntax 
        // could not be placed because of overlap
        if (Math.abs(p.x - x) < MAX_DISTANCE_BETWEEN && Math.abs(p.y - y) < MAX_DISTANCE_BETWEEN) {
          placed = false;
          break;
        }
      }

      if (placed) {
        const newAllCoords = this.state.allCoords.push({ x, y });
        this.setState({ allCoords: newAllCoords });
        coords.push({ x, y });
        count -= 1;
      }
    }
    return coords;
  }
  
  drawJunk() {
    const ctx = this.canvas.getContext('2d');
    for (const p of this.state.junkCoords) {
      ctx.beginPath();
      ctx.rect(p.x, p.y, JUNK_SIZE, JUNK_SIZE);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
    }
  }

  drawHoles() {
    this.state.holes.forEach(h => h.drawHole());
  }

  drawPlayers() {
    if (this.state.player) {
      this.state.player.drawPlayer();
    }

    // TODO: Draw other players
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
    // eslint-disable-next-line
    requestAnimationFrame(this.tick);
  }

  updateCanvas() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, width, height);
    this.drawPlayers();
    this.drawJunk();
    this.drawHoles();
    this.calculateNextState();
  }

  calculateNextState() {
    const PLAYER_RADIUS = 25;
    this.setState((prevState) => {
      const newState = prevState;
      if (this.state.leftPressed) newState.playerTheta = (prevState.playerTheta + 0.25) % 360;
      if (this.state.rightPressed) newState.playerTheta = (prevState.playerTheta - 0.25) % 360;
      if (this.state.downPressed) {
        newState.playerY = prevState.playerY + (0.5 * (PLAYER_RADIUS * Math.cos(prevState.playerTheta)));
        newState.playerX = prevState.playerX + (0.5 * (PLAYER_RADIUS * Math.sin(prevState.playerTheta)));
      }
      if (this.state.upPressed) {
        newState.playerY = prevState.playerY - (0.5 * (PLAYER_RADIUS * Math.cos(prevState.playerTheta)));
        newState.playerX = prevState.playerX - (0.5 * (PLAYER_RADIUS * Math.sin(prevState.playerTheta)));
      }

      if (newState.playerX + PLAYER_RADIUS > (width - 20)) {
        newState.playerX = width - 20 - PLAYER_RADIUS;
      } else if (newState.playerX - PLAYER_RADIUS < 0) {
        newState.playerX = PLAYER_RADIUS;
      }
      if (newState.playerY + PLAYER_RADIUS > (height - 20)) {
        newState.playerY = height - 20 - PLAYER_RADIUS;
      } else if (newState.playerY - PLAYER_RADIUS < 0) {
        newState.playerY = PLAYER_RADIUS;
      }

      return newState;
    });
  }

  keyDownHandler(e) {
    if (e.keyCode === 39) {
      this.setState({
        rightPressed: true,
      });
    } else if (e.keyCode === 37) {
      this.setState({
        leftPressed: true,
      });
    } else if (e.keyCode === 38) {
      this.setState({
        upPressed: true,
      });
    } else if (e.keyCode === 40) {
      this.setState({
        downPressed: true,
      });
    }
  }

  keyUpHandler(e) {
    if (e.keyCode === 39) {
      this.setState({
        rightPressed: false,
      });
    } else if (e.keyCode === 37) {
      this.setState({
        leftPressed: false,
      });
    } else if (e.keyCode === 38) {
      this.setState({
        upPressed: false,
      });
    } else if (e.keyCode === 40) {
      this.setState({
        downPressed: false,
      });
    }
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

