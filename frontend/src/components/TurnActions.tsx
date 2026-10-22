import { useState, type ReactNode } from 'react';
import { createPortal } from 'react-dom';
import { IntentKind, PromptKind, Purpose, type TurnActions as Allowed } from '../actions';
import { useKey } from '../keys';
import type { Intent, Prompt } from '../protocol';
import { useApp } from '../store';
import './Actions.css';
import './Overlays.css';

interface Props {
  allowed: Allowed;
  // disabled: the player may not act now (behind, or already sent).
  disabled: boolean;
  // step is the update's step; a press shows until the next one.
  step: number;
  waiting: Prompt;
  you: number;
  // left is how many buys and attacks are left this turn; only on the
  // viewer's own turn (15.24).
  left?: { buy: number; attack: number };
}

interface Tool {
  kind: number;
  label: string;
  key: string;
  icon: ReactNode;
  intent?: Intent;
  // why says why it can't be used now (15.29).
  why?: string | undefined;
  // choosing: the player is picking what this action is for.
  choosing: boolean;
  choosingLabel: string;
  left?: number;
}

// TurnActions are buy and attack as icon buttons by own name (15.19,
// F-11): lit when the engine allows them, dim when not (F-12). A press
// shows at once, before the server answers (V-15). While the player
// picks what to buy or attack, that icon stays pressed (F-13). On own
// turn a badge on each icon counts what is left (15.24). A click on a
// dim icon says why it can't be used (15.29, F-18).
export function TurnActions({ allowed, disabled, step, waiting, you, left }: Props) {
  const act = useApp((s) => s.act);
  const showWhy = useApp((s) => s.showWhy);
  const [pressed, setPressed] = useState<{ kind: number; step: number } | null>(null);
  const picking = waiting.kind === PromptKind.choose && waiting.player === you ? (waiting.purpose ?? -1) : -1;
  const tools: Tool[] = [
    {
      kind: IntentKind.purchase,
      label: 'Buy',
      key: 'b',
      icon: <BagIcon />,
      intent: allowed.buy,
      why: allowed.why?.buy,
      choosing: picking === Purpose.purchase,
      choosingLabel: 'Choosing what to buy',
      left: left?.buy,
    },
    {
      kind: IntentKind.attack,
      label: 'Attack',
      key: 'a',
      icon: <SwordIcon />,
      intent: allowed.attack,
      why: allowed.why?.attack,
      choosing: picking === Purpose.attackTarget,
      choosingLabel: 'Choosing what to attack',
      left: left?.attack,
    },
  ];
  const press = (t: Tool) => {
    if (!t.intent || disabled) {
      return;
    }
    setPressed({ kind: t.kind, step });
    act(t.intent);
  };

  // Keys: B buys, A attacks; the tooltip names them.
  const [buy, attack] = tools as [Tool, Tool];
  useKey(buy.key, () => {
    press(buy);
  });
  useKey(attack.key, () => {
    press(attack);
  });

  return (
    <div className="turn-actions" role="toolbar" aria-label="Turn actions">
      {tools.map((t) => {
        const lit = t.intent !== undefined && !disabled;
        const down = t.choosing || (pressed?.kind === t.kind && pressed.step === step);
        // A dim icon with a reason still takes a click, to say it.
        const why = !lit && !disabled && !t.choosing ? t.why : undefined;
        return (
          <button
            key={t.kind}
            type="button"
            className={`turn-action${lit ? ' lit' : ''}${down ? ' down' : ''}`}
            aria-label={t.choosing ? t.choosingLabel : t.label}
            aria-pressed={down}
            aria-disabled={why ? true : undefined}
            title={`${t.choosing ? t.choosingLabel : t.label} (${t.key.toUpperCase()})${why ? ` · ${why}` : ''}`}
            disabled={!lit && !why}
            onClick={(e) => {
              if (why) {
                showWhy(why, e.clientX, e.clientY);
                return;
              }
              press(t);
            }}
          >
            {t.icon}
            {t.left !== undefined && <Left n={t.left} />}
          </button>
        );
      })}
    </div>
  );
}

// Left is the count badge on a turn icon; past 99 it says "99+".
function Left({ n }: { n: number }) {
  const text = n > 99 ? '99+' : String(n);
  const size = text.length > 2 ? ' long' : text.length > 1 ? ' wide' : '';
  return (
    <span className={`turn-action-left${size}${n === 0 ? ' none' : ''}`} aria-label={`${String(n)} left`}>
      {text}
    </span>
  );
}

interface EndProps {
  // intent ends the turn, when the engine allows it now.
  intent?: Intent;
  // idle: nothing is left to do but end the turn (H-20).
  idle: boolean;
  // left names what the player may still do; the question lists it (V-46).
  left?: string[] | undefined;
  // own: it is the viewer's turn.
  own: boolean;
  disabled: boolean;
  step: number;
  // why says why the player can't end the turn now (15.29).
  why?: string | undefined;
}

// EndTurn is the one big End turn button at the board's right edge
// (15.20, H-04): yellow on own turn, grey "Their turn" otherwise
// (H-19), glowing green once nothing is left to do (H-20). E presses
// it (F-15). On own turn, a click on it while it is dim says why
// (15.29, F-18). While plays are left it asks first, unless the
// settings say not to (V-46).
export function EndTurn({ intent, idle, left, own, disabled, step, why: reason }: EndProps) {
  const act = useApp((s) => s.act);
  const showWhy = useApp((s) => s.showWhy);
  const ask = useApp((s) => !s.settings.no_end_turn_ask) && !idle;
  const [pressed, setPressed] = useState(-1);
  // asking is the step the question was asked at; a new update drops it.
  const [asking, setAsking] = useState(-1);
  const lit = own && intent !== undefined && !disabled;
  const open = lit && asking === step;
  const end = () => {
    if (lit) {
      setAsking(-1);
      setPressed(step);
      act(intent);
    }
  };
  const press = () => {
    if (lit && ask && !open) {
      setAsking(step);
      return;
    }
    end();
  };
  useKey('e', press);
  const label = own ? 'End turn' : 'Their turn';
  const state = !own ? 'theirs' : lit && idle ? 'own idle' : 'own';
  const why = own && !lit && !disabled ? reason : undefined;
  return (
    <>
      <button
        type="button"
        className={`end-turn ${state}${pressed === step ? ' down' : ''}`}
        title={own ? `End turn (E)${why ? ` · ${why}` : ''}` : 'Their turn'}
        aria-disabled={why ? true : undefined}
        disabled={!lit && !why}
        onClick={(e) => {
          if (why) {
            showWhy(why, e.clientX, e.clientY);
            return;
          }
          press();
        }}
      >
        {label}
      </button>
      {open &&
        createPortal(
          <EndTurnQuestion
            left={left ?? []}
            end={end}
            cancel={() => {
              setAsking(-1);
            }}
          />,
          document.body,
        )}
    </>
  );
}

// EndTurnQuestion asks before a turn ends while plays are left (V-46).
// E confirms, as does Enter on its focused End turn; Esc or a click
// outside keeps playing.
function EndTurnQuestion({ left, end, cancel }: { left: string[]; end: () => void; cancel: () => void }) {
  return (
    <div
      className="overlay end-turn-question"
      role="dialog"
      aria-modal="true"
      aria-label="End turn?"
      onClick={(e) => {
        if (e.target === e.currentTarget) {
          cancel();
        }
      }}
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          e.preventDefault();
          cancel();
        } else if (e.key.toLowerCase() === 'e' && !e.repeat && !e.ctrlKey && !e.metaKey && !e.altKey) {
          e.preventDefault();
          end();
        }
      }}
    >
      <div className="overlay-box">
        <h2>End turn?</h2>
        {left.length > 0 ? (
          <>
            <p className="muted">You can still:</p>
            <ul className="end-turn-left">
              {left.map((l) => (
                <li key={l}>{l}</li>
              ))}
            </ul>
          </>
        ) : (
          <p className="muted">You can still play.</p>
        )}
        <div className="overlay-buttons">
          <button type="button" className="primary" autoFocus onClick={end}>
            End turn
          </button>
          <button type="button" onClick={cancel}>
            Keep playing
          </button>
        </div>
      </div>
    </div>
  );
}

// SwordIcon is the game's pixel sword (15.24), or a drawn one without
// the private images or when the image fails to load.
function SwordIcon() {
  const image = useApp((s) => s.icons.sword);
  const [broken, setBroken] = useState('');
  if (image !== '' && image !== broken) {
    return (
      <img
        className="turn-action-image"
        src={image}
        alt=""
        aria-hidden="true"
        draggable={false}
        onError={() => {
          setBroken(image);
        }}
      />
    );
  }
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M19.5 3 21 4.5 11 14.5 9.5 13Z" />
      <path d="M6.5 13.5 10.5 17.5 9 19 5 15Z" />
      <path d="M7 17 4.2 19.8a1 1 0 0 0 1.4 1.4L8.4 18.4Z" />
    </svg>
  );
}

function BagIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M5 8h14l-1.2 12.1a1 1 0 0 1-1 .9H7.2a1 1 0 0 1-1-.9Z" />
      <path d="M9 9V7a3 3 0 0 1 6 0v2" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  );
}
