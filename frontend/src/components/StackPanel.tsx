import type { GameView } from '../protocol';
import { stackLine, stackWait, type Names } from '../stack';
import { Card } from './Card';
import './Actions.css';

// StackPanel lists what waits on the stack, the top (next to resolve)
// first (7.8), and whom it waits for (V-11).
export function StackPanel({ view, you, names }: { view: GameView; you: number; names: Names }) {
  const stack = view.stack ?? [];
  if (stack.length === 0) {
    return null;
  }
  const wait = stackWait(view, you, names);
  return (
    <aside className="stack" aria-label="Stack">
      <h2>Stack</h2>
      <ol>
        {[...stack].reverse().map((it) => (
          <li key={it.id}>
            {it.card && <Card cardRef={it.card} size="small" />}
            <span>{stackLine(it, view, names)}</span>
          </li>
        ))}
      </ol>
      {wait.prompt && (
        <p className="stack-wait" aria-label="Waiting for">
          {wait.prompt}
        </p>
      )}
      {wait.responding.length > 0 && (
        <p className="stack-wait" aria-label="Waiting for">
          Waiting for {wait.responding.join(', ')}
        </p>
      )}
      {wait.passed.length > 0 && (
        <p className="stack-wait stack-passed" aria-label="Passed">
          Passed: {wait.passed.join(', ')}
        </p>
      )}
    </aside>
  );
}
