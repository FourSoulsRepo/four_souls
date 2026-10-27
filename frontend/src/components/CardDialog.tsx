import type { ReactNode } from 'react';
import { abilityLine } from '../cards';
import { t } from '../i18n';
import { Card } from './Card';
import './Actions.css';

// CardLine is one line of a card dialog: its words, and what a click
// does. A line with no pick stays on the list, greyed; why says so.
export interface CardLine {
  text: string;
  pick?: () => void;
  why: string;
}

// CardDialog shows a card large beside its lines; the player picks one
// line (C-05). The Bone's abilities (15.24) and a "choose one-" card's
// modes (15.58) both ask this way. Cancel and a click outside call
// onCancel; buttons stand beside Cancel. backdrop is the backdrop's
// class: a menu the player opened, or a prompt the game waits on.
export function CardDialog({
  cardRef,
  name,
  label,
  lines,
  onCancel,
  backdrop,
  buttons,
}: {
  cardRef: string | undefined;
  name: string;
  label: string;
  lines: CardLine[];
  onCancel: () => void;
  backdrop: string;
  buttons?: ReactNode;
}) {
  return (
    <div className={backdrop} onClick={onCancel}>
      <div
        className="ability-dialog"
        role="dialog"
        aria-label={label}
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        {cardRef && <Card cardRef={cardRef} size="zoom" noZoom />}
        <div className="ability-lines">
          <h2>{name}</h2>
          {lines.map((line, i) => (
            <button
              // A line's index is its place on the card; it never moves.
              // eslint-disable-next-line @eslint-react/no-array-index-key
              key={i}
              type="button"
              className="ability-line"
              aria-label={line.text}
              disabled={!line.pick}
              title={line.pick ? undefined : line.why}
              onClick={line.pick}
            >
              <LineText line={line.text} />
            </button>
          ))}
          <div className="ability-buttons">
            {buttons}
            <button type="button" className="ability-cancel" onClick={onCancel}>
              {t('actions.cancel')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

// LineText draws an ability line like the card does: the cost icon in a
// box at the left, the cost in bold, then what it does.
function LineText({ line }: { line: string }) {
  const { icon, cost, effect } = abilityLine(line);
  return (
    <>
      {icon && (
        <span className={`ability-icon ability-${icon}`} aria-hidden="true">
          {icon === 'tap' ? (
            // The tap arrow, drawn: card fonts have no ↷.
            <svg viewBox="0 0 24 24">
              <path d="M5 16.5A7.5 7.5 0 1 1 17.5 11" />
              <path d="M13.5 9.5L17.5 14L21.5 9.5" />
            </svg>
          ) : (
            '$'
          )}
        </span>
      )}
      <span className="ability-words">
        {cost && <b>{cost}: </b>}
        {effect}
      </span>
    </>
  );
}
