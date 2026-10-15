// Where each player sits on screen (V-02). Pure, so it is tested
// without a browser.

import { Role, type CardView, type PlayerView } from './protocol';

export type Position = 'bottom' | 'left' | 'top' | 'right';

export interface Placed {
  position: Position;
  player: number;
}

// seatPositions puts the viewer at the bottom and the others in turn
// order, clockwise: left, top, right. Two players face each other;
// three leave the top free.
export function seatPositions(players: number, you: number): Placed[] {
  const order: Position[][] = [
    [],
    ['bottom'],
    ['bottom', 'top'],
    ['bottom', 'left', 'right'],
    ['bottom', 'left', 'top', 'right'],
  ];
  const positions = order.at(players) ?? [];
  return positions.map((position, i) => ({ position, player: (you + i) % players }));
}

// Each seat's own color: names on the stack (S-03) and the dot by a
// player's name. Red is left for targets.
const playerColors = ['#5ab0ff', '#5fd47a', '#c48cff', '#ff9a3c'];

// playerColor is a player's color; a game item (no player) has none.
export function playerColor(player: number): string | undefined {
  return player < 0 ? undefined : playerColors[player % playerColors.length];
}

// The win target when a view has none: an older server or record (R-WIN-01).
const defaultWinSouls = 4;

// soulSlots are a player's soul crystals (V-30): one slot per soul the
// game needs, the gained ones filled first. Souls past the target, at
// the end of a game, add filled slots.
export function soulSlots(souls: number, winSouls = 0): boolean[] {
  const target = winSouls > 0 ? winSouls : defaultWinSouls;
  return Array.from({ length: Math.max(target, souls) }, (_, i) => i < souls);
}

// isSoul: a card in play that is one of its owner's souls; it shows
// as crystals, not in the items row (V-43).
export function isSoul(c: CardView): boolean {
  return c.role === Role.soul;
}

// soulCards are the cards behind a player's filled crystals, bottom up
// (V-43): the character's own soul value first, then each soul card,
// once per soul it is worth. A view without the worth counts one.
export function soulCards(player: PlayerView): CardView[] {
  const cards: CardView[] = [];
  const add = (c: CardView, n: number) => {
    for (let i = 0; i < n; i++) {
      cards.push(c);
    }
  };
  add(player.character, player.character.souls ?? 0);
  for (const c of (player.in_play ?? []).filter(isSoul)) {
    add(c, Math.max(c.souls ?? 1, 1));
  }
  return cards;
}
