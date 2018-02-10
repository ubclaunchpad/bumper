import React from 'react';

const address = 'ws://localhost:9090/connect'


export default class App extends React.Component {
  constructor(props) {
    super(props);

    if(window.WebSocket){
      this.socket = new WebSocket(address);
      this.socket.onmessage = (event) => console.log(event.data);
    }
    else{
      console.log('websocket not available');
    }
    

    this.state = {};
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
