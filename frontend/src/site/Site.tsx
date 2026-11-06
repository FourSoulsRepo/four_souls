// The rules website (ADR 026): a header, the view the route names, and
// a footer with the fan-game notice. The engine loads in the background;
// the rules text needs no engine.
import { useEffect, useState, useSyncExternalStore } from 'react';
import { setLanguage, shownLanguage, t } from '../i18n';
import { CheckView } from './CheckView';
import { LoadError, loadEngine, type Loaded, type NoticeLine } from './engine';
import { checkHash, parseRoute, routeHash, rulesHash } from './route';
import { RulesView } from './RulesView';

function subscribe(onChange: () => void): () => void {
  window.addEventListener('hashchange', onChange);
  return () => {
    window.removeEventListener('hashchange', onChange);
  };
}

function hashNow(): string {
  return window.location.hash;
}

type Engine = { state: 'loading'; loaded: number; total: number } | { state: 'ready'; loaded: Loaded } | { state: 'failed'; error: LoadError };

// One engine per page, however often the shell mounts.
let starting: Promise<Loaded> | undefined;

function noticeKey(line: NoticeLine): string {
  return line.parts.map((p) => p.text).join('');
}

export function Site() {
  const hash = useSyncExternalStore(subscribe, hashNow);
  const { route, canonical } = parseRoute(hash, shownLanguage(''));
  setLanguage(route.lang);
  const [engine, setEngine] = useState<Engine>({ state: 'loading', loaded: 0, total: 0 });

  useEffect(() => {
    if (!canonical) {
      window.history.replaceState(null, '', routeHash(route));
    }
    document.documentElement.lang = route.lang;
    document.title = t('site.title');
  });

  useEffect(() => {
    starting ??= loadEngine(import.meta.env.BASE_URL, (loaded, total) => {
      setEngine({ state: 'loading', loaded, total });
    });
    starting.then(
      (loaded) => {
        setEngine({ state: 'ready', loaded });
      },
      (e: unknown) => {
        setEngine({ state: 'failed', error: e instanceof LoadError ? e : new LoadError('failed', String(e)) });
      },
    );
  }, []);

  const loaded = engine.state === 'ready' ? engine.loaded : undefined;

  return (
    <div className="site">
      <header className="site-header">
        <a className="site-name" href={checkHash(route.lang, '')}>
          {t('site.title')}
        </a>
        <nav>
          <a href={checkHash(route.lang, '')} aria-current={route.view === 'check' ? 'page' : undefined}>
            {t('site.nav.check')}
          </a>
          <a href={rulesHash(route.lang, '')} aria-current={route.view === 'rules' ? 'page' : undefined}>
            {t('site.nav.rules')}
          </a>
        </nav>
      </header>
      <main>
        {route.view === 'check' && engine.state === 'loading' && (
          <div className="loading" role="status">
            <span>{t('site.load.loading')}</span>
            <progress value={engine.total > 0 ? engine.loaded : undefined} max={engine.total > 0 ? engine.total : undefined} />
          </div>
        )}
        {engine.state === 'failed' && (
          <p className="note warn" role="alert">
            {engine.error.reason === 'no_decompression' ? t('site.load.no_decompression') : t('site.load.failed', { error: engine.error.message })}
          </p>
        )}
        {route.view === 'check' ? <CheckView lang={route.lang} data={route.data} loaded={loaded} /> : <RulesView lang={route.lang} rule={route.rule} />}
      </main>
      <footer className="site-footer">
        {loaded && (
          <>
            <div className="notice">
              {loaded.notice.map((line) => (
                <p key={noticeKey(line)}>
                  {line.parts.map((p) =>
                    p.url ? (
                      <a key={p.text} href={p.url} target="_blank" rel="noopener noreferrer">
                        {p.text}
                      </a>
                    ) : (
                      <span key={p.text}>{p.text}</span>
                    ),
                  )}
                </p>
              ))}
            </div>
            <p className="version">{t('site.footer.version', { engine: loaded.version.engine, app: loaded.version.app })}</p>
          </>
        )}
      </footer>
    </div>
  );
}
