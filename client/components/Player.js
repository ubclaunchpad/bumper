export default class Player {
	constructor(props) { 
		this.mass = props.mass || 10;
		this.color = props.color || '#FF0000'; // Red 
		this.name = props.name || 'Default name';
		this.canvas = props.canvas;
		this.position = props.position || { x: (window.innerWidth / 2), y: (window.innerHeight) / 2) }; 
		this.velocity = { dx: 0, dy: 0 };
		this.alive = true;  
		
		this.drawPlayer = this.drawPlayer.bind(this);
	 }
		
	drawPlayer() {
		
	}
}
