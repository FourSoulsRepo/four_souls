import { useEffect, useState } from 'react';
import { deckRatiosBlocked, getCardSets } from '../bridge';
import { t } from '../i18n';
import { useEscape } from '../keys';
import type { Options, RollbackMode } from '../protocol';
import { baseGame, setName, sortSets } from '../sets';
import { useApp } from '../store';
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
// Games list creates it on that server (N-17, 15.98).
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
  const [options, setOptions] = useState<Options>({
    ...simple,
    response_timer: last.response_timer,
    no_bonus_souls: last.no_bonus_souls,
    no_deck_ratios: last.no_deck_ratios,
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

  // The deck ratios need enough cards in each group (GS-07).
  useEffect(() => {
    let live = true;
    void deckRatiosBlocked(chosen, seats).then((why) => {
      if (live) {
        setBlocked(why);
      }
    });
    return () => {
      live = false;
    };
  }, [chosen, seats]);

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
                  rollback: options.rollback ?? '',
                },
              });
              const game = { seats, sets: chosen, options: detailed ? options : undefined, password: password || undefined };
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
        <label className="field">
          {t('host.players')}
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
            {/* Draft and bans need the ban and pick screen, which comes with
                the tournament tools (roadmap, Future); until then a game
                with them could not start. */}
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
            {/* Off and locked when the sets are too small; the tooltip
                names the short groups. */}
            <label className="check" title={blocked ?? undefined}>
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
