import { describe, expect, it } from 'vitest';
import { PromptKind } from './actions';
import type { GameView, StackItem } from './protocol';
import { findCard, secondsLeft, StackKind, stackLine, stackWait, waitText } from './stack';

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
      '2 damage to GAPER',
    );
    expect(stackLine(item({ kind: StackKind.death, target: { player: 1, is_player: true } }), view, names)).toBe(
      'Bo dies',
    );
  });
});

describe('findCard', () => {
  it('finds characters and monsters by object', () => {
    expect(findCard(view, 2)).toBe('cain');
    expect(findCard(view, 9)).toBe('gaper');
    expect(findCard(view, 99)).toBeUndefined();
  });
});

describe('waitText', () => {
  it('says whose turn it is, or who may respond', () => {
    expect(waitText(view, 0, names)).toBe('Your turn');
    expect(waitText(view, 1, names)).toBe("Ann's turn");
    const stacked = { ...view, stack: [item({})], waiting: { kind: PromptKind.priority, player: 1 } };
    expect(waitText(stacked, 0, names)).toBe('Bo may respond');
    expect(waitText(stacked, 1, names)).toBe('Respond, or pass');
  });

  it('says who is choosing or discarding', () => {
    const choosing = { ...view, waiting: { kind: PromptKind.choose, player: 1 } };
    expect(waitText(choosing, 0, names)).toBe('Bo is choosing');
    expect(waitText(choosing, 1, names, 'Choose a target')).toBe('Choose a target');
    expect(waitText({ ...view, over: true }, 0, names)).toBe('Game over');
  });
});

describe('stackWait', () => {
  it('names who may respond and who passed, as the server says', () => {
    const stacked = { ...view, stack: [item({})], passed: [1], responding: [0] };
    expect(stackWait(stacked, 0, names)).toEqual({ responding: ['You'], passed: ['Bo'] });
    expect(stackWait(stacked, 1, names)).toEqual({ responding: ['Ann'], passed: ['You'] });
  });

  it('names who must answer a prompt, and what', () => {
    const choosing = { ...view, waiting: { kind: PromptKind.choose, player: 1, purpose: 6 } };
    expect(stackWait(choosing, 0, names).prompt).toBe('Bo: Choose a target');
    const asked = { ...view, waiting: { kind: PromptKind.choose, player: 0, text: 'Pay 1¢?' } };
    expect(stackWait(asked, 0, names).prompt).toBe('You: Pay 1¢?');
    const discarding = { ...view, waiting: { kind: PromptKind.discard, player: 1, count: 2 } };
    expect(stackWait(discarding, 0, names).prompt).toBe('Bo: Discard 2 loot cards');
  });
});

describe('secondsLeft', () => {
  it('rounds up and stops at zero', () => {
    expect(secondsLeft(10_500, 0)).toBe(11);
    expect(secondsLeft(0, 5000)).toBe(0);
  });
});
