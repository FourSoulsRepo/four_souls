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
  const saveNickname = useApp((s) => s.saveNickname);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [seats, setSeats] = useState(2);
  const [detailed, setDetailed] = useState(false);
  const [options, setOptions] = useState<Options>(simple);
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
              await host({ seats, options: detailed ? options : undefined });
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
            <div className="row">
              <label className="field">
                Characters
                <select
                  value={options.picking}
                  onChange={(e) => {
                    option({ picking: e.target.value === 'draft' ? 'draft' : 'random' });
                  }}
                >
                  <option value="random">Random</option>
                  <option value="draft">Draft</option>
                </select>
              </label>
              {options.picking === 'draft' && (
                <label className="field">
                  Offered to each
                  <select
                    value={options.draft_size}
                    onChange={(e) => {
                      option({ draft_size: Number(e.target.value) });
                    }}
                  >
                    {[2, 3, 4, 5].map((n) => (
                      <option key={n} value={n}>
                        {n}
                      </option>
                    ))}
                  </select>
                </label>
              )}
            </div>
            <div className="row">
              <label className="field">
                Ban rounds
                <select
                  value={options.ban_rounds}
                  onChange={(e) => {
                    option({ ban_rounds: Number(e.target.value) });
                  }}
                >
                  {[0, 1, 2, 3].map((n) => (
                    <option key={n} value={n}>
                      {n}
                    </option>
                  ))}
                </select>
              </label>
              <label className="field">
                Ban timer
                <select
                  value={options.ban_timer}
                  onChange={(e) => {
                    option({ ban_timer: Number(e.target.value) });
                  }}
                >
                  <option value={0}>Off</option>
                  <option value={15}>15 s</option>
                  <option value={30}>30 s</option>
                  <option value={60}>60 s</option>
                </select>
              </label>
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
