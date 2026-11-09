import { describe, expect, it } from 'vitest';
import { Purpose, type Choice } from './actions';
import { aimClick, arc, arcLift, cancelOption, chainOf, dashCount, dashFrames, dashOpacity, marchTime, promptAim, targetAt } from './aim';
import { feel } from './feel';
import type { OptionView } from './protocol';

const monster: OptionView = { object: 40, player: -1 };
const ann: OptionView = { player: 0, is_player: true };
const roll: OptionView = { stack_id: 3, player: 1 };

describe('targetAt', () => {
  const targets = [monster, ann, roll];

  it('finds a card, a player, or a stack item', () => {
    expect(targetAt(['obj-40', 'board'], targets)).toBe(monster);
    expect(targetAt(['obj-12', 'player-0'], targets)).toBe(ann);
    expect(targetAt(['stack-3'], targets)).toBe(roll);
  });

  it('prefers the topmost zone', () => {
    expect(targetAt(['obj-40', 'player-0'], targets)).toBe(monster);
  });

  it("finds the monster deck as an attack's target (20.23)", () => {
    const deck: OptionView = { player: -1, deck: true };
    expect(targetAt(['deck-monster'], [monster, deck])).toBe(deck);
    expect(targetAt(['deck-monster'], targets)).toBeUndefined();
  });

  it('finds nothing off the targets', () => {
    expect(targetAt(['obj-41', 'player-1', 'deck-loot'], targets)).toBeUndefined();
    expect(targetAt([], targets)).toBeUndefined();
  });

  it('never mistakes a player for an object with the same number', () => {
    expect(targetAt(['obj-0'], [ann])).toBeUndefined();
    expect(targetAt(['player-40'], [monster])).toBeUndefined();
  });
});

describe('aimClick', () => {
  const targets = [monster, ann];

  it('picks a target', () => {
    expect(aimClick(['obj-40', 'board'], 7, targets)).toEqual({ kind: 'pick', target: monster });
  });

  it('cancels on empty space', () => {
    expect(aimClick([], 7, targets)).toEqual({ kind: 'cancel' });
  });

  it('does nothing on the source, even over a player it may target', () => {
    expect(aimClick(['obj-7', 'hand-0', 'player-0'], 7, targets)).toEqual({ kind: 'none' });
  });

  it('calls any other card, player, pile, or stack item wrong', () => {
    expect(aimClick(['obj-8', 'obj-7', 'hand-0'], 7, targets)).toEqual({ kind: 'wrong' });
    expect(aimClick(['player-1'], 7, targets)).toEqual({ kind: 'wrong' });
    expect(aimClick(['deck-monster'], 7, targets)).toEqual({ kind: 'wrong' });
    expect(aimClick(['stack-3'], 7, targets)).toEqual({ kind: 'wrong' });
  });

  it('lets a source that is its own target pick itself', () => {
    expect(aimClick(['obj-40'], 40, targets)).toEqual({ kind: 'pick', target: monster });
  });
});

// The arrow's way (20.23, ADR 027 "Target arrow"): a curve bowing up
// the screen, dashes along it, marching and fading at its ends.
describe('arc', () => {
  const { lift, liftShare } = feel.aim;

  it('bows a level arrow up by lift at its middle', () => {
    const a = arc({ x: 0, y: 0 }, { x: 600, y: 0 });
    expect(arcLift({ x: 0, y: 0 }, { x: 600, y: 0 })).toBe(lift);
    const mid = a.at(a.length / 2);
    expect(mid.x).toBeCloseTo(300, 0);
    expect(mid.y).toBeCloseTo(-lift, 0);
    // Level at the middle; it rises from the source and falls to the head.
    expect(Math.abs(mid.angle)).toBeLessThan(2);
    expect(a.at(1).angle).toBeLessThan(-20);
    expect(a.end.angle).toBeGreaterThan(20);
    expect(a.length).toBeGreaterThan(600);
  });

  it('bows up whichever way a level arrow points', () => {
    const a = arc({ x: 600, y: 0 }, { x: 0, y: 0 });
    expect(a.at(a.length / 2).y).toBeCloseTo(-lift, 0);
  });

  it('bows a short arrow less', () => {
    expect(arcLift({ x: 0, y: 0 }, { x: 200, y: 0 })).toBe(200 * liftShare);
  });

  it('keeps an upright arrow straight, and bends a leaning one a little', () => {
    const up = arc({ x: 0, y: 0 }, { x: 0, y: -400 });
    expect(arcLift({ x: 0, y: 0 }, { x: 0, y: -400 })).toBe(0);
    expect(up.length).toBeCloseTo(400, 3);
    expect(up.end.angle).toBe(-90);
    // Just either side of upright the bow is small and on either side:
    // it fades through zero, never flips at full size.
    const right = arcLift({ x: 0, y: 0 }, { x: 10, y: -400 });
    const left = arcLift({ x: 0, y: 0 }, { x: -10, y: -400 });
    expect(right).toBeGreaterThan(0);
    expect(right).toBeLessThan(10);
    expect(left).toBeCloseTo(right, 6);
    expect(arc({ x: 0, y: 0 }, { x: 10, y: -400 }).at(200).x).toBeLessThan(0);
    expect(arc({ x: 0, y: 0 }, { x: -10, y: -400 }).at(200).x).toBeGreaterThan(0);
  });

  it('ends at the pointer and keeps poses on the curve', () => {
    const a = arc({ x: 10, y: 20 }, { x: 410, y: -80 });
    expect(a.at(a.length)).toMatchObject({ x: 410, y: -80 });
    expect(a.at(-5)).toMatchObject({ x: 10, y: 20 });
    expect(a.at(a.length + 50)).toMatchObject({ x: 410, y: -80 });
  });
});

describe('dashes', () => {
  const { step, head, stops, headFade, tailFade } = feel.aim;

  it('fill the chain short of the head, one per step begun', () => {
    const a = arc({ x: 0, y: 0 }, { x: 0, y: -300 });
    expect(chainOf(a)).toBeCloseTo(300 - head, 3);
    expect(dashCount(chainOf(a))).toBe(Math.ceil((300 - head) / step));
    expect(dashCount(0)).toBe(0);
    expect(chainOf(arc({ x: 0, y: 0 }, { x: 0, y: -20 }))).toBe(0);
  });

  it('fade in at the tail and out by the head, on their curves', () => {
    expect(dashOpacity(0, 300)).toBe(0);
    expect(dashOpacity(step / 2, 300)).toBeCloseTo(Math.pow(0.5, tailFade), 6);
    expect(dashOpacity(150, 300)).toBe(1);
    expect(dashOpacity(300 - step / 2, 300)).toBeCloseTo(Math.pow(0.5, headFade), 6);
    expect(dashOpacity(300, 300)).toBe(0);
    expect(dashOpacity(320, 300)).toBe(0);
  });

  it('march one step toward the head; each ends where the next starts', () => {
    const from = { x: 0, y: 0 };
    const a = arc(from, { x: 500, y: -100 });
    const first = dashFrames(a, from, 0);
    const second = dashFrames(a, from, 1);
    expect(first).toHaveLength(stops + 1);
    expect(first[0]?.offset).toBe(0);
    expect(first.at(-1)?.offset).toBe(1);
    expect(first.at(-1)?.transform).toBe(second[0]?.transform);
    expect(first.at(-1)?.opacity).toBe(second[0]?.opacity);
    // A dash starts at the source, faded out.
    expect(first[0]?.transform).toMatch(/^translate\(0px, 0px\)/);
    expect(first[0]?.opacity).toBe(0);
  });

  it('take a step at march px a second, stretched by the speed', () => {
    expect(marchTime(1)).toBeCloseTo((step / feel.aim.march) * 1000, 6);
    expect(marchTime(2)).toBeCloseTo(marchTime(1) * 2, 6);
  });
});

describe('promptAim', () => {
  const choose = (choice: number) => ({ player: 0, kind: 4, choice });
  const target: Choice = {
    title: 'Choose a target',
    purpose: Purpose.target,
    options: [
      { label: 'Ann', intent: choose(0), target: ann },
      { label: 'cancel', intent: choose(1) },
    ],
  };

  it('aims a placed card at its target prompt', () => {
    expect(promptAim(7, target)).toEqual({ source: 7, choice: target, targets: [ann] });
  });

  it('aims from object 0 too', () => {
    expect(promptAim(0, target)?.source).toBe(0);
  });

  it('keeps the dialog when the card is unknown or nothing aims', () => {
    expect(promptAim(undefined, target)).toBeNull();
    expect(promptAim(7, null)).toBeNull();
    expect(promptAim(7, { ...target, purpose: Purpose.answer })).toBeNull();
    expect(promptAim(7, { ...target, options: target.options.slice(1) })).toBeNull();
  });

  it('cancels by the option that points at nothing', () => {
    expect(cancelOption(target)).toEqual(choose(1));
  });
});
