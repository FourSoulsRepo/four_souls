// Animation claims (19.6, ADR 012). Every animation that brings a card
// to its place or shows it large claims the card first. Two animations
// that claim one card for the same job clash: the card would appear or
// show twice. A clash is logged and kept in `clashes`; e2e tests fail
// on any (e2e/tests/fixtures.ts). docs/wiki/Animations.md lists every
// animation and what it claims.
//
// Time is counted in beats: a beat is one update, from its first effect
// until the next update's effects start, so it covers the cards its view
// mounts too.
import type { Effect } from './playback';

// arrive: the card comes to its place on the table. show: it shows
// large in the middle of the screen.
export type Job = 'arrive' | 'show';

export interface Clash {
  card: string;
  job: Job;
  // first and second name the two animations, as in Animations.md.
  first: string;
  second: string;
}

// Where tells claims of one card apart, when known. to: where an
// arrival brings the card (19.134). object: which copy of the card, or
// which turn of it on the table: the engine gives a card a new object
// at each move (19.136, 19.191).
export interface Where {
  to?: string | undefined;
  object?: number | undefined;
}

interface Claim extends Where {
  by: string;
  beat: number;
}

let beat = 0;
// claims keeps every live claim of a card for a job, oldest first, so
// a claim of one copy never hides another copy's (19.191).
const claims = new Map<string, Claim[]>();
export const clashes: Clash[] = [];

// How many beats back a claim still clashes: a card shown large in the
// update before is not shown again (19.4). A choice's deal shows the
// cards one player looks at, in a window, from the deck: a reveal of one
// in a later update is a public event of its own (19.139).
function reach(job: Job, by: string): number {
  return job === 'show' && by !== 'choice-deal' ? 1 : 0;
}

// nextBeat starts the next update's animations.
export function nextBeat(): void {
  beat++;
  for (const [key, list] of claims) {
    const live = list.filter((c) => beat - c.beat <= 1);
    if (live.length === 0) {
      claims.delete(key);
    } else {
      claims.set(key, live);
    }
  }
}

// claim records that animation `by` does `job` to card. It returns false,
// and records a clash, when another animation already did the same job
// to it within reach. One animation may bring a card again: two copies
// of a card, or a card looted and discarded in one update. Arrivals at
// two known, different places are two jobs: a card flies to the loot
// discard and is drawn back from it in one update (19.134). Two known,
// different objects are two cards: two copies of a loot card played in
// two updates in a row (19.136); or one card revealed on its deck, then
// revealed again into a slot after it moved (19.191).
export function claim(card: string, job: Job, by: string, where: Where = {}): boolean {
  const key = `${job}:${card}`;
  const list = claims.get(key) ?? [];
  const old = list.find((c) => clashesWith(c, job, by, where));
  list.push({ by, beat, ...where });
  claims.set(key, list);
  if (!old) {
    return true;
  }
  const clash = { card, job, first: old.by, second: by };
  clashes.push(clash);
  console.error(`animation clash: ${by} ${job}s ${card} after ${old.by}`);
  return false;
}

// clashesWith tells if a new claim clashes with an older one of the
// same card and job.
function clashesWith(old: Claim, job: Job, by: string, where: Where): boolean {
  if (beat - old.beat > reach(job, old.by) || (job === 'arrive' && old.by === by)) {
    return false;
  }
  if (job === 'arrive' && apart(where.to, old.to)) {
    return false;
  }
  return !apart(where.object, old.object);
}

// apart: both are known and differ.
function apart<T>(a: T | undefined, b: T | undefined): boolean {
  return a !== undefined && b !== undefined && a !== b;
}

// claimEffect claims the card an effect brings or shows; false means
// it clashes and is skipped.
export function claimEffect(e: Effect): boolean {
  switch (e.kind) {
    case 'fly':
    case 'place':
      return e.card === undefined || claim(e.card, 'arrive', e.kind, { to: e.to[0], object: e.object });
    case 'draw':
      // An own draw brings its card to the hand (20.7); another
      // player's card is never seen.
      return !e.own || e.card === undefined || claim(e.card, 'arrive', 'loot-draw', { to: e.to[0], object: e.object });
    case 'reveal':
      return claim(e.card, 'show', 'reveal', { object: e.object });
    case 'played':
      return claim(e.card, 'show', 'played', { object: e.object });
    default:
      return true;
  }
}

// dealt keeps, per prompt step, the cards a deal claimed, so a second
// mount of the same choice (React's StrictMode) does not claim again.
const dealt = new Map<number, Set<string>>();

// dealKey names one dealt card: its object, when known, so two copies
// in one deal are two cards (19.191); else its name.
export function dealKey(card: string, object: number | undefined): string {
  return object === undefined ? card : `${card}#${String(object)}`;
}

// claimDeal claims each card a choice's deal shows (`choice-deal`): it
// returns the keys (dealKey) of the cards that may fly; a card another
// animation showed just now comes at once instead. A card's object,
// when known, tells copies apart (19.141).
export function claimDeal(step: number, cards: { card: string; object?: number | undefined }[]): Set<string> {
  let ok = dealt.get(step);
  if (!ok) {
    const each = new Map(cards.map((c) => [dealKey(c.card, c.object), c]));
    ok = new Set([...each].filter(([, c]) => claim(c.card, 'show', 'choice-deal', { object: c.object })).map(([key]) => key));
    dealt.clear();
    dealt.set(step, ok);
  }
  return ok;
}

// forgetDeals forgets the deals of a left game: the next game counts its
// steps from 0 again, and an old deal at the same step would keep its
// cards from flying in (19.227).
export function forgetDeals(): void {
  dealt.clear();
}

// claimed tells if an animation already brought the card in this beat:
// the view's pop-in then stays still. It keys on the card's name: a
// copy that no effect brought then stays still too, which shows nothing
// twice.
export function claimed(card: string, job: Job): boolean {
  return (claims.get(`${job}:${card}`) ?? []).some((c) => c.beat === beat);
}

// arrivedBy is the animation that last brought the card in this beat,
// if any.
export function arrivedBy(card: string): string | undefined {
  return (claims.get(`arrive:${card}`) ?? []).findLast((c) => c.beat === beat)?.by;
}

// resetClaims forgets every claim and clash, for tests.
export function resetClaims(): void {
  beat = 0;
  claims.clear();
  dealt.clear();
  clashes.length = 0;
}
