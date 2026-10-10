import { useApp } from '../store';
import './screens.css';

// Stands in for the table until step 7.5.
export function GameScreen() {
  const quit = useApp((s) => s.quit);
  return (
    <div className="screen">
      <h1>The game has started</h1>
      <p className="muted">The game table comes in step 7.5.</p>
      <button
        type="button"
        onClick={() => {
          void quit();
        }}
      >
        Leave
      </button>
    </div>
  );
}
