import { describe, expect, it } from 'vitest';
import { languageList, languageName, localeOf, pseudo, pseudoLocale, resolveLanguage, textOf, translate, type Locale } from './locale';

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
    expect(languageList(['en'], false)).toEqual(['en', 'ru']);
    expect(languageList(['en', 'de', 'ru'], true)).toEqual(['en', 'de', 'ru', 'qps']);
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
