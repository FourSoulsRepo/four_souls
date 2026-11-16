// The replay player's pure parts (10.2, ADR 031): a record's step as a
// playback update, where a move goes, which turn a step is in, and the
// autoplay's pause.
import type { RecordStep, RecordSummary } from './bridge';
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

// Autoplay waits this long after a step's effects at 1×, in ms (ADR 031,
// 4.6): after a step with events, and after one without.
export const pauseAfterEvents = 800;
export const pauseAfterQuiet = 200;

// autoplayPause is the wait after a step with events events at speed.
export function autoplayPause(events: number, speed: ReplaySpeed): number {
  return (events > 0 ? pauseAfterEvents : pauseAfterQuiet) / speed;
}
