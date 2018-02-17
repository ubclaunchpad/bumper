export default class Player {
	constructor(props) { 
		this.mass = props.mass;
		this.color = props.color; 
		this.name = props.name;
		this.canvas = props.canvas;
		this.position = props.position;
		this.velocity = [0, 0];
		this.alive = true; 
		
		this.drawPlayer = this.drawPlayer.bind(this);
	 }
		
	drawPlayer() {
		
	}
}