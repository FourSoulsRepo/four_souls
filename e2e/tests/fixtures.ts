// Every spec imports test and expect from here, not from
// @playwright/test (19.6, ADR 012). After each test, the automatic
// fixture fails it if any app page logged an animation clash: two
// animations brought or showed one card (frontend/src/animations.ts).
import { test as base, expect, type Page } from '@playwright/test';

export { expect };

// checked are the pages already looked at: game pages stay open after
// their test, and one clash must fail only the test that made it.
const checked = new WeakSet<Page>();

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
      for (const page of [...browser.contexts(), context].flatMap((c) => c.pages())) {
        if (checked.has(page)) {
          continue;
        }
        checked.add(page);
        expect(await clashesOf(page), `animation clashes on ${page.url()}`).toEqual([]);
      }
    },
    { auto: true },
  ],
});
