import { describe, expect, it } from 'vitest';
import { airGain, airStep, flat, rest, smoothDamp, tiltStep, type Damped, type Tilt } from './drag';
import { feel } from './feel';

const frame = 1 / 60;

// still runs the tilt with no motion for the given seconds.
function still(t: Tilt, seconds: number): Tilt {
  for (let s = 0; s < seconds; s += frame) {
    t = tiltStep(t, 0, 0, frame);
  }
  return t;
}

describe('smoothDamp', () => {
  it('eases toward the target and never overshoots it', () => {
    let d: Damped = { value: 10, velocity: 0 };
    let last = d.value;
    for (let i = 0; i < 120; i++) {
      d = smoothDamp(d, 0, 0.2, frame);
      expect(d.value).toBeGreaterThanOrEqual(0);
      expect(d.value).toBeLessThanOrEqual(last);
      last = d.value;
    }
    expect(d.value).toBeLessThan(0.01);
  });

  it('stands still on a frame of no time', () => {
    const d = { value: 3, velocity: 1 };
    expect(smoothDamp(d, 0, 0.2, 0)).toBe(d);
  });

  it('caps its speed', () => {
    const d = smoothDamp({ value: 0, velocity: 0 }, 10, 0.04, frame, 1);
    expect(d.value).toBeLessThanOrEqual(0.04);
  });
});

describe('tiltStep', () => {
  it('stays flat when the card is still', () => {
    expect(tiltStep(flat, 0, 0, frame)).toEqual({ pitch: rest, roll: rest });
  });

  it('rolls right edge away on a move to the right', () => {
    // One unit right in one frame: 25°, less one frame of the spring.
    const t = tiltStep(flat, feel.drag.unit, 0, frame);
    expect(t.pitch.value).toBe(0);
    expect(t.roll.value).toBeCloseTo(24.69, 2);
  });

  it('pitches the top edge away on a move up', () => {
    const t = tiltStep(flat, 0, -feel.drag.unit, frame);
    expect(t.pitch.value).toBeCloseTo(24.69, 2);
    expect(tiltStep(flat, 0, feel.drag.unit, frame).pitch.value).toBeCloseTo(-24.69, 2);
  });

  it('keeps within ±40° pitch and ±45° roll', () => {
    let t = flat;
    for (let i = 0; i < 30; i++) {
      t = tiltStep(t, -10 * feel.drag.unit, 10 * feel.drag.unit, frame);
      expect(Math.abs(t.pitch.value)).toBeLessThanOrEqual(feel.drag.pitch);
      expect(Math.abs(t.roll.value)).toBeLessThanOrEqual(feel.drag.roll);
    }
    expect(t.pitch.value).toBeLessThan(-35);
    expect(t.roll.value).toBeLessThan(-40);
  });

  it('springs back to flat in about 0.2 s when still', () => {
    const tilted = tiltStep(flat, feel.drag.unit, -feel.drag.unit, frame);
    const soon = still(tilted, 0.2);
    expect(Math.abs(soon.roll.value)).toBeLessThan(0.5 * tilted.roll.value);
    const later = still(tilted, 0.6);
    expect(Math.abs(later.roll.value)).toBeLessThan(0.5);
    expect(Math.abs(later.pitch.value)).toBeLessThan(0.5);
  });

  it('tilts alike at 60 and 120 frames a second', () => {
    const at60 = still(tiltStep(flat, 40, 0, frame), 0.1);
    const half = tiltStep(tiltStep(flat, 20, 0, frame / 2), 20, 0, frame / 2);
    const at120 = (() => {
      let t = half;
      for (let s = 0; s < 0.1; s += frame / 2) {
        t = tiltStep(t, 0, 0, frame / 2);
      }
      return t;
    })();
    expect(Math.abs(at120.roll.value - at60.roll.value)).toBeLessThan(1);
  });
});

describe('air loop', () => {
  it('is silent while the card is still', () => {
    let a = rest;
    for (let i = 0; i < 30; i++) {
      a = airStep(a, 0, frame);
    }
    expect(airGain(a)).toBe(0);
  });

  it('gets louder with speed, up to 50 %', () => {
    const level = (px: number) => {
      let a = rest;
      for (let i = 0; i < 60; i++) {
        a = airStep(a, px, frame);
      }
      return airGain(a);
    };
    expect(level(5)).toBe(0);
    expect(level(30)).toBeGreaterThan(0);
    expect(level(60)).toBeGreaterThan(level(30));
    expect(level(1000)).toBe(feel.drag.air);
  });

  it('counts speed per 60 fps frame, whatever the frame rate', () => {
    let a60 = rest;
    let a120 = rest;
    for (let i = 0; i < 60; i++) {
      a60 = airStep(a60, 40, frame);
      a120 = airStep(airStep(a120, 20, frame / 2), 20, frame / 2);
    }
    expect(airGain(a120)).toBeCloseTo(airGain(a60), 2);
  });
});
