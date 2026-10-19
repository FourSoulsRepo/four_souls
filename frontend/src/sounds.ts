// Damage sounds (15.50), all from Hearthstone. An attack lands with an
// unarmed hit, card to card. Other damage is a fireball: its cast as the
// shot leaves and its impact as it lands. A thrown die clatters on wood
// as it lands (15.69). The files live only in the private images;
// without them the game is silent.
import type { Sounds } from './bridge';

export type Sound = 'attackHit' | 'effectCast' | 'effectImpact' | 'diceRoll';

const volume = 0.6;

let urls = new Map<Sound, string>();

// loadSounds takes the sound URLs; '' is no file.
export function loadSounds(s: Sounds): void {
  urls = new Map([
    ['attackHit', s.attackHit],
    ['effectCast', s.effectCast],
    ['effectImpact', s.effectImpact],
    ['diceRoll', s.diceRoll],
  ]);
}

// playSound plays a sound once; sounds may overlap. A sound without a
// file, or one that fails to load or is blocked, stays silent.
export function playSound(sound: Sound): void {
  const url = urls.get(sound) ?? '';
  if (url === '' || typeof Audio === 'undefined') {
    return;
  }
  const a = new Audio(url);
  a.volume = volume;
  a.play().catch(() => undefined);
}
