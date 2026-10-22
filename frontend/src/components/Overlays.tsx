import { useEffect, useState } from 'react';
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
    <div className="overlay" role="dialog" aria-label="Connection lost">
      <div className="overlay-box">
        <h2>Connection lost</h2>
        <p className="muted">{error}</p>
        <p>Your seat waits for you. The others may vote to go on without you.</p>
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
            Reconnect
          </button>
          <button
            type="button"
            onClick={() => {
              void quit();
            }}
          >
            Leave
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
    <div className="reconnect-badge" role="status" aria-label="Reconnecting">
      {restarting ? 'The server is restarting. Reconnecting…' : 'Reconnecting…'}
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
            OK
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
    <div className="overlay" role="dialog" aria-label="Paused">
      <div className="overlay-box">
        <h2>Waiting for {pause.away.map(name).join(', ')}</h2>
        <p className="muted">
          {grace
            ? 'The server restarted. The others are coming back.'
            : 'They left the game. Wait for them to come back, or go on without them.'}
        </p>
        <ul className="overlay-votes">
          {votes.map((v) => (
            <li key={v.seat}>
              {name(v.seat)}: {v.kick ? 'go on without them' : 'wait'}
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
            Wait
          </button>
          <button
            type="button"
            className={mine?.kick ? 'chosen' : ''}
            disabled={grace}
            onClick={() => {
              vote(true);
            }}
          >
            Go on without them
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
    ? 'You win!'
    : winners.length > 0
      ? `${winners.map(name).join(' and ')} win${winners.length === 1 ? 's' : ''}!`
      : 'Game over';
  return (
    <div className="overlay" role="dialog" aria-label="Game over">
      <div className="overlay-box">
        <h2>{title}</h2>
        <p className="muted">Keep the match record to watch it again later.</p>
        <div className="overlay-buttons">
          <button
            type="button"
            className="primary"
            onClick={() => {
              void saveRecord().then(setSaved);
            }}
          >
            Save record
          </button>
          <button
            type="button"
            onClick={() => {
              void quit();
            }}
          >
            Back to menu
          </button>
        </div>
        <p className="overlay-status" role="status">
          {saved ? `Saved to ${saved}` : error}
        </p>
      </div>
    </div>
  );
}
