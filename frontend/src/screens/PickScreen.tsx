import { useState } from 'react';
import { Card } from '../components/Card';
import { pickToggle } from '../choice';
import { t } from '../i18n';
import type { Setup } from '../protocol';
import { playSound } from '../sounds';
import { seatOf, useApp } from '../store';
import './screens.css';

// PickScreen is a draft's pick (GS-01, 14.11): the player sees only
// their own offers, each with its starting item, picks one and confirms.
// A character without a fixed item (Eden) says it picks one at the
// start (R-SETUP-09). Others see who has picked, not what.
export function PickScreen({ setup }: { setup: Setup }) {
  const you = useApp(seatOf);
  const pick = useApp((s) => s.pick);
  const quit = useApp((s) => s.quit);
  const error = useApp((s) => s.error);
  const cards = useApp((s) => s.cards);
  const name = (ref: string) => cards.get(ref)?.name ?? ref;
  const [picked, setPicked] = useState<number | null>(null);
  const [sent, setSent] = useState(false);
  const offers = setup.offers ?? [];
  const items = setup.items ?? [];
  // done: the server has this seat's pick, or it is on its way. A
  // refusal (an error) lets the player pick again.
  const done = (setup.picked?.at(you) ?? false) || (sent && !error);
  const card = picked === null ? undefined : offers.at(picked);
  return (
    <div className="screen pick">
      <h1>{t('pick.title')}</h1>
      <p className="muted">{done ? t('pick.waiting') : t('pick.caption')}</p>
      {!done && (
        <div className="pick-offers">
          {offers.map((c, i) => {
            const item = items.at(i) ?? '';
            return (
              <div key={c} className="pick-offer">
                <button
                  type="button"
                  className={`pick-card${picked === i ? ' picked' : ''}`}
                  aria-label={name(c)}
                  aria-pressed={picked === i}
                  data-testid={`offer-${c}`}
                  onClick={() => {
                    const next = pickToggle(picked, i);
                    playSound(next.sound);
                    setPicked(next.picked);
                  }}
                >
                  <Card cardRef={c} size="choice" highlight="pick" />
                </button>
                <div className="pick-item" data-testid={`item-${c}`}>
                  {item ? (
                    <>
                      <span className="muted">{t('pick.item')}</span>
                      <Card cardRef={item} size="table" />
                      <span className="pick-item-name">{name(item)}</span>
                    </>
                  ) : (
                    <span className="muted pick-item-later">{t('pick.item_later')}</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
      {!done && (
        <button
          type="button"
          className="primary"
          disabled={!card}
          onClick={() => {
            if (!card) {
              return;
            }
            playSound('choiceConfirm');
            useApp.setState({ error: '' });
            setSent(true);
            pick(card);
          }}
        >
          {t('pick.confirm')}
        </button>
      )}
      <ul className="pick-seats" aria-label={t('pick.players')}>
        {setup.seats.map((s) => (
          <li key={s.seat} className={setup.picked?.at(s.seat) ? 'picked' : ''}>
            {s.name}: {setup.picked?.at(s.seat) ? t('pick.has_picked') : t('pick.choosing')}
          </li>
        ))}
      </ul>
      <p className="status" role="status">
        {error}
      </p>
      <button
        type="button"
        onClick={() => {
          void quit();
        }}
      >
        {t('game.leave')}
      </button>
    </div>
  );
}
