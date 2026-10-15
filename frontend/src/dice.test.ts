import { describe, expect, it } from 'vitest';
import { faceOf, tumbleFaces } from './dice';

describe('faceOf', () => {
  it('finds each face of a d6 in the texture', () => {
    const cells = [1, 2, 3, 4, 5, 6].map((v) => faceOf(v));
    expect(cells).toEqual([
      { col: 0, row: 2 },
      { col: 0, row: 1 },
      { col: 1, row: 2 },
      { col: 1, row: 1 },
      { col: 2, row: 1 },
      { col: 2, row: 2 },
    ]);
    expect(new Set(cells.map((c) => JSON.stringify(c))).size).toBe(6);
  });

  it('has no face for modified values', () => {
    expect(faceOf(0)).toBeNull();
    expect(faceOf(7)).toBeNull();
  });
});

describe('tumbleFaces', () => {
  it('shows two other faces, then lands on the value (15.62)', () => {
    for (let v = 1; v <= 6; v++) {
      const faces = tumbleFaces(v);
      expect(faces.at(-1)).toBe(v);
      expect(new Set(faces).size).toBe(3);
      expect(faces.every((f) => f >= 1 && f <= 6)).toBe(true);
    }
  });
});
