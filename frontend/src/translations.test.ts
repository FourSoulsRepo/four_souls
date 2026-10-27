// The translation checks on the real files (ADR 021, 10 and 16), run by
// make i18n. Without I18N_LOCALE: every language file has no stray key,
// wrong shape or placeholder; make lint runs it. With I18N_LOCALE=xx:
// merge .cache/translations/xx/done.json into xx.json, write todo.json,
// and fail on every line still to do. The checks are in locale.ts; their
// tests, on test tables, in locale.test.ts.
import { existsSync, mkdirSync, readdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import { areaOf, checkLanguage, gapLines, languageOfFile, mergeDone, testLanguage, todoOf, type Entry, type Glossary } from './locale';

const locales = new URL('./locales/', import.meta.url);
const glossaries = new URL('./locales/glossary/', import.meta.url);
const cache = new URL('../../.cache/translations/', import.meta.url);

const read = (url: URL): unknown => JSON.parse(readFileSync(url, 'utf8'));
const write = (url: URL, data: unknown) => {
  writeFileSync(url, `${JSON.stringify(data, null, 2)}\n`);
};
const table = (url: URL): Record<string, unknown> => (existsSync(url) ? (read(url) as Record<string, unknown>) : {});

const en = read(new URL('en.json', locales)) as Record<string, Entry>;
const enTerms = read(new URL('en.json', glossaries)) as Glossary;
const files = readdirSync(locales).flatMap((name) => {
  const lang = languageOfFile(name);
  return lang ? [lang] : [];
});
const lang = import.meta.env.I18N_LOCALE ?? '';
const area = import.meta.env.I18N_AREA ?? '';

describe.runIf(lang === '')('every language file', () => {
  it.each(files)('%s.json: no stray key, shape or placeholder', (l) => {
    const gaps = checkLanguage(en, l, table(new URL(`${l}.json`, locales)), false);
    expect(gapLines(gaps)).toEqual([]);
  });

  it('every glossary term is in glossary/en.json, each with a note', () => {
    const wrong = readdirSync(glossaries).flatMap((name) => {
      const terms = Object.entries(table(new URL(name, glossaries)));
      const stray = name === 'en.json' ? [] : terms.filter(([t]) => !Object.hasOwn(enTerms, t));
      const empty = terms.filter(([, v]) => typeof v !== 'string' || v === '');
      return [...stray.map(([t]) => `${name}: ${t}: not in glossary/en.json`), ...empty.map(([t]) => `${name}: ${t}: empty`)];
    });
    expect(wrong).toEqual([]);
  });
});

describe.runIf(lang !== '')(`make i18n LOCALE=${lang}`, () => {
  it('has every line translated', () => {
    if (!/^[a-z]{2,3}$/.test(lang) || languageOfFile(`${lang}.json`) === undefined) {
      expect.fail(lang === testLanguage ? 'qps is built from en.json at run time; it has no file' : `${lang}: give a language code other than en, as de or ru`);
    }
    const file = new URL(`${lang}.json`, locales);
    const dir = new URL(`${lang}/`, cache);
    const done = new URL('done.json', dir);
    mkdirSync(dir, { recursive: true });
    let lines = table(file);
    const out: string[] = [];
    if (existsSync(done)) {
      const add = table(done);
      const merged = mergeDone(en, lines, add);
      lines = merged.file;
      write(file, lines);
      rmSync(done);
      out.push(`merged ${String(Object.keys(add).length - merged.unknown.length)} lines of done.json into frontend/src/locales/${lang}.json`);
      out.push(...merged.unknown.map((k) => `done.json: ${k}: en.json has no such key; left out`));
    }
    const all = checkLanguage(en, lang, lines, true);
    const gaps = area === '' ? all : all.filter((g) => areaOf(g.key) === area);
    const terms = table(new URL(`${lang}.json`, glossaries)) as Glossary;
    write(new URL('todo.json', dir), todoOf(en, lang, lines, gaps, enTerms, terms, area));
    const scope = area === '' ? '' : ` in area ${area}`;
    out.push(`${lang}: ${String(gaps.length)} lines to do${scope}, of ${String(Object.keys(en).length)}; todo: .cache/translations/${lang}/todo.json`);
    if (area !== '' && all.length > gaps.length) {
      out.push(`other areas: ${String(all.length - gaps.length)} lines to do`);
    }
    out.push(...gapLines(gaps));
    console.log(out.join('\n'));
    if (gaps.length > 0) {
      expect.fail(`${String(gaps.length)} lines to do in ${lang}${scope}: see the list above`);
    }
  });
});
