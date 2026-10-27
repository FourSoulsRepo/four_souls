import { describe, expect, it } from 'vitest';
import type { CardInfo } from './cards';
import { createPlayback, effectsFor, scaleFor, slotZone, type Death, type Effect } from './playback';
import type { GameEvent, GameView, Update } from './protocol';

const ctx = {
  card: (ref: string): CardInfo | undefined =>
    ref === 'cain'
      ? { ref, set: 'b2', type: 'Character Card', name: 'Cain', effects: [], footnotes: [], stats: [], rewards: [], image: '' }
      : undefined,
  name: (p: number) => ['Ann', 'Bo'].at(p) ?? '?',
  character: (p: number) => (p === 1 ? 31 : undefined),
  slot: (id: number) => (id === 50 ? 'monster-slot-1' : undefined),
  before: (): string | undefined => undefined,
  diedAt: new Map<number, Death>(),
  seen: (card: string) => card === 'pin',
  roll: (stack: number) => (stack === 7 ? 2 : undefined),
  you: 0,
};

// monster knows every card as a monster (19.73).
const monster = (ref: string): CardInfo => ({ ref, set: 'b2', type: 'Monster Card', name: ref, effects: [], footnotes: [], stats: [], rewards: [], image: '' });

describe('effectsFor', () => {
  it('throws each first-player die, then names a tie or who goes first (15.43)', () => {
    const events: GameEvent[] = [
      { kind: 'dice_rolled', player: 1, amount: 2, text: 'first', label: { key: 'event.first_player_roll' } },
      { kind: 'first_roll_tied', player: -1, amount: 2, text: '0,1' },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'dice', value: 2, label: 'first player roll', player: 1, thrown: true },
      { kind: 'banner', text: 'Tie at 2: You and Bo roll again', long: true },
    ]);
    expect(effectsFor([{ kind: 'first_player', player: 1, amount: 1 }], ctx)).toEqual([{ kind: 'banner', text: 'Bo goes first — rolled 1', long: true }]);
    expect(effectsFor([{ kind: 'first_player', player: 0, card: 'cain', text: 'character' }], ctx)).toEqual([
      { kind: 'banner', text: 'You go first — Cain', long: true },
    ]);
  });

  // Each card from the loot deck is dealt with its sound (18.7).
  it('flies looted cards from the loot deck to the hand', () => {
    const e: GameEvent = { kind: 'looted', player: 1, object: 9, card: 'a_dime', private: true };
    expect(effectsFor([e, { ...e, object: 10, card: 'a_penny' }], ctx)).toEqual([
      { kind: 'fly', from: ['deck-loot'], to: ['hand-1', 'player-1'], card: 'a_dime', back: 'loot', sound: 'lootDraw' },
      { kind: 'fly', from: ['deck-loot'], to: ['hand-1', 'player-1'], card: 'a_penny', back: 'loot', sound: 'lootDraw' },
    ]);
  });

  // A card taken from the loot discard makes the same sound (18.7).
  it('flies a card looted from the discard with its sound', () => {
    const e: GameEvent = { kind: 'looted', player: 1, object: 9, card: 'a_dime', text: 'from the discard' };
    expect(effectsFor([e], ctx)).toEqual([
      { kind: 'fly', from: ['discard-loot'], to: ['hand-1', 'player-1'], card: 'a_dime', back: 'loot', sound: 'lootDraw' },
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
      { kind: 'dice_rolled', player: 0, amount: 5, label: { key: 'event.attack_roll' } },
      { kind: 'gained_cents', player: 0, amount: 3 },
      { kind: 'damage_pending', player: 1, amount: 2 },
      { kind: 'turn_started', player: 1, amount: 2 },
      { kind: 'passed_priority', player: 0 },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'dice', value: 5, label: 'attack roll', player: 0, thrown: true },
      { kind: 'popup', at: ['player-0'], text: '+3¢', tone: 'good' },
      { kind: 'banner', text: "Bo's turn" },
    ]);
  });

  it("shows another player's played loot card large, not the viewer's own (15.31)", () => {
    const events: GameEvent[] = [
      { kind: 'loot_played', player: 1, object: 70, card: 'bomb' },
      { kind: 'loot_played', player: 0, object: 71, card: 'a_penny' },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'played', card: 'bomb', who: 'Bo', from: ['hand-1', 'player-1'], to: ['stack'] },
    ]);
  });

  it('moves stolen cents from the victim to the thief (15.75)', () => {
    const e: GameEvent = { kind: 'stole', player: 1, amount: 2, text: '0' };
    expect(effectsFor([e], ctx)).toEqual([{ kind: 'steal', from: ['player-0'], to: ['player-1'], amount: 2 }]);
    // Without a victim only the thief's text shows.
    expect(effectsFor([{ ...e, text: '' }], ctx)).toEqual([
      { kind: 'popup', at: ['player-1'], text: '+2¢ stolen', tone: 'good' },
    ]);
  });

  it('flies a loot card that leaves the stack from the stack (15.74)', () => {
    const e: GameEvent = { kind: 'discarded', player: 1, object: 80, prev: 70, card: 'bomb', text: 'stack' };
    expect(effectsFor([e], ctx)).toEqual([
      { kind: 'fly', from: ['played-70', 'stack-src-70', 'obj-70', 'deck-loot'], to: ['discard-loot'], card: 'bomb', back: 'loot', leaves: 70 },
    ]);
  });

  it('plays the hits and deaths of a resolved loot before its discard flight (15.51)', () => {
    const discarded: GameEvent = { kind: 'discarded', player: 0, object: 80, prev: 70, card: 'bomb', text: 'stack' };
    const damaged: GameEvent = { kind: 'damaged', player: 1, amount: 1, source: 70, by: 1 };
    const died: GameEvent = { kind: 'died', player: 1 };
    const kinds = (events: GameEvent[]) => effectsFor(events, ctx).map((e) => e.kind);
    expect(kinds([discarded, damaged, died])).toEqual(['hits', 'deaths', 'fly']);
    // The shot flies from the card, which is still there.
    expect(effectsFor([discarded, damaged], ctx)[0]).toMatchObject({ kind: 'hits', hits: [{ from: ['obj-70', 'played-70', 'stack-src-70', 'player-0'] }] });
    expect(kinds([damaged, died, discarded])).toEqual(['hits', 'deaths', 'fly']);
    // Damage from another source does not hold the card.
    expect(kinds([discarded, { ...damaged, source: 9 }])).toEqual(['fly', 'hits']);
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
    const from = ['obj-7', 'played-7', 'stack-src-7'];
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

  it('flies damage from its source, else from the player who dealt it (15.77)', () => {
    const hits = (e: GameEvent) => {
      const first = effectsFor([e], ctx).at(0);
      return first?.kind === 'hits' ? first.hits.at(0)?.from : undefined;
    };
    // Bo (by 2) has the character 31; Ann (by 1) has none known.
    expect(hits({ kind: 'damaged', player: 0, amount: 1, source: 7, by: 2 })).toEqual([
      'obj-7',
      'played-7',
      'stack-src-7',
      'obj-31',
      'player-1',
    ]);
    expect(hits({ kind: 'damaged', player: 1, amount: 1, by: 1 })).toEqual(['player-0']);
    expect(hits({ kind: 'damaged', player: 1, amount: 1 })).toEqual([]);
  });

  it('marks dying cards with skulls, where they were', () => {
    const events: GameEvent[] = [
      { kind: 'died', player: -1, object: 90, prev: 40, card: 'fly' },
      { kind: 'died', player: 1, object: 31, card: 'cain' },
    ];
    expect(effectsFor(events, ctx)).toEqual([{ kind: 'deaths', at: [['obj-40'], ['obj-31']] }]);
  });

  it('throws a fresh roll; a changed roll only spins in (15.62)', () => {
    const events: GameEvent[] = [
      { kind: 'dice_rolled', player: 0, amount: 2, text: 'reroll' },
      { kind: 'roll_changed', player: 0, amount: 3 },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'dice', value: 2, label: 'reroll', player: 0, thrown: true },
      { kind: 'dice', value: 3, label: "Ann's roll: 3", player: 0 },
    ]);
  });

  it('throws every die of a roll of several dice (19.39)', () => {
    const events: GameEvent[] = [{ kind: 'dice_rolled', player: 1, amount: 2, dice: [2, 6, 3, 5], text: 'reroll' }];
    expect(effectsFor(events, ctx)).toEqual([{ kind: 'dice', value: 2, label: 'reroll', player: 1, thrown: true, dice: [2, 6, 3, 5] }]);
  });

  it('changes a roll from its old value, from the card that changed it (15.76)', () => {
    // The shown stack holds roll 7 at 2; roll 8 is thrown in this update.
    const events: GameEvent[] = [
      { kind: 'roll_changed', player: 1, amount: 3, stack_id: 7, source: 44 },
      { kind: 'roll_changed', player: 1, amount: 4, stack_id: 7 },
      { kind: 'dice_rolled', player: 0, amount: 5, stack_id: 8 },
      { kind: 'roll_changed', player: 0, amount: 6, stack_id: 8, source: 45 },
    ];
    expect(effectsFor(events, ctx).filter((e) => e.kind === 'dice' && !e.thrown)).toEqual([
      { kind: 'dice', value: 3, label: "Bo's roll: 2 → 3", player: 1, old: 2, source: ['obj-44', 'stack-src-44'] },
      { kind: 'dice', value: 4, label: "Bo's roll: 3 → 4", player: 1, old: 3 },
      { kind: 'dice', value: 6, label: "Ann's roll: 5 → 6", player: 0, old: 5, source: ['obj-45', 'stack-src-45'] },
    ]);
  });

  it('places a card put into a slot before the attack roll after it (19.3)', () => {
    const events: GameEvent[] = [
      { kind: 'card_revealed', player: -1, object: 50, card: 'fly' },
      { kind: 'attack_target', player: 0, object: 50, card: 'fly' },
      { kind: 'dice_rolled', player: 0, amount: 4, label: { key: 'event.attack_roll' } },
    ];
    expect(effectsFor(events, ctx)).toEqual([
      { kind: 'reveal', card: 'fly', deck: ['deck-loot'], hold: true },
      { kind: 'place', from: ['deck-loot'], to: ['monster-slot-1'], card: 'fly', back: 'loot', middle: true },
      { kind: 'dice', value: 4, label: 'attack roll', player: 0, thrown: true },
    ]);
  });

  it('shows a card revealed in the update before large only once (19.4)', () => {
    const e: GameEvent = { kind: 'card_revealed', player: -1, object: 50, card: 'pin' };
    expect(effectsFor([e], ctx)).toEqual([
      { kind: 'place', from: ['deck-loot'], to: ['monster-slot-1'], card: 'pin', back: 'loot' },
    ]);
  });

  it('neither reveals nor places a card a dead monster uncovered (19.72)', () => {
    const e: GameEvent = { kind: 'card_revealed', player: -1, object: 50, prev: 41, card: 'fly' };
    expect(effectsFor([e], ctx)).toEqual([]);
  });

  it('flies a dead monster from its slot to the discard or its player, also in a later update (19.73)', () => {
    // Monster 40 dies in slot 0 (held object 90); one with a soul, 41,
    // dies in slot 2 (held 91) and becomes a soul after a reward roll.
    const before = (id: number) => new Map([[40, 'monster-slot-0'], [41, 'monster-slot-2']]).get(id);
    const c = { ...ctx, card: monster, before, diedAt: new Map<number, Death>() };
    const deaths: GameEvent[] = [
      { kind: 'died', player: -1, object: 90, prev: 40, card: 'fly' },
      { kind: 'died', player: -1, object: 91, prev: 41, card: 'boss' },
      { kind: 'discarded', player: -1, object: 92, prev: 90, card: 'fly' },
    ];
    expect(effectsFor(deaths, c)).toEqual([
      { kind: 'deaths', at: [['obj-40'], ['obj-41']] },
      { kind: 'fly', from: ['monster-slot-0'], to: ['discard-monster'], card: 'fly', back: 'monster', lifts: 40 },
    ]);
    expect(effectsFor([{ kind: 'dice_rolled', player: 1, amount: 4, label: { key: 'event.reward_roll' } }], c)).toHaveLength(1);
    const soul: GameEvent = { kind: 'gained_soul', player: 1, object: 93, prev: 91, card: 'boss' };
    expect(effectsFor([soul], c)).toEqual([
      { kind: 'fly', from: ['monster-slot-2'], to: ['player-1'], card: 'boss', back: 'monster', lifts: 41 },
    ]);
    // Each is flown once.
    expect(c.diedAt.size).toBe(0);
  });

  it('flies another soul from its place; a bonus soul not at all (19.73)', () => {
    const c = { ...ctx, diedAt: new Map<number, Death>() };
    const chest: GameEvent = { kind: 'gained_soul', player: 0, object: 60, prev: 55, card: 'holy_chest' };
    expect(effectsFor([chest], c)).toEqual([
      { kind: 'fly', from: ['obj-55'], to: ['player-0'], card: 'holy_chest', back: 'monster', lifts: 55 },
    ]);
    const bonus: GameEvent = { kind: 'gained_soul', player: 0, object: 61, card: 'soul_of_envy', text: 'bonus soul' };
    expect(effectsFor([bonus], c)).toEqual([]);
  });

  it('only reveals a card that is in no slot', () => {
    const e: GameEvent = { kind: 'card_revealed', player: 0, object: 51, card: 'fly' };
    expect(effectsFor([e], ctx)).toEqual([{ kind: 'reveal', card: 'fly', deck: ['deck-loot'] }]);
  });

  it('holds a reveal only for a move of its card right after it (15.73)', () => {
    const reveal: GameEvent = { kind: 'card_revealed', player: -1, object: 51, card: 'fly', text: 'loot deck' };
    const discard = (card: string): GameEvent => ({ kind: 'discarded', player: -1, object: 60, card });
    expect(effectsFor([reveal, discard('fly')], ctx)).toEqual([
      { kind: 'reveal', card: 'fly', deck: ['deck-loot'], hold: true },
      { kind: 'fly', from: ['obj-60', 'deck-loot'], to: ['discard-loot'], card: 'fly', back: 'loot', middle: true },
    ]);
    expect(effectsFor([reveal, discard('cain')], ctx)[0]).toEqual({ kind: 'reveal', card: 'fly', deck: ['deck-loot'] });
    const later: GameEvent[] = [reveal, { kind: 'gained_cents', player: 0, amount: 1 }, discard('fly')];
    expect(effectsFor(later, ctx)[0]).toEqual({ kind: 'reveal', card: 'fly', deck: ['deck-loot'] });
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
      begin: () => undefined,
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
    playback.push({ ...update(1), events: [] });
    playback.push(update(2));
    await playback.idle();
    expect(shown).toEqual([1, 2]);
    expect(played).toEqual([{ kind: 'dice', value: 2, label: '', player: 0, thrown: true }]);
  });

  it("plays the first update's effects once it is drawn: Cain goes first (15.43)", async () => {
    const shown: number[] = [];
    const order: string[] = [];
    let current: Update | null = null;
    const playback = createPlayback({
      begin: () => undefined,
      show: (u) => {
        current = u;
        shown.push(u.step);
        order.push('show');
      },
      run: (e) => {
        order.push(e.kind);
        return Promise.resolve();
      },
      effects: (u) => effectsFor(u.events ?? [], ctx),
      speed: () => 'normal',
      shown: () => current,
    });
    playback.push({ ...update(0), events: [{ kind: 'first_player', player: 1, card: 'cain', text: 'character' }] });
    await playback.idle();
    expect(shown).toEqual([0]);
    expect(order).toEqual(['show', 'banner']);
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
    expect(played.at(-1)).toEqual({ kind: 'dice', value: 30, label: '', player: 0, thrown: true });
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
