import { motion } from 'motion/react';
import { useEffect, useRef, useState, type MouseEvent, type PointerEvent } from 'react';
import { claimed } from '../animations';
import { iconText, kindOf, textParts, unknownCard, type CardInfo, type Kind } from '../cards';
import { scaleFor } from '../playback';
import type { MonsterStats } from '../protocol';
import { useApp } from '../store';
import { boxOf } from '../zoom';
import './Card.css';

// How long a touch must rest on a card to zoom it (B-06).
const LONG_PRESS_MS = 450;

export type CardSize = 'hand' | 'table' | 'small' | 'zoom' | 'choice';

interface Props {
  // cardRef is the card to show; undefined shows a face-down card.
  cardRef?: string;
  size: CardSize;
  // back is the deck of a face-down card.
  back?: Kind;
  // damage is shown as a counter on the card.
  damage?: number;
  // spent dims a used (uncharged) card.
  spent?: boolean;
  // objectId is the card's object on the table; effects find it by it,
  // and a new one pops in, unless an effect just brought it (19.6).
  objectId?: number;
  // player marks a character card as its player: aiming at the player
  // picks the character (15.33).
  player?: number;
  // counters is how many counters the card carries (7.9).
  counters?: number;
  // highlight marks a card the player may act on or pick (7.7).
  highlight?: 'act' | 'pick' | 'aim';
  // stats shows live health and attack as big numbers (H-15).
  stats?: { hp: number; atk: number };
  // monster shows the stats that differ from the print (C-03).
  monster?: MonsterStats;
  // noZoom: the card does not zoom itself; its row does (15.21).
  noZoom?: boolean;
  // rise: the zoom rises from the card, as from a hand (15.27, H-13).
  rise?: boolean;
  // why says why the player can't use the card now: its zoom says it,
  // and so does a click on it (15.29, F-17, F-18).
  why?: string;
  onClick?: () => void;
}

// Card is the one card component (7.4): the card image, or a blank card
// with its text when the image is missing (A-11).
export function Card({
  cardRef,
  size,
  back = 'loot',
  damage = 0,
  spent = false,
  counters = 0,
  objectId,
  player,
  highlight,
  stats,
  monster,
  noZoom = false,
  rise = false,
  why,
  onClick,
}: Props) {
  const info = useApp((s) => (cardRef ? s.cards.get(cardRef) : undefined));
  const backImage = useApp((s) => s.backs.get(back) ?? '');
  const scale = useApp((s) => scaleFor(s.settings.animation, 0, 0));
  const [brokenImage, setBrokenImage] = useState('');
  const zoomable = size !== 'zoom' && !noZoom;
  const zoomHandlers = useZoom(cardRef, zoomable, rise, why);
  const showWhy = useApp((s) => s.showWhy);
  const click =
    onClick ??
    (why
      ? (e: MouseEvent) => {
          showWhy(why, e.clientX, e.clientY);
        }
      : undefined);

  if (cardRef === undefined) {
    return (
      <div
        className={`card size-${size} kind-${back}${highlight ? ` hl-${highlight}` : ''}`}
        role={onClick ? 'button' : 'img'}
        aria-label={`Face-down ${back} card`}
        onClick={onClick}
      >
        {backImage !== '' && backImage !== brokenImage ? (
          <img
            src={backImage}
            alt="" // the card itself carries the label
            draggable={false}
            onError={() => {
              setBrokenImage(backImage);
            }}
          />
        ) : (
          <div className="card-back-face">{back}</div>
        )}
      </div>
    );
  }

  const card = info ?? unknownCard(cardRef);
  const image = card.image !== '' && card.image !== brokenImage;
  // The zoomed card shows no badges: its words stay readable.
  const badges = size !== 'zoom';

  const pop = objectId !== undefined && scale > 0 && !claimed(card.ref, 'arrive');
  return (
    // anim: card-pop
    <motion.div
      initial={pop ? { scale: 0.6, opacity: 0 } : false}
      animate={{ scale: 1, opacity: 1 }}
      transition={{ duration: 0.25 * scale }}
      data-zone={objectId === undefined ? undefined : `obj-${String(objectId)}`}
      data-player={player}
      className={`card size-${size} kind-${kindOf(card.type)}${image ? '' : ' blank'}${counters > 0 ? ' has-counters' : ''}${spent ? ' spent' : ''}${highlight ? ` hl-${highlight}` : ''}`}
      role={onClick ? 'button' : 'img'}
      aria-label={card.name}
      onClick={click}
      {...zoomHandlers}
    >
      {image ? (
        <img
          src={card.image}
          alt="" // the card itself carries the name
          draggable={false}
          onError={() => {
            setBrokenImage(card.image);
          }}
        />
      ) : (
        <BlankFace card={card} full={size === 'zoom'} />
      )}
      {badges && counters > 0 && <Badge kind="counters" value={counters} label="counters" />}
      {/* Damage shows as the lower, red health; the mark is for the rest. */}
      {badges && damage > 0 && !monster && !stats && <Badge kind="damage" value={damage} label="damage" />}
      {badges && stats && (
        <>
          <Badge kind="hp" value={stats.hp} label="health" title="Health" hurt={damage > 0} />
          <Badge kind="atk" value={stats.atk} label="attack" title="Attack" />
        </>
      )}
      {badges && monster && <MonsterBadges stats={monster} hurt={damage > 0} />}
    </motion.div>
  );
}

// useZoom zooms a card from the element it is spread on: on mouse
// hover, beside it or rising from it, or after a long touch (B-06). A
// card and a soul crystal share it (V-43). Nothing zooms without a
// card ref or when zoomable is false.
export function useZoom(cardRef: string | undefined, zoomable: boolean, rise = false, why?: string) {
  const setZoom = useApp((s) => s.setZoom);
  const pressRef = useRef(0);
  const downRef = useRef<{ x: number; y: number } | null>(null);
  const on = zoomable && cardRef !== undefined;

  // A card that leaves the screen while hovered gets no pointer-leave:
  // the card that started a mouse zoom ends it then.
  const zoomedRef = useRef(false);
  useEffect(
    () => () => {
      const { zoom, setZoom: clear } = useApp.getState();
      if (zoomedRef.current && zoom && zoom.ref === cardRef && !zoom.touch) {
        clear(null);
      }
    },
    [cardRef],
  );

  const cancelPress = () => {
    window.clearTimeout(pressRef.current);
    downRef.current = null;
  };
  return {
    onPointerEnter: (e: PointerEvent<HTMLElement>) => {
      if (on && e.pointerType === 'mouse') {
        zoomedRef.current = true;
        setZoom({ ref: cardRef, touch: false, at: boxOf(e.currentTarget), rise, ...(why ? { why } : {}) });
      }
    },
    onPointerLeave: (e: PointerEvent<HTMLElement>) => {
      cancelPress();
      if (on && e.pointerType === 'mouse') {
        zoomedRef.current = false;
        setZoom(null);
      }
    },
    onPointerDown: (e: PointerEvent<HTMLElement>) => {
      if (on && e.pointerType !== 'mouse') {
        const at = boxOf(e.currentTarget);
        downRef.current = { x: e.clientX, y: e.clientY };
        pressRef.current = window.setTimeout(() => {
          setZoom({ ref: cardRef, touch: true, at, ...(why ? { why } : {}) });
        }, LONG_PRESS_MS);
      }
    },
    onPointerMove: (e: PointerEvent<HTMLElement>) => {
      // A finger that moves is dragging, not pressing.
      const d = downRef.current;
      if (d && Math.hypot(e.clientX - d.x, e.clientY - d.y) > 8) {
        cancelPress();
      }
    },
    onPointerUp: cancelPress,
    onPointerCancel: cancelPress,
    onContextMenu: (e: MouseEvent<HTMLElement>) => {
      e.preventDefault(); // a long press must not open the menu
    },
  };
}

// MonsterBadges shows a monster's health, dice and attack, each only
// while it differs from the print (C-03): the card shows the rest.
// Health also shows, in red, while the monster is hurt (T-25).
function MonsterBadges({ stats, hurt }: { stats: MonsterStats; hurt: boolean }) {
  const printed = (n: number) => `printed ${String(n)}`;
  return (
    <>
      {(stats.hp !== stats.printed_hp || hurt) && (
        <Badge kind="hp" value={stats.hp} label="health" title={`Health (${printed(stats.printed_hp)})`} hurt={hurt} />
      )}
      {stats.dc !== stats.printed_dc && (
        <Badge kind="dc" value={stats.dc} label="dice" title={`Dice to hit (${printed(stats.printed_dc)})`} />
      )}
      {stats.atk !== stats.printed_atk && (
        <Badge kind="atk" value={stats.atk} label="attack" title={`Attack (${printed(stats.printed_atk)})`} />
      )}
    </>
  );
}

// Badge shapes, drawn like the game's cut-out art (15.23): a heart for
// health, a die for the dice, and lumpy blobs for the rest.
const blob = 'M52 6C78 5 95 24 94 50C93 77 74 95 48 94C22 93 5 76 6 49C7 22 26 7 52 6Z';
type BadgeKind = 'hp' | 'dc' | 'atk' | 'counters' | 'damage';
const shapes = new Map<BadgeKind, string>([
  ['hp', 'M50 92C24 72 5 55 6 33C7 15 22 6 35 7C43 8 48 13 50 20C52 13 57 8 65 7C78 6 93 15 94 33C95 55 76 72 50 92Z'],
  ['dc', 'M22 8L78 6C88 6 94 12 94 22L93 78C93 88 87 94 77 94L23 93C13 93 7 87 7 77L6 23C6 13 12 8 22 8Z'],
]);

interface BadgeProps {
  kind: BadgeKind;
  value: number;
  label: string;
  title?: string;
  hurt?: boolean;
}

// Badge is one number on a card: its shape in ink with a hard shadow,
// the number in the game's hand-drawn letters. A hurt health number is
// red while below full (T-25).
function Badge({ kind, value, label, title, hurt = false }: BadgeProps) {
  return (
    <span className={`card-badge card-${kind}${hurt ? ' hurt' : ''}`} title={title} aria-label={`${String(value)} ${label}`}>
      <svg className="card-badge-shape" viewBox="0 0 100 100" aria-hidden="true">
        <path d={shapes.get(kind) ?? blob} />
        <path className="card-badge-shine" d="M30 26C36 18 46 16 52 18C44 22 38 28 34 36C31 34 29 30 30 26Z" />
      </svg>
      <span className="card-badge-value">{value}</span>
    </span>
  );
}

// BlankFace shows a card's text; full adds the effects. The zoom shows
// the keywords beside the card (15.27).
function BlankFace({ card, full }: { card: CardInfo; full: boolean }) {
  return (
    <div className="card-blank">
      <div className="card-name">{card.name}</div>
      <div className="card-type">{card.type}</div>
      {full && (
        <div className="card-text">
          {card.effects.map((effect) => (
            <p key={effect}>
              <CardText text={effect} />
            </p>
          ))}
        </div>
      )}
      {card.rewards.length > 0 && (
        <div className="card-rewards">{card.rewards.map((r) => `${r.count} ${r.kind}`).join(' · ')}</div>
      )}
      {card.stats.length > 0 && (
        <div className="card-stats">
          {card.stats.map((s) => (
            <span key={s.name}>
              {s.name} {s.value}
            </span>
          ))}
        </div>
      )}
    </div>
  );
}

export function CardText({ text }: { text: string }) {
  return (
    <>
      {textParts(text).map((p, i) =>
        'icon' in p ? (
          // The parts of one text never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <span key={i} className="card-icon" title={p.icon}>
            {iconText(p.icon)}
          </span>
        ) : (
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <span key={i}>{p.text}</span>
        ),
      )}
    </>
  );
}
