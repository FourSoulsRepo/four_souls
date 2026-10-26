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

// cardLanguages have card names in card_db but no interface file yet;
// they are listed so a player can pick them (16.2, owner: Russian).
export const cardLanguages = ['ru'];

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

// languageList is what the settings offer: English first, then each
// language with an interface file or card names, then the test language
// in development (ADR 021, 11).
export function languageList(files: readonly string[], dev: boolean): string[] {
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
