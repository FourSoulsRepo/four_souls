import { useState, type ReactNode } from 'react';
import { IntentKind, PromptKind, Purpose, type TurnActions as Allowed } from '../actions';
import { useKey } from '../keys';
import type { Intent, Prompt } from '../protocol';
import { useApp } from '../store';
import './Actions.css';

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
  // choosing: the player is picking what this action is for.
  choosing: boolean;
  choosingLabel: string;
  left?: number;
}

// TurnActions are buy and attack as icon buttons by own name (15.19,
// F-11): lit when the engine allows them, dim when not (F-12). A press
// shows at once, before the server answers (V-15). While the player
// picks what to buy or attack, that icon stays pressed (F-13). On own
// turn a badge on each icon counts what is left (15.24).
export function TurnActions({ allowed, disabled, step, waiting, you, left }: Props) {
  const act = useApp((s) => s.act);
  const [pressed, setPressed] = useState<{ kind: number; step: number } | null>(null);
  const picking = waiting.kind === PromptKind.choose && waiting.player === you ? (waiting.purpose ?? -1) : -1;
  const tools: Tool[] = [
    {
      kind: IntentKind.purchase,
      label: 'Buy',
      key: 'b',
      icon: <BagIcon />,
      intent: allowed.buy,
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
        return (
          <button
            key={t.kind}
            type="button"
            className={`turn-action${lit ? ' lit' : ''}${down ? ' down' : ''}`}
            aria-label={t.choosing ? t.choosingLabel : t.label}
            aria-pressed={down}
            title={`${t.choosing ? t.choosingLabel : t.label} (${t.key.toUpperCase()})`}
            disabled={!lit}
            onClick={() => {
              press(t);
            }}
          >
            {t.icon}
            {t.left !== undefined && (
              <span className={`turn-action-left${t.left === 0 ? ' none' : ''}`} aria-label={`${String(t.left)} left`}>
                {t.left}
              </span>
            )}
          </button>
        );
      })}
    </div>
  );
}

interface EndProps {
  // intent ends the turn, when the engine allows it now.
  intent?: Intent;
  // idle: nothing is left to do but end the turn (H-20).
  idle: boolean;
  // own: it is the viewer's turn.
  own: boolean;
  disabled: boolean;
  step: number;
}

// EndTurn is the one big End turn button at the board's right edge
// (15.20, H-04): yellow on own turn, grey "Their turn" otherwise
// (H-19), glowing green once nothing is left to do (H-20). E presses
// it (F-15).
export function EndTurn({ intent, idle, own, disabled, step }: EndProps) {
  const act = useApp((s) => s.act);
  const [pressed, setPressed] = useState(-1);
  const lit = own && intent !== undefined && !disabled;
  const press = () => {
    if (lit) {
      setPressed(step);
      act(intent);
    }
  };
  useKey('e', press);
  const label = own ? 'End turn' : 'Their turn';
  const state = !own ? 'theirs' : lit && idle ? 'own idle' : 'own';
  return (
    <button
      type="button"
      className={`end-turn ${state}${pressed === step ? ' down' : ''}`}
      title={own ? 'End turn (E)' : 'Their turn'}
      disabled={!lit}
      onClick={press}
    >
      {label}
    </button>
  );
}

function SwordIcon() {
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
