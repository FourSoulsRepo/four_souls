import { useState } from 'react';
import type { Speed } from '../playback';
import { useApp } from '../store';
import './screens.css';

function speedOf(v: string): Speed {
  return v === 'fast' || v === 'slow' || v === 'off' ? v : 'normal';
}

export function SettingsScreen() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const save = useApp((s) => s.saveSettings);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [animation, setAnimation] = useState(settings.animation);
  const [saved, setSaved] = useState(false);

  return (
    <div className="screen">
      <h1>Settings</h1>
      <form
        className="row"
        onSubmit={(e) => {
          e.preventDefault();
          void save({ ...settings, nickname, animation }).then((ok) => {
            setSaved(ok);
            if (ok) {
              setNickname(useApp.getState().settings.nickname);
            }
          });
        }}
      >
        <label className="field">
          Nickname
          <input
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
              setSaved(false);
            }}
          />
        </label>
        <label className="field">
          Animations
          <select
            value={animation}
            onChange={(e) => {
              setAnimation(speedOf(e.target.value));
              setSaved(false);
            }}
          >
            <option value="normal">Normal</option>
            <option value="fast">Fast</option>
            <option value="slow">Slow</option>
            <option value="off">Off</option>
          </select>
        </label>
        <button type="submit">Save</button>
      </form>
      <p className="status" role="status">
        {error || (saved ? 'Saved.' : '')}
      </p>
      <button
        type="button"
        onClick={() => {
          show('menu');
        }}
      >
        Back
      </button>
    </div>
  );
}
