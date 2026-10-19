import { describe, expect, it } from 'vitest';
import { setList, setName, sortSets } from './sets';

describe('sets', () => {
  it('names known sets and keeps unknown codes', () => {
    expect(setName('g2')).toBe('Gold Box');
    expect(setName('zz')).toBe('zz');
  });
  it('sorts by release, unknown last', () => {
    expect(sortSets(['zz', 'r', 'fsp2', 'b2', 'g2'])).toEqual(['b2', 'g2', 'fsp2', 'r', 'zz']);
  });
  it('lists names for the game list', () => {
    expect(setList(['fsp2', 'b2'])).toBe('Base Game, Four Souls+');
    expect(setList(null)).toBe('');
  });
});
