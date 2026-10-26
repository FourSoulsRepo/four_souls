import { afterEach, describe, expect, it } from 'vitest';
import en from './locales/en.json';
import { language, languages, setLanguage, shownLanguage, t } from './i18n';

afterEach(() => {
  setLanguage('en');
});

describe('t', () => {
  it('reads en.json', () => {
    setLanguage('en');
    expect(t('settings.save')).toBe('Save');
  });

  it('switches at once, with no reload', () => {
    setLanguage('qps');
    expect(language()).toBe('qps');
    expect(t('settings.save')).toBe('[!! Šáṽé ~~ !!]');
    setLanguage('en');
    expect(t('settings.save')).toBe('Save');
  });

  it('falls back to English for Russian, which has no file yet', () => {
    setLanguage('ru');
    expect(language()).toBe('ru');
    expect(t('settings.save')).toBe('Save');
  });
});

describe('languages', () => {
  it('lists English and Russian, and qps in development', () => {
    expect(languages).toEqual(['en', 'ru', 'qps']);
  });

  it('a language this build lacks follows the OS, else English', () => {
    expect(['en', 'ru']).toContain(shownLanguage('xx'));
  });
});

describe('en.json', () => {
  it('has a text and a note for every key, sorted by key', () => {
    const keys = Object.keys(en);
    expect(keys).toEqual([...keys].sort());
    for (const [key, entry] of Object.entries(en)) {
      expect(key).toMatch(/^[a-z_]+\.[a-z_.]+$/);
      expect(entry.note, key).not.toBe('');
      expect(entry.text, key).not.toBe('');
    }
  });
});
