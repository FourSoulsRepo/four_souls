import type { GameView } from '../protocol';
import { stackLine, type Names } from '../stack';
import { Card } from './Card';
import './Actions.css';

// StackPanel lists what waits on the stack, the top (next to resolve)
// first (7.8).
export function StackPanel({ view, names }: { view: GameView; names: Names }) {
  const stack = view.stack ?? [];
  if (stack.length === 0) {
    return null;
  }
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
    </aside>
  );
}
