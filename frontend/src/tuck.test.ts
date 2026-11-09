import { describe, expect, it } from 'vitest';
import { feel } from './feel';
import { raised, tuckClip, tuckPose, tuckSide, tuckTime } from './tuck';

// A pile on the left of a 1000 px wide screen, and a card held over it.
const pile = { x: 200, y: 300, w: 80, rotate: 0 };
const from = raised(pile, 0, 1);

describe('tuck (19.147)', () => {
  it('slides a card out toward the middle of the screen', () => {
    expect(tuckSide(pile, 1000)).toBe(1);
    expect(tuckSide({ ...pile, x: 900 }, 1000)).toBe(-1);
  });

  it('holds the chooser\'s card over the pile, smaller', () => {
    expect(from.y).toBeLessThan(pile.y);
    expect(from.w).toBeLessThan(pile.w);
    const [a, b] = [raised(pile, 0, 2), raised(pile, 1, 2)];
    expect(a.x).toBeLessThan(b.x);
    expect(a.rotate).toBe(-b.rotate);
  });

  it('drops a kept card onto the pile and fades it there', () => {
    const k = feel.tuck;
    expect(tuckTime(false)).toBe(k.drop);
    expect(tuckPose(0, from, pile, false, 1)).toMatchObject({ x: from.x, y: from.y, hide: 0, opacity: 1 });
    expect(tuckPose(k.drop, from, pile, false, 1)).toMatchObject({ x: pile.x, y: pile.y, w: pile.w, rotate: 0, hide: 0, opacity: 0 });
  });

  it('slides a bottomed card out beside the pile, then under it', () => {
    const k = feel.tuck;
    expect(tuckTime(true)).toBe(k.out + k.under);
    const out = tuckPose(k.out - 1e-9, from, pile, true, 1);
    expect(out.x).toBeCloseTo(pile.x + pile.w * 1.1, 0);
    expect(out.hide).toBe(0);
    // Halfway under, the pile covers part of it, on the pile's side.
    const mid = tuckPose(k.out + k.under / 2, from, pile, true, 1);
    expect(mid.hide).toBeGreaterThan(0);
    expect(mid.hide).toBeLessThan(1);
    expect(mid.hide).toBeCloseTo(1 - (mid.x - pile.x) / pile.w, 5);
    const end = tuckPose(k.out + k.under, from, pile, true, 1);
    expect(end).toMatchObject({ x: pile.x, hide: 1, opacity: 0 });
  });

  it('clips the side toward the pile', () => {
    expect(tuckClip(0.25, 1)).toBe('inset(0 0 0 25.00%)');
    expect(tuckClip(0.25, -1)).toBe('inset(0 25.00% 0 0)');
  });
});
