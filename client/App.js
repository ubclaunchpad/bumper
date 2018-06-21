import React from 'react';

import GameOverModal from './components/GameOverModal';
import WelcomeModal from './components/WelcomeModal';
import { drawPlayer, drawHole, drawJunk } from './components/GameObjects';
import Minimap from './components/Minimap';

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


    const mapProps = {
      canvas: this.canvas,
      arena: { width: data.arenaWidth, height: data.arenaHeight },
    };
    this.Minimap = new Minimap(mapProps);
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

    // The minimap requires the update data before it is translated
    this.Minimap.update(data);

    let playerPosition = null;
    let playerOffset = null;
    data.players.forEach((player) => {
      if (player.name !== '' && player.id === this.state.player.id) {
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

    if (playerPosition != null) {
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
        if (player.name !== '' && player.id !== this.state.player.id) {
          player.position.x -= playerPosition.x;
          player.position.y -= playerPosition.y;
          player.position.x += playerOffset.x;
          player.position.y += playerOffset.y;
        }
      });
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
    this.Minimap.drawMap();

    this.state.holes.forEach((h) => {
      drawHole(h, this.canvas);
    });
    this.state.junk.forEach((j) => {
      drawJunk(j, this.canvas, 1);
    });
    this.state.players.forEach((p) => {
      drawPlayer(p, this.canvas, 1);
    });

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
            name={this.state.player.name}
            onSubmit={e => this.spawnPlayer(e)}
          />
        }
        {
          this.state.showGameOverModal &&
          <GameOverModal
            {...this.state.gameOverData}
            onRestart={() => {
              this.setState({ showWelcomeModal: true, showGameOverModal: false });
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
