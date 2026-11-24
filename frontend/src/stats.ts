// Stat sources (V-41, 15.92): why a card's health, dice or attack is not
// what it prints, and which effects of this turn change it. Built from
// the view's sources (15.54).
import { say, t } from './i18n';
import type { OddsView, RewardView, StatSource } from './protocol';

export type Trend = 'up' | 'down' | undefined;

// Name gives a source card's name; an unknown or hidden one has none.
export type Name = (card: string | undefined) => string;

// trend says whether a stat went up or down from its print.
export function trend(live: number, printed: number): Trend {
  if (live > printed) return 'up';
  if (live < printed) return 'down';
  return undefined;
}

// sumOf adds up the sources of one stat.
export function sumOf(sources: StatSource[] | undefined, stat: StatSource['stat']): number {
  return (sources ?? []).filter((s) => s.stat === stat).reduce((n, s) => n + amountOf(s), 0);
}

// amountOf is a source's amount; a shield of the whole next damage has
// none (18.30).
function amountOf(s: StatSource): number {
  return s.amount ?? 0;
}

// shields are the shield keys (ADR 037 item 1.8): one source each, never
// summed, and a whole-damage shield has no amount.
const shields = new Set<StatSource['stat']>(['shield', 'shield_cap']);

// Printed are the stats a card prints; dc only on monsters.
export interface Printed {
  hp: number;
  atk: number;
  dc?: number;
}

// Live are the stats now, as the engine counts them.
export interface Live {
  atk: number;
  dc?: number;
}

// merged adds up one source's parts of one stat (18.29, ADR 037): two
// "+1" of the Trick card this turn are one "+2". Parts that last or
// apply differently stay apart, and so does each shield (18.31). The
// order is each first part's.
export function merged(sources: StatSource[]): StatSource[] {
  const sums = new Map<string, StatSource>();
  sources.forEach((s, i) => {
    if (shields.has(s.stat)) {
      sums.set(`shield-${String(i)}`, s);
      return;
    }
    const key = JSON.stringify([s.stat, s.source ?? null, s.card ?? '', s.lasts, s.when ?? null]);
    const known = sums.get(key);
    sums.set(key, known ? { ...known, amount: amountOf(known) + amountOf(s) } : s);
  });
  return [...sums.values()].filter((s) => shields.has(s.stat) || amountOf(s) !== 0);
}

// statLines are the sums the zoom shows beside a card: "ATK 4 = 1 printed
// + 2 Polyphemus + 1 XI. Strength (this turn)", one line per stat that
// has sources, then the roll changes, then the damage the player takes.
// Empty without sources.
export function statLines(sources: StatSource[] | undefined, printed: Printed, live: Live, name: Name): string[] {
  // Lasting sources first, then this turn's: the print, then what stays.
  const list = merged([...(sources ?? [])].sort((a, b) => Number(a.lasts !== 'in_play') - Number(b.lasts !== 'in_play')));
  const lines: string[] = [];
  const formula = (stat: StatSource['stat'], base: number, now: number, line: (v: { n: number; base: number; terms: string }) => string) => {
    const parts = list.filter((s) => s.stat === stat);
    if (parts.length === 0) return;
    const terms = parts.map((s) => `${amountOf(s) < 0 ? '−' : '+'} ${String(Math.abs(amountOf(s)))} ${name(s.card)}${suffix(s)}`);
    lines.push(line({ n: now, base, terms: terms.join(' ') }));
  };
  formula('hp', printed.hp, printed.hp + sumOf(list, 'hp'), (v) => t('stats.sum_hp', v));
  formula('atk', printed.atk, live.atk, (v) => t('stats.sum_atk', v));
  if (printed.dc !== undefined && live.dc !== undefined) {
    const clamped = printed.dc + sumOf(list, 'dc') !== live.dc;
    formula('dc', printed.dc, live.dc, (v) => (clamped ? t('stats.sum_dc_clamped', v) : t('stats.sum_dc', v)));
  }
  const rolls = (stat: StatSource['stat'], line: (list: string) => string, unit = '') => {
    const parts = list.filter((s) => s.stat === stat);
    if (parts.length > 0) lines.push(line(parts.map((s) => `${signed(amountOf(s))}${unit} ${name(s.card)}${suffix(s)}`).join(', ')));
  };
  rolls('roll', (l) => t('stats.rolls', { list: l }));
  rolls('attack_roll', (l) => t('stats.attack_rolls', { list: l }));
  rolls('attacked_dc', (l) => t('stats.attacked', { list: l }), ` ${t('stats.dc')}`);
  const taken = takenLine(list, name);
  if (taken) lines.push(taken);
  return lines;
}

// takenLine is the damage a player takes more (18.29, owner 19.232):
// "Damage taken +3 = +2 Trick (this turn) + 1 Daddy Haunt". This turn's
// sources come first, as the owner wrote it. '' without any.
function takenLine(list: StatSource[], name: Name): string {
  const parts = list.filter((s) => s.stat === 'damage_taken').sort((a, b) => Number(a.lasts !== 'turn') - Number(b.lasts !== 'turn'));
  if (parts.length === 0) return '';
  const terms = parts.map((s, i) => `${i === 0 ? signed(amountOf(s)) : ` ${amountOf(s) < 0 ? '−' : '+'} ${String(Math.abs(amountOf(s)))}`} ${name(s.card)}${suffix(s)}`);
  return t('stats.sum_damage_taken', { n: signed(sumOf(parts, 'damage_taken')), terms: terms.join('') });
}

// damageTaken is how much more of each damage the player takes now: the
// sum of the view's sources (18.29, ADR 037). The engine adds no less
// than 0 (combat.go), so neither does this.
export function damageTaken(sources: StatSource[] | undefined): number {
  return Math.max(sumOf(sources, 'damage_taken'), 0);
}

// EffectIcon is one icon on a card (ADR 037 item 3): one source card's
// effects, one that waits for the next attack roll (SB-04), or all the
// shields (item 3.15). lines say what it does; number is the sum an
// icon of one number stat shows, or how many shields; title names an
// icon that stands for more than one card.
export interface EffectIcon {
  key: string;
  card?: string;
  nextRoll: boolean;
  lines: string[];
  number?: string;
  title?: string;
}

// SHIELDS is the key of the one icon all of a player's shields share.
const SHIELDS = 'shields';

// effectIcons are a card's icons: hidden effects only (ADR 037 item 3,
// owner, 2026-11-24). A lasting effect never gets one, wherever its card
// lies: the card on the table shows it (item 3.1). This turn's and the
// next roll's do, one icon per source; one source's changes of one stat
// are one line: "+2 damage taken", not two "+1" (18.29). All shields
// share one icon whose number is how many (item 3.15): each blocks one
// damage, oldest first (R-MECH-46), and the server sends them in that
// order. Its lines name each card; its art and zoom are the first known
// one's.
export function effectIcons(sources: StatSource[] | undefined, name: Name): EffectIcon[] {
  const out: (EffectIcon & { parts: StatSource[] })[] = [];
  merged(sources ?? []).forEach((s, i) => {
    if (s.lasts === 'in_play') return;
    const text = effect(s);
    if (text === '') return; // a stat without words: no icon line
    const nextRoll = s.lasts === 'next_roll';
    const when = whenText(s);
    const line = `${s.lasts === 'turn' ? t('stats.this_turn_effect', { effect: text }) : text}${when ? ` (${when})` : ''}`;
    const shield = shields.has(s.stat);
    const key = shield ? SHIELDS : `${s.source !== undefined ? String(s.source) : `x${String(i)}`}-${s.lasts}`;
    const known = out.find((t) => t.key === key);
    const shown = shield ? t('stats.effect_from', { effect: line, card: name(s.card) }) : line;
    if (known) {
      known.lines.push(shown);
      known.parts.push(s);
      if (known.card === undefined && s.card) known.card = s.card;
      return;
    }
    out.push({ key, ...(s.card ? { card: s.card } : {}), nextRoll, lines: [shown], parts: [s] });
  });
  return out.map(({ parts, ...icon }) => {
    if (icon.key === SHIELDS) return { ...icon, number: String(parts.length), title: t('stats.shields', { n: parts.length }) };
    const n = iconNumber(parts);
    return n === undefined ? icon : { ...icon, number: n };
  });
}

// sumStats have a sum line beside the zoomed card (statLines): the
// effect list leaves them out, so nothing shows twice (ADR 037 item 4.4).
const sumStats = new Set<StatSource['stat']>(['hp', 'atk', 'dc', 'roll', 'attack_roll', 'attacked_dc', 'damage_taken']);

// EffectList is the zoomed character's list of the hidden effects on
// its player (ADR 037 item 4): this turn's lines. A lasting effect has
// none (item 4.1), so the list has no "while in play" part (18.34).
export interface EffectList {
  turn: string[];
}

// effectList lists the effects on a player that the icons show (ADR 037
// item 4.1): "Prevents the next 1 damage — Soul Heart", what it does,
// then its card; a hidden card is named so. In the server's order;
// shields in the order they block. One source's parts of one stat add
// up (merged); shields stay apart. Lasting effects, stats with a sum
// line and stats without words are left out.
export function effectList(sources: StatSource[] | undefined, name: Name): EffectList {
  const list: EffectList = { turn: [] };
  for (const s of merged(sources ?? [])) {
    const text = s.lasts === 'in_play' || sumStats.has(s.stat) ? '' : effect(s);
    if (text === '') continue;
    const when = whenText(s);
    list.turn.push(t('stats.effect_from', { effect: `${text}${when ? ` (${when})` : ''}`, card: name(s.card) }));
  }
  return list;
}

// MAX_EFFECT_LINES is how many lines the zoom's effect list shows; a last
// line counts the rest (18.32). Eight lines of two rows each, their
// titles and the count stand no taller than the zoomed card in a 720 px
// high window.
export const MAX_EFFECT_LINES = 8;

// shownEffects cuts the list to MAX_EFFECT_LINES and says how many
// more there are. Nine lines are eight and "1 more": the count takes the
// ninth one's place, as the "⋯ K" chip does.
export function shownEffects(list: EffectList): EffectList & { more: number } {
  if (list.turn.length <= MAX_EFFECT_LINES) return { ...list, more: 0 };
  return { turn: list.turn.slice(0, MAX_EFFECT_LINES), more: list.turn.length - MAX_EFFECT_LINES };
}

// MAX_ICONS is how many icons a card shows; the rest are a "⋯ K" chip
// (ADR 037 items 3.6, 3.16).
export const MAX_ICONS = 4;

// shownIcons splits a card's icons into the ones it draws and how many
// more the "⋯ K" chip stands for. Five icons are four and "⋯ 1": the chip
// takes the fifth one's place.
export function shownIcons<T>(icons: T[]): { shown: T[]; more: number } {
  if (icons.length <= MAX_ICONS) return { shown: icons, more: 0 };
  return { shown: icons.slice(0, MAX_ICONS), more: icons.length - MAX_ICONS };
}

// chipText is the chip's text: "⋯ 2", a pill with no card art, never
// "+2", which an icon of one number stat shows (ADR 037 item 3.16).
export function chipText(more: number): string {
  return `⋯ ${String(more)}`;
}

// numberStats change a number by their amount (ADR 037 item 3.5). Not
// in it: flags (amount 1), caps (max hand size), the DC to be attacked
// (the DC itself) and shields (an amount prevented, not a change).
const numberStats = new Set<StatSource['stat']>([
  'hp',
  'atk',
  'dc',
  'roll',
  'attack_roll',
  'attacked_dc',
  'damage_taken',
  'shop_price',
  'loot_plays',
  'attacks',
  'purchases',
  'loot_step',
  'paid_discount',
  'attack_tax',
]);

// iconNumber is the number an icon shows: the sum of its lines when all
// are one number stat, as a change ("+2", "−5"); undefined otherwise.
// Attacks of 98 or more are "any number" and a shop price of −99 or
// less is free: no number then.
function iconNumber(parts: StatSource[]): string | undefined {
  const stat = parts.at(0)?.stat;
  if (stat === undefined || !numberStats.has(stat) || parts.some((p) => p.stat !== stat || anyAttacks(p) || freeShop(p))) return undefined;
  const sum = parts.reduce((n, p) => n + amountOf(p), 0);
  if (sum === 0) return undefined;
  // A $ ability's discount lowers its cost: −1, as a shop price.
  return signed(stat === 'paid_discount' ? -sum : sum);
}

// anyAttacks: Skeleton Key's 98 and Bum-Bo's 99 more attacks mean any
// number of them (18.30).
function anyAttacks(s: StatSource): boolean {
  return s.stat === 'attacks' && amountOf(s) >= 98;
}

// freeShop: Pity for the Poor's −1000 makes a shop item 0¢.
function freeShop(s: StatSource): boolean {
  return s.stat === 'shop_price' && amountOf(s) <= -99;
}

// effect says what one source does: "+1 ATK", "Can't gain souls"; ''
// for a stat it has no words for (a newer server's).
export function effect(s: StatSource): string {
  const n = Math.abs(amountOf(s));
  const amount = signed(amountOf(s));
  switch (s.stat) {
    case 'hp':
      return t('stats.effect.hp', { amount });
    case 'atk':
      return t('stats.effect.atk', { amount });
    case 'dc':
      return t('stats.effect.dc', { amount });
    case 'roll':
      return t('stats.effect.roll', { amount });
    case 'attack_roll':
      return s.lasts === 'next_roll' ? t('stats.effect.next_attack_roll', { amount }) : t('stats.effect.attack_roll', { amount });
    case 'attacked_dc':
      return t('stats.effect.attacked_dc', { amount });
    case 'damage_taken':
      return t('stats.effect.damage_taken', { amount });
    case 'shop_price':
      return freeShop(s) ? t('stats.effect.shop_free') : t('stats.effect.shop_price', { amount });
    case 'loot_plays':
      return t('stats.effect.loot_plays', { amount, n });
    case 'attacks':
      return anyAttacks(s) ? t('stats.effect.attacks_any') : t('stats.effect.attacks', { amount, n });
    case 'purchases':
      return t('stats.effect.purchases', { amount, n });
    case 'loot_step':
      return t('stats.effect.loot_step', { amount, n });
    case 'lock_others':
    case 'lock_turn':
      return t('stats.effect.lock_others');
    case 'loot_double':
      return t('stats.effect.loot_double');
    case 'death_ward':
      return t('stats.effect.death_ward');
    case 'recharge_limit':
      return t('stats.effect.recharge_limit');
    case 'hand_cap':
      return t('stats.effect.hand_cap', { n });
    case 'no_char_recharge':
      return t('stats.effect.no_char_recharge');
    case 'reveal_hand':
      return t('stats.effect.reveal_hand');
    case 'choose_attack':
      return t('stats.effect.choose_attack');
    case 'no_forced_loss':
      return t('stats.effect.no_forced_loss');
    case 'paid_discount':
      return t('stats.effect.paid_discount', { n });
    case 'no_souls':
      return t('stats.effect.no_souls');
    case 'penalty_item_to':
      return t('stats.effect.penalty_item_to');
    case 'player_attack_dc':
      return t('stats.effect.player_attack_dc', { n });
    case 'one_ability':
      return t('stats.effect.one_ability');
    case 'one_loot':
      return t('stats.effect.one_loot');
    case 'attack_tax':
      return t('stats.effect.attack_tax', { n });
    case 'shield':
      return s.amount === undefined ? t('stats.effect.shield_all') : t('stats.effect.shield', { n });
    case 'shield_cap':
      return t('stats.effect.shield_cap', { n });
    default:
      return ''; // a newer server's stat: no icon line (18.30)
  }
}

// whenText is when a source applies, from the server (ADR 021), or ''.
function whenText(s: StatSource): string {
  return say(s.when, { card: (ref) => ref, player: (p) => String(p + 1) });
}

// suffix says how long a source lasts and when it applies.
function suffix(s: StatSource): string {
  const notes = [s.lasts === 'turn' ? t('stats.this_turn') : s.lasts === 'next_roll' ? t('stats.next_roll') : '', whenText(s)].filter((x) => x !== '');
  return notes.length > 0 ? ` (${notes.join(', ')})` : '';
}

// hitsText says which die faces hit (15.94): "Hits on 4+", or the faces
// when they do not run up to 6 (a 2 or 3 that becomes a 1).
export function hitsText(hits: number[]): string {
  if (hits.length === 0) return t('stats.cant_hit');
  if (hits.length === 6) return t('stats.always_hits');
  const first = hits[0] ?? 7;
  const upward = hits.every((d, i) => d === first + i) && hits.at(-1) === 6;
  return upward ? t('stats.hits_from', { n: first }) : t('stats.hits_on', { faces: hits.join(', ') });
}

const rollStats = new Set<StatSource['stat']>(['roll', 'attack_roll']);

// oddsLine is the open attack's fight line (15.94): "Hits on 4+ (DC 5 −
// 1 Meat!) · ATK 3 = 1 printed + 1 Champion Belt (first attack roll)".
// A roll change lowers the face needed, so it is taken off the DC.
export function oddsLine(odds: OddsView, printedATK: number, name: Name): string {
  const list = odds.sources ?? [];
  const rolls = list.filter((s) => rollStats.has(s.stat));
  const terms = rolls.map((s) => ` ${amountOf(s) > 0 ? '−' : '+'} ${String(Math.abs(amountOf(s)))} ${name(s.card)}${suffix(s)}`).join('');
  const hit = t('stats.odds_hit', { hits: hitsText(odds.hits), dc: odds.dc, terms });
  const atk = list.filter((s) => s.stat === 'atk');
  const sum = atk.map((s) => ` ${amountOf(s) < 0 ? '−' : '+'} ${String(Math.abs(amountOf(s)))} ${name(s.card)}${suffix(s)}`).join('');
  return `${hit} · ${atk.length > 0 ? t('stats.odds_atk_sum', { n: odds.atk, base: printedATK, terms: sum }) : t('stats.odds_atk', { n: odds.atk })}`;
}

// rollChanges names what changes a roll: "+1 Meat!, −2 Round Worm (next
// roll)"; empty without any.
export function rollChanges(sources: StatSource[] | undefined, name: Name): string {
  return (sources ?? [])
    .filter((s) => rollStats.has(s.stat))
    .map((s) => `${signed(amountOf(s))} ${name(s.card)}${suffix(s)}`)
    .join(', ');
}

function signed(n: number): string {
  return `${n < 0 ? '−' : '+'}${String(Math.abs(n))}`;
}

// rewardsLine is a monster's changed rewards (14.33, R-TRICK-13,
// R-TRICK-15): "Rewards this turn: 1¢ · 1 loot"; '' when none changed.
export function rewardsLine(rewards: RewardView[] | undefined): string {
  const text = rewardsText(rewards);
  return text === '' ? '' : t('stats.rewards_now', { rewards: text });
}

// rewardsText is the rewards alone, "1¢ · 1 loot", for the mark on
// the card; '' when none changed. Fixed amounts of one kind add up:
// 2 loot and loot 1 added are "3 loot".
export function rewardsText(rewards: RewardView[] | undefined): string {
  const sums: RewardView[] = [];
  for (const r of rewards ?? []) {
    const same = r.roll ? undefined : sums.find((s) => s.kind === r.kind && !s.roll);
    if (same) {
      same.amount = (same.amount ?? 0) + (r.amount ?? 0);
    } else {
      sums.push({ ...r });
    }
  }
  return sums.map(rewardText).join(' · ');
}

function rewardText(r: RewardView): string {
  const n = r.amount ?? 0;
  switch (r.kind) {
    case 'loot':
      return r.roll ? t('stats.reward.roll_loot') : t('stats.reward.loot', { n });
    case 'treasure':
      return r.roll ? t('stats.reward.roll_treasure') : t('stats.reward.treasure', { n });
    default:
      return r.roll ? t('stats.reward.roll_cents') : t('stats.reward.cents', { n });
  }
}
