import { useState } from 'react';
import type { Settings } from '../bridge';
import { Die } from '../components/Die';
import { diceChoices } from '../dice';
import { matChoices, matLook, matOf } from '../mats';
import type { Speed } from '../playback';
import { useApp } from '../store';
import './screens.css';

function speedOf(v: string): Speed {
  return v === 'fast' || v === 'slow' || v === 'off' ? v : 'normal';
}

// SettingsScreen: ST-01, ST-03, ST-05, ST-06, the dice look
// (V-19, 15.64), and the question before ending the turn (V-46). Saved
// settings apply at once; the dice look goes out in the next hello.
export function SettingsScreen() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const save = useApp((s) => s.saveSettings);
  const show = useApp((s) => s.show);
  const mats = useApp((s) => s.mats);
  const dice = useApp((s) => s.dice);
  const diceGroups = diceChoices(dice);
  const [form, setForm] = useState<Settings>(settings);
  const [saved, setSaved] = useState(false);
  const preview = matLook(form.mat, mats);
  const diceImage = dice.find((d) => d.id === form.dice)?.image ?? '';
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
        <label className="check">
          <input
            type="checkbox"
            checked={!form.no_end_turn_ask}
            onChange={(e) => {
              change({ no_end_turn_ask: !e.target.checked });
            }}
          />
          Ask before ending the turn while plays are left
        </label>
        <div className="row">
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
          {diceGroups.length > 0 && (
            <label className="field">
              Dice
              <select
                value={form.dice}
                onChange={(e) => {
                  change({ dice: e.target.value });
                }}
              >
                {diceGroups.map((g) => (
                  <optgroup key={g.set} label={g.set}>
                    {g.looks.map((d) => (
                      <option key={d.id} value={d.id}>
                        {d.name}
                      </option>
                    ))}
                  </optgroup>
                ))}
              </select>
            </label>
          )}
        </div>
        <div className={`mat-preview ${preview.className}`} style={preview.style} aria-label="Game mat preview" />
        {diceGroups.length > 0 && (
          <div className="dice-preview" role="group" aria-label="Dice preview">
            {[1, 2, 3, 4, 5, 6].map((v) => (
              <Die key={v} value={v} image={diceImage} />
            ))}
          </div>
        )}
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
