import { describe, expect, it } from 'vitest';
import { bezier, feel } from '../feel';
import { turnFrames } from './turn';

describe('turnFrames (20.8)', () => {
  it('grows in, punches, holds and shrinks out over its life', () => {
    const f = turnFrames(false);
    const total = feel.turn.hide + feel.turn.out;
    expect(f.map((k) => k.offset)).toEqual([0, 0.25 / total, 0.4 / total, 1.5 / total, 1]);
    expect(f.map((k) => k.transform)).toEqual(['scale(0.1)', 'scale(1)', 'scale(0.98)', 'scale(0.98)', 'scale(0.1)']);
    expect(f.map((k) => k.opacity)).toEqual([0, 1, 1, 1, 0]);
    for (const k of f) {
      if (k.easing) {
        expect(() => bezier(k.easing ?? '')).not.toThrow();
      }
    }
  });

  it('only fades with less motion', () => {
    expect(turnFrames(true).every((k) => k.transform === 'scale(1)')).toBe(true);
  });

  it('lets playback go on before the banner ends', () => {
    expect(feel.turn.wait).toBeLessThan(feel.turn.hide + feel.turn.out);
  });
});
