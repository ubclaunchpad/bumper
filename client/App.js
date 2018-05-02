import React from 'react';

import GameOverModal from './components/GameOverModal';
import WelcomeModal from './components/WelcomeModal';

const PLAYER_RADIUS = 25;
const JUNK_SIZE = 15;

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
      playerRank: 0,
      timeStarted: null,
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
    const thisPlayer = this.state.players.find(p => p.color === this.state.playerID);

    this.setState({
      showGameOverModal: true,
      gameOverData: {
        finalTime: new Date((new Date() - this.state.timeStarted)),
        finalPoints: thisPlayer ? thisPlayer.points : 0,
        finalRanking: this.state.playerRank,
      },
    });
  }

  sendSubmitPlayerID(inputName) {
    console.log(window.WebSocket);
    if (window.WebSocket) {
      this.socket = new WebSocket(`${address}?name=${inputName}`);
      this.socket.onopen = () => {
        console.log('Connection with server open.');
        this.socket.onmessage = event => this.handleMessage(JSON.parse(event.data));
        this.sendSpawnMessage(inputName);
      };
    }
    this.setState({
      showWelcomeModal: false,
      playerName: inputName,
    }); //  Close Modal
  }

  sendSpawnMessage(inputName) {
    const spawnMessage = {
      name: inputName,
    };
    const message = {
      type: 'spawn',
      data: JSON.stringify(spawnMessage),
    };

    if (this.socket.readyState === 1) {
      this.socket.send(JSON.stringify(message));
    }
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
        console.log('RECEIVED INITIAL MESSAGE from server');
        this.initializeArena(msg.data);
        break;
      case 'death':
        console.log('RECEIVED DEATH MESSAGE from server');
        this.openGameOverModal();
        break;
      case 'update':
        console.log('RECEIVED UPDATE MESSAGE from server');
        console.log(msg);
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
      timeStarted: new Date(),
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
    let playerOffest = null;

    data.players.forEach((player) => {
      if (player.color === this.state.playerID) {
        playerPosition = player.position;
        this.setState({ playerAbsolutePosition: playerPosition });

        player.position = { x: playerPosition.x, y: playerPosition.y };
        playerOffest = { x: playerPosition.x, y: playerPosition.y };
        if (player.position.x > this.canvas.width / 2) {
          if ((player.position.x < this.state.arena.width - (this.canvas.width / 2))) {
            player.position.x = this.canvas.width / 2;
            playerOffest.x = this.canvas.width / 2;
          } else {
            playerOffest.x = player.position.x - (this.state.arena.width - this.canvas.width);
            player.position.x -= (this.state.arena.width - this.canvas.width);
          }
        }
        if (player.position.y > this.canvas.height / 2) {
          if ((player.position.y < this.state.arena.height - (this.canvas.height / 2))) {
            player.position.y = this.canvas.height / 2;
            playerOffest.y = this.canvas.height / 2;
          } else {
            playerOffest.y = player.position.y - (this.state.arena.height - this.canvas.height);
            player.position.y -= (this.state.arena.height - this.canvas.height);
          }
        }
      }
    });

    data.junk.forEach((junk) => {
      junk.position.x -= playerPosition.x;
      junk.position.y -= playerPosition.y;
      junk.position.x += playerOffest.x;
      junk.position.y += playerOffest.y;
    });
    data.holes.forEach((hole) => {
      hole.position.x -= playerPosition.x;
      hole.position.y -= playerPosition.y;
      hole.position.x += playerOffest.x;
      hole.position.y += playerOffest.y;
    });
    data.players.forEach((player) => {
      if (player.color !== this.state.playerID) {
        player.position.x -= playerPosition.x;
        player.position.y -= playerPosition.y;
        player.position.x += playerOffest.x;
        player.position.y += playerOffest.y;
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

    const thisPlayer = rankedPlayers.find(p => p.color === this.state.playerID);
    this.setState({
      playerRank: rankedPlayers.indexOf(thisPlayer) + 1,
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
      for (let i = 0; i < 720; i++) {
        let angle = 0.1 * i;
        let x = h.position.x + (1 + 1 * angle) * Math.cos(angle);
        let y = h.position.y + (1 + 1 * angle) * Math.sin(angle);

        // Find distance between the point (x, y) and the point (h.position.x, h.position.y)
        let x1 = Math.abs(h.position.x - x);
        let y1 = Math.abs(h.position.y - y);
        let distance = Math.sqrt(Math.pow(x1,2)+Math.pow(y1,2));

        // Only draw the line segment if it will correspond to a spiral with the correct radius
        if(distance <= h.radius) { 
          ctx.lineTo(x,y);
        }
      }
      ctx.strokeStyle = h.islive ? 'white' : 'rgba(255, 255, 255, 0.5)';
      ctx.lineWidth=1;
      ctx.stroke();
      ctx.closePath();
    });
  }

  drawJunk() {
    this.state.junk.forEach((j) => {
      const ctx = this.canvas.getContext('2d');
      ctx.beginPath();
      ctx.rect(j.position.x - (JUNK_SIZE / 2), j.position.y - (JUNK_SIZE / 2), JUNK_SIZE, JUNK_SIZE);
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
            {...this.state.gameOverData}
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
