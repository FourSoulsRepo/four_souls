// t() and the language files (ADR 021). Every file in locales/ named by
// a language code is a language; adding a language is adding a file.
// The store sets the language from the settings (setLanguage).
import en from './locales/en.json';
import { languageList, languageOfFile, localeOf, msgValues, pseudoLocale, resolveLanguage, testLanguage, textOf, translate, type Locale, type Msg, type Text, type Values } from './locale';

// Key is a line of en.json, so a typo fails tsc.
export type Key = keyof typeof en;

const english = textOf(en);
const files = import.meta.glob<Record<string, Text>>('./locales/*.json', { eager: true, import: 'default' });

const locales = new Map<string, Locale>([['en', english]]);
for (const [path, table] of Object.entries(files)) {
  // Other files sit there too: server keys (16.4); the glossary is a
  // folder, which the pattern skips (16.7). A language is a code.
  const lang = languageOfFile(path.split('/').pop() ?? '');
  if (lang) {
    locales.set(lang, localeOf(table));
  }
}
// The test language exists in development builds only (ADR 021, 12).
if (import.meta.env.DEV) {
  locales.set(testLanguage, pseudoLocale(english));
}

const fileLanguages = [...locales.keys()].filter((l) => l !== testLanguage);
let list = languageList(fileLanguages, [], import.meta.env.DEV);

// languages are the choices in Settings, English first: the language
// files, then the card languages once the cards load (16.10).
export function languages(): readonly string[] {
  return list;
}

function osLanguage(): string {
  return typeof navigator === 'undefined' ? 'en' : navigator.language;
}

// shownLanguage is the language a saved setting shows: '' or one this
// build lacks follows the OS, else English.
export function shownLanguage(saved: string): string {
  return resolveLanguage(saved, list, osLanguage());
}

let saved = '';
let current = shownLanguage(saved);

// setLanguage switches every later t() to the saved setting's language.
export function setLanguage(setting: string): void {
  saved = setting;
  current = shownLanguage(saved);
}

// setCardLanguages adds the languages the cards have names in; the
// saved setting may name one, so it is resolved again. It gives the new
// list for Settings.
export function setCardLanguages(langs: readonly string[]): readonly string[] {
  list = languageList(fileLanguages, langs, import.meta.env.DEV);
  current = shownLanguage(saved);
  return list;
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
