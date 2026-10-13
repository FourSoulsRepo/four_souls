import { describe, expect, it } from 'vitest';
import { actionsFor, attackStatus, IntentKind, PromptKind } from './actions';
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
      update({ waiting: { kind: PromptKind.choose, player: 0, purpose: 9, options: ['yes', 'no'], text: 'Pay 3¢?' } }, [
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
});

describe('attackStatus', () => {
  const me = (attacks: number) => ({
    id: 0,
    character: { id: 1, card: 'isaac', charged: true },
    hp: 2,
    atk: 1,
    cents: 0,
    souls: 0,
    in_play: null,
    hand_size: 0,
    loot_plays: 0,
    attacks,
    purchases: 0,
  });

  it('is ready while the viewer may attack, asleep with no attack left (T-11)', () => {
    const ready = update({ players: [me(1)] }, [{ player: 0, kind: IntentKind.attack }]);
    expect(attackStatus(actionsFor(ready, isCard), ready.view, 0)).toBe('ready');
    const spent = update({ players: [me(0)] }, [{ player: 0, kind: IntentKind.endTurn }]);
    expect(attackStatus(actionsFor(spent, isCard), spent.view, 0)).toBe('sleep');
  });

  it('marks nothing while an attack waits, or on another player’s turn', () => {
    const waiting = update({ players: [me(1)] }, [{ player: 0, kind: IntentKind.pass }]);
    expect(attackStatus(actionsFor(waiting, isCard), waiting.view, 0)).toBeUndefined();
    const theirs = update({ players: [me(0), { ...me(1), id: 1 }] }, null);
    expect(attackStatus(actionsFor(theirs, isCard), theirs.view, 1)).toBeUndefined();
  });
});
