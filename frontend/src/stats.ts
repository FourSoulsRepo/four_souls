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
  return (sources ?? []).filter((s) => s.stat === stat).reduce((n, s) => n + s.amount, 0);
}

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
// apply differently stay apart. The order is each first part's.
export function merged(sources: StatSource[]): StatSource[] {
  const sums = new Map<string, StatSource>();
  for (const s of sources) {
    const key = JSON.stringify([s.stat, s.source ?? null, s.card ?? '', s.lasts, s.when ?? null]);
    const known = sums.get(key);
    sums.set(key, known ? { ...known, amount: known.amount + s.amount } : s);
  }
  return [...sums.values()].filter((s) => s.amount !== 0);
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
    const terms = parts.map((s) => `${s.amount < 0 ? '−' : '+'} ${String(Math.abs(s.amount))} ${name(s.card)}${suffix(s)}`);
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
    if (parts.length > 0) lines.push(line(parts.map((s) => `${signed(s.amount)}${unit} ${name(s.card)}${suffix(s)}`).join(', ')));
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
  const terms = parts.map((s, i) => `${i === 0 ? signed(s.amount) : ` ${s.amount < 0 ? '−' : '+'} ${String(Math.abs(s.amount))}`} ${name(s.card)}${suffix(s)}`);
  return t('stats.sum_damage_taken', { n: signed(sumOf(parts, 'damage_taken')), terms: terms.join('') });
}

// damageTaken is how much more of each damage the player takes now: the
// sum of the view's sources (18.29, ADR 037). The engine adds no less
// than 0 (combat.go), so neither does this.
export function damageTaken(sources: StatSource[] | undefined): number {
  return Math.max(sumOf(sources, 'damage_taken'), 0);
}

// TempSource is one icon on a card: an effect of this turn or one that
// waits for the next attack roll (SB-04). lines say what it does.
export interface TempSource {
  key: string;
  card?: string;
  nextRoll: boolean;
  lines: string[];
}

// tempSources are a card's icons, one per source and lasting: lasting
// sources show only in the sums (owner, 2026-10-22). One source's
// changes of one stat are one line: "+2 damage taken", not two "+1"
// (18.29).
export function tempSources(sources: StatSource[] | undefined): TempSource[] {
  const out: TempSource[] = [];
  merged(sources ?? []).forEach((s, i) => {
    if (s.lasts === 'in_play') return;
    const nextRoll = s.lasts === 'next_roll';
    const key = `${s.source !== undefined ? String(s.source) : `x${String(i)}`}-${s.lasts}`;
    const text = effect(s);
    if (text === '') return; // no words for this stat yet (18.31)
    const when = whenText(s);
    const line = `${nextRoll ? text : t('stats.this_turn_effect', { effect: text })}${when ? ` (${when})` : ''}`;
    const known = out.find((t) => t.key === key);
    if (known) {
      known.lines.push(line);
      return;
    }
    out.push({ key, ...(s.card ? { card: s.card } : {}), nextRoll, lines: [line] });
  });
  return out;
}

// effect says what one source does: "+1 ATK", "+1 to the next attack roll";
// '' for a stat it has no words for.
function effect(s: StatSource): string {
  const amount = signed(s.amount);
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
    default:
      return ''; // a stat without words yet: no icon line (18.30)
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
  const terms = rolls.map((s) => ` ${s.amount > 0 ? '−' : '+'} ${String(Math.abs(s.amount))} ${name(s.card)}${suffix(s)}`).join('');
  const hit = t('stats.odds_hit', { hits: hitsText(odds.hits), dc: odds.dc, terms });
  const atk = list.filter((s) => s.stat === 'atk');
  const sum = atk.map((s) => ` ${s.amount < 0 ? '−' : '+'} ${String(Math.abs(s.amount))} ${name(s.card)}${suffix(s)}`).join('');
  return `${hit} · ${atk.length > 0 ? t('stats.odds_atk_sum', { n: odds.atk, base: printedATK, terms: sum }) : t('stats.odds_atk', { n: odds.atk })}`;
}

// rollChanges names what changes a roll: "+1 Meat!, −2 Round Worm (next
// roll)"; empty without any.
export function rollChanges(sources: StatSource[] | undefined, name: Name): string {
  return (sources ?? [])
    .filter((s) => rollStats.has(s.stat))
    .map((s) => `${signed(s.amount)} ${name(s.card)}${suffix(s)}`)
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
