import { memo, useEffect, useLayoutEffect, useMemo, useRef, useState, type CSSProperties, type PointerEvent, type RefObject } from 'react';
import { feel } from '../feel';
import { heldTiles, nearestTile, stripEntries, tileSize, type Entry } from '../history';
import { t } from '../i18n';
import type { GameView } from '../protocol';
import { playSound } from '../sounds';
import { stackBadge, stackCard, stackIcon, stackLine, type Names } from '../stack';
import { useApp, viewScale, type Zoom } from '../store';
import { playerColor } from '../table';
import { boxOf } from '../zoom';
import { Card } from './Card';
import { still } from './lunge';
import { swordPath } from './marks';
import { ItemIcon, itemZoom, SkullIcon } from './StackPanel';
import './InfoPanels.css';

interface Props {
  view: GameView;
  names: Names;
  // open: the full History panel is shown; toggle shows or hides it (V-03).
  open: boolean;
  toggle: () => void;
}

// HistoryStrip is a thin column of the history's headlines at the
// board's left edge (15.22, 15.89, H-07): resolved stack items with the
// stack's icons, buys and kills, newest on top (S-10, F-04). A badge
// marks each kind (H-08): a sword, a bolt, a coin or a skull. Rolls and
// damage are half height; a thin frame has the owner's color (20.21).
// A line marks each new turn. A new tile slides in from the left while
// the others move down (20.21, ADR 027 "History"). The strip is one
// hit area: the tile nearest the pointer shows its card large (S-11,
// H-09) with a soft sound, and the board dims; nothing moves while the
// pointer is on it. Its foot fades when more tiles are below. A
// rollback greys the undone tiles (ADR 020). The button on top opens
// the full History panel.
export function HistoryStrip({ view, names, open, toggle }: Props) {
  const history = useApp((s) => s.history);
  const setZoom = useApp((s) => s.setZoom);
  const scale = useApp(viewScale);
  const off = useApp((s) => s.settings.animation === 'off');
  const live = useMemo(() => stripEntries(history), [history]);
  // held: the tiles as the pointer found them; null while it is away.
  const [held, setHeld] = useState<Entry[] | null>(null);
  const tiles = useMemo(() => heldTiles(held, live), [held, live]);
  const [hover, setHover] = useState<number | null>(null);
  const boxRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLOListElement>(null);
  // latest is what a hover needs, as last drawn; y is the pointer's.
  const latestRef = useRef({ view, names, tiles, hover, y: null as number | null });
  useLayoutEffect(() => {
    latestRef.current = { ...latestRef.current, view, names, tiles, hover };
  });
  const more = useMore(boxRef, listRef);
  useSlide(listRef, tiles, scale);

  // show zooms the tile nearest y, once per tile: a newly hovered tile
  // sounds (ADR 027 "History" 11.9).
  const show = (y: number, touch: boolean) => {
    const list = listRef.current;
    const now = latestRef.current;
    const boxes = list ? [...list.children] : [];
    const i = nearestTile(
      boxes.map((b) => b.getBoundingClientRect()),
      y,
    );
    const el = i < 0 ? undefined : boxes.at(i);
    const entry = i < 0 ? undefined : now.tiles.at(i);
    if (!el || !entry || (!touch && entry.key === now.hover)) {
      return;
    }
    const zoom = tileZoom(entry, el, now.view, now.names, touch);
    if (!touch) {
      latestRef.current = { ...now, hover: entry.key };
      setHover(entry.key);
      playSound('historyHover');
    }
    setZoom(zoom);
  };
  const leave = () => {
    latestRef.current = { ...latestRef.current, hover: null, y: null };
    setHeld(null);
    setHover(null);
    setZoom(null);
  };
  const dim = {
    '--dim': String(feel.history.dim),
    '--dim-in': `${String(off ? 0 : feel.history.dimIn)}s`,
    '--dim-in-curve': feel.history.dimInCurve,
    '--dim-out': `${String(off ? 0 : feel.history.dimOut)}s`,
    '--dim-out-curve': feel.history.dimOutCurve,
  } as CSSProperties;
  return (
    <>
      <div className={`strip-dim${held ? ' on' : ''}`} style={dim} aria-hidden="true" data-testid="strip-dim" />
      <aside className="history-strip" aria-label={t('history.strip')} style={{ '--fade': String(feel.history.fade) } as CSSProperties}>
        <button
          type="button"
          className="history-toggle"
          aria-label={t('history.title')}
          aria-pressed={open}
          title={open ? t('history.hide_tip') : t('history.show_tip')}
          onClick={toggle}
        >
          <ScrollIcon />
        </button>
        <div
          ref={boxRef}
          className={`strip-tiles${more ? ' strip-more' : ''}`}
          data-testid="strip-tiles"
          onPointerEnter={(e: PointerEvent) => {
            if (e.pointerType === 'mouse') {
              setHeld(latestRef.current.tiles);
              latestRef.current.y = e.clientY;
              show(e.clientY, false);
            }
          }}
          onPointerMove={(e: PointerEvent) => {
            if (e.pointerType === 'mouse') {
              latestRef.current.y = e.clientY;
              show(e.clientY, false);
            }
          }}
          onPointerLeave={(e: PointerEvent) => {
            if (e.pointerType === 'mouse') {
              leave();
            }
          }}
          onPointerUp={(e: PointerEvent) => {
            if (e.pointerType !== 'mouse') {
              show(e.clientY, true);
            }
          }}
          onScroll={() => {
            // The tiles moved under a resting pointer.
            const y = latestRef.current.y;
            if (y !== null) {
              show(y, false);
            }
          }}
        >
          <ol ref={listRef} className="strip-list">
            {tiles.map((e, i) => (
              <StripItem key={e.key} entry={e} view={view} names={names} newTurn={(tiles.at(i + 1)?.turn ?? e.turn) !== e.turn} hovered={hover === e.key} />
            ))}
          </ol>
        </div>
      </aside>
    </>
  );
}

// tileZoom is a tile's popout (15.90): the play, the dealer, the killer,
// the buyer or the attacker large, then what it did; swords before an
// attack's (20.21). It opens beside the tile el.
function tileZoom(entry: Entry, el: Element, view: GameView, names: Names, touch: boolean): Zoom {
  const { plain, tile, line } = tileParts(entry, view, names);
  const outcome = entry.outcome && entry.outcome.length > 0 ? entry.outcome : undefined;
  const ref = (outcome && entry.from) ?? tile ?? '';
  return {
    ...(plain ? itemZoom(plain, ref, view, names, el, touch) : { ref, touch, at: boxOf(el), caption: { line: [{ text: line }] } }),
    ...(outcome ? { outcome } : {}),
    ...(entry.attack ? { swords: true } : {}),
  };
}

// tileParts are what a tile shows: a buy or a kill its card with the
// line; a stack item as on the stack.
function tileParts(entry: Entry, view: GameView, names: Names) {
  const it = entry.item;
  const plain = it && !entry.badge ? it : undefined;
  const tile = plain ? (stackCard(plain, view) ?? entry.card) : entry.card;
  const line = plain ? stackLine(plain, view, names) : entry.text;
  return { plain, tile, line };
}

// StripItem is one tile. newTurn: it is the first of its turn, so a line
// goes below it. It draws again only when its entry, its turn line or
// its hover changes: a new view or new names alone never redraw the
// strip (19.138); a resolved item's words do not change.
const StripItem = memo(
  function StripItem({ entry, view, names, newTurn, hovered }: { entry: Entry; view: GameView; names: Names; newTurn: boolean; hovered: boolean }) {
    const { plain, tile, line } = tileParts(entry, view, names);
    const badge = entry.badge ?? (plain ? stackBadge(plain) : undefined);
    const color = entry.owner === undefined ? undefined : playerColor(entry.owner);
    const size = tileSize(entry);
    return (
      <li
        className={`strip-item stack-${plain ? stackIcon(plain) : 'card'} strip-${size}${newTurn ? ' strip-turn' : ''}${entry.undone ? ' undone' : ''}${hovered ? ' hovered' : ''}${color ? ' owned' : ''}`}
        style={color ? ({ '--owner': color } as CSSProperties) : undefined}
        data-key={entry.key}
        aria-label={entry.undone ? t('history.undone', { line }) : line}
        title={entry.undone ? t('history.undone', { line: t('history.row', { n: entry.turn, line }) }) : t('history.row', { n: entry.turn, line })}
      >
        {plain ? (
          <ItemIcon it={plain} />
        ) : (
          <span className="stack-icon" aria-hidden="true">
            {tile && <Card cardRef={tile} size="small" noZoom />}
          </span>
        )}
        {badge && (
          <span className={`strip-badge strip-${badge}`} aria-hidden="true">
            {badge === 'attack' ? <SwordIcon /> : badge === 'spell' ? <BoltIcon /> : badge === 'buy' ? <CoinIcon /> : <SkullIcon />}
          </span>
        )}
      </li>
    );
  },
  (a, b) => a.entry === b.entry && a.newTurn === b.newTurn && a.hovered === b.hovered,
);

// useMore tells whether more tiles are below the strip's foot (20.21,
// ADR 027 "History" 11.5): the foot then fades. It reads the sizes only
// when they change or the strip scrolls, never on each update.
function useMore(boxRef: RefObject<HTMLDivElement | null>, listRef: RefObject<HTMLOListElement | null>): boolean {
  const [more, setMore] = useState(false);
  useEffect(() => {
    const b = boxRef.current;
    const l = listRef.current;
    if (!b || !l) {
      return;
    }
    const check = () => {
      setMore(b.scrollHeight - b.scrollTop - b.clientHeight > 1);
    };
    check();
    b.addEventListener('scroll', check, { passive: true });
    const ro = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(check);
    ro?.observe(b);
    ro?.observe(l);
    return () => {
      b.removeEventListener('scroll', check);
      ro?.disconnect();
    };
  }, [boxRef, listRef]);
  return more;
}

// useSlide slides new tiles in (20.21, ADR 027 "History" 11.6): each new
// tile comes from the left, level with its place, while the list moves
// down by their height, all in feel.history.slide with its curve. The
// list moves as one: the old tiles are never measured or moved one by
// one. A tile that comes while others still move takes over from where
// they are. Speed "off", a rollback's table and reduced motion place the
// tiles at once.
function useSlide(listRef: RefObject<HTMLOListElement | null>, tiles: Entry[], scale: number): void {
  const beforeRef = useRef<Set<number> | null>(null);
  const movingRef = useRef(new Map<Element, Animation>());
  useLayoutEffect(() => {
    const l = listRef.current;
    const keys = tiles.map((e) => e.key);
    const prev = beforeRef.current;
    beforeRef.current = new Set(keys);
    if (!l || !prev || prev.size === 0) {
      return;
    }
    // k: how many tiles are new on top.
    const k = keys.findIndex((key) => prev.has(key));
    if (k <= 0) {
      return;
    }
    const running = movingRef.current;
    if (scale <= 0 || still()) {
      for (const a of running.values()) {
        a.cancel();
      }
      running.clear();
      return;
    }
    const items = [...l.children] as HTMLElement[];
    const first = items.at(0);
    const old = items.at(k);
    if (!first || !old) {
      return;
    }
    const h = old.offsetTop - first.offsetTop;
    const x = first.offsetWidth * feel.history.from;
    // Where each moving part is now, before it is taken over.
    const at = new Map<Element, DOMMatrixReadOnly>();
    for (const [el, a] of running) {
      if (el.isConnected) {
        at.set(el, new DOMMatrixReadOnly(getComputedStyle(el).transform));
      }
      a.cancel();
    }
    running.clear();
    const r = at.get(l)?.m42 ?? 0;
    const ms = feel.history.slide * 1000 * scale;
    const go = (el: Element, fromX: number, fromY: number) => {
      // anim: strip-slide
      const a = el.animate([{ transform: `translate(${String(fromX)}px, ${String(fromY)}px)` }, { transform: 'none' }], {
        duration: ms,
        easing: feel.history.slideCurve,
      });
      running.set(el, a);
      a.onfinish = () => {
        running.delete(el);
      };
    };
    go(l, 0, r - h);
    // The new tiles stay level with their places as the list moves.
    items.slice(0, k).forEach((el) => {
      go(el, -x, h - r);
    });
    // Tiles that were still coming in go on from where they are.
    for (const [el, m] of at) {
      if (el !== l && items.indexOf(el as HTMLElement) >= k) {
        go(el, m.m41, m.m42);
      }
    }
  }, [listRef, tiles, scale]);
}

function ScrollIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M6 3h11a3 3 0 0 1 3 3v1h-3v11a3 3 0 0 1-3 3H5a3 3 0 0 1-3-3v-1h3V6a3 3 0 0 1 1-3Zm2 4v2h7V7Zm0 4v2h7v-2Zm0 4v2h5v-2Z" />
    </svg>
  );
}

function SwordIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d={swordPath} />
    </svg>
  );
}

function CoinIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M12 2a10 10 0 1 1 0 20 10 10 0 0 1 0-20Zm1 4h-2v1.1A3.5 3.5 0 0 0 8 10.5c0 2 1.6 2.8 3.4 3.3 1.5.4 2.1.7 2.1 1.4 0 .6-.6 1-1.5 1-1 0-1.6-.4-1.8-1.2H8.1c.2 1.7 1.3 2.8 2.9 3.1V19h2v-1c1.8-.3 3-1.6 3-3.3 0-2.2-1.8-2.9-3.5-3.4-1.4-.4-2-.6-2-1.2 0-.5.5-.9 1.3-.9.9 0 1.4.4 1.5 1h2.1c-.1-1.5-1.1-2.6-2.4-3Z" />
    </svg>
  );
}

function BoltIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M13.5 2 4 14h6.5L9.5 22 20 9.5h-6.6Z" />
    </svg>
  );
}
