import type { CSSProperties } from 'react';
import type { Kind } from '../cards';
import { t } from '../i18n';
import type { CardView } from '../protocol';
import { Card, type CardSize } from './Card';
import type { Handle } from './Interact';
import './Table.css';

// deckLabel names a deck and its count: "Loot deck: 3 cards".
function deckLabel(kind: Kind, n: number): string {
  switch (kind) {
    case 'monster':
      return t('deck.monster', { n });
    case 'treasure':
      return t('deck.treasure', { n });
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
        <Card cardRef={peek.card} size={size} highlight={handle.highlight} onClick={handle.onClick} />
      ) : count > 0 ? (
        <Card size={size} back={kind} highlight={handle.highlight} onClick={handle.onClick} />
      ) : (
        <div className={`pile-empty size-${size}`} />
      )}
      <span className="pile-count">{count}</span>
    </div>
  );
}

// Discard shows the top card of a discard pile, face up.
export function Discard({ kind, cards, size }: { kind: Kind; cards: CardView[] | null; size: CardSize }) {
  const top = cards?.at(-1);
  const label = discardLabel(kind, cards?.length ?? 0);
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`discard-${kind}`}>
      {top ? <Card cardRef={top.card} size={size} objectId={top.id} /> : <div className={`pile-empty size-${size}`}>{t('deck.discard')}</div>}
    </div>
  );
}

// Most backs a hidden hand fans out; the count tells the rest.
const fannedBacks = 8;

// HiddenHand is an opponent's hand: fanned loot backs with the card
// count (V-26, F-08, H-17). The backs never tell the cards apart. Its
// zone is the player's hand, so drawn cards fly to it; an empty hand
// keeps the zone with only the count.
export function HiddenHand({ playerId, count, size }: { playerId: number; count: number; size: CardSize }) {
  const label = t('table.hand', { n: count });
  const shown = Math.min(count, fannedBacks);
  return (
    <div className="hidden-hand" role="group" aria-label={label} title={label} data-zone={`hand-${String(playerId)}`}>
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
