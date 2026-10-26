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
