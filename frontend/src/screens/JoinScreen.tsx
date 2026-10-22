import { useEffect, useState } from 'react';
import { useEscape } from '../keys';
import { canCreate, lanRows, serverLabel } from '../lobby';
import { protocolVersion } from '../protocol';
import { setList } from '../sets';
import { useApp } from '../store';
import './screens.css';

// lanEvery is the pause in ms between two LAN searches: the LAN screen
// searches all the time (15.97).
const lanEvery = 2000;

// Join: how the player came from Play (M-03), then the server's games
// and saved games.
export function JoinScreen() {
  const connected = useApp((s) => s.connected);
  const path = useApp((s) => s.path);
  if (connected) {
    return <GameList />;
  }
  return path === 'lan' ? <LanConnect /> : <Connect />;
}

// LanConnect: host a game in this app, or join one on this network
// (15.97). The list fills by itself; a started game this app held a
// seat in offers Rejoin. The address field keeps only the last address
// typed, and nothing goes to Recent.
function LanConnect() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const joinServer = useApp((s) => s.joinServer);
  const saveNickname = useApp((s) => s.saveNickname);
  const saveSettings = useApp((s) => s.saveSettings);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [address, setAddress] = useState(settings.lan_address);
  const [busy, setBusy] = useState(false);
  const lan = useApp((s) => s.lan);
  const rows = lanRows(lan.found, protocolVersion, settings.rejoin);
  useLanSearch();
  useEscape(() => {
    show('menu');
  });

  // go joins to; a typed address is kept, a rejoin sends the seat.
  const go = (to: string, opts: { typed?: boolean; rejoin?: boolean }) => {
    setBusy(true);
    void (async () => {
      if (await saveNickname(nickname)) {
        const s = useApp.getState().settings;
        if (opts.typed && s.lan_address !== to.trim()) {
          await saveSettings({ ...s, lan_address: to });
        }
        await joinServer(to, { recent: false, token: opts.rejoin ? s.rejoin.token : undefined });
      }
      setBusy(false);
    })();
  };

  return (
    <div className="screen">
      <h1>LAN</h1>
      <button
        type="button"
        className="primary"
        onClick={() => {
          show('host');
        }}
      >
        Host a game
      </button>
      <section className="panel" aria-label="On your network">
        <label className="field">
          Nickname
          <input
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
            }}
          />
        </label>
        <span className="muted">On your network</span>
        {rows.length === 0 && (
          <span className="muted selectable">{lan.error || 'No games found yet. Still looking…'}</span>
        )}
        <ul className="list">
          {rows.map((r) => (
            <li key={r.key}>
              <span className="selectable">{r.text}</span>
              <button
                type="button"
                className={r.rejoin ? 'primary' : undefined}
                disabled={busy || !r.joinable}
                onClick={() => {
                  go(r.address, { rejoin: r.rejoin });
                }}
              >
                {r.rejoin ? 'Rejoin' : 'Join'}
              </button>
            </li>
          ))}
        </ul>
        <form
          className="row"
          onSubmit={(e) => {
            e.preventDefault();
            go(address, { typed: true });
          }}
        >
          <label className="field grow">
            Host address
            <input
              value={address}
              placeholder="192.168.1.5"
              onChange={(e) => {
                setAddress(e.target.value);
              }}
            />
          </label>
          <button type="submit" disabled={busy}>
            Connect
          </button>
        </form>
      </section>
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

// useLanSearch searches the LAN again and again while the screen shows.
function useLanSearch() {
  const findLanGames = useApp((s) => s.findLanGames);
  useEffect(() => {
    let live = true;
    let timer = 0;
    const loop = async () => {
      await findLanGames();
      if (live) {
        timer = window.setTimeout(() => void loop(), lanEvery);
      }
    };
    void loop();
    return () => {
      live = false;
      window.clearTimeout(timer);
    };
  }, [findLanGames]);
}

function Connect() {
  const settings = useApp((s) => s.settings);
  const error = useApp((s) => s.error);
  const joinServer = useApp((s) => s.joinServer);
  const saveNickname = useApp((s) => s.saveNickname);
  const forgetServer = useApp((s) => s.forgetServer);
  const show = useApp((s) => s.show);
  const [nickname, setNickname] = useState(settings.nickname);
  const [address, setAddress] = useState(settings.recent[0]?.address ?? '');
  const [busy, setBusy] = useState(false);
  const lan = useApp((s) => s.lan);
  const findLanGames = useApp((s) => s.findLanGames);
  const rows = lanRows(lan.found, protocolVersion, settings.rejoin);

  useEffect(() => {
    void findLanGames();
  }, [findLanGames]);

  const go = (to: string) => {
    setBusy(true);
    void (async () => {
      if (await saveNickname(nickname)) {
        await joinServer(to);
      }
      setBusy(false);
    })();
  };

  return (
    <div className="screen">
      <h1>Join a game</h1>
      <form
        className="panel"
        onSubmit={(e) => {
          e.preventDefault();
          go(address);
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
          Host address
          <input
            value={address}
            placeholder="192.168.1.5"
            onChange={(e) => {
              setAddress(e.target.value);
            }}
          />
        </label>
        <button type="submit" className="primary" disabled={busy}>
          Connect
        </button>
        {settings.recent.length > 0 && (
          <>
            <span className="muted">Recent</span>
            {/* A server shows by its name, the address on hover (N-18). */}
            <ul className="row recent" aria-label="Recent">
              {settings.recent.map((r) => (
                <li key={r.address}>
                  <button
                    type="button"
                    title={r.address}
                    disabled={busy}
                    onClick={() => {
                      setAddress(r.address);
                      go(r.address);
                    }}
                  >
                    {serverLabel(r)}
                  </button>
                  <button
                    type="button"
                    className="forget"
                    aria-label={`Remove ${serverLabel(r)} from Recent`}
                    title="Remove from Recent"
                    disabled={busy}
                    onClick={() => {
                      void forgetServer(r.address);
                    }}
                  >
                    ×
                  </button>
                </li>
              ))}
            </ul>
          </>
        )}
      </form>
      <section className="panel" aria-label="On your network">
        <div className="row lan-head">
          <span className="muted">On your network</span>
          <button
            type="button"
            disabled={lan.searching}
            onClick={() => {
              void findLanGames();
            }}
          >
            {lan.searching ? 'Searching…' : 'Search again'}
          </button>
        </div>
        {rows.length === 0 && !lan.searching && (
          <span className="muted selectable">{lan.error || 'No games found. Type the host address above.'}</span>
        )}
        <ul className="list">
          {rows.map((r) => (
            <li key={r.key}>
              <span className="selectable">{r.text}</span>
              <button
                type="button"
                disabled={busy || !r.joinable}
                onClick={() => {
                  setAddress(r.address);
                  go(r.address);
                }}
              >
                Join
              </button>
            </li>
          ))}
        </ul>
      </section>
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

function GameList() {
  const [passwords, setPasswords] = useState<Record<string, string>>({});
  const lobby = useApp((s) => s.lobby);
  const error = useApp((s) => s.error);
  const joinGame = useApp((s) => s.joinGame);
  const loadSave = useApp((s) => s.loadSave);
  const quit = useApp((s) => s.quit);
  const show = useApp((s) => s.show);
  const create = useApp((s) => canCreate(s.lobby, s.hosting !== null));
  const server = useApp((s) => s.lobby.welcome?.name ?? '');
  useEscape(() => {
    void quit('join');
  });

  return (
    <div className="screen">
      <h1>Games</h1>
      {server && <p className="server-name">{server}</p>}
      {lobby.saved && <p className="muted">The game was saved. Continue it from this list later.</p>}
      <div className="panel">
        {lobby.games.length === 0 && (
          <span className="muted">{create ? 'No games yet. Create one.' : 'No games yet. Ask the host to open one.'}</span>
        )}
        <ul className="list">
          {lobby.games.map((g) => (
            <li key={g.id} data-game={g.id}>
              <span>
                {g.locked && <span aria-label="Locked">🔒 </span>}
                {g.host}&apos;s game · {g.taken}/{g.seats} players{g.started ? ' · playing' : ''}
                {g.sets && g.sets.length > 0 && <span className="muted"> · {setList(g.sets)}</span>}
              </span>
              {g.locked && !g.started && g.taken < g.seats && (
                <input
                  type="password"
                  className="game-password"
                  aria-label={`Password for ${g.host}'s game`}
                  placeholder="Password"
                  autoComplete="off"
                  value={passwords[g.id] ?? ''}
                  onChange={(e) => {
                    setPasswords({ ...passwords, [g.id]: e.target.value });
                  }}
                />
              )}
              <button
                type="button"
                disabled={g.started || g.taken >= g.seats}
                onClick={() => {
                  joinGame(g.id, passwords[g.id]);
                }}
              >
                Join
              </button>
            </li>
          ))}
        </ul>
        {create && (
          <button
            type="button"
            className="primary"
            onClick={() => {
              show('host');
            }}
          >
            Create a game
          </button>
        )}
      </div>
      {lobby.saves.length > 0 && (
        <div className="panel">
          <h2>Saved games</h2>
          <ul className="list">
            {lobby.saves.map((s) => (
              <li key={s.id}>
                <span>
                  {(s.seats ?? []).join(', ')} · {new Date(s.saved).toLocaleString()}
                </span>
                <button
                  type="button"
                  onClick={() => {
                    loadSave(s.id);
                  }}
                >
                  Continue
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
      <p className="status" role="status">
        {error}
      </p>
      <button
        type="button"
        onClick={() => {
          void quit('join');
        }}
      >
        Back
      </button>
    </div>
  );
}
