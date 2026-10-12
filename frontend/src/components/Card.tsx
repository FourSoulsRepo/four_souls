import { motion } from 'motion/react';
import { useEffect, useRef, useState, type PointerEvent } from 'react';
import { iconText, kindOf, textParts, unknownCard, type CardInfo, type Kind } from '../cards';
import { scaleFor } from '../playback';
import { useApp } from '../store';
import './Card.css';

// How long a touch must rest on a card to zoom it (B-06).
const LONG_PRESS_MS = 450;

export type CardSize = 'hand' | 'table' | 'small' | 'zoom';

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
  // and a new one pops in.
  objectId?: number;
  // counters is how many counters the card carries (7.9).
  counters?: number;
  // highlight marks a card the player may act on or pick (7.7).
  highlight?: 'act' | 'pick';
  // stats shows live health and attack as big numbers (H-15).
  stats?: { hp: number; atk: number };
  // noZoom: the card does not zoom itself; its row does (15.21).
  noZoom?: boolean;
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
  highlight,
  stats,
  noZoom = false,
  onClick,
}: Props) {
  const info = useApp((s) => (cardRef ? s.cards.get(cardRef) : undefined));
  const backImage = useApp((s) => s.backs.get(back) ?? '');
  const setZoom = useApp((s) => s.setZoom);
  const scale = useApp((s) => scaleFor(s.settings.animation, 0, 0));
  const pressRef = useRef(0);
  const downRef = useRef<{ x: number; y: number } | null>(null);
  const [brokenImage, setBrokenImage] = useState('');

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
  const zoomable = size !== 'zoom' && !noZoom;
  const side = (e: PointerEvent) => (e.clientX < window.innerWidth / 2 ? 'right' : 'left');
  const cancelPress = () => {
    window.clearTimeout(pressRef.current);
    downRef.current = null;
  };

  const pop = objectId !== undefined && scale > 0;
  return (
    <motion.div
      initial={pop ? { scale: 0.6, opacity: 0 } : false}
      animate={{ scale: 1, opacity: 1 }}
      transition={{ duration: 0.25 * scale }}
      data-zone={objectId === undefined ? undefined : `obj-${String(objectId)}`}
      className={`card size-${size} kind-${kindOf(card.type)}${spent ? ' spent' : ''}${stats ? ' has-stats' : ''}${highlight ? ` hl-${highlight}` : ''}`}
      role={onClick ? 'button' : 'img'}
      aria-label={card.name}
      onClick={onClick}
      onPointerEnter={(e) => {
        if (zoomable && e.pointerType === 'mouse') {
          zoomedRef.current = true;
          setZoom({ ref: card.ref, touch: false, side: side(e) });
        }
      }}
      onPointerLeave={(e) => {
        cancelPress();
        if (zoomable && e.pointerType === 'mouse') {
          zoomedRef.current = false;
          setZoom(null);
        }
      }}
      onPointerDown={(e) => {
        if (zoomable && e.pointerType !== 'mouse') {
          const s = side(e);
          downRef.current = { x: e.clientX, y: e.clientY };
          pressRef.current = window.setTimeout(() => {
            setZoom({ ref: card.ref, touch: true, side: s });
          }, LONG_PRESS_MS);
        }
      }}
      onPointerMove={(e) => {
        // A finger that moves is dragging, not pressing.
        const d = downRef.current;
        if (d && Math.hypot(e.clientX - d.x, e.clientY - d.y) > 8) {
          cancelPress();
        }
      }}
      onPointerUp={cancelPress}
      onPointerCancel={cancelPress}
      onContextMenu={(e) => {
        e.preventDefault(); // a long press must not open the menu
      }}
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
      {counters > 0 && (
        <span className="card-counters" aria-label={`${String(counters)} counters`}>
          {counters}
        </span>
      )}
      {damage > 0 && (
        <span className="card-damage" aria-label={`${String(damage)} damage`}>
          {damage}
        </span>
      )}
      {stats && (
        <>
          <span className="card-stat card-hp" title="Health" aria-label={`${String(stats.hp)} health`}>
            {stats.hp}
          </span>
          <span className="card-stat card-atk" title="Attack" aria-label={`${String(stats.atk)} attack`}>
            {stats.atk}
          </span>
        </>
      )}
    </motion.div>
  );
}

// BlankFace shows a card's text; full adds the effects and keywords.
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
          {card.footnotes.map((f) => (
            <p key={f.title} className="card-footnote">
              <b>{f.title}</b> <CardText text={f.text} />
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

function CardText({ text }: { text: string }) {
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
