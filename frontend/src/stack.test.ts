import { describe, expect, it } from 'vitest';
import { PromptKind } from './actions';
import type { GameView, StackItem } from './protocol';
import {
  declaration,
  findCard,
  playedCard,
  secondsLeft,
  stackBadge,
  stackCard,
  stackIcon,
  StackKind,
  stackLine,
  stackParts,
  stackRows,
  stackWait,
  waitTurn,
} from './stack';

const view: GameView = {
  players: [
    { id: 0, character: { id: 1, card: 'isaac', charged: true }, hp: 2, atk: 1, cents: 3, souls: 0, in_play: null, hand_size: 3, loot_plays: 0, attacks: 0, purchases: 0 },
    { id: 1, character: { id: 2, card: 'cain', charged: true }, hp: 2, atk: 1, cents: 3, souls: 0, in_play: null, hand_size: 3, loot_plays: 0, attacks: 0, purchases: 0 },
  ],
  shop: null,
  monsters: [{ top: { id: 9, card: 'gaper', charged: true } }],
  deck_sizes: [1, 1, 1, 0],
  discards: [null, null, null, null],
  stack: null,
  turn: { active: 0, number: 3, step: 4, loot_plays: 0, attacks: 0, purchases: 0 },
  waiting: { kind: PromptKind.priority, player: 0 },
};

const names = {
  player: (p: number) => ['Ann', 'Bo'].at(p) ?? '?',
  card: (ref: string) => ref.toUpperCase(),
};

function item(partial: Partial<StackItem>): StackItem {
  return {
    id: 1,
    kind: StackKind.loot,
    controller: 0,
    source: 5,
    target: { player: 0 },
    ability: { card: '', index: 0 },
    roll_for: { card: '', index: 0 },
    ...partial,
  };
}

describe('stackLine', () => {
  it('names loot, rolls, damage and deaths', () => {
    expect(stackLine(item({ card: 'a_dime' }), view, names)).toBe('Ann plays A_DIME');
    expect(stackLine(item({ kind: StackKind.roll, roll: 4, attack: true, controller: 1 }), view, names)).toBe(
      'Attack: Bo rolls 4',
    );
    expect(stackLine(item({ kind: StackKind.damage, amount: 2, target: { player: -1, object: 9 } }), view, names)).toBe(
      'Ann deals 2 damage to GAPER',
    );
    expect(stackLine(item({ kind: StackKind.death, target: { player: 1, is_player: true } }), view, names)).toBe(
      'Bo dies',
    );
  });

  it('names what changes a thrown roll, and its result (15.94)', () => {
    const roll = item({
      kind: StackKind.roll,
      roll: 2,
      attack: true,
      target: { player: -1, object: 9 },
      sources: [{ stat: 'attack_roll', source: 5, card: 'meat', amount: 1, lasts: 'in_play' }],
      final: 3,
    });
    expect(stackLine(roll, view, names)).toBe('Attack: Ann rolls 2 against GAPER (+1 MEAT = 3)');
    // Changes that cancel out: the result is the die.
    const even = item({
      kind: StackKind.roll,
      roll: 4,
      sources: [
        { stat: 'roll', amount: 1, lasts: 'turn' },
        { stat: 'roll', amount: -1, lasts: 'turn' },
      ],
    });
    expect(stackLine(even, view, names)).toBe('Ann rolls 4 (+1 a hidden card (this turn), −1 a hidden card (this turn) = 4)');
  });

  it('names whom an attack is against and who hits back (15.22)', () => {
    const roll = item({ kind: StackKind.roll, roll: 3, attack: true, target: { player: -1, object: 9 } });
    expect(stackLine(roll, view, names)).toBe('Attack: Ann rolls 3 against GAPER');
    const rune = item({ kind: StackKind.roll, roll: 1, label: 'roll', roll_for: { card: 'blank_rune', index: 0 } });
    expect(stackLine(rune, view, names)).toBe('Ann rolls 1 for BLANK_RUNE');
    // A monster that died since: the item still names its card.
    const gone = item({ kind: StackKind.damage, amount: 1, attack: true, target: { player: -1, object: 77 }, target_card: 'fly' });
    expect(stackLine(gone, view, names)).toBe('Attack: Ann deals 1 damage to FLY');
    const back = item({ kind: StackKind.damage, amount: 1, attack: true, controller: -1, source: 9, target: { player: 0, is_player: true } });
    expect(stackLine(back, view, names)).toBe('Attack: GAPER deals 1 damage to Ann');
  });

  it('names the card when a player takes damage (15.47)', () => {
    // Boom Fly: each player controls the damage they take; the fly is gone.
    const fly = item({ kind: StackKind.damage, amount: 1, controller: 1, source: 307, source_card: 'boom_fly', target: { player: 1, is_player: true } });
    expect(stackLine(fly, view, names)).toBe('BOOM_FLY deals 1 damage to Bo');
    const bare = item({ kind: StackKind.damage, amount: 3, controller: 0, source: 0, target: { player: 0, is_player: true } });
    expect(stackLine(bare, view, names)).toBe('Ann takes 3 damage');
    // Damage to someone else still names the player who deals it.
    const other = item({ kind: StackKind.damage, amount: 1, controller: 0, source_card: 'mr_boom', target: { player: 1, is_player: true } });
    expect(stackLine(other, view, names)).toBe('Ann deals 1 damage to Bo');
  });
});

describe('stackParts', () => {
  it('keeps the player and the target apart from the words (S-03)', () => {
    expect(stackParts(item({ card: 'a_dime' }), view, names)).toEqual([
      { text: 'Ann', role: 'player', player: 0 },
      { text: ' plays ' },
      { text: 'A_DIME', role: 'card' },
    ]);
    expect(stackParts(item({ kind: StackKind.damage, amount: 1, target: { player: 1, is_player: true } }), view, names)).toEqual([
      { text: 'Ann', role: 'player', player: 0 },
      { text: ' deals ' },
      { text: '1 damage to ' },
      { text: 'Bo', role: 'target', player: 1 },
    ]);
  });
});

describe('stackIcon and stackCard', () => {
  it('picks an icon by kind (S-04)', () => {
    expect(stackIcon(item({ card: 'a_dime' }))).toBe('card');
    expect(stackIcon(item({ kind: StackKind.roll }))).toBe('roll');
    expect(stackIcon(item({ kind: StackKind.damage }))).toBe('damage');
    expect(stackIcon(item({ kind: StackKind.penalty }))).toBe('death');
  });

  it('shows the item card, or the card it comes from (S-06)', () => {
    expect(stackCard(item({ card: 'a_dime' }), view)).toBe('a_dime');
    expect(stackCard(item({ kind: StackKind.roll, source: 9 }), view)).toBe('gaper');
    expect(stackCard(item({ kind: StackKind.roll, source: 99 }), view)).toBeUndefined();
    // An attack roll has no source: its monster shows (15.22).
    expect(stackCard(item({ kind: StackKind.roll, source: 0, attack: true, target: { player: -1, object: 9 } }), view)).toBe('gaper');
    expect(stackCard(item({ kind: StackKind.roll, source: 0, target_card: 'fly', target: { player: -1, object: 77 } }), view)).toBe('fly');
  });
});

describe('stackBadge', () => {
  it('marks attacks with a sword and loot plays with a bolt (H-08)', () => {
    expect(stackBadge(item({ kind: StackKind.roll, attack: true }))).toBe('attack');
    expect(stackBadge(item({ kind: StackKind.damage, attack: true }))).toBe('attack');
    expect(stackBadge(item({ card: 'a_dime' }))).toBe('spell');
    expect(stackBadge(item({ kind: StackKind.ability }))).toBeUndefined();
    expect(stackBadge(item({ kind: StackKind.roll }))).toBeUndefined();
  });
});

describe('playedCard', () => {
  it("is the newest loot card another player has on the stack (15.31)", () => {
    const bomb = item({ id: 1, controller: 1, card: 'bomb' });
    const penny = item({ id: 2, controller: 0, card: 'a_penny' });
    const roll = item({ id: 3, controller: 1, kind: StackKind.roll });
    expect(playedCard({ ...view, stack: [bomb, penny, roll] }, 0)).toBe(bomb);
    expect(playedCard({ ...view, stack: [bomb, penny, roll] }, 1)).toBe(penny);
    expect(playedCard({ ...view, stack: [penny] }, 0)).toBeUndefined();
    expect(playedCard(view, 0)).toBeUndefined();
  });
});

describe('findCard', () => {
  it('finds characters and monsters by object', () => {
    expect(findCard(view, 2)).toBe('cain');
    expect(findCard(view, 9)).toBe('gaper');
    expect(findCard(view, 99)).toBeUndefined();
  });
});

describe('waitTurn', () => {
  it('is offered on another turn until the turn before ends (N-16)', () => {
    expect(waitTurn(view, 0, false)).toBe('none');
    expect(waitTurn(view, 1, false)).toBe('offer');
    expect(waitTurn(view, 1, true)).toBe('on');
    const ending = { ...view, turn: { ...view.turn, end_declared: true } };
    expect(waitTurn(ending, 1, false)).toBe('none');
    expect(waitTurn({ ...view, turn: { ...view.turn, step: 5 } }, 1, false)).toBe('none');
    expect(waitTurn({ ...view, over: true }, 1, false)).toBe('none');
  });
});

describe('stackWait', () => {
  it('names who may respond and who passed, as the server says', () => {
    const stacked = { ...view, stack: [item({})], passed: [1], responding: [0] };
    expect(stackWait(stacked, 0, names)).toEqual({ responding: ['You'], passed: ['Bo'] });
    expect(stackWait(stacked, 1, names)).toEqual({ responding: ['Ann'], passed: ['You'] });
  });
});

describe('secondsLeft', () => {
  it('rounds up and stops at zero', () => {
    expect(secondsLeft(10_500, 0)).toBe(11);
    expect(secondsLeft(0, 5000)).toBe(0);
  });
});

describe('declaration', () => {
  const line = (v: GameView) =>
    declaration(v, names)
      ?.map((p) => p.text)
      .join('');

  it('says an open attack, then its target once aimed (15.91)', () => {
    expect(line({ ...view, attack: {} })).toBe('Ann attacks…');
    expect(line({ ...view, attack: { target: 9 } })).toBe('Ann attacks GAPER');
    expect(line({ ...view, attack: { deck: true } })).toBe('Ann attacks the monster deck');
  });

  it('names an attacked player, not their character (R-ATK-20)', () => {
    const parts = declaration({ ...view, attack: { target: 2 } }, names);
    expect(parts?.at(-1)).toEqual({ text: 'Bo', role: 'target', player: 1 });
  });

  it('says an open buy', () => {
    expect(line({ ...view, buying: true })).toBe('Ann buys…');
    expect(line({ ...view, turn: { ...view.turn, active: 1 }, buying: true })).toBe('Bo buys…');
  });

  it('says a declared end of turn until the end phase; not a death end (19.74)', () => {
    const ending = { ...view, turn: { ...view.turn, end_declared: true } };
    expect(line(ending)).toBe('Ann ends the turn');
    expect(declaration(ending, names)?.[0]).toEqual({ text: 'Ann', role: 'player', player: 0 });
    expect(line({ ...view, turn: { ...view.turn, active: 1, end_declared: true } })).toBe('Bo ends the turn');
    expect(declaration({ ...ending, turn: { ...ending.turn, step: 5 } }, names)).toBeNull();
    expect(declaration({ ...view, turn: { ...view.turn, death_end: true } }, names)).toBeNull();
    expect(stackRows(ending, names).map((r) => r.key)).toEqual(['declared']);
  });

  it('is nothing without one, after the game, and before turn 1', () => {
    expect(declaration(view, names)).toBeNull();
    expect(declaration({ ...view, buying: true, over: true }, names)).toBeNull();
    expect(declaration({ ...view, attack: {}, first_roll: { round: 1, players: [0, 1] } }, names)).toBeNull();
  });
});

describe('stackRows', () => {
  const keys = (v: GameView) => stackRows(v, names).map((r) => r.key);

  it('keeps every row key while the stack grows and shrinks (15.51)', () => {
    const a = item({ id: 4, card: 'a_dime' });
    const b = item({ id: 7, kind: StackKind.roll, roll: 3 });
    const c = item({ id: 9, kind: StackKind.damage, amount: 1 });
    expect(keys({ ...view, stack: [a] })).toEqual(['item-4']);
    expect(keys({ ...view, stack: [a, b, c] })).toEqual(['item-9', 'item-7', 'item-4']);
    // The top resolves: the others keep their keys, so they move.
    expect(keys({ ...view, stack: [a, b] })).toEqual(['item-7', 'item-4']);
    // An item in the middle goes too, e.g. a cancelled roll.
    expect(keys({ ...view, stack: [a, c] })).toEqual(['item-9', 'item-4']);
  });

  it('puts an open declaration on top, under its own key', () => {
    const a = item({ id: 4, kind: StackKind.roll, roll: 2, attack: true });
    const rows = stackRows({ ...view, attack: { target: 9 }, stack: [a] }, names);
    expect(rows.map((r) => r.key)).toEqual(['declared', 'item-4']);
    expect(keys({ ...view, attack: { target: 9 } })).toEqual(['declared']);
  });
});
