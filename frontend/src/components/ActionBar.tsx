import { useEffect, useState } from 'react';
import { IntentKind, type Action } from '../actions';
import { useKey } from '../keys';
import { secondsLeft } from '../stack';
import { useApp } from '../store';
import './Actions.css';

interface Props {
  status: string;
  buttons: Action[];
  disabled: boolean;
  // deadline is when the response timer answers, Unix ms; 0: no timer.
  deadline: number;
  // skip: "Skip all" is offered, or on (N-06).
  skip: 'offer' | 'on' | 'none';
}

// ActionBar shows whose move it is and why the game waits, the actions
// that need no card, the response timer, and why the last action was
// refused (7.7, 7.8). Space passes (F-15).
export function ActionBar({ status, buttons, disabled, deadline, skip }: Props) {
  const act = useApp((s) => s.act);
  const skipAll = useApp((s) => s.skipAll);
  const error = useApp((s) => s.error);
  const pass = buttons.find((b) => b.intent.kind === IntentKind.pass);
  useKey(' ', () => {
    if (pass && !disabled) {
      act(pass.intent);
    }
  });
  return (
    <div className="action-bar" role="region" aria-label="Actions">
      <div className="action-status">
        {status}
        {deadline > 0 && <Timer deadline={deadline} />}
      </div>
      {(buttons.length > 0 || skip !== 'none') && (
        <div className="action-buttons">
          {buttons.map((b) => (
            <button
              key={b.label}
              type="button"
              title={b.intent.kind === IntentKind.pass ? `${b.label} (Space)` : undefined}
              disabled={disabled}
              onClick={() => {
                act(b.intent);
              }}
            >
              {b.label}
            </button>
          ))}
          {skip !== 'none' && (
            <button
              type="button"
              className="secondary"
              onClick={() => {
                skipAll(skip === 'offer');
              }}
            >
              {skip === 'offer' ? 'Skip all' : 'Cancel skip'}
            </button>
          )}
        </div>
      )}
      <div className="action-error" role="status">
        {error}
      </div>
    </div>
  );
}

// Timer counts down the response timer (N-07).
function Timer({ deadline }: { deadline: number }) {
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
    <span className="action-timer" aria-label="Time left">
      {' '}
      · {secondsLeft(deadline, now)} s
    </span>
  );
}
