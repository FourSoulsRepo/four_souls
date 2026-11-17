import { BotIcon } from '../components/BotIcon';
import { t } from '../i18n';
import { botLevelName, botLevels, isAway, optionsText, seatName } from '../lobby';
import { defaultPort, type Seat } from '../protocol';
import { useApp } from '../store';
import './screens.css';

// The table before the game: seats, ready, bots, and how friends join.
export function TableScreen() {
  const table = useApp((s) => s.lobby.table);
  const welcome = useApp((s) => s.lobby.welcome);
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
  const levels = botLevels(table, welcome);

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
        <ul className="list seats" aria-label={t('lobby.seats.seats')}>
          {seats.map((s) => (
            <li key={s.seat} data-seat={s.seat}>
              {s.bot ? (
                <BotSeat seat={s} levels={levels} />
              ) : (
                <span>
                  {s.name || <span className="muted">{t('lobby.seats.waiting')}</span>}
                  {s.seat === table.you && <span className="muted"> {t('lobby.seats.you')}</span>}
                  {isAway(s) && <span className="muted"> · {table.loaded ? t('lobby.seats.not_back') : t('lobby.seats.away')}</span>}
                  {table.loaded && s.connected && s.seat !== table.you && <span className="muted"> · {t('lobby.seats.back')}</span>}
                </span>
              )}
              {!s.name && levels.length > 0 ? <AddBot seat={s.seat} levels={levels} /> : <span className="seat-state">{s.ready ? t('lobby.seats.is_ready') : ''}</span>}
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

// AddBot is a free seat's buttons for the host: one per bot level the
// server allows (ADR 033, Decision 10.10).
function AddBot({ seat, levels }: { seat: number; levels: string[] }) {
  const seatBot = useApp((s) => s.seatBot);
  return (
    <span className="seat-bot" role="group" aria-label={t('lobby.bots.add')}>
      <span className="muted">{t('lobby.bots.add')}</span>
      {levels.map((level) => (
        <button
          key={level}
          type="button"
          aria-label={t('lobby.bots.add_level', { level: botLevelName(level) })}
          onClick={() => {
            seatBot(seat, level);
          }}
        >
          {botLevelName(level)}
        </button>
      ))}
    </span>
  );
}

// BotSeat is a bot's seat: the robot, its name and level. The host also
// changes its level or takes it away; levels is empty for the others.
function BotSeat({ seat, levels }: { seat: Seat; levels: string[] }) {
  const seatBot = useApp((s) => s.seatBot);
  const level = seat.bot ?? '';
  if (levels.length === 0) {
    return (
      <span>
        <BotIcon level={level} />
        {seatName(seat)}
      </span>
    );
  }
  const name = t('bot.number', { n: seat.bot_n ?? 0 });
  // A level the server no longer offers still shows as the bot's own.
  const options = levels.includes(level) ? levels : [...levels, level];
  return (
    <span className="seat-bot">
      <span>
        <BotIcon level={level} />
        {name}
      </span>
      <select
        aria-label={t('lobby.bots.level', { name })}
        value={level}
        onChange={(e) => {
          seatBot(seat.seat, e.target.value);
        }}
      >
        {options.map((l) => (
          <option key={l} value={l}>
            {botLevelName(l)}
          </option>
        ))}
      </select>
      <button
        type="button"
        className="forget"
        aria-label={t('lobby.bots.remove', { name })}
        title={t('lobby.bots.remove', { name })}
        onClick={() => {
          seatBot(seat.seat, '');
        }}
      >
        ×
      </button>
    </span>
  );
}
