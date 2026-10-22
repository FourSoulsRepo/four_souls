import { useEffect, useState } from 'react';
import type { Banner } from '../prompt';
import { secondsLeft } from '../stack';
import { useApp } from '../store';
import { playerColor } from '../table';
import './Actions.css';

interface Props {
  // banner is what the game waits for (15.28); none once it is over.
  banner: Banner | null;
  paused: boolean;
  name: (player: number) => string;
  // deadline is when the response timer answers, Unix ms; 0: no timer.
  deadline: number;
  // wait: "Wait for my turn" is offered, or on (N-16).
  wait: 'offer' | 'on' | 'none';
}

// ActionBar shows the prompt banner with the response timer, "Wait for
// my turn", and why the last action was refused (7.7, 7.8). Pass lives
// in the stack panel (15.21).
export function ActionBar({ banner, paused, name, deadline, wait }: Props) {
  const waitTurn = useApp((s) => s.waitTurn);
  // A lost connection says so in its badge or dialog, not here (N-20).
  const error = useApp((s) => (s.connected ? s.error : ''));
  return (
    <div className="action-bar" role="region" aria-label="Actions">
      {(paused || banner) && (
        <PromptBanner banner={paused ? null : banner} name={name} deadline={paused ? 0 : deadline} />
      )}
      {wait !== 'none' && (
        <div className="action-buttons">
          <button
            type="button"
            className="secondary"
            title="Pass until the turn before yours ends, or someone targets you"
            onClick={() => {
              waitTurn(wait === 'offer');
            }}
          >
            {wait === 'offer' ? 'Wait for my turn' : 'Cancel wait'}
          </button>
        </div>
      )}
      <div className="action-error" role="status">
        {error}
      </div>
    </div>
  );
}

// PromptBanner says what the game waits for (15.28, V-20): what to do,
// in a few words, when it waits for this player; whom it waits for and
// what for, in their color, when it waits for someone else. No banner
// means the game is paused.
function PromptBanner({ banner, name, deadline }: { banner: Banner | null; name: (player: number) => string; deadline: number }) {
  const mine = banner?.mine ?? false;
  return (
    <section
      className={`prompt-banner${mine ? ' mine' : ''}`}
      aria-label="Prompt"
      aria-live="polite"
      data-kind={banner?.kind}
    >
      <div className="prompt-title">
        {!banner ? (
          'Paused'
        ) : mine ? (
          banner.title
        ) : (
          <>
            Waiting for <b style={{ color: playerColor(banner.player) }}>{name(banner.player)}</b>
          </>
        )}
        {deadline > 0 && <Timer deadline={deadline} />}
      </div>
      {banner?.hint && <div className="prompt-hint">{banner.hint}</div>}
    </section>
  );
}

// Timer counts down the response timer (N-07).
function Timer({ deadline }: { deadline: number }) {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const t = window.setInterval(() => {
      setNow(Date.now());
    }, 250);
    return () => {
      window.clearInterval(t);
    };
  }, []);
  return (
    <span className="action-timer" aria-label="Time left">
      {' '}
      · {secondsLeft(deadline, now)} s
    </span>
  );
}
