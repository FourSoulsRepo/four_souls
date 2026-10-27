// The language the tests run in: LOCALE (make screenshots LOCALE=qps),
// English without it. tr gives a line as that language shows it, so a
// test finds what the player reads (ADR 021); cardName does the same for
// a card's name (16.5).
import { existsSync, readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { cardLanguagesOf, languageList, languageName, languageOfFile, localeOf, pseudoLocale, requiredForms, textOf, translate, type Entry, type Locale, type Text, type Values } from '../../frontend/src/locale';

const root = join(import.meta.dirname, '..', '..');
const dir = join(root, 'frontend', 'src', 'locales');
const english = textOf(JSON.parse(readFileSync(join(dir, 'en.json'), 'utf8')) as Record<string, Entry>);

// locale is the language every app page starts in.
export const locale = process.env.LOCALE || 'en';

const locales = new Map<string, Locale>([
  ['en', english],
  ['qps', pseudoLocale(english)],
]);
// A language with an interface file shows its lines (16.7): the tests
// read the same file the app does; a line it lacks stays English.
const own = join(dir, `${locale}.json`);
if (languageOfFile(`${locale}.json`) && existsSync(own)) {
  locales.set(locale, localeOf(JSON.parse(readFileSync(own, 'utf8')) as Record<string, Text>));
}

interface DataCard {
  id: string;
  version?: number;
  reprint?: boolean;
  text: Record<string, { name?: string } | undefined>;
  back?: DataCard; // a two-sided card's other face (ADR 023)
}

// The fake bridge's cards are the Base Game's and Requiem's
// (fakeBridge.ts); a later version's ref is id@N, as Go's, and a back
// face is a card of its own.
function dataCards(set: string): DataCard[] {
  return (JSON.parse(readFileSync(join(root, 'pkg', 'card_db', 'data', `${set}.json`), 'utf8')) as { cards: DataCard[] }).cards;
}
const b2 = { cards: dataCards('b2') };
const all = [...b2.cards, ...dataCards('r')].filter((c) => !c.reprint).flatMap((c) => (c.back ? [c, c.back] : [c]));
const cards = new Map(all.map((c) => [c.version && c.version > 1 ? `${c.id}@${String(c.version)}` : c.id, c]));

// languageNames are the Settings choices of a development build, as
// i18n.ts lists them: English, each language file and each language the
// fake bridge's cards have names in (16.10), qps. extra adds card
// languages a test injects (?cardlang).
export function languageNames(...extra: string[]): string[] {
  const names = b2.cards.map((c) => ({ names: Object.fromEntries(Object.entries(c.text).flatMap(([l, t]) => (t?.name ? [[l, t.name]] : []))) }));
  const files = readdirSync(dir).flatMap((name) => languageOfFile(name) ?? []);
  return languageList(files, [...cardLanguagesOf(names), ...extra], true).map(languageName);
}

// cardName is the card ref's name in the tests' language: card_db's
// text.<LOCALE>.name, else English, as the store's inLanguage picks it.
// With LOCALE=ru the names are Russian on English screens (16.5).
export function cardName(ref: string): string {
  const card = cards.get(ref);
  const name = card?.text[locale]?.name ?? card?.text.en?.name;
  if (name === undefined) {
    throw new Error(`no card ${ref} in b2.json or r.json`);
  }
  return name;
}

// tr is the line key in the tests' language.
export function tr(key: string, values?: Values): string {
  if (!english.has(key)) {
    throw new Error(`no line ${key} in en.json`);
  }
  return translate(locales, locale, key, values);
}

// trIn is the line key in another language: a page switched to it
// (languages.spec.ts, 16.4).
export function trIn(lang: string, key: string, values?: Values): string {
  if (!english.has(key)) {
    throw new Error(`no line ${key} in en.json`);
  }
  return translate(locales, lang, key, values);
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
  // A file lacking a form the language needs shows English for it.
  const en = english.get(key);
  const lacks = requiredForms(locale).some((f) => !(f in text));
  const texts = [...Object.values(text), ...(lacks && typeof en === 'object' ? Object.values(en) : [])];
  const forms = [...new Set(texts)].map((f) =>
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
