import { useState, type Ref } from 'react';
import { cubeFaces, faceOf, faceStyle } from '../dice';
import { t } from '../i18n';
import './Die.css';

// Die draws one d6 value with the player's dice look (V-19). Without a
// texture, or for values no face shows, it is a plain die with a number.
// label replaces "Rolled N": a die not rolled yet names whose it is, a
// tumbling one its result (15.61, 15.62). A hidden plain die shows "?".
// A bevel of light and shade makes it look solid (15.69).
export function Die({ value, image, label, hidden }: { value: number; image: string; label?: string; hidden?: boolean }) {
  const [broken, setBroken] = useState('');
  const named = label ?? t('dice.rolled', { n: value });
  const textured = faceOf(value) !== null && image !== '' && broken !== image;
  return (
    <div className={textured ? 'die' : 'die die-plain'} aria-label={named} data-testid={textured ? 'die-image' : undefined}>
      <Face value={value} image={textured ? image : ''} hidden={hidden} onBroken={setBroken} />
    </div>
  );
}

// Face is what one face shows: its cell of the texture, or the number
// on a plain die.
function Face({ value, image, hidden, onBroken }: { value: number; image: string; hidden?: boolean; onBroken: (image: string) => void }) {
  const face = faceOf(value);
  if (!face || !image) {
    return hidden ? '?' : value;
  }
  return (
    <img
      src={image}
      alt=""
      draggable={false}
      style={faceStyle(face)}
      onError={() => {
        onBroken(image);
      }}
    />
  );
}

// DieCube is a die as a solid cube of six faces (15.69), for a throw: it
// turns in 3D and comes to rest with the value on top (15.82). The cube is the
// element to turn; its faces are hidden from screen readers.
export function DieCube({ value, image, label, ref }: { value: number; image: string; label?: string; ref?: Ref<HTMLDivElement> }) {
  const [broken, setBroken] = useState('');
  const look = broken === image ? '' : image;
  return (
    <div ref={ref} className="die die-cube" role="img" aria-label={label ?? t('dice.rolled', { n: value })}>
      {Object.entries(cubeFaces(value)).map(([side, n]) => {
        const textured = faceOf(n) !== null && look !== '';
        return (
          <div key={side} className={`die-face die-${side}${textured ? '' : ' die-plain'}`} aria-hidden>
            <Face value={n} image={textured ? look : ''} onBroken={setBroken} />
          </div>
        );
      })}
    </div>
  );
}
