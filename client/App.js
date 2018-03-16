import React from 'react';

const DRAW = 0;
const CLEAR = 1;
const NUM_RANKS = 6;

let printedPlayerRank = false;

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

  /*
   * Performs on operation on the leaderboard
   * Param operation: CLEAR clears the leaderboard from the black canvas
   *                  DRAW draws an updated leaderboard on the canvas
   * Requires: topFivePlayers[], an array with the name strings of the top five players, in order from 1 to 5
   *           playerRank, an int of the rank of the player
   *           currPlayer, a string of the player's name
   *           playerColor, some representation of the player's color
   *           topFivePoints[], an array of ints corresponding to points of the top five players, in order from 1 to 5
   *           currPlayerPoints, an int of the points the player currently has
   */
  leaderboard() {
    const currPlayer = 'Player G';
    const currPlayerPoints = 100;
    const playerRank = 7; // This should be calculated based on an iteration through points
    const topFivePlayers = ['Player A', 'Player B', 'Player C', 'Player D', 'Player E'];
    const topFivePoints = [700, 600, 500, 400, 300, 200, 100];
    const playerColor = '#1702ff';
    const ctx = this.canvas.getContext('2d');
    ctx.beginPath();
    const rectHeight = 130;
    const rectWidth = 170;
    const rectX = window.innerWidth - rectWidth;
    const rectY = 0;
    ctx.rect(rectX, rectY, rectWidth, rectHeight);
    ctx.fillStyle = 'rgba(255,0,0,0.3)';
    ctx.fill();

    // Print leaderboard data:
    // Draw the leaderboard title:
    ctx.font = '16px Lucida Sans Unicode';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillStyle = '#FFFFFF';
    ctx.fillText('Leaderboard', rectX + (rectWidth / 2) - 10, rectY + (rectHeight / 2) - 45);
    
    // Draw the ranks with corresponding player names and points:
    ctx.font = '10px Lucida Sans Unicode';
    let index;
    for (let currRank = 1; currRank < NUM_RANKS; currRank++) {
      index = currRank - 1;
      printedPlayerRank = false;

      if (playerRank === currRank) {  // If player is in the top 5, print its rank in its player color
        printedPlayerRank = true;
        ctx.fillStyle = playerColor;
        ctx.textAlign = 'left'; 
        ctx.fillText(currRank + '. ' + topFivePlayers[index], rectX + (rectWidth / 2) - 80, rectY + (rectHeight / 2) - 25 + 15 * index);
        ctx.textAlign = 'right';
        ctx.fillText(topFivePoints[index], rectX + (rectWidth / 2) + 60, rectY + (rectHeight / 2) - 25 + 15 * index);
        ctx.fillStyle = '#FFFFFF';
      }
      else {   // Else, just print the rank
        ctx.textAlign = 'left';
        ctx.fillText(currRank + '. ' + topFivePlayers[index], rectX + (rectWidth / 2) - 80, rectY + (rectHeight / 2) - 25 + 15 * index);
        ctx.textAlign = 'right';
        ctx.fillText(topFivePoints[index], rectX + (rectWidth / 2) + 60, rectY + (rectHeight / 2) - 25 + 15 * index);
      }
    }

    if (!printedPlayerRank) { // Print the player's rank if it hasn't already been printed
      index = NUM_RANKS - 1;
      ctx.fillStyle = playerColor;
      ctx.textAlign = 'left';
      ctx.fillText(playerRank + '. ' + currPlayer, rectX + (rectWidth / 2) - 80, rectY + (rectHeight / 2) - 25 + 15 * index);
      ctx.textAlign = 'right';
      ctx.fillText(currPlayerPoints, rectX + (rectWidth / 2) + 60, rectY + (rectHeight / 2) - 25 + 15 * index);
      ctx.fillStyle = '#FFFFFF';
    }

    ctx.closePath();
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

  draw() {
    const ctx = this.canvas.getContext('2d');
    ctx.fillStyle = '#000000';
    ctx.fill();
    this.drawHoles();
    this.drawJunk();
    this.drawPlayers();
    this.leaderboard();
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

