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
      playerName: '',
      showWelcomeModal: true,
      showGameOverModal: false,
      isInitialized: false,
      junk: null,
      holes: null,
      players: null,
      playerAbsolutePosition: null,
      playerID: null,
      arena: null,
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
      this.socket = new WebSocket(`${address}?name=${inputName}`);
      this.socket.onopen = () => {
        this.socket.onmessage = event => this.handleMessage(JSON.parse(event.data));
      };
      this.socket.onclose = () => {
        this.openGameOverModal();
      };
    } else {
      console.log('websocket not available');
      return;
    }

    this.setState({
      showWelcomeModal: false,
      playerName: inputName,
    });
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
        this.initializeArena(msg.data);
        break;
      case 'update':
        this.update(msg.data);
        break;
      default:
        console.log(`unknown msg type ${msg.type}`);
        break;
    }
  }

  initializeArena(data) {
    this.setState({
      arena: { width: data.arenawidth, height: data.arenaheight },
      playerID: data.playerid,
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
    let playerOffset = null;

    data.players.forEach((player) => {
      if (player.color === this.state.playerID) {
        playerPosition = player.position;
        this.setState({ playerAbsolutePosition: playerPosition });

        player.position = { x: playerPosition.x, y: playerPosition.y };
        playerOffset = { x: playerPosition.x, y: playerPosition.y };
        if (player.position.x > this.canvas.width / 2) {
          if ((player.position.x < this.state.arena.width - (this.canvas.width / 2))) {
            player.position.x = this.canvas.width / 2;
            playerOffset.x = this.canvas.width / 2;
          } else {
            playerOffset.x = player.position.x - (this.state.arena.width - this.canvas.width);
            player.position.x -= (this.state.arena.width - this.canvas.width);
          }
        }
        if (player.position.y > this.canvas.height / 2) {
          if ((player.position.y < this.state.arena.height - (this.canvas.height / 2))) {
            player.position.y = this.canvas.height / 2;
            playerOffset.y = this.canvas.height / 2;
          } else {
            playerOffset.y = player.position.y - (this.state.arena.height - this.canvas.height);
            player.position.y -= (this.state.arena.height - this.canvas.height);
          }
        }
      }
    });

    data.junk.forEach((junk) => {
      junk.position.x -= playerPosition.x;
      junk.position.y -= playerPosition.y;
      junk.position.x += playerOffset.x;
      junk.position.y += playerOffset.y;
    });
    data.holes.forEach((hole) => {
      hole.position.x -= playerPosition.x;
      hole.position.y -= playerPosition.y;
      hole.position.x += playerOffset.x;
      hole.position.y += playerOffset.y;
    });
    data.players.forEach((player) => {
      if (player.color !== this.state.playerID) {
        player.position.x -= playerPosition.x;
        player.position.y -= playerPosition.y;
        player.position.x += playerOffset.x;
        player.position.y += playerOffset.y;
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
      ctx.fillText(`${i + 1}. ${player.name}`, xPos, yPos);
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
      for (let i = 0; i < 720; i += 1) {
        const angle = 0.1 * i;
        const x = h.position.x + (1 + 1 * angle) * Math.cos(angle);
        const y = h.position.y + (1 + 1 * angle) * Math.sin(angle);

        // Find distance between the point (x, y) and the point (h.position.x, h.position.y)
        const x1 = Math.abs(h.position.x - x);
        const y1 = Math.abs(h.position.y - y);
        const distance = Math.sqrt(Math.pow(x1, 2) + Math.pow(y1, 2));

        // Only draw the line segment if it will correspond to a spiral with the correct radius
        if (distance <= h.radius) {
          ctx.lineTo(x, y);
        }
      }
      ctx.strokeStyle = 'white';
      ctx.lineWidth = 1;
      ctx.stroke();
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

      // Constants
      const sinAngle = Math.sin(p.angle);
      const cosAngle = Math.cos(p.angle);

      // const frontCenterX = x + (PLAYER_RADIUS * sinAngle);
      // const frontCenterY = y + (PLAYER_RADIUS * cosAngle);

      const backCenterX = x - ((PLAYER_RADIUS * sinAngle) / 4 * 3);
      const backCenterY = y - ((PLAYER_RADIUS * cosAngle) / 4 * 3);

      const backLength = (PLAYER_RADIUS / 2);

      const wingTopX = x - ((PLAYER_RADIUS * sinAngle) / 2);
      const wingTopY = y - ((PLAYER_RADIUS * cosAngle) / 2);

      // TESTING
      // Circle
      // ctx.beginPath();
      // ctx.arc(x, y, PLAYER_RADIUS, 0, Math.PI * 2);
      // ctx.fillStyle = '#FFFFFF';
      // ctx.fill();
      // ctx.closePath();

      // TESTING
      // Center Line
      // ctx.beginPath();
      // ctx.moveTo(frontCenterX, frontCenterY);
      // ctx.lineTo(x - (PLAYER_RADIUS * sinAngle), y - (PLAYER_RADIUS * cosAngle));
      // ctx.strokeStyle = '#FFFFFF';
      // ctx.strokeWidth = 5;
      // ctx.stroke();

      /*
      Start drawing Rocket Chassis, starts from bottom right to the bottom left,
      draw toward the rocket tip then back to the bottom right to complete the shape and fill
      */
      // Coordinates of the Rocket Tip
      const rocketTipX = x + (PLAYER_RADIUS * sinAngle * 1.2);
      const rocketTipY = y + (PLAYER_RADIUS * cosAngle * 1.2);
      // Control Points for Bezier Curve from/toward the Rocket Tip
      const rocketTipModifierRightX = x + (PLAYER_RADIUS * Math.sin(p.angle - Math.PI / 4));
      const rocketTipModifierRightY = y + (PLAYER_RADIUS * Math.cos(p.angle - Math.PI / 4));
      const rocketTipModifierLeftX = x + (PLAYER_RADIUS * Math.sin(p.angle + Math.PI / 4));
      const rocketTipModifierLeftY = y + (PLAYER_RADIUS * Math.cos(p.angle + Math.PI / 4));
      // Center-Right Coordinates of Rocket
      const rightCenterX = x + (PLAYER_RADIUS * Math.sin(p.angle - Math.PI / 2));
      const rightCenterY = y + (PLAYER_RADIUS * Math.cos(p.angle - Math.PI / 2));
      // Center-Left Coordinates of Rocket
      const leftCenterX = x + (PLAYER_RADIUS * Math.sin(p.angle + Math.PI / 2));
      const leftCenterY = y + (PLAYER_RADIUS * Math.cos(p.angle + Math.PI / 2));
      // Base Coordinates
      const rocketBottomRightX = backCenterX - (backLength * cosAngle);
      const rocketBottomRightY = backCenterY + (backLength * sinAngle);
      const rocketBottomLeftX = backCenterX + (backLength * cosAngle);
      const rocketBottomLeftY = backCenterY - (backLength * sinAngle);
      // Rocket Base
      ctx.beginPath();
      ctx.moveTo(rocketBottomRightX, rocketBottomRightY); // bottom right side
      ctx.lineTo(rocketBottomLeftX, rocketBottomLeftY); // bottom left side
      // Left Side
      ctx.bezierCurveTo(leftCenterX, leftCenterY, rocketTipModifierLeftX, rocketTipModifierLeftY, rocketTipX, rocketTipY); // chassis left side
      // Right Side
      ctx.bezierCurveTo(rocketTipModifierRightX, rocketTipModifierRightY, rightCenterX, rightCenterY, rocketBottomRightX, rocketBottomRightY); // chassis right side
      ctx.fillStyle = p.color;
      ctx.fill();
      ctx.closePath();

      const wingTopRightX = wingTopX - (backLength * cosAngle);
      const wingTopRightY = wingTopY + (backLength * sinAngle);
      const wingBotRightX = backCenterX - (backLength * cosAngle);
      const wingBotRightY = backCenterY + (backLength * sinAngle);
      // TODO: Rocket Right Wing
      ctx.beginPath();
      ctx.moveTo(wingTopRightX, wingTopRightY);
      ctx.lineTo(wingBotRightX, wingBotRightY);
      ctx.strokeStyle = '#FFFFFF';
      ctx.strokeWidth = 5;
      ctx.stroke();
      ctx.closePath();
      // TODO: Rocket Left Wing
      // TODO: Rocket Bottom piece
      // TODO: Rocket Window
    });
  }

  drawWalls() {
    if (this.state.playerAbsolutePosition) {
      if (this.state.playerAbsolutePosition.x < (this.canvas.width / 2)) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(0, 0, 10, this.state.arena.height);
        ctx.fillStyle = 'yellow';
        ctx.fill();
        ctx.closePath();
      }
      if (this.state.playerAbsolutePosition.x > this.state.arena.width - (this.canvas.width / 2)) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(this.canvas.width - 10, 0, 10, this.state.arena.height);
        ctx.fillStyle = 'yellow';
        ctx.fill();
        ctx.closePath();
      }
      if (this.state.playerAbsolutePosition.y < (this.canvas.height / 2)) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(0, 0, this.state.arena.width, 10);
        ctx.fillStyle = 'yellow';
        ctx.fill();
        ctx.closePath();
      }
      if (this.state.playerAbsolutePosition.y > this.state.arena.height - (this.canvas.height / 2)) {
        const ctx = this.canvas.getContext('2d');
        ctx.beginPath();
        ctx.rect(0, this.canvas.height - 10, this.state.arena.width, 10);
        ctx.fillStyle = 'yellow';
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
    return (
      <div style={styles.canvasContainer}>
        <canvas id="ctx" style={styles.canvas} display="inline" width={window.innerWidth - 20} height={window.innerHeight - 20} margin={0} />
        {
          this.state.showWelcomeModal &&
          <WelcomeModal
            name={this.state.playerName}
            onSubmit={e => this.sendSubmitPlayerID(e)}
          />
        }
        {
          this.state.showGameOverModal &&
          <GameOverModal
            data={this.state.gameOverData}
            onRestart={() => this.setState({ showWelcomeModal: true, showGameOverModal: false })}
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
