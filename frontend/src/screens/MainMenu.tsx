import { useEffect, useState } from 'react';
import { getVersions, type Versions } from '../bridge';
import './MainMenu.css';

export function MainMenu() {
  const [versions, setVersions] = useState<Versions | null>(null);

  useEffect(() => {
    getVersions().then(setVersions, () => {
      setVersions(null);
    });
  }, []);

  return (
    <div className="menu">
      <h1>Four Souls</h1>
      <nav className="menu-buttons">
        <button type="button" disabled>Host a game</button>
        <button type="button" disabled>Join a game</button>
        <button type="button" disabled>Records</button>
        <button type="button" disabled>Settings</button>
      </nav>
      {versions && (
        <div className="menu-versions">
          app {versions.app} · rules engine {versions.engine}
        </div>
      )}
    </div>
  );
}
