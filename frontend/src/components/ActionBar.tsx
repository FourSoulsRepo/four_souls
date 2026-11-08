import { useEffect, useRef, useState, type CSSProperties } from 'react';
import { feel } from '../feel';
import { t } from '../i18n';
import type { Banner } from '../prompt';
import { playSound } from '../sounds';
import { partsOf, secondsLeft } from '../stack';
import { useApp } from '../store';
import { playerColor } from '../table';
import { nextChange, tickSound, timerShows } from '../timer';
import { still } from './lunge';
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
      data-deadline={deadline || undefined}
    >
      <div className="prompt-title">
        {!banner ? (
          paused === 'vote' ? t('banner.vote') : t('banner.paused')
        ) : mine ? (
          banner.title
        ) : (
          <Waiting player={banner.player} name={name(banner.player)} />
        )}
        {deadline > 0 && <ResponseTimer key={deadline} deadline={deadline} ticks={mine} />}
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

// Timer counts down a Soft vote's timer; the response timer is
// ResponseTimer (20.15).
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

// ResponseTimer is the response timer (N-07) near its end (20.15, ADR
// 027 principle 9): nothing until its last feel.timer.show seconds;
// then the seconds left, and a fuse along the prompt's foot that burns
// down to its deadline. Each second it then loses ticks, and 0 ticks
// the last, for the player it waits for (ticks) only. A new deadline
// is a new timer (its key). Speed "off" and reduced motion: the fuse
// steps once a second.
function ResponseTimer({ deadline, ticks }: { deadline: number; ticks: boolean }) {
  const off = useApp((s) => s.settings.animation === 'off');
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const wait = nextChange(deadline, now);
    if (wait === undefined) {
      return;
    }
    const id = window.setTimeout(() => {
      setNow(Date.now());
    }, wait);
    return () => {
      window.clearTimeout(id);
    };
  }, [deadline, now]);
  const shows = timerShows(deadline, now);
  const left = secondsLeft(deadline, now);
  const lastRef = useRef<number | undefined>(undefined);
  useEffect(() => {
    if (!shows) {
      return;
    }
    const sound = lastRef.current === undefined ? undefined : tickSound(lastRef.current, left);
    lastRef.current = left;
    if (sound && ticks) {
      playSound(sound);
    }
  }, [shows, left, ticks]);
  if (!shows) {
    return null;
  }
  const step = off || still();
  return (
    <>
      <span className="action-timer" aria-label={t('actions.time_left')}>
        {' '}
        · {t('actions.seconds', { n: left })}
      </span>
      <Fuse deadline={deadline} left={step ? left / feel.timer.show : undefined} />
    </>
  );
}

// Fuse is the timer's fuse along the prompt's foot. It burns from where
// it is as it shows to empty at deadline, feel.timer.curve; its start
// is fixed as it mounts, so a render never moves it. left: it stands
// at that part instead, for speed "off" and reduced motion.
function Fuse({ deadline, left }: { deadline: number; left: number | undefined }) {
  const [shown] = useState(() => Date.now());
  const { show, curve } = feel.timer;
  const style: CSSProperties =
    left === undefined
      ? {
          animationDuration: `${String(show)}s`,
          animationDelay: `${String((deadline - shown) / 1000 - show)}s`,
          animationTimingFunction: curve,
        }
      : { transform: `scaleX(${String(left)})` };
  return <span className={`timer-fuse${left === undefined ? '' : ' still'}`} style={style} aria-hidden />;
}
