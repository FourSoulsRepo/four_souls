import { describe, expect, it } from 'vitest';
import { bezier, countTimes, curves, elastic, elasticOut, fanPose, feel, hitSize } from './feel';

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
    const parts = [feel.lunge.pull, ...feel.banner.times, ...feel.banner.longTimes];
    for (const p of parts) {
      expect(p).toBeGreaterThanOrEqual(0);
      expect(p).toBeLessThanOrEqual(1);
    }
  });

  it('grows the damage number, holds it and fades it (ADR 027 "Damage and death")', () => {
    const { grow, hold, fade } = feel.burst;
    expect(grow).toBe(1);
    expect(hold).toBe(2);
    expect(fade).toBe(1);
    expect(grow).toBeLessThanOrEqual(hold);
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

describe('hand fan (20.5)', () => {
  const fan = (n: number) => Array.from({ length: n }, (_, i) => fanPose(i, n));

  it('lays 1–3 cards straight', () => {
    for (const n of [1, 2, 3]) {
      expect(fan(n)).toEqual(Array.from({ length: n }, () => ({ angle: 0, drop: 0 })));
    }
  });

  it('spans 40° plus 2° per card from 4 cards, spread evenly', () => {
    for (const [n, span] of [
      [4, 48],
      [5, 50],
      [7, 54],
      [10, 60],
      [12, 64],
    ] as const) {
      const poses = fan(n);
      const angles = poses.map((p) => p.angle);
      expect((angles.at(-1) ?? 0) - (angles[0] ?? 0), `n=${String(n)}`).toBeCloseTo(span);
      const step = span / (n - 1);
      angles.forEach((a, i) => {
        expect(a, `n=${String(n)} i=${String(i)}`).toBeCloseTo(step * (i - (n - 1) / 2));
      });
    }
  });

  it('drops the ends on a parabola, (i + 0.5 − n/2)² / (6n)', () => {
    const unit = feel.hand.drop;
    const drops = (n: number, want: number[]) => {
      fan(n).forEach((p, i) => {
        expect(p.drop / unit, `n=${String(n)} i=${String(i)}`).toBeCloseTo(want.at(i) ?? NaN);
      });
    };
    drops(4, [2.25 / 24, 0.25 / 24, 0.25 / 24, 2.25 / 24]);
    drops(5, [4 / 30, 1 / 30, 0, 1 / 30, 4 / 30]);
    expect((fan(10)[0]?.drop ?? 0) / unit).toBeCloseTo(20.25 / 60);
    expect((fan(12).at(-1)?.drop ?? 0) / unit).toBeCloseTo(30.25 / 72);
  });

  it('is symmetric about the middle', () => {
    for (const n of [4, 5, 8, 11]) {
      const poses = fan(n);
      poses.forEach((p, i) => {
        const q = poses.at(n - 1 - i);
        expect(q?.angle).toBeCloseTo(-p.angle);
        expect(q?.drop).toBeCloseTo(p.drop);
      });
    }
  });
});

describe('zoom (20.5)', () => {
  it('waits 0.4 s for a table card and grows in 0.15 s (ADR 027 principle 2)', () => {
    expect(feel.zoom.delay).toBe(0.4);
    expect(feel.zoom.grow).toBe(0.15);
    expect(() => bezier(feel.zoom.growCurve)).not.toThrow();
    expect(feel.zoom.grace).toBeGreaterThan(0);
  });
});

describe('counters (20.11)', () => {
  it('counts 1 step at once', () => {
    expect(countTimes(1)).toEqual([0]);
    expect(countTimes(-1)).toEqual([0]);
    expect(countTimes(0)).toEqual([]);
  });

  it('counts 5 steps 0.2 s apart', () => {
    const times = countTimes(5);
    expect(times).toHaveLength(5);
    [0, 0.2, 0.4, 0.6, 0.8].forEach((t, i) => {
      expect(times.at(i)).toBeCloseTo(t);
    });
  });

  it('speeds 20 steps up to fit 1 s', () => {
    const times = countTimes(20);
    expect(times).toHaveLength(20);
    times.forEach((t, i) => {
      expect(t).toBeCloseTo(i * 0.05);
    });
    expect(times.at(-1)).toBeLessThan(feel.count.fit);
  });

  it('flips each step in 0.1 s', () => {
    expect(feel.count.flip).toBe(0.1);
    expect(() => bezier(feel.count.flipCurve)).not.toThrow();
  });
});
