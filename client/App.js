import React from 'react';

const PLAYER_RADIUS = 25;
const JUNK_SIZE = 15;

const address = process.env.NODE_ENV === 'production'
  ? 'ws://ec2-18-188-53-231.us-east-2.compute.amazonaws.com:9090/connect'
  : 'ws://localhost:9090/connect';

export default class App extends React.Component {
  constructor(props) {
    super(props);
    if (window.WebSocket) {
      this.socket = new WebSocket(address);
      this.socket.onopen = () => {
        this.socket.onmessage = event => this.handleMessage(JSON.parse(event.data));
      };
    } else {
      console.log('websocket not available');
      return;
    }

    this.state = {
      isInitialized: false,
      junk: null,
      holes: null,
      players: null,
      player: null,
    };

    this.handleMessage = this.handleMessage.bind(this);
    this.initializeGame = this.initializeGame.bind(this);
    this.sendKeyPress = this.sendKeyPress.bind(this);
    this.update = this.update.bind(this);
    this.tick = this.tick.bind(this);
    this.draw = this.draw.bind(this);
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);

    window.addEventListener('keydown', this.keyDownHandler);
    window.addEventListener('keyup', this.keyUpHandler);
  }

  async componentDidMount() {
    this.canvas = document.getElementById('ctx');
  }

  sendKeyPress(keyPressed, isPressed) {
    const pressMessage = {
      playerID: 1, // TODO with player ID
      key: keyPressed,
      pressed: isPressed,
    };
    const message = {
      type: 'keyHandler',
      data: JSON.stringify(pressMessage),
    };

    this.socket.send(JSON.stringify(message));
  }

  handleMessage(msg) {
    switch (msg.type) {
      case 'initial':
        this.setState({ player: msg.data });
        break;
      case 'update':
        this.update(msg.data);
        break;
      default:
        console.log(`unknown msg type ${msg.type}`);
        break;
    }
  }

  initializeGame(data) {
    this.setState({
      junk: data.junk,
      holes: data.holes,
      players: data.players,
      isInitialized: true,
    }, () => this.tick());
  }

  update(data) {
    if (!this.state.isInitialized) {
      this.initializeGame(data);
      return;
    }

    this.setState({
      junk: data.junk,
      holes: data.holes,
      players: data.players,
    });
  }

  tick() {
    this.draw();
    // eslint-disable-next-line
    requestAnimationFrame(this.tick);
  }

  drawHoles() {
    this.state.holes.forEach((h) => {
      const ctx = this.canvas.getContext('2d');
      ctx.beginPath();
      ctx.arc(h.position.x, h.position.y, h.radius, 0, Math.PI * 2);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
    });
  }

  drawJunk() {
    this.state.junk.forEach((j) => {
      const ctx = this.canvas.getContext('2d');
      ctx.beginPath();
      ctx.rect(j.position.x, j.position.y, JUNK_SIZE, JUNK_SIZE);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
    });
  }

  drawPlayers() {
    this.state.players.forEach((p) => {
      const ctx = this.canvas.getContext('2d');
      const { x, y } = p.position;
      ctx.beginPath();
      ctx.arc(x, y, PLAYER_RADIUS, 0, Math.PI * 2);
      ctx.fillStyle = p.color;
      ctx.fill();
      ctx.closePath();
    });
  }

  drawPlayerPoints() {
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

  draw() {
    this.drawHoles();
    this.drawJunk();
    this.drawPlayers();
    this.drawPlayerPoints();
  }

  keyDownHandler(e) {
    this.sendKeyPress(e.keyCode, true);
  }

  keyUpHandler(e) {
    this.sendKeyPress(e.keyCode, false);
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

