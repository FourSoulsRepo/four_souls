import { useRef, useState, type CSSProperties } from 'react';
import { Role, type CardView, type Intent, type PlayerView } from '../protocol';
import { useApp } from '../store';
import type { Position } from '../table';
import { Card, type CardSize } from './Card';
import type { Interact } from './Interact';
import { HiddenHand } from './Pile';
import './Table.css';

interface Props {
  player: PlayerView;
  name: string;
  position: Position;
  you: boolean;
  active: boolean;
  ui: Interact;
}

// counted is how many counters a card carries, of any kind.
function counted(c: CardView): number {
  return (c.counters ?? []).reduce((n, x) => n + x.count, 0);
}

// spent: a character or item that is used up until it recharges.
function spent(c: CardView): boolean {
  return (c.role === Role.character || c.role === Role.item) && !c.charged;
}

// PlayerArea is one player's place at the table: name and counters,
// character and cards in play, and the hand for the viewer (V-02).
export function PlayerArea({ you, position, ...rest }: Props) {
  return you ? <OwnArea {...rest} /> : <OpponentArea position={position} {...rest} />;
}

// OpponentArea is another player's place: a header of stats, then the
// hand's back, the character, and the cards in play.
function OpponentArea({ player, name, position, active, ui }: Omit<Props, 'you'>) {
  const size: CardSize = position === 'left' || position === 'right' ? 'small' : 'table';
  const pick = ui.player(player.id);
  return (
    <section
      className={`player player-${position}${active ? ' active' : ''}${player.dead ? ' dead' : ''}${pick.highlight ? ' hl-pick' : ''}`}
      aria-label={name}
      data-zone={`player-${String(player.id)}`}
    >
      <header className="player-header">
        <b>{name}</b>
        <span className="player-stats">
          HP {player.hp} · ATK {player.atk} · {player.cents}¢ · {player.souls} souls · {player.hand_size} in hand
          {player.dead && ' · dead'}
        </span>
        {active && <Allowance player={player} />}
        {pick.onClick && (
          <button type="button" className="player-pick" onClick={pick.onClick}>
            Choose {name}
          </button>
        )}
      </header>
      <div className="player-cards">
        <HiddenHand playerId={player.id} count={player.hand_size} size={size} />
        {[player.character, ...(player.in_play ?? [])].map((c) => (
          <PlayCard key={c.id} card={c} size={size} ui={ui} />
        ))}
      </div>
    </section>
  );
}

// Coin is the game's coin image (15.38), or a drawn ¢ coin without the
// private images or when the image fails to load.
function Coin() {
  const image = useApp((s) => s.icons.coin);
  const [broken, setBroken] = useState('');
  if (image !== '' && image !== broken) {
    return (
      <img
        className="coin coin-image"
        src={image}
        alt=""
        aria-hidden="true"
        draggable={false}
        onError={() => {
          setBroken(image);
        }}
      />
    );
  }
  return (
    <span className="coin" aria-hidden="true">
      ¢
    </span>
  );
}

// OwnArea is the viewer's place at the bottom edge (15.17): name and
// big coins at the left, the character centered with its live health
// and attack, items to its right, and the hand fanned below, partly
// off screen (F-05, F-06, H-02, H-03, H-15, H-18).
function OwnArea({ player, name, active, ui }: Omit<Props, 'you' | 'position'>) {
  const label = `${name} (you)`;
  const pick = ui.player(player.id);
  const hand = player.hand ?? [];
  const { character } = player;
  return (
    <section
      className={`player player-bottom own${active ? ' active' : ''}${player.dead ? ' dead' : ''}${pick.highlight ? ' hl-pick' : ''}`}
      aria-label={label}
      data-zone={`player-${String(player.id)}`}
    >
      <div className="own-status">
        <b className="own-name">{label}</b>
        <span className="own-coins" title="Coins" aria-label={`${String(player.cents)} cents`}>
          <Coin />
          {player.cents}
        </span>
        <span className="own-souls" title="Souls">
          <b>{player.souls}</b> {player.souls === 1 ? 'soul' : 'souls'}
          {player.dead && ' · dead'}
        </span>
        {active && <Allowance player={player} />}
        {pick.onClick && (
          <button type="button" className="player-pick" onClick={pick.onClick}>
            Choose {label}
          </button>
        )}
      </div>
      <div className="player-cards own-character">
        <Card
          cardRef={character.card}
          size="table"
          damage={character.damage}
          spent={spent(character)}
          counters={counted(character)}
          stats={{ hp: player.hp, atk: player.atk }}
          objectId={character.id}
          {...ui.object(character.id)}
        />
      </div>
      <div className="player-cards own-items" role="group" aria-label="Your items">
        {(player.in_play ?? []).map((c) => (
          <PlayCard key={c.id} card={c} size="table" ui={ui} />
        ))}
      </div>
      <div
        className="hand"
        role="group"
        aria-label="Your hand"
        data-zone={`hand-${String(player.id)}`}
        style={fan(hand.length)}
      >
        {hand.map((c, i) => (
          <HandCard key={c.id} card={c} ui={ui} spread={i - (hand.length - 1) / 2} />
        ))}
      </div>
    </section>
  );
}

// Allowance is what the active player may still do this turn.
function Allowance({ player }: { player: PlayerView }) {
  return (
    <span className="player-allowance" aria-label="Left this turn">
      Loot {player.loot_plays} · Attacks {player.attacks} · Buys {player.purchases}
    </span>
  );
}

// PlayCard is a character or item in play.
function PlayCard({ card, size, ui }: { card: CardView; size: CardSize; ui: Interact }) {
  return (
    <Card
      cardRef={card.card}
      size={size}
      damage={card.damage}
      spent={spent(card)}
      counters={counted(card)}
      objectId={card.id}
      {...ui.object(card.id)}
    />
  );
}

// fan spreads a hand of n cards: a smaller angle and more overlap the
// more cards there are, so a big hand still fits (H-03).
function fan(n: number): CSSProperties {
  const step = Math.min(7, 40 / Math.max(n, 1));
  const overlap = n <= 3 ? 20 : Math.min(64, 20 + (n - 3) * 7);
  return { '--fan-step': `${String(step)}deg`, '--fan-overlap': `${String(overlap)}px` } as CSSProperties;
}

// onPlayArea: the point is over the board, where loot cards are played.
function onPlayArea(x: number, y: number): boolean {
  return document.elementsFromPoint(x, y).some((el) => el.closest('[data-drop="play"]') !== null);
}

// How far a pointer moves before a press becomes a drag, in pixels.
const dragThreshold = 8;

// HandCard is a card in the viewer's hand, spread from the hand's middle
// by its place in the fan. A playable one can be dragged
// onto the board to play it; clicking it opens its actions instead. The
// drag is plain pointer events with pointer capture, so mouse and touch
// share it and the release is never missed (ADR 008).
function HandCard({ card, ui, spread }: { card: CardView; ui: Interact; spread: number }) {
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
      style={
        offset
          ? { transform: `translate(${String(offset.x)}px, ${String(offset.y)}px) scale(1.08)` }
          : ({ '--spread': String(spread) } as CSSProperties)
      }
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
