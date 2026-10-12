import { useState } from 'react';
import { faceOf } from '../dice';
import './Die.css';

// Die draws one d6 value with the player's dice look (V-19). Without a
// texture, or for values no face shows, it is a plain die with a number.
export function Die({ value, image }: { value: number; image: string }) {
  const [broken, setBroken] = useState('');
  const face = faceOf(value);
  const label = `Rolled ${String(value)}`;
  if (!face || !image || broken === image) {
    return (
      <div className="die die-plain" aria-label={label}>
        {value}
      </div>
    );
  }
  return (
    <div className="die" aria-label={label} data-testid="die-image">
      <img
        src={image}
        alt=""
        draggable={false}
        style={{ left: `${String(-face.col * 100)}%`, top: `${String(-face.row * 100)}%` }}
        onError={() => {
          setBroken(image);
        }}
      />
    </div>
  );
}
