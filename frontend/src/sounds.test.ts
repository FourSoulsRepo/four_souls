import { afterEach, describe, expect, it, vi } from 'vitest';
import { loadSounds, playSound } from './sounds';

const hit = '/cardimg/assets/sounds/attack-hit.mp3';
const cast = '/cardimg/assets/sounds/effect-cast.mp3';
const impact = '/cardimg/assets/sounds/effect-impact.mp3';
const dice = '/cardimg/assets/sounds/dice-roll.mp3';
const deal = '/cardimg/assets/sounds/loot-draw.mp3';
const death = '/cardimg/assets/sounds/player-died.mp3';

// fakeAudio records each played URL; play fails with the given error.
function fakeAudio(fail?: DOMException, volumes: number[] = []): string[] {
  const played: string[] = [];
  vi.stubGlobal(
    'Audio',
    class {
      volume = 0.5;
      constructor(public src: string) {}
      play(): Promise<void> {
        played.push(this.src);
        volumes.push(this.volume);
        return fail ? Promise.reject(fail) : Promise.resolve();
      }
    },
  );
  return played;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('playSound', () => {
  it('plays each sound from its file', () => {
    loadSounds({ attackHit: hit, effectCast: cast, effectImpact: impact, diceRoll: dice, lootDraw: deal, playerDied: death });
    const played = fakeAudio();
    playSound('attackHit');
    playSound('effectCast');
    playSound('effectImpact');
    playSound('diceRoll');
    playSound('lootDraw');
    playSound('playerDied');
    expect(played).toEqual([hit, cast, impact, dice, deal, death]);
  });

  // The files are levelled to one loudness (ADR 022); none gets its own volume.
  it('plays every sound at full volume', () => {
    loadSounds({ attackHit: hit, effectCast: cast, effectImpact: impact, diceRoll: dice, lootDraw: deal, playerDied: death });
    const volumes: number[] = [];
    fakeAudio(undefined, volumes);
    playSound('attackHit');
    playSound('effectCast');
    playSound('effectImpact');
    playSound('diceRoll');
    playSound('lootDraw');
    playSound('playerDied');
    expect(volumes).toEqual([1, 1, 1, 1, 1, 1]);
  });

  it('is silent without the private files', () => {
    loadSounds({ attackHit: '', effectCast: '', effectImpact: '', diceRoll: '', lootDraw: '', playerDied: '' });
    const played = fakeAudio();
    playSound('attackHit');
    playSound('effectImpact');
    playSound('diceRoll');
    playSound('lootDraw');
    playSound('playerDied');
    expect(played).toEqual([]);
  });

  it('swallows a file that fails to load', async () => {
    loadSounds({ attackHit: hit, effectCast: cast, effectImpact: impact, diceRoll: dice, lootDraw: deal, playerDied: death });
    const played = fakeAudio(new DOMException('no file', 'NotSupportedError'));
    playSound('attackHit');
    await Promise.resolve();
    expect(played).toEqual([hit]);
  });
});
