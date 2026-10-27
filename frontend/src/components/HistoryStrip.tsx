import type { PointerEvent } from 'react';
import { stripEntries, type Entry } from '../history';
import { t } from '../i18n';
import type { GameView } from '../protocol';
import { stackBadge, stackCard, stackIcon, stackLine, type Names } from '../stack';
import { useApp, type Zoom } from '../store';
import { boxOf } from '../zoom';
import { Card } from './Card';
import { Row, Rows } from './Rows';
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
// marks each kind (H-08): a sword, a bolt, a coin or a skull. Hover shows
// the large card (S-11, H-09). A line marks each new turn. A new tile
// fades in and the others slide down (15.51). A rollback greys the
// undone tiles (ADR 020). The button on top opens the full History panel.
export function HistoryStrip({ view, names, open, toggle }: Props) {
  const history = useApp((s) => s.history);
  const items = stripEntries(history);
  return (
    <aside className="history-strip" aria-label={t('history.strip')}>
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
      <ol>
        <Rows>
          {items.map((e, i) => (
            <StripItem key={e.key} entry={e} view={view} names={names} newTurn={(items.at(i + 1)?.turn ?? e.turn) !== e.turn} />
          ))}
        </Rows>
      </ol>
    </aside>
  );
}

// StripItem is one resolved item. newTurn: it is the first of its turn,
// so a line goes below it.
function StripItem({ entry, view, names, newTurn }: { entry: Entry; view: GameView; names: Names; newTurn: boolean }) {
  const setZoom = useApp((s) => s.setZoom);
  const it = entry.item;
  // A buy or a kill shows its card with the line; a stack item as on the stack.
  const plain = it && !entry.badge ? it : undefined;
  const tile = plain ? (stackCard(plain, view) ?? entry.card) : entry.card;
  const badge = entry.badge ?? (plain ? stackBadge(plain) : undefined);
  const line = plain ? stackLine(plain, view, names) : entry.text;
  // The popout (15.90): the play, the dealer, the killer or the buyer
  // large, then what it did.
  const outcome = entry.outcome && entry.outcome.length > 0 ? entry.outcome : undefined;
  const card = (outcome && entry.from) ?? tile;
  const zoom = (ref: string, e: PointerEvent, touch: boolean): Zoom => ({
    ...(plain ? itemZoom(plain, ref, view, names, e, touch) : { ref, touch, at: boxOf(e.currentTarget), caption: { line: [{ text: line }] } }),
    ...(outcome ? { outcome } : {}),
  });
  return (
    <Row
      className={`strip-item stack-${plain ? stackIcon(plain) : 'card'}${newTurn ? ' strip-turn' : ''}${entry.undone ? ' undone' : ''}`}
      aria-label={entry.undone ? t('history.undone', { line }) : line}
      title={entry.undone ? t('history.undone', { line: t('history.row', { n: entry.turn, line }) }) : t('history.row', { n: entry.turn, line })}
      onPointerEnter={(e) => {
        if (card && e.pointerType === 'mouse') {
          setZoom(zoom(card, e, false));
        }
      }}
      onPointerLeave={(e) => {
        if (card && e.pointerType === 'mouse') {
          setZoom(null);
        }
      }}
      onPointerUp={(e) => {
        if (card && e.pointerType !== 'mouse') {
          setZoom(zoom(card, e, true));
        }
      }}
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
    </Row>
  );
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
      <path d="M20.5 2.5 21 7l-9.6 9.6 1.4 1.4-1.4 1.4-2.1-2.1-3.2 3.2a1.5 1.5 0 0 1-2.1-2.1l3.2-3.2-2.1-2.1 1.4-1.4 1.4 1.4L17.5 3Z" />
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
