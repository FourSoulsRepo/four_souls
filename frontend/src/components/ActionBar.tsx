import { useEffect, useState } from 'react';
import { t } from '../i18n';
import type { Banner } from '../prompt';
import { partsOf, secondsLeft } from '../stack';
import { useApp } from '../store';
import { playerColor } from '../table';
import './Actions.css';

interface Props {
  // banner is what the game waits for (15.28); none once it is over.
  banner: Banner | null;
  // paused: the game waits for a player who left, or for a rollback
  // vote; it names which.
  paused: '' | 'paused' | 'vote';
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
    <div className="action-bar" role="region" aria-label={t('actions.label')}>
      {(paused || banner) && (
        <PromptBanner banner={paused ? null : banner} paused={paused} name={name} deadline={paused ? 0 : deadline} />
      )}
      {wait !== 'none' && (
        <div className="action-buttons">
          <button
            type="button"
            className="secondary"
            title={t('actions.wait_tip')}
            onClick={() => {
              waitTurn(wait === 'offer');
            }}
          >
            {wait === 'offer' ? t('actions.wait') : t('actions.cancel_wait')}
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
function PromptBanner({
  banner,
  paused,
  name,
  deadline,
}: {
  banner: Banner | null;
  paused: '' | 'paused' | 'vote';
  name: (player: number) => string;
  deadline: number;
}) {
  const mine = banner?.mine ?? false;
  return (
    <section
      className={`prompt-banner${mine ? ' mine' : ''}`}
      aria-label={t('banner.label')}
      aria-live="polite"
      data-kind={banner?.kind}
    >
      <div className="prompt-title">
        {!banner ? (
          paused === 'vote' ? t('banner.vote') : t('banner.paused')
        ) : mine ? (
          banner.title
        ) : (
          <Waiting player={banner.player} name={name(banner.player)} />
        )}
        {deadline > 0 && <Timer deadline={deadline} />}
      </div>
      {banner?.hint && <div className="prompt-hint">{banner.hint}</div>}
    </section>
  );
}

// Waiting is "Waiting for Ann", the name in the player's color.
function Waiting({ player, name }: { player: number; name: string }) {
  return (
    <>
      {partsOf(t('banner.waiting_for'), { player: { text: name, role: 'player', player } }).map((p, i) =>
        p.role ? (
          // A line's parts never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <b key={i} style={{ color: playerColor(player) }}>
            {p.text}
          </b>
        ) : (
          p.text
        ),
      )}
    </>
  );
}

// Timer counts down the response timer (N-07), or a Soft vote's.
export function Timer({ deadline }: { deadline: number }) {
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
    <span className="action-timer" aria-label={t('actions.time_left')}>
      {' '}
      · {t('actions.seconds', { n: secondsLeft(deadline, now) })}
    </span>
  );
}
