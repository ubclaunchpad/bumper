export default class Socket {
  constructor(props) {
    if (!window.WebSocket) {
      console.log('WebSocket not available');
      return;
    }

    this.socket = new WebSocket(props.address);
  }
}
