import { Game } from './main.js';

const address = 'localhost:9090';

function initializeArena(data) {
    Game.player.id = data.playerID;
    Object.assign(Game, {
        arena: { width: data.arenaWidth, height: data.arenaHeight },
        player: data.player,
        timeStarted: new Date(),
    });
}

function update(data) {
    console.log(Game);
}

function handleMessage(msg) {
    switch (msg.type) {
      case 'initial':
        initializeArena(msg.data);
        break;
      case 'update':
        update(msg.data);
        break;
      default:
        break;
    }
}

export default async function connectPlayer() {
    const response = await fetch(`http://${address}/start`);
    const res = await response.json();

    // Address of lobby to connect to
    console.log(res.location);

    if (window.WebSocket) {
        let socket = new WebSocket(`ws://${address}/connect`);
        socket.onopen = () => {
            socket.onmessage = event => handleMessage(JSON.parse(event.data));
        };
    }
}