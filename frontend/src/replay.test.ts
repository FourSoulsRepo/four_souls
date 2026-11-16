import { describe, expect, it } from 'vitest';
import type { RecordStep } from './bridge';
import type { GameView } from './protocol';
import type { CardInfo } from './cards';
import type { GameEvent, PlayerView, Update } from './protocol';
import { animates, autoplayPause, clampStep, hideHands, replayUpdate, sliderFill, turnAt, withoutRecordCards, withRecordCards } from './replay';

const turns = [
  { number: 0, active: 0, step: 0 },
  { number: 1, active: 1, step: 2 },
  { number: 2, active: 0, step: 20 },
];

// The replay player's pure parts (10.2, ADR 031).
describe('replay', () => {
  it('makes a step an update with no intents, aims or deadline', () => {
    const step: RecordStep = { index: 4, n: 7, events: [{ kind: 'gained_cents', player: 1, amount: 3 }], view: { turn: { number: 2 } } as GameView };
    const u = replayUpdate(step, { seats: [{ seat: 0, name: 'Ann' }, { seat: 1, name: 'Bo' }] });
    expect(u).toEqual({
      step: 7,
      events: step.events,
      view: step.view,
      allowed: null,
      seats: [
        { seat: 0, name: 'Ann', connected: true },
        { seat: 1, name: 'Bo', connected: true },
      ],
    });
  });

  it('keeps a move inside the record', () => {
    expect(clampStep(-1, 30)).toBe(0);
    expect(clampStep(30, 30)).toBe(29);
    expect(clampStep(12.4, 30)).toBe(12);
  });

  it('plays effects only one step forward; back and jumps show at once', () => {
    expect(animates(4, 5)).toBe(true);
    expect(animates(5, 4)).toBe(false);
    expect(animates(4, 9)).toBe(false);
    expect(animates(4, 4)).toBe(false);
  });

  it('finds the turn a step is in', () => {
    expect(turnAt(turns, 0)).toBe(0);
    expect(turnAt(turns, 1)).toBe(0);
    expect(turnAt(turns, 2)).toBe(1);
    expect(turnAt(turns, 25)).toBe(2);
    expect(turnAt([{ number: 1, active: 0, step: 3 }], 1)).toBe(-1);
  });

  it('fills the slider from the first step to the last (10.7)', () => {
    expect(sliderFill(0, 5)).toBe(0);
    expect(sliderFill(2, 5)).toBe(50);
    expect(sliderFill(4, 5)).toBe(100);
    expect(sliderFill(9, 5)).toBe(100);
    expect(sliderFill(0, 1)).toBe(100);
  });

  it('pauses 0.8 s after events and 0.2 s after none at 1×, divided by the speed', () => {
    expect(autoplayPause(2, 1)).toBe(800);
    expect(autoplayPause(0, 1)).toBe(200);
    expect(autoplayPause(1, 4)).toBe(200);
    expect(autoplayPause(0, 0.5)).toBe(400);
  });

  it('hides every hand, the peeked treasure and private cards (RP-04)', () => {
    const card = (id: number, ref: string) => ({ id, card: ref });
    const player = (id: number): PlayerView =>
      ({ id, hand_size: 2, hand: [card(10 + id, 'a_penny'), card(20 + id, 'bomb')], treasure_top: card(30, 'the_d6') }) as unknown as PlayerView;
    const looted: GameEvent = { kind: 'looted', player: 1, object: 21, card: 'bomb', private: true };
    const gained: GameEvent = { kind: 'gained_cents', player: 0, amount: 3 };
    const u: Update = {
      step: 9,
      events: [looted, gained],
      allowed: null,
      seats: [],
      view: {
        players: [player(0), player(1)],
        waiting: { kind: 4, player: 1, options: [{ key: 'ask.yes' }] },
        subject: card(40, 'a_dime'),
        ordering: [card(41, 'a_nickel')],
        options: [{ object: 21 }],
        activating: 21,
        modes: [{ text: 'Gain 3¢.', option: 0 }],
      } as unknown as GameView,
    };
    const hidden = hideHands(u);
    expect(hidden.view.players.map((p) => [p.hand, p.hand_size, p.treasure_top])).toEqual([
      [undefined, 2, undefined],
      [undefined, 2, undefined],
    ]);
    expect(hidden.events).toEqual([{ kind: 'looted', player: 1, private: true }, gained]);
    expect(hidden.view).toMatchObject({ waiting: { kind: 4, player: 1 } });
    expect(hidden.view.waiting.options).toBeUndefined();
    const { subject, ordering, options, activating, modes } = hidden.view;
    expect([subject, ordering, options, activating, modes]).toEqual([undefined, undefined, undefined, undefined, undefined]);
    // The record's update itself keeps every hand.
    expect(u.view.players[0]?.hand).toHaveLength(2);
    expect(hideHands({ ...u, events: null }).events).toBeNull();
  });

  it("keeps a prompt's options that name no hidden card", () => {
    const u = { step: 1, events: null, allowed: null, view: { players: [], waiting: { kind: 1, player: 0, options: [{ key: 'ask.yes' }] } } } as unknown as Update;
    expect(hideHands(u).view.waiting.options).toEqual([{ key: 'ask.yes' }]);
  });

  it("adds the record's text for cards the app lacks; the app's data wins (RP-10)", () => {
    const own: CardInfo = { ref: 'a_penny', set: 'b2', type: 'Loot Card', name: 'A Penny!', effects: ['Gain 1¢.'], footnotes: [], stats: [], rewards: [], image: 'penny.png' };
    const app = new Map([[own.ref, own]]);
    const header = [
      { ref: 'a_penny', name: 'Old Penny', type: 'Loot Card', effects: ['Gain 2¢.'] },
      { ref: 'zz_future', name: 'Future Card', type: 'Treasure Card', effects: ['Gain 9¢.'] },
    ];
    const merged = withRecordCards(app, header);
    expect(merged.get('a_penny')).toBe(own);
    expect(merged.get('zz_future')).toEqual({
      ref: 'zz_future',
      set: '',
      type: 'Treasure Card',
      name: 'Future Card',
      effects: ['Gain 9¢.'],
      footnotes: [],
      stats: [],
      rewards: [],
      image: '',
      record: true,
    });
    expect(app.has('zz_future')).toBe(false);
    expect(withRecordCards(app, header.slice(0, 1))).toBe(app);
    // Leaving the replay drops them again.
    const back = withoutRecordCards(merged);
    expect([...back.keys()]).toEqual(['a_penny']);
    expect(withoutRecordCards(app)).toBe(app);
  });
});
