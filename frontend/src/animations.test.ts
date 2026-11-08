import { beforeEach, describe, expect, it } from 'vitest';
import doc from '../../docs/wiki/Animations.md?raw';
import { arrivedBy, claim, claimEffect, claimed, clashes, nextBeat, resetClaims } from './animations';

describe('claims', () => {
  beforeEach(() => {
    resetClaims();
  });

  it('lets one animation bring a card, and the pop-in sees it', () => {
    expect(claim('pin', 'arrive', 'place')).toBe(true);
    expect(claimed('pin', 'arrive')).toBe(true);
    nextBeat();
    expect(claimed('pin', 'arrive')).toBe(false);
    expect(clashes).toEqual([]);
  });

  it('finds two animations that bring one card in one update', () => {
    expect(claimEffect({ kind: 'place', from: ['deck-monster'], to: ['monster-slot-0'], card: 'pin', back: 'monster' })).toBe(true);
    expect(claimEffect({ kind: 'fly', from: ['deck-monster'], to: ['monster-slot-0'], card: 'pin', back: 'monster' })).toBe(false);
    expect(clashes).toEqual([{ card: 'pin', job: 'arrive', first: 'place', second: 'fly' }]);
  });

  it('lets the same animation bring two copies, and a card arrive again later', () => {
    expect(claim('a_penny', 'arrive', 'fly')).toBe(true);
    expect(claim('a_penny', 'arrive', 'fly')).toBe(true);
    nextBeat();
    expect(claim('a_penny', 'arrive', 'place')).toBe(true);
    expect(clashes).toEqual([]);
  });

  it('finds a card shown large twice in one update', () => {
    expect(claimEffect({ kind: 'reveal', card: 'pin', deck: ['deck-monster'] })).toBe(true);
    expect(claimEffect({ kind: 'reveal', card: 'pin', deck: ['deck-monster'] })).toBe(false);
    expect(clashes).toEqual([{ card: 'pin', job: 'show', first: 'reveal', second: 'reveal' }]);
  });

  it('finds a card shown large in two updates in a row (19.4)', () => {
    expect(claimEffect({ kind: 'reveal', card: 'pin', deck: ['deck-monster'] })).toBe(true);
    nextBeat();
    expect(claimEffect({ kind: 'reveal', card: 'pin', deck: ['deck-monster'] })).toBe(false);
    nextBeat();
    nextBeat();
    expect(claimEffect({ kind: 'reveal', card: 'pin', deck: ['deck-monster'] })).toBe(true);
    expect(clashes).toEqual([{ card: 'pin', job: 'show', first: 'reveal', second: 'reveal' }]);
  });

  it('claims an own draw\'s card, and only that (20.7)', () => {
    expect(claimEffect({ kind: 'draw', from: ['deck-loot'], to: ['hand-0'], card: 'a_penny', own: true, pace: 1 })).toBe(true);
    expect(arrivedBy('a_penny')).toBe('loot-draw');
    expect(claimEffect({ kind: 'draw', from: ['deck-loot'], to: ['hand-1'], own: false, pace: 1 })).toBe(true);
    nextBeat();
    expect(arrivedBy('a_penny')).toBeUndefined();
    expect(clashes).toEqual([]);
  });

  // A card may fly to the loot discard and be drawn back from it in one
  // update: two arrivals at two places, no clash (19.134).
  it('lets a card fly to the discard and be drawn back in one update', () => {
    expect(claimEffect({ kind: 'fly', from: ['hand-1'], to: ['discard-loot'], card: 'three_cents', back: 'loot' })).toBe(true);
    expect(claimEffect({ kind: 'draw', from: ['discard-loot'], to: ['hand-0', 'player-0'], card: 'three_cents', own: true, pace: 1 })).toBe(true);
    expect(arrivedBy('three_cents')).toBe('loot-draw');
    expect(clashes).toEqual([]);
  });

  // Another player may play two copies of a loot card in two updates in
  // a row: two objects, two cards, no clash (19.136). One object shown
  // twice, or a show with no object, still clashes (19.4).
  it('lets two copies of a played card show in two updates in a row', () => {
    const played = (object?: number) => ({ kind: 'played' as const, card: 'two_cents', who: 'Bob', from: ['hand-1'], to: ['stack'], ...(object === undefined ? {} : { object }) });
    expect(claimEffect(played(512))).toBe(true);
    nextBeat();
    expect(claimEffect(played(514))).toBe(true);
    expect(clashes).toEqual([]);
    nextBeat();
    expect(claimEffect(played(514))).toBe(false);
    nextBeat();
    nextBeat();
    expect(claimEffect(played(520))).toBe(true);
    nextBeat();
    expect(claimEffect(played())).toBe(false);
    expect(clashes).toEqual([
      { card: 'two_cents', job: 'show', first: 'played', second: 'played' },
      { card: 'two_cents', job: 'show', first: 'played', second: 'played' },
    ]);
  });

  it('claims nothing for effects that bring no card', () => {
    expect(claimEffect({ kind: 'dice', value: 3, label: '', player: 0 })).toBe(true);
    expect(claimEffect({ kind: 'fly', from: ['deck-loot'], to: ['hand-1'], back: 'loot' })).toBe(true);
    expect(clashes).toEqual([]);
  });
});

// The inventory: every animation in the code carries an "anim: <id>"
// tag, and docs/wiki/Animations.md has a row for each id (19.6).
describe('a choice deal, then a reveal (19.139)', () => {
  it('lets a card dealt in a look be revealed in the next update', () => {
    resetClaims();
    expect(claim('stoney', 'show', 'choice-deal')).toBe(true);
    nextBeat();
    expect(claim('stoney', 'show', 'reveal')).toBe(true);
    expect(clashes).toEqual([]);
  });

  it('still finds a reveal in the same update, or a deal right after a reveal', () => {
    resetClaims();
    claim('stoney', 'show', 'choice-deal');
    expect(claim('stoney', 'show', 'reveal')).toBe(false);
    resetClaims();
    claim('fly', 'show', 'reveal');
    nextBeat();
    expect(claim('fly', 'show', 'choice-deal')).toBe(false);
  });
});

describe('animation inventory', () => {
  const sources = import.meta.glob<string>(['./**/*.{ts,tsx,css}', '!./**/*.test.ts'], {
    query: '?raw',
    import: 'default',
    eager: true,
  });
  const site = /@keyframes|\.animate\(|<motion\./;
  const tag = /anim: ([a-z-]+)/;

  function codeIds(): { ids: Set<string>; untagged: string[] } {
    const ids = new Set<string>();
    const untagged: string[] = [];
    for (const [file, text] of Object.entries(sources)) {
      const lines = text.split('\n');
      lines.forEach((line, i) => {
        if (!site.test(line)) {
          return;
        }
        const m = tag.exec(lines[i - 1] ?? '');
        if (m?.[1]) {
          ids.add(m[1]);
        } else {
          untagged.push(`${file}:${String(i + 1)}`);
        }
      });
    }
    return { ids, untagged };
  }

  // docIds are the ids in the first column of the inventory table.
  function docIds(): Set<string> {
    const table = doc.split('## Inventory')[1]?.split('\n## ')[0] ?? '';
    return new Set([...table.matchAll(/^\| `([a-z-]+)` \|/gm)].map((m) => m[1]));
  }

  it('tags every animation in the code', () => {
    expect(codeIds().untagged).toEqual([]);
  });

  it('lists every tagged animation in Animations.md, and no other', () => {
    expect([...docIds()].sort()).toEqual([...codeIds().ids].sort());
  });
});
