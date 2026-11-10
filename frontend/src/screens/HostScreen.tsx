import { useEffect, useState } from 'react';
import { deckRatiosBlocked, getCardSets, getChallenges } from '../bridge';
import { challengeSeats, createOptions, pickChallenge, type ChallengeChoice } from '../hostSetup';
import { t } from '../i18n';
import { useEscape } from '../keys';
import { challengeName, difficultyName } from '../lobby';
import type { Options, RollbackMode } from '../protocol';
import { baseGame, requiem, setName, sortSets } from '../sets';
import { deckRatiosText, useApp } from '../store';
import './screens.css';

// rollbackHint says what each rollback mode allows (ADR 020, item 1).
function rollbackHint(mode: '' | RollbackMode): string {
  return mode === 'soft' ? t('host.rollback_hint.soft') : mode === 'god' ? t('host.rollback_hint.god') : t('host.rollback_hint.off');
}

// The detailed setup's starting values: the same as the simple setup.
const simple: Options = { picking: 'random', draft_size: 3, ban_rounds: 0, ban_timer: 0, response_timer: 0 };

// Host: a simple setup, with "Detailed" for the rest (GS-10). It starts
// from the last hosted setup; the first time every card set is ticked.
// LAN's Host a game hosts in this app; Create a game in a server's
// Games list creates it on that server (N-17, 15.98). The game mode is
// a regular game or a co-op challenge (ADR 028): a challenge has 2
// seats and no bonus souls.
export function HostScreen() {
  const settings = useApp((s) => s.settings);
  const saveSettings = useApp((s) => s.saveSettings);
  const last = settings.host;
  const error = useApp((s) => s.error);
  const host = useApp((s) => s.host);
  const hostSaved = useApp((s) => s.hostSaved);
  const createGame = useApp((s) => s.createGame);
  const remote = useApp((s) => s.connected && !s.hosting);
  const address = useApp((s) => s.address);
  const saveNickname = useApp((s) => s.saveNickname);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [seats, setSeats] = useState(last.seats || 2);
  const [detailed, setDetailed] = useState(last.detailed);
  // The game mode, and the challenges this app knows (ADR 028).
  const [challengeMode, setChallengeMode] = useState(last.mode === 'challenge');
  const [challenges, setChallenges] = useState<ChallengeChoice[]>([]);
  const [challenge, setChallenge] = useState(last.challenge);
  const [difficulty, setDifficulty] = useState(last.difficulty);
  const [options, setOptions] = useState<Options>({
    ...simple,
    // Random stays the default; a draft is remembered (GS-01, 14.11).
    picking: last.picking === 'draft' ? 'draft' : 'random',
    draft_size: last.draft_size || 3,
    response_timer: last.response_timer,
    no_bonus_souls: last.no_bonus_souls,
    no_deck_ratios: last.no_deck_ratios,
    room_deck: last.room_deck,
    random_seats: last.random_seats,
    rollback: last.rollback === 'soft' || last.rollback === 'god' ? last.rollback : undefined,
  });
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);
  // available are the sets this app has; chosen always holds the Base Game.
  const [available, setAvailable] = useState<string[]>([baseGame]);
  const [chosen, setChosen] = useState<string[]>([baseGame]);
  // blocked says why the deck ratios can't be used; null while unknown.
  const [blocked, setBlocked] = useState<string | null>(null);
  // One level up: the LAN screen, or the server's games.
  useEscape(() => {
    show('join');
  });

  useEffect(() => {
    void getCardSets().then((sets) => {
      const all = sortSets(sets.includes(baseGame) ? sets : [baseGame, ...sets]);
      setAvailable(all);
      // The core sets are on by default; a remembered choice wins.
      const kept = last.sets === null ? all : all.filter((c) => last.sets?.includes(c));
      setChosen(kept.includes(baseGame) ? kept : [baseGame, ...kept]);
    });
  }, [last.sets]);

  useEffect(() => {
    void getChallenges().then(setChallenges);
  }, []);

  // pick is the challenge to send: the chosen one, or the first the app
  // knows; null in a regular game or when the app knows none.
  const known = pickChallenge(challenges, challenge, difficulty);
  const pick = challengeMode ? known : null;
  const playing = pick ? challengeSeats : seats;
  const levels = challenges.find((c) => c.id === known?.challenge)?.difficulties ?? [];

  // The deck ratios need enough cards in each group (GS-07).
  useEffect(() => {
    let live = true;
    void deckRatiosBlocked(chosen, playing).then((why) => {
      if (live) {
        setBlocked(why);
      }
    });
    return () => {
      live = false;
    };
  }, [chosen, playing]);

  // Only Requiem has rooms: the room deck shows with it (14.13).
  const rooms = chosen.includes(requiem);

  const option = (o: Options) => {
    setOptions({ ...options, ...o });
  };

  return (
    <div className="screen">
      <h1>{remote ? t('host.title_create') : t('host.title_host')}</h1>
      {remote && <p className="muted">{t('host.on_server', { address })}</p>}
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
                  no_deck_ratios: options.no_deck_ratios ?? false,
                  room_deck: options.room_deck ?? false,
                  random_seats: options.random_seats ?? false,
                  rollback: options.rollback ?? '',
                  picking: options.picking === 'draft' ? 'draft' : '',
                  draft_size: options.draft_size ?? 3,
                  mode: challengeMode ? 'challenge' : '',
                  challenge: known?.challenge ?? '',
                  difficulty: known?.difficulty ?? '',
                },
              });
              const game = { seats: playing, sets: chosen, options: createOptions(detailed, options, rooms, pick), password: password || undefined };
              if (remote) {
                createGame(game);
              } else {
                await host(game);
              }
            }
            setBusy(false);
          })();
        }}
      >
        {/* On a server the nickname went with hello already. */}
        {!remote && (
          <label className="field">
            {t('lobby.nickname')}
            <input
              value={nickname}
              onChange={(e) => {
                setNickname(e.target.value);
              }}
            />
          </label>
        )}
        {/* A challenge: which one and how hard, beside the mode. The
            aria-labels keep each name exact: "Challenge" is also a mode. */}
        <div className="row">
          <label className="field">
            {t('host.mode')}
            <select
              aria-label={t('host.mode')}
              value={pick ? 'challenge' : 'regular'}
              onChange={(e) => {
                setChallengeMode(e.target.value === 'challenge');
              }}
            >
              <option value="regular">{t('host.mode_regular')}</option>
              <option value="challenge" disabled={known === null}>
                {t('host.mode_challenge')}
              </option>
            </select>
          </label>
          {pick && (
            <>
              <label className="field">
                {t('host.challenge')}
                <select
                  aria-label={t('host.challenge')}
                  value={pick.challenge}
                  onChange={(e) => {
                    setChallenge(e.target.value);
                  }}
                >
                  {challenges.map((c) => (
                    <option key={c.id} value={c.id}>
                      {challengeName(c.id)}
                    </option>
                  ))}
                </select>
              </label>
              <label className="field">
                {t('host.difficulty')}
                <select
                  aria-label={t('host.difficulty')}
                  value={pick.difficulty}
                  onChange={(e) => {
                    setDifficulty(e.target.value);
                  }}
                >
                  {levels.map((d) => (
                    <option key={d} value={d}>
                      {difficultyName(d)}
                    </option>
                  ))}
                </select>
              </label>
            </>
          )}
        </div>
        {/* A challenge is played by exactly 2 (R-COOP, ADR 028). */}
        <label className="field">
          {t('host.players')}
          <select
            value={playing}
            disabled={pick !== null}
            aria-describedby={pick ? 'challenge-seats' : undefined}
            onChange={(e) => {
              setSeats(Number(e.target.value));
            }}
          >
            <option value={2}>2</option>
            <option value={3}>3</option>
            <option value={4}>4</option>
          </select>
        </label>
        {pick && (
          <p id="challenge-seats" className="muted hint">
            {t('host.challenge_seats', { n: challengeSeats })}
          </p>
        )}
        <fieldset className="sets">
          <legend>{t('host.sets')}</legend>
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
          {t('host.password')}
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
          {t('host.detailed')}
        </label>
        {detailed && (
          <>
            {/* One random character each by default; a draft offers each
                player 2–5 and they pick one (GS-01, 14.11). Bans need
                the ban screen, which comes with the tournament tools. */}
            <div className="row">
              <label className="field">
                {t('host.picking')}
                <select
                  value={options.picking}
                  onChange={(e) => {
                    option({ picking: e.target.value === 'draft' ? 'draft' : 'random' });
                  }}
                >
                  <option value="random">{t('host.picking_random')}</option>
                  <option value="draft">{t('host.picking_draft')}</option>
                </select>
              </label>
              {options.picking === 'draft' && (
                <label className="field">
                  {t('host.draft_size')}
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
                {t('host.timer')}
                <select
                  value={options.response_timer}
                  onChange={(e) => {
                    option({ response_timer: Number(e.target.value) });
                  }}
                >
                  <option value={0}>{t('host.timer_off')}</option>
                  <option value={60}>{t('host.seconds', { n: 60 })}</option>
                  <option value={90}>{t('host.seconds', { n: 90 })}</option>
                  <option value={120}>{t('host.seconds', { n: 120 })}</option>
                  <option value={180}>{t('host.seconds', { n: 180 })}</option>
                </select>
              </label>
            </div>
            {/* A challenge plays without them (ADR 028). */}
            {!pick && (
              <label className="check">
                <input
                  type="checkbox"
                  checked={!options.no_bonus_souls}
                  onChange={(e) => {
                    option({ no_bonus_souls: !e.target.checked });
                  }}
                />
                {t('host.bonus_souls')}
              </label>
            )}
            {/* Off: the players sit in the order they joined. */}
            <label className="check">
              <input
                type="checkbox"
                checked={!!options.random_seats}
                onChange={(e) => {
                  option({ random_seats: e.target.checked });
                }}
              />
              {t('host.random_seats')}
            </label>
            {/* Off and locked when the sets are too small; the tooltip
                names the short groups. */}
            <label className="check" title={blocked ? deckRatiosText(blocked) : undefined}>
              <input
                type="checkbox"
                checked={blocked === '' && !options.no_deck_ratios}
                disabled={blocked !== ''}
                onChange={(e) => {
                  option({ no_deck_ratios: !e.target.checked });
                }}
              />
              {t('host.deck_ratios')}
            </label>
            {/* Off by default; only with Requiem (R-SETUP-16, 14.13). */}
            {rooms && (
              <label className="check">
                <input
                  type="checkbox"
                  checked={!!options.room_deck}
                  onChange={(e) => {
                    option({ room_deck: e.target.checked });
                  }}
                />
                {t('host.room_deck')}
              </label>
            )}
            {/* Disabled by default (ADR 020, item 1.2). */}
            <div className="row">
              <label className="field">
                {t('host.rollback')}
                <select
                  value={options.rollback ?? ''}
                  aria-describedby="rollback-hint"
                  onChange={(e) => {
                    option({ rollback: (e.target.value || undefined) as RollbackMode | undefined });
                  }}
                >
                  <option value="">{t('host.rollback_off')}</option>
                  <option value="soft">{t('host.rollback_soft')}</option>
                  <option value="god">{t('host.rollback_god')}</option>
                </select>
              </label>
            </div>
            <p id="rollback-hint" className="muted hint">
              {rollbackHint(options.rollback ?? '')}
            </p>
          </>
        )}
        <button type="submit" className="primary" disabled={busy}>
          {remote ? t('host.create') : t('host.host')}
        </button>
        {/* A server's saved games are in its Games list already. */}
        {!remote && (
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
            {t('host.continue_saved')}
          </button>
        )}
      </form>
      <p className="status" role="status">
        {error}
      </p>
      <button
        type="button"
        onClick={() => {
          show('join');
        }}
      >
        {t('lobby.back')}
      </button>
    </div>
  );
}
