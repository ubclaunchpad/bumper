import React from 'react';

const PLAYER_RADIUS = 5;
const JUNK_COUNT = 10;
const JUNK_RADIUS = 10;
const HOLE_COUNT = 5;
const HOLE_RADIUS = 25;
const MAX_RADIUS = 25;

var width = window.innerWidth;
var height = window.innerHeight;
var randomizationFactor = 50;



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


    //var ctx = canvas.getContext("2d");
    // this.setState({canvasWidth: canvas.width});
    // this.setState({canvasHeight: canvas.height});
    this.generateJunkCoordinates();
    this.generateHoleCoordinates();
    this.generatePlayerCoordinates();
    // add to canvas

  }

  // can be less evenly spaced
  generateJunkCoordinates() {
    var newCoords = this.generateCoords(JUNK_COUNT);

    this.setState({ junkCoords: newCoords }); 

  }

  // should be kind of evenly spaced
  generateHoleCoordinates() {
    var newCoords = this.generateCoords(HOLE_COUNT); 
    this.setState({
      holeCoords: newCoords
    });
  }

  // should appear somewhere in the centre 
  generatePlayerCoordinates() {
    var newCoords = this.generateCoords(1,width ); 
    this.setState({ playerCoords: newCoords }, this.drawObjects);
  }

  drawObjects() {

    var canvas = document.getElementById("ctx");
    var ctx = canvas.getContext("2d");
    for (var point of this.state.junkCoords) {
      ctx.beginPath();
      ctx.arc(point.x, point.y, JUNK_RADIUS, 0, Math.PI*2);
      ctx.fillStyle = "white";
      ctx.fill();
      ctx.closePath();
    }
    for (var point of this.state.holeCoords) {
      ctx.beginPath();
      ctx.arc(point.x, point.y, HOLE_RADIUS, 0, Math.PI*2);
      ctx.fillStyle = "pink";
      ctx.fill();
      ctx.closePath();
    }

    ctx.beginPath();
    ctx.arc(this.state.playerCoords[0].x, this.state.playerCoords[0].y, PLAYER_RADIUS, 0, Math.PI*2);
    ctx.fillStyle = "green";
    ctx.fill();
    ctx.closePath();
    
  }

  generateCoords(num) {
    var coords = [];
    while(num > 0) {
      var x = Math.floor(Math.random()*width);
      var y = Math.floor(Math.random()*height);
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
