import { describe, expect, it } from 'vitest';
import type { CardInfo } from './cards';
import { createPlayback, effectsFor, scaleFor, type Effect } from './playback';
import type { GameEvent, Update } from './protocol';

const ctx = {
  card: (ref: string): CardInfo | undefined =>
    ref === 'cain'
      ? { ref, set: 'b2', type: 'Character Card', name: 'Cain', effects: [], footnotes: [], stats: [], rewards: [], image: '' }
      : undefined,
  name: (p: number) => ['Ann', 'Bo'].at(p) ?? '?',
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
      { kind: 'popup', at: ['player-1'], text: '−2', tone: 'bad' },
      { kind: 'banner', text: "Bo's turn" },
    ]);
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
