import { useEffect, useState } from 'react';
import { getVersions, quitApp, type Versions } from '../bridge';
import { useApp } from '../store';
import './MainMenu.css';

export function MainMenu() {
  const [versions, setVersions] = useState<Versions | null>(null);
  const show = useApp((s) => s.show);

  useEffect(() => {
    getVersions().then(setVersions, () => {
      setVersions(null);
    });
  }, []);

  return (
    <div className="menu">
      <Title />
      <nav className="menu-buttons">
        <button
          type="button"
          onClick={() => {
            show('host');
          }}
        >
          Host a game
        </button>
        <button
          type="button"
          onClick={() => {
            show('join');
          }}
        >
          Join a game
        </button>
        <button type="button" disabled>Records</button>
        <button
          type="button"
          onClick={() => {
            show('settings');
          }}
        >
          Settings
        </button>
        <button type="button" onClick={quitApp}>
          Exit
        </button>
      </nav>
      {versions && (
        <div className="menu-versions">
          app {versions.app} · rules engine {versions.engine}
        </div>
      )}
    </div>
  );
}

// Title is the foursouls.ru banner (15.41), or the text title without the
// private images or when the banner fails to load.
function Title() {
  const image = useApp((s) => s.icons.banner);
  const [broken, setBroken] = useState('');
  if (image === '' || image === broken) {
    return <h1>Four Souls</h1>;
  }
  return (
    <h1 className="menu-banner">
      <img
        src={image}
        alt="Four Souls"
        draggable={false}
        onError={() => {
          setBroken(image);
        }}
      />
    </h1>
  );
}
