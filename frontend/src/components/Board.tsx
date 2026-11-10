import type { ReactNode } from 'react';
import type { Kind } from '../cards';
import { presents, timerText, type Present } from '../challenge';
import { t } from '../i18n';
import { Deck as DeckKind, type CardView, type GameView, type SlotView } from '../protocol';
import { useApp } from '../store';
import { Card } from './Card';
import { Count } from './Count';
import type { Interact } from './Interact';
import { Deck, Discard } from './Pile';
import './Table.css';

// playerN names a player by seat when no names are given.
const playerN = (p: number) => t('table.player_n', { n: p + 1 });

// Board is the middle of the table: one framed board with a fixed shape
// (F-01, H-06). Left to right: the history strip (strip, H-07), the info
// column (children: turn, stack, history), monster and shop slots, bonus souls, and the decks with
// their discards at the right edge (H-16), then the End turn button
// (end, H-04). An empty zone keeps its place, so the board never
// changes size (V-23). A game with the room deck has one more column,
// after the slots: its rooms, landscape, over the room deck and its
// discard (14.13); a game without one never shows it.
// A challenge (14.25, ADR 028 Decision 11) adds its final boss left of
// the monster slots, a minion row under them, and a challenge panel
// instead of the bonus souls: the rule card, the timer and the
// presents. Its board keeps one shape too: the minion row keeps its
// place while it is empty, and shows nothing then.
// The board is also where loot cards are dropped to play them (7.7).
// glow: the game waits for the viewer, so the board glows (15.95).
// name names a player, for a present a player gained.
export function Board({
  view,
  ui,
  glow = false,
  strip,
  end,
  name = playerN,
  children,
}: {
  view: GameView;
  ui: Interact;
  glow?: boolean;
  strip?: ReactNode;
  end?: ReactNode;
  name?: (p: number) => string;
  children?: ReactNode;
}) {
  const rooms = view.rooms ?? [];
  const challenge = view.challenge;
  const minions = view.minions ?? [];
  const monsters = (view.monsters ?? []).map((slot, i) => (
    <Slot key={slotKey(slot, i)} slot={slot} zone={`monster-slot-${String(i)}`} ui={ui} />
  ));
  return (
    <div className="board-cell">
      <section
        className={`board${rooms.length > 0 ? ' has-rooms' : ''}${challenge ? ' has-challenge' : ''}${glow ? ' my-move' : ''}`}
        aria-label={t('board.label')}
        data-drop="play"
      >
        <div className="board-strip">{strip}</div>
        <div className="board-info">{children}</div>
        {challenge ? (
          <div className="board-top">
            <div className="board-slots board-bosses" role="group" aria-label={t('board.final_boss')}>
              {(view.bosses ?? []).map((slot, i) => (
                <Slot key={slotKey(slot, i)} slot={slot} zone={`boss-slot-${String(i)}`} ui={ui} />
              ))}
            </div>
            <div className="board-slots board-monsters" role="group" aria-label={t('board.monsters')}>
              {monsters}
            </div>
          </div>
        ) : (
          <div className="board-slots board-monsters" role="group" aria-label={t('board.monsters')}>
            {monsters}
          </div>
        )}
        {challenge && (
          // The row keeps its place while empty: a minion flies to it
          // before the view draws its slot (minion_arrived).
          <div className="board-slots board-minions" data-zone="minions">
            {minions.length > 0 && (
              <div className="board-minion-row" role="group" aria-label={t('board.minions')}>
                {minions.map((slot, i) => (
                  <Slot key={slotKey(slot, i)} slot={slot} zone={`minion-slot-${String(i)}`} ui={ui} size="small" />
                ))}
              </div>
            )}
          </div>
        )}
        <div className="board-slots board-shop" role="group" aria-label={t('board.shop')}>
          {(view.shop ?? []).map((slot, i) => (
            <Slot key={slotKey(slot, i)} slot={slot} zone={`shop-slot-${String(i)}`} ui={ui} />
          ))}
        </div>
        {rooms.length > 0 && (
          <div className="board-rooms" role="group" aria-label={t('board.rooms')}>
            {rooms.map((slot, i) => (
              <Slot key={slotKey(slot, i)} slot={slot} zone={`room-slot-${String(i)}`} ui={ui} />
            ))}
            <Deck kind="room" count={view.deck_sizes[DeckKind.room] ?? 0} size="small" handle={ui.deck('room')} />
            <Discard kind="room" top={view.discards[DeckKind.room] ?? null} count={view.discard_sizes[DeckKind.room] ?? 0} size="small" />
          </div>
        )}
        {challenge ? (
          <div className="board-challenge" role="group" aria-label={t('challenge.panel')}>
            <Card cardRef={challenge.card.card} size="small" objectId={challenge.card.id} counters={challenge.card.counters} />
            <Timer rounds={challenge.timer} />
            <div className="board-presents" role="group" aria-label={t('challenge.presents')}>
              {presents(view, name).map((p) => (
                <PresentCard key={p.kind} present={p} />
              ))}
            </div>
          </div>
        ) : (
          <div className="bonus-souls" role="group" aria-label={t('board.bonus_souls')}>
            {view.bonus_souls?.map((c) => (
              <Card key={c.id} cardRef={c.card} size="small" counters={c.counters} objectId={c.id} />
            ))}
          </div>
        )}
        <div className="board-decks">
          <DeckPair
            kind="monster"
            count={view.deck_sizes[DeckKind.monster] ?? 0}
            discard={view.discards[DeckKind.monster] ?? null}
            discardCount={view.discard_sizes[DeckKind.monster] ?? 0}
            ui={ui}
            top={view.revealed}
          />
          <DeckPair
            kind="treasure"
            count={view.deck_sizes[DeckKind.treasure] ?? 0}
            discard={view.discards[DeckKind.treasure] ?? null}
            discardCount={view.discard_sizes[DeckKind.treasure] ?? 0}
            ui={ui}
            peek={view.players.find((p) => p.treasure_top)?.treasure_top}
          />
          <DeckPair
            kind="loot"
            count={view.deck_sizes[DeckKind.loot] ?? 0}
            discard={view.discards[DeckKind.loot] ?? null}
            discardCount={view.discard_sizes[DeckKind.loot] ?? 0}
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
  discardCount,
  ui,
  top,
  peek,
}: {
  kind: Kind;
  count: number;
  discard: CardView | null;
  discardCount: number;
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
      <Discard kind={kind} top={discard} count={discardCount} size="small" />
    </div>
  );
}

// A slot's key is its top card; an empty slot keeps its place.
function slotKey(slot: SlotView, i: number): string {
  return slot.top ? `card-${String(slot.top.id)}` : `empty-${String(i)}`;
}

// Slot is one shop, monster, room, final boss or minion slot: its top
// card and how many it covers. zone names the slot for effects that
// place a card on it (19.3). Minions are small cards (14.25).
function Slot({ slot, zone, ui, size = 'table' }: { slot: SlotView; zone: string; ui: Interact; size?: 'table' | 'small' }) {
  const covered = slot.covered?.length ?? 0;
  // A dead card flew off its slot: the slot shows empty until the next
  // view (19.73).
  const lifted = useApp((s) => !!slot.top && s.lifted.includes(slot.top.id));
  if (!slot.top || lifted) {
    return <div className={`pile-empty size-${size}`} data-zone={zone} />;
  }
  return (
    <div className="pile" data-zone={zone}>
      <Card
        cardRef={slot.top.card}
        size={size}
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

// Timer is a challenge's timer (R-COOP-08): a d8 with the rounds left;
// its count flips as it ticks down (count-flip).
function Timer({ rounds }: { rounds: number }) {
  const text = timerText(rounds);
  return (
    <div className="challenge-timer" role="img" aria-label={text} title={text} data-zone="challenge-timer">
      <svg className="challenge-d8" viewBox="0 0 100 100" aria-hidden="true">
        <path className="challenge-d8-face" d="M50 4L94 50L50 96L6 50Z" />
        <path className="challenge-d8-edge" d="M6 50L94 50" />
      </svg>
      <Count value={rounds} className="challenge-d8-value" />
      <span className="challenge-timer-label" aria-hidden="true">
        {t('challenge.timer_label')}
      </span>
    </div>
  );
}

// PresentCard is one present (R-KRAMPUS-05): face down on the treasure
// back with its counters, face up once flipped, or an empty place
// marked given or destroyed. Its kind is written under it. zone lets a
// counter popup or the card's flight start at it.
function PresentCard({ present: p }: { present: Present }) {
  const zone = `present-${p.kind}`;
  let card: ReactNode;
  if (p.state === 'down') {
    card = <Card size="small" back="treasure" label={p.label} counters={[{ count: p.counters }]} />;
  } else if (p.state === 'up' && p.card) {
    card = <Card cardRef={p.card} size="small" />;
  } else {
    card = (
      <div className={`pile-empty size-small present-gone present-${p.state}`} role="img" aria-label={p.label}>
        <span aria-hidden="true">{p.state === 'destroyed' ? '✗' : '✓'}</span>
      </div>
    );
  }
  return (
    <figure className={`present present-${p.state}`} data-zone={zone} title={p.label}>
      {card}
      <figcaption>{p.caption}</figcaption>
    </figure>
  );
}
