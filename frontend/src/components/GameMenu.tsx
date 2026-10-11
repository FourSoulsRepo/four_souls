import { useState } from 'react';
import { useApp } from '../store';
import './Overlays.css';

// GameMenu is the in-game menu: the host may save the game to continue
// later (N-11); anyone may leave.
export function GameMenu({ host }: { host: boolean }) {
  const saveGame = useApp((s) => s.saveGame);
  const quit = useApp((s) => s.quit);
  const [open, setOpen] = useState(false);
  if (!open) {
    return (
      <button
        type="button"
        className="game-menu-button"
        onClick={() => {
          setOpen(true);
        }}
      >
        Menu
      </button>
    );
  }
  return (
    <div className="overlay" role="dialog" aria-label="Game menu">
      <div className="overlay-box">
        <h2>Menu</h2>
        {host && <p className="muted">Saving stops the game for everyone; continue it later from Host a game.</p>}
        <div className="overlay-buttons">
          {host && (
            <button
              type="button"
              className="primary"
              onClick={() => {
                setOpen(false);
                saveGame();
              }}
            >
              Save and leave
            </button>
          )}
          <button
            type="button"
            onClick={() => {
              void quit();
            }}
          >
            Leave game
          </button>
          <button
            type="button"
            onClick={() => {
              setOpen(false);
            }}
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}
