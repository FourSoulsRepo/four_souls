import { useRef, useState } from 'react';
import { Role, type CardView, type Intent, type PlayerView } from '../protocol';
import { useApp } from '../store';
import type { Position } from '../table';
import { Card, type CardSize } from './Card';
import type { Interact } from './Interact';
import './Table.css';

interface Props {
  player: PlayerView;
  name: string;
  position: Position;
  you: boolean;
  active: boolean;
  ui: Interact;
}

// spent: a character or item that is used up until it recharges.
function spent(c: CardView): boolean {
  return (c.role === Role.character || c.role === Role.item) && !c.charged;
}

// PlayerArea is one player's place at the table: name and counters,
// character and cards in play, and the hand for the viewer (V-02).
export function PlayerArea({ player, name, position, you, active, ui }: Props) {
  const size: CardSize = position === 'left' || position === 'right' ? 'small' : 'table';
  const label = you ? `${name} (you)` : name;
  const pick = ui.player(player.id);
  return (
    <section
      className={`player player-${position}${active ? ' active' : ''}${player.dead ? ' dead' : ''}${pick.highlight ? ' hl-pick' : ''}`}
      aria-label={label}
      data-zone={`player-${String(player.id)}`}
    >
      <header className="player-header">
        <b>{label}</b>
        <span className="player-stats">
          HP {player.hp} · ATK {player.atk} · {player.cents}¢ · {player.souls} souls
          {!you && ` · ${String(player.hand_size)} in hand`}
          {player.dead && ' · dead'}
        </span>
        {pick.onClick && (
          <button type="button" className="player-pick" onClick={pick.onClick}>
            Choose {label}
          </button>
        )}
      </header>
      <div className="player-cards">
        {[player.character, ...(player.in_play ?? [])].map((c) => (
          <Card
            key={c.id}
            cardRef={c.card}
            size={size}
            damage={c.damage}
            spent={spent(c)}
            objectId={c.id}
            {...ui.object(c.id)}
          />
        ))}
      </div>
      {you && (
        <div className="hand" role="group" aria-label="Your hand" data-zone={`hand-${String(player.id)}`}>
          {(player.hand ?? []).map((c) => (
            <HandCard key={c.id} card={c} ui={ui} />
          ))}
        </div>
      )}
    </section>
  );
}

// onPlayArea: the point is over the board, where loot cards are played.
function onPlayArea(x: number, y: number): boolean {
  return document.elementsFromPoint(x, y).some((el) => el.closest('[data-drop="play"]') !== null);
}

// How far a pointer moves before a press becomes a drag, in pixels.
const dragThreshold = 8;

// HandCard is a card in the viewer's hand. A playable one can be dragged
// onto the board to play it; clicking it opens its actions instead. The
// drag is plain pointer events with pointer capture, so mouse and touch
// share it and the release is never missed (ADR 008).
function HandCard({ card, ui }: { card: CardView; ui: Interact }) {
  const act = useApp((s) => s.act);
  const setZoom = useApp((s) => s.setZoom);
  const [offset, setOffset] = useState<{ x: number; y: number } | null>(null);
  const startRef = useRef<{ x: number; y: number; pointer: number } | null>(null);
  // draggedRef swallows the click that ends a drag.
  const draggedRef = useRef(false);
  const play: Intent | undefined = ui.play(card.id);
  const handle = ui.object(card.id);
  const reset = () => {
    startRef.current = null;
    setOffset(null);
  };
  return (
    <div
      className={`hand-card${play ? ' playable' : ''}${offset ? ' dragging' : ''}`}
      style={offset ? { transform: `translate(${String(offset.x)}px, ${String(offset.y)}px) scale(1.08)` } : undefined}
      onPointerDown={(e) => {
        if (play && e.button === 0) {
          startRef.current = { x: e.clientX, y: e.clientY, pointer: e.pointerId };
          draggedRef.current = false;
        }
      }}
      onPointerMove={(e) => {
        const start = startRef.current;
        if (start?.pointer !== e.pointerId) {
          return;
        }
        const x = e.clientX - start.x;
        const y = e.clientY - start.y;
        if (!offset) {
          if (Math.hypot(x, y) < dragThreshold) {
            return;
          }
          e.currentTarget.setPointerCapture(e.pointerId);
          draggedRef.current = true;
          setZoom(null);
        }
        setOffset({ x, y });
      }}
      onPointerUp={(e) => {
        if (offset && play && onPlayArea(e.clientX, e.clientY)) {
          act(play);
        }
        reset();
      }}
      onPointerCancel={reset}
      onClickCapture={(e) => {
        if (draggedRef.current) {
          draggedRef.current = false;
          e.stopPropagation();
        }
      }}
    >
      <Card cardRef={card.card} size="hand" objectId={card.id} {...handle} />
    </div>
  );
}
