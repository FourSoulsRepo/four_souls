import type { PointerEvent } from 'react';
import { useKey } from '../keys';
import type { GameView, Intent, StackItem } from '../protocol';
import { stackCard, stackIcon, stackParts, stackWait, type LinePart, type Names } from '../stack';
import { useApp, type Zoom } from '../store';
import { playerColor } from '../table';
import { Card } from './Card';
import { Die } from './Die';
import './Actions.css';

interface Props {
  view: GameView;
  you: number;
  names: Names;
  // pass passes priority, when the engine allows it now.
  pass?: Intent;
  disabled: boolean;
}

// StackPanel is the stack inside the board (15.21, F-03): newest on
// top, one icon and one short line per item (S-01 – S-04), whom it
// waits for (V-11), and the Pass button; Space presses it (S-08). An
// empty stack says so (S-09).
export function StackPanel({ view, you, names, pass, disabled }: Props) {
  const act = useApp((s) => s.act);
  const stack = view.stack ?? [];
  const wait = stackWait(view, you, names);
  const lit = pass !== undefined && !disabled;
  const press = () => {
    if (lit) {
      act(pass);
    }
  };
  useKey(' ', press);
  return (
    <aside className={`stack${stack.length === 0 ? ' stack-empty' : ''}`} aria-label="Stack">
      <header>
        <h2>Stack</h2>
        <button type="button" className="stack-pass" title="Pass (Space)" disabled={!lit} onClick={press}>
          Pass
        </button>
      </header>
      {stack.length === 0 ? (
        <p className="stack-nothing">Nothing on the stack</p>
      ) : (
        <ol>
          {[...stack].reverse().map((it) => (
            <StackRow key={it.id} it={it} view={view} names={names} />
          ))}
        </ol>
      )}
      {wait.prompt && (
        <p className="stack-wait" aria-label="Waiting for">
          {wait.prompt}
        </p>
      )}
      {wait.responding.length > 0 && (
        <p className="stack-wait" aria-label="Waiting for">
          Waiting for {wait.responding.join(', ')}
        </p>
      )}
      {wait.passed.length > 0 && (
        <p className="stack-wait stack-passed" aria-label="Passed">
          Passed: {wait.passed.join(', ')}
        </p>
      )}
    </aside>
  );
}

// StackRow is one item. Hover shows its card large, the used ability
// highlighted and the owner below (S-06, S-07); a tap does the same.
function StackRow({ it, view, names }: { it: StackItem; view: GameView; names: Names }) {
  const setZoom = useApp((s) => s.setZoom);
  const card = stackCard(it, view);
  const parts = stackParts(it, view, names);
  const show = (e: PointerEvent, touch: boolean) => {
    if (card) {
      setZoom(itemZoom(it, card, names, e, touch));
    }
  };
  return (
    <li
      className={`stack-row stack-${stackIcon(it)}`}
      aria-label={parts.map((p) => p.text).join('')}
      onPointerEnter={(e) => {
        if (e.pointerType === 'mouse') {
          show(e, false);
        }
      }}
      onPointerLeave={(e) => {
        if (card && e.pointerType === 'mouse') {
          setZoom(null);
        }
      }}
      onPointerUp={(e) => {
        if (e.pointerType !== 'mouse') {
          show(e, true);
        }
      }}
    >
      <ItemIcon it={it} />
      <span className="stack-line">
        {parts.map((p, i) => (
          // A line's parts never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <Part key={i} part={p} />
        ))}
      </span>
    </li>
  );
}

// itemZoom is the large card of a stack item: the used ability's words
// and the owner in their color below it (S-06, S-07). It opens on the
// side away from the pointer.
export function itemZoom(it: StackItem, card: string, names: Names, e: PointerEvent, touch: boolean): Zoom {
  const owner = it.controller >= 0 ? names.player(it.controller) : undefined;
  return {
    ref: card,
    touch,
    side: e.clientX < window.innerWidth / 2 ? 'right' : 'left',
    caption: { owner, color: playerColor(it.controller), ability: it.text },
  };
}

// ItemIcon is a stack item's square icon (S-04): card art, a die in the
// viewer's dice look, a blood drop with the damage, or a skull.
export function ItemIcon({ it }: { it: StackItem }) {
  const diceImage = useApp((s) => s.dice.find((d) => d.id === s.settings.dice)?.image ?? '');
  const icon = stackIcon(it);
  return (
    <span className="stack-icon" aria-hidden="true">
      {icon === 'card' ? (
        it.card && <Card cardRef={it.card} size="small" noZoom />
      ) : icon === 'roll' ? (
        it.roll === undefined ? (
          <span className="stack-die">?</span>
        ) : (
          <Die value={it.roll} image={diceImage} />
        )
      ) : icon === 'damage' ? (
        <span className="stack-mark">
          <DamageIcon />
          <b>{it.amount ?? 0}</b>
        </span>
      ) : (
        <span className="stack-mark">
          <SkullIcon />
        </span>
      )}
    </span>
  );
}

function Part({ part }: { part: LinePart }) {
  if (part.role === 'target') {
    return <b className="stack-target">{part.text}</b>;
  }
  if (part.role === 'player') {
    return (
      <b className="stack-player" style={{ color: playerColor(part.player ?? -1) }}>
        {part.text}
      </b>
    );
  }
  if (part.role === 'card') {
    return <b className="stack-card">{part.text}</b>;
  }
  return <>{part.text}</>;
}

// DamageIcon is a blood drop (S-04).
function DamageIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M12 2.5C9 7 6 10.5 6 14.5a6 6 0 0 0 12 0C18 10.5 15 7 12 2.5Z" />
    </svg>
  );
}

function SkullIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M12 2.5a8 8 0 0 0-8 8c0 2.8 1.4 4.6 3 5.6V19a1.5 1.5 0 0 0 1.5 1.5h7A1.5 1.5 0 0 0 17 19v-2.9c1.6-1 3-2.8 3-5.6a8 8 0 0 0-8-8Zm-3.2 7.2a1.9 1.9 0 1 1 0 3.8 1.9 1.9 0 0 1 0-3.8Zm6.4 0a1.9 1.9 0 1 1 0 3.8 1.9 1.9 0 0 1 0-3.8ZM12 14l1.3 2.2h-2.6Z" />
    </svg>
  );
}
