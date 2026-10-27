import { describe, expect, it } from 'vitest';
import {
  cardLanguagesOf,
  checkLanguage,
  formExamples,
  gapLines,
  languageList,
  languageName,
  languageOfFile,
  localeOf,
  mergeDone,
  placeholdersOf,
  pseudo,
  pseudoLocale,
  requiredForms,
  resolveLanguage,
  termsOf,
  textOf,
  todoOf,
  translate,
  type Entry,
  type Locale,
} from './locale';

const en = localeOf({
  'lobby.join': 'Join',
  'refuse.attack_cost': 'Attacking costs {cents}¢',
  'hand.count': { one: '{n} card', other: '{n} cards' },
  'only.english': 'Only in English',
});
const ru = localeOf({
  'lobby.join': 'Войти',
  'hand.count': { one: '{n} карта', few: '{n} карты', many: '{n} карт' },
});
const locales = new Map<string, Locale>([
  ['en', en],
  ['ru', ru],
  ['qps', pseudoLocale(en)],
]);

describe('translate', () => {
  it('finds the chosen language, then English, then the key', () => {
    expect(translate(locales, 'ru', 'lobby.join')).toBe('Войти');
    expect(translate(locales, 'ru', 'only.english')).toBe('Only in English');
    expect(translate(locales, 'de', 'lobby.join')).toBe('Join');
    expect(translate(locales, 'ru', 'no.such_key')).toBe('no.such_key');
  });

  it('fills placeholders and keeps unknown ones', () => {
    expect(translate(locales, 'en', 'refuse.attack_cost', { cents: 3 })).toBe('Attacking costs 3¢');
    expect(translate(locales, 'en', 'refuse.attack_cost')).toBe('Attacking costs {cents}¢');
  });

  it('picks the plural form by the language rules', () => {
    expect(translate(locales, 'en', 'hand.count', { n: 1 })).toBe('1 card');
    expect(translate(locales, 'en', 'hand.count', { n: 5 })).toBe('5 cards');
    expect(translate(locales, 'ru', 'hand.count', { n: 1 })).toBe('1 карта');
    expect(translate(locales, 'ru', 'hand.count', { n: 3 })).toBe('3 карты');
    expect(translate(locales, 'ru', 'hand.count', { n: 5 })).toBe('5 карт');
    expect(translate(locales, 'ru', 'hand.count', { n: 21 })).toBe('21 карта');
  });

  it('falls back to other, then to English, for a missing form', () => {
    const partial = new Map([...locales, ['ru', localeOf({ 'hand.count': { other: '{n} шт.' } })]]);
    expect(translate(partial, 'ru', 'hand.count', { n: 5 })).toBe('5 шт.');
    const none = new Map([...locales, ['ru', localeOf({ 'hand.count': {} })]]);
    expect(translate(none, 'ru', 'hand.count', { n: 5 })).toBe('5 cards');
    // A language without the line uses English with English rules.
    expect(translate(locales, 'de', 'hand.count', { n: 1 })).toBe('1 card');
  });

  it('gives the test language English plural rules', () => {
    expect(translate(locales, 'qps', 'hand.count', { n: 1 })).toBe('[!! 1 çáŕð ~~~ !!]');
    expect(translate(locales, 'qps', 'hand.count', { n: 2 })).toBe('[!! 2 çáŕðš ~~~ !!]');
  });
});

describe('pseudo', () => {
  it('accents, lengthens and brackets the text', () => {
    expect(pseudo('Join')).toBe('[!! Ĵôîñ ~~ !!]');
    expect(pseudo('')).toBe('');
  });

  it('keeps placeholders', () => {
    expect(pseudo('Attacking costs {cents}¢')).toBe('[!! Åţţáçķîñĝ çôšţš {cents}¢ ~~~~~~~~ !!]');
  });

  it('is about 30% longer', () => {
    const text = 'Ask before ending the turn while plays are left';
    const out = pseudo(text);
    expect(out.length).toBeGreaterThanOrEqual(text.length * 1.3);
  });

  it('builds every line and plural form', () => {
    const qps = pseudoLocale(en);
    expect(qps.get('lobby.join')).toBe('[!! Ĵôîñ ~~ !!]');
    expect(qps.get('hand.count')).toEqual({ one: '[!! {n} çáŕð ~~~ !!]', other: '[!! {n} çáŕðš ~~~ !!]' });
  });
});

describe('languages', () => {
  it('lists English first, files and card languages, and qps in development', () => {
    expect(languageList(['en'], [], false)).toEqual(['en']);
    expect(languageList(['en'], ['en', 'ru'], false)).toEqual(['en', 'ru']);
    expect(languageList(['en', 'de', 'ru'], ['ru', 'xx'], true)).toEqual(['en', 'de', 'ru', 'xx', 'qps']);
  });

  // Card-only languages come from the names Go sends (16.10).
  it('finds the card languages in the card names', () => {
    const cards: { names?: Record<string, string> }[] = [{ names: { en: 'Isaac', ru: 'Исаак' } }, { names: { en: 'Cain', xx: 'Xcain' } }, {}];
    expect(cardLanguagesOf(cards)).toEqual(['en', 'ru', 'xx']);
    expect(cardLanguagesOf([])).toEqual([]);
  });

  it('resolves the saved language, then the OS one, then English', () => {
    const list = ['en', 'ru', 'qps'];
    expect(resolveLanguage('ru', list, 'en-US')).toBe('ru');
    expect(resolveLanguage('', list, 'ru-RU')).toBe('ru');
    expect(resolveLanguage('', list, 'de-DE')).toBe('en');
    expect(resolveLanguage('de', list, 'ru')).toBe('ru');
    expect(resolveLanguage('qps', ['en', 'ru'], 'en-GB')).toBe('en');
    expect(resolveLanguage('', list, 'qps')).toBe('en');
  });

  it('names each language in its own language', () => {
    expect(languageName('en')).toBe('English');
    expect(languageName('ru')).toBe('Русский');
    expect(languageName('qps')).toBe('qps (test)');
  });
});

describe('textOf', () => {
  it('drops the notes', () => {
    expect(textOf({ 'a.b': { text: 'A', note: 'n' } })).toEqual(new Map([['a.b', 'A']]));
  });
});

// The translation checks (16.7) on test tables: a real trial language
// stays out of locales/, where every file is a choice in Settings.
const source: Record<string, Entry> = {
  'hand.count': { text: { one: '{n} card', other: '{n} cards' }, note: 'Under a hand. n: cards in it.' },
  'hand.title': { text: 'Hand', note: 'Title over the hand.' },
  'lobby.join': { text: 'Join', note: 'Button on a game row.' },
  'lobby.wait': { text: 'Waiting for {player}', note: 'Lobby. player: the host.' },
  'refuse.attack_cost': { text: 'Attacking costs {cents}¢', note: 'Attack button hint.' },
  'stack.loot': { text: '{player} plays {card}', note: 'Stack row.' },
  'table.souls': { text: { one: '{n} soul of {goal}', other: '{n} souls of {goal}' }, note: 'Status. n: souls, goal: souls to win.' },
};

describe('languageOfFile', () => {
  it('is a language code; other files in locales/ are none', () => {
    expect(languageOfFile('de.json')).toBe('de');
    expect(languageOfFile('fil.json')).toBe('fil');
    expect(languageOfFile('en.json')).toBeUndefined();
    expect(languageOfFile('qps.json')).toBeUndefined();
    expect(languageOfFile('server-keys.json')).toBeUndefined();
    expect(languageOfFile('glossary')).toBeUndefined();
  });
});

describe('requiredForms', () => {
  it('are the forms whole numbers 0–1000 pick', () => {
    expect(requiredForms('en')).toEqual(['one', 'other']);
    expect(requiredForms('ru')).toEqual(['one', 'few', 'many']);
    expect(requiredForms('ja')).toEqual(['other']);
    expect(requiredForms('ar')).toEqual(['zero', 'one', 'two', 'few', 'many', 'other']);
    expect(requiredForms('qps')).toEqual(['one', 'other']);
  });

  it('show the numbers each form is for', () => {
    expect(formExamples('ru')).toEqual({ one: '1, 21, 31, 41, 51', few: '2, 3, 4, 22, 23', many: '0, 5, 6, 7, 8' });
  });
});

describe('placeholdersOf', () => {
  it('lists the names of every form, once each', () => {
    expect(placeholdersOf('{player} plays {card}')).toEqual(['card', 'player']);
    expect(placeholdersOf({ one: 'one soul of {goal}', other: '{n} souls of {goal}' })).toEqual(['goal', 'n']);
    expect(placeholdersOf('Join')).toEqual([]);
  });
});

describe('checkLanguage', () => {
  // Half translated: every kind of gap once.
  const half = {
    'hand.count': { one: '{n} карта', few: '{n} карты' },
    'hand.title': { text: 'Рука', note: 'copied from en.json' },
    'lobby.join': 'Войти',
    'lobby.wait': 'Ждём {host}',
    'refuse.attack_cost': 'Атака стоит ¢',
    'stack.loot': '',
    'table.souls': 'душ: {n}',
    'old.key': 'Старое',
  };

  it('lists every gap of a half-translated file', () => {
    expect(checkLanguage(source, 'ru', half, true).map((g) => [g.key, g.problem])).toEqual([
      ['hand.count', 'forms'],
      ['hand.title', 'shape'],
      ['lobby.wait', 'placeholders'],
      ['old.key', 'stray'],
      ['refuse.attack_cost', 'placeholders'],
      ['stack.loot', 'missing'],
      ['table.souls', 'shape'],
    ]);
  });

  it('says each gap in plain words', () => {
    const lines = gapLines(checkLanguage(source, 'ru', half, true));
    expect(lines).toContain('forms        hand.count          needs one, few, many; lacks many');
    expect(lines).toContain('placeholders lobby.wait          has {host}; English has {player}');
    expect(lines).toContain('placeholders refuse.attack_cost  has none; English has {cents}');
    expect(lines).toContain('shape        hand.title          holds text and note; a language file holds the text only');
    expect(lines).toContain('shape        table.souls         a plural line: give one text per form: one, few, many');
  });

  it('lists keys not there yet as missing, with the English', () => {
    expect(checkLanguage(source, 'ru', {}, true)).toContainEqual({ key: 'lobby.join', problem: 'missing', detail: 'not translated: Join' });
  });

  it('without LOCALE: only stray keys, shapes and placeholders', () => {
    expect(checkLanguage(source, 'ru', half, false).map((g) => g.problem)).toEqual(['shape', 'placeholders', 'stray', 'placeholders', 'shape']);
  });

  it('passes a full file, where one form leaves out {n}', () => {
    const full = {
      'hand.count': { one: 'одна карта', few: '{n} карты', many: '{n} карт', other: '{n} карты' },
      'hand.title': 'Рука',
      'lobby.join': 'Войти',
      'lobby.wait': 'Ждём {player}',
      'refuse.attack_cost': 'Атака стоит {cents}¢',
      'stack.loot': '{player} играет {card}',
      'table.souls': { one: '{n} душа из {goal}', few: '{n} души из {goal}', many: '{n} душ из {goal}' },
    };
    expect(checkLanguage(source, 'ru', full, true)).toEqual([]);
  });

  it('fails a form the language does not have, or a name English lacks', () => {
    const wrong = { 'hand.count': { one: '{n} карта', two: '{n} карты', few: '{n} карты', many: '{n} карт' }, 'hand.title': 'Рука {n}' };
    expect(gapLines(checkLanguage(source, 'ru', wrong, false))).toEqual([
      'shape        hand.count  ru has no plural form two; it uses one, few, many, other',
      'placeholders hand.title  has {n}; English has none',
    ]);
  });
});

describe('termsOf', () => {
  const terms = ['loot', 'soul', 'roll', 'cent', 'bonus soul'];

  it('finds a term in any case, plural or not', () => {
    expect(termsOf('Loot 2 cards', terms)).toEqual(['loot']);
    expect(termsOf({ one: '{n} soul', other: '{n} souls' }, terms)).toEqual(['soul']);
    expect(termsOf('Take the bonus soul', terms)).toEqual(['soul', 'bonus soul']);
  });

  it('skips placeholders and longer words', () => {
    expect(termsOf('Costs {cents}¢', terms)).toEqual([]);
    expect(termsOf('Rollback, unrolled', terms)).toEqual([]);
  });
});

describe('todoOf', () => {
  const half = { 'lobby.join': 'Войти', 'lobby.wait': 'Ждём {host}' };
  const gaps = checkLanguage(source, 'ru', half, true);
  const enTerms = { soul: 'What players collect.', loot: 'Hand cards.' };

  it('groups the lines by area, with all an agent needs', () => {
    const todo = todoOf(source, 'ru', half, gaps, enTerms, { soul: 'душа' }, '');
    expect(Object.keys(todo.areas)).toEqual(['hand', 'lobby', 'refuse', 'stack', 'table']);
    expect(todo.plural_forms).toEqual(formExamples('ru'));
    expect(todo.glossary).toEqual({ soul: 'What players collect.' });
    expect(todo.areas.table).toEqual([
      {
        key: 'table.souls',
        english: { one: '{n} soul of {goal}', other: '{n} souls of {goal}' },
        note: 'Status. n: souls, goal: souls to win.',
        placeholders: ['goal', 'n'],
        forms: ['one', 'few', 'many'],
        glossary: { soul: 'душа' },
        problem: 'missing',
        detail: 'not translated: {"one":"{n} soul of {goal}","other":"{n} souls of {goal}"}',
      },
    ]);
    expect(todo.areas.lobby[0]).toMatchObject({ key: 'lobby.wait', problem: 'placeholders', now: 'Ждём {host}' });
    expect(todo.stray).toEqual([]);
  });

  it('lists stray keys apart', () => {
    const gone = { 'old.key': 'Старое' };
    expect(todoOf(source, 'ru', gone, checkLanguage(source, 'ru', gone, true), enTerms, {}, '').stray).toEqual(['old.key']);
  });

  it('holds one area with AREA', () => {
    expect(Object.keys(todoOf(source, 'ru', half, gaps, enTerms, {}, 'hand').areas)).toEqual(['hand']);
  });
});

describe('mergeDone', () => {
  it('adds the done lines, sorted by key, deletes nulls, and leaves out unknown keys', () => {
    const file = { 'lobby.wait': 'Ждём {host}', 'hand.title': 'Рука', 'old.key': 'Старое' };
    const merged = mergeDone(source, file, { 'lobby.wait': 'Ждём {player}', 'lobby.join': 'Войти', 'no.such': 'x', 'old.key': null });
    expect(Object.entries(merged.file)).toEqual([
      ['hand.title', 'Рука'],
      ['lobby.join', 'Войти'],
      ['lobby.wait', 'Ждём {player}'],
    ]);
    expect(merged.unknown).toEqual(['no.such']);
  });
});
