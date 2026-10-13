import { useEffect, useState } from 'react';
import { getCardSets } from '../bridge';
import type { Options } from '../protocol';
import { baseGame, setName, sortSets } from '../sets';
import { useApp } from '../store';
import './screens.css';

// The detailed setup's starting values: the same as the simple setup.
const simple: Options = { picking: 'random', draft_size: 3, ban_rounds: 0, ban_timer: 0, response_timer: 0 };

// Host: a simple setup, with "Detailed" for the rest (GS-10). It starts
// from the last hosted setup; the first time every card set is ticked.
export function HostScreen() {
  const settings = useApp((s) => s.settings);
  const saveSettings = useApp((s) => s.saveSettings);
  const last = settings.host;
  const error = useApp((s) => s.error);
  const host = useApp((s) => s.host);
  const hostSaved = useApp((s) => s.hostSaved);
  const saveNickname = useApp((s) => s.saveNickname);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [seats, setSeats] = useState(last.seats || 2);
  const [detailed, setDetailed] = useState(last.detailed);
  const [options, setOptions] = useState<Options>({
    ...simple,
    response_timer: last.response_timer,
    no_bonus_souls: last.no_bonus_souls,
  });
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);
  // available are the sets this app has; chosen always holds the Base Game.
  const [available, setAvailable] = useState<string[]>([baseGame]);
  const [chosen, setChosen] = useState<string[]>([baseGame]);

  useEffect(() => {
    void getCardSets().then((sets) => {
      const all = sortSets(sets.includes(baseGame) ? sets : [baseGame, ...sets]);
      setAvailable(all);
      // The core sets are on by default; a remembered choice wins.
      const kept = last.sets === null ? all : all.filter((c) => last.sets?.includes(c));
      setChosen(kept.includes(baseGame) ? kept : [baseGame, ...kept]);
    });
  }, [last.sets]);

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
              // Remember the setup for the next host; never the password.
              const s = useApp.getState().settings;
              await saveSettings({
                ...s,
                host: {
                  seats,
                  sets: chosen,
                  detailed,
                  response_timer: options.response_timer ?? 0,
                  no_bonus_souls: options.no_bonus_souls ?? false,
                },
              });
              await host({ seats, sets: chosen, options: detailed ? options : undefined, password: password || undefined });
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
        <fieldset className="sets">
          <legend>Card sets</legend>
          {available.map((code) => (
            <label key={code} className="check">
              <input
                type="checkbox"
                checked={chosen.includes(code)}
                disabled={code === baseGame}
                onChange={(e) => {
                  setChosen(e.target.checked ? sortSets([...chosen, code]) : chosen.filter((c) => c !== code));
                }}
              />
              {setName(code)}
            </label>
          ))}
        </fieldset>
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
