import React from 'react';

import GameOverModal from './components/GameOverModal';
import WelcomeModal from './components/WelcomeModal';

const PLAYER_RADIUS = 25;
const JUNK_SIZE = 15;

// Testing constants:
const FINAL_TIME = 100;
const FINAL_POINTS = 200;
const FINAL_RANKING = 1;

const address = process.env.NODE_ENV === 'production'
  ? 'ws://ec2-54-193-127-203.us-west-1.compute.amazonaws.com/connect'
  : 'ws://localhost:9090/connect';

export default class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      showWelcomeModal: true,
      showGameOverModal: false,
      isInitialized: false,
      junk: null,
      holes: null,
      players: null,
      player: null,
      playerAbsolutePosition: null,
      arena: null,
      center: null,
    };

    this.sendSubmitPlayerID = this.sendSubmitPlayerID.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.initializeGame = this.initializeGame.bind(this);
    this.sendKeyPress = this.sendKeyPress.bind(this);
    this.update = this.update.bind(this);
    this.tick = this.tick.bind(this);
    this.draw = this.draw.bind(this);
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);
  }

  async componentDidMount() {
    this.canvas = document.getElementById('ctx');
  }

  openGameOverModal() {
    this.setState({
      showGameOverModal: true,
      gameOverData: {
        finalTime: FINAL_TIME,
        finalPoints: FINAL_POINTS,
        finalRanking: FINAL_RANKING,
      },
    });
  }

  sendSubmitPlayerID(inputName) {
    if (window.WebSocket) {
      this.socket = new WebSocket(address + "?name=" + inputName);
      this.socket.onopen = () => {
        this.socket.onmessage = event => this.handleMessage(JSON.parse(event.data));
      };
    } else {
      console.log('websocket not available');
      return;
    }

    this.setState({ showWelcomeModal: false }); //  Close Modal
  }

  sendKeyPress(keyPressed, isPressed) {
    const pressMessage = {
      key: keyPressed,
      pressed: isPressed,
    };
    const message = {
      type: 'keyHandler',
      data: JSON.stringify(pressMessage),
    };

    if (this.socket.readyState === 1) {
      this.socket.send(JSON.stringify(message));
    }
  }

  handleMessage(msg) {
    switch (msg.type) {
      case 'initial':
        this.initializePlayerAndArena(msg.data);
        break;
      case 'update':
        this.update(msg.data);
        break;
      default:
        console.log(`unknown msg type ${msg.type}`);
        break;
    }
  }

  initializePlayerAndArena(data) {
    this.setState({
      arena: data.arena, 
      player: data.player,
      center: { x: data.arena.width / 2, y: data.arena.height / 2}
    });
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
      window.addEventListener('keydown', this.keyDownHandler);
      window.addEventListener('keyup', this.keyUpHandler);
      return;
    }

    let playerPosition = null;
    data.players.forEach((player) => {
      // console.log(player.id);
      // console.log(this.state.player);
      if (player.id === this.state.player.id) {
        // console.log('found player');
        playerPosition = player.position;
        this.setState({ playerAbsolutePosition: playerPosition });
        player.position = this.state.center;
      }
    });

    // console.log(playerPosition);

    data.junk.forEach((junk) => {
      junk.position.x = junk.position.x - playerPosition.x;
      junk.position.y = junk.position.y - playerPosition.y;
      junk.position.x = junk.position.x + this.state.center.x;
      junk.position.y = junk.position.y + this.state.center.y;
    });
    data.holes.forEach((hole) => {
      hole.position.x = hole.position.x - playerPosition.x;
      hole.position.y = hole.position.y - playerPosition.y;
      hole.position.x = hole.position.x + this.state.center.x;
      hole.position.y = hole.position.y + this.state.center.y;
    });
    data.players.forEach((player) => {
      if (player.id !== this.state.player.id) {
        player.position.x = player.position.x - playerPosition.x;
        player.position.y = player.position.y - playerPosition.y;
        player.position.x = player.position.x + this.state.center.x;
        player.position.y = player.position.y + this.state.center.y;
      }
    });
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
   * DRAW draws an updated leaderboard on the canvas
   * Requires: this.state.players an array with the players
   * // TODO identify current player
   */

  drawLeaderboard() {
    const rankedPlayers = this.state.players.sort((a, b) => {
      if (b.points < a.points) return -1;
      if (b.points > a.points) return 1;
      if (a.color < b.color) return -1; // sort by color on ties
      if (a.color > b.color) return 1;
      return 0;
    });

    const ctx = this.canvas.getContext('2d');
    ctx.beginPath();
    const rectHeight = 130;
    const rectWidth = 170;
    const rectX = window.innerWidth - rectWidth;
    const rectY = 0;
    let xPos;
    let yPos;
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
    rankedPlayers.forEach((player, i) => {
      ctx.fillStyle = player.color;
      ctx.textAlign = 'left';
      xPos = rectX + (rectWidth / 2) - 80;
      yPos = rectY + (rectHeight / 2) - 25 + 15 * i;
      ctx.fillText(`${i + 1}. Player ${player.color}`, xPos, yPos);
      ctx.textAlign = 'right';
      xPos = rectX + (rectWidth / 2) + 60;
      yPos = rectY + (rectHeight / 2) - 25 + 15 * i;
      ctx.fillText(player.points, xPos, yPos);
      ctx.fillStyle = '#FFFFFF';
    });
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
      ctx.fillStyle = j.color;
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

      ctx.beginPath();
      ctx.moveTo(x + (PLAYER_RADIUS * Math.sin(p.angle)), y + (PLAYER_RADIUS * Math.cos(p.angle)));
      ctx.lineTo(x - (PLAYER_RADIUS * Math.sin(p.angle)), y - (PLAYER_RADIUS * Math.cos(p.angle)));
      ctx.strokeStyle = '#000000';
      ctx.strokeWidth = 5;
      ctx.stroke();

      const backCenterX = x - ((PLAYER_RADIUS * Math.sin(p.angle)) / 2);
      const backCenterY = y - ((PLAYER_RADIUS * Math.cos(p.angle)) / 2);
      const backLength = (2.5 * ((PLAYER_RADIUS / 2) / Math.tan(45)));
      ctx.beginPath();
      ctx.moveTo(backCenterX - (backLength * Math.cos(p.angle)), backCenterY + (backLength * Math.sin(p.angle)));
      ctx.lineTo(backCenterX + (backLength * Math.cos(p.angle)), backCenterY - (backLength * Math.sin(p.angle)));
      ctx.strokeStyle = '#0000000';
      ctx.strokeWidth = 5;
      ctx.stroke();
    });
  }

  drawWalls() {
    if (this.state.playerAbsolutePosition) {
      if (this.state.playerAbsolutePosition.x < this.state.center.x) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(0, 0, this.state.center.x - this.state.playerAbsolutePosition.x, this.state.arena.height);
        ctx.fillStyle = 'white';
        ctx.fill();
        ctx.closePath();
      }
      if (this.state.playerAbsolutePosition.x > this.state.arena.width - this.state.center.x) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(this.state.arena.width, 0, this.state.center.x - this.state.playerAbsolutePosition.x, this.state.arena.height);
        ctx.fillStyle = 'white';
        ctx.fill();
        ctx.closePath();
      }
      if (this.state.playerAbsolutePosition.y < this.state.center.y) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(0, 0, this.state.arena.width, this.state.center.y - this.state.playerAbsolutePosition.y);
        ctx.fillStyle = 'white';
        ctx.fill();
        ctx.closePath();
      }
      if (this.state.playerAbsolutePosition.y > this.state.arena.height - this.state.center.y) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(0, this.state.arena.height, this.state.arena.width, this.state.center.y - this.state.playerAbsolutePosition.y);
        ctx.fillStyle = 'white';
        ctx.fill();
        ctx.closePath();
      }
    }
  }

  draw() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    this.drawHoles();
    this.drawJunk();
    this.drawPlayers();
    this.drawLeaderboard();
    this.drawWalls();
  }

  keyDownHandler(e) {
    this.sendKeyPress(e.keyCode, true);
  }

  keyUpHandler(e) {
    this.sendKeyPress(e.keyCode, false);
  }

  render() {
    if (this.state.showGameOverModal) {
      return <GameOverModal data={this.state.gameOverData} />;
    }

    return (
      <div style={styles.canvasContainer}>
        <canvas id="ctx" style={styles.canvas} display="inline" width={window.innerWidth - 20} height={window.innerHeight - 20} margin={0} />
        {
          this.state.showWelcomeModal &&
          <WelcomeModal
            onSubmit={e => this.sendSubmitPlayerID(e)}
          />
        }
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