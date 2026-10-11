import { useState } from 'react';
import type { Update } from '../protocol';
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

// PauseOverlay: the game waits for players who left; the others vote
// to wait or to kick them (N-08).
export function PauseOverlay({ game, you, name }: { game: Update; you: number; name: (p: number) => string }) {
  const vote = useApp((s) => s.vote);
  const pause = game.pause;
  if (!pause) {
    return null;
  }
  const votes = pause.votes ?? [];
  const mine = votes.find((v) => v.seat === you);
  return (
    <div className="overlay" role="dialog" aria-label="Paused">
      <div className="overlay-box">
        <h2>Waiting for {pause.away.map(name).join(', ')}</h2>
        <p className="muted">They left the game. Wait for them to come back, or go on without them.</p>
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
