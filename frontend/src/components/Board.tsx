import type { ReactNode } from 'react';
import { Deck as DeckKind, type GameView, type SlotView } from '../protocol';
import { Card } from './Card';
import type { Interact } from './Interact';
import { Deck, Discard } from './Pile';
import './Table.css';

// Board is the middle of the table: monsters, shop, decks, discards and
// bonus souls (V-01).
// The board is also where loot cards are dropped to play them (7.7).
export function Board({ view, ui, children }: { view: GameView; ui: Interact; children?: ReactNode }) {
  return (
    <section className="board" aria-label="Board" data-drop="play">
      {children}
      <div className="board-row" role="group" aria-label="Monsters">
        <Deck
          kind="monster"
          count={view.deck_sizes[DeckKind.monster] ?? 0}
          size="table"
          handle={ui.deck('monster')}
        />
        {(view.monsters ?? []).map((slot, i) => (
          <Slot key={slotKey(slot, i)} slot={slot} ui={ui} />
        ))}
        <Discard kind="monster" cards={view.discards[DeckKind.monster] ?? null} size="table" />
        {(view.bonus_souls ?? []).length > 0 && (
          <div className="bonus-souls" role="group" aria-label="Bonus souls">
            {view.bonus_souls?.map((c) => (
              <Card key={c.id} cardRef={c.card} size="small" objectId={c.id} />
            ))}
          </div>
        )}
      </div>
      <div className="board-row" role="group" aria-label="Shop">
        <Deck
          kind="treasure"
          count={view.deck_sizes[DeckKind.treasure] ?? 0}
          size="table"
          handle={ui.deck('treasure')}
        />
        {(view.shop ?? []).map((slot, i) => (
          <Slot key={slotKey(slot, i)} slot={slot} ui={ui} />
        ))}
        <Discard kind="treasure" cards={view.discards[DeckKind.treasure] ?? null} size="table" />
      </div>
      <div className="board-side" role="group" aria-label="Loot">
        <Deck kind="loot" count={view.deck_sizes[DeckKind.loot] ?? 0} size="table" />
        <Discard kind="loot" cards={view.discards[DeckKind.loot] ?? null} size="table" />
      </div>
    </section>
  );
}

// A slot's key is its top card; an empty slot keeps its place.
function slotKey(slot: SlotView, i: number): string {
  return slot.top ? `card-${String(slot.top.id)}` : `empty-${String(i)}`;
}

// Slot is one shop or monster slot: its top card and how many it covers.
function Slot({ slot, ui }: { slot: SlotView; ui: Interact }) {
  const covered = slot.covered?.length ?? 0;
  if (!slot.top) {
    return <div className="pile-empty size-table" />;
  }
  return (
    <div className="pile">
      <Card
        cardRef={slot.top.card}
        size="table"
        damage={slot.top.damage}
        objectId={slot.top.id}
        {...ui.object(slot.top.id)}
      />
      {covered > 0 && <span className="pile-count">+{covered}</span>}
    </div>
  );
}
