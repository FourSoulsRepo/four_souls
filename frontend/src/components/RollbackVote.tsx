import { moveText } from '../history';
import type { Update } from '../protocol';
import { useApp } from '../store';
import { Timer } from './ActionBar';
import './Overlays.css';

// RollbackVotePanel is the open Soft vote (ADR 020, item 5). The game
// waits; the panel leaves the table and the history in view, so the
// voters can see what the rollback undoes. Voters answer yes or no and
// may change it; whoever opened the vote may cancel it; the rest see
// how it stands.
export function RollbackVotePanel({ game, you, name }: { game: Update; you: number; name: (p: number) => string }) {
  const moves = useApp((s) => s.moves);
  const vote = useApp((s) => s.rollbackVote);
  const cancel = useApp((s) => s.rollbackCancel);
  const v = game.rollback_vote;
  if (!v) {
    return null;
  }
  const who = (p: number) => (p === you ? 'You' : name(p));
  const votes = v.votes ?? [];
  const mine = votes.find((b) => b.seat === you);
  const yes = votes.filter((b) => b.yes && v.voters.includes(b.seat)).length;
  const voter = v.voters.includes(you);
  const title =
    v.seat === v.author
      ? `${who(v.author)} ${v.author === you ? 'ask' : 'asks'} to undo a move`
      : `${who(v.seat)} ${v.seat === you ? 'want' : 'wants'} to undo ${v.author === you ? 'your' : `${name(v.author)}'s`} move`;
  const steps = `${String(v.undo)} ${v.undo === 1 ? 'step' : 'steps'}`;
  // The move is named as on the Undo button; the open vote's target is
  // still listed (17.8).
  const target = game.rollback_steps?.find((t) => t.step === v.step) ?? { step: v.step, seat: v.author };
  return (
    <section className="overlay-box rollback-vote" role="dialog" aria-label="Rollback vote">
      <h2>{title}</h2>
      <p>
        “{moveText(moves, target, you, name)}”: undoes {steps}.
      </p>
      <p className="muted">
        {String(yes)} of {String(v.need)} yes needed
        {v.deadline ? <Timer deadline={v.deadline} /> : null}
      </p>
      <ul className="overlay-votes">
        {v.voters.map((p) => {
          const b = votes.find((x) => x.seat === p);
          return (
            <li key={p}>
              {who(p)}: {b ? (b.yes ? 'yes' : 'no') : '…'}
            </li>
          );
        })}
      </ul>
      {(voter || v.seat === you) && (
        <div className="overlay-buttons">
          {voter && (
            <>
              <button
                type="button"
                className={mine?.yes ? 'chosen' : ''}
                onClick={() => {
                  vote(true);
                }}
              >
                Yes, undo
              </button>
              <button
                type="button"
                className={mine && !mine.yes ? 'chosen' : ''}
                onClick={() => {
                  vote(false);
                }}
              >
                No
              </button>
            </>
          )}
          {v.seat === you && (
            <button
              type="button"
              className="secondary"
              onClick={() => {
                cancel();
              }}
            >
              Cancel vote
            </button>
          )}
        </div>
      )}
    </section>
  );
}
