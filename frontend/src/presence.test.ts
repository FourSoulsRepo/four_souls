import { describe, expect, it } from 'vitest';
import { backSlot, empty, fromFrame, marksOf, onFrame, presenceFor, same, sendGap, spotAt, spotOf } from './presence';

const hand = [70, 71, 72];

describe('spotOf', () => {
  it('turns zones into spots', () => {
    expect(spotOf('obj-40', hand)).toEqual({ object: 40 });
    expect(spotOf('player-1', hand)).toEqual({ player: 1 });
    expect(spotOf('stack-3', hand)).toEqual({ stack_id: 3 });
    expect(spotOf('deck-monster', hand)).toEqual({ deck: 2 });
  });

  it('shares an own hand card as its slot, never its object', () => {
    expect(spotOf('obj-71', hand)).toEqual({ hand: 2 });
  });

  it('skips zones that are no spot', () => {
    for (const z of ['hand-1', 'discard-loot', 'monster-slot-0', 'stack-src-5', 'deck-nothing']) {
      expect(spotOf(z, hand), z).toBeUndefined();
    }
  });
});

describe('spotAt', () => {
  it('takes the topmost zone that is a spot', () => {
    expect(spotAt(['hand-1', 'obj-40', 'player-0'], hand)).toEqual({ object: 40 });
    expect(spotAt(['hand-1'], hand)).toBeUndefined();
  });
});

describe('presenceFor', () => {
  const point = { x: 0.5, y: 0.25 };

  it('shares an aim from a hand slot to the spot under the pointer', () => {
    expect(presenceFor({ source: 72, active: false }, ['obj-40'], hand, point)).toEqual({ from: { hand: 3 }, at: { object: 40 } });
  });

  it('shares the pointer while an aim is over no spot', () => {
    expect(presenceFor({ source: 40, active: true }, [], hand, point)).toEqual({ from: { object: 40 }, point });
  });

  it('shares a hover only on the own turn (owner, 2026-10-20)', () => {
    expect(presenceFor({ active: true }, ['obj-71'], hand, point)).toEqual({ at: { hand: 2 } });
    expect(presenceFor({ active: true }, [], hand, point)).toEqual({});
    expect(presenceFor({ active: false }, ['obj-40'], hand, point)).toEqual({});
  });

  it('shares a dragged card alone, as its slot, on any turn (20.13)', () => {
    expect(presenceFor({ held: 71, active: false }, ['obj-40'], hand, point)).toEqual({ held: { hand: 2 } });
    expect(presenceFor({ held: 72, active: true }, [], hand, point)).toEqual({ held: { hand: 3 } });
  });
});

describe('sendGap (20.13)', () => {
  it('waits 0.35 s, or 0.25 s while aiming (ADR 027 "Other players live")', () => {
    expect(sendGap({ at: { object: 40 } })).toBe(350);
    expect(sendGap({ held: { hand: 1 } })).toBe(350);
    expect(sendGap({ from: { hand: 1 }, point: { x: 0, y: 0 } })).toBe(250);
  });
});

describe('marksOf (20.13)', () => {
  it('marks the target of an aim red, a plain hover in the player color', () => {
    expect(marksOf({ from: { object: 5 }, at: { object: 40 } })).toEqual([[{ object: 40 }, 'at']]);
    expect(marksOf({ at: { object: 40 } })).toEqual([[{ object: 40 }, 'hover']]);
    expect(marksOf({ at: { hand: 2 } })).toEqual([[{ hand: 2 }, 'hover']]);
  });

  it('lifts the hand slot aimed with or held', () => {
    expect(marksOf({ from: { hand: 1 }, point: { x: 0, y: 0 } })).toEqual([[{ hand: 1 }, 'from']]);
    expect(marksOf({ held: { hand: 3 } })).toEqual([[{ hand: 3 }, 'held']]);
    expect(marksOf({})).toEqual([]);
  });
});

describe('same and empty', () => {
  it('counts a pointer that barely moved as the same', () => {
    const a = { from: { object: 1 }, point: { x: 0.5, y: 0.5 } };
    expect(same(a, { from: { object: 1 }, point: { x: 0.501, y: 0.5 } })).toBe(true);
    expect(same(a, { from: { object: 1 }, point: { x: 0.6, y: 0.5 } })).toBe(false);
    expect(same(a, { from: { object: 1 }, at: { player: 0 } })).toBe(false);
    expect(same({}, {})).toBe(true);
    expect(same({ held: { hand: 1 } }, { held: { hand: 2 } })).toBe(false);
  });

  it('knows a presence that points at nothing', () => {
    expect(empty({})).toBe(true);
    expect(empty({ at: { deck: 0 } })).toBe(false);
    expect(empty({ held: { hand: 1 } })).toBe(false);
  });
});

describe('board frame points', () => {
  const box = { left: 100, top: 50, width: 400, height: 200 };

  it('go to the frame and back', () => {
    expect(onFrame(300, 100, box)).toEqual({ x: 0.5, y: 0.25 });
    expect(fromFrame({ x: 0.5, y: 0.25 }, box)).toEqual({ x: 300, y: 100 });
  });

  it('stay inside the frame', () => {
    expect(onFrame(0, 999, box)).toEqual({ x: 0, y: 1 });
  });
});

describe('backSlot', () => {
  it('puts a slot past the fanned backs on the last one', () => {
    expect(backSlot(1, 8)).toBe(0);
    expect(backSlot(10, 8)).toBe(7);
  });
});
