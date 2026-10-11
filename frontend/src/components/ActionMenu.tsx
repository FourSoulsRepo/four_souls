import { useState } from 'react';
import type { Action } from '../actions';
import { useApp } from '../store';
import './Actions.css';

// ActionMenu lists what the selected card can do, next to the card
// (7.7). A click anywhere else closes it.
export function ActionMenu({ actions }: { actions: Map<number, Action[]> }) {
  const selected = useApp((s) => s.selected);
  const select = useApp((s) => s.select);
  const list = selected === null ? undefined : actions.get(selected);
  if (selected === null || !list) {
    return null;
  }
  return (
    <Menu
      key={selected}
      id={selected}
      list={list}
      onClose={() => {
        select(null);
      }}
    />
  );
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
