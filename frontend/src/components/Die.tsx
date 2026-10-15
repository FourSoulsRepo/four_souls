import { useState } from 'react';
import { faceOf } from '../dice';
import './Die.css';

// Die draws one d6 value with the player's dice look (V-19). Without a
// texture, or for values no face shows, it is a plain die with a number.
// label replaces "Rolled N": a die not rolled yet names whose it is, a
// tumbling one its result (15.61, 15.62). A hidden plain die shows "?".
export function Die({ value, image, label, hidden }: { value: number; image: string; label?: string; hidden?: boolean }) {
  const [broken, setBroken] = useState('');
  const face = faceOf(value);
  const named = label ?? `Rolled ${String(value)}`;
  if (!face || !image || broken === image) {
    return (
      <div className="die die-plain" aria-label={named}>
        {hidden ? '?' : value}
      </div>
    );
  }
  return (
    <div className="die" aria-label={named} data-testid="die-image">
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
