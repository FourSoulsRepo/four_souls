import { useEffect, useState } from 'react';
import { t } from '../i18n';
import type { Update } from '../protocol';
import { noticeText } from '../reconnect';
import { useApp } from '../store';
import './Overlays.css';

// LostOverlay: the connection dropped during the game. The seat waits
// on the server; reconnecting takes it back (N-08).
export function LostOverlay() {
  const error = useApp((s) => s.error);
  const reconnect = useApp((s) => s.reconnect);
  const quit = useApp((s) => s.quit);
  const [busy, setBusy] = useState(false);
  return (
    <div className="overlay" role="dialog" aria-label={t('lost.title')}>
      <div className="overlay-box">
        <h2>{t('lost.title')}</h2>
        <p className="muted">{error}</p>
        <p>{t('lost.text')}</p>
        <div className="overlay-buttons">
          <button
            type="button"
            className="primary"
            disabled={busy}
            onClick={() => {
              setBusy(true);
              void reconnect().finally(() => {
                setBusy(false);
              });
            }}
          >
            {t('lost.reconnect')}
          </button>
          <button
            type="button"
            onClick={() => {
              void quit();
            }}
          >
            {t('lost.leave')}
          </button>
        </div>
      </div>
    </div>
  );
}

// ReconnectBadge: the connection dropped and the app reconnects by
// itself; the table stays usable to look at (N-20). The "Connection
// lost" dialog replaces it after a while.
export function ReconnectBadge() {
  const restarting = useApp((s) => s.restarting);
  return (
    <div className="reconnect-badge" role="status" aria-label={t('lost.reconnecting_label')}>
      {restarting ? t('lost.restarting') : t('lost.reconnecting')}
    </div>
  );
}

// NoticeOverlay: a server notice, on any screen (N-20).
export function NoticeOverlay() {
  const notice = useApp((s) => s.notice);
  const endNotice = useApp((s) => s.endNotice);
  if (!notice) {
    return null;
  }
  const { title, text } = noticeText(notice);
  return (
    <div className="overlay" role="dialog" aria-label={title}>
      <div className="overlay-box">
        <h2>{title}</h2>
        <p>{text}</p>
        <div className="overlay-buttons">
          <button type="button" className="primary" onClick={endNotice}>
            {t('notice.ok')}
          </button>
        </div>
      </div>
    </div>
  );
}

// useNow is the time, renewed at until so the screen changes then.
function useNow(until: number): number {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    if (until <= now) {
      return;
    }
    const t = window.setTimeout(() => {
      setNow(Date.now());
    }, until - now);
    return () => {
      window.clearTimeout(t);
    };
  }, [until, now]);
  return now;
}

// PauseOverlay: the game waits for players who left; the others vote
// to wait or to kick them (N-08). Right after a server restart there is
// no kick vote for a while: the players are coming back (N-20).
export function PauseOverlay({ game, you, name }: { game: Update; you: number; name: (p: number) => string }) {
  const vote = useApp((s) => s.vote);
  const pause = game.pause;
  const now = useNow(pause?.grace ?? 0);
  if (!pause) {
    return null;
  }
  const votes = pause.votes ?? [];
  const mine = votes.find((v) => v.seat === you);
  const grace = (pause.grace ?? 0) > now;
  return (
    <div className="overlay" role="dialog" aria-label={t('pause.label')}>
      <div className="overlay-box">
        <h2>{t('pause.title', { players: pause.away.map(name).join(', ') })}</h2>
        <p className="muted">{grace ? t('pause.restarted') : t('pause.left')}</p>
        <ul className="overlay-votes">
          {votes.map((v) => (
            <li key={v.seat}>
              {name(v.seat)}: {v.kick ? t('pause.vote_kick') : t('pause.vote_wait')}
            </li>
          ))}
        </ul>
        <div className="overlay-buttons">
          <button
            type="button"
            className={mine && !mine.kick ? 'chosen' : ''}
            onClick={() => {
              vote(false);
            }}
          >
            {t('pause.wait')}
          </button>
          <button
            type="button"
            className={mine?.kick ? 'chosen' : ''}
            disabled={grace}
            onClick={() => {
              vote(true);
            }}
          >
            {t('pause.kick')}
          </button>
        </div>
      </div>
    </div>
  );
}

// EndOverlay: the winner, and the match record to keep (RP-05).
export function EndOverlay({ game, you, name }: { game: Update; you: number; name: (p: number) => string }) {
  const saveRecord = useApp((s) => s.saveRecord);
  const quit = useApp((s) => s.quit);
  const error = useApp((s) => s.error);
  const [saved, setSaved] = useState('');
  const winners = game.view.winners ?? [];
  const title = winners.includes(you)
    ? t('over.you_win')
    : winners.length === 1
      ? t('over.wins', { player: name(winners[0] ?? 0) })
      : winners.length > 1
        ? t('over.win', { players: winners.slice(0, -1).map(name).join(', '), last: name(winners.at(-1) ?? 0) })
        : t('over.title');
  return (
    <div className="overlay" role="dialog" aria-label={t('over.title')}>
      <div className="overlay-box">
        <h2>{title}</h2>
        <p className="muted">{t('over.hint')}</p>
        <div className="overlay-buttons">
          <button
            type="button"
            className="primary"
            onClick={() => {
              void saveRecord().then(setSaved);
            }}
          >
            {t('over.save')}
          </button>
          <button
            type="button"
            onClick={() => {
              void quit();
            }}
          >
            {t('over.menu')}
          </button>
        </div>
        <p className="overlay-status" role="status">
          {saved ? t('over.saved', { path: saved }) : error}
        </p>
      </div>
    </div>
  );
}
