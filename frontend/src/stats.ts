// Stat sources (V-41, 15.92): why a card's health, dice or attack is not
// what it prints, and which effects of this turn change it. Built from
// the view's sources (15.54).
import { t } from './i18n';
import type { OddsView, StatSource } from './protocol';

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

// statLines are the sums the zoom shows beside a card: "ATK 4 = 1 printed
// + 2 Polyphemus + 1 XI. Strength (this turn)", one line per stat that
// has sources, then the roll changes. Empty without sources.
export function statLines(sources: StatSource[] | undefined, printed: Printed, live: Live, name: Name): string[] {
  // Lasting sources first, then this turn's: the print, then what stays.
  const list = [...(sources ?? [])].sort((a, b) => Number(a.lasts !== 'in_play') - Number(b.lasts !== 'in_play'));
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
  return lines;
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
// sources show only in the sums (owner, 2026-10-22).
export function tempSources(sources: StatSource[] | undefined): TempSource[] {
  const out: TempSource[] = [];
  (sources ?? []).forEach((s, i) => {
    if (s.lasts === 'in_play') return;
    const nextRoll = s.lasts === 'next_roll';
    const key = `${s.source !== undefined ? String(s.source) : `x${String(i)}`}-${s.lasts}`;
    const line = `${nextRoll ? effect(s) : t('stats.this_turn_effect', { effect: effect(s) })}${s.when ? ` (${s.when})` : ''}`;
    const known = out.find((t) => t.key === key);
    if (known) {
      known.lines.push(line);
      return;
    }
    out.push({ key, ...(s.card ? { card: s.card } : {}), nextRoll, lines: [line] });
  });
  return out;
}

// effect says what one source does: "+1 ATK", "+1 to the next attack roll".
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
  }
}

// suffix says how long a source lasts and when it applies.
function suffix(s: StatSource): string {
  const notes = [s.lasts === 'turn' ? t('stats.this_turn') : s.lasts === 'next_roll' ? t('stats.next_roll') : '', s.when ?? ''].filter((x) => x !== '');
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
