// The check view (12.2): a situation payload, pasted or from a link,
// and the engine's answer with its rule. The refusal is a message key
// with values, shown through say() like the app's server messages.
import { useEffect, useMemo, useState } from 'react';
import { inLanguage } from '../cards';
import { language, say, t } from '../i18n';
import { plainNames, type MsgNames } from '../protocol';
import { parse, type Answer, type Loaded } from './engine';
import { compact, decodePayload, encodePayload, linkWarnAt } from './link';
import { checkHash, rulesHash } from './route';

// examples are the payloads internal/website's tests check, by title.
const exampleFiles = import.meta.glob<string>('../../../internal/website/testdata/*.json', { query: '?raw', import: 'default', eager: true });
const examples = Object.values(exampleFiles).map((text) => ({ text, title: titleOf(text) }));

// titleOf is a payload's title, '' when it has none or is no JSON.
function titleOf(payload: string): string {
  try {
    const v = JSON.parse(payload) as { title?: unknown } | null;
    return typeof v?.title === 'string' ? v.title : '';
  } catch {
    return '';
  }
}

// pretty indents a JSON payload for the text box.
function pretty(payload: string): string {
  try {
    return JSON.stringify(JSON.parse(payload), null, 2);
  } catch {
    return payload;
  }
}

type Result = { answer: Answer } | { error: string };

function run(loaded: Loaded, payload: string): Result {
  try {
    return { answer: parse(loaded.engine.check(payload)) as Answer };
  } catch (e) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

// linkOf is the full address of a check link.
function linkOf(hash: string): string {
  return `${window.location.href.split('#')[0] ?? ''}${hash}`;
}

interface Props {
  lang: string;
  // data is the link's <data>; '' for no link.
  data: string;
  loaded: Loaded | undefined;
}

export function CheckView({ lang, data, loaded }: Props) {
  const [text, setText] = useState('');
  // checked is the payload the answer is for.
  const [checked, setChecked] = useState<string>();
  const [badLink, setBadLink] = useState(false);
  const [status, setStatus] = useState('');
  const [link, setLink] = useState('');

  useEffect(() => {
    if (data === '') {
      return;
    }
    let live = true;
    decodePayload(data).then(
      (payload) => {
        if (live) {
          setBadLink(false);
          setText(pretty(payload));
          setChecked(payload);
        }
      },
      () => {
        if (live) {
          setBadLink(true);
        }
      },
    );
    return () => {
      live = false;
    };
  }, [data]);

  const result = useMemo(() => (loaded && checked !== undefined ? run(loaded, checked) : undefined), [loaded, checked]);

  const names: MsgNames = {
    player: plainNames.player,
    card: (ref) => {
      const card = loaded?.cards.find((c) => c.ref === ref);
      return card ? inLanguage(card, language()).name : ref;
    },
  };

  const check = (payload: string) => {
    setChecked(payload);
    setLink('');
    setStatus('');
    // The address becomes the situation's link, so a reload keeps it.
    void encodePayload(compact(payload)).then((d) => {
      window.history.replaceState(null, '', checkHash(lang, d));
    });
  };

  const copy = async (what: string, done: string) => {
    try {
      await navigator.clipboard.writeText(what);
      setStatus(done);
    } catch {
      setStatus(t('site.check.copy_failed'));
    }
  };

  const copyLink = async () => {
    const url = linkOf(checkHash(lang, await encodePayload(compact(text))));
    setLink(url);
    await copy(url, t('site.check.copied_link'));
  };

  const title = checked === undefined ? '' : titleOf(checked);

  return (
    <div className="check">
      <p className="intro">{t('site.check.intro')}</p>
      {badLink && <p className="note warn">{t('site.check.bad_link')}</p>}
      <label className="field">
        <span>{t('site.check.payload')}</span>
        <textarea
          value={text}
          spellCheck={false}
          rows={12}
          onChange={(e) => {
            setText(e.target.value);
          }}
        />
      </label>
      <div className="buttons">
        <button
          type="button"
          className="primary"
          disabled={text.trim() === ''}
          onClick={() => {
            check(text);
          }}
        >
          {t('site.check.run')}
        </button>
        <button type="button" disabled={text.trim() === ''} onClick={() => void copyLink()}>
          {t('site.check.copy_link')}
        </button>
        <button type="button" disabled={text.trim() === ''} onClick={() => void copy(text, t('site.check.copied_payload'))}>
          {t('site.check.copy_payload')}
        </button>
      </div>
      {status !== '' && (
        <p className="note" role="status">
          {status}
        </p>
      )}
      {link !== '' && (
        <div className="link">
          <input
            readOnly
            value={link}
            aria-label={t('site.check.link')}
            onFocus={(e) => {
              e.target.select();
            }}
          />
          {link.length > linkWarnAt && <p className="note warn">{t('site.check.long_link', { n: link.length })}</p>}
        </div>
      )}
      {examples.length > 0 && (
        <div className="examples">
          <span>{t('site.check.examples')}</span>
          {examples.map((ex) => (
            <button
              key={ex.title}
              type="button"
              className="example"
              onClick={() => {
                setText(ex.text.trim());
                check(ex.text);
              }}
            >
              {ex.title}
            </button>
          ))}
        </div>
      )}
      {checked !== undefined && result === undefined && <p className="note">{t('site.check.waiting')}</p>}
      {result && <Answered result={result} title={title} lang={lang} names={names} />}
    </div>
  );
}

function Answered({ result, title, lang, names }: { result: Result; title: string; lang: string; names: MsgNames }) {
  if ('error' in result) {
    return (
      <section className="answer error" aria-live="polite">
        <h2>{t('site.check.error')}</h2>
        <p>
          <code>{result.error}</code>
        </p>
      </section>
    );
  }
  const a = result.answer;
  const events = a.events ?? [];
  return (
    <section className={a.allowed ? 'answer ok' : 'answer no'} aria-live="polite">
      {title !== '' && <p className="situation">{title}</p>}
      <h2>{a.allowed ? t('site.check.allowed') : t('site.check.refused')}</h2>
      {!a.allowed && <p className="reason">{say(a.say, names) || a.reason}</p>}
      {a.rule && (
        <p className="rule-line">
          {t('site.check.rule')}{' '}
          <a className="rule-ref" href={rulesHash(lang, a.rule)}>
            {a.rule}
          </a>
        </p>
      )}
      {a.matches !== undefined && (
        <div className={a.matches ? 'expect ok' : 'expect no'}>
          <p>{a.matches ? t('site.check.matches') : t('site.check.mismatch')}</p>
          {(a.mismatch ?? []).length > 0 && (
            <ul>
              {(a.mismatch ?? []).map((m) => (
                <li key={m}>{m}</li>
              ))}
            </ul>
          )}
        </div>
      )}
      {events.length > 0 && (
        <details className="events">
          <summary>{t('site.check.events', { n: events.length })}</summary>
          <ol>
            {events.map((e, i) => (
              // Events never reorder: their place is their key.
              // eslint-disable-next-line @eslint-react/no-array-index-key
              <li key={i}>
                {names.player(e.player)} · <code>{e.kind}</code>
                {e.label?.key ? ` · ${say(e.label, names)}` : ''}
                {e.card ? ` · ${names.card(e.card)}` : ''}
                {e.amount ? ` · ${String(e.amount)}` : ''}
              </li>
            ))}
          </ol>
        </details>
      )}
    </section>
  );
}
