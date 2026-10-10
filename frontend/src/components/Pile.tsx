import type { Kind } from '../cards';
import type { CardView } from '../protocol';
import { Card, type CardSize } from './Card';
import './Table.css';

// title is a deck name for labels: "loot" -> "Loot".
function title(kind: Kind): string {
  return kind.charAt(0).toUpperCase() + kind.slice(1);
}

// Deck is a face-down pile with its card count.
export function Deck({ kind, count, size }: { kind: Kind; count: number; size: CardSize }) {
  const label = `${title(kind)} deck: ${String(count)} cards`;
  return (
    <div className="pile" role="group" aria-label={label} title={label} data-zone={`deck-${kind}`}>
      {count > 0 ? <Card size={size} back={kind} /> : <div className={`pile-empty size-${size}`} />}
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
