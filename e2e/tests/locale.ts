// The language the tests run in: LOCALE (make screenshots LOCALE=qps),
// English without it. tr gives a line as that language shows it, so a
// test finds what the player reads (ADR 021).
import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { pseudoLocale, textOf, translate, type Entry, type Values } from '../../frontend/src/locale';

const file = join(import.meta.dirname, '..', '..', 'frontend', 'src', 'locales', 'en.json');
const english = textOf(JSON.parse(readFileSync(file, 'utf8')) as Record<string, Entry>);

// locale is the language every app page starts in.
export const locale = process.env.LOCALE || 'en';

const locales = new Map([
  ['en', english],
  ['qps', pseudoLocale(english)],
]);

// tr is the line key in the tests' language.
export function tr(key: string, values?: Values): string {
  if (!english.has(key)) {
    throw new Error(`no line ${key} in en.json`);
  }
  return translate(locales, locale, key, values);
}

// trRe is the line key in the tests' language as a pattern, to find it
// inside a longer text or with values a test can't know: given values
// fill their places, a place left out matches any text. A plural line
// with no n matches any of its forms. flags as for RegExp, e.g. 'i'.
export function trRe(key: string, values: Values = {}, flags = ''): RegExp {
  const text = locales.get(locale)?.get(key) ?? english.get(key);
  if (typeof text !== 'object' || 'n' in values) {
    return lineRe(tr(key, values), flags);
  }
  const given = new Map(Object.entries(values));
  const forms = [...new Set(Object.values(text))].map((f) =>
    lineRe(f.replace(/\{(\w+)\}/g, (all, name: string) => String(given.get(name) ?? all))).source,
  );
  return new RegExp(forms.map((f) => `(?:${f})`).join('|'), flags);
}

// lineRe is a line from tr as a pattern, as trRe: for a line built of
// lines, e.g. tr('stack.attack', { line: tr('stack.damage', { n: 2 }) }).
export function lineRe(line: string, flags = ''): RegExp {
  const escaped = line.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  return new RegExp(escaped.replace(/\\\{\w+\\\}/g, '.*?'), flags);
}

// whole is a pattern that must match the whole text.
export function whole(re: RegExp): RegExp {
  return new RegExp(`^(?:${re.source})$`, re.flags);
}

// anyOf matches any of the patterns.
export function anyOf(...res: RegExp[]): RegExp {
  return new RegExp(res.map((r) => `(?:${r.source})`).join('|'));
}
