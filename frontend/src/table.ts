// Where each player sits on screen (V-02). Pure, so it is tested
// without a browser.

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
