// t() and the language files (ADR 021). Every file in locales/ named by
// a language code is a language; adding a language is adding a file.
// The store sets the language from the settings (setLanguage).
import en from './locales/en.json';
import { languageList, localeOf, msgValues, pseudoLocale, resolveLanguage, testLanguage, textOf, translate, type Locale, type Msg, type Text, type Values } from './locale';

// Key is a line of en.json, so a typo fails tsc.
export type Key = keyof typeof en;

const english = textOf(en);
const files = import.meta.glob<Record<string, Text>>('./locales/*.json', { eager: true, import: 'default' });

const locales = new Map<string, Locale>([['en', english]]);
for (const [path, table] of Object.entries(files)) {
  // Other files may sit there (server keys, 16.4); a language is a code.
  const lang = /\/([a-z]{2,3})\.json$/.exec(path)?.[1];
  if (lang && lang !== 'en' && lang !== testLanguage) {
    locales.set(lang, localeOf(table));
  }
}
// The test language exists in development builds only (ADR 021, 12).
if (import.meta.env.DEV) {
  locales.set(testLanguage, pseudoLocale(english));
}

// languages are the choices in Settings, English first.
export const languages = languageList(
  [...locales.keys()].filter((l) => l !== testLanguage),
  import.meta.env.DEV,
);

function osLanguage(): string {
  return typeof navigator === 'undefined' ? 'en' : navigator.language;
}

// shownLanguage is the language a saved setting shows: '' or one this
// build lacks follows the OS, else English.
export function shownLanguage(saved: string): string {
  return resolveLanguage(saved, languages, osLanguage());
}

let current = shownLanguage('');

// setLanguage switches every later t() to the saved setting's language.
export function setLanguage(saved: string): void {
  current = shownLanguage(saved);
}

// language is the language t() uses now.
export function language(): string {
  return current;
}

// say gives a line the server sent (ADR 021, 9): its key in the chosen
// language, with its cards and players named by names. No line is ''.
export function say(m: Msg | undefined, names: { card: (ref: string) => string; player: (p: number) => string }): string {
  if (!m?.key) {
    return '';
  }
  if (import.meta.env.DEV && !english.has(m.key)) {
    console.warn(`no English line for ${m.key}`);
  }
  return translate(locales, current, m.key, msgValues(m, names));
}

// t gives a line in the chosen language, then English, then the key.
// Pass the key as a string literal, so checks can find it.
export function t(key: Key, values?: Values): string {
  if (import.meta.env.DEV && !english.has(key)) {
    console.warn(`no English line for ${key}`);
  }
  return translate(locales, current, key, values);
}
