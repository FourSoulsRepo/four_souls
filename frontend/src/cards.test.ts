import { describe, expect, it } from 'vitest';
import { iconText, kindOf, textParts, unknownCard } from './cards';

describe('kindOf', () => {
  it('sorts every Base Game type into its deck family', () => {
    const cases: Record<string, string> = {
      'Character Card': 'character',
      'Bonus Soul Card': 'soul',
      'Eternal Treasure Card': 'treasure',
      'Soul Treasure Card': 'treasure',
      'Basic Monster Card': 'monster',
      'Epic Boss Card': 'monster',
      'Holy/Charmed Monster Card': 'monster',
      'Good Event Card': 'monster',
      'Curse Card': 'monster',
      'Trinket Card': 'loot',
      'Pill/Rune Card': 'loot',
      '1 Cent Card': 'loot',
      'Lost Soul Card': 'loot',
    };
    for (const [type, kind] of Object.entries(cases)) {
      expect(kindOf(type), type).toBe(kind);
    }
  });
});

describe('textParts', () => {
  it('splits icons out of the text', () => {
    expect(textParts('{Tap Effect}Roll. Gain {HP} now.')).toEqual([
      { icon: 'Tap Effect' },
      { text: 'Roll. Gain ' },
      { icon: 'HP' },
      { text: ' now.' },
    ]);
  });

  it('keeps text without icons, or with a stray brace, as is', () => {
    expect(textParts('Plain.')).toEqual([{ text: 'Plain.' }]);
    expect(textParts('a { b')).toEqual([{ text: 'a { b' }]);
    expect(textParts('')).toEqual([]);
  });
});

describe('iconText', () => {
  it('shortens effect icons', () => {
    expect(iconText('Tap Effect')).toBe('↷');
    expect(iconText('LV10 Effect')).toBe('LV10');
    expect(iconText('ATK')).toBe('ATK');
  });
});

describe('unknownCard', () => {
  it('shows the ref as its name', () => {
    expect(unknownCard('fan_card').name).toBe('fan_card');
  });
});
