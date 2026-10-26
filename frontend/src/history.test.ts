import { describe, expect, it } from 'vitest';
import kindList from '../../pkg/rules_engine/testdata/event_kinds.txt?raw';
import { PromptKind } from './actions';
import {
  entriesFor as entries,
  historyKinds,
  lookGroups,
  lookedAt,
  moveOf,
  moveText,
  restartCut,
  rolledBack,
  stepName,
  stripEntries,
  undoButton,
  voteEnded,
  type Entry,
  type Move,
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
        { kind: 'dice_rolled', player: 0, amount: 2, text: 'first', label: { key: 'event.first_player_roll' } },
        { kind: 'dice_rolled', player: 1, amount: 2, text: 'first', label: { key: 'event.first_player_roll' } },
        { kind: 'first_roll_tied', player: -1, amount: 2, text: '0,1' },
        { kind: 'dice_rolled', player: 1, amount: 1, text: 'first', label: { key: 'event.first_player_roll' } },
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
    expect(entriesFor([{ kind: 'dice_rolled', player: 0, amount: 4, label: { key: 'event.attack_roll' } }], 1, 0, n, 0)).toEqual([]);
  });

  it('makes one entry per resolved stack item and marks turns', () => {
    const entries = entriesFor(
      [
        { kind: 'turn_started', player: 1, amount: 4 },
        { kind: 'passed_priority', player: 0 },
        { kind: 'stack_resolved', player: 1, object: 7, card: 'a_dime' },
        { kind: 'stack_resolved', player: 1, object: 8, amount: 5, label: { key: 'event.attack_roll' } },
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
        { kind: 'stack_resolved', player: 1, amount: 6, label: { key: 'event.attack_roll' }, stack_id: 3, item },
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
        { kind: 'stack_resolved', player: 1, amount: 6, label: { key: 'event.attack_roll' } },
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
        { kind: 'stack_resolved', player: 1, amount: 1, label: { key: 'event.attack_roll' }, stack_id: 3, item },
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
      { ...base, kind: StackKind.damage, amount: 2, label: { key: 'event.damage' }, target: gaper },
      { ...base, kind: StackKind.damage, controller: -1, amount: 1, label: { key: 'event.damage' }, source_card: 'gaper', target: { player: 1, is_player: true } },
      { ...base, kind: StackKind.death, controller: -1, label: { key: 'event.death' }, target: gaper },
      { ...base, kind: StackKind.penalty, controller: -1, label: { key: 'event.death_penalty' }, target: { player: 1, is_player: true } },
      { ...base, kind: StackKind.deathStep, controller: -1, label: { key: 'event.soul' }, target: gaper },
      { ...base, kind: StackKind.roll, roll: 5, attack: true, label: { key: 'event.attack_roll' }, target: gaper },
      { ...base, kind: StackKind.roll, roll: 3, label: { key: 'event.roll' }, target: { player: -1 } },
    ];
    for (const item of items) {
      const [entry] = entriesFor([{ kind: 'stack_resolved', player: item.controller, amount: item.roll, ...(item.label ? { label: item.label } : {}), stack_id: 7, item }], 1, 0, n, 0);
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

  it('keeps whose hand a card is from; a deck card has no owner (19.71)', () => {
    const events = [
      { kind: 'looked_at', player: 0, object: 7, card: 'a_dime', private: true, by: 2 },
      { kind: 'looked_at', player: 0, object: 8, card: 'a_penny', private: true, by: 3 },
      { kind: 'looked_at', player: 0, object: 9, card: 'a_nickel', private: true },
    ];
    const looked = lookedAt(events, 0);
    expect(looked).toEqual([
      { id: 7, card: 'a_dime', owner: 1 },
      { id: 8, card: 'a_penny', owner: 2 },
      { id: 9, card: 'a_nickel' },
    ]);
    expect(looked[2]).not.toHaveProperty('owner');
  });

  it('groups hand cards by owner in event order; deck cards in one group (19.71)', () => {
    const groups = lookGroups([
      { id: 1, card: 'a', owner: 2 },
      { id: 2, card: 'b', owner: 1 },
      { id: 3, card: 'c', owner: 2 },
      { id: 4, card: 'd' },
    ]);
    expect(groups.map((g) => [g.owner, g.cards.map((c) => c.id)])).toEqual([
      [2, [1, 3]],
      [1, [2]],
      [undefined, [4]],
    ]);
  });
});

describe('historyKinds', () => {
  const kinds = kindList.split('\n').filter(Boolean);
  // sample: an event of a told kind that gives a line; others need only
  // their kind and a player.
  const sample = new Map<string, Partial<GameEvent>>([['dice_rolled', { text: 'first', label: { key: 'event.first_player_roll' }, amount: 3 }]]);

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

  it('greys the undone rows and tells who rolled back', () => {
    const after = rolledBack(history, { seat: 1, to: 3, undid: 2 }, 4, 1, 0, n);
    expect(after.filter((e) => e.undone).map((e) => e.key)).toEqual([4, 5]);
    expect(after.at(-1)).toEqual({ key: 6, turn: 1, icon: '↶', text: 'Bo rolls back 2 steps', step: 4, rollback: true });
    // The viewer's own rollback, of one step, past the first one.
    const again = rolledBack(after, { seat: 0, to: 1, undid: 1 }, 5, 1, 0, n);
    expect(again.at(-1)?.text).toBe('You roll back 1 step');
    expect(again.filter((e) => e.undone).map((e) => e.key)).toEqual([2, 3, 4, 5, 6]);
  });

  it('adds a row when a vote ends with no rollback', () => {
    const after = voteEnded(history, { step: 1, seat: 1, why: 'declined' }, 1, 0, n);
    expect(after.at(-1)).toEqual({ key: 6, turn: 1, icon: '↶', text: 'The rollback was declined', rollback: true });
    expect(voteEnded(history, { step: 1, seat: 1, why: 'cancelled' }, 1, 0, n).at(-1)?.text).toBe('Bo cancelled the rollback vote');
    expect(voteEnded(history, { step: 1, seat: 0, why: 'cancelled' }, 1, 0, n).at(-1)?.text).toBe('You cancelled the rollback vote');
    expect(voteEnded(history, { step: 1, seat: 0, why: 'timeout' }, 1, 0, n).at(-1)?.text).toBe('The rollback vote ran out of time');
  });

  it('adds one row where a server restart cut the history (ADR 020, item 10.6)', () => {
    const after = restartCut(history, 3, 1);
    expect(after.at(-1)).toEqual({ key: 6, turn: 1, icon: '⟳', text: 'The server restarted; no rollback before this', step: 3, rollback: true, restarted: true });
    // Every update after the restart says so; the row shows once.
    expect(restartCut(after, 3, 1)).toBe(after);
    // A later restart cuts again.
    expect(restartCut(after, 7, 1)).toHaveLength(after.length + 1);
  });
});

describe('moves (17.7)', () => {
  const loot = (player: number, card: string): GameEvent => ({ kind: 'loot_played', player, object: 40, card });

  it('names a play by the card, though the history has no line for it (15.88)', () => {
    const events = [loot(1, 'a_penny')];
    expect(entriesFor(events, 1, 0, n, 0)).toEqual([]);
    expect(moveOf(events, [], 5, 5, 0, n)).toEqual({ from: 5, to: 5, text: 'Bo plays A_PENNY' });
    expect(moveOf([loot(0, 'bomb')], [], 5, 5, 0, n).text).toBe('You play BOMB');
    const use = { kind: 'activated', player: 1, object: 41, card: 'sack_of_pennies' };
    expect(moveOf([use], [], 5, 5, 0, n).text).toBe('Bo uses SACK_OF_PENNIES');
  });

  it('names an ended turn, else the first line that is no turn line', () => {
    expect(moveOf([{ kind: 'turn_ended', player: 1 }], [], 2, 2, 0, n).text).toBe('Bo ends the turn');
    const lines: Entry[] = [
      { key: 1, turn: 2, icon: '▶', text: 'Turn 2: Bo' },
      { key: 2, turn: 2, icon: '⚔', text: 'Bo declares an attack' },
    ];
    expect(moveOf([], lines, 3, 3, 0, n).text).toBe('Bo declares an attack');
    expect(moveOf([], [], 3, 3, 0, n)).toEqual({ from: 3, to: 3 });
  });

  it('names the barrier it holds', () => {
    expect(moveOf([{ kind: 'dice_rolled', player: 1, amount: 4 }], [], 3, 4, 0, n).barrier).toBe('a die was rolled');
    expect(moveOf([{ kind: 'turn_ended', player: 1 }, { kind: 'dice_rolled', player: 1 }], [], 3, 3, 0, n).barrier).toBe('the turn ended');
    expect(moveOf([{ kind: 'gained_cents', player: 1, private: true }], [], 3, 3, 0, n).barrier).toBe('a hidden card was seen');
    // A hand card put on the loot deck is no barrier (ADR 020, item 4.22).
    expect(moveOf([{ kind: 'moved_to_deck', player: 1, private: true }], [], 3, 3, 0, n).barrier).toBeUndefined();
  });

  it("names a target by its update, or its play's answer up to its last step", () => {
    const moves: Move[] = [
      { from: 4, to: 4, text: 'Bo attacks' },
      { from: 5, to: 5 },
      { from: 6, to: 7, text: 'Bo plays BOMB' },
    ];
    expect(moveText(moves, { step: 4, seat: 1 }, 0, n.player)).toBe('Bo attacks');
    // The play asked its target; the answer (step 6) names it (17.7).
    expect(moveText(moves, { step: 5, seat: 1, last: 6 }, 0, n.player)).toBe('Bo plays BOMB');
    expect(moveText(moves, { step: 7, seat: 1 }, 0, n.player)).toBe('Bo plays BOMB');
  });

  it("names a move without a line of its own Bo's choice, never the next move (17.8)", () => {
    // Step 5 put a hand card back on the deck: no line. Step 6 is a
    // move of its own, after it.
    const moves: Move[] = [
      { from: 5, to: 5 },
      { from: 6, to: 6, text: 'Ann plays BOMB' },
    ];
    expect(moveText(moves, { step: 5, seat: 1 }, 0, n.player)).toBe("Bo's choice");
    expect(moveText(moves, { step: 5, seat: 0 }, 0, n.player)).toBe('your choice');
    // A play still waiting for its answer has no line yet either.
    expect(moveText(moves.slice(0, 1), { step: 5, seat: 1, last: 5 }, 0, n.player)).toBe("Bo's choice");
    expect(moveText(moves, { step: 9, seat: 2 }, 0, n.player)).toBe("Cy's choice");
  });
});

describe('the Undo button (17.7)', () => {
  const moves: Move[] = [
    { from: 1, to: 1, text: 'Ann attacks', barrier: 'a die was rolled' },
    { from: 2, to: 2, text: 'Bo plays BOMB' },
    { from: 3, to: 3, text: 'Ann passes' },
  ];
  // The server lists steps 1 and 2; step 3 is a lone pass.
  const soft = {
    rollback_mode: 'soft' as const,
    rollback_from: 1,
    rollback_steps: [
      { step: 1, seat: 0 },
      { step: 2, seat: 1 },
    ],
    view: {},
  };

  it('is hidden in Disabled mode', () => {
    expect(undoButton({ ...soft, rollback_mode: undefined }, 0, moves, n.player)).toBeNull();
  });

  it('targets the latest move and names it', () => {
    expect(undoButton(soft, 0, moves, n.player)).toEqual({ label: 'Undo', step: 2, ask: false, title: 'Undo: Bo plays BOMB' });
    // Godmode: anyone, the author too, with a plain Undo.
    expect(undoButton({ ...soft, rollback_mode: 'god' }, 1, moves, n.player)).toEqual({ label: 'Undo', step: 2, ask: false, title: 'Undo: Bo plays BOMB' });
  });

  it("asks for the Soft target's own author", () => {
    expect(undoButton(soft, 1, moves, n.player)).toEqual({ label: 'Ask to undo', step: 2, ask: true, title: 'Ask to undo: Bo plays BOMB' });
  });

  it("names a target without a line Bo's choice, yours your choice (17.8)", () => {
    const pick = [...moves, { from: 4, to: 4 }];
    const u = { ...soft, rollback_steps: [...soft.rollback_steps, { step: 4, seat: 1 }] };
    expect(undoButton(u, 0, pick, n.player)?.title).toBe("Undo: Bo's choice");
    expect(undoButton(u, 1, pick, n.player)?.title).toBe('Ask to undo: your choice');
    expect(undoButton({ ...u, rollback_mode: 'god' }, 1, pick, n.player)?.title).toBe('Undo: your choice');
  });

  it('is greyed with the reason', () => {
    const why = (u: Parameters<typeof undoButton>[0], you = 0) => undoButton(u, you, moves, n.player)?.why;
    const vote = { step: 1, seat: 1, author: 0, undo: 3, voters: [1], votes: null, need: 1 };
    expect(why(soft, -1)).toBe('Spectators cannot undo');
    expect(why({ ...soft, view: { over: true } })).toBe('The game is over');
    expect(why({ ...soft, rollback_vote: vote })).toBe('A rollback vote is open');
    expect(why({ ...soft, pause: { away: [1], votes: null } })).toBe('The game is paused');
    // Nothing after the barrier: the barrier's step is named.
    expect(why({ ...soft, rollback_from: 2, rollback_steps: [] })).toBe('Nothing to undo since a die was rolled');
    expect(why({ ...soft, rollback_from: 3, rollback_steps: [] })).toBe('Nothing to undo yet');
    expect(why({ ...soft, rollback_from: 2, rollback_steps: [], rollback_cut: 1 })).toBe('Nothing to undo since the server restarted');
    // Godmode has no barriers: nothing is left since the base.
    expect(why({ ...soft, rollback_mode: 'god', rollback_from: 2, rollback_steps: [] })).toBe('Nothing to undo yet');
    // A greyed button names no step and says why in its tooltip.
    expect(undoButton({ ...soft, view: { over: true } }, 1, moves, n.player)).toEqual({ label: 'Ask to undo', ask: true, why: 'The game is over', title: 'The game is over' });
  });
});
