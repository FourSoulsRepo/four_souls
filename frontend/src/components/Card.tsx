import { useRef, useState, type PointerEvent } from 'react';
import { iconText, kindOf, textParts, unknownCard, type CardInfo, type Kind } from '../cards';
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
}

// Card is the one card component (7.4): the card image, or a blank card
// with its text when the image is missing (A-11).
export function Card({ cardRef, size, back = 'loot', damage = 0, spent = false }: Props) {
  const info = useApp((s) => (cardRef ? s.cards.get(cardRef) : undefined));
  const setZoom = useApp((s) => s.setZoom);
  const pressRef = useRef(0);
  const [brokenImage, setBrokenImage] = useState('');

  if (cardRef === undefined) {
    return (
      <div className={`card size-${size} kind-${back}`} role="img" aria-label={`Face-down ${back} card`}>
        <div className="card-back-face">{back}</div>
      </div>
    );
  }

  const card = info ?? unknownCard(cardRef);
  const image = card.image !== '' && card.image !== brokenImage;
  const zoomable = size !== 'zoom';
  const side = (e: PointerEvent) => (e.clientX < window.innerWidth / 2 ? 'right' : 'left');
  const cancelPress = () => {
    window.clearTimeout(pressRef.current);
  };

  return (
    <div
      className={`card size-${size} kind-${kindOf(card.type)}${spent ? ' spent' : ''}`}
      role="img"
      aria-label={card.name}
      onPointerEnter={(e) => {
        if (zoomable && e.pointerType === 'mouse') {
          setZoom({ ref: card.ref, touch: false, side: side(e) });
        }
      }}
      onPointerLeave={(e) => {
        cancelPress();
        if (zoomable && e.pointerType === 'mouse') {
          setZoom(null);
        }
      }}
      onPointerDown={(e) => {
        if (zoomable && e.pointerType !== 'mouse') {
          const s = side(e);
          pressRef.current = window.setTimeout(() => {
            setZoom({ ref: card.ref, touch: true, side: s });
          }, LONG_PRESS_MS);
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
      {damage > 0 && (
        <span className="card-damage" aria-label={`${String(damage)} damage`}>
          {damage}
        </span>
      )}
    </div>
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
