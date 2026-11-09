// Where a looked-at deck card went (19.147, ADR 027 "Choosing cards"
// 14): kept on top, it drops onto the pile; put on the bottom, it
// slides out beside the pile, then back under it, hidden where the
// pile covers it. Times are in seconds at normal speed (feel.tuck); a
// pose is sampled at a time, so one clock drives the move and the clip.
import { cubicBezier } from 'motion/react';
import { bezier, feel } from './feel';

// Spot is a card's centre and width on the screen, in px, and its tilt
// in degrees.
export interface Spot {
  x: number;
  y: number;
  w: number;
  rotate: number;
}

// TuckPose is a card at one moment: its spot, the part of its width
// hidden under the pile (0 to 1, from the pile's side), and its opacity.
export interface TuckPose extends Spot {
  hide: number;
  opacity: number;
}

// Side is where a bottomed card slides out: 1 right of the pile, -1 left.
export type Side = 1 | -1;

const ease = (css: string) => cubicBezier(...bezier(css));

function part(t: number, at: number, time: number): number {
  return time <= 0 ? 1 : Math.min(1, Math.max(0, (t - at) / time));
}

const mix = (a: number, b: number, u: number) => a + (b - a) * u;

// tuckSide slides a card out toward the middle of the screen, so it
// stays on screen beside a pile at either edge.
export function tuckSide(pile: Spot, width: number): Side {
  return pile.x < width / 2 ? 1 : -1;
}

// raised is card i of n held over the pile while it is looked at, when
// no back over the deck shows it (the chooser's own look): a small
// fan above the pile, as the backs others see (20.18).
export function raised(pile: Spot, i: number, n: number): Spot {
  const spread = i - (n - 1) / 2;
  return { x: pile.x + spread * pile.w * 0.12, y: pile.y - pile.w * 0.35, w: pile.w * 0.6, rotate: spread * 7 };
}

// tuckTime is how long a card's move lasts; out is the slide out, by
// default feel.tuck.out (a chooser's card flies out over choice.back).
export function tuckTime(bottom: boolean, out: number = feel.tuck.out): number {
  const k = feel.tuck;
  return bottom ? out + k.under : k.drop;
}

// tuckPose is the card t seconds after it leaves `from`.
export function tuckPose(t: number, from: Spot, pile: Spot, bottom: boolean, side: Side, out: number = feel.tuck.out): TuckPose {
  const k = feel.tuck;
  if (!bottom) {
    const u = ease(k.dropCurve)(part(t, 0, k.drop));
    const fade = part(t, k.drop * 0.6, k.drop * 0.4);
    return { x: mix(from.x, pile.x, u), y: mix(from.y, pile.y, u), w: mix(from.w, pile.w, u), rotate: mix(from.rotate, 0, u), hide: 0, opacity: 1 - fade };
  }
  const beside = { x: pile.x + side * pile.w * 1.1, y: pile.y };
  if (t < out) {
    const u = ease(k.outCurve)(part(t, 0, out));
    return { x: mix(from.x, beside.x, u), y: mix(from.y, beside.y, u), w: mix(from.w, pile.w, u), rotate: mix(from.rotate, 0, u), hide: 0, opacity: 1 };
  }
  const u = ease(k.underCurve)(part(t, out, k.under));
  const x = mix(beside.x, pile.x, u);
  const hide = Math.min(1, Math.max(0, 1 - Math.abs(x - pile.x) / Math.max(pile.w, 1)));
  return { x, y: pile.y, w: pile.w, rotate: 0, hide, opacity: hide >= 1 ? 0 : 1 };
}

// tuckClip is a pose's clip-path: the hidden part is on the pile's
// side, the side away from where the card slid out.
export function tuckClip(hide: number, side: Side): string {
  const p = `${(hide * 100).toFixed(2)}%`;
  return side === 1 ? `inset(0 0 0 ${p})` : `inset(0 ${p} 0 0)`;
}
