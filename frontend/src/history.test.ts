import { describe, expect, it } from 'vitest';
import { entriesFor, stepName } from './history';

const n = { player: (p: number) => ['Ann', 'Bo'].at(p) ?? '?', card: (ref: string) => ref.toUpperCase() };

describe('entriesFor', () => {
  it('makes one entry per resolved stack item and marks turns', () => {
    const entries = entriesFor(
      [
        { kind: 'turn_started', player: 1, amount: 4 },
        { kind: 'passed_priority', player: 0 },
        { kind: 'stack_resolved', player: 1, object: 7, card: 'a_dime' },
        { kind: 'stack_resolved', player: 1, object: 8, amount: 5, text: 'attack roll' },
      ],
      3,
      0,
      n,
      10,
    );
    expect(entries).toEqual([
      { key: 10, turn: 4, icon: '▶', text: 'Turn 4: Bo' },
      { key: 11, turn: 4, icon: '✦', text: 'A_DIME', card: 'a_dime' },
      { key: 12, turn: 4, icon: '✦', text: 'Roll — attack roll (rolled 5)' },
    ]);
  });

  it('keeps hidden loots hidden', () => {
    const mine = entriesFor([{ kind: 'looted', player: 0, card: 'a_dime', private: true }], 1, 0, n, 0);
    const theirs = entriesFor([{ kind: 'looted', player: 1, private: true }], 1, 0, n, 0);
    expect(mine[0]?.text).toBe('You loot A_DIME');
    expect(theirs[0]?.text).toBe('Bo loots a card');
    expect(theirs[0]?.card).toBeUndefined();
  });
});

describe('stepName', () => {
  it('names the turn steps', () => {
    expect(stepName(4)).toBe('Action');
    expect(stepName(99)).toBe('');
  });
});
