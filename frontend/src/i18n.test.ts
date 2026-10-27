import { afterEach, describe, expect, it } from 'vitest';
import en from './locales/en.json';
import { language, languages, setCardLanguages, setLanguage, shownLanguage, t } from './i18n';

afterEach(() => {
  setCardLanguages([]);
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
    setCardLanguages(['en', 'ru']);
    setLanguage('ru');
    expect(language()).toBe('ru');
    expect(t('settings.save')).toBe('Save');
  });
});

describe('languages', () => {
  // A language file adds its own choice, so a new one needs no change
  // here; the glossary folder and server-keys.json add none (16.7).
  it('lists English, the files, and qps in development', () => {
    setCardLanguages([]);
    expect(languages()[0]).toBe('en');
    expect(languages()).not.toContain('xx');
    expect(languages().at(-1)).toBe('qps');
    expect(languages()).not.toContain('glossary');
    expect(languages()).not.toContain('server-keys');
  });

  // A card-only language joins when the cards load (16.10); a setting
  // saved before that shows it then.
  it('adds the card languages; a saved one shows once they load', () => {
    setCardLanguages([]);
    setLanguage('xx');
    expect(language()).not.toBe('xx');
    expect(setCardLanguages(['en', 'ru', 'xx'])).toEqual(expect.arrayContaining(['en', 'ru', 'xx', 'qps']));
    expect(language()).toBe('xx');
    expect(t('settings.save')).toBe('Save');
  });

  it('a language this build lacks follows the OS, else English', () => {
    setCardLanguages(['ru']);
    expect(['en', 'ru']).toContain(shownLanguage('xx'));
  });
});

describe('en.json', () => {
  it('has a text and a note for every key, sorted by key', () => {
    const keys = Object.keys(en);
    expect(keys).toEqual([...keys].sort());
    for (const [key, entry] of Object.entries(en)) {
      expect(key).toMatch(/^[a-z0-9_]+\.[a-z0-9_.]+$/);
      expect(entry.note, key).not.toBe('');
      expect(entry.text, key).not.toBe('');
    }
  });
});

// Every key in en.json is used, and every key used is there (ADR 021,
// 10). Frontend code passes keys to t() as string literals (8.4); the
// server's keys are listed in locales/server-keys.json (16.4).
describe('en.json keys', () => {
  const sources = import.meta.glob<string>(['./**/*.{ts,tsx}', '!./**/*.test.ts', '!./i18n.ts'], {
    query: '?raw',
    import: 'default',
    eager: true,
  });
  const serverFiles = import.meta.glob<unknown>('./locales/server-keys.json', { import: 'default', eager: true });
  const server = Object.values(serverFiles).flatMap((f) => (Array.isArray(f) ? f.map(String) : Object.keys(f ?? {})));
  const used = new Set<string>();
  const computed: string[] = [];
  for (const [file, text] of Object.entries(sources)) {
    for (const m of text.matchAll(/\bt\(\s*([^)\s,]*)/g)) {
      // "t()" with nothing in it is a comment that names t.
      if (m[1] === '') continue;
      const literal = /^'([^']+)'$/.exec(m[1]);
      if (literal?.[1]) {
        used.add(literal[1]);
      } else {
        computed.push(`${file}: t(${m[1]}`);
      }
    }
  }

  it('passes every key as a string literal', () => {
    expect(computed).toEqual([]);
  });

  it('has every key the code uses', () => {
    expect([...used].filter((k) => !Object.hasOwn(en, k))).toEqual([]);
  });

  it('uses every key it has', () => {
    const all = new Set([...used, ...server]);
    expect(Object.keys(en).filter((k) => !all.has(k))).toEqual([]);
  });
});
