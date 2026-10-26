import { useEffect, useRef } from 'react';
import { stepName, undoButton, type Entry, type Undo } from '../history';
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
// card on hover or tap. In a game with rollback one Undo button under
// the title undoes the latest move (17.7). Undone rows are greyed, and
// a row tells who rolled back (ADR 020).
export function HistoryPanel({ hide }: { hide: () => void }) {
  const history = useApp((s) => s.history);
  const moves = useApp((s) => s.moves);
  const setZoom = useApp((s) => s.setZoom);
  const game = useApp((s) => s.game);
  const waiting = useApp((s) => s.behind || (s.game !== null && s.sent === s.game.step));
  const you = useApp(seatOf);
  const undo = game ? undoButton(game, you, moves) : null;
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
          <span className="history-mode" title="Any player may undo, at once, one move per press">
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
      {undo && <UndoButton undo={undo} disabled={waiting} />}
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
          </li>
        ))}
      </ol>
    </aside>
  );
}

// UndoButton is the History's one Undo button (17.7, ADR 020): it
// undoes the latest move, at once in Godmode, by a vote in Soft; the
// Soft author asks instead. Its tooltip names the move. Greyed, it
// still takes a click, to say why (15.29).
function UndoButton({ undo, disabled }: { undo: Undo; disabled: boolean }) {
  const rollBack = useApp((s) => s.rollBack);
  const showWhy = useApp((s) => s.showWhy);
  const setZoom = useApp((s) => s.setZoom);
  const { why, step } = undo;
  return (
    <div className="history-undo">
      <button
        type="button"
        title={undo.title}
        aria-disabled={why ? true : undefined}
        disabled={disabled && !why}
        onClick={(ev) => {
          if (why) {
            showWhy(why, ev.clientX, ev.clientY);
            return;
          }
          if (step !== undefined) {
            setZoom(null);
            rollBack(step, undo.ask);
          }
        }}
      >
        <span aria-hidden="true">↶ </span>
        {undo.label}
      </button>
    </div>
  );
}

// rowClass marks turn lines, undone rows and rollback rows.
function rowClass(e: Entry): string | undefined {
  const out = [e.icon === '▶' ? 'history-turn' : '', e.undone ? 'undone' : '', e.rollback ? 'history-rolled' : ''].filter(Boolean);
  return out.length > 0 ? out.join(' ') : undefined;
}
