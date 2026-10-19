import { useLayoutEffect, useRef, useState } from 'react';
import { useApp } from '../store';
import { keyword, placeZoom, type Place } from '../zoom';
import { Card, CardText } from './Card';
import { Part } from './StackPanel';
import './Card.css';

// CardZoom shows the zoomed card. A mouse zoom opens right beside the
// hovered card, or rises from a hand card, and lets clicks through
// (15.27); a touch zoom stays centered until the next tap. A caption
// names the owner and highlights the used ability (S-06, S-07); the
// card's keywords stand beside it (H-13). A card the player can't use
// now says why under it (15.29, F-17, C-06).
export function CardZoom() {
  const zoom = useApp((s) => s.zoom);
  const setZoom = useApp((s) => s.setZoom);
  const footnotes = useApp((s) => (zoom ? s.cards.get(zoom.ref)?.footnotes : undefined));
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
        {footnotes && footnotes.length > 0 && (
          <aside className="card-zoom-tips" role="note" aria-label="Keywords">
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
