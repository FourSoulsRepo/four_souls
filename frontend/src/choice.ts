// Choosing cards (20.18, ADR 027 "Choosing cards"): the row the cards
// lie in, and their motion. A look at a deck's top cards deals them
// from that deck: each rises toward a point above the deck, flies to
// its place and turns face up on the way. Other card choices grow at
// their places. After the pick the cards fly back to the deck. Times
// are in seconds at normal speed (feel.choice); a pose is sampled at a
// time, so one clock drives move, turn and size.
import { cubicBezier } from 'motion/react';
import { PromptKind, Purpose } from './actions';
import { bezier, feel } from './feel';
import { t } from './i18n';
import type { GameView, Prompt } from './protocol';
import type { Sound } from './sounds';

export interface Row {
  // scale is the cards' size, as a part of a full choice card.
  scale: number;
  // gap is the space between two cards, in card widths.
  gap: number;
  // wrap: more cards than one row holds; they wrap to a second row.
  wrap: boolean;
}

// choiceRow is the row of n cards (ADR 027 "Choosing cards" 2, 12.4):
// one row up to feel.choice.row cards, smaller from 5 on.
export function choiceRow(n: number): Row {
  const { sizes, row } = feel.choice;
  const step = sizes.find((s) => n <= s.upTo) ?? sizes.at(-1) ?? { scale: 1, gap: 0.25 };
  return { scale: step.scale, gap: step.gap, wrap: n > row };
}

// dealDelays is when each of n dealt cards leaves its deck, in seconds:
// the first at once, the second feel.choice.first later, each next one
// feel.choice.gap after the one before.
export function dealDelays(n: number): number[] {
  const { first, gap } = feel.choice;
  return Array.from({ length: n }, (_, i) => (i === 0 ? 0 : first + gap * (i - 1)));
}

export interface Point {
  x: number;
  y: number;
}

// DealPose is a dealt card at one moment: its offset from its place in
// px, its scale, and its turn in degrees about its vertical axis: 180
// shows its back, 0 its face.
export interface DealPose {
  x: number;
  y: number;
  scale: number;
  turn: number;
}

const ease = (css: string) => cubicBezier(...bezier(css));

function part(t: number, at: number, time: number): number {
  return time <= 0 ? 1 : Math.min(1, Math.max(0, (t - at) / time));
}

// dealPose is a card dealt from a deck at offset `from` (the deck's
// centre less the card's, in px) and scale `scale`, t seconds after it
// leaves: on a curve that first rises toward feel.choice.lift px above
// the deck, to its place, expo-out, while it turns face up.
export function dealPose(t: number, from: Point, scale: number): DealPose {
  const c = feel.choice;
  const u = ease(c.dealCurve)(part(t, 0, c.deal));
  const v = 1 - u;
  const lift = { x: from.x, y: from.y - c.lift };
  const turn = ease(c.turnCurve)(part(t, c.turnAt, c.turn));
  return {
    x: v * v * from.x + 2 * v * u * lift.x,
    y: v * v * from.y + 2 * v * u * lift.y,
    scale: scale + (1 - scale) * u,
    turn: 180 * (1 - turn),
  };
}

// backDelays is when each card flies back to its deck after the pick,
// in seconds, by its place in `order`: the first at once, the rest
// feel.choice.backFirst later and feel.choice.backGap apart.
export function backDelays(n: number): number[] {
  const { backFirst, backGap } = feel.choice;
  return Array.from({ length: n }, (_, i) => (i === 0 ? 0 : backFirst + backGap * (i - 1)));
}

// Start-of-game choices (20.22, ADR 027 "Start of game"): the choice
// opens as a screen of its own; a click picks a card and Confirm sends
// it. Every player with such a choice makes it at the same time.

// confirmFirst: a pick of this purpose waits for Confirm before it is
// sent. A start-of-game pick lasts the whole game, so a misclick costs
// much (16.3). The draft pick (14.11) can use it too.
export function confirmFirst(purpose: number): boolean {
  return purpose === Purpose.startingItem;
}

// Toggle is a click on a card before Confirm: the card picked now, if
// any, and the sound the click plays.
export interface Toggle {
  picked: number | null;
  sound: Sound;
}

// pickToggle is a click on option i while picked is picked: it picks
// i, or unpicks it when it is the picked one (16.4).
export function pickToggle(picked: number | null, i: number): Toggle {
  return picked === i ? { picked: null, sound: 'choiceUnpick' } : { picked: i, sound: 'choicePick' };
}

// choosers are all who may answer the prompt now: its players on
// start-of-game choices (20.22), else its player.
export function choosers(w: Prompt): number[] {
  return w.players && w.players.length > 0 ? w.players : [w.player];
}

// startChoice: the prompt is a start-of-game choice (Eden's).
export function startChoice(w: Prompt): boolean {
  return w.kind === PromptKind.choose && w.purpose === Purpose.startingItem;
}

// startWaits are the lines the others see while players make their
// start-of-game choices, one per chooser: "Bo is choosing" (16.5). A
// chooser sees their own screen instead; nothing for other prompts.
export function startWaits(view: GameView, you: number, name: (p: number) => string): string[] {
  const w = view.waiting;
  if (view.over || !startChoice(w)) {
    return [];
  }
  const who = choosers(w);
  return who.includes(you) ? [] : who.map((p) => t('start.choosing', { player: name(p) }));
}

// keptZone is where a chooser's picked card waits for its flight to
// them after Confirm (20.22): playback flies it from there.
export function keptZone(player: number): string {
  return `choice-kept-${String(player)}`;
}
