import { describe, expect, it } from 'vitest';
import { feel } from './feel';
import { answer, dragged, gapMoves, into, release } from './hold';

describe('picking a hand card up (20.20)', () => {
  it('drags only past the threshold, either way', () => {
    expect(dragged(0, 0)).toBe(false);
    expect(dragged(feel.play.drag - 1, 0)).toBe(false);
    expect(dragged(0, -feel.play.drag)).toBe(true);
    expect(dragged(6, 6)).toBe(true);
  });

  it('plays a held card by a click on the board, never on a button', () => {
    expect(release(true, false)).toBe('play');
    expect(release(true, true)).toBe('back');
    expect(release(false, false)).toBe('back');
  });
});

describe('the answer to a played card (20.20)', () => {
  it('waits while the server has not answered', () => {
    expect(answer({ step: 4, now: 4, sent: 4, more: false })).toBe('wait');
  });

  it('goes home when refused', () => {
    expect(answer({ step: 4, now: 4, sent: null, more: false })).toBe('refused');
  });

  it('goes home when the card stays in the hand after the answer', () => {
    expect(answer({ step: 4, now: 5, sent: null, more: false })).toBe('home');
  });

  it('waits on while its aimed target is picked next', () => {
    expect(answer({ step: 4, now: 5, sent: null, more: true })).toBe('wait');
  });
});

describe('the hand closing its gap (20.20)', () => {
  it('moves the cards that stay from their old places', () => {
    const before = new Map([
      [1, 0],
      [2, 100],
      [3, 200],
    ]);
    const after = new Map([
      [1, 50],
      [3, 150],
    ]);
    expect(gapMoves(before, after)).toEqual([
      { id: 1, dx: -50 },
      { id: 3, dx: 50 },
    ]);
  });

  it('moves nothing when no card left', () => {
    expect(gapMoves(new Map([[1, 0]]), new Map([[1, 40], [2, 80]]))).toEqual([]);
  });
});

describe('a played card into its stack row (20.20)', () => {
  const natural = { w: 100, h: 140 };

  it('ends as tall as the row, at its middle', () => {
    const p = into({ x: 500, y: 300, w: 108, h: 151.2 }, natural, { x: 100, y: 50, w: 200, h: 28 });
    expect(p.from).toBeCloseTo(1.08);
    expect(p.to).toBeCloseTo(0.2);
    // The top left ends where the shrunk card is centred on the row.
    expect(500 + p.dx + (natural.w * p.to) / 2).toBeCloseTo(200);
    expect(300 + p.dy + (natural.h * p.to) / 2).toBeCloseTo(64);
  });

  it('shrinks in place with no row', () => {
    const p = into({ x: 0, y: 0, w: 100, h: 140 }, natural, null);
    expect(p.to).toBe(0.5);
    expect(p.dx + (natural.w * p.to) / 2).toBeCloseTo(50);
    expect(p.dy + (natural.h * p.to) / 2).toBeCloseTo(70);
  });
});
