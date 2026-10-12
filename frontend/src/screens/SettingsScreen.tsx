import { useState } from 'react';
import type { Settings, Theme } from '../bridge';
import { matChoices, matLook, matOf } from '../mats';
import type { Speed } from '../playback';
import { useApp } from '../store';
import './screens.css';

function speedOf(v: string): Speed {
  return v === 'fast' || v === 'slow' || v === 'off' ? v : 'normal';
}

function themeOf(v: string): Theme {
  return v === 'light' ? 'light' : 'dark';
}

// SettingsScreen: ST-01, ST-02, ST-03, ST-05, ST-06. Saved settings
// apply at once.
export function SettingsScreen() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const save = useApp((s) => s.saveSettings);
  const show = useApp((s) => s.show);
  const mats = useApp((s) => s.mats);
  const [form, setForm] = useState<Settings>(settings);
  const [saved, setSaved] = useState(false);
  const preview = matLook(form.mat, mats);
  const change = (s: Partial<Settings>) => {
    setForm({ ...form, ...s });
    setSaved(false);
  };

  return (
    <div className="screen">
      <h1>Settings</h1>
      <form
        className="panel"
        onSubmit={(e) => {
          e.preventDefault();
          void save(form).then((ok) => {
            setSaved(ok);
            if (ok) {
              setForm(useApp.getState().settings);
            }
          });
        }}
      >
        <label className="field">
          Nickname
          <input
            value={form.nickname}
            onChange={(e) => {
              change({ nickname: e.target.value });
            }}
          />
        </label>
        <label className="check">
          <input
            type="checkbox"
            checked={form.fullscreen}
            onChange={(e) => {
              change({ fullscreen: e.target.checked });
            }}
          />
          Fullscreen
        </label>
        <div className="row">
          <label className="field">
            Theme
            <select
              value={form.theme}
              onChange={(e) => {
                change({ theme: themeOf(e.target.value) });
              }}
            >
              <option value="dark">Dark</option>
              <option value="light">Light</option>
            </select>
          </label>
          <label className="field">
            Game mat
            <select
              value={matOf(form.mat, mats)}
              onChange={(e) => {
                change({ mat: e.target.value });
              }}
            >
              {matChoices(mats).map((m) => (
                <option key={m.id} value={m.id}>
                  {m.name}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            Animations
            <select
              value={form.animation}
              onChange={(e) => {
                change({ animation: speedOf(e.target.value) });
              }}
            >
              <option value="normal">Normal</option>
              <option value="fast">Fast</option>
              <option value="slow">Slow</option>
              <option value="off">Off</option>
            </select>
          </label>
        </div>
        <div className={`mat-preview ${preview.className}`} style={preview.style} aria-label="Game mat preview" />
        <button type="submit" className="primary">
          Save
        </button>
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
