import { useApp } from '../store';
import { Card } from './Card';
import './Card.css';

// CardZoom shows the zoomed card. A mouse zoom follows the hover and
// lets clicks through; a touch zoom stays until the next tap.
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
      <Card cardRef={zoom.ref} size="zoom" />
    </div>
  );
}
