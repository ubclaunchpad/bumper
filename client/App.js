import React from 'react';

const PLAYER_RADIUS = 25;
const JUNK_COUNT = 10;
const JUNK_SIZE = 15;
const HOLE_COUNT = 10;
const HOLE_RADIUS = 25;
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
      playerX: 200,
      playerY: 200,
      rightPressed: false,
      leftPressed: false,
      upPressed: false,
      downPressed: false,
      allCoords: [],
      junkCoords: [],
      holeCoords: [],
    };

    this.resizeCanvas = this.resizeCanvas.bind(this);
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);
    this.drawObjects = this.drawObjects.bind(this);
    this.tick = this.tick.bind(this);
  }

  componentDidMount() {
    this.generateJunkCoordinates();
    this.generateHoleCoordinates();
    this.generatePlayerCoordinates();

    this.canvas = document.getElementById('ctx');
    window.addEventListener('resize', this.resizeCanvas);
    window.addEventListener('keydown', this.keyDownHandler);
    window.addEventListener('keyup', this.keyUpHandler);
    this.tick();
  }

  generateJunkCoordinates() {
    const newCoords = this.generateCoords(JUNK_COUNT);
    this.setState({ junkCoords: newCoords });
  }

  generateHoleCoordinates() {
    const newCoords = this.generateCoords(HOLE_COUNT);
    this.setState({
      holeCoords: newCoords,
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

  drawObjects() {
    this.drawJunk();
    this.drawHoles();
    this.drawPlayer();
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
    const ctx = this.canvas.getContext('2d');
    for (const p of this.state.holeCoords) {
      ctx.beginPath();
      ctx.arc(p.x, p.y, HOLE_RADIUS, 0, Math.PI * 2);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
    }
  }

  drawPlayer() {
    const ctx = this.canvas.getContext('2d');
    ctx.beginPath();
    ctx.arc(this.state.playerX, this.state.playerY, PLAYER_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = 'green';
    ctx.fill();
    ctx.closePath();
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
    this.drawObjects();

    const speed = 10;
    if (this.state.rightPressed) {
      this.setState(prevState => ({
        playerX: prevState.playerX + speed,
      }));
    }
    if (this.state.leftPressed) {
      this.setState(prevState => ({
        playerX: prevState.playerX - speed,
      }));
    }
    if (this.state.upPressed) {
      this.setState(prevState => ({
        playerY: prevState.playerY - speed,
      }));
    }
    if (this.state.downPressed) {
      this.setState(prevState => ({
        playerY: prevState.playerY + speed,
      }));
    }
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

