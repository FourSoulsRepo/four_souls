import { describe, expect, it } from 'vitest';
import { counterColor, playerCounters } from './counters';

describe('counter colours (18.15)', () => {
  it('a general counter is green and a gold counter is gold', () => {
    expect(counterColor('')).toBe('#5cb83a');
    expect(counterColor()).toBe('#5cb83a');
    expect(counterColor('gold')).toBe('#e0b43a');
  });

  it('a kind with no colour yet is grey until the owner picks one', () => {
    expect(counterColor('silver')).toBe('#8a8f98');
  });
});

describe('playerCounters (14.33)', () => {
  it('lists the kinds a player has; fear says what it does', () => {
    expect(playerCounters(undefined)).toEqual([]);
    expect(playerCounters([{ name: 'fear', count: 0 }])).toEqual([]);
    expect(playerCounters([{ name: 'fear', count: 3 }, { name: 'luck', count: 1 }])).toEqual([
      { name: 'fear', count: 3, word: 'Fear', title: 'Fear: 3 — monsters +1 evasion per 2 on the active player' },
      { name: 'luck', count: 1, word: 'luck', title: 'luck: 1' },
    ]);
  });
});
