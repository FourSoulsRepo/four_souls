import { useRef, useState, type CSSProperties, type ReactNode } from 'react';
import { counted, Role, type CardView, type PlayerView } from '../protocol';
import { useApp } from '../store';
import type { LinePart } from '../stack';
import { isCurse, isSoul, playerColor, soulCards, soulSlots, type Position } from '../table';
import { Card, useZoom, type CardSize } from './Card';
import type { Interact } from './Interact';
import { HiddenHand, RevealedHand } from './Pile';
import { Part } from './StackPanel';
import './Table.css';

interface Props {
  player: PlayerView;
  name: string;
  position: Position;
  you: boolean;
  active: boolean;
  // asked: the game waits for this player now (15.95).
  asked: boolean;
  // winSouls is the soul value that wins; 0 when the view has none (V-30).
  winSouls: number;
  ui: Interact;
  // tools are the viewer's turn actions, shown by own name (15.19).
  tools?: ReactNode;
  // doing is what this active player does now, for the others (15.30).
  doing?: LinePart[] | null;
}

// spent: a character or item that is used up until it recharges.
function spent(c: CardView): boolean {
  return (c.role === Role.character || c.role === Role.item) && !c.charged;
}

// PlayerArea is one player's place at the table: name and counters,
// character and cards in play, and the hand for the viewer (V-02).
export function PlayerArea({ you, position, tools, ...rest }: Props) {
  return you ? <OwnArea tools={tools} {...rest} /> : <OpponentArea position={position} {...rest} />;
}

// OpponentArea is another player's place (15.18): a small name with
// coins, the souls and the character with its live health and attack, the
// items, and the hand as fanned backs (F-07, F-08, H-17, H-21), or face
// up when the view sends it (19.5). Curses stand at the far left
// (15.55). At the top it mirrors the own area, the hand hanging past the top edge; in
// a side column everything stacks in a narrow column (F-09).
function OpponentArea({ player, name, position, active, asked, winSouls, ui, doing }: Omit<Props, 'you' | 'tools'>) {
  const size: CardSize = position === 'left' || position === 'right' ? 'small' : 'table';
  const pick = ui.player(player.id);
  const { character } = player;
  return (
    <section
      className={`player player-${position} opponent${active ? ' active' : ''}${player.dead ? ' dead' : ''}${pick.highlight ? ` hl-${pick.highlight}` : ''}`}
      aria-label={name}
      data-zone={`player-${String(player.id)}`}
    >
      {player.hand ? (
        <RevealedHand playerId={player.id} cards={player.hand} size="small" />
      ) : (
        <HiddenHand playerId={player.id} count={player.hand_size} size="small" />
      )}
      <div className="opp-left">
        <Curses player={player} label={`${name}: curses`} size={size} ui={ui} />
        <div className="opp-status">
          <b className="opp-name">
            <PlayerDot player={player.id} />
            {name}
            {asked && <Hourglass player={player.id} />}
          </b>
          <span className="opp-line">
            <span className="opp-coins" title="Coins" aria-label={`${String(player.cents)} cents`}>
              <Coin />
              {player.cents}
            </span>
            {player.dead && <span className="player-dead">dead</span>}
          </span>
          {active && doing && <Activity parts={doing} />}
          {active && <Allowance player={player} />}
          {pick.onClick && (
            <button type="button" className="player-pick" onClick={pick.onClick}>
              Choose {name}
            </button>
          )}
        </div>
      </div>
      <div className="player-cards opp-character">
        <Souls player={player} winSouls={winSouls} ui={ui} />
        <Card
          cardRef={character.card}
          size={size}
          damage={character.damage}
          spent={spent(character)}
          counters={counted(character)}
          stats={statsOf(player)}
          objectId={character.id}
          player={player.id}
          {...ui.character(player.id, character.id)}
        />
      </div>
      <div className="player-cards opp-items" role="group" aria-label={`${name}: items`}>
        {(player.in_play ?? [])
          .filter((c) => !isSoul(c) && !isCurse(c))
          .map((c) => (
            <PlayCard key={c.id} card={c} owner={player.id} size={size} ui={ui} />
          ))}
      </div>
    </section>
  );
}

// Coin is the game's coin image (15.38), or a drawn ¢ coin without the
// private images or when the image fails to load.
export function Coin() {
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

// OwnArea is the viewer's place at the bottom edge (15.17): name, big
// coins and the turn actions at the left, any curses at the far left
// (15.55), the souls and the character centered with its live health
// and attack, items to its right, and the hand fanned below, partly
// off screen (F-05, F-06, H-02, H-03, H-15, H-18).
function OwnArea({ player, name, active, asked, winSouls, ui, tools }: Omit<Props, 'you' | 'position'>) {
  const label = `${name} (you)`;
  const pick = ui.player(player.id);
  const hand = player.hand ?? [];
  const { character } = player;
  return (
    <section
      className={`player player-bottom own${active ? ' active' : ''}${player.dead ? ' dead' : ''}${pick.highlight ? ` hl-${pick.highlight}` : ''}`}
      aria-label={label}
      data-zone={`player-${String(player.id)}`}
    >
      <div className="own-left">
        <Curses player={player} label="Your curses" size="table" ui={ui} />
        <div className="own-status">
          <b className="own-name">
            <PlayerDot player={player.id} />
            {label}
            {asked && <Hourglass player={player.id} />}
          </b>
          <span className="own-coins" title="Coins" aria-label={`${String(player.cents)} cents`}>
            <Coin />
            {player.cents}
          </span>
          {player.dead && <span className="player-dead">dead</span>}
          {active && (
            <span className="player-allowance" aria-label="Loot plays left">
              Loot {player.loot_plays}
            </span>
          )}
          {pick.onClick && (
            <button type="button" className="player-pick" onClick={pick.onClick}>
              Choose {label}
            </button>
          )}
        </div>
        {tools}
      </div>
      <div className="player-cards own-character">
        <Souls player={player} winSouls={winSouls} ui={ui} />
        <Card
          cardRef={character.card}
          size="table"
          damage={character.damage}
          spent={spent(character)}
          counters={counted(character)}
          stats={statsOf(player)}
          objectId={character.id}
          player={player.id}
          {...ui.character(player.id, character.id)}
        />
      </div>
      <div className="player-cards own-items" role="group" aria-label="Your items">
        {(player.in_play ?? [])
          .filter((c) => !isSoul(c) && !isCurse(c))
          .map((c) => (
            <PlayCard key={c.id} card={c} owner={player.id} size="table" ui={ui} />
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

// Souls are a player's soul crystals left of the character (V-30,
// V-31): filled slots are gained souls, empty ones what is left to
// win. They fill from the bottom up. Each filled one holds its soul
// card, which is not on the table (V-43).
function Souls({ player, winSouls, ui }: { player: PlayerView; winSouls: number; ui: Interact }) {
  const cards = soulCards(player);
  const slots = soulSlots(Math.max(player.souls, cards.length), winSouls);
  const label = `Souls: ${String(player.souls)} of ${String(slots.length)}`;
  return (
    <div className="souls" role="group" aria-label={label} title={label}>
      {slots.map((filled, i) => {
        const card = cards.at(i);
        return (
          // A slot is only its place in the row, so the index is its key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <SoulCrystal key={i} filled={filled} card={card} own={card?.id === player.character.id} ui={ui} />
        );
      })}
    </div>
  );
}

// SoulCrystal is one soul slot. Hovering a filled one, or a long touch,
// zooms its card, a killed monster or a bonus soul (V-43); an effect
// that aims at the soul aims at the crystal. A character's own soul
// value (own) only zooms: the character card stays its zone.
function SoulCrystal({ filled, card, own, ui }: { filled: boolean; card?: CardView; own: boolean; ui: Interact }) {
  const soul = card && !own ? card : undefined;
  const pick = soul ? ui.object(soul.id) : {};
  const zoom = useZoom(card?.card, card !== undefined);
  const name = useApp((s) => (card ? (s.cards.get(card.card)?.name ?? card.card) : ''));
  return (
    <span
      className={`soul-slot${filled ? ' filled' : ''}${pick.highlight ? ` hl-${pick.highlight}` : ''}`}
      data-zone={soul ? `obj-${String(soul.id)}` : undefined}
      role={pick.onClick ? 'button' : card ? 'img' : undefined}
      aria-label={card ? `Soul: ${name}` : undefined}
      // An empty title hides the row's tooltip: the zoom shows instead.
      title={card ? '' : undefined}
      onClick={pick.onClick}
      {...zoom}
    />
  );
}

// Curses are the curses a player carries, apart from the items, at the
// far left of the player (V-42). None, no group: the area keeps its
// size (V-23).
function Curses({ player, label, size, ui }: { player: PlayerView; label: string; size: CardSize; ui: Interact }) {
  const curses = (player.in_play ?? []).filter(isCurse);
  if (curses.length === 0) {
    return null;
  }
  return (
    <div className="player-cards curses" role="group" aria-label={label}>
      {curses.map((c) => (
        <PlayCard key={c.id} card={c} owner={player.id} size={size} ui={ui} />
      ))}
    </div>
  );
}

// Activity is what the active player does now, as the stack words it:
// "Ann attacks Pooter" (15.30, V-21).
function Activity({ parts }: { parts: LinePart[] }) {
  return (
    <span className="player-activity" role="status" aria-label={parts.map((p) => p.text).join('')}>
      {parts.map((p, i) => (
        // A line's parts never change order, so the index is a stable key.
        // eslint-disable-next-line @eslint-react/no-array-index-key
        <Part key={i} part={p} />
      ))}
    </span>
  );
}

// Allowance is what the active player may still do this turn.
function Allowance({ player }: { player: PlayerView }) {
  return (
    <span className="player-allowance" aria-label="Left this turn">
      <span>Loot {player.loot_plays}</span> · <span>Attacks {player.attacks}</span> ·{' '}
      <span>Shop {player.purchases}</span>
    </span>
  );
}

// statsOf is a character's live stats, its print and their sources
// (15.92, V-41).
function statsOf(player: PlayerView) {
  return { hp: player.hp, atk: player.atk, printed_hp: player.printed_hp, printed_atk: player.printed_atk, sources: player.sources };
}

// PlayCard is a character or item in play; owner's dice look draws the
// numbers it keeps (15.67).
function PlayCard({ card, owner, size, ui }: { card: CardView; owner: number; size: CardSize; ui: Interact }) {
  return (
    <Card
      cardRef={card.card}
      size={size}
      damage={card.damage}
      spent={spent(card)}
      counters={counted(card)}
      dice={card.dice}
      owner={owner}
      monster={card.stats}
      objectId={card.id}
      {...ui.object(card.id)}
    />
  );
}

// fan spreads a hand of n cards: a smaller angle and more overlap the
// more cards there are, so a big hand still fits (H-03). The ends drop
// at most fanDrop pixels, so their tops stay in reach (19.15).
function fan(n: number): CSSProperties {
  const step = Math.min(7, 40 / Math.max(n, 1));
  const overlap = n <= 3 ? 20 : Math.min(64, 20 + (n - 3) * 7);
  const half = Math.max((n - 1) / 2, 1);
  const drop = Math.min(4, fanDrop / (half * half));
  return {
    '--fan-step': `${String(step)}deg`,
    '--fan-overlap': `${String(overlap)}px`,
    '--fan-drop': `${String(drop)}px`,
  } as CSSProperties;
}

// fanDrop is how far the end cards of a big hand drop, in pixels.
const fanDrop = 24;

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
// A card that needs a target is not dragged: a click starts aiming and
// the card stays lifted in the fan meanwhile (15.25).
function HandCard({ card, ui, spread }: { card: CardView; ui: Interact; spread: number }) {
  const act = useApp((s) => s.act);
  const setZoom = useApp((s) => s.setZoom);
  const aiming = useApp((s) => s.aiming?.source === card.id);
  const [offset, setOffset] = useState<{ x: number; y: number } | null>(null);
  const startRef = useRef<{ x: number; y: number; pointer: number } | null>(null);
  // draggedRef swallows the click that ends a drag.
  const draggedRef = useRef(false);
  const action = ui.play(card.id);
  const play = action && !action.targets?.length ? action.intent : undefined;
  const handle = ui.object(card.id);
  const reset = () => {
    startRef.current = null;
    setOffset(null);
  };
  return (
    <div
      className={`hand-card${play ? ' playable' : ''}${handle.highlight === 'act' ? ' usable' : ''}${offset ? ' dragging' : ''}${aiming ? ' aiming' : ''}`}
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
      <Card cardRef={card.card} size="hand" objectId={card.id} rise {...handle} />
    </div>
  );
}

// PlayerDot is the player's own color by their name; the stack names
// them in it (15.21, S-03).
function PlayerDot({ player }: { player: number }) {
  return <span className="player-dot" style={{ background: playerColor(player) }} aria-hidden="true" />;
}

// Hourglass marks the player the game waits for, on every screen
// (15.95, V-52). A shape in the player's color, not a hue alone, so
// color-blind players read it too.
function Hourglass({ player }: { player: number }) {
  return (
    <svg className="player-asked" viewBox="0 0 12 16" role="img" aria-label="The game waits for this player" style={{ color: playerColor(player) }}>
      <path d="M1 1h10M1 15h10M2 1c0 4 3 5 3 7s-3 3-3 7M10 1c0 4-3 5-3 7s3 3 3 7" />
      <path className="player-asked-sand" d="M3.5 4h5L6 7zM3 14.5c0-2 3-3 3-3s3 1 3 3z" />
    </svg>
  );
}
