import { useState, type ReactNode } from 'react';
import { discardIntent, modeLines, orderIntent, outcomeLines, type Actions } from '../actions';
import { lookGroups, type Looked } from '../history';
import { t } from '../i18n';
import type { CardView, ModeView } from '../protocol';
import { useApp } from '../store';
import { Card } from './Card';
import { CardDialog } from './CardDialog';
import './Actions.css';

// ChoiceDialog asks the player to pick one option. Card options show as
// cards, Hearthstone-style (V-07); options on the table can be clicked
// there too. subject is the card the choice is about, e.g. the monster
// a deck attack drew (19.9) or the deck top a player looks at (19.11);
// it shows beside the options; option labels name it as a card (ADR 021).
// An option with an outcome says what picking it gives (19.17); a
// trigger option says its ability's words (19.19).
//
// A choice shows big cards over the dimmed board (15.29, H-11); "View
// board" hides it to look at the table, where its targets still glow
// and answer, and "Back to the choice" brings it back (F-20).
export function ChoiceDialog({ choice, subject }: { choice: NonNullable<Actions['choice']>; subject?: string | undefined }) {
  const act = useApp((s) => s.act);
  const cardName = (ref: string) => useApp.getState().cards.get(ref)?.name ?? ref;
  return (
    <ChoiceFrame label={choice.title}>
      {subject && (
        <div className="choice-subject">
          <Card cardRef={subject} size="choice" />
        </div>
      )}
      <div className="choice-options">
        {choice.options.map((o) => (
          <div key={o.intent.choice ?? 0} className="choice-option">
            {o.card && (
              <Card
                cardRef={o.card}
                size="choice"
                highlight="pick"
                onClick={() => {
                  act(o.intent);
                }}
              />
            )}
            <button
              type="button"
              onClick={() => {
                act(o.intent);
              }}
            >
              {o.card ? cardName(o.card) : o.label}
            </button>
            {o.line && <p className="choice-line">{o.line}</p>}
            {o.outcome && (
              <ul className="choice-outcome" aria-label={t('choice.result')}>
                {outcomeLines(o.outcome, cardName).map((l) => (
                  <li key={l}>→ {l}</li>
                ))}
              </ul>
            )}
          </div>
        ))}
      </div>
    </ChoiceFrame>
  );
}

// ChoiceFrame is a choice over the dimmed board, or, while the player
// views the board, a small bar to bring it back (15.29, F-20, H-11).
// title is the heading when it says more than the label.
function ChoiceFrame({ label, title, className = '', children }: { label: string; title?: ReactNode; className?: string; children: ReactNode }) {
  const [board, setBoard] = useState(false);
  if (board) {
    return (
      <Peek
        title={title ?? label}
        back={() => {
          setBoard(false);
        }}
      />
    );
  }
  return (
    <div className="choice-backdrop">
      <div className={`choice${className ? ` ${className}` : ''}`} role="dialog" aria-label={label}>
        <h2>{title ?? label}</h2>
        {children}
        <button
          type="button"
          className="choice-view"
          onClick={() => {
            setBoard(true);
          }}
        >
          {t('choice.view_board')}
        </button>
      </div>
    </div>
  );
}

// Peek is a hidden choice: its title and a button to bring it back.
function Peek({ title, back }: { title: ReactNode; back: () => void }) {
  return (
    <div className="choice-peek" role="group" aria-label={t('choice.viewing')}>
      <span>{title}</span>
      <button type="button" onClick={back}>
        {t('choice.back')}
      </button>
    </div>
  );
}

// ModeDialog asks a "choose one-" card's mode the way The Bone asks its
// ability: the card large, one line per mode (15.58, V-45). A mode with
// no valid target now stays on the list, greyed. Cancel and a click
// outside send the prompt's cancel; "View board" hides it, like any
// choice (F-20).
export function ModeDialog({ choice, card, modes }: { choice: NonNullable<Actions['choice']>; card: string | undefined; modes: ModeView[] }) {
  const act = useApp((s) => s.act);
  const [board, setBoard] = useState(false);
  const name = useApp((s) => (card ? (s.cards.get(card)?.name ?? card) : choice.title));
  if (board) {
    return (
      <Peek
        title={choice.title}
        back={() => {
          setBoard(false);
        }}
      />
    );
  }
  const { lines, cancel } = modeLines(choice, modes);
  return (
    <CardDialog
      cardRef={card}
      name={name}
      label={choice.title}
      lines={lines.map(({ text, intent }) => ({
        text,
        why: t('choice.no_target'),
        ...(intent
          ? {
              pick: () => {
                act(intent);
              },
            }
          : {}),
      }))}
      onCancel={() => {
        if (cancel) {
          act(cancel);
        }
      }}
      backdrop="choice-backdrop"
      buttons={
        <button
          type="button"
          className="ability-view"
          onClick={() => {
            setBoard(true);
          }}
        >
          {t('choice.view_board')}
        </button>
      }
    />
  );
}

// LookDialog shows the cards an effect let this player look at,
// Hearthstone-style (V-07), until they close it. Hand cards come in
// one labelled row per hand, so the player sees whose they are (19.71).
export function LookDialog({ cards, you, name }: { cards: Looked[]; you: number; name: (seat: number) => string }) {
  const endLook = useApp((s) => s.endLook);
  const row = (cs: Looked[], key?: string) => (
    <div key={key} className="choice-options">
      {cs.map((c) => (
        <Card key={c.id} cardRef={c.card} size="hand" />
      ))}
    </div>
  );
  return (
    <div className="choice look" role="dialog" aria-label={t('look.title')}>
      <h2>{t('look.title')}</h2>
      {lookGroups(cards).map((g) => {
        if (g.owner === undefined) {
          return row(g.cards, 'deck');
        }
        const label = g.owner === you ? t('look.your_hand') : t('look.hand_of', { player: name(g.owner) });
        return (
          <section key={g.owner} className="look-hand" role="group" aria-label={label}>
            <h3>{label}</h3>
            {row(g.cards)}
          </section>
        );
      })}
      <button type="button" onClick={endLook}>
        {t('look.ok')}
      </button>
    </div>
  );
}

// OrderDialog puts the cards of an order prompt back on the deck (19.16).
// A click puts a card in the next free place, top first; a click on a
// placed card takes it back. Nothing is sent until every card has a
// place, so the player sees the whole order before it counts.
export function OrderDialog({ order, cards }: { order: NonNullable<Actions['order']>; cards: CardView[] }) {
  const act = useApp((s) => s.act);
  const [placed, setPlaced] = useState<number[]>([]);
  const left = cards.filter((c) => !placed.includes(c.id));
  // One place per card, numbered from the top.
  const places = cards.map((_, i) => ({ n: i + 1, card: cards.find((x) => x.id === placed.at(i)) }));
  return (
    <ChoiceFrame label={order.title} className="order">
      <p className="order-hint">{t('order.hint')}</p>
      <div className="order-deck">
        {places.map(({ n, card: c }) => (
          <div key={n} className="order-place" role="group" aria-label={t('order.place', { n })}>
            {c ? (
              <Card
                cardRef={c.card}
                size="table"
                onClick={() => {
                  setPlaced(placed.filter((id) => id !== c.id));
                }}
              />
            ) : (
              <div className="order-empty" />
            )}
            <span>{n === 1 ? t('order.top') : String(n)}</span>
          </div>
        ))}
      </div>
      {left.length > 0 && (
        <div className="choice-options" role="group" aria-label={t('order.not_placed')}>
          {left.map((c) => (
            <Card
              key={c.id}
              cardRef={c.card}
              size="table"
              highlight="pick"
              onClick={() => {
                setPlaced([...placed, c.id]);
              }}
            />
          ))}
        </div>
      )}
      <button
        type="button"
        disabled={left.length > 0}
        onClick={() => {
          act(orderIntent(order.intent, placed));
        }}
      >
        {t('order.put_back')}
      </button>
    </ChoiceFrame>
  );
}

// DiscardDialog asks for exactly count cards from the hand.
export function DiscardDialog({
  discard,
  hand,
}: {
  discard: NonNullable<Actions['discard']>;
  hand: CardView[];
}) {
  const act = useApp((s) => s.act);
  const [chosen, setChosen] = useState<number[]>([]);
  const cards = hand.filter((c) => discard.from.includes(c.id));
  const toggle = (id: number) => {
    setChosen(chosen.includes(id) ? chosen.filter((x) => x !== id) : [...chosen, id]);
  };
  return (
    <ChoiceFrame label={t('discard.label')} title={t('discard.title', { n: discard.count })}>
      <div className="choice-options">
        {cards.map((c) => (
          <Card
            key={c.id}
            cardRef={c.card}
            size="choice"
            highlight={chosen.includes(c.id) ? 'pick' : 'act'}
            onClick={() => {
              toggle(c.id);
            }}
          />
        ))}
      </div>
      <button
        type="button"
        disabled={chosen.length !== discard.count}
        onClick={() => {
          act(discardIntent(discard.intent, chosen));
        }}
      >
        {t('discard.button', { chosen: chosen.length, n: discard.count })}
      </button>
    </ChoiceFrame>
  );
}
