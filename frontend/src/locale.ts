// Pure parts of the translations (ADR 021): lookup, plurals, values,
// the test language and the language list. i18n.ts loads the files and
// holds the chosen language; the e2e tests import this file too.

// Forms is a plural line: one text per CLDR form (ADR 021, 7).
export type Forms = Partial<Record<Intl.LDMLPluralRule, string>>;

// Text is one line: plain, or one text per plural form.
export type Text = string | Forms;

// Locale is one language: key → text. en.json also has notes; see Entry.
export type Locale = ReadonlyMap<string, Text>;

// Entry is a line of en.json: the text and where it shows.
export interface Entry {
  text: Text;
  note: string;
}

// Values fill a line's {name} placeholders; n also picks the plural form.
export type Values = Record<string, string | number>;

// MsgValue is one value of a server message (ADR 021, 6): a number, a
// plain string, a card ({card: ref}) or a player ({player: seat}).
export type MsgValue = number | string | { card: string } | { player: number };

// Msg is a line the server sends: a key of en.json and its values.
export interface Msg {
  key?: string;
  args?: Record<string, MsgValue>;
}

// msgValues turns a message's values into plain ones: a card by its
// name, a player by theirs.
export function msgValues(m: Msg, names: { card: (ref: string) => string; player: (p: number) => string }): Values {
  return Object.fromEntries(
    Object.entries(m.args ?? {}).map(([name, v]) => [name, typeof v !== 'object' ? v : 'card' in v ? names.card(v.card) : names.player(v.player)]),
  );
}

// testLanguage is built from English at run time; development only.
export const testLanguage = 'qps';

// textOf drops the notes of en.json.
export function textOf(entries: Record<string, Entry>): Locale {
  return new Map(Object.entries(entries).map(([k, e]) => [k, e.text]));
}

// localeOf reads another language's file.
export function localeOf(file: Record<string, Text>): Locale {
  return new Map(Object.entries(file));
}

// translate finds key in lang, then in English, then gives the key
// itself (ADR 021, 8).
export function translate(locales: ReadonlyMap<string, Locale>, lang: string, key: string, values: Values = {}): string {
  const own = locales.get(lang)?.get(key);
  const en = locales.get('en')?.get(key);
  const text = own ?? en;
  const vals = new Map(Object.entries(values));
  if (text === undefined) {
    return key;
  }
  if (typeof text === 'string') {
    return fill(text, vals);
  }
  // A plural line: the language's form, its other, then English.
  const n = Number(vals.get('n') ?? 0);
  const forms = new Map(Object.entries(text));
  const form = new Intl.PluralRules(pluralLang(own === undefined ? 'en' : lang)).select(n);
  const enForms = new Map(Object.entries(typeof en === 'object' ? en : {}));
  const enForm = new Intl.PluralRules('en').select(n);
  return fill(forms.get(form) ?? forms.get('other') ?? enForms.get(enForm) ?? enForms.get('other') ?? key, vals);
}

// pluralLang is the language whose plural rules a language uses.
function pluralLang(lang: string): string {
  return lang === testLanguage ? 'en' : lang;
}

// fill puts the values into {name}; an unknown name stays as it is.
function fill(text: string, values: ReadonlyMap<string, string | number>): string {
  return text.replace(/\{(\w+)\}/g, (all, name: string) => {
    const v = values.get(name);
    return v === undefined ? all : String(v);
  });
}

const plain = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
const accents = new Map(Array.from('áƀçðéƒĝĥîĵķļɱñôþǫŕšţûṽŵẋýžÅƁÇÐÉƑĜĤÎĴĶĻṀÑÖÞǪŔŠŢÛṼŴẊÝŽ', (a, i) => [plain.charAt(i), a]));

// pseudo turns English into the test language (ADR 021, 12):
// "Join" becomes "[!! Ĵôîñ ~~ !!]", accented and about 30% longer.
// Placeholders stay as they are.
export function pseudo(text: string): string {
  if (text === '') {
    return '';
  }
  const body = text
    .split(/(\{\w+\})/)
    .map((part) => (part.startsWith('{') ? part : Array.from(part, (c) => accents.get(c) ?? c).join('')))
    .join('');
  return `[!! ${body} ${'~'.repeat(Math.ceil(text.length * 0.3))} !!]`;
}

// pseudoLocale is the whole test language, from English.
export function pseudoLocale(en: Locale): Locale {
  return new Map(
    [...en].map(([k, text]) => [k, typeof text === 'string' ? pseudo(text) : Object.fromEntries(Object.entries(text).map(([form, f]) => [form, pseudo(f)]))]),
  );
}

// cardLanguagesOf lists the languages the cards have names in, from the
// names Go sends (App.Cards, 16.10): a language with no interface file
// is still a choice, with English screens and its card names.
export function cardLanguagesOf(cards: readonly { names?: Record<string, string> }[]): string[] {
  return [...new Set(cards.flatMap((c) => Object.keys(c.names ?? {})))].sort();
}

// languageList is what the settings offer: English first, then each
// language with an interface file or card names, then the test language
// in development (ADR 021, 11).
export function languageList(files: readonly string[], cardLanguages: readonly string[], dev: boolean): string[] {
  const others = [...new Set([...files, ...cardLanguages])].filter((l) => l !== 'en' && l !== testLanguage).sort();
  return ['en', ...others, ...(dev ? [testLanguage] : [])];
}

// resolveLanguage is the language to show: the saved one if listed,
// else the OS language if listed, else English.
export function resolveLanguage(saved: string, list: readonly string[], os: string): string {
  if (saved !== '' && list.includes(saved)) {
    return saved;
  }
  const primary = os.toLowerCase().split('-')[0] ?? '';
  return primary !== testLanguage && list.includes(primary) ? primary : 'en';
}

// languageName is a language in its own language: "English", "Русский".
export function languageName(code: string): string {
  if (code === testLanguage) {
    return 'qps (test)';
  }
  try {
    const name = new Intl.DisplayNames([code], { type: 'language' }).of(code) ?? code;
    return name.charAt(0).toLocaleUpperCase(code) + name.slice(1);
  } catch {
    return code;
  }
}

// Checks for translators (ADR 021, 10 and 16). translations.test.ts
// runs them on the files: make i18n, with or without LOCALE.

// pluralForms are the CLDR forms in their usual order.
const pluralForms: readonly Intl.LDMLPluralRule[] = ['zero', 'one', 'two', 'few', 'many', 'other'];

// languageOfFile is the language a file in locales/ holds: a name that
// is a language code (de.json). Other files there are no language:
// server-keys.json, and glossary/ is a folder. en.json is the source.
export function languageOfFile(name: string): string | undefined {
  const lang = /^([a-z]{2,3})\.json$/.exec(name)?.[1];
  return lang === undefined || lang === 'en' || lang === testLanguage ? undefined : lang;
}

// areaOf is a key's area: its first word (ADR 021, 5).
export function areaOf(key: string): string {
  return key.split('.')[0] ?? key;
}

// requiredForms are the plural forms a language must fill: those its
// rules pick for whole numbers 0–1000 (ADR 021, 10). Russian: one, few,
// many; its other is for fractions only.
export function requiredForms(lang: string): Intl.LDMLPluralRule[] {
  return Object.keys(formExamples(lang)) as Intl.LDMLPluralRule[];
}

// formExamples gives each required form the first numbers that pick
// it, so a translator sees what the form is for: one: "1, 21, 31".
export function formExamples(lang: string): Partial<Record<Intl.LDMLPluralRule, string>> {
  const rules = new Intl.PluralRules(pluralLang(lang));
  const found = new Map<Intl.LDMLPluralRule, number[]>();
  for (let n = 0; n <= 1000; n++) {
    const form = rules.select(n);
    found.set(form, [...(found.get(form) ?? []), n]);
  }
  return Object.fromEntries(pluralForms.filter((f) => found.has(f)).map((f) => [f, (found.get(f) ?? []).slice(0, 5).join(', ')]));
}

// placeholdersOf lists a line's {names}, of all its forms together.
export function placeholdersOf(text: Text): string[] {
  const all = typeof text === 'string' ? [text] : Object.values(text);
  return [...new Set(all.flatMap((t) => Array.from(t.matchAll(/\{(\w+)\}/g), (m) => m[1])))].sort();
}

// Problem is what is wrong with one line of a language file.
export type Problem = 'missing' | 'stray' | 'shape' | 'placeholders' | 'forms';

// Gap is one problem of one line; detail says it in plain words.
export interface Gap {
  key: string;
  problem: Problem;
  detail: string;
}

// lineGap is what is wrong with one translated line, or nothing; all
// also asks for every required plural form (make i18n LOCALE=xx).
// value is as read from JSON, so any shape may come.
export function lineGap(english: Text, value: unknown, lang: string, all: boolean): Omit<Gap, 'key'> | undefined {
  const needed = requiredForms(lang);
  if (typeof english === 'string') {
    if (typeof value !== 'string') {
      return { problem: 'shape', detail: holdsNote(value) ?? 'English is one text, so this must be one text too' };
    }
    if (value === '') {
      return { problem: 'missing', detail: 'empty text' };
    }
  } else {
    if (typeof value !== 'object' || value === null || Array.isArray(value)) {
      return { problem: 'shape', detail: `a plural line: give one text per form: ${needed.join(', ')}` };
    }
    const forms = Object.entries(value as Record<string, unknown>);
    const used = new Set(new Intl.PluralRules(pluralLang(lang)).resolvedOptions().pluralCategories);
    const bad = forms.filter(([f, t]) => !used.has(f as Intl.LDMLPluralRule) || typeof t !== 'string');
    if (bad.length > 0) {
      return { problem: 'shape', detail: holdsNote(value) ?? `${lang} has no plural form ${bad.map(([f]) => f).join(', ')}; it uses ${pluralForms.filter((f) => used.has(f)).join(', ')}` };
    }
    const have = new Set(forms.filter(([, t]) => t !== '').map(([f]) => f));
    const lacking = needed.filter((f) => !have.has(f));
    if (lacking.length === needed.length) {
      return { problem: 'missing', detail: 'no plural form filled' };
    }
    if (all && lacking.length > 0) {
      return { problem: 'forms', detail: `needs ${needed.join(', ')}; lacks ${lacking.join(', ')}` };
    }
  }
  const want = placeholdersOf(english);
  const got = placeholdersOf(value);
  const extra = got.filter((p) => !want.includes(p));
  const lost = want.filter((p) => !got.includes(p));
  if (extra.length > 0 || lost.length > 0) {
    const list = (ps: string[]) => (ps.length === 0 ? 'none' : ps.map((p) => `{${p}}`).join(' '));
    return { problem: 'placeholders', detail: `has ${list(got)}; English has ${list(want)}` };
  }
  return undefined;
}

// holdsNote catches an entry copied from en.json with its note.
function holdsNote(value: unknown): string | undefined {
  return typeof value === 'object' && value !== null && 'note' in value ? 'holds text and note; a language file holds the text only' : undefined;
}

// checkLanguage lists the gaps of one language file, sorted by key.
// Every file: stray keys, wrong shapes and placeholders. all
// (make i18n LOCALE=xx) also lists missing keys and plural forms.
export function checkLanguage(en: Record<string, Entry>, lang: string, file: Record<string, unknown>, all: boolean): Gap[] {
  const gaps: Gap[] = [];
  const lines = new Map(Object.entries(file));
  for (const [key, entry] of Object.entries(en)) {
    if (!lines.has(key)) {
      if (all) {
        gaps.push({ key, problem: 'missing', detail: `not translated: ${typeof entry.text === 'string' ? entry.text : JSON.stringify(entry.text)}` });
      }
      continue;
    }
    const gap = lineGap(entry.text, lines.get(key), lang, all);
    if (gap && (all || gap.problem !== 'missing')) {
      gaps.push({ key, ...gap });
    }
  }
  for (const key of Object.keys(file)) {
    if (!Object.hasOwn(en, key)) {
      gaps.push({ key, problem: 'stray', detail: 'en.json has no such key; delete it' });
    }
  }
  return gaps.sort((a, b) => (a.key < b.key ? -1 : a.key > b.key ? 1 : 0));
}

// Glossary is the game's terms (ADR 021, 16): in glossary/en.json each
// English term with a note; in glossary/<lang>.json its translation.
export type Glossary = Record<string, string>;

// termsOf lists the glossary terms an English line uses, in any case,
// plural or not; placeholders are not words.
export function termsOf(text: Text, terms: readonly string[]): string[] {
  const words = (typeof text === 'string' ? [text] : Object.values(text)).join(' ').replace(/\{\w+\}/g, ' ').toLowerCase();
  return terms.filter((term) => hasWord(words, term.toLowerCase()));
}

// hasWord finds word in text as a whole word, or with s or es after it.
function hasWord(text: string, word: string): boolean {
  const letter = (c: string) => c >= 'a' && c <= 'z';
  for (let at = text.indexOf(word); at >= 0; at = text.indexOf(word, at + 1)) {
    const rest = text.slice(at + word.length);
    const ends = ['', 's', 'es'].some((end) => rest.startsWith(end) && !letter(rest.charAt(end.length)));
    if (!letter(text.charAt(at - 1)) && ends) {
      return true;
    }
  }
  return false;
}

// TodoEntry is one line for a translator or an agent to do.
export interface TodoEntry {
  key: string;
  english: Text;
  note: string;
  placeholders: string[];
  forms?: Intl.LDMLPluralRule[];
  glossary?: Glossary;
  problem: Problem;
  detail: string;
  now?: unknown;
}

// Todo is .cache/translations/<lang>/todo.json, grouped by area.
export interface Todo {
  language: string;
  how: string[];
  plural_forms: Partial<Record<Intl.LDMLPluralRule, string>>;
  glossary: Glossary;
  areas: Record<string, TodoEntry[]>;
  stray: string[];
}

// todoOf turns the gaps into a todo, of one area or of all (''). Each entry
// has what is needed to translate it; glossary terms are context only,
// with no check of the words (owner, 2026-10-27). enTerms: term → note;
// terms: term → this language's word.
export function todoOf(en: Record<string, Entry>, lang: string, file: Record<string, unknown>, gaps: readonly Gap[], enTerms: Glossary, terms: Glossary, area: string): Todo {
  const areas = new Map<string, TodoEntry[]>();
  const usedTerms = new Set<string>();
  const entries = new Map(Object.entries(en));
  const lines = new Map(Object.entries(file));
  const words = new Map(Object.entries(terms));
  const stray: string[] = [];
  for (const gap of gaps) {
    if (gap.problem === 'stray') {
      stray.push(gap.key);
      continue;
    }
    const entry = entries.get(gap.key);
    const at = areaOf(gap.key);
    if (entry === undefined || (area !== '' && at !== area)) {
      continue;
    }
    const found = termsOf(entry.text, Object.keys(enTerms));
    found.forEach((t) => usedTerms.add(t));
    const item: TodoEntry = {
      key: gap.key,
      english: entry.text,
      note: entry.note,
      placeholders: placeholdersOf(entry.text),
      ...(typeof entry.text === 'object' ? { forms: requiredForms(lang) } : {}),
      ...(found.length > 0 ? { glossary: Object.fromEntries(found.map((t) => [t, words.get(t) ?? ''])) } : {}),
      problem: gap.problem,
      detail: gap.detail,
      ...(lines.has(gap.key) ? { now: lines.get(gap.key) } : {}),
    };
    areas.set(at, [...(areas.get(at) ?? []), item]);
  }
  const done = `.cache/translations/${lang}/done.json`;
  return {
    language: lang,
    how: [
      `Translate each entry's english into the language ${lang}; its note says where the line shows and what each {value} is.`,
      `Write only ${done}: one object, key → text, e.g. {"lobby.join": "…"}.`,
      'A line with forms is a plural line: give an object with exactly those forms, e.g. {"one": "…", "few": "…"}; plural_forms shows which numbers pick each form.',
      'Keep every {placeholder} exactly as written; together the forms use every one. One plural form may leave out {n}.',
      'Names of players and cards stay in their dictionary form: "{player} → {monster}: attack" (ADR 021, 15).',
      'glossary: the game terms a line uses and this language\'s word for each; an empty word is not fixed yet: pick one and keep it the same everywhere.',
      'An entry with now has a translation that is wrong: detail says why.',
      'stray: keys en.json no longer has; give each null in done.json to delete it.',
      `Then run make i18n LOCALE=${lang} (AREA=<area> for one area): it merges done.json into frontend/src/locales/${lang}.json and lists what is left.`,
      'make screenshots shows every screen in English into .cache/screenshots/; LOCALE=' + lang + ' shows the translation.',
    ],
    plural_forms: formExamples(lang),
    glossary: Object.fromEntries(Object.entries(enTerms).filter(([t]) => usedTerms.has(t))),
    areas: Object.fromEntries(areas),
    stray,
  };
}

// mergeDone puts an agent's done.json into a language file, sorted by
// key (ADR 021, 16). A key en.json lacks stays out and is listed; a key
// set to null is deleted, as a stray one.
export function mergeDone(en: Record<string, Entry>, file: Record<string, unknown>, done: Record<string, unknown>): { file: Record<string, unknown>; unknown: string[] } {
  const merged = new Map(Object.entries(file));
  const unknown: string[] = [];
  for (const [key, text] of Object.entries(done)) {
    if (text === null) {
      merged.delete(key);
    } else if (Object.hasOwn(en, key)) {
      merged.set(key, text);
    } else {
      unknown.push(key);
    }
  }
  return { file: Object.fromEntries([...merged].sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0))), unknown };
}

// gapLines prints gaps as plain lines for people.
export function gapLines(gaps: readonly Gap[]): string[] {
  const width = Math.max(0, ...gaps.map((g) => g.key.length));
  return gaps.map((g) => `${g.problem.padEnd(12)} ${g.key.padEnd(width)}  ${g.detail}`);
}
