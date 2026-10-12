import { describe, expect, it } from 'vitest';
import { actionsFor, IntentKind, PromptKind } from './actions';
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
