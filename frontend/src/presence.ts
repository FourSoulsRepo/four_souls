// Presence (15.36, ADR 016): what this player points at goes to the
// others, and what they point at shows here. Pure, except spotElement
// and boardBox, which read the page.
import type { Point } from './aim';
import { Deck, type Presence, type Spot } from './protocol';

// What a sender may share (owner, 2026-10-20): an aim always; a hover
// only on the sender's own turn, so a waiting table stays quiet.
export interface Pointing {
  // source is the card the player aims with, if aiming.
  source?: number;
  // active: it is the player's turn.
  active: boolean;
}

// spotOf turns a zone under the pointer into a spot. A card in the
// player's own hand becomes its slot, never its object (A-16).
export function spotOf(zone: string, hand: number[]): Spot | undefined {
  const m = /^(obj|player|stack|deck)-(\w+)$/.exec(zone);
  if (!m?.[2]) {
    return undefined;
  }
  const n = Number(m[2]);
  switch (m[1]) {
    case 'obj':
      return cardSpot(n, hand);
    case 'player':
      return { player: n };
    case 'stack':
      return { stack_id: n };
    default:
      return m[2] in Deck ? { deck: Deck[m[2] as keyof typeof Deck] } : undefined;
  }
}

// cardSpot is a card's spot: its hand slot, or the card itself.
export function cardSpot(id: number, hand: number[]): Spot {
  const i = hand.indexOf(id);
  return i >= 0 ? { hand: i + 1 } : { object: id };
}

// spotAt is the first zone, topmost first, that is a spot.
export function spotAt(zones: string[], hand: number[]): Spot | undefined {
  for (const z of zones) {
    const s = spotOf(z, hand);
    if (s) {
      return s;
    }
  }
  return undefined;
}

// presenceFor is what the player shares now, from the zones under the
// pointer and the pointer on the board frame. A pointer over no spot
// counts only while aiming, for the arrow's head.
export function presenceFor(p: Pointing, zones: string[], hand: number[], point: Point | undefined): Omit<Presence, 'seat'> {
  if (p.source === undefined && !p.active) {
    return {};
  }
  const at = spotAt(zones, hand);
  if (p.source === undefined) {
    return at ? { at } : {};
  }
  const from = cardSpot(p.source, hand);
  if (at) {
    return { from, at };
  }
  return point ? { from, point } : { from };
}

// empty: the presence points at nothing.
export function empty(p: Omit<Presence, 'seat'>): boolean {
  return !p.from && !p.at && !p.point;
}

// same: two presences point at the same spots. Points close together
// count as the same, so a still pointer sends nothing new.
export function same(a: Omit<Presence, 'seat'>, b: Omit<Presence, 'seat'>): boolean {
  if (JSON.stringify(a.from) !== JSON.stringify(b.from) || JSON.stringify(a.at) !== JSON.stringify(b.at)) {
    return false;
  }
  if (!a.point || !b.point) {
    return !a.point && !b.point;
  }
  return Math.abs(a.point.x - b.point.x) < 0.002 && Math.abs(a.point.y - b.point.y) < 0.002;
}

// onFrame is a screen point on a box, 0 to 1 each way, kept inside it.
export function onFrame(x: number, y: number, box: { left: number; top: number; width: number; height: number }): Point {
  const clamp = (v: number) => Math.min(Math.max(v, 0), 1);
  return { x: clamp((x - box.left) / box.width), y: clamp((y - box.top) / box.height) };
}

// fromFrame is a point on a box back on the screen.
export function fromFrame(p: Point, box: { left: number; top: number; width: number; height: number }): Point {
  return { x: box.left + p.x * box.width, y: box.top + p.y * box.height };
}

// backSlot is the back a hand slot shows on: a hidden hand fans only
// so many backs, so a later slot shows on the last one.
export function backSlot(slot: number, backs: number): number {
  return Math.min(slot, backs) - 1;
}

// deckName is the zone name of a deck kind.
function deckName(kind: number): string | undefined {
  return Object.entries(Deck).find(([, v]) => v === kind)?.[0];
}

// spotElement finds a spot on this page; seat is the sender, whose
// hand a hand slot is in. A spot that is gone finds nothing.
export function spotElement(s: Spot | undefined, seat: number): Element | null {
  if (!s) {
    return null;
  }
  if (s.hand) {
    const backs = document.querySelectorAll(`[data-zone="hand-${String(seat)}"] .hidden-back`);
    return backs.item(backSlot(s.hand, backs.length));
  }
  if (s.player !== undefined) {
    return document.querySelector(`[data-player="${String(s.player)}"]`);
  }
  if (s.stack_id) {
    return document.querySelector(`[data-stack="${String(s.stack_id)}"]`);
  }
  if (s.deck !== undefined) {
    const name = deckName(s.deck);
    return name ? document.querySelector(`[data-zone="deck-${name}"]`) : null;
  }
  return s.object ? document.querySelector(`[data-zone="obj-${String(s.object)}"]`) : null;
}

// board is the board frame's box on screen, if drawn.
export function boardBox(): DOMRect | undefined {
  return document.querySelector('.board')?.getBoundingClientRect();
}
