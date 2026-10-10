import type { Action } from '../actions';
import { useApp } from '../store';
import './Actions.css';

// ActionBar shows whose move it is, the actions that need no card, and
// why the last action was refused (7.7).
export function ActionBar({ status, buttons, disabled }: { status: string; buttons: Action[]; disabled: boolean }) {
  const act = useApp((s) => s.act);
  const error = useApp((s) => s.error);
  return (
    <div className="action-bar" role="region" aria-label="Actions">
      <div className="action-status">{status}</div>
      {buttons.length > 0 && (
        <div className="action-buttons">
          {buttons.map((b) => (
            <button
              key={b.label}
              type="button"
              disabled={disabled}
              onClick={() => {
                act(b.intent);
              }}
            >
              {b.label}
            </button>
          ))}
        </div>
      )}
      <div className="action-error" role="status">
        {error}
      </div>
    </div>
  );
}
