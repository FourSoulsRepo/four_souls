import { describe, expect, it } from 'vitest';
import { cubeFaces, diceChoices, diceDelays, diceImage, faceOf, faceStyle, firstRollRows, lastRollAfter, lastRollText } from './dice';

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

describe('cubeFaces', () => {
  it('puts the value on top of a real d6 (15.69, 15.82)', () => {
    for (let v = 1; v <= 6; v++) {
      const f = cubeFaces(v);
      expect(f.top).toBe(v);
      expect(new Set(Object.values(f))).toEqual(new Set([1, 2, 3, 4, 5, 6]));
      expect([f.front + f.back, f.top + f.bottom, f.left + f.right]).toEqual([7, 7, 7]);
    }
  });

  it('shows a modified value on top of a one', () => {
    expect(cubeFaces(7)).toEqual({ ...cubeFaces(1), top: 7 });
    expect(cubeFaces(0).top).toBe(0);
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

describe('faceStyle', () => {
  it('turns the face half a turn about its own middle (15.64)', () => {
    expect(faceStyle({ col: 2, row: 1 })).toEqual({
      left: '-200%',
      top: '-100%',
      transform: 'rotate(180deg)',
      transformOrigin: '83.33333333333334% 50%',
    });
  });
});

describe('diceDelays', () => {
  it('starts one die at once and several at small uneven gaps (19.39)', () => {
    expect(diceDelays(1)).toEqual([0]);
    const d = diceDelays(4);
    const gaps = d.slice(1).map((x, i) => x - (d.at(i) ?? 0));
    expect(d[0]).toBe(0);
    expect(gaps.every((g) => g > 0.05 && g < 0.2)).toBe(true);
    expect(new Set(gaps).size).toBe(gaps.length);
  });
});

describe('firstRollRows', () => {
  it('lists this round in seat order, with the dice so far and the next roller (15.43)', () => {
    expect(firstRollRows({ round: 1, players: [0, 1, 2], rolls: [4] })).toEqual([
      { player: 0, value: 4, next: false },
      { player: 1, next: true },
      { player: 2, next: false },
    ]);
    expect(firstRollRows({ round: 2, players: [1, 2] })).toEqual([
      { player: 1, next: true },
      { player: 2, next: false },
    ]);
  });
});
