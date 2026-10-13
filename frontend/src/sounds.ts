// Damage sounds (15.50). An attack lands with a card hit, a card slapped
// on another (Kenney, CC0). Other damage is a fireball: Hearthstone's
// cast as the shot leaves and its impact as it lands. The fireball sounds
// live only in the private images; without them the impact is the card
// hit and the cast is silent.
import cardHit from './assets/sounds/card-hit.mp3';
import type { Sounds } from './bridge';

export type Sound = 'hit' | 'cast' | 'impact';

const volume = 0.6;

let urls = new Map<Sound, string>([['hit', cardHit]]);

// broken are the URLs that failed to load; their sound falls back.
const broken = new Set<string>();

// loadSounds takes the private sound URLs.
export function loadSounds(s: Sounds): void {
  urls = new Map([
    ['hit', cardHit],
    ['cast', s.fireballCast],
    ['impact', s.fireballImpact],
  ]);
}

// urlFor is the file a sound plays, '' for none.
export function urlFor(sound: Sound): string {
  const url = urls.get(sound) ?? '';
  if (url !== '' && !broken.has(url)) {
    return url;
  }
  return sound === 'impact' ? urlFor('hit') : '';
}

// playSound plays a sound once; sounds may overlap. A file that fails to
// load is marked broken, and its fallback plays instead.
export function playSound(sound: Sound): void {
  const url = urlFor(sound);
  if (url === '' || typeof Audio === 'undefined') {
    return;
  }
  const a = new Audio(url);
  a.volume = volume;
  a.play().catch((e: unknown) => {
    // A missing file rejects with NotSupportedError; a blocked autoplay
    // with NotAllowedError, which is not the file's fault.
    if (e instanceof DOMException && e.name === 'NotSupportedError') {
      broken.add(url);
      if (urlFor(sound) !== '') {
        playSound(sound);
      }
    }
  });
}
