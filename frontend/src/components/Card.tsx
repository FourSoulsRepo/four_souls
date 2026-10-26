import { AnimatePresence, motion, useIsPresent } from 'motion/react';
import { useEffect, useRef, useState, type MouseEvent, type PointerEvent } from 'react';
import { claimed } from '../animations';
import { iconText, kindOf, textParts, unknownCard, type CardInfo, type Kind } from '../cards';
import { diceImage } from '../dice';
import { t } from '../i18n';
import { scaleFor } from '../playback';
import { fighters } from '../prompt';
import type { MonsterStats, StatSource } from '../protocol';
import { statLines, sumOf, tempSources, trend, type Name, type TempSource, type Trend } from '../stats';
import { useApp, viewScale } from '../store';
import { boxOf } from '../zoom';
import { Die } from './Die';
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
  // dice are the die numbers the card keeps, drawn in owner's dice look
  // (15.67).
  dice?: number[];
  owner?: number;
  // highlight marks a card the player may act on or pick (7.7).
  highlight?: 'act' | 'pick' | 'aim';
  // stats shows live health and attack as big numbers (H-15); the print
  // and the sources say why they changed (15.92, V-41).
  stats?: { hp: number; atk: number; printed_hp?: number; printed_atk?: number; sources?: StatSource[] };
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
  dice,
  owner = -1,
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
  const scale = useApp(viewScale);
  const [brokenImage, setBrokenImage] = useState('');
  const zoomable = size !== 'zoom' && !noZoom;
  const cards = useApp((s) => s.cards);
  const name: Name = (ref) => (ref ? (cards.get(ref)?.name ?? ref) : t('stats.hidden_card'));
  const sources = stats?.sources ?? monster?.sources;
  const lines = sources ? statLines(sources, printedOf(stats, monster), liveOf(stats, monster), name) : [];
  const zoomHandlers = useZoom(cardRef, zoomable, rise, why, lines);
  const cardElRef = useRef<HTMLDivElement>(null);
  // Leaving an icon for the card zooms the card again.
  const backToCard = (e: PointerEvent<HTMLElement>) => {
    const el = cardElRef.current;
    if (el && e.relatedTarget instanceof Node && el.contains(e.relatedTarget)) {
      zoomHandlers.onPointerEnter({ ...e, currentTarget: el });
    }
  };
  const showWhy = useApp((s) => s.showWhy);
  // A card in a fight has a red frame for everyone (15.31, F-19).
  const fight = useApp((s) => objectId !== undefined && !!s.game && fighters(s.game.view).includes(objectId));
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
        aria-label={faceDown(back)}
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
          <div className="card-back-face">{backWord(back)}</div>
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
      ref={cardElRef}
      initial={pop ? { scale: 0.6, opacity: 0 } : false}
      animate={{ scale: 1, opacity: 1 }}
      transition={{ duration: 0.25 * scale }}
      data-zone={objectId === undefined ? undefined : `obj-${String(objectId)}`}
      data-player={player}
      className={`card size-${size} kind-${kindOf(card.type)}${image ? '' : ' blank'}${counters > 0 ? ' has-counters' : ''}${spent ? ' spent' : ''}${fight ? ' fight' : ''}${highlight ? ` hl-${highlight}` : ''}`}
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
      {badges && counters > 0 && <Badge kind="counters" value={counters} label={t('badge.counters', { n: counters })} />}
      {badges && dice && dice.length > 0 && <DieBadges values={dice} owner={owner} />}
      {/* Damage shows as the lower, red health; the mark is for the rest. */}
      {badges && damage > 0 && !monster && !stats && <Badge kind="damage" value={damage} label={t('badge.damage', { n: damage })} />}
      {badges && stats && (
        <>
          <Badge kind="hp" value={stats.hp} label={t('badge.hp', { n: stats.hp })} title={t('badge.hp_title')} hurt={damage > 0} trend={trend(sumOf(stats.sources, 'hp'), 0)} />
          <Badge kind="atk" value={stats.atk} label={t('badge.atk', { n: stats.atk })} title={t('badge.atk_title')} trend={stats.printed_atk === undefined ? undefined : trend(stats.atk, stats.printed_atk)} />
        </>
      )}
      {badges && monster && <MonsterBadges stats={monster} hurt={damage > 0} />}
      {badges && <SourceIcons list={tempSources(sources)} name={name} onLeave={backToCard} />}
    </motion.div>
  );
}

// useZoom zooms a card from the element it is spread on: on mouse
// hover, beside it or rising from it, or after a long touch (B-06). A
// card and a soul crystal share it (V-43). Nothing zooms without a
// card ref or when zoomable is false.
export function useZoom(cardRef: string | undefined, zoomable: boolean, rise = false, why?: string, stats: string[] = []) {
  const setZoom = useApp((s) => s.setZoom);
  const pressRef = useRef(0);
  const downRef = useRef<{ x: number; y: number } | null>(null);
  const on = zoomable && cardRef !== undefined;

  // A card that leaves the screen while hovered gets no pointer-leave:
  // the card that started a mouse zoom ends it then.
  const zoomedRef = useRef(false);
  // Nor does a card whose hovered part was redrawn, such as an image
  // that failed and became a drawn face (19.50): the zoom also ends on
  // the first pointer-over outside the card. unwatch stops that watch.
  const unwatchRef = useRef<(() => void) | null>(null);
  const unwatch = () => {
    unwatchRef.current?.();
    unwatchRef.current = null;
  };
  useEffect(
    () => () => {
      unwatch();
      const { zoom, setZoom: clear } = useApp.getState();
      if (zoomedRef.current && zoom && zoom.ref === cardRef && !zoom.touch) {
        clear(null);
      }
    },
    [cardRef],
  );
  const watch = (el: HTMLElement) => {
    unwatch();
    const over = (e: Event) => {
      if (e.target instanceof Node && el.contains(e.target)) {
        return;
      }
      unwatch();
      zoomedRef.current = false;
      const { zoom, setZoom: clear } = useApp.getState();
      if (zoom && zoom.ref === cardRef && !zoom.touch) {
        clear(null);
      }
    };
    document.addEventListener('pointerover', over, true);
    unwatchRef.current = () => {
      document.removeEventListener('pointerover', over, true);
    };
  };

  const cancelPress = () => {
    window.clearTimeout(pressRef.current);
    downRef.current = null;
  };
  return {
    onPointerEnter: (e: PointerEvent<HTMLElement>) => {
      if (on && e.pointerType === 'mouse') {
        zoomedRef.current = true;
        watch(e.currentTarget);
        setZoom({ ref: cardRef, touch: false, at: boxOf(e.currentTarget), rise, ...(why ? { why } : {}), ...(stats.length > 0 ? { stats } : {}) });
      }
    },
    onPointerLeave: (e: PointerEvent<HTMLElement>) => {
      cancelPress();
      if (on && e.pointerType === 'mouse') {
        unwatch();
        zoomedRef.current = false;
        setZoom(null);
      }
    },
    onPointerDown: (e: PointerEvent<HTMLElement>) => {
      if (on && e.pointerType !== 'mouse') {
        const at = boxOf(e.currentTarget);
        downRef.current = { x: e.clientX, y: e.clientY };
        pressRef.current = window.setTimeout(() => {
          setZoom({ ref: cardRef, touch: true, at, ...(why ? { why } : {}), ...(stats.length > 0 ? { stats } : {}) });
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
  return (
    <>
      {(stats.hp !== stats.printed_hp || hurt) && (
        <Badge kind="hp" value={stats.hp} label={t('badge.hp', { n: stats.hp })} title={t('badge.hp_printed', { n: stats.printed_hp })} hurt={hurt} trend={trend(sumOf(stats.sources, 'hp'), 0)} />
      )}
      {stats.dc !== stats.printed_dc && (
        <Badge kind="dc" value={stats.dc} label={t('badge.dc', { n: stats.dc })} title={t('badge.dc_printed', { n: stats.printed_dc })} trend={trend(stats.dc, stats.printed_dc)} />
      )}
      {stats.atk !== stats.printed_atk && (
        <Badge kind="atk" value={stats.atk} label={t('badge.atk', { n: stats.atk })} title={t('badge.atk_printed', { n: stats.printed_atk })} trend={trend(stats.atk, stats.printed_atk)} />
      )}
    </>
  );
}

// DieBadges are the die numbers a card keeps (15.67): a stored roll or
// a chosen number, each a die in its owner's dice look, on the art.
function DieBadges({ values, owner }: { values: number[]; owner: number }) {
  const image = useApp((s) => diceImage(s.dice, s.game?.seats, owner));
  return (
    <span className="card-dice">
      {values.map((v) => (
        <span key={v} className="card-die" title={t('badge.kept_die_title', { n: v })}>
          <Die value={v} image={image} label={t('badge.kept_die', { n: v })} />
        </span>
      ))}
    </span>
  );
}

// printedOf and liveOf are a character's or monster's stats for its sums.
function printedOf(stats: Props['stats'], monster: MonsterStats | undefined) {
  if (monster) return { hp: monster.printed_hp, atk: monster.printed_atk, dc: monster.printed_dc };
  return { hp: stats?.printed_hp ?? stats?.hp ?? 0, atk: stats?.printed_atk ?? stats?.atk ?? 0 };
}

function liveOf(stats: Props['stats'], monster: MonsterStats | undefined) {
  if (monster) return { atk: monster.atk, dc: monster.dc };
  return { atk: stats?.atk ?? 0 };
}

// SourceIcons are the effects of this turn on a character or monster
// (15.92, SB-01): a column of small icons cut from each source card's
// art; lasting sources show only in the zoom's sums. A new icon fades
// in, a leaving one fades out (15.93); with animations off they come
// and go at once.
function SourceIcons({ list, name, onLeave }: { list: TempSource[]; name: Name; onLeave: (e: PointerEvent<HTMLElement>) => void }) {
  const scale = useApp(viewScale);
  // A rollback's table keeps the icons mounted; they only jump.
  const on = useApp((s) => scaleFor(s.settings.animation, 0, 0) > 0);
  const icons = list.map((t) => <SourceIcon key={t.key} source={t} name={name(t.card)} onLeave={onLeave} scale={scale} />);
  if (!on) {
    return list.length === 0 ? null : <span className="card-sources">{icons}</span>;
  }
  return (
    <span className="card-sources">
      <AnimatePresence initial={false}>{icons}</AnimatePresence>
    </span>
  );
}

// SourceIcon is one effect: hover zooms its card with what it does
// (SB-03). A text card shows its initial; a hidden one a question mark.
// A boost that waits for the next attack roll has a small die; used by
// the roll, it pops and fades; at the end of the turn an icon shrinks
// away (15.93, SB-12, SB-13). A leaving icon takes no pointer and
// screen readers skip it.
function SourceIcon({ source, name, onLeave, scale }: { source: TempSource; name: string; onLeave: (e: PointerEvent<HTMLElement>) => void; scale: number }) {
  const image = useApp((s) => (source.card ? (s.cards.get(source.card)?.image ?? '') : ''));
  const [broken, setBroken] = useState(false);
  const present = useIsPresent();
  const zoom = useZoom(source.card, source.card !== undefined && present, false, undefined, source.lines);
  return (
    // anim: source-icon
    <motion.span
      className={`card-source${source.nextRoll ? ' next-roll' : ''}`}
      role="img"
      aria-label={`${name}: ${source.lines.join(', ')}`}
      aria-hidden={present ? undefined : true}
      style={present ? undefined : { pointerEvents: 'none' }}
      initial={scale > 0 ? { opacity: 0, scale: 0.4 } : false}
      animate={{ opacity: 1, scale: 1 }}
      exit={source.nextRoll ? { opacity: 0, scale: 1.8 } : { opacity: 0, scale: 0.5 }}
      transition={{ duration: (source.nextRoll ? 0.35 : 0.6) * scale }}
      {...zoom}
      onPointerLeave={(e) => {
        zoom.onPointerLeave(e);
        onLeave(e);
      }}
    >
      {image !== '' && !broken ? (
        <img
          src={image}
          alt=""
          draggable={false}
          onError={() => {
            setBroken(true);
          }}
        />
      ) : (
        <span className="card-source-mark">{source.card ? name.charAt(0) : '?'}</span>
      )}
      {source.nextRoll && <span className="card-source-die" aria-hidden="true" />}
    </motion.span>
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
  // label is the screen-reader name: the number and what it counts.
  label: string;
  title?: string;
  hurt?: boolean;
  // trend colors a stat that went up or down from its print (15.92).
  trend?: Trend;
}

// Badge is one number on a card: its shape in ink with a hard shadow,
// the number in the game's hand-drawn letters. A hurt health number is
// red while below full (T-25); a stat above its print is green, below
// it red (15.92).
function Badge({ kind, value, label, title, hurt = false, trend: dir }: BadgeProps) {
  return (
    <span className={`card-badge card-${kind}${hurt ? ' hurt' : ''}${dir ? ` ${dir}` : ''}`} title={title} aria-label={label}>
      <svg className="card-badge-shape" viewBox="0 0 100 100" aria-hidden="true">
        <path d={shapes.get(kind) ?? blob} />
        <path className="card-badge-shine" d="M30 26C36 18 46 16 52 18C44 22 38 28 34 36C31 34 29 30 30 26Z" />
      </svg>
      <span className="card-badge-value">{value}</span>
    </span>
  );
}

// backWord names a deck on a card back drawn without its image.
function backWord(kind: Kind): string {
  switch (kind) {
    case 'character':
      return t('deck.back.character');
    case 'treasure':
      return t('deck.back.treasure');
    case 'monster':
      return t('deck.back.monster');
    case 'soul':
      return t('deck.back.soul');
    default:
      return t('deck.back.loot');
  }
}

// faceDown is a card back's screen-reader name.
function faceDown(kind: Kind): string {
  switch (kind) {
    case 'character':
      return t('deck.face_down.character');
    case 'treasure':
      return t('deck.face_down.treasure');
    case 'monster':
      return t('deck.face_down.monster');
    case 'soul':
      return t('deck.face_down.soul');
    default:
      return t('deck.face_down.loot');
  }
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
