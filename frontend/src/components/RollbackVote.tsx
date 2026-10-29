import { moveText } from '../history';
import { t } from '../i18n';
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
  const who = (p: number) => (p === you ? t('rollback.you') : name(p));
  const votes = v.votes ?? [];
  const mine = votes.find((b) => b.seat === you);
  const yes = votes.filter((b) => b.yes && v.voters.includes(b.seat)).length;
  const voter = v.voters.includes(you);
  const title =
    v.seat === v.author
      ? v.author === you
        ? t('rollback.ask.you')
        : t('rollback.ask.other', { player: name(v.author) })
      : v.seat === you
        ? t('rollback.want.you', { author: name(v.author) })
        : v.author === you
          ? t('rollback.want.yours', { player: name(v.seat) })
          : t('rollback.want.other', { player: name(v.seat), author: name(v.author) });
  // The move is named as on the Undo button; the open vote's target is
  // still listed (17.8).
  const target = { step: v.step, seat: v.author, last: v.last };
  return (
    <section className="overlay-box rollback-vote" role="dialog" aria-label={t('rollback.label')}>
      <h2>{title}</h2>
      <p>{t('rollback.undoes', { move: moveText(moves, target, you, name), n: v.undo })}</p>
      <p className="muted">
        {t('rollback.needed', { yes, n: v.need })}
        {v.deadline ? <Timer deadline={v.deadline} /> : null}
      </p>
      <ul className="overlay-votes">
        {v.voters.map((p) => {
          const b = votes.find((x) => x.seat === p);
          return (
            <li key={p}>
              {who(p)}: {b ? (b.yes ? t('rollback.voted_yes') : t('rollback.voted_no')) : '…'}
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
                {t('rollback.yes')}
              </button>
              <button
                type="button"
                className={mine && !mine.yes ? 'chosen' : ''}
                onClick={() => {
                  vote(false);
                }}
              >
                {t('rollback.no')}
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
              {t('rollback.cancel')}
            </button>
          )}
        </div>
      )}
    </section>
  );
}
