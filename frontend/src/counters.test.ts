import { describe, expect, it } from 'vitest';
import { counterColor, fearOf, playerCounters } from './counters';

describe('counter colours (18.15)', () => {
  it('a general counter is green and a gold counter is gold', () => {
    expect(counterColor('')).toBe('#5cb83a');
    expect(counterColor()).toBe('#5cb83a');
    expect(counterColor('gold')).toBe('#e0b43a');
  });

  it('a kind with no colour yet is grey until the owner picks one', () => {
    expect(counterColor('silver')).toBe('#8a8f98');
  });

  it('fear is ghost white (18.26)', () => {
    expect(counterColor('fear')).toBe('#e6ecf5');
  });
});

describe('fearOf (18.26)', () => {
  it("is the player's fear counters, 0 when none", () => {
    expect(fearOf(undefined)).toBe(0);
    expect(fearOf([])).toBe(0);
    expect(fearOf([{ count: 4 }])).toBe(0);
    expect(fearOf([{ name: 'luck', count: 1 }, { name: 'fear', count: 3 }])).toBe(3);
  });
});

describe('playerCounters (14.31)', () => {
  it('lists the kinds a player has; fear shows on the character card instead (18.26)', () => {
    expect(playerCounters(undefined)).toEqual([]);
    expect(playerCounters([{ name: 'luck', count: 0 }])).toEqual([]);
    expect(playerCounters([{ name: 'fear', count: 3 }, { name: 'luck', count: 1 }])).toEqual([{ name: 'luck', count: 1, word: 'luck', title: 'luck: 1' }]);
  });
});
