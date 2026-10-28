import { describe, expect, it } from 'vitest';
import { counterColor } from './counters';

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
