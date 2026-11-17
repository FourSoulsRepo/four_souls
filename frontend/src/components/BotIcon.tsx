import './BotIcon.css';

// BotIcon is the robot before a bot seat's name (ADR 030, Decision 10.3).
// The name says "Bot" already, so screen readers skip it.
export function BotIcon() {
  return (
    <svg className="bot-icon" viewBox="0 0 24 24" aria-hidden="true" data-testid="bot-icon">
      <path
        fillRule="evenodd"
        d="M12 1.5a1.5 1.5 0 0 1 .75 2.8V7H18a3 3 0 0 1 3 3v7a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3v-7a3 3 0 0 1 3-3h5.25V4.3A1.5 1.5 0 0 1 12 1.5ZM8.5 10.75a1.75 1.75 0 1 0 0 3.5 1.75 1.75 0 0 0 0-3.5Zm7 0a1.75 1.75 0 1 0 0 3.5 1.75 1.75 0 0 0 0-3.5ZM8 16v1.5h8V16ZM1 11h1.5v5H1Zm20.5 0H23v5h-1.5Z"
      />
    </svg>
  );
}
