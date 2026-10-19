import { describe, expect, it } from 'vitest';
import { backSlot, empty, fromFrame, onFrame, presenceFor, same, spotAt, spotOf } from './presence';

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
});

describe('same and empty', () => {
  it('counts a pointer that barely moved as the same', () => {
    const a = { from: { object: 1 }, point: { x: 0.5, y: 0.5 } };
    expect(same(a, { from: { object: 1 }, point: { x: 0.501, y: 0.5 } })).toBe(true);
    expect(same(a, { from: { object: 1 }, point: { x: 0.6, y: 0.5 } })).toBe(false);
    expect(same(a, { from: { object: 1 }, at: { player: 0 } })).toBe(false);
    expect(same({}, {})).toBe(true);
  });

  it('knows a presence that points at nothing', () => {
    expect(empty({})).toBe(true);
    expect(empty({ at: { deck: 0 } })).toBe(false);
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
