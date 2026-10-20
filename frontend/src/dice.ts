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

// faceStyle places the texture, three dice wide, so one face fills the
// die. The textures are upside down, so the face turns half a turn
// about its own middle (15.64).
export function faceStyle(face: { col: number; row: number }) {
  const middle = (n: number) => `${String(((n + 0.5) / 3) * 100)}%`;
  return {
    left: `${String(-face.col * 100)}%`,
    top: `${String(-face.row * 100)}%`,
    transform: 'rotate(180deg)',
    transformOrigin: `${middle(face.col)} ${middle(face.row)}`,
  };
}

// DiceGroup is one set of dice looks in the settings menu (15.64).
export interface DiceGroup<T> {
  set: string;
  looks: T[];
}

// diceChoices groups the looks that have a texture by their set, in
// order. Without the private images there are none: dice stay plain.
export function diceChoices<T extends { set: string; image: string }>(looks: readonly T[]): DiceGroup<T>[] {
  const groups: DiceGroup<T>[] = [];
  for (const look of looks) {
    if (!look.image) {
      continue;
    }
    const group = groups.find((g) => g.set === look.set);
    if (group) {
      group.looks.push(look);
    } else {
      groups.push({ set: look.set, looks: [look] });
    }
  }
  return groups;
}

// CubeFaces are the six faces of a thrown die (15.69), as a d6: front
// and back add up to seven, and so do top and bottom, left and right.
export type CubeFaces = Record<'front' | 'back' | 'top' | 'bottom' | 'left' | 'right', number>;

// cubeFaces puts the value on top, as a die lies on a table (15.82). A
// value no face shows (a modified roll) still goes on top, around the
// faces of a one.
export function cubeFaces(value: number): CubeFaces {
  const base = faceOf(value) ? value : 1;
  // The sides come in pairs that add up to seven: the first two are
  // never a pair.
  const [front = 2, right = 3] = [1, 2, 3, 4, 5, 6].filter((n) => n !== base && n !== 7 - base);
  return { top: value, bottom: 7 - base, front, back: 7 - front, right, left: 7 - right };
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

// Gaps between dice of one throw, as parts of a throw: small and uneven,
// as dice leave a hand one by one (19.39).
const diceGaps = [0.13, 0.08, 0.17, 0.1, 0.15];

// diceDelays is when each die of a throw starts, as parts of a throw:
// the first at once, each next one a gap later.
export function diceDelays(count: number): number[] {
  const delays: number[] = [];
  let at = 0;
  for (let i = 0; i < count; i++) {
    delays.push(at);
    at += diceGaps.at(i % diceGaps.length) ?? 0;
  }
  return delays;
}
