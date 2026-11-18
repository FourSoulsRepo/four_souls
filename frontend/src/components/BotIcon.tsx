import { botLead, namePieces, seatMarks, seatName, type NamedSeat } from '../lobby';
import { useApp } from '../store';
import './BotIcon.css';

// BotLook is how a bot level's robot looks; a level this app does not
// know (a newer server's) looks standard.
export type BotLook = 'standard' | 'hard';

// botLook is the robot of a bot level (ADR 033, Decision 10.3).
export function botLook(level: string): BotLook {
  return level === 'hard' ? 'hard' : 'standard';
}

// botPaths draw the robots, one per look, in a 24×24 box with the even-odd
// rule; turn.ts builds the turn banner's copy by hand. Both share the
// head and the ears, so they read as one family. Standard has round
// eyes, a straight mouth and an antenna; hard has slanted eyes under
// brows, a jagged mouth and two horns.
export const botPaths: Record<BotLook, string> = {
  standard:
    'M12 1.5a1.5 1.5 0 0 1 .75 2.8V7H18a3 3 0 0 1 3 3v7a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3v-7a3 3 0 0 1 3-3h5.25V4.3A1.5 1.5 0 0 1 12 1.5ZM8.5 10.75a1.75 1.75 0 1 0 0 3.5 1.75 1.75 0 0 0 0-3.5Zm7 0a1.75 1.75 0 1 0 0 3.5 1.75 1.75 0 0 0 0-3.5ZM8 16v1.5h8V16ZM1 11h1.5v5H1Zm20.5 0H23v5h-1.5Z',
  hard: 'M3 10a3 3 0 0 1 3-3L4.5 1.5 10 7h4l5.5-5.5L18 7a3 3 0 0 1 3 3v7a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3ZM5.5 9.2l5.5 2.6-.4 1.4-1.2-.5-.6 1.8-2.6-1.1Zm13 0-5.5 2.6.4 1.4 1.2-.5.6 1.8 2.6-1.1ZM7 18.6V16l2.5 1.3L12 16l2.5 1.3L17 16v2.6ZM1 11h1.5v5H1Zm20.5 0H23v5h-1.5Z',
};

// BotIcon is the robot before a bot seat's name (ADR 033, Decision 10.3),
// drawn for the seat's level. The name says "Bot" already, so screen
// readers skip it.
export function BotIcon({ level }: { level: string }) {
  const look = botLook(level);
  return (
    <svg className="bot-icon" viewBox="0 0 24 24" aria-hidden="true" data-testid="bot-icon" data-level={look}>
      <path fillRule="evenodd" d={look === 'hard' ? botPaths.hard : botPaths.standard} />
    </svg>
  );
}

// BotName is a bot's name after its level's robot, which keeps to its
// first word (botLead).
function BotName({ name, level }: { name: string; level: string }) {
  const [lead, rest] = botLead(name);
  return (
    <span className="bot-name">
      <span className="bot-lead">
        <BotIcon level={level} />
        {lead}
      </span>
      {rest}
    </span>
  );
}

// SeatName is a seat's name as seatName gives it; a bot's after its
// robot. Being a bot and its level come from the seat (19.214).
export function SeatName({ seat }: { seat: NamedSeat }) {
  // A new language renames the bots: draw again.
  useApp((s) => s.settings.language);
  return seat.bot ? <BotName name={seatName(seat)} level={seat.bot} /> : <>{seat.name}</>;
}

// PlayerName is the name of the game's seat player: SeatName, or
// fallback for a seat the game does not list.
export function PlayerName({ player, fallback }: { player: number; fallback: string }) {
  const seat = useApp((s) => s.game?.seats?.find((x) => x.seat === player));
  return seat ? <SeatName seat={seat} /> : <>{fallback}</>;
}

// BotText draws a line with its level's robot before every bot's name in
// it (23.12): a history row, a window's title. seats are the seats the
// names come from; the game's by default. A line is split at every
// seat's name (seatMarks), so a player whose name looks like a bot's
// never gets the robot (19.214).
export function BotText({ text, seats }: { text: string; seats?: readonly NamedSeat[] | null }) {
  const gameSeats = useApp((s) => s.game?.seats);
  // A new language renames the bots: draw again.
  useApp((s) => s.settings.language);
  return (
    <>
      {namePieces(text, seatMarks(seats ?? gameSeats)).map((p, i) =>
        p.bot ? (
          // A line's pieces never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <BotName key={i} name={p.text} level={p.bot} />
        ) : (
          p.text
        ),
      )}
    </>
  );
}
