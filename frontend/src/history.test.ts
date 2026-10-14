import { describe, expect, it } from 'vitest';
import { entriesFor, lookedAt, stepName } from './history';

const n = { player: (p: number) => ['Ann', 'Bo', 'Cy'].at(p) ?? '?', card: (ref: string) => ref.toUpperCase() };

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

  it('keeps the resolved stack item for the strip (15.22)', () => {
    const item = { id: 3, kind: 0, controller: 1, source: 7, card: 'a_dime', target: { player: 0 }, ability: { card: 'a_dime', index: 0 }, roll_for: { card: '', index: 0 } };
    const [entry] = entriesFor([{ kind: 'stack_resolved', player: 1, object: 7, card: 'a_dime', stack_id: 3, item }], 1, 0, n, 0);
    expect(entry.item).toEqual(item);
    expect(entriesFor([{ kind: 'looted', player: 0, card: 'a_dime' }], 1, 0, n, 0)[0]?.item).toBeUndefined();
  });

  it('shows a changed roll result next to its die', () => {
    const item = { id: 3, kind: 3, controller: 1, source: 0, roll: 6, attack: true, target: { player: 0 }, ability: { card: '', index: 0 }, roll_for: { card: '', index: 0 } };
    const [changed] = entriesFor(
      [
        { kind: 'stack_resolved', player: 1, amount: 6, text: 'attack roll', stack_id: 3, item },
        { kind: 'roll_resolved', player: 1, amount: 4, text: 'attack' },
      ],
      1,
      0,
      n,
      0,
    );
    expect(changed.text).toBe('Roll — attack roll (rolled 6 → 4)');
    expect(changed.item?.roll).toBe(4);
    expect(item.roll).toBe(6);
    const [same] = entriesFor(
      [
        { kind: 'stack_resolved', player: 1, amount: 6, text: 'attack roll' },
        { kind: 'roll_resolved', player: 1, amount: 6, text: 'attack' },
      ],
      1,
      0,
      n,
      0,
    );
    expect(same.text).toBe('Roll — attack roll (rolled 6)');
  });

  it('keeps hidden loots hidden', () => {
    const mine = entriesFor([{ kind: 'looted', player: 0, card: 'a_dime', private: true }], 1, 0, n, 0);
    const theirs = entriesFor([{ kind: 'looted', player: 1, private: true }], 1, 0, n, 0);
    expect(mine[0]?.text).toBe('You loot A_DIME');
    expect(theirs[0]?.text).toBe('Bo loots a card');
    expect(theirs[0]?.card).toBeUndefined();
  });

  it('shows a stolen card to the thief and the victim only (V-14)', () => {
    const shown = { kind: 'gave_card', player: 1, card: 'bomb', text: '0', private: true, seen_by: [0] };
    const hidden = { kind: 'gave_card', player: 1, text: '2', private: true, seen_by: [2] };
    expect(entriesFor([shown], 1, 1, n, 0)[0]?.text).toBe('You get BOMB from Ann');
    expect(entriesFor([shown], 1, 0, n, 0)[0]?.text).toBe('Bo gets your BOMB');
    const other = entriesFor([hidden], 1, 0, n, 0).at(0);
    expect(other?.text).toBe('Bo gets a card from Cy');
    expect(other?.card).toBeUndefined();
  });
});

describe('stepName', () => {
  it('names the turn steps', () => {
    expect(stepName(4)).toBe('Action');
    expect(stepName(99)).toBe('');
  });
});

describe('lookedAt', () => {
  it('gives the looker the cards, and others only that a card was seen (V-07)', () => {
    const events = [
      { kind: 'looked_at', player: 0, object: 9, card: 'a_dime', private: true },
      { kind: 'looked_at', player: 1, private: true },
      { kind: 'looted', player: 0, card: 'a_nickel' },
    ];
    expect(lookedAt(events, 0)).toEqual([{ id: 9, card: 'a_dime' }]);
    expect(lookedAt(events, 1)).toEqual([]);
    const mine = entriesFor(events, 1, 0, n, 0);
    expect(mine[0]).toMatchObject({ text: 'You look at A_DIME', card: 'a_dime' });
    expect(mine[1]?.text).toBe('Bo looks at a card');
  });
});
