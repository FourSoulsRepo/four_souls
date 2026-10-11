import { useEffect, useRef, useState } from 'react';
import { stepName } from '../history';
import type { GameView } from '../protocol';
import type { Names } from '../stack';
import { useApp } from '../store';
import './InfoPanels.css';

// TurnBadge shows the turn number, whose turn it is and its step (V-05).
export function TurnBadge({ view, names }: { view: GameView; names: Names }) {
  const { turn } = view;
  return (
    <div className="turn-badge" role="status" aria-label="Turn">
      Turn {turn.number} · {names.player(turn.active)} · {stepName(turn.step)}
    </div>
  );
}

// HistoryPanel lists the match so far, newest at the bottom (V-03). It
// can be hidden; card names zoom their card on hover or tap.
export function HistoryPanel() {
  const history = useApp((s) => s.history);
  const setZoom = useApp((s) => s.setZoom);
  const [open, setOpen] = useState(true);
  const listRef = useRef<HTMLOListElement>(null);
  useEffect(() => {
    const el = listRef.current;
    if (el) {
      el.scrollTop = el.scrollHeight;
    }
  }, [history.length, open]);
  if (!open) {
    return (
      <button
        type="button"
        className="history-show"
        onClick={() => {
          setOpen(true);
        }}
      >
        History
      </button>
    );
  }
  return (
    <aside className="history" aria-label="History">
      <header>
        <h2>History</h2>
        <button
          type="button"
          onClick={() => {
            setOpen(false);
          }}
        >
          Hide
        </button>
      </header>
      <ol ref={listRef}>
        {history.map((e) => (
          <li
            key={e.key}
            className={e.icon === '▶' ? 'history-turn' : undefined}
            title={`Turn ${String(e.turn)}: ${e.text}`}
            onPointerEnter={(ev) => {
              if (e.card && ev.pointerType === 'mouse') {
                setZoom({ ref: e.card, touch: false, side: 'right' });
              }
            }}
            onPointerLeave={(ev) => {
              if (e.card && ev.pointerType === 'mouse') {
                setZoom(null);
              }
            }}
            onClick={() => {
              if (e.card) {
                setZoom({ ref: e.card, touch: true, side: 'right' });
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
