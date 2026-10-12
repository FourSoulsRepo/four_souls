// Dice looks (V-19). A dice texture is a 3×3 grid of faces, as in
// Tabletop Simulator: the top row is unused, then 2 4 5 and 1 3 6.

const faces = new Map([
  [2, { col: 0, row: 1 }],
  [4, { col: 1, row: 1 }],
  [5, { col: 2, row: 1 }],
  [1, { col: 0, row: 2 }],
  [3, { col: 1, row: 2 }],
  [6, { col: 2, row: 2 }],
]);

// faceOf finds a value's cell in the texture; null for values a d6 has
// no face for (a modified roll can be 0 or 7), drawn as plain numbers.
export function faceOf(value: number): { col: number; row: number } | null {
  return faces.get(value) ?? null;
}
