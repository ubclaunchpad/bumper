const PLAYER_RADIUS = 25;

// eslint-disable-next-line
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
