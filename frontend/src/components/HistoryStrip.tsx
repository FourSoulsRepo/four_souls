import type { Entry } from '../history';
import type { GameView } from '../protocol';
import { stackBadge, stackCard, stackIcon, stackLine, type Names } from '../stack';
import { useApp } from '../store';
import { ItemIcon, itemZoom } from './StackPanel';
import './InfoPanels.css';

interface Props {
  view: GameView;
  names: Names;
  // open: the full History panel is shown; toggle shows or hides it (V-03).
  open: boolean;
  toggle: () => void;
}

// HistoryStrip is a thin column of resolved stack items at the board's
// left edge (15.22, H-07): the stack's icons, newest on top (S-10, F-04),
// a sword or bolt badge (H-08), and the same large card on hover (S-11,
// H-09). A line marks each new turn. The button on top opens the full
// History panel.
export function HistoryStrip({ view, names, open, toggle }: Props) {
  const history = useApp((s) => s.history);
  const items = history.filter((e) => e.item).reverse();
  return (
    <aside className="history-strip" aria-label="History strip">
      <button
        type="button"
        className="history-toggle"
        aria-label="History"
        aria-pressed={open}
        title={open ? 'Hide the history' : 'Show the history'}
        onClick={toggle}
      >
        <ScrollIcon />
      </button>
      <ol>
        {items.map((e, i) => (
          <StripItem key={e.key} entry={e} view={view} names={names} newTurn={(items.at(i + 1)?.turn ?? e.turn) !== e.turn} />
        ))}
      </ol>
    </aside>
  );
}

// StripItem is one resolved item. newTurn: it is the first of its turn,
// so a line goes below it.
function StripItem({ entry, view, names, newTurn }: { entry: Entry; view: GameView; names: Names; newTurn: boolean }) {
  const setZoom = useApp((s) => s.setZoom);
  const it = entry.item;
  if (!it) {
    return null;
  }
  const card = stackCard(it, view) ?? entry.card;
  const badge = stackBadge(it);
  const line = stackLine(it, view, names);
  return (
    <li
      className={`strip-item stack-${stackIcon(it)}${newTurn ? ' strip-turn' : ''}`}
      aria-label={line}
      title={`Turn ${String(entry.turn)}: ${line}`}
      onPointerEnter={(e) => {
        if (card && e.pointerType === 'mouse') {
          setZoom(itemZoom(it, card, view, names, e, false));
        }
      }}
      onPointerLeave={(e) => {
        if (card && e.pointerType === 'mouse') {
          setZoom(null);
        }
      }}
      onPointerUp={(e) => {
        if (card && e.pointerType !== 'mouse') {
          setZoom(itemZoom(it, card, view, names, e, true));
        }
      }}
    >
      <ItemIcon it={it} />
      {badge && (
        <span className={`strip-badge strip-${badge}`} aria-hidden="true">
          {badge === 'attack' ? <SwordIcon /> : <BoltIcon />}
        </span>
      )}
    </li>
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

function BoltIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M13.5 2 4 14h6.5L9.5 22 20 9.5h-6.6Z" />
    </svg>
  );
}
