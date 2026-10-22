import type { PointerEvent } from 'react';
import { diceImage, firstRollRows, lastRollText } from '../dice';
import { useKey } from '../keys';
import type { FirstRoll, GameView, Intent, StackItem } from '../protocol';
import { playedCard, stackCard, stackIcon, stackParts, stackRows, stackWait, type LinePart, type Names } from '../stack';
import { useApp, type Zoom } from '../store';
import { playerColor } from '../table';
import { boxOf } from '../zoom';
import { Card } from './Card';
import { Die } from './Die';
import { Row, Rows } from './Rows';
import { skullPath } from './marks';
import './Actions.css';

interface Props {
  view: GameView;
  you: number;
  names: Names;
  // pass passes priority, when the engine allows it now.
  pass?: Intent;
  disabled: boolean;
  // prompt is the banner's prompt, when it is not priority (15.28).
  prompt?: string | undefined;
  // answer: Pass answers the banner's prompt, so it glows (15.28).
  answer: boolean;
}

// StackPanel is the stack inside the board (15.21, F-03): newest on
// top, one icon and one short line per item (S-01 – S-04), whom it
// waits for (V-11), and the Pass button; Space presses it (S-08). An
// open attack or buy is a row on top (15.91). Rows fade in and out and
// slide (15.51). An empty stack says so (S-09). Pass glows when it
// answers the prompt.
export function StackPanel({ view, you, names, pass, disabled, prompt, answer }: Props) {
  const act = useApp((s) => s.act);
  const rows = stackRows(view, names);
  const empty = rows.length === 0;
  const wait = stackWait(view, you, names);
  const lit = pass !== undefined && !disabled;
  const press = () => {
    if (lit) {
      act(pass);
    }
  };
  useKey(' ', press);
  return (
    <aside className={`stack${empty ? ' stack-empty' : ''}`} aria-label="Stack" data-zone="stack">
      <header>
        <h2>Stack</h2>
        <LastRollDie names={names} />
        <button
          type="button"
          className={`stack-pass${lit && answer ? ' answer' : ''}`}
          title="Pass (Space)"
          disabled={!lit}
          onClick={press}
        >
          Pass
        </button>
      </header>
      {view.first_roll ? (
        <FirstRollList roll={view.first_roll} names={names} />
      ) : empty ? (
        <p className="stack-nothing">Nothing on the stack</p>
      ) : (
        <ol>
          <Rows>
            {rows.map((r) =>
              'declared' in r ? <DeclaredRow key={r.key} parts={r.declared} odds={r.odds} /> : <StackRow key={r.key} it={r.item} view={view} names={names} />,
            )}
          </Rows>
        </ol>
      )}
      {prompt && (
        <p className="stack-wait" aria-label="Waiting for">
          {prompt}
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

// FirstRollList is the first-player roll before turn 1 (15.43, GS-06):
// each player of this round with their die once thrown, in their own
// dice look, and "…" for the next to throw. After a tie, it says so.
function FirstRollList({ roll, names }: { roll: FirstRoll; names: Names }) {
  return (
    <div className="first-roll" role="group" aria-label="First player roll">
      <p className="first-roll-title">{roll.round > 1 ? 'Tie: roll again, lowest goes first' : 'Lowest roll goes first'}</p>
      <ol>
        {firstRollRows(roll).map((r) => (
          <FirstRollRow key={r.player} player={r.player} value={r.value} next={r.next} names={names} />
        ))}
      </ol>
    </div>
  );
}

function FirstRollRow({ player, value, next, names }: { player: number; value: number | undefined; next: boolean; names: Names }) {
  const image = useApp((s) => diceImage(s.dice, s.game?.seats, player));
  const who = names.player(player);
  return (
    <li className={`first-roll-row${next ? ' next' : ''}`} aria-label={value === undefined ? `${who}: not rolled yet` : `${who}: ${String(value)}`}>
      <b className="stack-player" style={{ color: playerColor(player) }}>
        {who}
      </b>
      {value === undefined ? <span className="first-roll-wait">{next ? 'throwing…' : '…'}</span> : <Die value={value} image={image} />}
    </li>
  );
}

// PlayedCard is another player's loot card on the stack, large under
// the stack until it resolves and its effects have played (15.31,
// H-10, T-16). The view changes only after an update's effects, so it
// stays while they play.
export function PlayedCard({ view, you, names }: { view: GameView; you: number; names: Names }) {
  const it = playedCard(view, you);
  // It hides once its card flies from it to the discard (15.74), until
  // the next view (19.41).
  const leaving = useApp((s) => !!it && s.left.includes(it.source));
  if (!it?.card) {
    return null;
  }
  const who = names.player(it.controller);
  return (
    <figure
      className="played"
      aria-label={`${who} plays ${names.card(it.card)}`}
      data-zone={`played-${String(it.source)}`}
      style={leaving ? { visibility: 'hidden' } : undefined}
    >
      <Card cardRef={it.card} size="hand" />
      <figcaption>
        <b className="stack-player" style={{ color: playerColor(it.controller) }}>
          {who}
        </b>{' '}
        plays
      </figcaption>
    </figure>
  );
}

// StackRow is one item. Hover shows its card large, the used ability
// highlighted and the owner below (S-06, S-07); a tap does the same.
function StackRow({ it, view, names }: { it: StackItem; view: GameView; names: Names }) {
  const setZoom = useApp((s) => s.setZoom);
  const card = stackCard(it, view);
  const parts = stackParts(it, view, names);
  const aimed = useApp((s) => s.aiming?.targets.some((t) => t.stack_id === it.id) ?? false);
  // Its card flew away; the row hides until the next view (19.41).
  const left = useApp((s) => !!it.source && s.left.includes(it.source));
  const show = (e: PointerEvent, touch: boolean) => {
    if (card) {
      setZoom(itemZoom(it, card, view, names, e, touch));
    }
  };
  return (
    <Row
      className={`stack-row stack-${stackIcon(it)}${aimed ? ' hl-aim' : ''}`}
      // Damage from this item's source flies from here (15.37).
      data-zone={it.source ? `stack-src-${String(it.source)}` : undefined}
      // A dragged card may aim at this item (15.25).
      data-stack={it.id}
      aria-label={parts.map((p) => p.text).join('')}
      style={left ? { visibility: 'hidden' } : undefined}
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
    </Row>
  );
}

// DeclaredRow is an open attack or buy (15.91, V-51). It is no stack
// item (R-ATK-02, R-SHOP-02), so it has no icon frame and a dashed edge.
// An aimed attack adds its fight line: what the roll needs (15.94).
function DeclaredRow({ parts, odds }: { parts: LinePart[]; odds?: string | undefined }) {
  const line = parts.map((p) => p.text).join('');
  return (
    <Row className="stack-declared" aria-label={odds ? `${line}. ${odds}` : line}>
      <span className="stack-line">
        {parts.map((p, i) => (
          // A line's parts never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <Part key={i} part={p} />
        ))}
      </span>
      {odds && <span className="stack-odds">{odds}</span>}
    </Row>
  );
}

// itemZoom is the large card of a stack item: the used ability's words
// and the owner in their color below it (S-06, S-07). An item without
// words says what it did instead, naming its target (15.22). It opens
// beside the row (15.27).
export function itemZoom(it: StackItem, card: string, view: GameView, names: Names, e: PointerEvent, touch: boolean): Zoom {
  const owner = it.controller >= 0 ? names.player(it.controller) : undefined;
  return {
    ref: card,
    touch,
    at: boxOf(e.currentTarget),
    caption: {
      owner,
      color: playerColor(it.controller),
      ...(it.text ? { ability: it.text } : { line: stackParts(it, view, names) }),
    },
  };
}

// ItemIcon is a stack item's square icon (S-04): card art, a die in the
// roller's dice look (15.6), a blood drop with the damage, or a skull.
export function ItemIcon({ it }: { it: StackItem }) {
  const image = useApp((s) => diceImage(s.dice, s.game?.seats, it.controller));
  const icon = stackIcon(it);
  return (
    <span className="stack-icon" aria-hidden="true">
      {icon === 'card' ? (
        it.card && <Card cardRef={it.card} size="small" noZoom />
      ) : icon === 'roll' ? (
        it.roll === undefined ? (
          <span className="stack-die">?</span>
        ) : (
          <Die value={it.roll} image={image} />
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

// Part is one piece of a stack line, colored by its role (S-03).
export function Part({ part }: { part: LinePart }) {
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

export function SkullIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d={skullPath} />
    </svg>
  );
}

// LastRollDie is the last roll as a small die by the title (15.63, F-21):
// its final value, after changes and rerolls. Hover or focus names the
// roller and what the roll was for.
function LastRollDie({ names }: { names: Names }) {
  const last = useApp((s) => s.lastRoll);
  const image = useApp((s) => diceImage(s.dice, s.game?.seats, s.lastRoll?.player ?? -1));
  if (!last) {
    return null;
  }
  const text = lastRollText(last, names.player);
  return (
    <span className="last-roll" role="group" aria-label="Last roll" tabIndex={0}>
      <Die value={last.value} image={image} />
      <span className="last-roll-tip" role="tooltip">
        {text}
      </span>
    </span>
  );
}
