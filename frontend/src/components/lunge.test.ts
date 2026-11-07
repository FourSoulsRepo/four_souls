import { describe, expect, it } from 'vitest';
import { hitSound, stackSpot } from './lunge';

describe('stackSpot (20.4)', () => {
  it('stacks bursts on one card: middle, above, below, then further out', () => {
    expect(stackSpot([])).toBe(0);
    expect(stackSpot([0])).toBe(-1);
    expect(stackSpot([0, -1])).toBe(1);
    expect(stackSpot([0, -1, 1])).toBe(-2);
    expect(stackSpot([0, -1, 1, -2])).toBe(2);
  });

  it('takes the first spot that freed up', () => {
    expect(stackSpot([-1, 1])).toBe(0);
    expect(stackSpot([0, 1])).toBe(-1);
  });
});

describe('hitSound (20.3)', () => {
  it('sounds by size: attack-hit, attack-hit-2, attack-hit-3', () => {
    expect([1, 2, 3, 4, 7].map(hitSound)).toEqual(['attackHit', 'attackHit2', 'attackHit2', 'attackHit3', 'attackHit3']);
  });
});
