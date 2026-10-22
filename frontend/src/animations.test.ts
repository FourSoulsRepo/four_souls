import { beforeEach, describe, expect, it } from 'vitest';
import doc from '../../docs/wiki/Animations.md?raw';
import { claim, claimEffect, claimed, clashes, nextBeat, resetClaims } from './animations';

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

  it('claims nothing for effects that bring no card', () => {
    expect(claimEffect({ kind: 'dice', value: 3, label: '', player: 0 })).toBe(true);
    expect(claimEffect({ kind: 'fly', from: ['deck-loot'], to: ['hand-1'], back: 'loot' })).toBe(true);
    expect(clashes).toEqual([]);
  });
});

// The inventory: every animation in the code carries an "anim: <id>"
// tag, and docs/wiki/Animations.md has a row for each id (19.6).
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
