import { describe, expect, it } from 'vitest';
import { diceChoices, diceImage, faceOf, lastRollAfter, lastRollText, tumbleFaces } from './dice';

describe('faceOf', () => {
  it('finds each face of a d6 in the texture', () => {
    const cells = [1, 2, 3, 4, 5, 6].map((v) => faceOf(v));
    expect(cells).toEqual([
      { col: 0, row: 2 },
      { col: 0, row: 1 },
      { col: 1, row: 2 },
      { col: 1, row: 1 },
      { col: 2, row: 1 },
      { col: 2, row: 2 },
    ]);
    expect(new Set(cells.map((c) => JSON.stringify(c))).size).toBe(6);
  });

  it('has no face for modified values', () => {
    expect(faceOf(0)).toBeNull();
    expect(faceOf(7)).toBeNull();
  });
});

describe('tumbleFaces', () => {
  it('shows two other faces, then lands on the value (15.62)', () => {
    for (let v = 1; v <= 6; v++) {
      const faces = tumbleFaces(v);
      expect(faces.at(-1)).toBe(v);
      expect(new Set(faces).size).toBe(3);
      expect(faces.every((f) => f >= 1 && f <= 6)).toBe(true);
    }
  });
});

describe('lastRollAfter', () => {
  it('keeps the final value after changes and rerolls (15.63)', () => {
    let last = lastRollAfter(null, [{ kind: 'dice_rolled', player: 1, amount: 2, text: 'attack roll' }]);
    expect(last).toEqual({ player: 1, value: 2, label: 'attack roll' });
    last = lastRollAfter(last, [{ kind: 'dice_rolled', player: 1, amount: 5, text: 'reroll' }]);
    expect(last).toEqual({ player: 1, value: 5, label: 'attack roll' });
    last = lastRollAfter(last, [
      { kind: 'roll_changed', player: 1, amount: 6 },
      { kind: 'roll_resolved', player: 1, amount: 4, text: 'attack' },
    ]);
    expect(last).toEqual({ player: 1, value: 4, label: 'attack roll' });
    expect(lastRollAfter(last, [{ kind: 'looted', player: 0 }])).toBe(last);
    // A new roll replaces it.
    expect(lastRollAfter(last, [{ kind: 'dice_rolled', player: 0, amount: 3, text: 'roll' }])).toEqual({ player: 0, value: 3, label: 'roll' });
  });

  it('names the roller and what the roll was for', () => {
    const name = (p: number) => ['Ann', 'Bo'].at(p) ?? '?';
    expect(lastRollText({ player: 0, value: 4, label: 'attack roll' }, name)).toBe('Ann rolled 4 — attack roll');
    expect(lastRollText({ player: 1, value: 2, label: '' }, name)).toBe('Bo rolled 2');
  });
});

describe('diceImage', () => {
  const looks = [
    { id: 'tear', image: '/tear.webp' },
    { id: 'gold', image: '/gold.webp' },
    { id: 'lost', image: '' },
  ];
  const seats = [
    { seat: 0, name: 'Ann', connected: true, dice: 'gold' },
    { seat: 1, name: 'Bot', connected: true },
    { seat: 2, name: 'Cy', connected: true, dice: 'd20' },
    { seat: 3, name: 'Di', connected: true, dice: 'lost' },
  ];

  it("draws each roller's own look (15.6)", () => {
    expect(diceImage(looks, seats, 0)).toBe('/gold.webp');
  });

  it('gives a seat without a look the default one', () => {
    expect(diceImage(looks, seats, 1)).toBe('/tear.webp');
    expect(diceImage(looks, null, 0)).toBe('/tear.webp');
  });

  it('draws an unknown look or a missing file plain', () => {
    expect(diceImage(looks, seats, 2)).toBe('');
    expect(diceImage(looks, seats, 3)).toBe('');
    expect(diceImage([], seats, 1)).toBe('');
  });
});

describe('diceChoices', () => {
  it('groups the looks with a texture by set, in order (15.64)', () => {
    const looks = [
      { id: 'tear', set: 'Holy', image: '/tear.webp' },
      { id: 'poop', set: 'Unholy', image: '/poop.webp' },
      { id: 'cent', set: 'Holy', image: '/cent.webp' },
      { id: 'lost', set: 'Unholy', image: '' },
    ];
    expect(diceChoices(looks).map((g) => [g.set, g.looks.map((d) => d.id)])).toEqual([
      ['Holy', ['tear', 'cent']],
      ['Unholy', ['poop']],
    ]);
  });

  it('offers nothing without the images', () => {
    expect(diceChoices([{ id: 'tear', set: 'Holy', image: '' }])).toEqual([]);
  });
});
