// Starts a real game on the test server with one page per player.
import { expect, type Browser, type Page } from '@playwright/test';
import { appURL, serverPort } from './server';
import { realAssets } from './shot';

// statMonsters change other monsters' stats (14.10). A test that checks
// a fresh monster bans them, or a lucky deal puts a badge on it.
export const statMonsters = ['delirium', 'stoney'];

// hostBans adds the cards to the host's create message on its way to
// the server: the host screen has no ban list yet (14.9). The getter it
// resolves to returns the bans the server's table names, as a check.
function hostBans(page: Page, bans: string[]): Promise<() => string[]> {
  let echoed: string[] = [];
  return page
    .routeWebSocket(`ws://127.0.0.1:${String(serverPort())}/ws`, (ws) => {
      const server = ws.connectToServer();
      ws.onMessage((msg) => {
        const env = JSON.parse(String(msg)) as { type: string; data?: { options?: Record<string, unknown> } };
        if (env.type === 'create' && env.data) {
          env.data.options = { ...env.data.options, host_bans: bans };
        }
        server.send(JSON.stringify(env));
      });
      server.onMessage((msg) => {
        const env = JSON.parse(String(msg)) as { type: string; data?: { options?: { host_bans?: string[] } } };
        if (env.type === 'table') {
          echoed = env.data?.options?.host_bans ?? [];
        }
        ws.send(msg);
      });
    })
    .then(() => () => echoed);
}

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
    url = appURL(),
    // setup changes the host form before Host is pressed.
    setup = (_: Page) => Promise.resolve(),
    // open runs on each new page before it loads the app.
    open = (_: Page) => Promise.resolve(),
    // bans are cards the host leaves out of the game (GS-13).
    bans = [] as string[],
  } = {},
): Promise<Page[]> {
  const pages: Page[] = [];
  for (const name of names) {
    void name;
    const page = await (await browser.newContext({ viewport, hasTouch: touch })).newPage();
    await realAssets(page);
    await open(page);
    pages.push(page);
  }
  const [host, ...guests] = pages;
  if (!host) {
    throw new Error('no players');
  }
  const banned = bans.length > 0 ? await hostBans(host, bans) : () => [];
  await host.goto(url);
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByLabel('Nickname').fill(names[0] ?? '');
  await host.getByLabel('Players').selectOption(String(names.length));
  // Game tests play the Base Game only: they were written for its cards,
  // and a forced attack (Krampus) could keep a turn from ending.
  const sets = host.getByRole('group', { name: 'Card sets' });
  for (const set of ['Gold Box', 'Four Souls+']) {
    await sets.getByLabel(set).uncheck();
  }
  await setup(host);
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toBeVisible();
  expect(banned()).toEqual(bans);

  for (const [i, page] of guests.entries()) {
    await page.goto(url);
    await page.getByRole('button', { name: 'Join a game' }).click();
    await page.getByLabel('Nickname').fill(names[i + 1] ?? '');
    await page.getByLabel('Host address').fill(`127.0.0.1:${String(serverPort())}`);
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
    for (const el of document.querySelectorAll('.player, .board, .hand')) {
      const r = el.getBoundingClientRect();
      if (r.left < 0 || r.top < 0 || r.right > w || r.bottom > h) {
        out.push(`${el.className} at ${JSON.stringify(r)}`);
      }
    }
    return out;
  });
  expect(outside).toEqual([]);
}
