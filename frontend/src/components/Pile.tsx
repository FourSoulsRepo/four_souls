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

// HiddenHand is an opponent's hand: one loot back with the card count
// (V-26, F-08). The backs never tell the cards apart. Its zone is the
// player's hand, so drawn cards fly to it.
export function HiddenHand({ playerId, count, size }: { playerId: number; count: number; size: CardSize }) {
  if (count === 0) {
    return null;
  }
  const label = `Hand: ${String(count)} cards`;
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`hand-${String(playerId)}`}>
      <Card size={size} back="loot" />
      <span className="pile-count">{count}</span>
    </div>
  );
}
