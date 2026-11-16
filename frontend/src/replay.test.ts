import { describe, expect, it } from 'vitest';
import type { RecordStep } from './bridge';
import type { GameView } from './protocol';
import { animates, autoplayPause, clampStep, replayUpdate, turnAt } from './replay';

const turns = [
  { number: 0, active: 0, step: 0 },
  { number: 1, active: 1, step: 2 },
  { number: 2, active: 0, step: 20 },
];

// The replay player's pure parts (10.2, ADR 031).
describe('replay', () => {
  it('makes a step an update with no intents, aims or deadline', () => {
    const step: RecordStep = { index: 4, n: 7, events: [{ kind: 'gained_cents', player: 1, amount: 3 }], view: { turn: { number: 2 } } as GameView };
    const u = replayUpdate(step, { seats: [{ seat: 0, name: 'Ann' }, { seat: 1, name: 'Bo' }] });
    expect(u).toEqual({
      step: 7,
      events: step.events,
      view: step.view,
      allowed: null,
      seats: [
        { seat: 0, name: 'Ann', connected: true },
        { seat: 1, name: 'Bo', connected: true },
      ],
    });
  });

  it('keeps a move inside the record', () => {
    expect(clampStep(-1, 30)).toBe(0);
    expect(clampStep(30, 30)).toBe(29);
    expect(clampStep(12.4, 30)).toBe(12);
  });

  it('plays effects only one step forward; back and jumps show at once', () => {
    expect(animates(4, 5)).toBe(true);
    expect(animates(5, 4)).toBe(false);
    expect(animates(4, 9)).toBe(false);
    expect(animates(4, 4)).toBe(false);
  });

  it('finds the turn a step is in', () => {
    expect(turnAt(turns, 0)).toBe(0);
    expect(turnAt(turns, 1)).toBe(0);
    expect(turnAt(turns, 2)).toBe(1);
    expect(turnAt(turns, 25)).toBe(2);
    expect(turnAt([{ number: 1, active: 0, step: 3 }], 1)).toBe(-1);
  });

  it('pauses 0.8 s after events and 0.2 s after none at 1×, divided by the speed', () => {
    expect(autoplayPause(2, 1)).toBe(800);
    expect(autoplayPause(0, 1)).toBe(200);
    expect(autoplayPause(1, 4)).toBe(200);
    expect(autoplayPause(0, 0.5)).toBe(400);
  });
});
