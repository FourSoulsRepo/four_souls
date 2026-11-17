import { botLead, botNames, namePieces } from '../lobby';
import type { Seat } from '../protocol';
import { useApp } from '../store';
import './BotIcon.css';

// botPath draws the robot; turn.ts builds the turn banner's copy by hand.
export const botPath =
  'M12 1.5a1.5 1.5 0 0 1 .75 2.8V7H18a3 3 0 0 1 3 3v7a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3v-7a3 3 0 0 1 3-3h5.25V4.3A1.5 1.5 0 0 1 12 1.5ZM8.5 10.75a1.75 1.75 0 1 0 0 3.5 1.75 1.75 0 0 0 0-3.5Zm7 0a1.75 1.75 0 1 0 0 3.5 1.75 1.75 0 0 0 0-3.5ZM8 16v1.5h8V16ZM1 11h1.5v5H1Zm20.5 0H23v5h-1.5Z';

// BotIcon is the robot before a bot seat's name (ADR 030, Decision 10.3).
// The name says "Bot" already, so screen readers skip it.
export function BotIcon() {
  return (
    <svg className="bot-icon" viewBox="0 0 24 24" aria-hidden="true" data-testid="bot-icon">
      <path fillRule="evenodd" d={botPath} />
    </svg>
  );
}

// BotName is a bot's name after the robot, which keeps to its first
// word (botLead).
function BotName({ name }: { name: string }) {
  const [lead, rest] = botLead(name);
  return (
    <span className="bot-name">
      <span className="bot-lead">
        <BotIcon />
        {lead}
      </span>
      {rest}
    </span>
  );
}

// BotText draws a line with the robot before every bot's name in it
// (21.12): a seat name, a history row, a window's title. seats are the
// seats the names come from; the game's by default.
export function BotText({ text, seats }: { text: string; seats?: readonly Seat[] | null }) {
  const gameSeats = useApp((s) => s.game?.seats);
  // A new language renames the bots: draw again.
  useApp((s) => s.settings.language);
  return (
    <>
      {namePieces(text, botNames(seats ?? gameSeats)).map((p, i) =>
        p.bot ? (
          // A line's pieces never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <BotName key={i} name={p.text} />
        ) : (
          p.text
        ),
      )}
    </>
  );
}
