const PLAYER_RADIUS = 25;

export default class Player {
  constructor(props) {
    this.canvas = props.canvas;
    this.position = props.position || { x: window.innerWidth / 2, y: window.innerHeight / 2 };
    this.velocity = { dx: 0, dy: 0 };
    this.theta = props.theta;

    this.mass = props.mass || 10;
    this.color = props.color || '#FF0000';
    this.name = props.name || 'Default name';
    this.alive = true;
	this.walls = {left: false, right: false, top: false, bottom: false};

    this.drawPlayer = this.drawPlayer.bind(this);

    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.keyUpHandler = this.keyUpHandler.bind(this);

    window.addEventListener('keydown', this.keyDownHandler);
    window.addEventListener('keyup', this.keyUpHandler);

    this.rightPressed = false;
    this.leftPressed = false;
    this.upPressed = false;
    this.downPressed = false;
  }

  drawPlayer() {
    const ctx = this.canvas.getContext('2d');
    const { x, y } = this.position;
    ctx.beginPath();
    ctx.arc(x, y, PLAYER_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = '#00FFFF';
    ctx.fill();
    ctx.closePath();

    ctx.beginPath();
    ctx.moveTo(x + (PLAYER_RADIUS * Math.sin(this.theta)), y + (PLAYER_RADIUS * Math.cos(this.theta)));
    ctx.lineTo(x - (PLAYER_RADIUS * Math.sin(this.theta)), y - (PLAYER_RADIUS * Math.cos(this.theta)));
    ctx.strokeStyle = '#000000';
    ctx.strokeWidth = 5;
    ctx.stroke();

    const backCenterX = x - ((PLAYER_RADIUS * Math.sin(this.theta)) / 2);
    const backCenterY = y - ((PLAYER_RADIUS * Math.cos(this.theta)) / 2);
    const backLength = (2.5 * ((PLAYER_RADIUS / 2) / Math.tan(45)));
    ctx.beginPath();
    ctx.moveTo(backCenterX - (backLength * Math.cos(this.theta)), backCenterY + (backLength * Math.sin(this.theta)));
    ctx.lineTo(backCenterX + (backLength * Math.cos(this.theta)), backCenterY - (backLength * Math.sin(this.theta)));
    ctx.strokeStyle = '#0000000';
    ctx.strokeWidth = 5;
    ctx.stroke();
  }

  updatePosition(screen) {
    const controlsVector = { dx: 0, dy: 0 };

    if (this.leftPressed) {
      this.theta = (this.theta + 0.1) % 360;
    }

    if (this.rightPressed) {
      this.theta = (this.theta - 0.1) % 360;
    }

    if (this.downPressed) {
      controlsVector.dy = (0.5 * (PLAYER_RADIUS * Math.cos(this.theta)));
      controlsVector.dx = (0.5 * (PLAYER_RADIUS * Math.sin(this.theta)));
      // Apply resultant vector
      this.position.x += controlsVector.dx;
      this.position.y += controlsVector.dy;
	}

    if (this.upPressed) {
      controlsVector.dy = (0.5 * (PLAYER_RADIUS * Math.cos(this.theta)));
      controlsVector.dx = (0.5 * (PLAYER_RADIUS * Math.sin(this.theta)));
	  // Apply resultant vector
      this.position.x -= controlsVector.dx;
      this.position.y -= controlsVector.dy;
	}

    // Validate position result
    if (this.position.x + PLAYER_RADIUS > (screen.width - 20)) {
      this.position.x = screen.width - 20 - PLAYER_RADIUS;
    } else if (this.position.x - PLAYER_RADIUS < 0) {
      this.position.x = PLAYER_RADIUS;
    }

    if (this.position.y + PLAYER_RADIUS > (screen.height - 20)) {
      this.position.y = screen.height - 20 - PLAYER_RADIUS;
    } else if (this.position.y - PLAYER_RADIUS < 0) {
      this.position.y = PLAYER_RADIUS;
    }
	
	var topWall = this.position.y - PLAYER_RADIUS == 0;
	var bottomWall = this.position.y + PLAYER_RADIUS + 20 == screen.height;
	var leftWall = this.position.x - PLAYER_RADIUS == 0;
	var rightWall = this.position.x + PLAYER_RADIUS + 20 == screen.width;
	
	// Wall detection for the player object:
	if(topWall) { this.walls.top = true; }
	  		else { this.walls.top = false; }	
	if(bottomWall) { this.walls.bottom = true; }
	 		else { this.walls.bottom = false; }
	if(leftWall) { this.walls.left = true; }
	  		else { this.walls.left = false; }
	if(rightWall) { this.walls.right = true; }
	  		else { this.walls.right = false; }

	// To bounce players off the wall they hit:
	if(topWall || bottomWall || leftWall || rightWall) {
		this.theta += Math.PI;
	}
	
  }

  keyDownHandler(e) {
    if (e.keyCode === 39) {
      this.rightPressed = true;
    } else if (e.keyCode === 37) {
      this.leftPressed = true;
    } else if (e.keyCode === 38) {
      this.upPressed = true;
    } else if (e.keyCode === 40) {
      this.downPressed = true;
    }
  }

  keyUpHandler(e) {
    if (e.keyCode === 39) {
      this.rightPressed = false;
    } else if (e.keyCode === 37) {
      this.leftPressed = false;
    } else if (e.keyCode === 38) {
      this.upPressed = false;
    } else if (e.keyCode === 40) {
      this.downPressed = false;
    }
  }
}