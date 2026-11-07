import { describe, expect, it } from 'vitest';
import { bezier, curves, elastic, elasticOut, feel, hitSize } from './feel';

describe('curves', () => {
  it('parses every curve string', () => {
    for (const css of Object.values(curves)) {
      const b = bezier(css);
      expect(b).toHaveLength(4);
      expect(b.every(Number.isFinite)).toBe(true);
    }
  });

  it('parses every easing in the times', () => {
    const { windupCurve, strikeCurve, settleCurve } = feel.lunge;
    for (const css of [windupCurve, strikeCurve, settleCurve, feel.fly.curve, ...elasticOut.map((s) => s.easing)]) {
      expect(() => bezier(css)).not.toThrow();
    }
  });

  it('reads CSS keywords and cubic-bezier', () => {
    expect(bezier('ease-in-out')).toEqual([0.42, 0, 0.58, 1]);
    expect(bezier(curves.expoOut)).toEqual([0.16, 1, 0.3, 1]);
  });

  it('refuses what CSS would not take', () => {
    for (const css of ['', 'easeInOut', 'cubic-bezier(1, 2, 3)', 'cubic-bezier(1.2, 0, 0.5, 1)', 'cubic-bezier(0.1, , 0.5, 1)', 'cubic-bezier(a, 0, 0.5, 1)']) {
      expect(() => bezier(css), css).toThrow();
    }
  });
});

describe('elastic-out', () => {
  it('runs from 0 to 1 and overshoots on the way', () => {
    expect(elastic(0)).toBe(0);
    expect(elastic(1)).toBe(1);
    expect(Math.max(...elasticOut.map((s) => s.value))).toBeGreaterThan(1.2);
  });

  it('is a keyframe list from offset 0 to 1', () => {
    expect(elasticOut[0]).toMatchObject({ offset: 0, value: 0 });
    expect(elasticOut.at(-1)).toMatchObject({ offset: 1, value: 1 });
    const offsets = elasticOut.map((s) => s.offset);
    expect(offsets).toEqual([...offsets].sort((a, b) => a - b));
  });
});

describe('feel', () => {
  it('keeps each part of a motion inside its time', () => {
    const parts = [feel.lunge.pull, feel.burst.pop, feel.burst.fade, ...feel.banner.times, ...feel.banner.longTimes];
    for (const p of parts) {
      expect(p).toBeGreaterThanOrEqual(0);
      expect(p).toBeLessThanOrEqual(1);
    }
    expect(feel.burst.pop).toBeLessThan(feel.burst.fade);
  });

  it('strikes fast and settles slow (ADR 027 "Attack")', () => {
    const { windup, strike, settle, strikeCurve, settleCurve } = feel.lunge;
    expect(strike).toBe(0.12);
    expect(strikeCurve).toBe(curves.linear);
    expect(settleCurve).toBe(curves.expoOut);
    expect(windup).toBeGreaterThan(strike);
    expect(settle).toBeGreaterThan(windup);
  });
});

describe('hit sizes (20.3)', () => {
  it('picks the size by damage: 1, 2–3, 4 and more', () => {
    expect([0, 1, 2, 3, 4, 5, 12].map(hitSize)).toEqual([1, 1, 2, 2, 3, 3, 3]);
  });

  it('grows with the size, with the shake strengths of the reference', () => {
    expect(feel.hit.push).toEqual([0.1, 0.5, 1]);
    expect(feel.hit.burst[0]).toBe(1);
    expect(feel.hit.burst).toEqual([...feel.hit.burst].sort((a, b) => a - b));
    expect(() => bezier(feel.hit.shakeCurve)).not.toThrow();
  });
});
