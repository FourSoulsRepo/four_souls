import { describe, expect, it } from 'vitest';
import {
  ActionKind,
  actionsFor,
  aimedOption,
  IntentKind,
  orderIntent,
  outcomeLines,
  PromptKind,
  Purpose,
  sameTarget,
} from './actions';
import type { Update } from './protocol';

function update(partial: Partial<Update['view']>, allowed: Update['allowed']): Update {
  return {
    step: 1,
    events: [],
    allowed,
    seats: null,
    view: {
      players: [],
      shop: null,
      monsters: null,
      deck_sizes: [0, 0, 0, 0],
      discards: [null, null, null, null],
      stack: null,
      turn: { active: 0, number: 1, step: 4, loot_plays: 0, attacks: 0, purchases: 0 },
      waiting: { kind: PromptKind.priority, player: 0 },
      ...partial,
    },
  };
}

const isCard = (label: string) => !label.includes(' ');

describe('actionsFor', () => {
  it('turns allowed intents into buttons and card actions', () => {
    const a = actionsFor(
      update({}, [
        { player: 0, kind: IntentKind.endTurn },
        { player: 0, kind: IntentKind.attack },
        { player: 0, kind: IntentKind.purchase },
        { player: 0, kind: IntentKind.playLoot, objects: [7] },
        { player: 0, kind: IntentKind.activate, objects: [3] },
        { player: 0, kind: IntentKind.activate, objects: [4] },
        { player: 0, kind: IntentKind.activate, objects: [4], choice: 1 },
      ]),
      isCard,
    );
    expect(a.buttons).toEqual([]);
    expect(a.turn.endTurn?.kind).toBe(IntentKind.endTurn);
    expect(a.turn.idle).toBeUndefined();
    expect(a.turn.attack?.kind).toBe(IntentKind.attack);
    expect(a.turn.buy?.kind).toBe(IntentKind.purchase);
    expect(a.cards.get(7)?.map((x) => x.label)).toEqual(['Play']);
    expect(a.cards.get(3)?.map((x) => x.label)).toEqual(['Use']);
    expect(a.cards.get(4)?.map((x) => x.label)).toEqual(['Use ability 1', 'Use ability 2']);
    expect(a.choice).toBeNull();
  });

  it('gives each ability of a card with many its line (C-05)', () => {
    const bone = { id: 4, card: 'the_bone', charged: true, abilities: ['↷: Put a counter on this.', '', 'Remove 2: Deal 1.'] };
    const player = {
      id: 0,
      character: { id: 1, card: 'isaac', charged: true },
      hp: 2,
      atk: 1,
      cents: 0,
      souls: 0,
      in_play: [bone],
      hand_size: 0,
      loot_plays: 0,
      attacks: 0,
      purchases: 0,
    };
    const a = actionsFor(
      update({ players: [player] }, [
        { player: 0, kind: IntentKind.activate, objects: [4] },
        { player: 0, kind: IntentKind.activate, objects: [4], choice: 2 },
      ]),
      isCard,
    );
    expect(a.cards.get(4)?.map((x) => x.line)).toEqual(['↷: Put a counter on this.', 'Remove 2: Deal 1.']);
  });

  it('is idle when only passing or ending the turn is left', () => {
    const a = actionsFor(
      update({}, [
        { player: 0, kind: IntentKind.pass },
        { player: 0, kind: IntentKind.endTurn },
      ]),
      isCard,
    );
    expect(a.buttons.map((b) => b.label)).toEqual(['Pass']);
    expect(a.turn.idle).toBe(true);
    expect(a.turn.left).toBeUndefined();
  });

  it('lists what is left before ending the turn (V-46)', () => {
    const all = actionsFor(
      update({}, [
        { player: 0, kind: IntentKind.activate, objects: [3] },
        { player: 0, kind: IntentKind.playLoot, objects: [7] },
        { player: 0, kind: IntentKind.playLoot, objects: [8] },
        { player: 0, kind: IntentKind.purchase },
        { player: 0, kind: IntentKind.attack },
        { player: 0, kind: IntentKind.endTurn },
      ]),
      isCard,
    );
    expect(all.turn.left).toEqual(['Attack', 'Buy', 'Play loot', 'Use abilities']);
    const loot = actionsFor(
      update({}, [
        { player: 0, kind: IntentKind.playLoot, objects: [7] },
        { player: 0, kind: IntentKind.endTurn },
      ]),
      isCard,
    );
    expect(loot.turn.left).toEqual(['Play loot']);
    const notMine = actionsFor(update({}, [{ player: 0, kind: IntentKind.playLoot, objects: [7] }]), isCard);
    expect(notMine.turn.left).toBeUndefined();
  });

  it('offers nothing when it is not the viewer’s move', () => {
    const a = actionsFor(update({}, null), isCard);
    expect(a.buttons).toEqual([]);
    expect(a.turn).toEqual({});
    expect(a.cards.size).toBe(0);
  });

  it('builds a choice with cards and table targets', () => {
    const a = actionsFor(
      update(
        {
          waiting: { kind: PromptKind.choose, player: 0, purpose: 1, options: ['gaper', 'monster deck'] },
          options: [
            { object: 42, player: -1 },
            { player: -1, deck: true },
          ],
        },
        [
          { player: 0, kind: IntentKind.choose },
          { player: 0, kind: IntentKind.choose, choice: 1 },
        ],
      ),
      isCard,
    );
    expect(a.choice).toEqual({
      title: 'Choose what to attack',
      purpose: 1,
      options: [
        {
          label: 'gaper',
          card: 'gaper',
          target: { object: 42, player: -1 },
          intent: { player: 0, kind: IntentKind.choose },
        },
        {
          label: 'monster deck',
          target: { player: -1, deck: true },
          intent: { player: 0, kind: IntentKind.choose, choice: 1 },
        },
      ],
    });
  });

  it('uses the question as the title of an answer', () => {
    const a = actionsFor(
      update({ waiting: { kind: PromptKind.choose, player: 0, purpose: 9, options: ['yes', 'No'], text: 'Pay 3¢?' } }, [
        { player: 0, kind: IntentKind.choose },
      ]),
      isCard,
    );
    expect(a.choice?.title).toBe('Pay 3¢?');
  });

  it('asks for a discard of a set count', () => {
    const a = actionsFor(
      update({ waiting: { kind: PromptKind.discard, player: 0, count: 2 } }, [
        { player: 0, kind: IntentKind.discard, objects: [5, 6, 7] },
      ]),
      isCard,
    );
    expect(a.discard).toMatchObject({ count: 2, from: [5, 6, 7] });
  });

  it('asks for an order with the question as its title (19.16)', () => {
    const text = 'Put the cards back in any order';
    const intent = { player: 0, kind: IntentKind.order, objects: [5, 6, 7] };
    const a = actionsFor(update({ waiting: { kind: PromptKind.order, player: 0, count: 3, text } }, [intent]), isCard);
    expect(a.order).toEqual({ title: text, intent });
    expect(orderIntent(intent, [7, 5, 6]).objects).toEqual([7, 5, 6]);
  });

  it('offers the throw of a roll, named by the prompt (15.61)', () => {
    const intent = { player: 0, kind: IntentKind.throw };
    const a = actionsFor(update({ waiting: { kind: PromptKind.throw, player: 0, text: 'attack roll' } }, [intent]), isCard);
    expect(a.throw).toEqual({ label: 'attack roll', intent });
    expect(a.buttons).toEqual([]);
  });
});

describe('aims', () => {
  const monster = { object: 40, player: -1 };
  const ann = { player: 0, is_player: true };

  it('gives a play its targets by its place in allowed', () => {
    const u = update({}, [
      { player: 0, kind: IntentKind.pass },
      { player: 0, kind: IntentKind.playLoot, objects: [7] },
      { player: 0, kind: IntentKind.playLoot, objects: [8] },
    ]);
    const a = actionsFor({ ...u, aims: [{ intent: 1, targets: [monster, ann] }] }, isCard);
    expect(a.cards.get(7)?.[0]?.targets).toEqual([monster, ann]);
    expect(a.cards.get(8)?.[0]?.targets).toBeUndefined();
  });

  it('keeps the reason an aim gives for wrong clicks', () => {
    const u = update({}, [{ player: 0, kind: IntentKind.playLoot, objects: [7] }]);
    const a = actionsFor({ ...u, aims: [{ intent: 0, targets: [monster], reason: 'Pick a monster.' }] }, isCard);
    expect(a.cards.get(7)?.[0]?.reason).toBe('Pick a monster.');
  });

  it('picks the dropped-on target from the target prompt', () => {
    const u = update({ waiting: { kind: PromptKind.choose, player: 0, purpose: Purpose.target, options: ['gaper', 'player 1', 'cancel'] }, options: [monster, ann, { player: -1 }] }, [
      { player: 0, kind: IntentKind.choose, choice: 0 },
      { player: 0, kind: IntentKind.choose, choice: 1 },
      { player: 0, kind: IntentKind.choose, choice: 2 },
    ]);
    const { choice } = actionsFor(u, isCard);
    expect(aimedOption(choice, ann)?.choice).toBe(1);
    expect(aimedOption(choice, { object: 41, player: -1 })).toBeUndefined();
    expect(aimedOption({ title: '', purpose: Purpose.attackTarget, options: choice?.options ?? [] }, ann)).toBeUndefined();
  });

  it('tells players, stack items and cards apart', () => {
    expect(sameTarget(ann, { player: 0, is_player: true })).toBe(true);
    expect(sameTarget(ann, { object: 0, player: 0 })).toBe(false);
    expect(sameTarget({ stack_id: 2, player: 1 }, { stack_id: 2, player: 0 })).toBe(true);
    expect(sameTarget({ player: -1 }, { player: -1 })).toBe(false);
  });
});

// 15.29: what the player can't use now says why.
describe('blocked', () => {
  it('maps blocked cards and turn actions to their reasons', () => {
    const u = update({}, [{ player: 0, kind: IntentKind.pass }]);
    const a = actionsFor(
      {
        ...u,
        blocked: [
          { kind: IntentKind.attack, object: 0, reason: 'No attacks left' },
          { kind: IntentKind.endTurn, object: 0, reason: 'You must attack first' },
          { kind: IntentKind.playLoot, object: 7, reason: 'No loot plays left' },
          { kind: IntentKind.activate, object: 9, reason: 'Used; it recharges on your turn' },
          { kind: IntentKind.activate, object: 9, reason: 'Not enough ¢' },
        ],
      },
      isCard,
    );
    expect(a.turn.why).toEqual({ attack: 'No attacks left', endTurn: 'You must attack first' });
    expect(a.blocked.get(7)).toBe('No loot plays left');
    expect(a.blocked.get(9)).toBe('Used; it recharges on your turn');
    expect(a.blocked.has(0)).toBe(false);
  });

  it('is empty when the server sends no reasons', () => {
    const a = actionsFor(update({}, []), isCard);
    expect(a.blocked.size).toBe(0);
    expect(a.turn.why).toBeUndefined();
  });
});

// 19.17: each order of replacement effects says what it gives.
describe('outcomeLines', () => {
  const names = new Map([
    ['bum_bo', 'Bum-Bo!'],
    ['counterfeit_penny', 'Counterfeit Penny'],
    ['poker_chip', 'Poker Chip'],
  ]);
  const name = (ref: string) => names.get(ref) ?? ref;

  it('names the card that levels up and the one that does nothing', () => {
    const o = { actions: [{ kind: ActionKind.addCounters, player: 0, amount: 3, card: 'bum_bo' }], unused: ['counterfeit_penny'] };
    expect(outcomeLines(o, name)).toEqual(['Bum-Bo! levels up by 3', 'Counterfeit Penny does nothing']);
  });

  it('says a gain, a roll, a prevention and a second pick', () => {
    expect(outcomeLines({ actions: [{ kind: ActionKind.gainCents, player: 0, amount: 4 }] }, name)).toEqual(['You gain 4¢']);
    expect(outcomeLines({ actions: [{ kind: ActionKind.pushLingering, player: 0, amount: 4, card: 'poker_chip' }] }, name)).toEqual([
      'Poker Chip rolls for 4¢',
    ]);
    expect(outcomeLines({}, name)).toEqual(['Nothing happens']);
    expect(outcomeLines({ again: true }, name)).toEqual(['Then you pick the next card']);
  });

  it('passes the outcome on to the choice option', () => {
    const outcome = { actions: [{ kind: ActionKind.gainCents, player: 0, amount: 4 }] };
    const u = update({ waiting: { kind: PromptKind.choose, player: 0, options: ['counterfeit_penny'] }, options: [{ object: 7, player: -1, outcome }] }, [
      { player: 0, kind: IntentKind.choose, choice: 0 },
    ]);
    expect(actionsFor(u, () => true).choice?.options[0]?.outcome).toEqual(outcome);
  });
});

// 19.19: a trigger option is its card with the ability's words.
describe('trigger order', () => {
  const choice = (options: string[]) =>
    actionsFor(
      update({ waiting: { kind: PromptKind.choose, player: 0, purpose: Purpose.triggerOrder, options } }, [
        { player: 0, kind: IntentKind.choose, choice: 0 },
      ]),
      (l) => l === 'cains_eye',
    ).choice;

  it('says the first one resolves last', () => {
    expect(choice(['cains_eye'])?.title).toMatch(/resolves last/);
  });

  it('splits the label into the card and its ability text', () => {
    expect(choice(['cains_eye: At the start of your turn, look.'])?.options[0]).toMatchObject({
      card: 'cains_eye',
      line: 'At the start of your turn, look.',
    });
  });

  it('keeps a label that names no card as text', () => {
    const o = choice(['unknown: text'])?.options[0];
    expect(o?.card).toBeUndefined();
    expect(o?.line).toBeUndefined();
  });
});
