import { useState } from 'react';
import { t } from '../i18n';
import { useApp } from '../store';
import './Overlays.css';

// GameMenu is the in-game menu: the host may save the game to continue
// later (N-11), or end it for everyone after a confirm step (19.84);
// anyone may leave.
export function GameMenu({ host }: { host: boolean }) {
  const saveGame = useApp((s) => s.saveGame);
  const endGame = useApp((s) => s.endGame);
  const quit = useApp((s) => s.quit);
  const [open, setOpen] = useState(false);
  const [ending, setEnding] = useState(false);
  if (!open) {
    return (
      <button
        type="button"
        className="game-menu-button"
        onClick={() => {
          setOpen(true);
        }}
      >
        {t('game_menu.button')}
      </button>
    );
  }
  if (ending) {
    return (
      <div className="overlay" role="dialog" aria-label={t('game_menu.label')}>
        <div className="overlay-box">
          <h2>{t('game_menu.end_title')}</h2>
          <p className="muted">{t('game_menu.end_hint')}</p>
          <div className="overlay-buttons">
            <button
              type="button"
              className="primary"
              onClick={() => {
                setOpen(false);
                setEnding(false);
                endGame();
              }}
            >
              {t('game_menu.end')}
            </button>
            <button
              type="button"
              onClick={() => {
                setEnding(false);
              }}
            >
              {t('game_menu.back')}
            </button>
          </div>
        </div>
      </div>
    );
  }
  return (
    <div className="overlay" role="dialog" aria-label={t('game_menu.label')}>
      <div className="overlay-box">
        <h2>{t('game_menu.title')}</h2>
        {host && <p className="muted">{t('game_menu.save_hint')}</p>}
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
              {t('game_menu.save')}
            </button>
          )}
          {host && (
            <button
              type="button"
              onClick={() => {
                setEnding(true);
              }}
            >
              {t('game_menu.end')}
            </button>
          )}
          <button
            type="button"
            onClick={() => {
              void quit();
            }}
          >
            {t('game_menu.leave')}
          </button>
          <button
            type="button"
            onClick={() => {
              setOpen(false);
            }}
          >
            {t('game_menu.cancel')}
          </button>
        </div>
      </div>
    </div>
  );
}
