import type { ReactNode } from 'react';
import type { Kind } from '../cards';
import { Deck as DeckKind, type CardView, type GameView, type SlotView } from '../protocol';
import { Card } from './Card';
import type { Interact } from './Interact';
import { Deck, Discard } from './Pile';
import './Table.css';

// Board is the middle of the table: one framed board with a fixed shape
// (F-01, H-06). Left to right: the history strip (strip, H-07), the info
// column (children: turn, stack, history), monster and shop slots, bonus souls, and the decks with
// their discards at the right edge (H-16), then the End turn button
// (end, H-04). An empty zone keeps its place, so the board never
// changes size (V-23).
// The board is also where loot cards are dropped to play them (7.7).
export function Board({
  view,
  ui,
  strip,
  end,
  children,
}: {
  view: GameView;
  ui: Interact;
  strip?: ReactNode;
  end?: ReactNode;
  children?: ReactNode;
}) {
  return (
    <div className="board-cell">
      <section className="board" aria-label="Board" data-drop="play">
        <div className="board-strip">{strip}</div>
        <div className="board-info">{children}</div>
        <div className="board-slots board-monsters" role="group" aria-label="Monsters">
          {(view.monsters ?? []).map((slot, i) => (
            <Slot key={slotKey(slot, i)} slot={slot} ui={ui} />
          ))}
        </div>
        <div className="board-slots board-shop" role="group" aria-label="Shop">
          {(view.shop ?? []).map((slot, i) => (
            <Slot key={slotKey(slot, i)} slot={slot} ui={ui} />
          ))}
        </div>
        <div className="bonus-souls" role="group" aria-label="Bonus souls">
          {view.bonus_souls?.map((c) => (
            <Card key={c.id} cardRef={c.card} size="small" objectId={c.id} />
          ))}
        </div>
        <div className="board-decks">
          <DeckPair
            kind="monster"
            count={view.deck_sizes[DeckKind.monster] ?? 0}
            discard={view.discards[DeckKind.monster] ?? null}
            ui={ui}
          />
          <DeckPair
            kind="treasure"
            count={view.deck_sizes[DeckKind.treasure] ?? 0}
            discard={view.discards[DeckKind.treasure] ?? null}
            ui={ui}
          />
          <DeckPair
            kind="loot"
            count={view.deck_sizes[DeckKind.loot] ?? 0}
            discard={view.discards[DeckKind.loot] ?? null}
            ui={ui}
          />
        </div>
        <div className="board-end">{end}</div>
      </section>
    </div>
  );
}

// DeckPair is a deck and its discard, side by side.
function DeckPair({ kind, count, discard, ui }: { kind: Kind; count: number; discard: CardView[] | null; ui: Interact }) {
  return (
    <div className="deck-pair">
      <Deck kind={kind} count={count} size="small" handle={ui.deck(kind)} />
      <Discard kind={kind} cards={discard} size="small" />
    </div>
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
