import type { CSSProperties } from 'react';
import type { Kind } from '../cards';
import type { CardView } from '../protocol';
import { Card, type CardSize } from './Card';
import type { Handle } from './Interact';
import './Table.css';

// title is a deck name for labels: "loot" -> "Loot".
function title(kind: Kind): string {
  return kind.charAt(0).toUpperCase() + kind.slice(1);
}

// Deck is a face-down pile with its card count.
const noHandle: Handle = {};

export function Deck({
  kind,
  count,
  size,
  handle = noHandle,
}: {
  kind: Kind;
  count: number;
  size: CardSize;
  handle?: Handle;
}) {
  const label = `${title(kind)} deck: ${String(count)} cards`;
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`deck-${kind}`}>
      {count > 0 ? (
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
  const label = `${title(kind)} discard: ${String(cards?.length ?? 0)} cards`;
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`discard-${kind}`}>
      {top ? <Card cardRef={top.card} size={size} objectId={top.id} /> : <div className={`pile-empty size-${size}`}>Discard</div>}
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
  const label = `Hand: ${String(count)} cards`;
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
