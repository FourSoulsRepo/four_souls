// Damage sounds (15.50), all from Hearthstone. An attack lands with an
// unarmed hit, card to card. Other damage is a fireball: its cast as the
// shot leaves and its impact as it lands. A thrown die clatters on wood
// as it lands (15.69), each die of several on its own (19.39). Each
// looted card is dealt with a card sound (18.7). A player's death sounds
// with its skull (18.8). Each gained treasure sounds as it flies in
// (18.9). The files live only in the private images; without them the
// game is silent.
// Every file is levelled to one loudness (ADR 022), so every sound plays
// at the same level. The player's volume (ST-09) scales them all alike.
import type { Sounds } from './bridge';

export type Sound = 'attackHit' | 'effectCast' | 'effectImpact' | 'diceRoll' | 'lootDraw' | 'playerDied' | 'treasureGain';

// level is the player's volume, 0 to 1 (ST-09).
let level = 1;

let urls = new Map<Sound, string>();

// loadSounds takes the sound URLs; '' is no file.
export function loadSounds(s: Sounds): void {
  urls = new Map([
    ['attackHit', s.attackHit],
    ['effectCast', s.effectCast],
    ['effectImpact', s.effectImpact],
    ['diceRoll', s.diceRoll],
    ['lootDraw', s.lootDraw],
    ['playerDied', s.playerDied],
    ['treasureGain', s.treasureGain],
  ]);
}

// setSoundVolume sets the player's volume, in percent; 0 is mute.
export function setSoundVolume(percent: number): void {
  level = percent / 100;
}

// playSound plays a sound once; sounds may overlap. A sound without a
// file, or one that fails to load or is blocked, stays silent.
export function playSound(sound: Sound): void {
  const url = urls.get(sound) ?? '';
  if (url === '' || typeof Audio === 'undefined') {
    return;
  }
  const a = new Audio(url);
  a.volume = level;
  a.play().catch(() => undefined);
}
