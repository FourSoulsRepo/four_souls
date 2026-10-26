import { t } from '../i18n';
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
      <h1>{table.loaded ? t('lobby.seats.saved_title') : t('lobby.seats.title')}</h1>
      {hosting && (
        <div className="panel">
          <span className="muted">{t('lobby.seats.join_at')}</span>
          <ul className="list selectable" aria-label={t('lobby.seats.addresses')}>
            {hosting.addresses.map((a) => (
              <li key={a}>{a + port}</li>
            ))}
          </ul>
        </div>
      )}
      <div className="panel">
        <ul className="list" aria-label={t('lobby.seats.seats')}>
          {seats.map((s) => (
            <li key={s.seat}>
              <span>
                {s.name || <span className="muted">{t('lobby.seats.waiting')}</span>}
                {s.seat === table.you && <span className="muted"> {t('lobby.seats.you')}</span>}
                {s.name && !s.connected && (
                  <span className="muted"> · {table.loaded ? t('lobby.seats.not_back') : t('lobby.seats.away')}</span>
                )}
                {table.loaded && s.connected && s.seat !== table.you && <span className="muted"> · {t('lobby.seats.back')}</span>}
              </span>
              <span>{s.ready ? t('lobby.seats.is_ready') : ''}</span>
            </li>
          ))}
        </ul>
        <span className="muted">{optionsText(table.options)}</span>
        <span className="muted">
          {table.loaded
            ? t('lobby.seats.saved_hint')
            : t('lobby.seats.hint')}
        </span>
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
          {ready ? t('lobby.seats.not_ready') : t('lobby.seats.ready')}
        </button>
        <button
          type="button"
          onClick={() => {
            void leaveTable();
          }}
        >
          {t('lobby.seats.leave')}
        </button>
      </div>
    </div>
  );
}
