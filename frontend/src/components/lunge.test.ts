import { describe, expect, it } from 'vitest';
import { hitSound, stackSpot, stopShort } from './lunge';

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

describe('stopShort (20.19)', () => {
  it('stops the attacker three quarters of its own length short of the target', () => {
    expect(stopShort(80, 110, 0, -300)).toBeCloseTo(82.5);
    expect(stopShort(80, 110, 300, 0)).toBeCloseTo(60);
    // On a slant the way leaves the card through its nearer side.
    expect(stopShort(80, 110, 300, -300)).toBeCloseTo(60 * Math.SQRT2);
    expect(stopShort(80, 110, 0, 0)).toBe(0);
  });
});
