import { useEffect, useState } from 'react';
import { canCreate, lanRows, serverLabel } from '../lobby';
import { protocolVersion } from '../protocol';
import { setList } from '../sets';
import { useApp } from '../store';
import './screens.css';

// Join: an address, then the server's games and saved games.
export function JoinScreen() {
  const connected = useApp((s) => s.connected);
  return connected ? <GameList /> : <Connect />;
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
  const rows = lanRows(lan.found, protocolVersion);

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
      <button
        type="button"
        className="primary"
        onClick={() => {
          show('host');
        }}
      >
        Host a game
      </button>
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
          void quit();
        }}
      >
        Back
      </button>
    </div>
  );
}
