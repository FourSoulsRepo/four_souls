// Every spec imports test and expect from here, not from
// @playwright/test (19.6, ADR 012). After each test, the automatic
// fixture fails it if any app page logged an animation clash: two
// animations brought or showed one card (frontend/src/animations.ts).
// Then it closes the contexts the test opened (startGame's players),
// so no game page outlives its test (19.38).
import { test as base, expect, type Page } from '@playwright/test';

export { expect };

async function clashesOf(page: Page): Promise<unknown[]> {
  try {
    return await page.evaluate(() => {
      const w = window as unknown as { fourSouls?: { clashes?: unknown[] } };
      return w.fourSouls?.clashes ?? [];
    });
  } catch {
    return []; // a closed or navigating page has nothing to check
  }
}

export const test = base.extend<{ clashCheck: undefined }>({
  clashCheck: [
    async ({ browser, context }, use) => {
      await use(undefined);
      // Playwright closes the test's own context; the rest are ours.
      const others = browser.contexts().filter((c) => c !== context);
      try {
        for (const page of [context, ...others].flatMap((c) => c.pages())) {
          expect(await clashesOf(page), `animation clashes on ${page.url()}`).toEqual([]);
        }
      } finally {
        for (const c of others) {
          await c.close();
        }
      }
    },
    { auto: true },
  ],
});
