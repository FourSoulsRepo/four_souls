// Aiming (15.25): what an aiming card points at, and its arrow's way
// (20.23). Pure, except zonesAt, onScreen, targetElement and the
// pointer, which read the page.
import { Purpose, type Choice } from './actions';
import { feel } from './feel';
import type { Intent, OptionView } from './protocol';
import { px } from './uiScale';

// zonesAt lists the zones under a point, topmost first: data-zone of
// cards ("obj-N") and piles, data-stack of stack rows. A player is aimed
// at by the character card (data-player, "player-N"), never by the area
// around it (15.33), so the areas' own zones are left out. The stack
// panel and the large played card are zones only for effects (15.74):
// a click there is a click on empty space.
const effectOnly = /^(player-|played-|stack$)/;

export function zonesAt(x: number, y: number): string[] {
  const out: string[] = [];
  for (const el of document.elementsFromPoint(x, y)) {
    const zone = el.getAttribute('data-zone');
    const stack = el.getAttribute('data-stack');
    const player = el.getAttribute('data-player');
    if (stack) {
      out.push(`stack-${stack}`);
    }
    if (zone && !effectOnly.test(zone)) {
      out.push(zone);
    }
    if (player) {
      out.push(`player-${player}`);
    }
  }
  return out;
}

// targetAt is the first of the targets the zones show, topmost first:
// a character card as an object wins over it as its player.
export function targetAt(zones: string[], targets: OptionView[]): OptionView | undefined {
  for (const zone of zones) {
    // The monster deck is an attack's target (20.23).
    const deck = zone === 'deck-monster' && targets.find((o) => o.deck);
    if (deck) {
      return deck;
    }
    const m = /^(obj|player|stack)-(\d+)$/.exec(zone);
    if (!m) {
      continue;
    }
    const n = Number(m[2]);
    const t = targets.find((o) => {
      switch (m[1]) {
        case 'obj':
          return !o.is_player && o.object === n;
        case 'player':
          return !!o.is_player && o.player === n;
        default:
          return o.stack_id === n;
      }
    });
    if (t) {
      return t;
    }
  }
  return undefined;
}

// AimClick is what a click while aiming does (15.33): pick a target;
// cancel on empty space (T-06); nothing on the source itself; or, on a
// card, pile, or stack item that is no target, say why (T-05).
export type AimClick = { kind: 'pick'; target: OptionView } | { kind: 'cancel' } | { kind: 'wrong' } | { kind: 'none' };

// aimClick decides a click by the zones under it, topmost first. The
// source does nothing, unless it is a target itself (an item that may
// pick any item, as The D20): then it picks itself. Owner, 2026-10-14.
export function aimClick(zones: string[], source: number, targets: OptionView[]): AimClick {
  const src = `obj-${String(source)}`;
  if (zones.find((z) => z.startsWith('obj-')) === src) {
    const self = targetAt([src], targets);
    return self ? { kind: 'pick', target: self } : { kind: 'none' };
  }
  const target = targetAt(zones, targets);
  if (target) {
    return { kind: 'pick', target };
  }
  return zones.length === 0 ? { kind: 'cancel' } : { kind: 'wrong' };
}

// PromptAim is a target prompt answered by aiming (15.34, T-13): the
// card already placed aims at the prompt's options.
export interface PromptAim {
  source: number;
  choice: Choice;
  targets: OptionView[];
}

// promptAim is the aim a target prompt opens, when the card that asks
// is known: the viewer's own play or use. Others keep the dialog.
export function promptAim(activating: number | undefined, choice: Choice | null): PromptAim | null {
  if (activating === undefined || choice?.purpose !== Purpose.target) {
    return null;
  }
  const targets = choice.options.flatMap((o) => (o.target ? [o.target] : []));
  return targets.length > 0 ? { source: activating, choice, targets } : null;
}

// cancelOption is the prompt's option that cancels the play: the one
// that points at nothing.
export function cancelOption(choice: Choice): Intent | undefined {
  return choice.options.find((o) => !o.target)?.intent;
}

// zoneOf is the page selector of a target, as zonesAt reads it.
function zoneOf(t: OptionView): string {
  if (t.deck) {
    return '[data-zone="deck-monster"]';
  }
  if (t.is_player) {
    return `[data-player="${String(t.player)}"]`;
  }
  if (t.stack_id) {
    return `[data-stack="${String(t.stack_id)}"]`;
  }
  return `[data-zone="obj-${String(t.object ?? 0)}"]`;
}

// onScreen: every target and the source can be clicked on the table;
// otherwise the dialog lists the options instead (15.34).
export function onScreen(a: PromptAim): boolean {
  return (
    document.querySelector(`[data-zone="obj-${String(a.source)}"]`) !== null &&
    a.targets.every((t) => document.querySelector(zoneOf(t)) !== null)
  );
}

// aimed: the target is one of the targets.
export function aimed(targets: OptionView[] | null, test: (o: OptionView) => boolean): boolean {
  return targets?.some(test) ?? false;
}

export interface Point {
  x: number;
  y: number;
}

// Pose is a place on an aiming arrow: where, and which way it runs
// there, in degrees from the x axis.
export interface Pose {
  x: number;
  y: number;
  angle: number;
}

// Arc is an aiming arrow's way (20.23, ADR 027 "Target arrow"): a
// curve from its source to the pointer. length is the curve's length;
// at(s) the pose s px along it, kept on the curve; end the pose at the
// pointer, turned along the curve's last stretch.
export interface Arc {
  length: number;
  at: (s: number) => Pose;
  end: Pose;
}

// arcLift is how far the middle of an arrow from a point to a point
// stands off the straight line, up the screen (feel.aim): lift, or
// liftShare of a short line, times how much of the lift lies across
// the line (|dx| / length × lean, at most 1). A line straight up or
// down gets none, so the bow fades out there and never flips side.
export function arcLift(from: Point, to: Point): number {
  const { liftShare, lean } = feel.aim;
  const lift = px(feel.aim.lift);
  const dx = to.x - from.x;
  const length = Math.hypot(dx, to.y - from.y);
  if (length === 0) {
    return 0;
  }
  return Math.min(lift, length * liftShare) * Math.min(1, (Math.abs(dx) / length) * lean);
}

// segments is how many straight pieces measure an arc.
const segments = 32;

// arc is the way of an aiming arrow: a quadratic curve whose middle
// stands arcLift off the line, on its upper side.
export function arc(from: Point, to: Point): Arc {
  const dx = to.x - from.x;
  const dy = to.y - from.y;
  const chord = Math.hypot(dx, dy);
  const h = arcLift(from, to);
  // The upper side's normal; (dy, −dx) points up for a line going right.
  const side = dx < 0 ? -1 : 1;
  const nx = chord > 0 ? (side * dy) / chord : 0;
  const ny = chord > 0 ? (-side * dx) / chord : 0;
  // A quadratic's middle is half way to its control point.
  const c = { x: (from.x + to.x) / 2 + 2 * h * nx, y: (from.y + to.y) / 2 + 2 * h * ny };
  const point = (t: number): Point => {
    const u = 1 - t;
    return { x: u * u * from.x + 2 * u * t * c.x + t * t * to.x, y: u * u * from.y + 2 * u * t * c.y + t * t * to.y };
  };
  const points = Array.from({ length: segments + 1 }, (_, i) => point(i / segments));
  const lengths = [0];
  points.forEach((b, i) => {
    const a = points.at(i - 1);
    if (i > 0 && a) {
      lengths.push((lengths.at(-1) ?? 0) + Math.hypot(b.x - a.x, b.y - a.y));
    }
  });
  const length = lengths.at(-1) ?? 0;
  const degrees = (x: number, y: number) => (Math.atan2(y, x) * 180) / Math.PI;
  const at = (s: number): Pose => {
    const d = Math.min(Math.max(s, 0), length);
    let i = 1;
    while (i < segments && (lengths.at(i) ?? 0) < d) {
      i++;
    }
    const a = points.at(i - 1) ?? from;
    const b = points.at(i) ?? to;
    const l0 = lengths.at(i - 1) ?? 0;
    const piece = (lengths.at(i) ?? 0) - l0;
    const k = piece > 0 ? (d - l0) / piece : 0;
    return { x: a.x + (b.x - a.x) * k, y: a.y + (b.y - a.y) * k, angle: degrees(b.x - a.x, b.y - a.y) };
  };
  // The head turns along the curve's end: from the control point on.
  const ex = to.x - c.x;
  const ey = to.y - c.y;
  const end = { x: to.x, y: to.y, angle: ex === 0 && ey === 0 ? degrees(dx, dy) : degrees(ex, ey) };
  return { length, at, end };
}

// dashOpacity is how much a dash shows s px along the chain, which ends
// chain px from the source, by the head (feel.aim): the tail's fades in
// as (s / step)^tailFade, the one by the head fades out over its last
// step as (left / step)^headFade; past the chain's end, none.
export function dashOpacity(s: number, chain: number): number {
  const { headFade, tailFade } = feel.aim;
  const step = px(feel.aim.step);
  const left = chain - s;
  if (s <= 0 || left <= 0) {
    return 0;
  }
  return Math.min(1, Math.pow(s / step, tailFade), Math.pow(left / step, headFade));
}

// chainOf is how long the dashed chain of an arc is: feel.aim.head
// short of the pointer, under the head's back.
export function chainOf(a: Arc): number {
  return Math.max(0, a.length - px(feel.aim.head));
}

// dashCount is how many dashes the chain holds: one per step begun.
export function dashCount(chain: number): number {
  return chain > 0 ? Math.ceil(chain / px(feel.aim.step)) : 0;
}

// DashFrame is a dash's pose at one point of its march, relative to the
// arrow's source.
export interface DashFrame {
  offset: number;
  transform: string;
  opacity: number;
}

// dashPose is dash i's pose part u (0 to 1) of the way through one
// march step: its middle (i + u) steps along the arc, turned with it,
// relative to the source.
export function dashPose(a: Arc, from: Point, i: number, u: number): DashFrame {
  const s = (i + u) * px(feel.aim.step);
  const p = a.at(s);
  const r = (n: number) => String(Math.round(n * 100) / 100);
  return {
    offset: u,
    transform: `translate(${r(p.x - from.x)}px, ${r(p.y - from.y)}px) rotate(${r(p.angle)}deg)`,
    opacity: Math.round(dashOpacity(s, chainOf(a)) * 1000) / 1000,
  };
}

// dashFrames are dash i's keyframes through one march step: it moves
// one step toward the head along the arc, sampled at feel.aim.stops
// points, fading as dashOpacity says. Repeated, the chain marches with
// no seam: each dash ends where the next one starts.
export function dashFrames(a: Arc, from: Point, i: number): DashFrame[] {
  const { stops } = feel.aim;
  return Array.from({ length: stops + 1 }, (_, k) => dashPose(a, from, i, k / stops));
}

// marchTime is one march step's time in ms at an animation speed's
// time scale: step px at march px a second.
export function marchTime(scale: number): number {
  return (feel.aim.step / feel.aim.march) * 1000 * scale;
}

// pointer is the last mouse place on the page, for an arrow drawn at
// once when an aim starts (20.23); null before any mouse moved, or once
// a touch or pen took over.
let pointer: Point | null = null;

if (typeof window !== 'undefined') {
  const track = (e: PointerEvent) => {
    pointer = e.pointerType === 'mouse' ? { x: e.clientX, y: e.clientY } : null;
  };
  window.addEventListener('pointermove', track, { capture: true, passive: true });
  window.addEventListener('pointerdown', track, { capture: true, passive: true });
}

// lastPointer is the last mouse place, if a mouse leads.
export function lastPointer(): Point | null {
  return pointer;
}

// arrowShows is true while the own arrow is drawn (20.23): a card under
// it zooms only once the arrow rests on it (useZoom).
let shows = false;

export function arrowShows(): boolean {
  return shows;
}

export function setArrowShows(on: boolean): void {
  shows = on;
}

// targetElement is a target's element on the page, if drawn.
export function targetElement(t: OptionView): Element | null {
  return document.querySelector(zoneOf(t));
}
