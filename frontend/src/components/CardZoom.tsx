import { useApp } from '../store';
import { Card } from './Card';
import { Part } from './StackPanel';
import './Card.css';

// CardZoom shows the zoomed card. A mouse zoom follows the hover and
// lets clicks through; a touch zoom stays until the next tap. A caption
// names the owner and highlights the used ability (S-06, S-07).
export function CardZoom() {
  const zoom = useApp((s) => s.zoom);
  const setZoom = useApp((s) => s.setZoom);
  if (!zoom) {
    return null;
  }
  return (
    <div
      className={`card-zoom card-zoom-${zoom.side}${zoom.touch ? ' card-zoom-touch' : ''}`}
      data-testid="card-zoom"
      onPointerDown={() => {
        setZoom(null);
      }}
    >
      <figure className="card-zoom-figure">
        <Card cardRef={zoom.ref} size="zoom" {...zoom.marks} />
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
    </div>
  );
}
