import { Game } from './main.js';

const address = 'localhost:9090';

export class Connection {
    constructor(params) {
        this.address = params.address;

        console.log(`Constructed with address: ${this.address}`);
    }

    async connectPlayer() {
        if (window.WebSocket) {
            const response = await fetch(`http://${this.address}/start`);
            const res = await response.json();
        
            // Address of lobby to connect to
            console.log(res.location);
            this.socket = new WebSocket(`ws://${this.address}/connect`);
            this.socket.onopen = () => {
                this.socket.onmessage = event => handleMessage(JSON.parse(event.data));
            };

            return true;
        }

        return false;
    }

    on(event, handler) {
        console.log(event, handler);
    }
}


// TODO: Should be in Game related functions
function initializeArena(data) {
    Game.player.id = data.playerID;
    Object.assign(Game, {
        arena: { width: data.arenaWidth, height: data.arenaHeight },
        player: data.player,
        timeStarted: new Date(),
    });
}

// TODO: hould be in Game related functions
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