// The replay player's pure parts (10.2, ADR 031): a record's step as a
// playback update, where a move goes, which turn a step is in, and the
// autoplay's pause.
import { PromptKind } from './actions';
import type { RecordStep, RecordSummary } from './bridge';
import type { CardInfo } from './cards';
import type { Update } from './protocol';

// replaySpeeds are the player's speeds (ADR 031, 4.5).
export const replaySpeeds = [0.5, 1, 2, 4] as const;
export type ReplaySpeed = (typeof replaySpeeds)[number];

// Replaying is the record on the table: its summary, the step shown
// (an index into the live steps), the speed, and whether it autoplays.
export interface Replaying {
  summary: RecordSummary;
  at: number;
  speed: ReplaySpeed;
  playing: boolean;
  // hands: every hand shows face up (RP-04); true at first (ADR 031, 4.7).
  hands: boolean;
  // bar: the replay bar shows; false while the viewer hides it (ADR
  // 031, 4.9). True at each Watch; it never pauses the autoplay.
  bar: boolean;
}

// replayUpdate is a step as playback takes it (ADR 031, 3.4, 3.5): its
// events and view, the seats from the header; no intents, no aims, no
// deadline, no rollback.
export function replayUpdate(step: RecordStep, summary: Pick<RecordSummary, 'seats'>): Update {
  return {
    step: step.n,
    events: step.events,
    view: step.view,
    allowed: null,
    seats: summary.seats.map((s) => ({ seat: s.seat, name: s.name, connected: true })),
  };
}

// clampStep keeps a step index inside a record of steps steps.
export function clampStep(index: number, steps: number): number {
  return Math.max(0, Math.min(steps - 1, Math.round(index)));
}

// animates: only one step forward plays its effects; back and any jump
// show the step's view at once (ADR 031, 4.2, 4.3).
export function animates(from: number, to: number): boolean {
  return to === from + 1;
}

// turnAt is the index into turns of the turn step index is in; -1
// before the first.
export function turnAt(turns: RecordSummary['turns'], index: number): number {
  return turns.findLastIndex((t) => t.step <= index);
}

// sliderFill is how far the slider is along, in percent: 0 at the
// first step, 100 at the last. The bar paints the track up to it.
export function sliderFill(at: number, steps: number): number {
  if (steps <= 1) {
    return 100;
  }
  return (Math.min(Math.max(at, 0), steps - 1) / (steps - 1)) * 100;
}

// Autoplay waits this long after a step's effects at 1×, in ms (ADR 031,
// 4.6): after a step with events, and after one without.
export const pauseAfterEvents = 800;
export const pauseAfterQuiet = 200;

// autoplayPause is the wait after a step with events events at speed.
export function autoplayPause(events: number, speed: ReplaySpeed): number {
  return (events > 0 ? pauseAfterEvents : pauseAfterQuiet) / speed;
}

// hideHands is an update as an onlooker with no seat sees it (RP-04,
// 10.3): no hand cards, no peeked treasure, none of the chooser's
// options, and private events (a looted card) without their card, as
// the engine's FilterEvents leaves them. The record keeps every hand;
// this only changes what the table draws.
export function hideHands(u: Update): Update {
  const { view } = u;
  return {
    ...u,
    events: u.events?.map((e) => (e.private ? { ...e, card: undefined, object: undefined } : e)) ?? null,
    view: {
      ...view,
      players: view.players.map((p) => ({ ...p, hand: undefined, treasure_top: undefined })),
      waiting: view.waiting.kind === PromptKind.choose ? { ...view.waiting, options: undefined } : view.waiting,
      subject: undefined,
      ordering: undefined,
      options: undefined,
      activating: undefined,
      modes: undefined,
    },
  };
}

// RecordCard is a card's text as the record's header keeps it (RP-10).
type RecordCard = RecordSummary['cards'][number];

// withRecordCards adds the record's card text for each ref the app has
// no data for, e.g. a card of a newer set (ADR 031, 2.4): it draws as a
// text card. The app's own data wins for refs it has.
export function withRecordCards(cards: Map<string, CardInfo>, header: RecordCard[]): Map<string, CardInfo> {
  const missing = header.filter((c) => !cards.has(c.ref));
  if (missing.length === 0) {
    return cards;
  }
  const more = missing.map((c): [string, CardInfo] => [
    c.ref,
    { ref: c.ref, set: '', type: c.type, name: c.name, effects: c.effects, footnotes: [], stats: [], rewards: [], image: '', record: true },
  ]);
  return new Map([...cards, ...more]);
}

// withoutRecordCards drops what withRecordCards added, when the replay
// closes.
export function withoutRecordCards(cards: Map<string, CardInfo>): Map<string, CardInfo> {
  const own = [...cards].filter(([, c]) => c.record !== true);
  return own.length === cards.size ? cards : new Map(own);
}
