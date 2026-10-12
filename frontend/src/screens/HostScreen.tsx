import { useState } from 'react';
import type { Options } from '../protocol';
import { useApp } from '../store';
import './screens.css';

// The detailed setup's starting values: the same as the simple setup.
const simple: Options = { picking: 'random', draft_size: 3, ban_rounds: 0, ban_timer: 0, response_timer: 0 };

// Host: a simple setup, with "Detailed" for the rest (GS-10).
export function HostScreen() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const host = useApp((s) => s.host);
  const hostSaved = useApp((s) => s.hostSaved);
  const saveNickname = useApp((s) => s.saveNickname);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [seats, setSeats] = useState(2);
  const [detailed, setDetailed] = useState(false);
  const [options, setOptions] = useState<Options>(simple);
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);

  const option = (o: Options) => {
    setOptions({ ...options, ...o });
  };

  return (
    <div className="screen">
      <h1>Host a game</h1>
      <form
        className="panel"
        onSubmit={(e) => {
          e.preventDefault();
          setBusy(true);
          void (async () => {
            if (await saveNickname(nickname)) {
              await host({ seats, options: detailed ? options : undefined, password: password || undefined });
            }
            setBusy(false);
          })();
        }}
      >
        <label className="field">
          Nickname
          <input
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
            }}
          />
        </label>
        <label className="field">
          Players
          <select
            value={seats}
            onChange={(e) => {
              setSeats(Number(e.target.value));
            }}
          >
            <option value={2}>2</option>
            <option value={3}>3</option>
            <option value={4}>4</option>
          </select>
        </label>
        <label className="field">
          Password (optional)
          <input
            type="password"
            autoComplete="off"
            maxLength={64}
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
          />
        </label>
        <label className="check">
          <input
            type="checkbox"
            checked={detailed}
            onChange={(e) => {
              setDetailed(e.target.checked);
            }}
          />
          Detailed setup
        </label>
        {detailed && (
          <>
            {/* Draft and bans need the ban and pick screen, which comes with
                the tournament tools (roadmap, Future); until then a game
                with them could not start. */}
            <div className="row">
              <label className="field">
                Response timer
                <select
                  value={options.response_timer}
                  onChange={(e) => {
                    option({ response_timer: Number(e.target.value) });
                  }}
                >
                  <option value={0}>Off</option>
                  <option value={60}>60 s</option>
                  <option value={90}>90 s</option>
                  <option value={120}>120 s</option>
                  <option value={180}>180 s</option>
                </select>
              </label>
            </div>
            <label className="check">
              <input
                type="checkbox"
                checked={!options.no_bonus_souls}
                onChange={(e) => {
                  option({ no_bonus_souls: !e.target.checked });
                }}
              />
              Bonus souls
            </label>
          </>
        )}
        <button type="submit" className="primary" disabled={busy}>
          Host
        </button>
        <button
          type="button"
          disabled={busy}
          onClick={() => {
            setBusy(true);
            void (async () => {
              if (await saveNickname(nickname)) {
                await hostSaved();
              }
              setBusy(false);
            })();
          }}
        >
          Continue a saved game
        </button>
      </form>
      <p className="status" role="status">
        {error}
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
