// Dice looks (V-19). A dice texture is a 3×3 grid of faces, as in
// Tabletop Simulator: the top row is unused, then 2 4 5 and 1 3 6.
import type { GameEvent, Seat } from './protocol';

const faces = new Map([
  [2, { col: 0, row: 1 }],
  [4, { col: 1, row: 1 }],
  [5, { col: 2, row: 1 }],
  [1, { col: 0, row: 2 }],
  [3, { col: 1, row: 2 }],
  [6, { col: 2, row: 2 }],
]);

// faceOf finds a value's cell in the texture; null for values a d6 has
// no face for (a modified roll can be 0 or 7), drawn as plain numbers.
export function faceOf(value: number): { col: number; row: number } | null {
  return faces.get(value) ?? null;
}

// diceImage is a player's dice texture (15.6, V-19): their seat's look,
// or the default (the first) for a seat without one, as a bot's. An
// unknown look or a missing file gives "": a plain die.
export function diceImage(looks: readonly { id: string; image: string }[], seats: readonly Seat[] | null | undefined, player: number): string {
  const id = seats?.find((s) => s.seat === player)?.dice;
  const look = id ? looks.find((d) => d.id === id) : looks[0];
  return look?.image ?? '';
}

// tumbleFaces are the faces a thrown die shows as it lands and bounces
// (15.62): two others, then the value it rolled.
export function tumbleFaces(value: number): number[] {
  const other = (n: number) => ((value + n - 1) % 6 + 6) % 6 + 1;
  return [other(2), other(4), value];
}

// LastRoll is the last roll, kept as a small die by the stack (15.63,
// F-21): who rolled, its final value, and what it was for.
export interface LastRoll {
  player: number;
  value: number;
  label: string;
}

// lastRollAfter follows one update's events: a throw starts a new last
// roll; a reroll, a change and the final result change its value.
export function lastRollAfter(prev: LastRoll | null, events: GameEvent[]): LastRoll | null {
  let last = prev;
  for (const e of events) {
    const value = e.amount ?? 0;
    const same = last !== null && last.player === e.player;
    switch (e.kind) {
      case 'dice_rolled':
        last = e.text === 'reroll' && last && same ? { ...last, value } : { player: e.player, value, label: e.text ?? '' };
        break;
      case 'roll_changed':
      case 'roll_resolved':
        if (last && same && value > 0) {
          last = { ...last, value };
        }
        break;
    }
  }
  return last;
}

// lastRollText names the roller and what the roll was for:
// "Ann rolled 4 — attack roll".
export function lastRollText(r: LastRoll, player: (p: number) => string): string {
  return `${player(r.player)} rolled ${String(r.value)}${r.label ? ` — ${r.label}` : ''}`;
}
