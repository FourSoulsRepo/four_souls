import { optionsText } from '../lobby';
import { defaultPort } from '../protocol';
import { useApp } from '../store';
import './screens.css';

// The table before the game: seats, ready, and how friends join.
export function TableScreen() {
  const table = useApp((s) => s.lobby.table);
  const hosting = useApp((s) => s.hosting);
  const error = useApp((s) => s.error);
  const setReady = useApp((s) => s.setReady);
  const leaveTable = useApp((s) => s.leaveTable);
  if (!table) {
    return null;
  }
  const seats = table.seats ?? [];
  const ready = seats[table.you]?.ready ?? false;
  const port = hosting && hosting.port !== defaultPort ? `:${String(hosting.port)}` : '';

  return (
    <div className="screen">
      <h1>{table.loaded ? 'Saved game' : 'Table'}</h1>
      {hosting && (
        <div className="panel">
          <span className="muted">Friends join at</span>
          <ul className="list">
            {hosting.addresses.map((a) => (
              <li key={a}>{a + port}</li>
            ))}
          </ul>
        </div>
      )}
      <div className="panel">
        <ul className="list" aria-label="Seats">
          {seats.map((s) => (
            <li key={s.seat}>
              <span>
                {s.name || <span className="muted">Waiting for a player…</span>}
                {s.seat === table.you && <span className="muted"> (you)</span>}
                {s.name && !s.connected && <span className="muted"> · away</span>}
              </span>
              <span>{s.ready ? 'Ready' : ''}</span>
            </li>
          ))}
        </ul>
        <span className="muted">{optionsText(table.options)}</span>
        <span className="muted">The game starts when every seat is taken and everyone is ready.</span>
      </div>
      <p className="status" role="status">
        {error}
      </p>
      <div className="row">
        <button
          type="button"
          className={ready ? '' : 'primary'}
          onClick={() => {
            setReady(!ready);
          }}
        >
          {ready ? 'Not ready' : 'Ready'}
        </button>
        <button
          type="button"
          onClick={() => {
            void leaveTable();
          }}
        >
          Leave
        </button>
      </div>
    </div>
  );
}
