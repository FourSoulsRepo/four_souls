import { useEffect, useRef } from 'react';
import { stepName } from '../history';
import type { GameView } from '../protocol';
import type { Names } from '../stack';
import { useApp } from '../store';
import { boxOf } from '../zoom';
import './InfoPanels.css';

// TurnBadge shows the turn number, whose turn it is and its step (V-05).
// Before turn 1 it names the first-player roll (15.43).
export function TurnBadge({ view, names }: { view: GameView; names: Names }) {
  const { turn } = view;
  return (
    <div className="turn-badge" role="status" aria-label="Turn">
      {view.first_roll ? (
        <>Before turn 1 · First player roll</>
      ) : (
        <>
          Turn {turn.number} · {names.player(turn.active)} · {stepName(turn.step)}
        </>
      )}
    </div>
  );
}

// HistoryPanel lists the match so far, newest at the bottom (V-03). The
// history strip opens it (15.22); hide closes it. Card names zoom their
// card on hover or tap.
export function HistoryPanel({ hide }: { hide: () => void }) {
  const history = useApp((s) => s.history);
  const setZoom = useApp((s) => s.setZoom);
  const listRef = useRef<HTMLOListElement>(null);
  useEffect(() => {
    const el = listRef.current;
    if (el) {
      el.scrollTop = el.scrollHeight;
    }
  }, [history.length]);
  return (
    <aside className="history" aria-label="History">
      <header>
        <h2>History</h2>
        <button type="button" onClick={hide}>
          Hide
        </button>
      </header>
      <ol ref={listRef} className="selectable">
        {history.map((e) => (
          <li
            key={e.key}
            className={e.icon === '▶' ? 'history-turn' : undefined}
            title={e.turn > 0 ? `Turn ${String(e.turn)}: ${e.text}` : e.text}
            onPointerEnter={(ev) => {
              if (e.card && ev.pointerType === 'mouse') {
                setZoom({ ref: e.card, touch: false, at: boxOf(ev.currentTarget) });
              }
            }}
            onPointerLeave={(ev) => {
              if (e.card && ev.pointerType === 'mouse') {
                setZoom(null);
              }
            }}
            onClick={(ev) => {
              if (e.card) {
                setZoom({ ref: e.card, touch: true, at: boxOf(ev.currentTarget) });
              }
            }}
          >
            <span className="history-icon" aria-hidden="true">
              {e.icon}
            </span>
            <span>{e.text}</span>
          </li>
        ))}
      </ol>
    </aside>
  );
}
