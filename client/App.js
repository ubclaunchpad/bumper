import React from 'react';

import GameOverModal from './components/GameOverModal';
import WelcomeModal from './components/WelcomeModal';
import { drawGame, drawWalls } from './components/GameObjects';
import Minimap from './components/Minimap';
import { registerNewTesterEvent, registerTesterUpdateEvent } from './database/database';


const address = process.env.NODE_ENV === 'production'
  ? 'ws://ec2-54-193-127-203.us-west-1.compute.amazonaws.com/connect'
  : 'ws://localhost:9090/connect';

export default class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      showMiniMap: false,
      showWelcomeModal: true,
      showGameOverModal: false,
      isInitialized: false,
      player: {},
      junk: null,
      holes: null,
      players: null,
      playerAbsolutePosition: null,
      timeStarted: null,
      arena: null,
    };

    this.spawnPlayer = this.spawnPlayer.bind(this);
    this.connectPlayer = this.connectPlayer.bind(this);
    this.sendReconnectMessage = this.sendReconnectMessage.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.initializeArena = this.initializeArena.bind(this);
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
    this.connectPlayer();
    registerNewTesterEvent();
    registerTesterUpdateEvent();
  }

  openGameOverModal() {
    const thisPlayer = this.state.players.find(p => p.id === this.state.player.id);

    this.setState({
      showGameOverModal: true,
      gameOverData: {
        finalTime: new Date((new Date() - this.state.timeStarted)),
        finalPoints: thisPlayer ? thisPlayer.points : 0,
        finalRanking: this.state.player.rank,
      },
    });
  }

  // connect player on load
  connectPlayer() {
    if (window.WebSocket) {
      this.socket = new WebSocket(address);
      this.socket.onopen = () => {
        this.socket.onmessage = event => this.handleMessage(JSON.parse(event.data));
      };
    }
  }

  // spawn player on submit
  spawnPlayer(inputName) {
    this.sendSpawnMessage(inputName);
    this.state.player.name = inputName;
    this.setState({
      showWelcomeModal: false,
      showMiniMap: true,
      player: this.state.player,
    });
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

  sendReconnectMessage() {
    const message = {
      type: 'reconnect',
      data: null,
    };

    if (this.socket.readyState === 1) {
      this.socket.send(JSON.stringify(message));
    }
  }

  sendKeyPress(key, isPressed) {
    const pressMessage = {
      key,
      isPressed,
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
      case 'death':
        this.sendReconnectMessage();
        this.openGameOverModal();
        break;
      case 'update':
        this.update(msg.data);
        break;
      default:
        break;
    }
  }

  initializeArena(data) {
    this.state.player.id = data.playerID;
    this.setState({
      arena: { width: data.arenaWidth, height: data.arenaHeight },
      player: this.state.player,
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

    const thisPlayer = rankedPlayers.find((p, idx) => {
      if (p.id === this.state.player.id) {
        this.state.player.rank = idx + 1;
        return true;
      }

      return false;
    });

    if (thisPlayer) {
      this.setState({
        player: this.state.player,
      });
    }

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
      if (player.name !== '') {
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
      }
    });
    ctx.closePath();
  }

  draw() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    drawGame(this.state, this.canvas);

    this.drawLeaderboard();

    // Drawing the walls requires the players position
    this.state.players.forEach((player) => {
      if (player.name !== '' && player.id === this.state.player.id) {
        drawWalls(player, this.state.arena, this.canvas);
      }
    });
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
          this.state.showMiniMap &&
          <Minimap
            arena={this.state.arena}
            junk={this.state.junk}
            players={this.state.players}
            holes={this.state.holes}
          />
        }
        {
          this.state.showWelcomeModal &&
          <WelcomeModal
            name={this.state.player.name}
            onSubmit={e => this.spawnPlayer(e)}
          />
        }
        {
          this.state.showGameOverModal &&
          <GameOverModal
            {...this.state.gameOverData}
            onRestart={() => {
              this.setState({ showMiniMap: false, showWelcomeModal: true, showGameOverModal: false });
              }
            }
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
