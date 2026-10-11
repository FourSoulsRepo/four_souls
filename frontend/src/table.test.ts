import { describe, expect, it } from 'vitest';
import { seatPositions } from './table';

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
