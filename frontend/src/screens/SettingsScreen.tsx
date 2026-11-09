import { useState } from 'react';
import type { Settings } from '../bridge';
import { Die } from '../components/Die';
import { diceChoices } from '../dice';
import { shownLanguage, t } from '../i18n';
import { languageName } from '../locale';
import { matChoices, matLook, matOf } from '../mats';
import type { Speed } from '../playback';
import { useApp } from '../store';
import { maxScale, minScale } from '../uiScale';
import './screens.css';

function speedOf(v: string): Speed {
  return v === 'fast' || v === 'slow' || v === 'off' ? v : 'normal';
}

// SettingsScreen: ST-01, ST-03, ST-05, ST-06, ST-09, ST-10, the dice look
// (V-19, 15.64), the question before ending the turn (V-46), and the
// language (ADR 021). Saved settings apply at once; the dice look goes
// out in the next hello.
export function SettingsScreen() {
  const show = useApp((s) => s.show);
  return (
    <div className="screen">
      <h1>{t('settings.title')}</h1>
      <SettingsForm
        inGame={false}
        onBack={() => {
          show('menu');
        }}
      />
    </div>
  );
}

// SettingsForm is the settings form of the main screen and of the game
// menu (ST-11, 18.19). In a game it leaves out what the running game
// would not take: the nickname and dice go out in the hello, and the
// language stays with the main screen.
export function SettingsForm({ inGame, onBack }: { inGame: boolean; onBack: () => void }) {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const save = useApp((s) => s.saveSettings);
  const mats = useApp((s) => s.mats);
  // Card languages join the list when the cards load (16.10).
  const languages = useApp((s) => s.languages);
  const dice = useApp((s) => s.dice);
  const diceGroups = diceChoices(dice);
  const showDice = !inGame && diceGroups.length > 0;
  const [form, setForm] = useState<Settings>(settings);
  const [saved, setSaved] = useState(false);
  const preview = matLook(form.mat, mats);
  const diceImage = dice.find((d) => d.id === form.dice)?.image ?? '';
  const change = (s: Partial<Settings>) => {
    setForm({ ...form, ...s });
    setSaved(false);
  };

  return (
    <>
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
        {!inGame && (
          <label className="field">
            {t('settings.nickname')}
            <input
              value={form.nickname}
              onChange={(e) => {
                change({ nickname: e.target.value });
              }}
            />
          </label>
        )}
        <label className="check">
          <input
            type="checkbox"
            checked={form.fullscreen}
            onChange={(e) => {
              change({ fullscreen: e.target.checked });
            }}
          />
          {t('settings.fullscreen')}
        </label>
        <label className="check">
          <input
            type="checkbox"
            checked={!form.no_end_turn_ask}
            onChange={(e) => {
              change({ no_end_turn_ask: !e.target.checked });
            }}
          />
          {t('settings.end_turn_ask')}
        </label>
        {/* Each language in its own name; the list is not translated. */}
        <div className="row">
          {!inGame && (
            <label className="field">
              {t('settings.language')}
              <select
                value={shownLanguage(form.language)}
                onChange={(e) => {
                  change({ language: e.target.value });
                }}
              >
                {languages.map((l) => (
                  <option key={l} value={l} lang={l}>
                    {languageName(l)}
                  </option>
                ))}
              </select>
            </label>
          )}
          <label className="field">
            {t('settings.animations')}
            <select
              value={form.animation}
              onChange={(e) => {
                change({ animation: speedOf(e.target.value) });
              }}
            >
              <option value="normal">{t('settings.speed.normal')}</option>
              <option value="fast">{t('settings.speed.fast')}</option>
              <option value="slow">{t('settings.speed.slow')}</option>
              <option value="off">{t('settings.speed.off')}</option>
            </select>
          </label>
        </div>
        <label className="field">
          {t('settings.volume')}
          <input
            type="range"
            min={0}
            max={100}
            step={5}
            value={form.volume}
            onChange={(e) => {
              change({ volume: Number(e.target.value) });
            }}
          />
          <span>{form.volume}%</span>
        </label>
        <label className="field">
          {t('settings.scale')}
          <input
            type="range"
            min={minScale}
            max={maxScale}
            step={5}
            value={form.scale}
            onChange={(e) => {
              change({ scale: Number(e.target.value) });
            }}
          />
          <span>{form.scale}%</span>
        </label>
        {/* The looks sit on one line, above their previews (15.80). */}
        <div className="row">
          <label className="field">
            {t('settings.mat')}
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
          {showDice && (
            <label className="field">
              {t('settings.dice')}
              <select
                value={form.dice}
                onChange={(e) => {
                  change({ dice: e.target.value });
                }}
              >
                {diceGroups.map((g) => (
                  <optgroup key={g.set} label={diceSet(g.set)}>
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
        <div className={`mat-preview ${preview.className}`} style={preview.style} aria-label={t('settings.mat_preview')} />
        {showDice && (
          <div className="dice-preview" role="group" aria-label={t('settings.dice_preview')}>
            {[1, 2, 3, 4, 5, 6].map((v) => (
              <Die key={v} value={v} image={diceImage} />
            ))}
          </div>
        )}
        <button type="submit" className="primary">
          {t('settings.save')}
        </button>
      </form>
      <p className="status" role="status">
        {error || (saved ? t('settings.saved') : '')}
      </p>
      <button type="button" onClick={onBack}>
        {t('settings.back')}
      </button>
    </>
  );
}

// diceSet names a group of dice looks; Go sends Holy or Unholy. The
// looks keep their item names, as cards do.
function diceSet(set: string): string {
  if (set === 'Holy') {
    return t('settings.dice_set.holy');
  }
  return set === 'Unholy' ? t('settings.dice_set.unholy') : set;
}
