// Playing a card from the own hand (20.20, ADR 027 "Playing from the
// hand"): the plain decisions, with no page. HandCard.tsx runs them.
import { feel } from './feel';
import { px } from './uiScale';

// dragged: a press moved far enough to pick the card up.
export function dragged(dx: number, dy: number): boolean {
  return Math.hypot(dx, dy) >= px(feel.play.drag);
}

// Release is what a click does to a card held by an earlier click or
// tap: play it there, or put it back.
export type Release = 'play' | 'back';

// release decides a click while a card is held: over the board it
// plays, unless it lands on a button (End turn is on the board); a click
// anywhere else puts the card back.
export function release(board: boolean, button: boolean): Release {
  return board && !button ? 'play' : 'back';
}

// Answer is what the server's answer means for a card played and held
// where it was let go: wait on, go home, or leave the hand.
export type Answer = 'wait' | 'home' | 'refused';

export interface Waiting {
  // step is the game's step when the card was played.
  step: number;
  // step and sent are the store's now; sent is null once refused.
  now: number;
  sent: number | null;
  // more: the play's answer goes on. The card asks its target next and
  // the target was aimed already (15.25), so the pick follows at once.
  more: boolean;
}

// answer tells a held card what to do while it is still in the hand: a
// new step is the answer; the card left the hand if it took it with it.
// A refusal keeps the step and clears sent.
export function answer(w: Waiting): Answer {
  if (w.now !== w.step) {
    return w.more ? 'wait' : 'home';
  }
  return w.sent === null ? 'refused' : 'wait';
}

// gapMoves are the shifts of the hand cards that stay when one or more
// leave (13 in ADR 027 "Playing from the hand"): each moves from its old
// place to its new one, so the fan closes the gap. before and after are
// each card's left edge by its id. Nothing moves when no card left, or
// for a card that did not move.
export function gapMoves(before: ReadonlyMap<number, number>, after: ReadonlyMap<number, number>): { id: number; dx: number }[] {
  if (![...before.keys()].some((id) => !after.has(id))) {
    return [];
  }
  const out: { id: number; dx: number }[] = [];
  for (const [id, x] of after) {
    const old = before.get(id);
    if (old !== undefined && old !== x) {
      out.push({ id, dx: old - x });
    }
  }
  return out;
}

export interface Rect {
  x: number;
  y: number;
  w: number;
  h: number;
}

// Into is the shrink of a played card into its stack row: the card,
// drawn at its natural size with its top left corner at the held
// card's, scales from `from` and moves by (dx, dy) as it scales to `to`.
export interface Into {
  from: number;
  to: number;
  dx: number;
  dy: number;
}

// into is the path of a played card from where it was held (card, as
// drawn) into its stack row (row), as one scale about its top left
// corner: it ends as tall as the row, at the row's middle. With no row
// it shrinks to half where it is.
export function into(card: Rect, natural: { w: number; h: number }, row: Rect | null): Into {
  const from = card.w / natural.w;
  const to = row ? row.h / natural.h : from / 2;
  const cx = row ? row.x + row.w / 2 : card.x + card.w / 2;
  const cy = row ? row.y + row.h / 2 : card.y + card.h / 2;
  return { from, to, dx: cx - (natural.w * to) / 2 - card.x, dy: cy - (natural.h * to) / 2 - card.y };
}
