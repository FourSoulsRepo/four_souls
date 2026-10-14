import { describe, expect, it } from 'vitest';
import type { CardInfo } from './cards';
import { createPlayback, effectsFor, scaleFor, slotZone, type Effect } from './playback';
import type { GameEvent, GameView, Update } from './protocol';

const ctx = {
  card: (ref: string): CardInfo | undefined =>
    ref === 'cain'
      ? { ref, set: 'b2', type: 'Character Card', name: 'Cain', effects: [], footnotes: [], stats: [], rewards: [], image: '' }
      : undefined,
  name: (p: number) => ['Ann', 'Bo'].at(p) ?? '?',
  character: (p: number) => (p === 1 ? 31 : undefined),
  slot: (id: number) => (id === 50 ? 'monster-slot-1' : undefined),
  seen: (card: string) => card === 'pin',
};

describe('effectsFor', () => {
  it('flies looted cards from the loot deck to the hand', () => {
    const e: GameEvent = { kind: 'looted', player: 1, object: 9, card: 'a_dime', private: true };
    expect(effectsFor([e], ctx)).toEqual([
      { kind: 'fly', from: ['deck-loot'], to: ['hand-1', 'player-1'], card: 'a_dime', back: 'loot' },
    ]);
  });

  it('flies a discarded card from where it was to its deck discard', () => {
    const e: GameEvent = { kind: 'discarded', player: -1, object: 12, prev: 7, card: 'cain' };
    expect(effectsFor([e], ctx)).toEqual([
      { kind: 'fly', from: ['obj-7', 'deck-character'], to: ['discard-character'], card: 'cain', back: 'character' },
    ]);
  });

  it('rolls dice, pops counters, and names the turn', () => {
    const events: GameEvent[] = [
      { kind: 'dice_rolled', player: 0, amount: 5, text: 'attack roll' },
      { kind: 'gained_cents', player: 0, amount: 3 },
      { kind: 'damage_pending', player: 1, amount: 2 },
      { kind: 'turn_started', player: 1, amount: 2 },
      { kind: 'passed_priority', player: 0 },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'dice', value: 5, label: 'attack roll' },
      { kind: 'popup', at: ['player-0'], text: '+3¢', tone: 'good' },
      { kind: 'banner', text: "Bo's turn" },
    ]);
  });

  it('lunges the attacker at what its attack damaged, with no pending popup', () => {
    const events: GameEvent[] = [
      { kind: 'damaged', player: -1, object: 40, amount: 1, source: 30, text: 'combat' },
      { kind: 'damaged', player: 1, amount: 2, source: 40, text: 'combat' },
      { kind: 'damaged', player: 0, amount: 1, source: 40, text: 'combat' },
      { kind: 'damage_pending', player: -1, object: 40, amount: 1, text: 'combat' },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'lunge', from: ['obj-30'], to: ['obj-40'], amount: 1 },
      { kind: 'lunge', from: ['obj-40'], to: ['obj-31', 'player-1'], amount: 2 },
      { kind: 'lunge', from: ['obj-40'], to: ['player-0'], amount: 1 },
    ]);
  });

  it('hits every card other damage lands on at once, from its source', () => {
    const events: GameEvent[] = [
      { kind: 'damage_pending', player: -1, object: 40, amount: 1 },
      { kind: 'damaged', player: -1, object: 40, amount: 1, source: 7 },
      { kind: 'damaged', player: 1, amount: 2, source: 7 },
      { kind: 'damaged', player: 0, amount: 1 },
      { kind: 'gained_cents', player: 0, amount: 1 },
      { kind: 'damaged', player: -1, object: 41, amount: 1, source: 7 },
    ];
    const from = ['obj-7', 'stack-src-7'];
    expect(effectsFor(events, ctx)).toEqual([
      {
        kind: 'hits',
        hits: [
          { from, to: ['obj-40'], amount: 1 },
          { from, to: ['obj-31', 'player-1'], amount: 2 },
          { from: [], to: ['player-0'], amount: 1 },
        ],
      },
      { kind: 'popup', at: ['player-0'], text: '+1¢', tone: 'good' },
      { kind: 'hits', hits: [{ from, to: ['obj-41'], amount: 1 }] },
    ]);
  });

  it('marks dying cards with skulls, where they were', () => {
    const events: GameEvent[] = [
      { kind: 'died', player: -1, object: 90, prev: 40, card: 'fly' },
      { kind: 'died', player: 1, object: 31, card: 'cain' },
    ];
    expect(effectsFor(events, ctx)).toEqual([{ kind: 'deaths', at: [['obj-40'], ['obj-31']] }]);
  });

  it('places a card put into a slot before the attack roll after it (19.3)', () => {
    const events: GameEvent[] = [
      { kind: 'card_revealed', player: -1, object: 50, card: 'fly' },
      { kind: 'attack_target', player: 0, object: 50, card: 'fly' },
      { kind: 'dice_rolled', player: 0, amount: 4, text: 'attack roll' },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'reveal', card: 'fly' },
      { kind: 'place', from: ['deck-loot'], to: ['monster-slot-1'], card: 'fly', back: 'loot' },
      { kind: 'dice', value: 4, label: 'attack roll' },
    ]);
  });

  it('shows a card revealed in the update before large only once (19.4)', () => {
    const e: GameEvent = { kind: 'card_revealed', player: -1, object: 50, card: 'pin' };
    expect(effectsFor([e], ctx)).toEqual([
      { kind: 'place', from: ['deck-loot'], to: ['monster-slot-1'], card: 'pin', back: 'loot' },
    ]);
  });

  it('only reveals a card that is in no slot', () => {
    const e: GameEvent = { kind: 'card_revealed', player: 0, object: 51, card: 'fly' };
    expect(effectsFor([e], ctx)).toEqual([{ kind: 'reveal', card: 'fly' }]);
  });
});

describe('slotZone', () => {
  it('names the monster or shop slot a card is the top of', () => {
    const view = {
      monsters: [{ top: { id: 7 } }, { top: { id: 8 } }],
      shop: [{}, { top: { id: 9 } }],
    } as unknown as GameView;
    expect(slotZone(view, 8)).toBe('monster-slot-1');
    expect(slotZone(view, 9)).toBe('shop-slot-1');
    expect(slotZone(view, 10)).toBeUndefined();
  });
});

describe('scaleFor', () => {
  it('follows the speed setting', () => {
    expect(scaleFor('normal', 0, 1)).toBe(1);
    expect(scaleFor('slow', 0, 1)).toBe(1.6);
    expect(scaleFor('off', 0, 1)).toBe(0);
  });

  it('speeds up bursts and skips a long queue', () => {
    expect(scaleFor('normal', 3, 1)).toBe(0.4);
    expect(scaleFor('normal', 0, 7)).toBe(0.4);
    expect(scaleFor('normal', 9, 1)).toBe(0);
  });
});

// update makes an update whose turn number tells it apart, with one dice
// roll in its events.
function update(n: number): Update {
  return {
    step: n,
    events: [{ kind: 'dice_rolled', player: 0, amount: n }],
    view: { turn: { number: n } } as Update['view'],
    allowed: null,
    seats: null,
  };
}

describe('createPlayback', () => {
  function setup() {
    const shown: number[] = [];
    const played: Effect[] = [];
    let current: Update | null = null;
    const playback = createPlayback({
      show: (u) => {
        current = u;
        shown.push(u.step);
      },
      run: (e) => {
        played.push(e);
        return new Promise((resolve) => setTimeout(resolve, 1));
      },
      effects: (u) => effectsFor(u.events ?? [], ctx),
      speed: () => 'normal',
      shown: () => current,
    });
    return { playback, shown, played };
  }

  it('shows the first update at once, then plays effects before each view', async () => {
    const { playback, shown, played } = setup();
    playback.push(update(1));
    playback.push(update(2));
    await playback.idle();
    expect(shown).toEqual([1, 2]);
    expect(played).toEqual([{ kind: 'dice', value: 2, label: '' }]);
  });

  it('shows every update of a fast burst, in order, ending on the last', async () => {
    const { playback, shown, played } = setup();
    for (let n = 1; n <= 30; n++) {
      playback.push(update(n));
    }
    await playback.idle();
    expect(shown).toEqual(Array.from({ length: 30 }, (_, i) => i + 1));
    // A long queue skips effects until it catches up.
    expect(played.length).toBeLessThan(29);
    expect(played.at(-1)).toEqual({ kind: 'dice', value: 30, label: '' });
  });

  it('clear drops what is still queued', async () => {
    const { playback, shown } = setup();
    playback.push(update(1));
    playback.push(update(2));
    playback.push(update(3));
    playback.clear();
    await playback.idle();
    expect(shown.length).toBeLessThan(3);
  });
});
