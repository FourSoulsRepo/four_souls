// A held hand card's feel (20.6, ADR 027 "Drag"): while the card is
// dragged, each frame tilts it by how far it moved (drag.ts) and sets
// the air loop by its speed; the bed loop fades in once. On drop both
// loops fade out and the card lies flat. Calm (animation speed "off" or
// reduced motion) keeps the card flat; the loops still play.
import { airGain, airStep, flat, rest, tiltStep, type Damped, type Tilt } from '../drag';
import { feel } from '../feel';
import { playLoop } from '../sounds';

export interface Held {
  // move is the pointer's place now.
  move(x: number, y: number): void;
  // end drops the card.
  end(): void;
}

// holdCard starts the feel of a card picked up at (x, y). tilt gets the
// card's pitch and roll in degrees each frame, and 0, 0 on drop.
export function holdCard(x: number, y: number, calm: boolean, tilt: (pitch: number, roll: number) => void): Held {
  const { bed: bedGain, bedIn, out } = feel.drag;
  const bed = playLoop('dragBed');
  const air = playLoop('dragAir');
  bed.set(bedGain, bedIn);
  let at = { x, y };
  let last = at;
  let turn: Tilt = flat;
  let level: Damped = rest;
  let prev: number | undefined;
  let frame = 0;
  let done = false;
  const tick = (now: number) => {
    if (prev !== undefined) {
      const dt = Math.min((now - prev) / 1000, 0.1);
      const dx = at.x - last.x;
      const dy = at.y - last.y;
      last = at;
      if (!calm) {
        turn = tiltStep(turn, dx, dy, dt);
        tilt(turn.pitch.value, turn.roll.value);
      }
      level = airStep(level, Math.hypot(dx, dy), dt);
      air.set(airGain(level), dt);
    }
    prev = now;
    frame = requestAnimationFrame(tick);
  };
  frame = requestAnimationFrame(tick);
  return {
    move(nx, ny) {
      at = { x: nx, y: ny };
    },
    end() {
      if (done) {
        return;
      }
      done = true;
      cancelAnimationFrame(frame);
      bed.stop(out);
      air.stop(out);
      tilt(0, 0);
    },
  };
}
