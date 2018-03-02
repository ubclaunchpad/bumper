// eslint-disable-next-line
export function generateRandomColor() {
  let c = '';
  while (c.length < 6) {
    c += (Math.random()).toString(16).substr(-6).substr(-1);
  }
  return `#${c}`;
}