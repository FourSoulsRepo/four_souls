import { useState } from 'react';
import { discardIntent, type Actions } from '../actions';
import type { Looked } from '../history';
import type { CardView } from '../protocol';
import { useApp } from '../store';
import { Card } from './Card';
import './Actions.css';

// ChoiceDialog asks the player to pick one option. Card options show as
// cards, Hearthstone-style (V-07); options on the table can be clicked
// there too. subject is the card the choice is about, e.g. the monster
// a deck attack drew (19.9) or the deck top a player looks at (19.11);
// it shows beside the options, and option labels name it by its title.
export function ChoiceDialog({ choice, subject }: { choice: NonNullable<Actions['choice']>; subject?: string | undefined }) {
  const act = useApp((s) => s.act);
  const subjectName = subject ? useApp.getState().cards.get(subject)?.name : undefined;
  const label = (l: string) => (subject && subjectName ? l.replaceAll(subject, subjectName) : l);
  return (
    <div className="choice" role="dialog" aria-label={choice.title}>
      <h2>{choice.title}</h2>
      {subject && (
        <div className="choice-subject">
          <Card cardRef={subject} size="table" />
        </div>
      )}
      <div className="choice-options">
        {choice.options.map((o) => (
          <div key={o.intent.choice ?? 0} className="choice-option">
            {o.card && (
              <Card
                cardRef={o.card}
                size="table"
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
              {o.card ? (useApp.getState().cards.get(o.card)?.name ?? o.label) : label(o.label)}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

// LookDialog shows the cards an effect let this player look at,
// Hearthstone-style (V-07), until they close it.
export function LookDialog({ cards }: { cards: Looked[] }) {
  const endLook = useApp((s) => s.endLook);
  return (
    <div className="choice look" role="dialog" aria-label="You look at">
      <h2>You look at</h2>
      <div className="choice-options">
        {cards.map((c) => (
          <Card key={c.id} cardRef={c.card} size="hand" />
        ))}
      </div>
      <button type="button" onClick={endLook}>
        OK
      </button>
    </div>
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
    <div className="choice" role="dialog" aria-label="Discard">
      <h2>
        Discard {discard.count} loot card{discard.count === 1 ? '' : 's'}
      </h2>
      <div className="choice-options">
        {cards.map((c) => (
          <Card
            key={c.id}
            cardRef={c.card}
            size="table"
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
        Discard ({chosen.length}/{discard.count})
      </button>
    </div>
  );
}
