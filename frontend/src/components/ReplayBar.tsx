import type { CSSProperties, ReactNode } from 'react';
import { t } from '../i18n';
import { useEscape } from '../keys';
import { replaySpeeds, sliderFill, turnAt } from '../replay';
import { useApp } from '../store';
import './Replay.css';

// Icon draws one of the bar's ink signs; the button around it carries
// the name (10.7).
function Icon({ children }: { children: ReactNode }) {
  return (
    <svg className="replay-icon" viewBox="0 0 16 16" aria-hidden="true">
      {children}
    </svg>
  );
}

const backIcon = (
  <Icon>
    <rect x="2" y="3" width="2.5" height="10" rx="0.6" />
    <path d="M14 3.2v9.6a.6.6 0 0 1-.9.5L5.6 8.5a.6.6 0 0 1 0-1l7.5-4.8a.6.6 0 0 1 .9.5z" />
  </Icon>
);
const forwardIcon = (
  <Icon>
    <rect x="11.5" y="3" width="2.5" height="10" rx="0.6" />
    <path d="M2 3.2v9.6a.6.6 0 0 0 .9.5l7.5-4.8a.6.6 0 0 0 0-1L2.9 2.7a.6.6 0 0 0-.9.5z" />
  </Icon>
);
const playIcon = (
  <Icon>
    <path d="M4 2.4v11.2a.7.7 0 0 0 1.1.6l8.6-5.6a.7.7 0 0 0 0-1.2L5.1 1.8a.7.7 0 0 0-1.1.6z" />
  </Icon>
);
const pauseIcon = (
  <Icon>
    <rect x="3" y="2.5" width="3.5" height="11" rx="0.8" />
    <rect x="9.5" y="2.5" width="3.5" height="11" rx="0.8" />
  </Icon>
);
const hideIcon = (
  <Icon>
    <path d="M3 5.5l5 5 5-5" fill="none" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round" strokeLinejoin="round" />
  </Icon>
);
const showIcon = (
  <Icon>
    <path d="M3 10.5l5-5 5 5" fill="none" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round" strokeLinejoin="round" />
  </Icon>
);

// ReplayBar steps through the record on the table (10.2, ADR 031, 4):
// play and pause, a step back or forward, a jump to a turn, a slider
// over all steps, the speed, and whether hands show (10.3). It wears
// the game's parchment buttons on a dark panel (10.7, ADR 031, 4.8).
// The viewer can hide it; a small button in the same corner brings it
// back, and the autoplay goes on meanwhile (4.9). Esc goes back to the
// records, shown or hidden.
export function ReplayBar() {
  const replay = useApp((s) => s.replay);
  const turn = useApp((s) => s.game?.view.turn.number ?? 0);
  const error = useApp((s) => s.error);
  const replayTo = useApp((s) => s.replayTo);
  const replayPlay = useApp((s) => s.replayPlay);
  const replaySpeed = useApp((s) => s.replaySpeed);
  const replayHands = useApp((s) => s.replayHands);
  const replayBar = useApp((s) => s.replayBar);
  const leave = useApp((s) => s.leaveReplay);
  useEscape(leave);
  if (!replay) {
    return null;
  }
  const { at, speed, playing, hands, summary, bar } = replay;
  if (!bar) {
    return (
      <section className="panel replay-bar replay-hidden" role="region" aria-label={t('replay.label')}>
        <button
          type="button"
          className="replay-show"
          aria-label={t('replay.show')}
          title={t('replay.show')}
          onClick={() => {
            replayBar(true);
          }}
        >
          {showIcon}
        </button>
      </section>
    );
  }
  const last = summary.steps - 1;
  const turns = summary.turns;
  const goTo = (index: number) => {
    void replayTo(index);
  };
  const fill = { '--fill': `${String(sliderFill(at, summary.steps))}%` } as CSSProperties;
  return (
    <section className="panel replay-bar" role="region" aria-label={t('replay.label')}>
      <div className="replay-row">
        <span className="replay-where" role="status">
          {t('replay.where', { i: at + 1, n: summary.steps, turn })}
        </span>
        <button
          type="button"
          className="replay-hide"
          aria-label={t('replay.hide')}
          title={t('replay.hide')}
          onClick={() => {
            replayBar(false);
          }}
        >
          {hideIcon}
        </button>
      </div>
      <input
        className="replay-slider"
        type="range"
        min={1}
        max={summary.steps}
        value={at + 1}
        style={fill}
        aria-label={t('replay.slider')}
        onChange={(e) => {
          goTo(Number(e.target.value) - 1);
        }}
      />
      <div className="replay-row">
        <button
          type="button"
          className="replay-icon-button"
          aria-label={t('replay.back')}
          title={t('replay.back')}
          disabled={at <= 0}
          onClick={() => {
            goTo(at - 1);
          }}
        >
          {backIcon}
        </button>
        <button
          type="button"
          className="replay-icon-button replay-play"
          aria-label={playing ? t('replay.pause') : t('replay.play')}
          title={playing ? t('replay.pause') : t('replay.play')}
          disabled={!playing && at >= last}
          onClick={() => {
            replayPlay(!playing);
          }}
        >
          {playing ? pauseIcon : playIcon}
        </button>
        <button
          type="button"
          className="replay-icon-button"
          aria-label={t('replay.forward')}
          title={t('replay.forward')}
          disabled={at >= last}
          onClick={() => {
            goTo(at + 1);
          }}
        >
          {forwardIcon}
        </button>
        <div className="replay-speeds" role="group" aria-label={t('replay.speed')}>
          {replaySpeeds.map((s) => (
            <button
              key={s}
              type="button"
              aria-pressed={s === speed}
              onClick={() => {
                replaySpeed(s);
              }}
            >
              {t('replay.speed_n', { n: s })}
            </button>
          ))}
        </div>
      </div>
      <div className="replay-row">
        <select
          className="replay-turns"
          aria-label={t('replay.turn_pick')}
          value={turnAt(turns, at)}
          onChange={(e) => {
            const i = Number(e.target.value);
            const to = turns.at(i);
            if (i >= 0 && to) {
              goTo(to.step);
            }
          }}
        >
          {turnAt(turns, at) < 0 && <option value={-1}>{t('replay.turn_pick')}</option>}
          {turns.map((tn, i) => (
            <option key={tn.step} value={i}>
              {tn.number > 0 ? t('replay.turn', { n: tn.number }) : t('replay.turn_zero')}
            </option>
          ))}
        </select>
        <button
          type="button"
          className="replay-hands"
          aria-pressed={hands}
          onClick={() => {
            void replayHands(!hands);
          }}
        >
          {t('replay.hands')}
        </button>
      </div>
      <div className="replay-row">
        <button type="button" className="replay-leave" onClick={leave}>
          {t('replay.leave')}
        </button>
      </div>
      {error && (
        <p className="replay-error" role="alert">
          {error}
        </p>
      )}
    </section>
  );
}
