
import { Connection, connectPlayer} from './connection.js';
export const Game = {
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
    arena: {
        width: null,
        height: null,
    },
}

//pixi 
let type = "WebGL"
if(!PIXI.utils.isWebGLSupported()){
  type = "canvas"
}
PIXI.utils.sayHello(type)

//Create a Pixi Application
let app = new PIXI.Application( { 
    width: 256,         // default: 800
    height: 256,        // default: 600
    antialias: true,    // default: false
    transparent: false, // default: false
    resolution: 1       // default: 1
    }
);
app.renderer.view.style.position = "absolute";
app.renderer.view.style.display = "block";
app.renderer.autoResize = true;
app.renderer.resize(window.innerWidth, window.innerHeight);
//Add the canvas that Pixi automatically created for you to the HTML document
document.body.appendChild(app.view);

// let texture = PIXI.utils.TextureCache["images/cat.png"];//cache a texture into WebGL for sprite loading
// let sprite = new PIXI.Sprite(texture);//assign sprite object from cached texture
PIXI.loader
    .add([{
            url: "https://cdn2.iconfinder.com/data/icons/outline-signs/350/spaceship-512.png",
            onComplete: function () {},
            crossOrigin: true,
          }
    ])
    .on("progress", loadProgressHandler)
    //list all files you want to load in an array inside a single add method or chain them
    .load(setup); //call setup when loading is finished

function loadProgressHandler(loader, resource) {
    console.log("loading: " + resource.url)
    console.log("progress: " + loader.progress + "%");
}

function setup() {
    //run when loader finish image loading
    let texture = PIXI.utils.TextureCache["https://cdn2.iconfinder.com/data/icons/outline-signs/350/spaceship-512.png"];
    let rocket = new PIXI.Sprite(texture);
    rocket.anchor.set(0.5,0.5);
    rocket.x = 96;
    rocket.y = 96;
    rocket.scale.set(0.2, 0.2); //scale of original png object
    rocket.rotation = 0.5; //radians
    app.stage.addChild(rocket);
}

window.onload = async () => {
    console.log('Starting...');

    const socket = new Connection({
        address: 'localhost:9090',
    });

    socket.connectPlayer();
};