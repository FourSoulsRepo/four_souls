import { useEffect, useState, type ReactNode } from 'react';
import { IntentKind, PromptKind, Purpose, type TurnActions as Allowed } from '../actions';
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
}

// TurnActions are buy and attack as icon buttons by own name (15.19,
// F-11): lit when the engine allows them, dim when not (F-12). A press
// shows at once, before the server answers (V-15). While the player
// picks what to buy or attack, that icon stays pressed (F-13).
export function TurnActions({ allowed, disabled, step, waiting, you }: Props) {
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
    },
    {
      kind: IntentKind.attack,
      label: 'Attack',
      key: 'a',
      icon: <SwordIcon />,
      intent: allowed.attack,
      choosing: picking === Purpose.attackTarget,
      choosingLabel: 'Choosing what to attack',
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
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.repeat || e.ctrlKey || e.metaKey || e.altKey || typing(e.target)) {
        return;
      }
      const t = tools.find((x) => x.key === e.key.toLowerCase());
      if (t) {
        press(t);
      }
    };
    window.addEventListener('keydown', onKey);
    return () => {
      window.removeEventListener('keydown', onKey);
    };
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
          </button>
        );
      })}
    </div>
  );
}

// typing: the key goes to a text field, not to the table.
function typing(target: EventTarget | null): boolean {
  return target instanceof HTMLElement && (target.isContentEditable || /^(INPUT|TEXTAREA|SELECT)$/.test(target.tagName));
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
