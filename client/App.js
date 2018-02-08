import React from 'react';

function drawBall(props) {
  const {ctx, x, y, ballRadius} = props;
  ctx.beginPath();
  ctx.arc(x, y, ballRadius, 0, Math.PI*2);
  ctx.fillStyle = '#FFFFFF';
  ctx.fill();
  ctx.closePath();
}

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      playerX: 200,
      playerY: 200,
      rightPressed: false,
      leftPressed: false,
      upPressed: false,
      downPressed: false,
    };
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);
  }

  componentDidMount() {
    window.addEventListener('keydown', this.keyDownHandler);
    window.addEventListener('keyup', this.keyUpHandler);
    this.timerID = setInterval(
      () => this.tick(),
      100,
    );
  }

  tick() {
    const ctx = this.refs.ctx.getContext('2d');
    //console.log('tick');
    
    this.updateCanvas();
  }

  updateCanvas() {
    const canvas = this.refs.ctx;
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    drawBall({
      ctx, x: this.state.playerX, y: this.state.playerY, ballRadius: 50,
    });

    if (this.state.rightPressed) {
      this.setState(prevState => ({
        playerX: prevState.playerX + 5,
      }));
    }
    if (this.state.leftPressed) {
      this.setState(prevState => ({
        playerX: prevState.playerX - 5,
      }));
    }
    if (this.state.upPressed) {
      this.setState(prevState => ({
        playerY: prevState.playerY - 5,
      }));
    }
    if (this.state.downPressed) {
      this.setState(prevState => ({
        playerY: prevState.playerY + 5,
      }));
    }
  }

  keyDownHandler(e) {
    // console.log(e.keyCode);
    if (e.keyCode === 39) {
      this.setState({
        rightPressed: true,
      });
    } else if (e.keyCode === 37) {
      this.setState({
        leftPressed: true,
      });
    } else if (e.keyCode === 38) {
      this.setState({
        upPressed: true,
      });
    } else if (e.keyCode === 40) {
      this.setState({
        downPressed: true,
      });
    }
  }

  keyUpHandler(e) {
    // console.log(e.keyCode);
    if (e.keyCode === 39) {
      this.setState({
        rightPressed: false,
      });
    } else if (e.keyCode === 37) {
      this.setState({
        leftPressed: false,
      });
    } else if (e.keyCode === 38) {
      this.setState({
        upPressed: false,
      });
    } else if (e.keyCode === 40) {
      this.setState({
        downPressed: false,
      });
    }
  }

  render() {
    return (
      <div>
        <canvas ref="ctx" style={styles.canvas} width={window.innerWidth} height={window.innerHeight} />
      </div>
    );
  }
}

const styles = {
  canvas: {
    background: '#000000',
  },
};
