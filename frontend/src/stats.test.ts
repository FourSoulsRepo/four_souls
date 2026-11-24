import { describe, expect, it } from 'vitest';
import protocol from './protocol.ts?raw';
import type { StatSource } from './protocol';
import en from './locales/en.json';
import { chipText, damageTaken, effect, effectIcons, effectList, hitsText, oddsLine, rewardsLine, rollChanges, shownEffects, shownIcons, statLines, sumOf, trend } from './stats';

const names = new Map([
  ['polyphemus', 'Polyphemus'],
  ['xi_strength', 'XI. Strength'],
  ['meat', 'Meat!'],
  ['round_worm', 'Round Worm'],
  ['champion_belt', 'Champion Belt'],
  ['curse', 'Curse of the Blind'],
  ['trick_treat_normhardultra', 'Trick'],
  ['daddy_haunt', 'Daddy Haunt'],
  ['steamy_sale', 'Steamy Sale!'],
  ['curse_of_the_soulless', 'Curse of the Soulless'],
  ['tax_for_the_mighty', 'Tax for the Mighty'],
  ['get_out_of_jail_card', 'Get Out of Jail Card'],
  ['soul_heart', 'Soul Heart'],
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
    const belt: StatSource = { stat: 'atk', source: 7, card: 'champion_belt', amount: 1, lasts: 'in_play', when: { key: 'ask.when.first_attack_roll' } };
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

// Damage taken (18.29). Trick: "1-2: Damage you would take is increased
// by 1 till end of turn", rolled twice on Hard; Daddy Haunt: "take that
// much damage +1 instead". So +3: +2 from Trick this turn, +1 lasting.
describe('damage taken (18.29)', () => {
  const trick: StatSource = { stat: 'damage_taken', source: 3, card: 'trick_treat_normhardultra', amount: 1, lasts: 'turn' };
  const haunt: StatSource = { stat: 'damage_taken', source: 8, card: 'daddy_haunt', amount: 1, lasts: 'in_play' };
  it('sums every source', () => {
    expect(damageTaken([trick, haunt, trick, poly])).toBe(3);
    expect(damageTaken([poly])).toBe(0);
    expect(damageTaken(undefined)).toBe(0);
  });
  it('says the sum with one line per source, this turn first (owner, 19.232)', () => {
    expect(statLines([haunt, trick, trick], { hp: 2, atk: 1 }, { atk: 1 }, name)).toEqual(['Damage taken +3 = +2 Trick (this turn) + 1 Daddy Haunt']);
    expect(statLines([trick], { hp: 2, atk: 1 }, { atk: 1 }, name)).toEqual(['Damage taken +1 = +1 Trick (this turn)']);
  });
  it('comes after the stats and rolls', () => {
    expect(statLines([trick, poly, meat], { hp: 2, atk: 1 }, { atk: 3 }, name)).toEqual([
      'ATK 3 = 1 printed + 2 Polyphemus',
      'Attack rolls: +1 Meat!',
      'Damage taken +1 = +1 Trick (this turn)',
    ]);
  });
  it('gives the card of this turn one icon with one line', () => {
    expect(effectIcons([trick, trick, haunt], name)).toEqual([{ key: '3-turn', card: 'trick_treat_normhardultra', nextRoll: false, lines: ['+2 damage taken this turn'], number: '+2' }]);
  });
  it('adds up one source in the other sums too', () => {
    expect(statLines([strength, strength], { hp: 2, atk: 1 }, { atk: 3 }, name)).toEqual(['ATK 3 = 1 printed + 2 XI. Strength (this turn)']);
  });
});

describe('effectIcons', () => {
  it('gives an icon to effects of this turn and the next roll, not to lasting ones', () => {
    expect(effectIcons([poly, strength, meat, worm], name)).toEqual([
      { key: '9-turn', card: 'xi_strength', nextRoll: false, lines: ['+1 ATK this turn'], number: '+1' },
      { key: '12-next_roll', card: 'round_worm', nextRoll: true, lines: ['+1 to the next attack roll'], number: '+1' },
    ]);
  });
  it('joins two changes of one source into one icon', () => {
    const hp: StatSource = { ...strength, stat: 'hp' };
    expect(effectIcons([strength, hp], name)).toEqual([{ key: '9-turn', card: 'xi_strength', nextRoll: false, lines: ['+1 ATK this turn', '+1 max HP this turn'] }]);
  });
  it('drops a next-roll icon once the roll used its boost (15.93)', () => {
    expect(effectIcons([meat, worm], name).map((t) => t.key)).toEqual(['12-next_roll']);
    // The view after the roll no longer sends the used boost.
    expect(effectIcons([meat], name)).toEqual([]);
  });
  it('keeps a hidden source apart, with no card', () => {
    expect(effectIcons([{ stat: 'atk', amount: 1, lasts: 'turn' }, { stat: 'hp', amount: 1, lasts: 'turn' }], name)).toEqual([
      { key: 'x0-turn', nextRoll: false, lines: ['+1 ATK this turn'], number: '+1' },
      { key: 'x1-turn', nextRoll: false, lines: ['+1 max HP this turn'], number: '+1' },
    ]);
  });
  it('gives words to every effect; the sums leave the new ones out (18.31)', () => {
    const jail: StatSource = { stat: 'lock_turn', source: 7, card: 'get_out_of_jail_card', amount: 1, lasts: 'turn' };
    const shield = { stat: 'shield', source: 8, card: 'the_poop', lasts: 'turn' } as StatSource; // no amount: the whole next damage
    expect(effectIcons([jail, shield, { ...jail, stat: 'hp' }], name)).toEqual([
      { key: '7-turn', card: 'get_out_of_jail_card', nextRoll: false, lines: ["Others can't play loot cards or activate items this turn", '+1 max HP this turn'] },
      { key: 'shields', card: 'the_poop', nextRoll: false, lines: ['Prevents the next damage this turn — the_poop'], number: '1', title: '1 shield' },
    ]);
    expect(statLines([jail, shield], { hp: 2, atk: 1 }, { atk: 1 }, name)).toEqual([]);
  });
});

// Effect icons for hidden effects only (18.31, 18.34, ADR 037 item 3).
describe('effectIcons on a player (18.34)', () => {
  // All Hallows' Eve: "Players can't gain souls", a room.
  const eve: StatSource = { stat: 'no_souls', source: 40, card: 'all_hallows_eve', amount: 1, lasts: 'in_play' };
  // Tax for the Mighty: "pay each other player 2¢ to attack", a room.
  const tax: StatSource = { stat: 'attack_tax', source: 41, card: 'tax_for_the_mighty', amount: 2, lasts: 'in_play' };
  // Steamy Sale!, an item in the player's own area: "Shop items you
  // purchase cost 5¢ less."
  const sale: StatSource = { stat: 'shop_price', source: 20, card: 'steamy_sale', amount: -5, lasts: 'in_play' };
  // Get Out of Jail Card, played and gone.
  const jail: StatSource = { stat: 'lock_turn', source: 7, card: 'get_out_of_jail_card', amount: 1, lasts: 'turn' };
  // Shields (R-MECH-46): Soul Heart "Prevent the next 1 damage"; The
  // Poop prevents the whole next damage (no amount); The Jar's cap.
  const heart: StatSource = { stat: 'shield', source: 46, card: 'soul_heart', amount: 1, lasts: 'turn' };
  const poop = { stat: 'shield', source: 8, card: 'the_poop', lasts: 'turn' } as StatSource;
  const jar: StatSource = { stat: 'shield_cap', source: 47, card: 'the_jar', amount: 1, lasts: 'turn' };
  it('gives no icon to a lasting effect: own area, room, monster alike (item 3.1)', () => {
    expect(effectIcons([poly, eve, tax, sale], name)).toEqual([]);
    const baby: StatSource = { stat: 'dc', source: 3, card: 'baby_haunt', amount: 1, lasts: 'in_play' };
    expect(effectIcons([baby], name)).toEqual([]);
  });
  it('keeps a this-turn icon of an own card and of a gone one (item 3.8)', () => {
    const ownTurn: StatSource = { ...strength, source: 5, card: 'polyphemus' };
    expect(effectIcons([ownTurn, sale, jail], name).map((i) => i.key)).toEqual(['5-turn', '7-turn']);
  });
  it('adds up one source of one stat, and shows the sum (items 3.4, 3.5)', () => {
    // Holy Chest's 1–2: "You may attack an additional time this turn",
    // twice from one source.
    const extra: StatSource = { stat: 'attacks', source: 42, card: 'holy_chest', amount: 1, lasts: 'turn' };
    expect(effectIcons([extra, extra], name)).toEqual([{ key: '42-turn', card: 'holy_chest', nextRoll: false, lines: ['+2 attacks this turn'], number: '+2' }]);
    // Two stats of one source: two lines, no number.
    const icons = effectIcons([strength, { ...strength, stat: 'hp' }], name);
    expect(icons).toHaveLength(1);
    expect(icons[0]?.number).toBeUndefined();
    // One stat with two conditions: one icon, its lines sum up.
    const first: StatSource = { ...strength, when: { key: 'ask.when.first_attack_roll' } };
    expect(effectIcons([strength, first], name)[0]?.number).toBe('+2');
  });
  it('shows no number for a flag, a cap, the attack DC or any number of attacks', () => {
    const turn = (s: StatSource): StatSource => ({ ...s, lasts: 'turn' });
    const tiny: StatSource = { stat: 'hand_cap', source: 43, card: 'curse_of_tiny_hands', amount: 2, lasts: 'turn' };
    const stab: StatSource = { stat: 'player_attack_dc', source: 44, card: 'backstabber', amount: 3, lasts: 'turn' };
    const key: StatSource = { stat: 'attacks', source: 45, card: 'skeleton_key', amount: 98, lasts: 'turn' };
    const icons = effectIcons([turn(eve), tiny, stab, key, jail], name);
    expect(icons.map((i) => i.number)).toEqual([undefined, undefined, undefined, undefined, undefined]);
    expect(icons.map((i) => i.lines[0])).toEqual([
      "Can't gain souls this turn",
      'Max hand size 2 this turn',
      'Can be attacked, at DC 3 this turn',
      'May attack any number of times this turn',
      "Others can't play loot cards or activate items this turn",
    ]);
  });
  it('gives three shields one icon "3", one line each in the order they block (item 3.15)', () => {
    // The server sends them oldest first, the order combat.go uses them.
    expect(effectIcons([jail, heart, poop, jar], name)).toEqual([
      { key: '7-turn', card: 'get_out_of_jail_card', nextRoll: false, lines: ["Others can't play loot cards or activate items this turn"] },
      {
        key: 'shields',
        card: 'soul_heart',
        nextRoll: false,
        lines: ['Prevents the next 1 damage this turn — Soul Heart', 'Prevents the next damage this turn — the_poop', 'The next damage is at most 1 this turn — the_jar'],
        number: '3',
        title: '3 shields',
      },
    ]);
  });
  it('keeps two shields of one card as two lines, never a sum', () => {
    expect(effectIcons([poop, poop], name)).toEqual([
      {
        key: 'shields',
        card: 'the_poop',
        nextRoll: false,
        lines: ['Prevents the next damage this turn — the_poop', 'Prevents the next damage this turn — the_poop'],
        number: '2',
        title: '2 shields',
      },
    ]);
  });
  it("takes the shields' art from the first known card; a hidden one keeps its place", () => {
    const hidden: StatSource = { stat: 'shield', amount: 2, lasts: 'turn' };
    const icons = effectIcons([hidden, heart], name);
    expect(icons.map((i) => i.card)).toEqual(['soul_heart']);
    expect(icons.map((i) => i.lines)).toEqual([['Prevents the next 2 damage this turn — a hidden card', 'Prevents the next 1 damage this turn — Soul Heart']]);
  });
  it('says what a shop or $ change does', () => {
    expect(effect({ stat: 'loot_plays', amount: 1, lasts: 'turn' })).toBe('+1 loot play');
    expect(effect({ stat: 'shop_price', amount: -1000, lasts: 'turn' })).toBe('May buy a shop item for 0¢');
    expect(effect({ stat: 'paid_discount', amount: 1, lasts: 'turn' })).toBe('$ abilities of your items cost 1¢ less');
  });
  it('has words for every stat the protocol names', () => {
    // The union is one line of protocol.ts (a Go test reads it too).
    const line = protocol.split('\n').find((l) => l.trim().startsWith('stat: "hp"')) ?? '';
    const stats = [...line.matchAll(/"([a-z_]+)"/g)].map((m) => m[1] as StatSource['stat']);
    expect(stats.length).toBeGreaterThan(30);
    for (const stat of stats) {
      expect(effect({ stat, amount: 1, lasts: 'turn' }), stat).not.toBe('');
    }
    expect(Object.keys(en).filter((k) => k.startsWith('stats.effect.')).length).toBeGreaterThan(30);
  });
});

describe('shownIcons: four, then a "⋯ K" chip (ADR 037 items 3.6, 3.16)', () => {
  it('shows up to four icons and no chip', () => {
    expect(shownIcons([1, 2, 3, 4])).toEqual({ shown: [1, 2, 3, 4], more: 0 });
    expect(shownIcons([])).toEqual({ shown: [], more: 0 });
  });
  it('shows the first four and counts the rest', () => {
    expect(shownIcons([1, 2, 3, 4, 5])).toEqual({ shown: [1, 2, 3, 4], more: 1 });
    expect(shownIcons([1, 2, 3, 4, 5, 6, 7])).toEqual({ shown: [1, 2, 3, 4], more: 3 });
  });
  it('writes the chip as "⋯ K", never an icon\'s "+N" (item 3.16)', () => {
    expect(chipText(2)).toBe('⋯ 2');
    expect(chipText(2)).not.toMatch(/^\+/);
  });
});

// The zoom's list of the hidden effects on a player (18.32, 18.34, ADR
// 037 item 4).
describe('effectList (18.34)', () => {
  // The own area. Steamy Sale!: "Shop items you purchase cost 5¢
  // less." Curse of the Soulless: "You can't gain souls."
  const sale: StatSource = { stat: 'shop_price', source: 20, card: 'steamy_sale', amount: -5, lasts: 'in_play' };
  const soulless: StatSource = { stat: 'no_souls', source: 21, card: 'curse_of_the_soulless', amount: 1, lasts: 'in_play' };
  // A room. Tax for the Mighty: "... must pay each other player 2¢ to
  // attack."
  const tax: StatSource = { stat: 'attack_tax', source: 41, card: 'tax_for_the_mighty', amount: 2, lasts: 'in_play' };
  // Played and gone. Get Out of Jail Card: "... other players can't play
  // loot cards or activate items" till end of turn; Soul Heart: "Prevent
  // the next 1 damage they would take this turn."
  const jail: StatSource = { stat: 'lock_turn', source: 7, card: 'get_out_of_jail_card', amount: 1, lasts: 'turn' };
  const heart: StatSource = { stat: 'shield', source: 46, card: 'soul_heart', amount: 1, lasts: 'turn' };
  it("lists this turn's hidden effects only, as the icons (item 4.1)", () => {
    expect(effectList([sale, jail, tax, soulless, heart], name)).toEqual({
      turn: ["Others can't play loot cards or activate items — Get Out of Jail Card", 'Prevents the next 1 damage — Soul Heart'],
    });
  });
  it('gives a lasting room or own-area effect no line (item 4.1)', () => {
    expect(effectList([sale, sale, tax, soulless], name)).toEqual({ turn: [] });
  });
  it('says what it does, then the card; a hidden one is "a hidden card" (items 4.3, 1.6)', () => {
    const hidden: StatSource = { stat: 'purchases', amount: 1, lasts: 'turn' };
    expect(effectList([hidden], name).turn).toEqual(['+1 purchase — a hidden card']);
  });
  it('repeats no sum line: HP, ATK, rolls, attacked DC, damage taken (item 4.4)', () => {
    const taken: StatSource = { stat: 'damage_taken', source: 30, card: 'trick_treat_normhardultra', amount: 1, lasts: 'turn' };
    const attacked: StatSource = { stat: 'attacked_dc', source: 31, card: 'curse', amount: -1, lasts: 'turn' };
    const roll: StatSource = { stat: 'roll', source: 32, card: 'meat', amount: 1, lasts: 'turn' };
    const hp: StatSource = { ...strength, stat: 'hp' };
    expect(effectList([poly, strength, hp, meat, worm, roll, attacked, taken], name)).toEqual({ turn: [] });
    expect(effectList(undefined, name)).toEqual({ turn: [] });
  });
  it('lists shields in the order they block, two of one card apart; one source adds up', () => {
    const poop = { stat: 'shield', source: 8, card: 'the_poop', lasts: 'turn' } as StatSource; // the whole next damage
    expect(effectList([poop, poop, heart], name).turn).toEqual([
      'Prevents the next damage — the_poop',
      'Prevents the next damage — the_poop',
      'Prevents the next 1 damage — Soul Heart',
    ]);
    // Holy Chest's 1–2: "You may attack an additional time this turn".
    const extra: StatSource = { stat: 'attacks', source: 50, card: 'holy_chest', amount: 1, lasts: 'turn' };
    expect(effectList([extra, extra], name).turn).toEqual(['+2 attacks — holy_chest']);
  });
});

describe('shownEffects: eight lines, then a count (18.32)', () => {
  const lines = (n: number, tag: string) => Array.from({ length: n }, (_, i) => `${tag}${String(i)}`);
  it('shows up to eight lines and no count', () => {
    expect(shownEffects({ turn: lines(8, 't') })).toEqual({ turn: lines(8, 't'), more: 0 });
  });
  it('counts the rest', () => {
    expect(shownEffects({ turn: lines(11, 't') })).toEqual({ turn: lines(8, 't'), more: 3 });
  });
});

describe('hitsText', () => {
  it('says the lowest face that hits, or the faces when a 2 or 3 turns into a 1', () => {
    expect(hitsText([4, 5, 6])).toBe('Hits on 4+');
    expect(hitsText([1, 2, 3, 4, 5, 6])).toBe('Always hits');
    expect(hitsText([])).toBe("Can't hit");
    expect(hitsText([1, 4, 5, 6])).toBe('Hits on 1, 4, 5, 6');
  });
});

describe('oddsLine', () => {
  it('has the DC and the ATK alone without changes', () => {
    expect(oddsLine({ dc: 4, hits: [4, 5, 6], atk: 1 }, 1, name)).toBe('Hits on 4+ (DC 4) · ATK 1');
  });
  it('takes a roll change off the DC', () => {
    expect(oddsLine({ dc: 5, hits: [4, 5, 6], atk: 1, sources: [meat] }, 1, name)).toBe('Hits on 4+ (DC 5 − 1 Meat!) · ATK 1');
  });
  it('names two roll changes and the ATK of a hit with its condition', () => {
    const belt: StatSource = { stat: 'atk', source: 7, card: 'champion_belt', amount: 1, lasts: 'in_play', when: { key: 'ask.when.first_attack_roll' } };
    const curse: StatSource = { stat: 'attack_roll', source: 8, card: 'curse', amount: -2, lasts: 'next_roll' };
    expect(oddsLine({ dc: 3, hits: [4, 5, 6], atk: 2, sources: [meat, curse, belt] }, 1, name)).toBe(
      'Hits on 4+ (DC 3 − 1 Meat! + 2 Curse of the Blind (next roll)) · ATK 2 = 1 printed + 1 Champion Belt (first attack roll)',
    );
  });
});

describe('rollChanges', () => {
  it('lists roll changes only', () => {
    expect(rollChanges([meat, worm, poly], name)).toBe('+1 Meat!, +1 Round Worm (next roll)');
    expect(rollChanges(undefined, name)).toBe('');
  });
});

describe('rewardsLine (14.33)', () => {
  it('says the rewards a monster gives this turn; none changed, no line', () => {
    expect(rewardsLine(undefined)).toBe('');
    expect(rewardsLine([])).toBe('');
    expect(rewardsLine([{ kind: 'cents', amount: 1 }])).toBe('Rewards this turn: 1¢');
    expect(
      rewardsLine([
        { kind: 'cents', amount: 3 },
        { kind: 'loot', amount: 1 },
        { kind: 'treasure', amount: 2 },
        { kind: 'cents', roll: true },
      ]),
    ).toBe('Rewards this turn: 3¢ · 1 loot · 2 treasures · ¢ by a roll');
  });

  it('adds up fixed amounts of one kind; a roll stays apart', () => {
    const rewards = [
      { kind: 'loot', amount: 2 },
      { kind: 'cents', roll: true },
      { kind: 'loot', amount: 1 },
      { kind: 'cents', amount: 3 },
    ] as const;
    expect(rewardsLine(rewards.map((r) => ({ ...r })))).toBe('Rewards this turn: 3 loot · ¢ by a roll · 3¢');
    expect(rewards[0].amount).toBe(2);
  });
});
