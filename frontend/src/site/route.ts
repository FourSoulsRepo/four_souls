// Routes of the rules website, in the URL fragment (ADR 026, 6), so
// GitHub Pages needs no rewrites and a payload never reaches a server:
//   #/<lang>/check/<data>    a situation and its answer
//   #/<lang>/rules/<rule ID> the rules text at that rule
// Without <lang> the page uses the browser's language (fallback).

export type Route = { lang: string; view: 'check'; data: string } | { lang: string; view: 'rules'; rule: string };

// A language is a code of two or three letters; faq is a view (12.6).
const langCode = /^[a-z]{2,3}$/;
const views = new Set(['check', 'rules', 'faq', 'stack']);

// parseRoute reads a location hash. fallback is the language when the
// hash names none. canonical is true when the hash is already the
// route's own form; else the page replaces it with routeHash.
export function parseRoute(hash: string, fallback: string): { route: Route; canonical: boolean } {
  const parts = hash.replace(/^#?\/?/, '').split('/');
  let lang = fallback;
  const first = parts[0] ?? '';
  if (langCode.test(first) && !views.has(first)) {
    lang = first;
    parts.shift();
  }
  const [view = '', arg = ''] = parts;
  let route: Route;
  if (view === 'rules') {
    route = { lang, view: 'rules', rule: arg };
  } else if (view === 'check') {
    route = { lang, view: 'check', data: arg };
  } else {
    route = { lang, view: 'check', data: '' };
  }
  return { route, canonical: routeHash(route) === hash };
}

// routeHash is the hash of a route.
export function routeHash(r: Route): string {
  const arg = r.view === 'rules' ? r.rule : r.data;
  return `#/${r.lang}/${r.view}${arg === '' ? '' : `/${arg}`}`;
}

// rulesHash links a rule ID in a language.
export function rulesHash(lang: string, rule: string): string {
  return routeHash({ lang, view: 'rules', rule });
}

// checkHash links a situation's data in a language.
export function checkHash(lang: string, data: string): string {
  return routeHash({ lang, view: 'check', data });
}
