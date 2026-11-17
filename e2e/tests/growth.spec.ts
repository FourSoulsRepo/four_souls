// The long-game check of games.spec.ts (22.9, growth.ts) on a plain
// page: one whose list grows each step fails it, one that stays flat
// passes. Chromium only: the script time comes from DevTools.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { describe, limits, measured, openMetrics, overLimits, sample, sampleEvery, type Sample } from './growth';

// steps: enough for the early window and a late one apart from it.
const steps = 1300;

// run plays steps on a page that starts with `base` list items, adds
// `add` a step and walks the whole list each step, as the history strip
// did (19.138).
async function run(page: Page, base: number, add: number): Promise<Sample[]> {
  await page.setContent('<ul id="list"></ul>');
  await page.evaluate(([b, n]) => {
    const list = document.getElementById('list') as HTMLElement;
    for (let i = 0; i < b; i++) {
      list.append(document.createElement('li'));
    }
    (window as unknown as { step: () => number }).step = () => {
      for (let i = 0; i < n; i++) {
        list.append(document.createElement('li'));
      }
      let sum = 0;
      for (const li of list.children) {
        for (let k = 0; k < 200; k++) {
          sum += li.tagName.charCodeAt(k % 2);
        }
      }
      return sum;
    };
  }, [base, add]);
  const cdp = await page.context().newCDPSession(page);
  await openMetrics(cdp);
  const samples = [await sample(page, cdp, 0)];
  for (let step = 1; step <= steps; step++) {
    // In a page task, as an update is: DevTools counts no script that
    // a test's evaluate runs itself.
    await page.evaluate(
      () =>
        new Promise<void>((done) => {
          setTimeout(() => {
            (window as unknown as { step: () => number }).step();
            done();
          }, 0);
        }),
    );
    if (step % sampleEvery === 0 || step === steps) {
      samples.push(await sample(page, cdp, step));
    }
  }
  return samples;
}

test.describe('long-game check', () => {
  test.skip(({ browserName }) => browserName !== 'chromium', 'script time comes from Chromium DevTools');

  test('a page whose list grows each step fails the check', async ({ page }) => {
    const m = measured('Grow', await run(page, 0, 4));
    console.log(describe(m));
    const problems = overLimits(m, 'test game, seed 1');
    expect(problems).toContainEqual(expect.stringMatching(/^Grow: elements at the end \d+, limit \d+; test game, seed 1$/));
    expect(problems).toContainEqual(expect.stringMatching(/^Grow: script growth late \/ early [\d.]+, limit [\d.]+; test game, seed 1$/));
  });

  test('a page that stays flat passes the check', async ({ page }) => {
    const m = measured('Flat', await run(page, 1000, 0));
    console.log(describe(m));
    expect(overLimits(m, 'test game, seed 1')).toEqual([]);
  });
});

test('each limit names the page, the value and the game', () => {
  const m = { page: 'Ann', steps: 3000, elements: limits.elements + 1, sounds: limits.sounds + 1, early: 10, late: limits.scriptMs + 5 };
  expect(overLimits(m, 'g, seed 7')).toEqual([
    `Ann: elements at the end ${String(limits.elements + 1)}, limit ${String(limits.elements)}; g, seed 7`,
    `Ann: live sounds ${String(limits.sounds + 1)}, limit ${String(limits.sounds)}; g, seed 7`,
    `Ann: script ms a step late ${(limits.scriptMs + 5).toFixed(1)}, limit ${String(limits.scriptMs)}; g, seed 7`,
    `Ann: script growth late / early ${((limits.scriptMs + 5) / 10).toFixed(2)}, limit ${String(limits.growth)}; g, seed 7`,
  ]);
});
