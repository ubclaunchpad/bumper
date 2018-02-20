import React from 'react';
import Hole from './components/Hole';

const PLAYER_RADIUS = 25;
const JUNK_COUNT = 10;
const JUNK_SIZE = 15;
const HOLE_COUNT = 10;
const MAX_DISTANCE_BETWEEN = 50;

const width = window.innerWidth;
const height = window.innerHeight;
const address = 'ws://localhost:9090/connect';

function drawPlayer(props) {
  const {
    ctx, x, y, ballRadius, theta,
  } = props;

  // don't draw player if dead
  // TODO remove when player objects are introduced
  if (x === null || y === null) {
    return;
  }
  ctx.beginPath();
  ctx.arc(x, y, ballRadius, 0, Math.PI * 2);
  ctx.fillStyle = '#00FFFF';
  ctx.fill();
  ctx.closePath();

  // Draw the flag
  ctx.beginPath();
  ctx.moveTo(x + (ballRadius * Math.sin(theta)), y + (ballRadius * Math.cos(theta)));
  ctx.lineTo(x - (ballRadius * Math.sin(theta)), y - (ballRadius * Math.cos(theta)));
  ctx.strokeStyle = '#000000';
  ctx.strokeWidth = 5;
  ctx.stroke();

  const backCenterX = x - ((ballRadius * Math.sin(theta)) / 2);
  const backCenterY = y - ((ballRadius * Math.cos(theta)) / 2);
  const backLength = (2.5 * ((ballRadius / 2) / Math.tan(45)));
  ctx.beginPath();
  ctx.moveTo(backCenterX - (backLength * Math.cos(theta)), backCenterY + (backLength * Math.sin(theta)));
  ctx.lineTo(backCenterX + (backLength * Math.cos(theta)), backCenterY - (backLength * Math.sin(theta)));
  ctx.strokeStyle = '#0000000';
  ctx.strokeWidth = 5;
  ctx.stroke();
}

function areCirclesColliding(x1, y1, r1, x2, y2, r2) {
  return (((x1 - x2) ** 2) + ((y1 - y2) ** 2)) <= ((r1 + r2) ** 2);
}

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
      playerX: 200,
      playerY: 200,
      playerTheta: 0,
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
    this.drawObjects = this.drawObjects.bind(this);
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
    this.setState({ playerX: x, playerY: y });
    this.state.allCoords.push({ x, y });
    this.setState({ allCoords: this.state.allCoords });
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

      // check whether area is available
      for (const p of this.state.allCoords) { //es-lint-disable no-restricted-syntax 
        // could not be placed because of overlap
        if (areCirclesColliding(p.x, p.y, MAX_DISTANCE_BETWEEN, x, y, MAX_DISTANCE_BETWEEN)) {
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

  drawObjects() {
    this.drawJunk();
    this.drawHoles();
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
    for (const h of this.state.holes) {
      h.drawHole();
    }
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
    // check for hole and player collistions
    // TODO check rest
    this.checkForCollisions();
    // eslint-disable-next-line
    requestAnimationFrame(this.tick);
  }

  checkForCollisions() {
    this.state.holes.forEach((hole) => {
      const { position, radius } = hole;
      // detect collision
      // (x2-x1)^2 + (y1-y2)^2 <= (r1+r2)^2
      if (areCirclesColliding(this.state.playerX, this.state.playerY, PLAYER_RADIUS, position.x, position.y, radius)) {
        console.log("ur dead");
        this.setState({
          playerX: null,
          playerY: null,
        });
      }
    });
  }

  

  updateCanvas() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, width, height);
    drawPlayer({
      ctx,
      x: this.state.playerX,
      y: this.state.playerY,
      ballRadius: PLAYER_RADIUS,
      theta: this.state.playerTheta,
    });
    this.drawObjects();

    this.calculateNextState();
  }

  calculateNextState() {
    // player is dead don't render
    // TODO remove when player objects introduced
    if (this.state.playerX === null || this.state.playerY === null) {
      return;
    }
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
  }

};

