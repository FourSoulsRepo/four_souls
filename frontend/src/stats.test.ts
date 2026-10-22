import { describe, expect, it } from 'vitest';
import type { StatSource } from './protocol';
import { statLines, sumOf, tempSources, trend } from './stats';

const names = new Map([
  ['polyphemus', 'Polyphemus'],
  ['xi_strength', 'XI. Strength'],
  ['meat', 'Meat!'],
  ['round_worm', 'Round Worm'],
  ['champion_belt', 'Champion Belt'],
  ['curse', 'Curse of the Blind'],
]);
const name = (card: string | undefined) => (card ? (names.get(card) ?? card) : 'a hidden card');

const poly: StatSource = { stat: 'atk', source: 5, card: 'polyphemus', amount: 2, lasts: 'in_play' };
const strength: StatSource = { stat: 'atk', source: 9, card: 'xi_strength', amount: 1, lasts: 'turn' };
const meat: StatSource = { stat: 'attack_roll', source: 6, card: 'meat', amount: 1, lasts: 'in_play' };
const worm: StatSource = { stat: 'attack_roll', source: 12, card: 'round_worm', amount: 1, lasts: 'next_roll' };

describe('trend', () => {
  it('says up, down, or nothing against the print', () => {
    expect(trend(3, 1)).toBe('up');
    expect(trend(1, 2)).toBe('down');
    expect(trend(2, 2)).toBeUndefined();
  });
});

describe('sumOf', () => {
  it('adds one stat only', () => {
    expect(sumOf([poly, strength, meat], 'atk')).toBe(3);
    expect(sumOf(undefined, 'atk')).toBe(0);
  });
});

describe('statLines', () => {
  it('sums ATK from its print, lasting and this turn', () => {
    expect(statLines([poly, strength], { hp: 2, atk: 1 }, { atk: 4 }, name)).toEqual(['ATK 4 = 1 printed + 2 Polyphemus + 1 XI. Strength (this turn)']);
  });
  it('lists roll changes after the stats, with the next roll and a condition', () => {
    const belt: StatSource = { stat: 'atk', source: 7, card: 'champion_belt', amount: 1, lasts: 'in_play', when: 'first attack roll' };
    expect(statLines([belt, meat, worm], { hp: 2, atk: 1 }, { atk: 2 }, name)).toEqual([
      'ATK 2 = 1 printed + 1 Champion Belt (first attack roll)',
      'Attack rolls: +1 Meat!, +1 Round Worm (next roll)',
    ]);
  });
  it('shows a monster DC down, and says when the dice clamp', () => {
    const blind: StatSource = { stat: 'dc', source: 3, card: 'curse', amount: -1, lasts: 'turn' };
    expect(statLines([blind], { hp: 3, atk: 1, dc: 4 }, { atk: 1, dc: 3 }, name)).toEqual(['DC 3 = 4 printed − 1 Curse of the Blind (this turn)']);
    const big: StatSource = { stat: 'dc', source: 3, card: 'polyphemus', amount: 2, lasts: 'in_play' };
    expect(statLines([big], { hp: 3, atk: 1, dc: 5 }, { atk: 1, dc: 6 }, name)).toEqual(['DC 6 = 5 printed + 2 Polyphemus (dice stay 1–6)']);
  });
  it('sums max HP and names a hidden source', () => {
    expect(statLines([{ stat: 'hp', amount: 1, lasts: 'turn' }], { hp: 2, atk: 1 }, { atk: 1 }, name)).toEqual(['Max HP 3 = 2 printed + 1 a hidden card (this turn)']);
  });
  it('says the DC change of attacked monsters', () => {
    expect(statLines([{ stat: 'attacked_dc', source: 4, card: 'meat', amount: -1, lasts: 'in_play' }], { hp: 2, atk: 1 }, { atk: 1 }, name)).toEqual([
      'Monsters you attack: −1 DC Meat!',
    ]);
  });
  it('is empty without sources', () => {
    expect(statLines(undefined, { hp: 2, atk: 1 }, { atk: 1 }, name)).toEqual([]);
  });
});

describe('tempSources', () => {
  it('gives an icon to effects of this turn and the next roll, not to lasting ones', () => {
    expect(tempSources([poly, strength, meat, worm])).toEqual([
      { key: '9-turn', card: 'xi_strength', nextRoll: false, lines: ['+1 ATK this turn'] },
      { key: '12-next_roll', card: 'round_worm', nextRoll: true, lines: ['+1 to the next attack roll'] },
    ]);
  });
  it('joins two changes of one source into one icon', () => {
    const hp: StatSource = { ...strength, stat: 'hp' };
    expect(tempSources([strength, hp])).toEqual([{ key: '9-turn', card: 'xi_strength', nextRoll: false, lines: ['+1 ATK this turn', '+1 max HP this turn'] }]);
  });
  it('drops a next-roll icon once the roll used its boost (15.93)', () => {
    expect(tempSources([meat, worm]).map((t) => t.key)).toEqual(['12-next_roll']);
    // The view after the roll no longer sends the used boost.
    expect(tempSources([meat])).toEqual([]);
  });
    it('keeps a hidden source apart, with no card', () => {
    expect(tempSources([{ stat: 'atk', amount: 1, lasts: 'turn' }, { stat: 'hp', amount: 1, lasts: 'turn' }])).toEqual([
      { key: 'x0-turn', nextRoll: false, lines: ['+1 ATK this turn'] },
      { key: 'x1-turn', nextRoll: false, lines: ['+1 max HP this turn'] },
    ]);
  });
});
