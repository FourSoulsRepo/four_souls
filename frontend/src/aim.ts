// Aiming (15.25): what an aiming card points at. Pure, except
// zonesAt, which reads the page.
import type { OptionView } from './protocol';

// zonesAt lists the zones under a point, topmost first: data-zone of
// cards ("obj-N") and piles, data-stack of stack rows. A player is aimed
// at by the character card (data-player, "player-N"), never by the area
// around it (15.33), so the areas' own zones are left out.
export function zonesAt(x: number, y: number): string[] {
  const out: string[] = [];
  for (const el of document.elementsFromPoint(x, y)) {
    const zone = el.getAttribute('data-zone');
    const stack = el.getAttribute('data-stack');
    const player = el.getAttribute('data-player');
    if (stack) {
      out.push(`stack-${stack}`);
    }
    if (zone && !zone.startsWith('player-')) {
      out.push(zone);
    }
    if (player) {
      out.push(`player-${player}`);
    }
  }
  return out;
}

// targetAt is the first of the targets the zones show, topmost first:
// a character card as an object wins over it as its player.
export function targetAt(zones: string[], targets: OptionView[]): OptionView | undefined {
  for (const zone of zones) {
    const m = /^(obj|player|stack)-(\d+)$/.exec(zone);
    if (!m) {
      continue;
    }
    const n = Number(m[2]);
    const t = targets.find((o) => {
      switch (m[1]) {
        case 'obj':
          return !o.is_player && o.object === n;
        case 'player':
          return !!o.is_player && o.player === n;
        default:
          return o.stack_id === n;
      }
    });
    if (t) {
      return t;
    }
  }
  return undefined;
}

// AimClick is what a click while aiming does (15.33): pick a target;
// cancel on the source (T-07) or on empty space (T-06); or, on a card,
// pile, or stack item that is no target, say why (T-05).
export type AimClick = { kind: 'pick'; target: OptionView } | { kind: 'cancel' } | { kind: 'wrong' };

// aimClick decides a click by the zones under it, topmost first. The
// source cancels, unless it is a target itself (an item that may pick
// any item, as The D20).
export function aimClick(zones: string[], source: number, targets: OptionView[]): AimClick {
  const src = `obj-${String(source)}`;
  if (zones.find((z) => z.startsWith('obj-')) === src) {
    const self = targetAt([src], targets);
    return self ? { kind: 'pick', target: self } : { kind: 'cancel' };
  }
  const target = targetAt(zones, targets);
  if (target) {
    return { kind: 'pick', target };
  }
  return zones.length === 0 ? { kind: 'cancel' } : { kind: 'wrong' };
}

// aimed: the target is one of the targets.
export function aimed(targets: OptionView[] | null, test: (o: OptionView) => boolean): boolean {
  return targets?.some(test) ?? false;
}

export interface Point {
  x: number;
  y: number;
}

// arrow is how an aiming arrow from a point to a point is drawn: its
// length, and its angle in degrees from the x axis.
export function arrow(from: Point, to: Point): { length: number; angle: number } {
  const dx = to.x - from.x;
  const dy = to.y - from.y;
  return { length: Math.hypot(dx, dy), angle: (Math.atan2(dy, dx) * 180) / Math.PI };
}

// blocks is how many blocks of an aiming arrow's chain fit between its
// source and its head: one every step pixels, the head's length short
// of the pointer.
export function blocks(length: number, step: number, head: number): number {
  return Math.max(0, Math.floor((length - head) / step));
}
