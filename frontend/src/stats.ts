// Stat sources (V-41, 15.92): why a card's health, dice or attack is not
// what it prints, and which effects of this turn change it. Built from
// the view's sources (15.54).
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
  const formula = (label: string, stat: StatSource['stat'], base: number, now: number, clamp = '') => {
    const parts = list.filter((s) => s.stat === stat);
    if (parts.length === 0) return;
    const terms = parts.map((s) => `${s.amount < 0 ? '−' : '+'} ${String(Math.abs(s.amount))} ${name(s.card)}${suffix(s)}`);
    lines.push(`${label} ${String(now)} = ${String(base)} printed ${terms.join(' ')}${clamp}`);
  };
  formula('Max HP', 'hp', printed.hp, printed.hp + sumOf(list, 'hp'));
  formula('ATK', 'atk', printed.atk, live.atk);
  if (printed.dc !== undefined && live.dc !== undefined) {
    const clamped = printed.dc + sumOf(list, 'dc') !== live.dc;
    formula('DC', 'dc', printed.dc, live.dc, clamped ? ' (dice stay 1–6)' : '');
  }
  const rolls = (label: string, stat: StatSource['stat'], unit = '') => {
    const parts = list.filter((s) => s.stat === stat);
    if (parts.length > 0) lines.push(`${label}: ${parts.map((s) => `${signed(s.amount)}${unit} ${name(s.card)}${suffix(s)}`).join(', ')}`);
  };
  rolls('Rolls', 'roll');
  rolls('Attack rolls', 'attack_roll');
  rolls('Monsters you attack', 'attacked_dc', ' DC');
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
    const line = `${effect(s)}${nextRoll ? '' : ' this turn'}${s.when ? ` (${s.when})` : ''}`;
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
  const n = signed(s.amount);
  switch (s.stat) {
    case 'hp':
      return `${n} max HP`;
    case 'atk':
      return `${n} ATK`;
    case 'dc':
      return `${n} DC`;
    case 'roll':
      return `${n} to rolls`;
    case 'attack_roll':
      return s.lasts === 'next_roll' ? `${n} to the next attack roll` : `${n} to attack rolls`;
    case 'attacked_dc':
      return `${n} DC of monsters you attack`;
  }
}

// suffix says how long a source lasts and when it applies.
function suffix(s: StatSource): string {
  const notes = [s.lasts === 'turn' ? 'this turn' : s.lasts === 'next_roll' ? 'next roll' : '', s.when ?? ''].filter((x) => x !== '');
  return notes.length > 0 ? ` (${notes.join(', ')})` : '';
}

// hitsText says which die faces hit (15.94): "Hits on 4+", or the faces
// when they do not run up to 6 (a 2 or 3 that becomes a 1).
export function hitsText(hits: number[]): string {
  if (hits.length === 0) return "Can't hit";
  if (hits.length === 6) return 'Always hits';
  const first = hits[0] ?? 7;
  const upward = hits.every((d, i) => d === first + i) && hits.at(-1) === 6;
  return upward ? `Hits on ${String(first)}+` : `Hits on ${hits.join(', ')}`;
}

const rollStats = new Set<StatSource['stat']>(['roll', 'attack_roll']);

// oddsLine is the open attack's fight line (15.94): "Hits on 4+ (DC 5 −
// 1 Meat!) · ATK 3 = 1 printed + 1 Champion Belt (first attack roll)".
// A roll change lowers the face needed, so it is taken off the DC.
export function oddsLine(odds: OddsView, printedATK: number, name: Name): string {
  const list = odds.sources ?? [];
  const rolls = list.filter((s) => rollStats.has(s.stat));
  const terms = rolls.map((s) => ` ${s.amount > 0 ? '−' : '+'} ${String(Math.abs(s.amount))} ${name(s.card)}${suffix(s)}`).join('');
  const hit = `${hitsText(odds.hits)} (DC ${String(odds.dc)}${terms})`;
  const atk = list.filter((s) => s.stat === 'atk');
  const sum = atk.map((s) => ` ${s.amount < 0 ? '−' : '+'} ${String(Math.abs(s.amount))} ${name(s.card)}${suffix(s)}`).join('');
  return `${hit} · ATK ${String(odds.atk)}${atk.length > 0 ? ` = ${String(printedATK)} printed${sum}` : ''}`;
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
