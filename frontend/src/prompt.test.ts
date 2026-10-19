import { describe, expect, it } from 'vitest';
import { IntentKind, PromptKind, type Actions } from './actions';
import type { GameView, Intent, StackItem } from './protocol';
import { picking, promptBanner, stackPrompt } from './prompt';
import { StackKind } from './stack';

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

const none: Actions = { buttons: [], turn: {}, cards: new Map(), blocked: new Map(), choice: null, discard: null, order: null, throw: null };
const intent = (kind: number): Intent => ({ player: 0, kind });

const loot: StackItem = {
  id: 1,
  kind: StackKind.loot,
  controller: 1,
  source: 5,
  card: 'a_penny',
  target: { player: 0 },
  ability: { card: '', index: 0 },
  roll_for: { card: '', index: 0 },
};

const banner = (v: GameView, you: number, actions = none) => promptBanner(v, you, names, actions);

describe('promptBanner', () => {
  it('says what the player may do on their own turn', () => {
    const all = { ...none, cards: new Map([[7, []]]), turn: { attack: intent(IntentKind.attack), buy: intent(IntentKind.purchase), endTurn: intent(IntentKind.endTurn) } };
    expect(banner(view, 0, all)).toEqual({ kind: 'turn', mine: true, player: 0, title: 'Your turn', hint: 'Play a card, attack, buy, or end the turn' });
    const attack = { ...none, turn: { attack: intent(IntentKind.attack) } };
    expect(banner(view, 0, attack)?.hint).toBe('Attack');
    const idle = { ...none, turn: { endTurn: intent(IntentKind.endTurn), idle: true } };
    expect(banner(view, 0, idle)?.hint).toBe('Nothing left to do: end the turn');
    expect(banner(view, 0)?.hint).toBe('Pass to go on');
  });

  it('names whom the game waits for, and what for', () => {
    expect(banner(view, 1)).toEqual({ kind: 'turn', mine: false, player: 0, title: 'Waiting for Ann', hint: 'Playing their turn' });
    const out = { ...view, waiting: { kind: PromptKind.priority, player: 1 } };
    expect(banner(out, 1)).toMatchObject({ kind: 'act', title: 'You may act', hint: 'Use an item or a loot card, or pass' });
    expect(banner(out, 0)).toMatchObject({ kind: 'act', title: 'Waiting for Bo', hint: 'May act, or pass' });
  });

  it('asks to respond to the top of the stack', () => {
    const stacked = { ...view, stack: [loot] };
    expect(banner(stacked, 0)).toMatchObject({ kind: 'respond', mine: true, title: 'Respond, or pass', hint: 'On the stack: Bo plays A_PENNY' });
    expect(banner(stacked, 1)).toMatchObject({ kind: 'respond', mine: false, title: 'Waiting for Ann', hint: 'Respond, or pass' });
  });

  it('covers choices, discards, orders and throws', () => {
    const at = (waiting: GameView['waiting']) => ({ ...view, waiting });
    expect(banner(at({ kind: PromptKind.choose, player: 0, purpose: 6 }), 0)).toMatchObject({ kind: 'choose', title: 'Choose a target' });
    expect(banner(at({ kind: PromptKind.choose, player: 0, text: 'Pay 1¢?' }), 1)).toMatchObject({ title: 'Waiting for Ann', hint: 'Pay 1¢?' });
    expect(banner(at({ kind: PromptKind.discard, player: 0, count: 1 }), 0)).toMatchObject({ kind: 'discard', title: 'Discard 1 loot card' });
    expect(banner(at({ kind: PromptKind.discard, player: 0, count: 2 }), 1)?.hint).toBe('Discard 2 loot cards');
    expect(banner(at({ kind: PromptKind.order, player: 0, purpose: 12 }), 0)).toMatchObject({ kind: 'order', title: 'Choose the card for the very bottom' });
    const throwing = at({ kind: PromptKind.throw, player: 0, text: 'attack roll' });
    expect(banner(throwing, 0)).toMatchObject({ kind: 'throw', title: 'Throw the die', hint: 'For the attack roll' });
    expect(banner(throwing, 1)?.hint).toBe('Throw the die for the attack roll');
  });

  it("says the player's own aim, and how to cancel it", () => {
    const aiming = { source: 9, targets: [] };
    expect(promptBanner(view, 0, names, none, aiming)).toMatchObject({ kind: 'aim', mine: true, title: 'Choose a target', hint: 'For GAPER. Esc cancels' });
  });

  it('shows nothing once the game is over', () => {
    expect(banner({ ...view, over: true }, 0)).toBeNull();
  });
});

describe('stackPrompt', () => {
  it('repeats prompts other than priority in the stack panel (V-11)', () => {
    const choosing = { ...view, waiting: { kind: PromptKind.choose, player: 1, purpose: 6 } };
    expect(stackPrompt(banner(choosing, 0), names)).toBe('Bo: Choose a target');
    expect(stackPrompt(banner(choosing, 1), names)).toBe('Choose a target');
    expect(stackPrompt(banner(view, 0), names)).toBeUndefined();
    expect(stackPrompt(null, names)).toBeUndefined();
  });
});

describe('picking', () => {
  it('is on while aiming or while a choice offers a table target', () => {
    expect(picking(none, null)).toBe(false);
    expect(picking(none, { source: 1, targets: [] })).toBe(true);
    const option = { label: 'gaper', intent: intent(IntentKind.choose) };
    const dialog = { ...none, choice: { title: 'Choose one', purpose: 8, options: [option] } };
    expect(picking(dialog, null)).toBe(false);
    const table = { ...none, choice: { title: 'Choose a target', purpose: 6, options: [{ ...option, target: { object: 9, player: 0 } }] } };
    expect(picking(table, null)).toBe(true);
  });
});
