import { useState } from 'react';
import type { Action } from '../actions';
import { useApp } from '../store';
import { CardDialog, type CardLine } from './CardDialog';
import './Actions.css';

// ActionMenu lists what the selected card can do, next to the card
// (7.7). A card with two or more abilities opens the ability dialog
// instead (C-05). A click anywhere else closes either.
export function ActionMenu({ actions }: { actions: Map<number, Action[]> }) {
  const selected = useApp((s) => s.selected);
  const select = useApp((s) => s.select);
  const list = selected === null ? undefined : actions.get(selected);
  if (selected === null || !list) {
    return null;
  }
  const close = () => {
    select(null);
  };
  if (list.filter((a) => a.line).length > 1) {
    return <AbilityDialog key={selected} id={selected} list={list} onClose={close} />;
  }
  return <Menu key={selected} id={selected} list={list} onClose={close} />;
}

function Menu({ id, list, onClose }: { id: number; list: Action[]; onClose: () => void }) {
  const act = useApp((s) => s.act);
  // Place the menu above the card, read once when it opens.
  const [at] = useState(() => {
    const r = document.querySelector(`[data-zone="obj-${String(id)}"]`)?.getBoundingClientRect();
    return r ? { left: r.left + r.width / 2, top: r.top } : { left: window.innerWidth / 2, top: window.innerHeight / 2 };
  });
  return (
    <div className="action-backdrop" onClick={onClose}>
      <div
        className="action-menu"
        role="menu"
        style={{ left: at.left, top: at.top }}
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        {list.map((a) => (
          <button
            key={a.label}
            type="button"
            role="menuitem"
            onClick={() => {
              act(a.intent);
            }}
          >
            {a.label}
          </button>
        ))}
        <button type="button" role="menuitem" onClick={onClose}>
          Cancel
        </button>
      </div>
    </div>
  );
}

// AbilityDialog shows the card large beside its ability lines; the
// player picks one line (C-05). Lines the player may not use now stay
// on the list, greyed, so the card reads in full.
function AbilityDialog({ id, list, onClose }: { id: number; list: Action[]; onClose: () => void }) {
  const act = useApp((s) => s.act);
  const card = useApp((s) =>
    s.game?.view.players.flatMap((p) => [p.character, ...(p.in_play ?? [])]).find((c) => c.id === id),
  );
  const name = useApp((s) => (card ? (s.cards.get(card.card)?.name ?? card.card) : ''));
  // An empty line is not an ability the player uses.
  const lines = (card?.abilities ?? []).flatMap((text, i): CardLine[] => {
    if (text === '') {
      return [];
    }
    const use = list.find((a) => (a.intent.choice ?? 0) === i);
    const pick = () => {
      if (use) {
        act(use.intent);
      }
    };
    return [{ text, why: 'Not usable now', ...(use ? { pick } : {}) }];
  });
  return (
    <CardDialog
      cardRef={card?.card}
      name={name}
      label={`Use ${name}`}
      lines={lines}
      onCancel={onClose}
      backdrop="action-backdrop ability-backdrop"
    />
  );
}
