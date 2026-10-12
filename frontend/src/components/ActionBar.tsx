import { useEffect, useState } from 'react';
import { secondsLeft } from '../stack';
import { useApp } from '../store';
import './Actions.css';

interface Props {
  status: string;
  // deadline is when the response timer answers, Unix ms; 0: no timer.
  deadline: number;
  // skip: "Skip all" is offered, or on (N-06).
  skip: 'offer' | 'on' | 'none';
}

// ActionBar shows whose move it is and why the game waits, the response
// timer, "Skip all", and why the last action was refused (7.7, 7.8).
// Pass lives in the stack panel (15.21).
export function ActionBar({ status, deadline, skip }: Props) {
  const skipAll = useApp((s) => s.skipAll);
  const error = useApp((s) => s.error);
  return (
    <div className="action-bar" role="region" aria-label="Actions">
      <div className="action-status">
        {status}
        {deadline > 0 && <Timer deadline={deadline} />}
      </div>
      {skip !== 'none' && (
        <div className="action-buttons">
          <button
            type="button"
            className="secondary"
            onClick={() => {
              skipAll(skip === 'offer');
            }}
          >
            {skip === 'offer' ? 'Skip all' : 'Cancel skip'}
          </button>
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
