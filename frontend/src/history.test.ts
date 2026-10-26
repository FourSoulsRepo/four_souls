import { describe, expect, it } from 'vitest';
import kindList from '../../pkg/rules_engine/testdata/event_kinds.txt?raw';
import { PromptKind } from './actions';
import {
  entriesFor as entries,
  historyKinds,
  lookedAt,
  restartCut,
  rollbackButtons,
  rollbackTargets,
  rolledBack,
  stepName,
  stepText,
  stripEntries,
  voteEnded,
  type Entry,
} from './history';
import type { GameEvent, GameView, StackItem } from './protocol';
import { StackKind, stackLine } from './stack';

const n = { player: (p: number) => ['Ann', 'Bo', 'Cy'].at(p) ?? '?', card: (ref: string) => ref.toUpperCase() };

const player = (id: number, card: string) => ({
  id,
  character: { id: 100 + id, card, charged: true },
  hp: 2,
  atk: 1,
  cents: 3,
  souls: 0,
  in_play: null,
  hand_size: 3,
  loot_plays: 0,
  attacks: 0,
  purchases: 0,
});

const view: GameView = {
  players: [player(0, 'isaac'), player(1, 'cain'), player(2, 'eve')],
  shop: null,
  monsters: [{ top: { id: 9, card: 'gaper', charged: true } }],
  deck_sizes: [1, 1, 1, 0],
  discards: [null, null, null, null],
  stack: null,
  turn: { active: 1, number: 1, step: 4, loot_plays: 0, attacks: 0, purchases: 0 },
  waiting: { kind: PromptKind.priority, player: 1 },
};

// entriesFor is the history's, with the table above.
const entriesFor = (events: GameEvent[], turn: number, you: number, names: typeof n, key: number) => entries(events, turn, you, names, key, view);

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
    expect(changed.text).toBe('Attack: Bo rolls 6 → 4');
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

  it('names what changed a roll (15.94)', () => {
    const item = { id: 3, kind: 3, controller: 1, source: 0, roll: 1, attack: true, target: { player: 0 }, ability: { card: '', index: 0 }, roll_for: { card: '', index: 0 } };
    const meat = { stat: 'attack_roll' as const, source: 5, card: 'meat', amount: 1, lasts: 'in_play' as const };
    const [changed] = entriesFor(
      [
        { kind: 'stack_resolved', player: 1, amount: 1, text: 'attack roll', stack_id: 3, item },
        { kind: 'roll_resolved', player: 1, amount: 2, text: 'attack', sources: [meat] },
      ],
      1,
      0,
      n,
      0,
    );
    expect(changed.text).toBe(`Attack: Bo rolls 1 → 2 (+1 ${n.card('meat')})`);
  });

  it("says what a card-less stack item did, in the strip's words (19.52)", () => {
    const none = { card: '', index: 0 };
    const base = { id: 7, controller: 1, source: 0, ability: none, roll_for: none };
    const gaper = { player: -1, object: 9 };
    const items: StackItem[] = [
      { ...base, kind: StackKind.damage, amount: 2, label: 'damage', target: gaper },
      { ...base, kind: StackKind.damage, controller: -1, amount: 1, label: 'damage', source_card: 'gaper', target: { player: 1, is_player: true } },
      { ...base, kind: StackKind.death, controller: -1, label: 'death', target: gaper },
      { ...base, kind: StackKind.penalty, controller: -1, label: 'death penalty', target: { player: 1, is_player: true } },
      { ...base, kind: StackKind.deathStep, controller: -1, label: 'soul', target: gaper },
      { ...base, kind: StackKind.roll, roll: 5, attack: true, label: 'attack roll', target: gaper },
      { ...base, kind: StackKind.roll, roll: 3, label: 'roll', target: { player: -1 } },
    ];
    for (const item of items) {
      const [entry] = entriesFor([{ kind: 'stack_resolved', player: item.controller, amount: item.roll, text: item.label, stack_id: 7, item }], 1, 0, n, 0);
      expect(entry.text).toBe(stackLine(item, view, n));
      expect(entry.text).not.toMatch(/^Roll/);
    }
    expect(entriesFor([{ kind: 'stack_resolved', player: 1, text: 'damage', item: items[0] }], 1, 0, n, 0)[0]?.text).toBe('Bo deals 2 damage to GAPER');
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
    expect(lines).toEqual([
      {
        key: 0,
        turn: 1,
        icon: '☠',
        text: 'Bo kills GAPER',
        card: 'gaper',
        item: death,
        headline: true,
        badge: 'kill',
        from: 'cain',
        outcome: [{ text: 'GAPER', card: 'gaper', mark: 'death' }],
      },
    ]);
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
      ['Ann dies', undefined],
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

describe('outcome', () => {
  const none = { card: '', index: 0 };
  const bomb: StackItem = { id: 3, kind: StackKind.loot, controller: 1, source: 40, card: 'bomb', target: { player: -1, object: 9 }, ability: { card: 'bomb', index: 0 }, roll_for: none };

  it('shows a played card, its target and what it gave right away (15.90)', () => {
    const penny: StackItem = { ...bomb, card: 'a_penny', target: { player: 1, is_player: true } };
    const [play] = entriesFor(
      [
        { kind: 'stack_resolved', player: 1, object: 40, card: 'a_penny', stack_id: 3, item: penny },
        { kind: 'gained_cents', player: 1, amount: 1 },
        { kind: 'looted', player: 1, private: true },
      ],
      1,
      0,
      n,
      0,
    );
    expect(play.from).toBe('a_penny');
    expect(play.outcome).toEqual([
      { text: 'Bo', card: 'cain' },
      { text: '+1¢', mark: 'cents', amount: 1 },
      { text: 'a loot card' },
    ]);
  });

  it('gives a later line to the newest headline: a trigger in the middle', () => {
    const trigger: StackItem = { ...bomb, id: 4, kind: StackKind.trigger, card: 'greed', target: { player: 0, is_player: true } };
    const lines = entriesFor(
      [
        { kind: 'stack_resolved', player: 1, object: 40, card: 'bomb', stack_id: 3, item: bomb },
        { kind: 'stack_added', player: 1 },
        { kind: 'gained_cents', player: 1, amount: 2 },
        { kind: 'stack_resolved', player: 1, object: 41, card: 'greed', stack_id: 4, item: trigger },
        { kind: 'lost_cents', player: 0, amount: 3 },
      ],
      1,
      0,
      n,
      0,
    );
    const [play, trig] = lines.filter((e) => e.headline);
    expect(play.outcome).toEqual([
      { text: 'GAPER', card: 'gaper' },
      { text: '+2¢', mark: 'cents', amount: 2 },
    ]);
    expect(trig.outcome).toEqual([
      { text: 'Ann', card: 'isaac' },
      { text: '−3¢', mark: 'cents', amount: -3 },
    ]);
    // A new update starts with no headline: its gains go nowhere.
    expect(entriesFor([{ kind: 'gained_cents', player: 1, amount: 1 }], 1, 0, n, 0)[0]?.outcome).toBeUndefined();
  });

  it('shows damage from its dealer, a kill from its killer with the rewards, a buy from its buyer', () => {
    const hit: StackItem = { ...bomb, id: 5, kind: StackKind.damage, card: undefined, source: 0, amount: 2, target: { player: -1, object: 9 } };
    const [damage] = entriesFor([{ kind: 'stack_resolved', player: 1, text: 'damage', stack_id: 5, item: hit }], 1, 0, n, 0);
    expect(damage.from).toBe('cain');
    expect(damage.outcome).toEqual([{ text: 'GAPER', card: 'gaper', mark: 'damage', amount: 2 }]);

    const death: StackItem = { ...bomb, id: 6, kind: StackKind.death, controller: -1, card: undefined, source: 0, target: { player: -1, object: 9 } };
    const kill = entriesFor(
      [
        { kind: 'stack_resolved', player: -1, text: 'death', stack_id: 6, item: death },
        { kind: 'died', player: -1, object: 50, card: 'gaper', prev: 9, by: 2 },
        { kind: 'rewards_gained', player: 1, object: 50, card: 'gaper' },
        { kind: 'gained_cents', player: 1, amount: 1 },
        { kind: 'gained_soul', player: 1, object: 50, card: 'gaper' },
      ],
      1,
      0,
      n,
      0,
    ).find((e) => e.badge === 'kill');
    expect(kill?.from).toBe('cain');
    expect(kill?.outcome).toEqual([
      { text: 'GAPER', card: 'gaper', mark: 'death' },
      { text: '+1¢', mark: 'cents', amount: 1 },
      { text: 'GAPER', card: 'gaper' },
    ]);

    const [buy] = entriesFor(
      [
        { kind: 'purchased', player: 2, object: 12, card: 'spoon_bender' },
        { kind: 'gained_item', player: 2, object: 12, card: 'spoon_bender' },
      ],
      1,
      0,
      n,
      0,
    );
    expect(buy.from).toBe('eve');
    expect(buy.outcome).toEqual([{ text: 'SPOON_BENDER', card: 'spoon_bender' }]);
  });
});


describe('stripEntries', () => {
  it('keeps every tile key when the history grows (15.51)', () => {
    const first = entriesFor([{ kind: 'purchased', player: 1, object: 12, card: 'spoon_bender' }], 1, 0, n, 0);
    const before = stripEntries(first).map((e) => e.key);
    const next = entriesFor([{ kind: 'died', player: -1, object: 9, card: 'gaper', prev: 5, by: 2 }], 1, 0, n, first.length);
    const after = stripEntries([...first, ...next]).map((e) => e.key);
    // The new tile is on top; the old ones keep their keys below it.
    expect(after.slice(1)).toEqual(before);
    expect(after).toHaveLength(before.length + 1);
  });
});

describe('rollback (ADR 020)', () => {
  // Steps 1 to 3 brought two rows, none, and two; the start brought one.
  const row = (key: number, step: number, text: string): Entry => ({ key, turn: 1, icon: '+', text, step });
  const history = [row(1, 0, 'start'), row(2, 1, 'Ann loots'), row(3, 1, 'Ann gains 1¢'), row(4, 3, 'Bo attacks'), row(5, 3, 'Bo rolls 4')];

  it('offers each step since the base once, on its first row', () => {
    expect([...rollbackTargets(history, 1)]).toEqual([2, 4]);
    // After a restart the base moved: older steps are no target.
    expect([...rollbackTargets(history, 2)]).toEqual([4]);
  });

  it('greys the undone rows and tells who rolled back', () => {
    const after = rolledBack(history, { seat: 1, to: 3, undid: 2 }, 4, 1, 0, n);
    expect(after.filter((e) => e.undone).map((e) => e.key)).toEqual([4, 5]);
    expect(after.at(-1)).toEqual({ key: 6, turn: 1, icon: '↶', text: 'Bo rolls back 2 steps', step: 4, rollback: true });
    // Undone rows and the rollback row are no target; the rest still are.
    expect([...rollbackTargets(after, 1)]).toEqual([2]);
    // The viewer's own rollback, of one step, past the first one.
    const again = rolledBack(after, { seat: 0, to: 1, undid: 1 }, 5, 1, 0, n);
    expect(again.at(-1)?.text).toBe('You roll back 1 step');
    expect(again.filter((e) => e.undone).map((e) => e.key)).toEqual([2, 3, 4, 5, 6]);
    expect(rollbackTargets(again, 1).size).toBe(0);
  });
});

describe('Soft rollback buttons (ADR 020, items 3 and 5)', () => {
  const row = (key: number, step: number, text: string): Entry => ({ key, turn: 1, icon: '+', text, step });
  const history = [row(1, 0, 'start'), row(2, 1, 'Ann attacks'), row(3, 2, 'Bo plays a card'), row(4, 3, 'Bo passes')];
  // The server lists steps 1 and 2: step 3 is a lone pass.
  const soft = {
    rollback_mode: 'soft' as const,
    rollback_from: 1,
    rollback_steps: [
      { step: 1, seat: 0 },
      { step: 2, seat: 1 },
    ],
    view: {},
  };

  it('gives the author an ask, the others a vote, and nothing else', () => {
    expect([...rollbackButtons(history, soft, 0)]).toEqual([
      [2, 'ask'],
      [3, 'vote'],
    ]);
    expect([...rollbackButtons(history, soft, 1)]).toEqual([
      [2, 'vote'],
      [3, 'ask'],
    ]);
    // A spectator has no buttons.
    expect(rollbackButtons(history, soft, -1).size).toBe(0);
  });

  it('offers nothing while a vote is open or once the game is over', () => {
    const vote = { step: 1, seat: 1, author: 0, undo: 3, voters: [1], votes: null, need: 1 };
    expect(rollbackButtons(history, { ...soft, rollback_vote: vote }, 1).size).toBe(0);
    expect(rollbackButtons(history, { ...soft, view: { over: true } }, 1).size).toBe(0);
    expect(rollbackButtons(history, { ...soft, rollback_mode: undefined }, 1).size).toBe(0);
  });

  it('keeps the barrier: steps before rollback_from have no button', () => {
    expect([...rollbackButtons(history, { ...soft, rollback_from: 2, rollback_steps: [{ step: 2, seat: 1 }] }, 0)]).toEqual([[3, 'vote']]);
  });

  it('gives every target a plain button in Godmode', () => {
    expect([...rollbackButtons(history, { rollback_mode: 'god', rollback_from: 1, view: {} }, 0).values()]).toEqual(['roll', 'roll', 'roll']);
  });

  it('adds a row when a vote ends with no rollback', () => {
    const after = voteEnded(history, { step: 1, seat: 1, why: 'declined' }, 1, 0, n);
    expect(after.at(-1)).toEqual({ key: 5, turn: 1, icon: '↶', text: 'The rollback was declined', rollback: true });
    expect(voteEnded(history, { step: 1, seat: 1, why: 'cancelled' }, 1, 0, n).at(-1)?.text).toBe('Bo cancelled the rollback vote');
    expect(voteEnded(history, { step: 1, seat: 0, why: 'cancelled' }, 1, 0, n).at(-1)?.text).toBe('You cancelled the rollback vote');
    expect(voteEnded(history, { step: 1, seat: 0, why: 'timeout' }, 1, 0, n).at(-1)?.text).toBe('The rollback vote ran out of time');
    // The row is no rollback target.
    expect(rollbackTargets(after, 1).has(5)).toBe(false);
  });

  it('adds one row where a server restart cut the history (ADR 020, item 10.6)', () => {
    const after = restartCut(history, 3, 1);
    expect(after.at(-1)).toEqual({ key: 5, turn: 1, icon: '⟳', text: 'The server restarted; no rollback before this', step: 3, rollback: true, restarted: true });
    // Every update after the restart says so; the row shows once.
    expect(restartCut(after, 3, 1)).toBe(after);
    // A later restart cuts again.
    expect(restartCut(after, 7, 1)).toHaveLength(after.length + 1);
    expect(rollbackTargets(after, 1).has(5)).toBe(false);
  });

  it('names the step a vote is about by its first row', () => {
    expect(stepText(history, 2)).toBe('Bo plays a card');
    expect(stepText(history, 9)).toBe('step 9');
  });
});

describe('a play has its own row (17.6)', () => {
  const penny: StackItem = {
    id: 3,
    kind: StackKind.loot,
    controller: 1,
    source: 40,
    card: 'a_penny',
    target: { player: 1, is_player: true },
    ability: { card: 'a_penny', index: 0 },
    roll_for: { card: '', index: 0 },
    text: 'Gain 1¢.',
  };
  // Step 5: Bo plays A Penny!; step 6 is Ann's pass; step 7 resolves it.
  const play = entriesFor([{ kind: 'loot_played', player: 1, object: 40, card: 'a_penny' }], 1, 0, n, 1).map((e) => ({ ...e, step: 5 }));
  const resolved = entriesFor(
    [
      { kind: 'stack_resolved', player: 1, card: 'a_penny', stack_id: 3, item: penny },
      { kind: 'gained_cents', player: 1, amount: 1 },
    ],
    1,
    0,
    n,
    1 + play.length,
  ).map((e) => ({ ...e, step: 7 }));
  const history = [...play, ...resolved];

  it('names who plays a loot card or uses an item; the strip keeps the result only', () => {
    expect(play).toEqual([{ key: 1, turn: 1, icon: '✦', text: 'Bo plays A_PENNY', card: 'a_penny', step: 5 }]);
    expect(entriesFor([{ kind: 'loot_played', player: 0, card: 'a_penny' }], 1, 0, n, 0)[0]?.text).toBe('You play A_PENNY');
    const [use] = entriesFor([{ kind: 'activated', player: 1, object: 41, card: 'sack_of_pennies', text: '↷: Gain 1¢.' }], 1, 0, n, 0);
    expect(use).toMatchObject({ icon: '✦', text: 'Bo uses SACK_OF_PENNIES', card: 'sack_of_pennies' });
    expect(use.headline).toBeUndefined();
    expect(stripEntries(history).map((e) => e.text)).toEqual(['A_PENNY']);
    // The play has no outcome: the result's gains go to the resolved line.
    expect(play[0]?.outcome).toBeUndefined();
  });

  it('gives the play row the button, in Soft while the card is on the stack and after', () => {
    const soft = { rollback_mode: 'soft' as const, rollback_from: 1, rollback_steps: [{ step: 5, seat: 1 }], view: {} };
    // On the stack: only the play row is there.
    expect([...rollbackButtons(play, soft, 0)]).toEqual([[1, 'vote']]);
    expect([...rollbackButtons(play, soft, 1)]).toEqual([[1, 'ask']]);
    // Resolved: the result came with a pass, no target; it has no button.
    expect([...rollbackButtons(history, soft, 0)]).toEqual([[1, 'vote']]);
  });

  it('gives the play row its own step in Godmode (19.69)', () => {
    const god = { rollback_mode: 'god' as const, rollback_from: 1, view: {} };
    const buttons = rollbackButtons(history, god, 0);
    expect(buttons.get(1)).toBe('roll');
    expect(history.find((e) => e.key === 1)?.step).toBe(5);
    // The result row undoes only the resolution: the card goes back on the stack.
    expect(buttons.get(2)).toBe('roll');
  });
});
