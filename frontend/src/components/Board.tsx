import type { ReactNode } from 'react';
import type { Kind } from '../cards';
import { t } from '../i18n';
import { Deck as DeckKind, type CardView, type GameView, type SlotView } from '../protocol';
import { useApp } from '../store';
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
// glow: the game waits for the viewer, so the board glows (15.95).
export function Board({
  view,
  ui,
  glow = false,
  strip,
  end,
  children,
}: {
  view: GameView;
  ui: Interact;
  glow?: boolean;
  strip?: ReactNode;
  end?: ReactNode;
  children?: ReactNode;
}) {
  return (
    <div className="board-cell">
      <section className={`board${glow ? ' my-move' : ''}`} aria-label={t('board.label')} data-drop="play">
        <div className="board-strip">{strip}</div>
        <div className="board-info">{children}</div>
        <div className="board-slots board-monsters" role="group" aria-label={t('board.monsters')}>
          {(view.monsters ?? []).map((slot, i) => (
            <Slot key={slotKey(slot, i)} slot={slot} zone={`monster-slot-${String(i)}`} ui={ui} />
          ))}
        </div>
        <div className="board-slots board-shop" role="group" aria-label={t('board.shop')}>
          {(view.shop ?? []).map((slot, i) => (
            <Slot key={slotKey(slot, i)} slot={slot} zone={`shop-slot-${String(i)}`} ui={ui} />
          ))}
        </div>
        <div className="bonus-souls" role="group" aria-label={t('board.bonus_souls')}>
          {view.bonus_souls?.map((c) => (
            <Card key={c.id} cardRef={c.card} size="small" counters={c.counters} objectId={c.id} />
          ))}
        </div>
        <div className="board-decks">
          <DeckPair
            kind="monster"
            count={view.deck_sizes[DeckKind.monster] ?? 0}
            discard={view.discards[DeckKind.monster] ?? null}
            ui={ui}
            top={view.revealed}
          />
          <DeckPair
            kind="treasure"
            count={view.deck_sizes[DeckKind.treasure] ?? 0}
            discard={view.discards[DeckKind.treasure] ?? null}
            ui={ui}
            peek={view.players.find((p) => p.treasure_top)?.treasure_top}
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

// DeckPair is a deck and its discard, side by side. top is a drawn card
// waiting on the deck for its slot (19.9); it leaves the deck once a
// place effect flies it to the slot. peek is the top card only the
// viewer may look at, e.g. under There's Options (19.30).
function DeckPair({
  kind,
  count,
  discard,
  ui,
  top,
  peek,
}: {
  kind: Kind;
  count: number;
  discard: CardView[] | null;
  ui: Interact;
  top?: CardView | undefined;
  peek?: CardView | undefined;
}) {
  const leaving = useApp(
    (s) =>
      !!top &&
      ((s.playing?.effect.kind === 'place' && s.playing.effect.card === top.card) ||
        s.placed.some((p) => p.card === top.card)),
  );
  return (
    <div className="deck-pair">
      <Deck kind={kind} count={count} size="small" handle={ui.deck(kind)} top={leaving ? undefined : top} peek={peek} />
      <Discard kind={kind} cards={discard} size="small" />
    </div>
  );
}

// A slot's key is its top card; an empty slot keeps its place.
function slotKey(slot: SlotView, i: number): string {
  return slot.top ? `card-${String(slot.top.id)}` : `empty-${String(i)}`;
}

// Slot is one shop or monster slot: its top card and how many it covers.
// zone names the slot for effects that place a card on it (19.3).
function Slot({ slot, zone, ui }: { slot: SlotView; zone: string; ui: Interact }) {
  const covered = slot.covered?.length ?? 0;
  // A dead card flew off its slot: the slot shows empty until the next
  // view (19.73).
  const lifted = useApp((s) => !!slot.top && s.lifted.includes(slot.top.id));
  if (!slot.top || lifted) {
    return <div className="pile-empty size-table" data-zone={zone} />;
  }
  return (
    <div className="pile" data-zone={zone}>
      <Card
        cardRef={slot.top.card}
        size="table"
        damage={slot.top.damage}
        counters={slot.top.counters}
        monster={slot.top.stats}
        objectId={slot.top.id}
        {...ui.object(slot.top.id)}
      />
      {covered > 0 && <span className="pile-count">+{covered}</span>}
    </div>
  );
}
