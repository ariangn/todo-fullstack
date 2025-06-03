// Produces a hex color string from an input string
export function randomColorFromString(str: string): string {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    // simple hash
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  // Pick a mustard-yellow/orange hue
  const hue = (hash % 60) + 30; // between 30° and 90° (yellow/orange)
  return `hsl(${hue}, 70%, 50%)`;
}
