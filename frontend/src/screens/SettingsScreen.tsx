import { useState } from 'react';
import { useApp } from '../store';
import './SettingsScreen.css';

export function SettingsScreen() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const save = useApp((s) => s.saveSettings);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [saved, setSaved] = useState(false);

  return (
    <div className="settings">
      <h1>Settings</h1>
      <form
        className="settings-form"
        onSubmit={(e) => {
          e.preventDefault();
          void save({ ...settings, nickname }).then((ok) => {
            setSaved(ok);
            if (ok) {
              setNickname(useApp.getState().settings.nickname);
            }
          });
        }}
      >
        <label>
          Nickname
          <input
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
              setSaved(false);
            }}
          />
        </label>
        <button type="submit">Save</button>
      </form>
      <p className="settings-status" role="status">
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
