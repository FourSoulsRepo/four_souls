import { describe, expect, it } from 'vitest';
import { drawTime, otherPose, ownPose, type DrawPath } from './draw';
import { drawPace, feel } from './feel';

describe('drawPace (20.7)', () => {
  it('keeps draws 1–3 at full time, then each 11 % faster, down to 0.556', () => {
    const pace = [0, 1, 2, 3, 4, 5, 6, 7, 12].map((i) => Number(drawPace(i).toFixed(3)));
    expect(pace).toEqual([1, 1, 1, 0.889, 0.778, 0.667, 0.556, 0.556, 0.556]);
  });
});

const path: DrawPath = {
  deck: { x: 900, y: 400 },
  deckScale: 0.25,
  show: { x: 640, y: 320 },
  showScale: 0.625,
  hand: { x: 640, y: 660 },
  handScale: 0.3125,
};

describe('ownPose (20.7)', () => {
  it('takes 1.5 s: deck, show spot, hand', () => {
    expect(drawTime(true)).toBe(1.5);
    expect(ownPose(0, path)).toMatchObject({ x: 900, y: 400, scale: 0.25, turn: 180 });
    expect(ownPose(feel.draw.show, path)).toMatchObject({ x: 640, y: 320, scale: 0.625, turn: 0 });
    expect(ownPose(1.5, path)).toMatchObject({ x: 640, y: 660, scale: 0.3125, turn: 0 });
  });

  it('rises from the deck first', () => {
    expect(ownPose(0.01, path).y).toBeLessThan(400);
  });

  it('turns face up and grows only after a short wait', () => {
    expect(ownPose(feel.draw.turnAt, path)).toMatchObject({ turn: 180, scale: 0.25 });
    expect(ownPose(0.5, path).turn).toBeLessThan(10);
  });

  it('holds at the show spot before it goes to the hand', () => {
    const p = ownPose(feel.draw.show + feel.draw.hand / 2, path);
    expect(p.y - 320).toBeLessThan(20);
  });
});

describe('otherPose (20.7)', () => {
  it('flies deck to hand, back up, with no show spot', () => {
    const other = { ...path, show: undefined, handScale: 0.25 };
    expect(drawTime(false)).toBe(0.8);
    expect(otherPose(0, other)).toMatchObject({ x: 900, y: 400, turn: 180 });
    expect(otherPose(0.8, other)).toMatchObject({ x: 640, y: 660, scale: 0.25, turn: 180 });
    for (let t = 0; t <= 0.8; t += 0.1) {
      expect(otherPose(t, other).turn).toBe(180);
    }
  });
});
