import { afterEach, describe, expect, it, vi } from 'vitest';
import { loadSounds, playSound, urlFor } from './sounds';

const cast = '/cardimg/assets/sounds/fireball-cast.mp3';
const impact = '/cardimg/assets/sounds/fireball-impact.mp3';

// fakeAudio records each played URL; play fails for the URLs in fail.
function fakeAudio(fail: string[]): string[] {
  const played: string[] = [];
  vi.stubGlobal(
    'Audio',
    class {
      volume = 1;
      constructor(public src: string) {}
      play(): Promise<void> {
        played.push(this.src);
        return fail.includes(this.src) ? Promise.reject(new DOMException('no file', 'NotSupportedError')) : Promise.resolve();
      }
    },
  );
  return played;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('urlFor', () => {
  it('plays the fireball with the private sounds', () => {
    loadSounds({ fireballCast: cast, fireballImpact: impact });
    expect(urlFor('cast')).toBe(cast);
    expect(urlFor('impact')).toBe(impact);
    expect(urlFor('hit')).toMatch(/card-hit\.mp3$/);
  });

  it('plays the card hit on impact without them, and no cast', () => {
    loadSounds({ fireballCast: '', fireballImpact: '' });
    expect(urlFor('cast')).toBe('');
    expect(urlFor('impact')).toBe(urlFor('hit'));
  });
});

describe('playSound', () => {
  it('falls back to the card hit when a file fails to load', async () => {
    const missing = '/cardimg/missing-impact.mp3';
    loadSounds({ fireballCast: '', fireballImpact: missing });
    const played = fakeAudio([missing]);
    playSound('impact');
    await vi.waitFor(() => {
      expect(played).toEqual([missing, urlFor('hit')]);
    });
    // The broken file is not tried again.
    playSound('impact');
    expect(played.at(-1)).toBe(urlFor('hit'));
  });

  it('keeps a file that autoplay blocked', async () => {
    loadSounds({ fireballCast: cast, fireballImpact: impact });
    const played: string[] = [];
    vi.stubGlobal(
      'Audio',
      class {
        volume = 1;
        constructor(public src: string) {}
        play(): Promise<void> {
          played.push(this.src);
          return Promise.reject(new DOMException('no gesture', 'NotAllowedError'));
        }
      },
    );
    playSound('cast');
    await Promise.resolve();
    expect(played).toEqual([cast]);
    expect(urlFor('cast')).toBe(cast);
  });
});
