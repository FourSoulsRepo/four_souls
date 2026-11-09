// A held hand card's feel (20.6, ADR 027 "Drag"): it tilts by its own
// motion, as if held in a hand, and springs back to flat when still;
// an "air" loop gets louder the faster it moves. These are the plain
// functions, one frame at a time; dragFeel.ts runs them while a card is
// dragged. Every number lives in feel.drag.
import { feel } from './feel';
import { px } from './uiScale';

// Damped is a value that eases toward a target, and its speed.
export interface Damped {
  value: number;
  velocity: number;
}

export const rest: Damped = { value: 0, velocity: 0 };

// smoothDamp moves d toward target over one frame of dt seconds, as a
// critically damped spring that settles in about `time` seconds; it
// never overshoots the target. maxSpeed caps its speed, in units per
// second. The usual game-engine smooth-damp (Game Programming Gems 4).
export function smoothDamp(d: Damped, target: number, time: number, dt: number, maxSpeed = Infinity): Damped {
  if (dt <= 0) {
    return d;
  }
  const t = Math.max(time, 1e-4);
  const omega = 2 / t;
  const x = omega * dt;
  const decay = 1 / (1 + x + 0.48 * x * x + 0.235 * x * x * x);
  const most = maxSpeed * t;
  const change = Math.min(Math.max(d.value - target, -most), most);
  const to = d.value - change;
  const temp = (d.velocity + omega * change) * dt;
  let velocity = (d.velocity - omega * temp) * decay;
  let value = to + (change + temp) * decay;
  if (target - d.value > 0 === value > target) {
    value = target;
    velocity = 0;
  }
  return { value, velocity };
}

function clamp(v: number, limit: number): number {
  return Math.min(Math.max(v, -limit), limit);
}

// Tilt is a held card's turn, in degrees: pitch about its horizontal
// axis (rotateX), roll about its vertical one (rotateY).
export interface Tilt {
  pitch: Damped;
  roll: Damped;
}

export const flat: Tilt = { pitch: rest, roll: rest };

// tiltStep is the tilt after one frame of dt seconds in which the card
// moved dx, dy px. Each move adds to the tilt, feel.drag.gain degrees
// per unit, up to the limits; the edge it moves toward leans away.
// Then the tilt springs back toward flat over feel.drag.rest.
export function tiltStep(t: Tilt, dx: number, dy: number, dt: number): Tilt {
  const { gain, pitch, roll, rest: back } = feel.drag;
  const unit = px(feel.drag.unit);
  const p = clamp(t.pitch.value - (dy / unit) * gain, pitch);
  const r = clamp(t.roll.value + (dx / unit) * gain, roll);
  return {
    pitch: smoothDamp({ value: p, velocity: t.pitch.velocity }, 0, back, dt),
    roll: smoothDamp({ value: r, velocity: t.roll.velocity }, 0, back, dt),
  };
}

// airStep is the air loop's level after one frame of dt seconds in
// which the card moved dist px. Its target grows with the log of the
// speed, as the reference's: still, it is below 0 and stays silent.
// The speed counts per frame at 60 fps, whatever the frame rate.
export function airStep(a: Damped, dist: number, dt: number): Damped {
  if (dt <= 0) {
    return a;
  }
  const { airScale, airBase, airEase, airRate } = feel.drag;
  const unit = px(feel.drag.unit);
  const perFrame = (dist / unit) * (1 / 60 / dt);
  return smoothDamp(a, Math.log(perFrame * airScale + airBase), airEase, dt, airRate);
}

// airGain is the air loop's mixer gain for a level: 0 to feel.drag.air.
export function airGain(a: Damped): number {
  return Math.min(Math.max(a.value, 0), feel.drag.air);
}
