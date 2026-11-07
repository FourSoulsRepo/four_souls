import { describe, expect, it } from 'vitest';
import { hitSound } from './lunge';

describe('hitSound (20.3)', () => {
  it('sounds by size: attack-hit, attack-hit-2, attack-hit-3', () => {
    expect([1, 2, 3, 4, 7].map(hitSound)).toEqual(['attackHit', 'attackHit2', 'attackHit2', 'attackHit3', 'attackHit3']);
  });
});
