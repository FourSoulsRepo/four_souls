// Media players (19.135): each played file holds one, and Chromium and
// WebView2 cap them per page. A long game plays thousands of sounds;
// past the cap the browser blocks new ones and logs a console error.
import { expect, test } from './fixtures';

// wav is a silent WAV file: 8 kHz, 16 bit, mono, 6 s. Each play lasts
// past the test's end, so none can give its player back by ending, and
// the page keeps every playing element.
function wav(): Buffer {
  const samples = 48_000;
  const b = Buffer.alloc(44 + samples * 2);
  b.write('RIFF', 0);
  b.writeUInt32LE(36 + samples * 2, 4);
  b.write('WAVEfmt ', 8);
  b.writeUInt32LE(16, 16);
  b.writeUInt16LE(1, 20);
  b.writeUInt16LE(1, 22);
  b.writeUInt32LE(8000, 24);
  b.writeUInt32LE(16000, 28);
  b.writeUInt16LE(2, 32);
  b.writeUInt16LE(16, 34);
  b.write('data', 36);
  b.writeUInt32LE(samples * 2, 40);
  return b;
}

test('1200 sounds at once never run out of media players', async ({ page }) => {
  const errors: string[] = [];
  page.on('console', (m) => {
    if (m.type() === 'error') {
      errors.push(m.text());
    }
  });
  await page.route(/\/test-sounds\/.+\.wav$/, (route) => route.fulfill({ body: wav(), contentType: 'audio/wav' }));
  await page.goto('/?fake&nosplash');
  await page.waitForFunction(() => 'fourSouls' in window);
  const live = await page.evaluate(async () => {
    interface Hook {
      sounds: {
        load: (s: Record<string, string[]>) => void;
        play: (s: string) => void;
        live: () => number;
      };
    }
    const { sounds } = (window as unknown as { fourSouls: Hook }).fourSouls;
    // Every one-shot sound at once; a sound starts again only after
    // 80 ms (20.14), so each round waits a little longer than that.
    const names = Array.from({ length: 20 }, (_, i) => `s${String(i)}`);
    sounds.load(Object.fromEntries(names.map((n) => [n, [`/test-sounds/${n}.wav`]])));
    let most = 0;
    for (let round = 0; round < 60; round++) {
      for (const n of names) {
        sounds.play(n);
      }
      most = Math.max(most, sounds.live());
      await new Promise((r) => setTimeout(r, 90));
    }
    return most;
  });
  expect(live).toBeLessThanOrEqual(16);
  expect(errors.filter((e) => e.includes('WebMediaPlayer'))).toEqual([]);
});
