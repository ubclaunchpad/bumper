const PLAYER_RADIUS = 25;
const JUNK_SIZE = 15;

export function drawGame(data, canvas) {
  const player = data.players.find(p => p.id === data.player.id);
  if (player.name !== '') {
    // Save a copy of the player's original position.
    const rawPlayerPosition = { x: player.position.x, y: player.position.y };
    // Offset defaults to position.
    const playerOffset = { x: player.position.x, y: player.position.y };

    // If the player's x coord is < half the canvas width we don't need to translate it
    // becuase then the x coord of the screen we're drawing is 0.
    if (player.position.x > canvas.width / 2) {
      // If the player is past half the screen lock them at half till they're close to the rightmost arena edge.
      if ((player.position.x < data.arena.width - (canvas.width / 2))) {
        player.position.x = canvas.width / 2;
        playerOffset.x = canvas.width / 2;
      } else {
        // The player is less than half a canvas width from the rightmost edge. Now let them approach the edge.
        playerOffset.x = player.position.x - (data.arena.width - canvas.width);
        player.position.x -= (data.arena.width - canvas.width);
      }
    }
    // The same logic applies in the Y direction locking the edges of the arena to the edges of the visible canvas.
    if (player.position.y > canvas.height / 2) {
      if ((player.position.y < data.arena.height - (canvas.height / 2))) {
        player.position.y = canvas.height / 2;
        playerOffset.y = canvas.height / 2;
      } else {
        playerOffset.y = player.position.y - (data.arena.height - canvas.height);
        player.position.y -= (data.arena.height - canvas.height);
      }
    }
    const objectXTranslation = playerOffset.x - rawPlayerPosition.x;
    const objectYTranslation = playerOffset.y - rawPlayerPosition.y;

    data.junk.map((j) => {
      const drawableJunk = { ...j };
      drawableJunk.position.x += objectXTranslation;
      drawableJunk.position.y += objectYTranslation;
      drawJunk(drawableJunk, canvas, 1);
      return true;
    });

    data.holes.map((h) => {
      const drawableHole = { ...h };
      drawableHole.position.x += objectXTranslation;
      drawableHole.position.y += objectYTranslation;
      drawHole(drawableHole, canvas);
      return true;
    });

    data.players.map((p) => {
      const drawablePlayer = { ...p };
      if (p.name !== '' && p.id !== data.player.id) {
        drawablePlayer.position.x += objectXTranslation;
        drawablePlayer.position.y += objectYTranslation;
        drawPlayer(drawablePlayer, canvas, 1);
      }
      return true;
    });
  }
}

export function drawJunk(j, canvas, scale) {
  const ctx = canvas.getContext('2d');
  const junkSize = JUNK_SIZE / scale;

  ctx.beginPath();
  ctx.rect(j.position.x - (junkSize / 2), j.position.y - (junkSize / 2), junkSize, junkSize);
  ctx.fillStyle = j.color;
  ctx.fill();
  ctx.closePath();
}

export function drawHole(h, canvas) {
  const ctx = canvas.getContext('2d');
  ctx.beginPath();
  for (let i = 0; i < 720; i += 1) {
    const angle = 0.1 * i;
    const x = h.position.x + ((1 + angle) * Math.cos(angle));
    const y = h.position.y + ((1 + angle) * Math.sin(angle));

    // Find distance between the point (x, y) and the point (h.position.x, h.position.y)
    const x1 = Math.abs(h.position.x - x);
    const y1 = Math.abs(h.position.y - y);
    const distance = Math.hypot(x1, y1);

    // Only draw the line segment if it will correspond to a spiral with the correct radius
    if (distance <= h.radius) {
      ctx.lineTo(x, y);
    }
  }
  ctx.strokeStyle = h.isAlive ? 'white' : 'rgba(255, 255, 255, 0.5)';
  ctx.lineWidth = 1;
  ctx.stroke();
  ctx.closePath();
}

export function drawMapHole(h, canvas, scale) {
  const ctx = canvas.getContext('2d');
  ctx.beginPath();
  ctx.arc(h.position.x, h.position.y, h.radius / scale, 0, 2 * Math.PI);
  ctx.fillStyle = 'rgb(255,225,225)';
  ctx.fill();
  ctx.stroke();
}

export function drawPlayer(p, canvas, scale) {
  const ctx = canvas.getContext('2d');
  const { x, y } = p.position;
  const playerSize = PLAYER_RADIUS / scale;

  if (p.name !== '') {
    // Proportions
    const proportionBackCenter = 3 / 4;
    const proportionWingOuterTop = 4 / 7;
    const proportionWingOuterBottom = 5 / 6;
    const proportionWingOuterDistance = 4 / 5;
    const proportionWingTopInnerDistance = 7 / 10;
    // Constants
    const sinAngle = Math.sin(p.angle);
    const cosAngle = Math.cos(p.angle);
    const playerRadiusSinAngle = playerSize * sinAngle;
    const playerRadiusCosAngle = playerSize * cosAngle;
    const backCenterX = x - (playerRadiusSinAngle * proportionBackCenter); // determines location of the base of the rocket
    const backCenterY = y - (playerRadiusCosAngle * proportionBackCenter);
    const backLength = (playerSize / 2);
    const backLengthSinAngle = backLength * sinAngle;
    const backLengthCosAngle = backLength * cosAngle;
    const wingTopX = x - (playerRadiusSinAngle / 3); // determines location of the top of the wing
    const wingTopY = y - (playerRadiusCosAngle / 3);
    /*
    Start drawing Rocket Chassis, starts from bottom right to the bottom left,
    draw toward the rocket tip then back to the bottom right to complete the shape and fill
    */
    // Coordinates of the Rocket Tip
    const rocketTipX = x + (playerRadiusSinAngle * 1.2);
    const rocketTipY = y + (playerRadiusCosAngle * 1.2);
    // Control Points for Bezier Curve from/toward the Rocket Tip
    const rocketTipModifierRightX = x + (playerSize * Math.sin(p.angle - Math.PI / 4));
    const rocketTipModifierRightY = y + (playerSize * Math.cos(p.angle - Math.PI / 4));
    const rocketTipModifierLeftX = x + (playerSize * Math.sin(p.angle + Math.PI / 4));
    const rocketTipModifierLeftY = y + (playerSize * Math.cos(p.angle + Math.PI / 4));
    // Center-Right Coordinates of Rocket
    const rightCenterX = x + (playerSize * Math.sin(p.angle - Math.PI / 2));
    const rightCenterY = y + (playerSize * Math.cos(p.angle - Math.PI / 2));
    // Center-Left Coordinates of Rocket
    const leftCenterX = x + (playerSize * Math.sin(p.angle + Math.PI / 2));
    const leftCenterY = y + (playerSize * Math.cos(p.angle + Math.PI / 2));
    // Base Coordinates
    const rocketBottomRightX = backCenterX - backLengthCosAngle;
    const rocketBottomRightY = backCenterY + backLengthSinAngle;
    const rocketBottomLeftX = backCenterX + backLengthCosAngle;
    const rocketBottomLeftY = backCenterY - backLengthSinAngle;
    // Draw Rocket Bottom
    ctx.beginPath();
    ctx.moveTo(rocketBottomRightX, rocketBottomRightY); // bottom right side
    ctx.lineTo(rocketBottomLeftX, rocketBottomLeftY); // bottom left side
    // Draw Left Side
    ctx.bezierCurveTo(leftCenterX, leftCenterY, rocketTipModifierLeftX, rocketTipModifierLeftY, rocketTipX, rocketTipY); // chassis left side
    // Draw Right Side
    ctx.bezierCurveTo(rocketTipModifierRightX, rocketTipModifierRightY, rightCenterX, rightCenterY, rocketBottomRightX, rocketBottomRightY); // chassis right side
    ctx.fillStyle = p.color;
    ctx.fill();
    ctx.closePath();
    /*
    Start drawing Rocket Wings, the top of the wing is drawn first, moving toward the base of the rocket and then
    toward the outer part of the wing before going back toward the front side and closing at the top of the wing again.
    */
    // Helper points along the vertical axis of the player model.
    const wingOuterTopX = x - (playerRadiusSinAngle * proportionWingOuterTop); // Point that sets the height level of the top outer part of the wings
    const wingOuterTopY = y - (playerRadiusCosAngle * proportionWingOuterTop);
    const wingOuterBottomX = x - (playerRadiusSinAngle * proportionWingOuterBottom);// Point that sets the height level of the bottom outer part of the wings
    const wingOuterBottomY = y - (playerRadiusCosAngle * proportionWingOuterBottom);
    // Exact points for the right side of the wing
    const wingTopRightX = wingTopX - (playerRadiusCosAngle * proportionWingTopInnerDistance); // inner top right corner
    const wingTopRightY = wingTopY + (playerRadiusSinAngle * proportionWingTopInnerDistance);
    const wingBotRightX = rocketBottomRightX; // inner bottom right corner
    const wingBotRightY = rocketBottomRightY;
    const wingOuterTopRightX = wingOuterTopX - (playerRadiusCosAngle * proportionWingOuterDistance); // outer top right corner
    const wingOuterTopRightY = wingOuterTopY + (playerRadiusSinAngle * proportionWingOuterDistance);
    const wingOuterBottomRightX = wingOuterBottomX - (playerRadiusCosAngle * proportionWingOuterDistance); // outer bottom right corner
    const wingOuterBottomRightY = wingOuterBottomY + (playerRadiusSinAngle * proportionWingOuterDistance);
    // Exact points for the left side of the wing
    const wingTopLeftX = wingTopX + (playerRadiusCosAngle * proportionWingTopInnerDistance); // inner top left corner
    const wingTopLeftY = wingTopY - (playerRadiusSinAngle * proportionWingTopInnerDistance);
    const wingBotLeftX = rocketBottomLeftX; // inner bottom left corner
    const wingBotLeftY = rocketBottomLeftY;
    const wingOuterTopLeftX = wingOuterTopX + (playerRadiusCosAngle * proportionWingOuterDistance); // outer top left corner
    const wingOuterTopLeftY = wingOuterTopY - (playerRadiusSinAngle * proportionWingOuterDistance);
    const wingOuterBottomLeftX = wingOuterBottomX + (playerRadiusCosAngle * proportionWingOuterDistance); // outer bottom left corner
    const wingOuterBottomLeftY = wingOuterBottomY - (playerRadiusSinAngle * proportionWingOuterDistance);
    // Draw Rocket Right Wing
    ctx.beginPath();
    ctx.moveTo(wingTopRightX, wingTopRightY);
    ctx.lineTo(wingBotRightX, wingBotRightY);
    ctx.lineTo(wingOuterBottomRightX, wingOuterBottomRightY);
    ctx.lineTo(wingOuterTopRightX, wingOuterTopRightY);
    ctx.fillStyle = p.color;
    ctx.fill();
    ctx.closePath();
    // Draw Rocket Left Wing
    ctx.beginPath();
    ctx.moveTo(wingTopLeftX, wingTopLeftY);
    ctx.lineTo(wingBotLeftX, wingBotLeftY);
    ctx.lineTo(wingOuterBottomLeftX, wingOuterBottomLeftY);
    ctx.lineTo(wingOuterTopLeftX, wingOuterTopLeftY);
    ctx.fillStyle = p.color;
    ctx.fill();
    ctx.closePath();

    // TODO: Rocket Bottom piece
    // TODO: Rocket Window
  }
}

export function drawWalls(player, arena, canvas) {
  if (player.position) {
    const ctx = canvas.getContext('2d');
    if (player.position.x < (canvas.width / 2)) {
      ctx.beginPath();
      ctx.rect(0, 0, 10, arena.height);
      ctx.fillStyle = 'yellow';
      ctx.fill();
      ctx.closePath();
    }
    if (player.position.x > arena.width - (canvas.width / 2)) {
      ctx.beginPath();
      ctx.rect(canvas.width - 10, 0, 10, arena.height);
      ctx.fillStyle = 'yellow';
      ctx.fill();
      ctx.closePath();
    }
    if (player.position.y < (canvas.height / 2)) {
      ctx.beginPath();
      ctx.rect(0, 0, arena.width, 10);
      ctx.fillStyle = 'yellow';
      ctx.fill();
      ctx.closePath();
    }
    if (player.position.y > arena.height - (canvas.height / 2)) {
      ctx.beginPath();
      ctx.rect(0, canvas.height - 10, arena.width, 10);
      ctx.fillStyle = 'yellow';
      ctx.fill();
      ctx.closePath();
    }
  }
}
