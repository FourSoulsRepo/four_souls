// Counter colours (18.15). A general counter is green, a gold counter
// is gold. The owner picks the colour of each new kind; until then a
// kind shows grey.
const colors = new Map<string, string>([
  ['', '#5cb83a'],
  ['gold', '#e0b43a'],
]);

// counterColor is the fill of a counter badge of this kind; '' is the
// general kind.
export function counterColor(name = ''): string {
  return colors.get(name) ?? '#8a8f98';
}
