// The challenge table (14.25, ADR 028 Decision 11): what the board
// shows of a co-op challenge. The view says it all; these helpers only
// name and sort it. Pure, for Vitest.
import { say, t } from './i18n';
import type { GameView, Msg, MsgNames, PresentView, SlotView } from './protocol';

// monsterRows are the slots whose cards are monsters in play: the
// monster slots, a challenge's final boss and its minions (ADR 028,
// Decision 7.3).
export function monsterRows(view: GameView): SlotView[] {
  return [...(view.monsters ?? []), ...(view.bosses ?? []), ...(view.minions ?? [])];
}

// presentKind names a present's kind on its card: "¢", "Loot".
export function presentKind(kind: string): string {
  switch (kind) {
    case 'loot':
      return t('challenge.present.loot');
    case 'treasure':
      return t('challenge.present.treasure');
    case 'soul':
      return t('challenge.present.soul');
    default:
      return t('challenge.present.cents');
  }
}

// Present is one present as the challenge panel draws it (R-KRAMPUS-05):
// face down with its counters; face up once Happy Holidays! flipped it
// and nobody has it yet; given, once a player gained it as an item; or
// destroyed, in the treasure discard.
export interface Present {
  kind: string;
  // caption is the kind under the card: "¢", "Loot".
  caption: string;
  state: 'down' | 'up' | 'given' | 'destroyed';
  // label says it all, for the screen reader and the tooltip.
  label: string;
  counters: number;
  card?: string;
}

// presents lists the challenge's presents in their order (¢, loot,
// treasure, soul), named; player names who gained one.
export function presents(view: GameView, player: (p: number) => string): Present[] {
  return (view.challenge?.presents ?? []).map((p) => present(p, player));
}

function present(p: PresentView, player: (p: number) => string): Present {
  const caption = presentKind(p.kind);
  const base = { kind: p.kind, caption, counters: p.counters, ...(p.card ? { card: p.card } : {}) };
  if (p.destroyed) {
    return { ...base, state: 'destroyed', label: t('challenge.present_destroyed', { present: caption }) };
  }
  if (p.flipped && p.player >= 0) {
    return { ...base, state: 'given', label: t('challenge.present_given', { present: caption, player: player(p.player) }) };
  }
  if (p.flipped) {
    return { ...base, state: 'up', label: t('challenge.present_up', { present: caption }) };
  }
  return { ...base, state: 'down', label: t('challenge.present_down', { present: caption, n: p.counters }) };
}

// timerText is the timer's tooltip and name: "8 rounds left".
export function timerText(n: number): string {
  return t('challenge.timer', { n });
}

// ChallengeEnd is the game-over title of a challenge, and why it ended.
export interface ChallengeEnd {
  title: string;
  reason: string;
}

// challengeEnd is a finished challenge's game-over text (R-CHAL-14,
// ADR 028 Decision 11.6): won by both players, or lost with the reason;
// null in a regular game or while it goes on.
export function challengeEnd(view: GameView, names: MsgNames): ChallengeEnd | null {
  const c = view.challenge;
  if (!c || !view.over) {
    return null;
  }
  if (c.lost) {
    return { title: t('over.challenge_lost'), reason: reasonText(c.reason, names) };
  }
  return { title: t('over.challenge_won'), reason: '' };
}

function reasonText(reason: Msg | undefined, names: MsgNames): string {
  return reason?.key ? say(reason, names) : '';
}
