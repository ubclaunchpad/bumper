import React from 'react';

import GameOverModal from './components/GameOverModal';
import WelcomeModal from './components/WelcomeModal';
import { drawGame, drawWalls } from './components/GameObjects';
import Leaderboard from './components/Leaderboard';

const address = 'ec2-34-220-30-193.us-west-2.compute.amazonaws.com';

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
    await this.connectPlayer();
    // registerNewTesterEvent();
    // registerTesterUpdateEvent();
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
  async connectPlayer() {
    const response = await fetch(`http://${address}/start`);
    const res = await response.json();

    // Address of lobby to connect to
    console.log(res.location);

    if (window.WebSocket) {
      this.socket = new WebSocket(`ws://${address}/connect`);
      this.socket.onopen = () => {
        this.socket.onmessage = event => this.handleMessage(JSON.parse(event.data));
      };
    }
  }

  // spawn player on submit
  spawnPlayer(name, country) {
    this.sendSpawnMessage(name, country);
    this.state.player.name = name;
    this.state.player.country = country;

    this.setState({
      showWelcomeModal: false,
      showMiniMap: true,
      player: this.state.player,
    });
  }

  sendSpawnMessage(name, country) {
    const spawnMessage = {
      name,
      country,
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

  draw() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    drawGame(this.state, this.canvas);

    // Drawing the walls requires the players position
    const player = this.state.players.find(p => p.id === this.state.player.id);
    if (player && player.name !== '') {
      drawWalls(player, this.state.arena, this.canvas);
    }
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
        <Leaderboard players={this.state.players} />
        <canvas id="ctx" style={styles.canvas} display="inline" width={window.innerWidth - 20} height={window.innerHeight - 20} margin={0} />
        {
          this.state.showWelcomeModal &&
          <WelcomeModal
            name={this.state.player.name}
            country={this.state.player.country}
            onSubmit={(inputName, country) => this.spawnPlayer(inputName, country)}
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
