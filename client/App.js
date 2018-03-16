import React from 'react';
import Player from './components/Player';
import Hole from './components/Hole';
import Junk from './components/Junk';

const DRAW = 0;
const CLEAR = 1;
const NUM_RANKS = 6;

let printedPlayerRank = false;

const PLAYER_RADIUS = 25;
const JUNK_COUNT = 10;
const JUNK_SIZE = 15;
const HOLE_COUNT = 10;
const MAX_DISTANCE_BETWEEN = 50;
const POINTS_PER_JUNK = 100;

const MIN_HOLE_RADIUS = 15;
const MAX_HOLE_RADIUS = 30;
const MIN_HOLE_LIFE = 25;
const MAX_HOLE_LIFE = 75;

const width = window.innerWidth;
const height = window.innerHeight;
const address = process.env.NODE_ENV === 'production'
  ? 'ws://ec2-18-188-53-231.us-east-2.compute.amazonaws.com:9090/connect'
  : 'ws://localhost:9090/connect';

// detect collision
// (x2-x1)^2 + (y1-y2)^2 <= (r1+r2)^2
function areCirclesColliding(p, r1, q, r2) {
  return (((p.x - q.x) ** 2) + ((p.y - q.y) ** 2)) <= ((r1 + r2) ** 2);
}

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
      allCoords: [],
      isInitialized: false,
      junk: null,
      holes: null,
      players: ["Player A","Player B","Player C","Player D","Player E", "Player F", "Player G"],
      currPlayer: "Player G",
      currPlayerPoints: 100,
      playerRank: 7, // This should be calculated based on an iteration through points
      allPoints: [700,600,500,400,300,200,100],
      topFivePlayers: ["Player A","Player B","Player C","Player D","Player E"],
      topFivePoints: [700,600,500,400,300,200,100],
      playerColor: '#1702ff',
    };

    this.handleMessage = this.handleMessage.bind(this);
    this.initializeGame = this.initializeGame.bind(this);
    this.update = this.update.bind(this);
    this.tick = this.tick.bind(this);
    this.draw = this.draw.bind(this);
  }

  async componentDidMount() {
    this.canvas = document.getElementById('ctx');
  }

  handleMessage(msg) {
    switch (msg.type) {
      case 'initial':
        console.log('initial msg received');
        // TODO: set player id
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
      player: data.players[0],
      isInitialized: true,
    }, () => this.tick());
  }
  
  update(data) {
    if (!this.state.isInitialized) {
      this.initializeGame(data);
      return;
    }
    
    // TODO: update objects accordingly
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
  leaderboard(operation) {
    // Draw the rectangle for the leaderboard:
	  const ctx = this.canvas.getContext('2d');
    ctx.beginPath();
    const rectHeight = 130;
    const rectWidth = 170;
    const rectX = window.innerWidth - rectWidth;
    const rectY = 0;
    ctx.rect(rectX, rectY, rectWidth, rectHeight);
    if(operation == DRAW) {
      ctx.fillStyle = 'rgba(255,0,0,0.3)';
    } else if(operation == CLEAR) {
      ctx.fillStyle = '#000000';
    }
    ctx.fill();
    
    // Print leaderboard data:
    if(operation == DRAW) {
      // Draw the leaderboard title:
      ctx.font = '16px Lucida Sans Unicode';
      ctx.textAlign = 'center'; 
      ctx.textBaseline = 'middle';
      ctx.fillStyle = '#FFFFFF';
      ctx.fillText('Leaderboard', rectX + (rectWidth / 2) - 10, rectY + (rectHeight / 2) - 45);
      
      // Draw the ranks with corresponding player names and points:
      ctx.font = '10px Lucida Sans Unicode';
      let index;
      for(let currRank = 1; currRank < NUM_RANKS; currRank++) {
        index = currRank - 1;
        printedPlayerRank = false;

        if(this.state.playerRank == currRank) {  // If player is in the top 5, print its rank in its player color
          printedPlayerRank = true;
          ctx.fillStyle = this.state.playerColor;
          ctx.textAlign = 'left'; 
          ctx.fillText(rank + '. ' + this.state.topFivePlayers[index], rectX + (rectWidth / 2) - 80, rectY + (rectHeight / 2) - 25 + 15 * index);
          ctx.textAlign = 'right';
          ctx.fillText(this.state.topFivePoints[index], rectX + (rectWidth / 2) + 60, rectY + (rectHeight / 2) - 25 + 15 * index);
          ctx.fillStyle = '#FFFFFF';
        }
        else {   // Else, just print the rank
          ctx.textAlign = 'left'; 
          ctx.fillText(currRank + '. ' + this.state.topFivePlayers[index], rectX + (rectWidth / 2) - 80, rectY + (rectHeight / 2) - 25 + 15 * index);
          ctx.textAlign = 'right';
          ctx.fillText(this.state.topFivePoints[index], rectX + (rectWidth / 2) + 60, rectY + (rectHeight / 2) - 25 + 15 * index);
        }
      }
      if(!printedPlayerRank) { // Print the player's rank if it hasn't already been printed
        index = NUM_RANKS - 1;
        ctx.fillStyle = this.state.playerColor;
        ctx.textAlign = 'left'; 
        ctx.fillText(this.state.playerRank + '. ' + this.state.currPlayer, rectX + (rectWidth / 2) - 80, rectY + (rectHeight / 2) - 25 + 15 * index);
        ctx.textAlign = 'right';
        ctx.fillText(this.state.currPlayerPoints, rectX + (rectWidth / 2) + 60, rectY + (rectHeight / 2) - 25 + 15 * index);
        ctx.fillStyle = '#FFFFFF';
      }

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
    this.drawHoles();
    this.drawJunk();
    this.drawPlayers();
    this.leaderboard(CLEAR);
    this.leaderboard(DRAW);
  }

  render() {
    return (
      <div style={styles.canvasContainer}>
        <canvas id="ctx" style={styles.canvas} display="inline" width={window.innerWidth - 20} height={window.innerHeight - 20} margin={0} />
      </div>
    );
  }

  generateJunk() {
    const newCoords = this.generateCoords(JUNK_COUNT);
    newCoords.forEach((coord) => {
      const props = {
        position: { x: coord.x, y: coord.y },
        canvas: this.canvas,
      };
      this.state.junk.push(new Junk(props));
    });
    this.setState(this.state);
  }

  generateHoles() {
    const newCoords = this.generateCoords(HOLE_COUNT);
    const newHoles = [];
    newCoords.forEach((coord) => {
      const props = {
        position: { x: coord.x, y: coord.y },
        radius: Math.floor(Math.random() * ((MAX_HOLE_RADIUS - MIN_HOLE_RADIUS) + 1)) + MIN_HOLE_RADIUS,
        lifespan: Math.floor(Math.random() * ((MAX_HOLE_LIFE - MIN_HOLE_LIFE) + 1)) + MIN_HOLE_LIFE,
        canvas: this.canvas,
      };
      const hole = new Hole(props);
      newHoles.push(hole);
    });
    this.setState({
      holes: newHoles,
    });
  }

  // TODO check for collisions
  generatePlayerCoords() {
    const maxWidth = (2 * width) / 3;
    const minWidth = width / 3;
    const maxHeight = (2 * height) / 3;
    const minHeight = height / 3;
    const x = Math.floor(Math.random() * ((maxWidth - minWidth) + 1)) + minWidth;
    const y = Math.floor(Math.random() * ((maxHeight - minHeight) + 1)) + minHeight;
    return {x, y};
  }

  // should appear somewhere in the centre
  generatePlayer() {
    const coords = this.generatePlayerCoords();
    const props = {
      x: coords.x,
      y: coords.y,
      canvas: this.canvas,
      theta: 0,
    };
    this.state.allCoords.push({ x: coords.x, y: coords.y });
    const player = new Player(props);
    this.setState({
      player,
      allCoords: this.state.allCoords,
    });
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
      const isColliding = this.state.allCoords.some((p) => {
        return areCirclesColliding(p.x, p.y, MAX_DISTANCE_BETWEEN, x, y, MAX_DISTANCE_BETWEEN);
      });

      if (!isColliding) {
        const newAllCoords = this.state.allCoords.push({ x, y });
        this.setState({ allCoords: newAllCoords });
        coords.push({ x, y });
        count -= 1;
      }
    }
    return coords;
  }

  generateNewHoleCoords() {
    // make sure object radius isn't outside of canvas
    const maxWidth = width - MAX_DISTANCE_BETWEEN;
    const minWidth = MAX_DISTANCE_BETWEEN;
    const maxHeight = height - MAX_DISTANCE_BETWEEN;
    const minHeight = MAX_DISTANCE_BETWEEN;

    const coords = { x: 0, y: 0 };
    while (true) {
      coords.x = Math.floor(Math.random() * ((maxWidth - minWidth) + 1)) + minWidth;
      coords.y = Math.floor(Math.random() * ((maxHeight - minHeight) + 1)) + minHeight;

      let isColliding = false;
      this.state.holes.forEach((hole) => {
        const { position } = hole;
        // Check every other
        if (areCirclesColliding(position, MAX_HOLE_RADIUS, coords, MAX_HOLE_RADIUS)) {
          isColliding = true;
        }
      });
      this.state.junk.forEach((junk) => {
        const { position } = junk;
        // Check every junk so we don't swallow them up
        if (areCirclesColliding(position, JUNK_SIZE, coords, MAX_HOLE_RADIUS)) {
          isColliding = true;
        }
      });
      // Check player to junk collisions
      if (this.state.player && !isColliding) {
        if (areCirclesColliding(this.state.player.position, PLAYER_RADIUS * 3, coords, MAX_HOLE_RADIUS)) {
          isColliding = true;
        }
      }

      // Dangerous infite loop?
      if (!isColliding) {
        break;
      }
    }
    return coords;
  }

  
  resizeCanvas() {
    const ctx = document.getElementById('ctx');
    ctx.width = window.innerWidth - 20;
    ctx.height = window.innerHeight - 20;
    ctx.textAlign = 'center';
    this.updateCanvas();
  }

  checkForCollisions() {
    // Check hole to player/junk collisions
    this.state.holes.forEach((hole) => {
      const { position, radius } = hole;
      // Check the player
      if (this.state.player) {
        if (areCirclesColliding(this.state.player.position, PLAYER_RADIUS, position, radius)) {
          this.setState({
            player: null,
          });
        }
      }

      // Check each junk
      this.state.junk.forEach((junk) => {
        if (areCirclesColliding(junk.position, JUNK_SIZE, position, radius)) {
          // Add points for the last bumper player here
          if (junk.lastHitBy !== null) {
            this.state.player.points += POINTS_PER_JUNK;
          }

          this.state.junk = this.state.junk.filter(j => j !== junk);
          this.setState(this.state);
        }
      });
    });

    // Check player to junk collisions
    this.state.junk.forEach((junk) => {
      const { position } = junk;
      if (this.state.player) {
        if (areCirclesColliding(this.state.player.position, PLAYER_RADIUS, position, JUNK_SIZE)) {
          junk.hitBy(this.state.player);
        }
      }
    });
  }

  updateCanvas() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, width, height);
    this.drawJunk();
    this.drawHoles();
    this.drawPlayers();
    this.drawPlayerPoints();
  }

  calculateNextState() {
    // TODO check all players
    if (!this.state.player) {
      return;
    }

    this.state.player.updatePosition();
    this.state.junk.forEach(j => j.updatePosition());
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

