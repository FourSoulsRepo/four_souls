import { describe, expect, it } from 'vitest';
import { faceOf } from './dice';

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
