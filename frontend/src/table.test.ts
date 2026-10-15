import { describe, expect, it } from 'vitest';
import { Role, type CardView, type PlayerView } from './protocol';
import { seatPositions, soulCards, soulSlots } from './table';

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

describe('soulCards', () => {
  const card = (id: number, role: number, souls?: number): CardView => ({ id, card: `c${String(id)}`, role, charged: false, souls });
  const player = (character: CardView, inPlay: CardView[]) => ({ id: 0, character, in_play: inPlay }) as unknown as PlayerView;
  const ids = (cards: CardView[]) => cards.map((c) => c.id);

  it('ties one crystal to each soul card, in order', () => {
    const p = player(card(1, Role.character), [card(2, Role.soul, 1), card(3, Role.item), card(4, Role.soul, 1)]);
    expect(ids(soulCards(p))).toEqual([2, 4]);
  });

  it('ties both crystals of a two-soul card to it', () => {
    const p = player(card(1, Role.character), [card(2, Role.soul, 2), card(3, Role.soul, 1)]);
    expect(ids(soulCards(p))).toEqual([2, 2, 3]);
  });

  it("puts a character's own soul value first", () => {
    const p = player(card(1, Role.character, 1), [card(2, Role.soul, 1)]);
    expect(ids(soulCards(p))).toEqual([1, 2]);
  });

  it('counts one soul when the view has no worth', () => {
    const p = player(card(1, Role.character), [card(2, Role.soul)]);
    expect(ids(soulCards(p))).toEqual([2]);
  });

  it('has none without souls', () => {
    expect(soulCards(player(card(1, Role.character), []))).toEqual([]);
  });
});
