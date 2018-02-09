import React from 'react';

const PLAYER_RADIUS = 20;
const JUNK_COUNT = 15;
const JUNK_SIZE = 15;
const HOLE_COUNT = 10;
const HOLE_RADIUS = 25;
const MAX_RADIUS = 50;

var width = window.innerWidth;
var height = window.innerHeight;

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      allCoords: [],
      junkCoords: [],
      holeCoords: [],
      playerCoords: []
    };

    this.drawObjects = this.drawObjects.bind(this);
  }

  componentDidMount() {

    this.generateJunkCoordinates();
    this.generateHoleCoordinates();
    this.generatePlayerCoordinates();

  }

  generateJunkCoordinates() {
    var newCoords = this.generateCoords(JUNK_COUNT);

    this.setState({ junkCoords: newCoords }); 

  }

  generateHoleCoordinates() {
    var newCoords = this.generateCoords(HOLE_COUNT); 
    this.setState({
      holeCoords: newCoords
    });
  }

  // should appear somewhere in the centre 
  generatePlayerCoordinates() {
    var maxWidth = (2*width)/3;
    var minWidth = width/3;
    var maxHeight = (2*height)/3;
    var minHeight = height/3
    var x = Math.floor(Math.random() * (maxWidth - minWidth + 1)) + minWidth;
    var y = Math.floor(Math.random() * (maxHeight - minHeight + 1)) + minHeight;
    this.setState({ playerCoords: { x: x, y: y} });
  }

  generateCoords(num) {
    var coords = [];
    while(num > 0) {
      var x = Math.floor(Math.random()*(width-MAX_RADIUS));
      var y = Math.floor(Math.random()*(height-MAX_RADIUS));
      var placed = true;

      // check whether area is available
      for (var point of this.state.allCoords) {
        // could not be placed because of overlap
        if (Math.abs(point.x-x) < MAX_RADIUS || Math.abs(point.y-y) < MAX_RADIUS) {
          placed = false;
          break;
        }
      }

      if (placed) {
        var newAllCoords = this.state.allCoords.push({ x: x, y: y });
        this.setState({ allCoords: newAllCoords});
        coords.push({ x: x, y: y });
        num = num -1;
      }
    }
    return coords;
  }

  drawObjects() {

    var canvas = document.getElementById("ctx");
    var ctx = canvas.getContext("2d");
    for (var point of this.state.junkCoords) {
      ctx.beginPath();
      ctx.rect(point.x, point.y, JUNK_SIZE, JUNK_SIZE);
      ctx.fillStyle = "white";
      ctx.fill();
      ctx.closePath();
    }
    for (var point of this.state.holeCoords) {
      ctx.beginPath();
      ctx.arc(point.x, point.y, HOLE_RADIUS, 0, Math.PI*2);
      ctx.fillStyle = "white";
      ctx.fill();
      ctx.closePath();
    }

    ctx.beginPath();
    ctx.arc(this.state.playerCoords.x, this.state.playerCoords.y, PLAYER_RADIUS, 0, Math.PI*2);
    ctx.fillStyle = "green";
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
