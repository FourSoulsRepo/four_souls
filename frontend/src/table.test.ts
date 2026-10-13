import { describe, expect, it } from 'vitest';
import { seatPositions, soulSlots } from './table';

describe('seatPositions', () => {
  it('puts the viewer at the bottom and the other player on top', () => {
    expect(seatPositions(2, 1)).toEqual([
      { position: 'bottom', player: 1 },
      { position: 'top', player: 0 },
    ]);
  });

  it('seats three players left and right, in turn order', () => {
    expect(seatPositions(3, 2)).toEqual([
      { position: 'bottom', player: 2 },
      { position: 'left', player: 0 },
      { position: 'right', player: 1 },
    ]);
  });

  it('seats four players left, top, right', () => {
    expect(seatPositions(4, 1).map((p) => `${p.position}:${String(p.player)}`)).toEqual([
      'bottom:1',
      'left:2',
      'top:3',
      'right:0',
    ]);
  });
});

describe('soulSlots', () => {
  it('fills one slot per soul, the rest left to win', () => {
    expect(soulSlots(0, 4)).toEqual([false, false, false, false]);
    expect(soulSlots(2, 4)).toEqual([true, true, false, false]);
  });

  it('takes its length from the win target', () => {
    expect(soulSlots(1, 6)).toHaveLength(6);
  });

  it('needs four souls when the view has no target', () => {
    expect(soulSlots(1)).toEqual([true, false, false, false]);
  });

  it('adds filled slots past the target', () => {
    expect(soulSlots(5, 4)).toEqual([true, true, true, true, true]);
  });
});
