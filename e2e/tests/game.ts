// Starts a real game on the test server with one page per player.
import { expect, type Browser, type Page } from '@playwright/test';
import { appURL, serverPort } from './server';

const address = `127.0.0.1:${String(serverPort)}`;

// startGame: the first name hosts a game with a seat per name, the others
// join. Tests share one server, so each test needs its own host name.
// Everyone gets ready and waits for the table. Each player has its
// own browser context, like its own app.
export async function startGame(
  browser: Browser,
  names: string[],
  {
    viewport = { width: 1280, height: 720 },
    touch = false,
    url = appURL,
    // setup changes the host form before Host is pressed.
    setup = (_: Page) => Promise.resolve(),
  } = {},
): Promise<Page[]> {
  const pages: Page[] = [];
  for (const name of names) {
    void name;
    pages.push(await (await browser.newContext({ viewport, hasTouch: touch })).newPage());
  }
  const [host, ...guests] = pages;
  if (!host) {
    throw new Error('no players');
  }
  await host.goto(url);
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByLabel('Nickname').fill(names[0] ?? '');
  await host.getByLabel('Players').selectOption(String(names.length));
  await setup(host);
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toBeVisible();

  for (const [i, page] of guests.entries()) {
    await page.goto(url);
    await page.getByRole('button', { name: 'Join a game' }).click();
    await page.getByLabel('Nickname').fill(names[i + 1] ?? '');
    await page.getByLabel('Host address').fill(address);
    await page.getByRole('button', { name: 'Connect' }).click();
    await page
      .getByRole('listitem')
      .filter({ hasText: `${names[0] ?? ''}'s game` })
      .getByRole('button', { name: 'Join' })
      .click();
    await expect(page.getByRole('list', { name: 'Seats' })).toBeVisible();
  }
  for (const page of pages) {
    await page.getByRole('button', { name: 'Ready', exact: true }).click();
  }
  for (const page of pages) {
    await expect(page.getByRole('main', { name: /^Turn / })).toBeVisible();
  }
  return pages;
}

// fitsScreen fails when the table scrolls or a part of it leaves the
// window.
export async function expectFits(page: Page): Promise<void> {
  const outside = await page.evaluate(() => {
    const out: string[] = [];
    const { innerWidth: w, innerHeight: h } = window;
    if (document.documentElement.scrollHeight > h || document.documentElement.scrollWidth > w) {
      out.push('page scrolls');
    }
    for (const el of document.querySelectorAll('.player, .board-row, .board-side, .hand')) {
      const r = el.getBoundingClientRect();
      if (r.left < 0 || r.top < 0 || r.right > w || r.bottom > h) {
        out.push(`${el.className} at ${JSON.stringify(r)}`);
      }
    }
    return out;
  });
  expect(outside).toEqual([]);
}
