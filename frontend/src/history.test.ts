import { describe, expect, it } from 'vitest';
import kindList from '../../pkg/rules_engine/testdata/event_kinds.txt?raw';
import { entriesFor, historyKinds, lookedAt, stepName } from './history';
import type { GameEvent } from './protocol';

const n = { player: (p: number) => ['Ann', 'Bo', 'Cy'].at(p) ?? '?', card: (ref: string) => ref.toUpperCase() };

describe('entriesFor', () => {
  it('keeps every first-player die, a tie and who goes first (15.43)', () => {
    const entries = entriesFor(
      [
        { kind: 'dice_rolled', player: 0, amount: 2, text: 'first player roll' },
        { kind: 'dice_rolled', player: 1, amount: 2, text: 'first player roll' },
        { kind: 'first_roll_tied', player: -1, amount: 2, text: '0,1' },
        { kind: 'dice_rolled', player: 1, amount: 1, text: 'first player roll' },
        { kind: 'first_player', player: 1, amount: 1 },
      ],
      0,
      0,
      n,
      0,
    );
    expect(entries.map((e) => e.text)).toEqual([
      'You roll 2 for the first player',
      'Bo rolls 2 for the first player',
      'Tie at 2: You and Bo roll again',
      'Bo rolls 1 for the first player',
      'Bo goes first — rolled 1',
    ]);
    expect(entriesFor([{ kind: 'first_player', player: 2, card: 'cain', text: 'character' }], 1, 0, n, 0)).toEqual([
      { key: 0, turn: 1, icon: '①', text: 'Cy goes first — CAIN', card: 'cain' },
    ]);
    // Other rolls are in their stack item's entry.
    expect(entriesFor([{ kind: 'dice_rolled', player: 0, amount: 4, text: 'attack roll' }], 1, 0, n, 0)).toEqual([]);
  });

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
      { key: 11, turn: 4, icon: '✦', text: 'A_DIME', card: 'a_dime', headline: true },
      { key: 12, turn: 4, icon: '✦', text: 'Roll — attack roll (rolled 5)', headline: true },
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

  it('shows a card looked at twice in one update once (19.13)', () => {
    const events = [
      { kind: 'looked_at', player: 0, object: 361, card: 'library_card', private: true },
      { kind: 'looked_at', player: 0, object: 361, card: 'library_card', private: true },
      { kind: 'looked_at', player: 0, object: 430, card: 'a_dime', private: true },
    ];
    expect(lookedAt(events, 0)).toEqual([
      { id: 361, card: 'library_card' },
      { id: 430, card: 'a_dime' },
    ]);
  });
});

describe('historyKinds', () => {
  const kinds = kindList.split('\n').filter(Boolean);
  // sample: an event of a told kind that gives a line; others need only
  // their kind and a player.
  const sample = new Map<string, Partial<GameEvent>>([['dice_rolled', { text: 'first player roll', amount: 3 }]]);

  it('tells or skips every engine event kind, with a reason (15.88)', () => {
    expect(Object.keys(historyKinds).sort()).toEqual(kinds);
    for (const [kind, why] of Object.entries(historyKinds)) {
      expect(why, kind).not.toBe('');
    }
  });

  it('gives a line for every told kind and none for a skipped one', () => {
    for (const [kind, why] of Object.entries(historyKinds)) {
      const lines = entriesFor([{ kind, player: 1, card: 'gaper', amount: 1, ...sample.get(kind) }], 1, 0, n, 0);
      if (why === 'line') {
        expect(lines.length, kind).toBeGreaterThan(0);
      } else if (kind !== 'roll_resolved') {
        expect(lines, kind).toEqual([]);
      }
    }
  });

  it('names who attacks, buys and kills, as headlines', () => {
    const lines = entriesFor(
      [
        { kind: 'attack_declared', player: 1 },
        { kind: 'attack_target', player: 1, object: 5, card: 'gaper' },
        { kind: 'died', player: -1, object: 9, card: 'gaper', prev: 5, by: 2 },
        { kind: 'rewards_gained', player: 1, object: 9, card: 'gaper' },
        { kind: 'gained_cents', player: 1, amount: 3 },
        { kind: 'purchase_declared', player: 1 },
        { kind: 'purchased', player: 1, object: 12, card: 'spoon_bender' },
        { kind: 'gained_item', player: 1, object: 12, card: 'spoon_bender' },
      ],
      1,
      0,
      n,
      0,
    );
    expect(lines.map((e) => [e.text, e.headline ?? false])).toEqual([
      ['Bo declares an attack', false],
      ['Bo attacks GAPER', false],
      ['Bo kills GAPER', true],
      ['Bo gains the rewards of GAPER', false],
      ['Bo gains 3¢', false],
      ['Bo declares a purchase', false],
      ['Bo buys SPOON_BENDER', true],
    ]);
    expect(lines[6]?.card).toBe('spoon_bender');
    expect(lines.map((e) => e.badge)).toEqual([undefined, undefined, 'kill', undefined, undefined, undefined, 'buy']);
  });

  it('turns a resolved death into the kill, once, with the dead card (15.89)', () => {
    const none = { card: '', index: 0 };
    const death = { id: 7, kind: 5, controller: -1, source: 0, text: 'death', target: { player: -1, object: 5 }, ability: none, roll_for: none };
    const lines = entriesFor(
      [
        { kind: 'stack_resolved', player: -1, text: 'death', stack_id: 7, item: death },
        { kind: 'died', player: -1, object: 9, card: 'gaper', prev: 5, by: 2 },
      ],
      1,
      0,
      n,
      0,
    );
    expect(lines).toEqual([{ key: 0, turn: 1, icon: '☠', text: 'Bo kills GAPER', card: 'gaper', item: death, headline: true, badge: 'kill' }]);
    // A player's death item ends only that player's death.
    const mine = { ...death, target: { player: 0, is_player: true } };
    const two = entriesFor(
      [
        { kind: 'stack_resolved', player: -1, text: 'death', stack_id: 7, item: mine },
        { kind: 'died', player: 1, object: 2, card: 'cain' },
      ],
      1,
      0,
      n,
      0,
    );
    expect(two.map((e) => [e.text, e.badge])).toEqual([
      ['Roll — death', undefined],
      ['Bo dies', 'kill'],
    ]);
  });

  it('gives a bought treasure deck top one line, though its gain comes first', () => {
    const lines = entriesFor(
      [
        { kind: 'gained_treasure', player: 0, object: 4, card: 'the_d6' },
        { kind: 'gained_item', player: 0, object: 4, card: 'the_d6' },
        { kind: 'entered_play', player: 0, object: 4, card: 'the_d6' },
        { kind: 'purchased', player: 0, object: 4, card: 'the_d6', text: 'treasure deck' },
      ],
      1,
      0,
      n,
      0,
    );
    expect(lines.map((e) => e.text)).toEqual(['You buy THE_D6 from the treasure deck']);
  });

  it('names gains, thefts and losses once each', () => {
    const lines = entriesFor(
      [
        { kind: 'gained_treasure', player: 0, object: 4, card: 'the_d6' },
        { kind: 'gained_item', player: 0, object: 4, card: 'the_d6' },
        { kind: 'gained_treasure', player: 1, object: 6, card: 'sack', text: 'stolen' },
        { kind: 'stole', player: 1, amount: 2, text: '0' },
        { kind: 'gave_cents', player: 0, amount: 1, by: 3 },
        { kind: 'lost_cents', player: 2, amount: 4 },
        { kind: 'gained_soul', player: 2, object: 8, card: 'gaper' },
        { kind: 'gained_soul', player: 0, object: 9, card: 'greed', text: 'bonus soul' },
        { kind: 'purchase_failed', player: 0, amount: 10 },
      ],
      1,
      0,
      n,
      0,
    );
    expect(lines.map((e) => e.text)).toEqual([
      'You gain THE_D6',
      'Bo steals SACK',
      'Bo steals 2¢ from you',
      'You give 1¢ to Cy',
      'Cy loses 4¢',
      'Cy gains the soul of GAPER',
      'You gain the bonus soul GREED',
      "You can't buy: the price is 10¢",
    ]);
  });

  it('names deaths, a destroyed item once, fizzles and cancels', () => {
    const lines = entriesFor(
      [
        { kind: 'died', player: 1, object: 2, card: 'cain', by: 1 },
        { kind: 'died', player: 0, object: 3, card: 'isaac' },
        { kind: 'object_died', player: -1, object: 20, card: 'guppys_head', prev: 7, text: 'item', by: 2 },
        { kind: 'destroyed', player: 1, object: 20, card: 'guppys_head', prev: 7 },
        { kind: 'destroyed', player: 1, object: 21, card: 'sack', prev: 8 },
        { kind: 'fizzled', player: 1, object: 5, card: 'bomb' },
        { kind: 'cancelled', player: 1, object: 5, card: 'bomb' },
        { kind: 'penalty_paid', player: 1 },
      ],
      1,
      0,
      n,
      0,
    );
    expect(lines.map((e) => e.text)).toEqual([
      'You kill Bo',
      'You die',
      'Bo destroys GUPPYS_HEAD',
      'SACK is destroyed',
      'BOMB fizzles: its target is gone',
      'BOMB is cancelled',
      'Bo pays the death penalty',
    ]);
  });

  it("shows a stack card's discard only as its resolved line", () => {
    const events = [
      { kind: 'discarded', player: 1, object: 30, card: 'bomb', text: 'stack' },
      { kind: 'discarded', player: 1, object: 31, card: 'a_penny' },
    ];
    expect(entriesFor(events, 1, 0, n, 0).map((e) => e.text)).toEqual(['Bo discards A_PENNY']);
  });
});

