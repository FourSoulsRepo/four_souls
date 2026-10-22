import { skyLayers, skyStyle } from '../sky';
import { useApp } from '../store';
import './Sky.css';

// Sky is the drifting night sky behind the menu screens (15.42); without
// the images, nothing, and the screens keep their plain color.
export function Sky() {
  const images = useApp((s) => s.icons.sky);
  const speed = useApp((s) => s.settings.animation);
  const layers = skyLayers(images, speed);
  if (layers.length === 0) {
    return null;
  }
  return (
    <div className="sky" aria-hidden="true">
      {layers.map((l) => (
        <div key={l.image} className="sky-layer" style={skyStyle(l)} />
      ))}
    </div>
  );
}
