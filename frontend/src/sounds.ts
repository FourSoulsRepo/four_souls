// Damage sounds (15.50), all from Hearthstone. An attack lands with an
// unarmed hit, card to card. Other damage is a fireball: its cast as the
// shot leaves and its impact as it lands. A thrown die clatters on wood
// as it lands (15.69), each die of several on its own (19.39). Each
// looted card is dealt with a card sound (18.7). A player's death sounds
// with its skull (18.8). Each gained treasure sounds as it flies in
// (18.9). Every soul gained sounds with one sound: as it flies to its
// player, a bonus soul too (18.13), and as the Lost Soul turns to one
// (18.14). The reference client's sounds of step 20 (20.17) are here too;
// each later sub-step of step 20 plays its own. A sound may have several
// files: each play picks one at random, all alike, as the reference
// does. The files live only in the private images; without them the
// game is silent.
// Every file is levelled to one loudness (ADR 022), so every sound plays
// at the same level. The player's volume (ST-09) scales them all alike.
import type { Sounds } from './bridge';

export type Sound = keyof Sounds;

// level is the player's volume, 0 to 1 (ST-09).
let level = 1;

let urls = new Map<Sound, string[]>();

// loadSounds takes the sound URLs, a list per sound; an empty list is no file.
export function loadSounds(s: Sounds): void {
  urls = new Map(Object.entries(s) as [Sound, string[]][]);
}

// setSoundVolume sets the player's volume, in percent; 0 is mute.
export function setSoundVolume(percent: number): void {
  level = percent / 100;
}

// playSound plays a sound once, one of its files at random; sounds may
// overlap. A sound without a file, or one that fails to load or is
// blocked, stays silent.
export function playSound(sound: Sound): void {
  const files = urls.get(sound) ?? [];
  if (files.length === 0 || typeof Audio === 'undefined') {
    return;
  }
  const a = new Audio(files[Math.floor(Math.random() * files.length)]);
  a.volume = level;
  a.play().catch(() => undefined);
}
