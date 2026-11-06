import { describe, expect, it } from 'vitest';
import { checkHash, parseRoute, routeHash, rulesHash } from './route';

describe('routes (ADR 026, 6)', () => {
  it('reads a check link and a rule link', () => {
    expect(parseRoute('#/en/check/abc_-1', 'en')).toEqual({ route: { lang: 'en', view: 'check', data: 'abc_-1' }, canonical: true });
    expect(parseRoute('#/ru/rules/R-CARD-08', 'en')).toEqual({ route: { lang: 'ru', view: 'rules', rule: 'R-CARD-08' }, canonical: true });
  });

  it('without a language uses the fallback and asks for the full form', () => {
    const got = parseRoute('#/rules/R-CARD-08', 'en');
    expect(got).toEqual({ route: { lang: 'en', view: 'rules', rule: 'R-CARD-08' }, canonical: false });
    expect(routeHash(got.route)).toBe('#/en/rules/R-CARD-08');
    expect(parseRoute('#/check/xyz', 'de').route).toEqual({ lang: 'de', view: 'check', data: 'xyz' });
  });

  it('an empty or unknown hash opens the check view', () => {
    for (const hash of ['', '#', '#/', '#/en', '#/en/nothing']) {
      const got = parseRoute(hash, 'en');
      expect(got.route, hash).toEqual({ lang: 'en', view: 'check', data: '' });
      expect(got.canonical, hash).toBe(false);
    }
    expect(parseRoute('#/en/check', 'en').canonical).toBe(true);
  });

  it('a view name is no language', () => {
    expect(parseRoute('#/faq', 'en').route.lang).toBe('en');
  });

  it('builds the hashes it reads', () => {
    expect(rulesHash('en', 'R-STACK-03')).toBe('#/en/rules/R-STACK-03');
    expect(rulesHash('en', '')).toBe('#/en/rules');
    expect(checkHash('en', '')).toBe('#/en/check');
    expect(parseRoute(checkHash('en', 'q1-_'), 'ru').route).toEqual({ lang: 'en', view: 'check', data: 'q1-_' });
  });
});
