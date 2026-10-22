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

interface Claim {
  by: string;
  beat: number;
}

let beat = 0;
const claims = new Map<string, Claim>();
export const clashes: Clash[] = [];

// How many beats back a claim still clashes: a card shown large in the
// update before is not shown again (19.4).
function reach(job: Job): number {
  return job === 'show' ? 1 : 0;
}

// nextBeat starts the next update's animations.
export function nextBeat(): void {
  beat++;
  for (const [key, c] of claims) {
    if (beat - c.beat > 1) {
      claims.delete(key);
    }
  }
}

// claim records that animation `by` does `job` to card. It returns false,
// and records a clash, when another animation already did the same job
// to it within reach. One animation may bring a card again: two copies
// of a card, or a card looted and discarded in one update.
export function claim(card: string, job: Job, by: string): boolean {
  const key = `${job}:${card}`;
  const old = claims.get(key);
  claims.set(key, { by, beat });
  if (!old || beat - old.beat > reach(job) || (job === 'arrive' && old.by === by)) {
    return true;
  }
  const clash = { card, job, first: old.by, second: by };
  clashes.push(clash);
  console.error(`animation clash: ${by} ${job}s ${card} after ${old.by}`);
  return false;
}

// claimEffect claims the card an effect brings or shows; false means
// it clashes and is skipped.
export function claimEffect(e: Effect): boolean {
  switch (e.kind) {
    case 'fly':
    case 'place':
      return e.card === undefined || claim(e.card, 'arrive', e.kind);
    case 'reveal':
      return claim(e.card, 'show', 'reveal');
    case 'played':
      return claim(e.card, 'show', 'played');
    default:
      return true;
  }
}

// claimed tells if an animation already brought the card in this beat:
// the view's pop-in then stays still.
export function claimed(card: string, job: Job): boolean {
  return claims.get(`${job}:${card}`)?.beat === beat;
}

// resetClaims forgets every claim and clash, for tests.
export function resetClaims(): void {
  beat = 0;
  claims.clear();
  clashes.length = 0;
}
