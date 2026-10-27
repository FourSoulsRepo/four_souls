import { useLayoutEffect, useRef, useState } from 'react';
import { language, t } from '../i18n';
import { useApp } from '../store';
import { keyword, placeZoom, type Place } from '../zoom';
import type { Outcome } from '../history';
import { translatorsOf } from '../cards';
import { Card, CardText } from './Card';
import { starPolygon } from './lunge';
import { skullPath } from './marks';
import { Part } from './StackPanel';
import './Card.css';

// CardZoom shows the zoomed card. A mouse zoom opens right beside the
// hovered card, or rises from a hand card, and lets clicks through
// (15.27); a touch zoom stays centered until the next tap. A caption
// names the owner and highlights the used ability (S-06, S-07); the
// card's keywords stand beside it (H-13). A card the player can't use
// now says why under it (15.29, F-17, C-06). A history headline's zoom
// shows what it did beside the card, after an arrow (15.90). A changed
// stat's sum goes under the card (15.92). A translated card credits its
// translators under it (CD-06, 16.5). A two-sided card shows its other
// side beside it, smaller (ADR 023).
export function CardZoom() {
  const zoom = useApp((s) => s.zoom);
  const setZoom = useApp((s) => s.setZoom);
  const info = useApp((s) => (zoom ? s.cards.get(zoom.ref) : undefined));
  const footnotes = info?.footnotes;
  const otherSide = info?.back ?? info?.front;
  // Who translated the card's name into the shown language (CD-06).
  const translators = info ? translatorsOf(info, language()) : [];
  const bodyRef = useRef<HTMLDivElement>(null);
  const [place, setPlace] = useState<Place | null>(null);
  // The zoom is measured before it paints, then placed by its size.
  useLayoutEffect(() => {
    const el = bodyRef.current;
    if (!zoom || zoom.touch || !el) {
      setPlace(null);
      return;
    }
    const size = { width: el.offsetWidth, height: el.offsetHeight };
    setPlace(placeZoom(zoom.at, size, { width: window.innerWidth, height: window.innerHeight }, zoom.rise ?? false));
  }, [zoom]);
  if (!zoom) {
    return null;
  }
  const placed = !zoom.touch && place;
  return (
    <div
      className={`card-zoom${zoom.touch ? ' card-zoom-touch' : ''}`}
      data-testid="card-zoom"
      onPointerDown={() => {
        setZoom(null);
      }}
    >
      <div
        ref={bodyRef}
        className={`card-zoom-body tips-${place?.tips ?? 'right'}`}
        style={
          zoom.touch
            ? undefined
            : { transform: placed ? `translate(${String(place.x)}px, ${String(place.y)}px)` : undefined, visibility: placed ? undefined : 'hidden' }
        }
      >
        <figure className="card-zoom-figure">
          <Card cardRef={zoom.ref} size="zoom" />
          {zoom.why && (
            <p className="card-zoom-why" role="note">
              {zoom.why}
            </p>
          )}
          {zoom.stats && zoom.stats.length > 0 && (
            <ul className="card-zoom-stats" aria-label={t('zoom.stats')}>
              {zoom.stats.map((line) => (
                <li key={line}>{line}</li>
              ))}
            </ul>
          )}
          {translators.length > 0 && <p className="card-zoom-credit">{t('zoom.translators', { names: translators.join(', ') })}</p>}
          {zoom.caption && (
            <figcaption className="card-zoom-caption">
              {zoom.caption.ability && <p className="card-zoom-ability">{zoom.caption.ability}</p>}
              {zoom.caption.line && (
                <p className="card-zoom-ability card-zoom-line">
                  {zoom.caption.line.map((p, i) => (
                    // A line's parts never change order, so the index is a stable key.
                    // eslint-disable-next-line @eslint-react/no-array-index-key
                    <Part key={i} part={p} />
                  ))}
                </p>
              )}
              {zoom.caption.owner && (
                <p className="card-zoom-owner" style={{ color: zoom.caption.color }}>
                  {zoom.caption.owner}
                </p>
              )}
            </figcaption>
          )}
        </figure>
        {otherSide && (
          <figure className="card-zoom-side">
            <Card cardRef={otherSide} size="hand" noZoom />
            <figcaption>{t('zoom.other_side')}</figcaption>
          </figure>
        )}
        {zoom.outcome && zoom.outcome.length > 0 && <Outcomes list={zoom.outcome} />}
        {footnotes && footnotes.length > 0 && (
          <aside className="card-zoom-tips" role="note" aria-label={t('zoom.keywords')}>
            {footnotes.map((f) => (
              <p key={f.title} className="card-zoom-tip">
                <b>{keyword(f.title)}</b>
                <CardText text={f.text} />
              </p>
            ))}
          </aside>
        )}
      </div>
    </div>
  );
}

// Outcomes is what a play did (15.90, V-51): an arrow, then each card it
// hit or gave, with a starburst and the damage, or a skull; cents are a
// plain coin with the amount.
function Outcomes({ list }: { list: Outcome[] }) {
  return (
    <div className="card-zoom-outcome" role="list" aria-label={t('zoom.outcome')}>
      <svg className="card-zoom-arrow" viewBox="0 0 24 24" aria-hidden="true">
        <path d="M3 10h12V5l7 7-7 7v-5H3Z" />
      </svg>
      <div className="card-zoom-results">
        {list.map((o, i) => (
          // A popout's results never change order, so the index is a stable key.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <figure key={i} className={`card-zoom-result${o.mark ? ` result-${o.mark}` : ''}`} role="listitem" aria-label={resultLabel(o)}>
            {o.card ? (
              <Card cardRef={o.card} size="hand" noZoom />
            ) : o.mark === 'cents' ? (
              <span className="result-coin" aria-hidden="true">
                {o.text}
              </span>
            ) : (
              <span className="result-text">{o.text}</span>
            )}
            {o.mark === 'damage' && (
              <span className="result-burst" aria-hidden="true">
                <svg viewBox="0 0 100 100">
                  <polygon points={starPolygon()} />
                </svg>
                <b>−{o.amount}</b>
              </span>
            )}
            {o.mark === 'death' && (
              <span className="result-skull" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d={skullPath} />
                </svg>
              </span>
            )}
          </figure>
        ))}
      </div>
    </div>
  );
}

// resultLabel names one result: "Gaper: 2 damage", "Gaper dies", "+3¢".
function resultLabel(o: Outcome): string {
  if (o.mark === 'damage') return t('zoom.result_damage', { target: o.text, n: o.amount ?? 0 });
  if (o.mark === 'death') return t('zoom.result_death', { target: o.text });
  return o.text;
}
