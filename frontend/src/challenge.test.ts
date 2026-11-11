import { describe, expect, it } from 'vitest';
import { challengeEnd, monsterRows, presents, timerText } from './challenge';
import type { GameView } from './protocol';
import { findCard } from './stack';

const names = { player: (p: number) => ['Ann', 'Bo'].at(p) ?? '?', card: (ref: string) => ref };

// krampus is a challenge view: the ¢ present face down, the loot one
// destroyed, the treasure one given to Bo, the soul one face up.
function krampus(over = false, lost = false): GameView {
  return {
    over,
    players: [],
    discards: [],
    monsters: [{ top: { id: 1, card: 'gaper' } }],
    bosses: [{ top: { id: 2, card: 'pilfering_krampus_norm' } }],
    minions: [{ top: { id: 3, card: 'fly' } }],
    challenge: {
      id: 'krampus',
      difficulty: 'normal',
      timer: 7,
      card: { id: 4, card: 'how_the_krampus_stole_christmas_norm', charged: true },
      presents: [
        { kind: 'cents', counters: 1, player: -1 },
        { kind: 'loot', counters: 0, destroyed: true, card: 'the_d6', player: -1 },
        { kind: 'treasure', counters: 0, flipped: true, card: 'the_d20', player: 1 },
        { kind: 'soul', counters: 0, flipped: true, card: 'the_d4', player: -1 },
      ],
      ...(lost ? { lost: true, reason: { key: 'event.challenge_lost.presents' } } : {}),
    },
  } as unknown as GameView;
}

describe('monsterRows', () => {
  it('lists the monster slots, then the final boss, then the minions (ADR 028, 7.3)', () => {
    expect(monsterRows(krampus()).map((s) => s.top?.id)).toEqual([1, 2, 3]);
    expect(monsterRows({ monsters: null } as unknown as GameView)).toEqual([]);
  });

  it('lets the stack and history name a boss, a minion and the rule card', () => {
    expect(findCard(krampus(), 2)).toBe('pilfering_krampus_norm');
    expect(findCard(krampus(), 3)).toBe('fly');
    expect(findCard(krampus(), 4)).toBe('how_the_krampus_stole_christmas_norm');
  });
});

describe('presents', () => {
  it('names each present by its kind and what became of it (R-KRAMPUS-05)', () => {
    expect(presents(krampus(), names.player)).toEqual([
      { kind: 'cents', caption: '¢', state: 'down', label: '¢ present: face down, 1 counter', counters: 1 },
      { kind: 'loot', caption: 'Loot', state: 'destroyed', label: 'Loot present: destroyed', counters: 0, card: 'the_d6' },
      { kind: 'treasure', caption: 'Treasure', state: 'given', label: 'Treasure present: Bo has it as an item', counters: 0, card: 'the_d20' },
      { kind: 'soul', caption: 'Soul', state: 'up', label: 'Soul present: face up, for a player to gain', counters: 0, card: 'the_d4' },
    ]);
  });

  it('has none in a regular game', () => {
    expect(presents({ players: [] } as unknown as GameView, names.player)).toEqual([]);
  });
});

describe('timerText', () => {
  it('counts the rounds left', () => {
    expect(timerText(8)).toBe('8 rounds left');
    expect(timerText(1)).toBe('1 round left');
  });
});

describe('challengeEnd', () => {
  it('says won, or lost with the reason, once the game is over (ADR 028, 11.6)', () => {
    expect(challengeEnd(krampus(), names)).toBeNull();
    expect(challengeEnd(krampus(true), names)).toEqual({ title: 'Challenge won!', reason: '' });
    expect(challengeEnd(krampus(true, true), names)).toEqual({ title: 'Challenge lost', reason: 'All 4 presents were destroyed' });
  });

  it('leaves a regular game to its winners', () => {
    expect(challengeEnd({ over: true, players: [] } as unknown as GameView, names)).toBeNull();
  });
});
