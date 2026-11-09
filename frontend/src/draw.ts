// A drawn loot card's path (20.7, ADR 027 "Draw"). An own draw rises
// from the deck, flies to the show spot while it turns face up and
// grows, then goes into the hand: the drawer sees each card before it
// joins the hand. Another player's draw flies from the deck to their
// hand, back only. Times are in seconds at normal speed (feel.draw);
// a pose is sampled at a time, so one clock drives move, turn and size.
import { cubicBezier } from 'motion/react';
import { bezier, feel } from './feel';
import { px } from './uiScale';

export interface Point {
  x: number;
  y: number;
}

// DrawPath is where a drawn card goes: centres in px, and its scale at
// each end, as a part of the size it is drawn at.
export interface DrawPath {
  deck: Point;
  deckScale: number;
  // show is the show spot; an other player's draw has none.
  show?: Point;
  showScale: number;
  hand: Point;
  handScale: number;
}

// DrawPose is the card at one moment: its centre, its scale, and its
// turn in degrees about its vertical axis: 180 shows its back, 0 its face.
export interface DrawPose {
  x: number;
  y: number;
  scale: number;
  turn: number;
}

const ease = (css: string) => cubicBezier(...bezier(css));

// part is how far t is through a part that starts at `at` and lasts
// `time`, from 0 to 1.
function part(t: number, at: number, time: number): number {
  return time <= 0 ? 1 : Math.min(1, Math.max(0, (t - at) / time));
}

function mix(a: number, b: number, u: number): number {
  return a + (b - a) * u;
}

// arc is the point u of the way from a to b on a curve that first
// rises toward a point lift px above a.
function arc(a: Point, b: Point, u: number): Point {
  const c = { x: a.x, y: a.y - px(feel.draw.lift) };
  const v = 1 - u;
  return { x: v * v * a.x + 2 * v * u * c.x + u * u * b.x, y: v * v * a.y + 2 * v * u * c.y + u * u * b.y };
}

// drawTime is a draw's time at normal speed, before its pace.
export function drawTime(own: boolean): number {
  const d = feel.draw;
  return own ? d.show + d.hand : d.other;
}

// ownPose is an own draw at t seconds: deck → show spot (expo-out, as
// the reference's path) while it turns face up and grows, then show
// spot → hand (expo-in: it holds still a moment, then goes).
export function ownPose(t: number, p: DrawPath): DrawPose {
  const d = feel.draw;
  const show = p.show ?? p.hand;
  if (t <= d.show) {
    const at = arc(p.deck, show, ease(d.showCurve)(part(t, 0, d.show)));
    const grow = ease(d.growCurve)(part(t, d.growAt, d.grow));
    const turn = ease(d.turnCurve)(part(t, d.turnAt, d.turn));
    return { ...at, scale: mix(p.deckScale, p.showScale, grow), turn: 180 * (1 - turn) };
  }
  const u = ease(d.handCurve)(part(t, d.show, d.hand));
  return { x: mix(show.x, p.hand.x, u), y: mix(show.y, p.hand.y, u), scale: mix(p.showScale, p.handScale, u), turn: 0 };
}

// otherPose is another player's draw at t seconds: deck → hand on the
// same rising curve, back up, at its size.
export function otherPose(t: number, p: DrawPath): DrawPose {
  const d = feel.draw;
  const u = ease(d.otherCurve)(part(t, 0, d.other));
  return { ...arc(p.deck, p.hand, u), scale: mix(p.deckScale, p.handScale, u), turn: 180 };
}
