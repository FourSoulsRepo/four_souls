import { Card } from '../components/Card';
import { useApp } from '../store';
import './screens.css';

// CardGallery shows every card in every size, for tests and screenshots
// (dev server only: ?screen=cards).
export function CardGallery() {
  const cards = [...useApp((s) => s.cards).values()];
  return (
    <div className="screen gallery">
      <h1>Cards</h1>
      <div className="row" aria-label="Hand size">
        {cards.slice(0, 6).map((c) => (
          <Card key={c.ref} cardRef={c.ref} size="hand" />
        ))}
        <Card size="hand" back="treasure" />
        <Card size="hand" back="loot" />
        <Card size="hand" back="monster" />
      </div>
      <div className="row" aria-label="Table size">
        {cards.slice(0, 12).map((c) => (
          <Card key={c.ref} cardRef={c.ref} size="table" />
        ))}
      </div>
    </div>
  );
}
