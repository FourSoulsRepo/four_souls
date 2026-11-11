import { AnimatePresence, motion } from 'motion/react';
import type { CSSProperties } from 'react';
import type { Kind } from '../cards';
import { dealDelays } from '../choice';
import { bezier, feel } from '../feel';
import { t } from '../i18n';
import { scaleFor } from '../playback';
import { Deck as DeckKind, type CardView } from '../protocol';
import { seatOf, useApp } from '../store';
import { Card, type CardSize } from './Card';
import type { Handle } from './Interact';
import { still } from './lunge';
import './Table.css';
import { px, rem } from '../uiScale';

// deckLabel names a deck and its count: "Loot deck: 3 cards".
function deckLabel(kind: Kind, n: number): string {
  switch (kind) {
    case 'monster':
      return t('deck.monster', { n });
    case 'treasure':
      return t('deck.treasure', { n });
    case 'room':
      return t('deck.room', { n });
    default:
      return t('deck.loot', { n });
  }
}

// discardLabel names a discard pile and its count.
function discardLabel(kind: Kind, n: number): string {
  switch (kind) {
    case 'monster':
      return t('deck.monster_discard', { n });
    case 'treasure':
      return t('deck.treasure_discard', { n });
    case 'room':
      return t('deck.room_discard', { n });
    default:
      return t('deck.loot_discard', { n });
  }
}

// Deck is a face-down pile with its card count. A drawn card that waits
// for its place (top) lies face up on it (19.9). A top card only the
// viewer may look at (peek, 19.30) also lies face up; the deck still
// takes clicks, e.g. to buy from it.
const noHandle: Handle = {};

export function Deck({
  kind,
  count,
  size,
  handle = noHandle,
  top,
  peek,
}: {
  kind: Kind;
  count: number;
  size: CardSize;
  handle?: Handle;
  top?: CardView | undefined;
  peek?: CardView | undefined;
}) {
  const label = deckLabel(kind, count);
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`deck-${kind}`}>
      {top ? (
        <Card cardRef={top.card} size={size} objectId={top.id} highlight="pick" />
      ) : peek ? (
        <Card cardRef={peek.card} size={size} highlight={handle.highlight} cost={handle.cost} onClick={handle.onClick} />
      ) : count > 0 ? (
        <Card size={size} back={kind} highlight={handle.highlight} cost={handle.cost} onClick={handle.onClick} />
      ) : (
        <div className={`pile-empty size-${size}`} />
      )}
      <LookBacks kind={kind} />
      <span className="pile-count">{count}</span>
    </div>
  );
}

// lookDeck is the engine's deck (protocol Deck) of a pile kind, or -1.
function lookDeck(kind: Kind): number {
  switch (kind) {
    case 'treasure':
      return DeckKind.treasure;
    case 'loot':
      return DeckKind.loot;
    case 'monster':
      return DeckKind.monster;
    case 'room':
      return DeckKind.room;
    default:
      return -1;
  }
}

// useOthersLook is how many top cards of this deck another player looks
// at now, and who (20.18); 0 when nobody else does.
function useOthersLook(kind: Kind): { n: number; who: number } {
  const n = useApp((s) => {
    const w = s.game?.view.waiting;
    const look = w?.look;
    return look && (look.deck ?? 0) === lookDeck(kind) && w.player !== seatOf(s) ? (look.count ?? 0) : 0;
  });
  const who = useApp((s) => s.game?.view.waiting.player ?? -1);
  return { n, who };
}

// LookBacks are the small backs over a deck another player looks at
// (20.18, ADR 027 "Choosing cards" 12.13): one per card, fanned on
// the pile under its count, faces hidden. They rise out of the deck as
// the look starts and sink back into it as it ends.
function LookBacks({ kind }: { kind: Kind }) {
  const { n, who } = useOthersLook(kind);
  const name = useApp((s) => s.game?.seats?.find((x) => x.seat === who)?.name ?? '');
  const calm = useApp((s) => scaleFor(s.settings.animation, 0, 0) === 0) || still();
  const scale = useApp((s) => scaleFor(s.settings.animation, 0, 0));
  const o = feel.choice;
  const backs = Array.from({ length: n }, (_, i) => i);
  return (
    <AnimatePresence>
      {n > 0 && (
        // anim: choice-others
        <motion.div key="look" className="deck-look" role="img" aria-label={t('deck.looking', { player: name, n })}>
          {backs.map((i) => {
            const spread = i - (n - 1) / 2;
            const at = { x: px(spread * 9), y: px(-12 + Math.abs(spread) * 2), rotate: spread * 7, scale: 1 };
            const home = { x: 0, y: 0, rotate: 0, scale: 0.6 };
            const delay = (dealDelays(n).at(i) ?? 0) * scale;
            return (
              // anim: choice-others
              <motion.div
                key={i}
                className="deck-look-back"
                initial={calm ? false : home}
                animate={at}
                exit={calm ? { opacity: 0, transition: { duration: 0 } } : { ...home, transition: { duration: o.back * scale, ease: bezier(o.backCurve) } }}
                transition={{ duration: o.others * scale, delay, ease: bezier(o.othersCurve) }}
              >
                <Card back={kind} size="small" />
              </motion.div>
            );
          })}
        </motion.div>
      )}
    </AnimatePresence>
  );
}

// Discard shows the top card of a discard pile, face up, and the pile's size.
export function Discard({ kind, top, count, size }: { kind: Kind; top: CardView | null; count: number; size: CardSize }) {
  const label = discardLabel(kind, count);
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`discard-${kind}`}>
      {top ? <Card cardRef={top.card} size={size} objectId={top.id} /> : <div className={`pile-empty size-${size}`}>{t('deck.discard')}</div>}
    </div>
  );
}

// Most backs a hidden hand fans out; the count tells the rest.
const fannedBacks = 8;

// othersFeel is feel.others as the CSS variables a hidden hand's backs
// lift, float and sway by (20.13, Table.css).
const othersFeel = (() => {
  const o = feel.others;
  return {
    '--others-hover': `${String(o.hover)}s`,
    '--others-hover-lift': rem(o.hoverLift),
    '--others-hover-curve': o.hoverCurve,
    '--others-hold': `${String(o.hold)}s`,
    '--others-hold-lift': rem(o.holdLift),
    '--others-hold-grow': String(o.holdGrow),
    '--others-hold-curve': o.holdCurve,
    '--others-sway': `${String(o.sway)}deg`,
    '--others-sway-period': `${String(o.swayTime * 2)}s`,
    '--others-sway-curve': o.swayCurve,
  } as CSSProperties;
})();

// HiddenHand is an opponent's hand: fanned loot backs with the card
// count (V-26, F-08, H-17). The backs never tell the cards apart. Its
// zone is the player's hand, so drawn cards fly to it; an empty hand
// keeps the zone with only the count. A back the player hovers or
// drags moves as Presence.tsx marks it (20.13); speed "off" and
// reduced motion (still) set it in place at once.
export function HiddenHand({ playerId, count, size }: { playerId: number; count: number; size: CardSize }) {
  const label = t('table.hand', { n: count });
  const shown = Math.min(count, fannedBacks);
  const calm = useApp((s) => scaleFor(s.settings.animation, 0, 0) === 0) || still();
  return (
    <div
      className={`hidden-hand${calm ? ' still' : ''}`}
      style={othersFeel}
      role="group"
      aria-label={label}
      title={label}
      data-zone={`hand-${String(playerId)}`}
    >
      {Array.from({ length: shown }, (_, i) => (
        <div key={i} className="hidden-back" style={{ '--spread': String(i - (shown - 1) / 2) } as CSSProperties}>
          <Card size={size} back="loot" />
        </div>
      ))}
      <span className="pile-count">{count}</span>
    </div>
  );
}

// RevealedHand is an opponent's hand the view sends face up, e.g.
// under Guppy's Eye (19.5): every card, fanned like the backs, each
// zoomable. Its zone is the player's hand, like a hidden one.
export function RevealedHand({ playerId, cards, size }: { playerId: number; cards: CardView[]; size: CardSize }) {
  const label = t('table.revealed_hand', { n: cards.length });
  return (
    <div className="hidden-hand revealed-hand" role="group" aria-label={label} title={label} data-zone={`hand-${String(playerId)}`}>
      {cards.map((c, i) => (
        <div key={c.id} className="hidden-back" style={{ '--spread': String(i - (cards.length - 1) / 2) } as CSSProperties}>
          <Card cardRef={c.card} size={size} objectId={c.id} />
        </div>
      ))}
      <span className="pile-count">{cards.length}</span>
    </div>
  );
}
