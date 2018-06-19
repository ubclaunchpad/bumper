export function magnitude(vector) {
  return Math.sqrt((vector.dx * vector.dx) + (vector.dy * vector.dy));
}

export function normalize(vector) {
  const mag = magnitude(vector);
  if (mag > 0) {
    return {
      dx: vector.dx / mag,
      dy: vector.dy / mag,
    };
  }

  return vector;
}
