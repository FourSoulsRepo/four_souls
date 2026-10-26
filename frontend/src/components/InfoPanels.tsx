import { useEffect, useRef } from 'react';
import { rollbackButtons, stepName, type Entry, type RollbackButton } from '../history';
import type { GameView } from '../protocol';
import type { Names } from '../stack';
import { seatOf, useApp } from '../store';
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
// card on hover or tap. In Godmode each step's first row has a rollback
// button; in Soft the rows after the last barrier, an "ask" button for
// the step's own author. Undone rows are greyed, and a row tells who
// rolled back (ADR 020).
export function HistoryPanel({ hide }: { hide: () => void }) {
  const history = useApp((s) => s.history);
  const setZoom = useApp((s) => s.setZoom);
  const rollBack = useApp((s) => s.rollBack);
  const game = useApp((s) => s.game);
  const waiting = useApp((s) => s.behind || (s.game !== null && s.sent === s.game.step));
  const you = useApp(seatOf);
  const buttons = game ? rollbackButtons(history, game, you) : new Map<number, RollbackButton>();
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
        {game?.rollback_mode === 'god' && (
          <span className="history-mode" title="Any player may roll back any step, at once">
            Godmode
          </span>
        )}
        {game?.rollback_mode === 'soft' && (
          <span className="history-mode" title="The other players vote on a rollback; dice, damage, hidden cards and a turn's end stop it">
            Soft rollback
          </span>
        )}
        <button type="button" onClick={hide}>
          Hide
        </button>
      </header>
      <ol ref={listRef} className="selectable">
        {history.map((e) => (
          <li
            key={e.key}
            className={rowClass(e)}
            title={`${e.undone ? 'Undone: ' : ''}${e.turn > 0 ? `Turn ${String(e.turn)}: ${e.text}` : e.text}`}
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
            {buttons.has(e.key) && e.step !== undefined && (
              <RollbackRowButton
                kind={buttons.get(e.key) ?? 'roll'}
                text={e.text}
                disabled={waiting}
                press={() => {
                  setZoom(null);
                  rollBack(e.step ?? 0, buttons.get(e.key) === 'ask');
                }}
              />
            )}
          </li>
        ))}
      </ol>
    </aside>
  );
}

// RollbackRowButton is a history row's rollback button: ↶ rolls back
// or opens a vote; the author's own step in Soft gets "Ask" instead.
function RollbackRowButton({ kind, text, disabled, press }: { kind: RollbackButton; text: string; disabled: boolean; press: () => void }) {
  const ask = kind === 'ask';
  return (
    <button
      type="button"
      className={`history-rollback${ask ? ' ask' : ''}`}
      aria-label={`${ask ? 'Ask to undo' : 'Roll back'}: ${text}`}
      title={
        ask
          ? 'Ask the others to undo this and everything after it'
          : kind === 'vote'
            ? 'Undo this and everything after it, if the others agree'
            : 'Undo this and everything after it'
      }
      disabled={disabled}
      onClick={(ev) => {
        ev.stopPropagation();
        press();
      }}
    >
      {ask ? 'Ask ↶' : '↶'}
    </button>
  );
}

// rowClass marks turn lines, undone rows and rollback rows.
function rowClass(e: Entry): string | undefined {
  const out = [e.icon === '▶' ? 'history-turn' : '', e.undone ? 'undone' : '', e.rollback ? 'history-rolled' : ''].filter(Boolean);
  return out.length > 0 ? out.join(' ') : undefined;
}
