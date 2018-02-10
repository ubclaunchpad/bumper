import React from 'react';

const PLAYER_RADIUS = 20;
const JUNK_COUNT = 5;
const JUNK_SIZE = 15;
const HOLE_COUNT = 5;
const HOLE_RADIUS = 25;
const MAX_RADIUS = 50;

const width = window.innerWidth;
const height = window.innerHeight;

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      allCoords: [],
      junkCoords: [],
      holeCoords: [],
      playerCoords: [],
    };
    this.drawObjects = this.drawObjects.bind(this);
  }

  componentDidMount() {
    this.generateJunkCoordinates();
    this.generateHoleCoordinates();
    this.generatePlayerCoordinates();
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
    this.setState({ playerCoords: { x, y } });
  }

  generateCoords(num) {
    let count = num;
    const coords = [];
    while (count > 0) {
      const x = Math.floor(Math.random() * (width - MAX_RADIUS));
      const y = Math.floor(Math.random() * (height - MAX_RADIUS));
      let placed = true;

      // check whether area is available
      for (const p of this.state.allCoords) { //es-lint-disable no-restricted-syntax 
        // could not be placed because of overlap
        if (Math.abs(p.x - x) < MAX_RADIUS || Math.abs(p.y - y) < MAX_RADIUS) {
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
    const canvas = document.getElementById('ctx');
    const ctx = canvas.getContext('2d');
    for (const p of this.state.junkCoords) {
      ctx.beginPath();
      ctx.rect(p.x, p.y, JUNK_SIZE, JUNK_SIZE);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
    }
    for (const p of this.state.holeCoords) {
      ctx.beginPath();
      ctx.arc(p.x, p.y, HOLE_RADIUS, 0, Math.PI * 2);
      ctx.fillStyle = 'white';
      ctx.fill();
      ctx.closePath();
    }

    ctx.beginPath();
    ctx.arc(this.state.playerCoords.x, this.state.playerCoords.y, PLAYER_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = 'green';
    ctx.fill();
    ctx.closePath();
  }
  componentDidUpdate() {
    this.drawObjects();
  }

  render() {
    return (
      <div>
        <canvas id="ctx" style={styles.canvas} width={window.innerWidth} height={window.innerHeight} />
      </div>
    );
  }
}

const styles = {
  canvas: {
    background: '#000',
  },
};
