// Starts a real game on the test server with one page per player.
import { expect, type Browser, type Locator, type Page } from '@playwright/test';
import { appURL, serverPort } from './server';
import { realAssets } from './shot';
import { enabled, seen, throwButton } from './turn';
import { openHost, openOnline } from './menu';
import { tr, trRe } from './locale';

// statMonsters change other monsters' stats (14.10). A test that checks
// a fresh monster bans them, or a lucky deal puts a badge on it.
export const statMonsters = ['delirium', 'stoney'];

// hostBans adds the cards to the host's create message on its way to
// the server: the host screen has no ban list yet (14.9). It adds the
// other options too: the rollback mode has no control yet (17.4). The
// getter it resolves to returns the bans the server's table names, as
// a check.
function hostBans(page: Page, bans: string[], port: number, options: Record<string, unknown> = {}): Promise<() => string[]> {
  let echoed: string[] = [];
  return page
    .routeWebSocket(`ws://127.0.0.1:${String(port)}/ws`, (ws) => {
      const server = ws.connectToServer();
      ws.onMessage((msg) => {
        const env = JSON.parse(String(msg)) as { type: string; data?: { options?: Record<string, unknown> } };
        if (env.type === 'create' && env.data) {
          env.data.options = { ...env.data.options, ...options, ...(bans.length > 0 ? { host_bans: bans } : {}) };
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

// gameId is the id of the game page sits at, as the server named it.
export async function gameId(page: Page): Promise<string> {
  const id = await page.evaluate(() => {
    interface S {
      lobby: { table?: { game: string } };
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().lobby.table?.game;
  });
  if (!id) {
    throw new Error('not at a table');
  }
  return id;
}

// lobbyGame is the game's row in page's Games list. Tests share one
// server and its list, so a row is found by its game id: host names
// may repeat (19.31).
export function lobbyGame(page: Page, id: string): Locator {
  return page.locator(`li[data-game="${id}"]`);
}

// startGame: the first name hosts a game with a seat per name, the others
// join it by its id, so names may repeat across tests (19.31).
// Everyone gets ready and waits for the table. Each player has its
// own browser context, like its own app.
export async function startGame(
  browser: Browser,
  names: string[],
  {
    viewport = { width: 1280, height: 720 },
    touch = false,
    // port is the server's; the project's test server by default.
    port = serverPort(),
    url = appURL(port),
    // setup changes the host form before Host is pressed.
    setup = (_: Page) => Promise.resolve(),
    // open runs on each new page before it loads the app.
    open = (_: Page) => Promise.resolve(),
    // bans are cards the host leaves out of the game (GS-13).
    bans = [] as string[],
    // options are more create options, e.g. the rollback mode (ADR 020).
    options = {} as Record<string, unknown>,
    // roll throws every first-player die (15.43), so the test starts in
    // turn 1; false leaves the roll to the test.
    roll = true,
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
  const banned = bans.length > 0 || Object.keys(options).length > 0 ? await hostBans(host, bans, port, options) : () => [];
  await host.goto(url);
  await openHost(host);
  await host.getByLabel(tr('lobby.nickname')).fill(names[0] ?? '');
  await host.getByLabel(tr('host.players')).selectOption(String(names.length));
  // Game tests play the Base Game only: they were written for its cards,
  // and a forced attack (Krampus) could keep a turn from ending.
  const sets = host.getByRole('group', { name: tr('host.sets') });
  for (const set of [tr('set.gold_box'), tr('set.four_souls_plus')]) {
    await sets.getByLabel(set).uncheck();
  }
  await setup(host);
  await host.getByRole('button', { name: tr('host.host'), exact: true }).click();
  await expect(host.getByRole('list', { name: tr('lobby.seats.seats') })).toBeVisible();
  expect(banned()).toEqual(bans);
  const id = await gameId(host);

  for (const [i, page] of guests.entries()) {
    await page.goto(url);
    await openOnline(page);
    await page.getByLabel(tr('lobby.nickname')).fill(names[i + 1] ?? '');
    await page.getByLabel(tr('lobby.online.address')).fill(`127.0.0.1:${String(port)}`);
    await page.getByRole('button', { name: tr('lobby.connect') }).click();
    // The click goes to the button, not to a point: a table above that
    // starts moves the row up under the pointer, and a click there hit
    // a started game's disabled Join (19.14).
    const join = lobbyGame(page, id).getByRole('button', { name: tr('lobby.join') });
    await expect(join).toBeEnabled();
    await join.dispatchEvent('click');
    await expect(page.getByRole('list', { name: tr('lobby.seats.seats') })).toBeVisible();
  }
  for (const page of pages) {
    await page.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  }
  for (const page of pages) {
    await expect(page.getByRole('main', { name: new RegExp(`^${trRe('game.turn').source}`) })).toBeVisible();
  }
  if (roll) {
    await throwFirstRoll(pages);
  }
  return pages;
}

// idle: the page plays no effect and shows the latest update.
function idle(): boolean {
  interface S {
    playing: unknown;
    behind: boolean;
  }
  const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
  return s.playing === null && !s.behind;
}

// throwFirstRoll throws every die of the first-player roll, ties too,
// until every page shows turn 1 (15.43). A Cain deal has no roll. It
// only throws: the game waits in turn 1 as it started.
export async function throwFirstRoll(pages: Page[]): Promise<void> {
  for (let i = 0; i < 600; i++) {
    const states = await Promise.all(pages.map(seen));
    if (states.every((s) => s.turn >= 1)) {
      // Cain's banner plays on the drawn table: wait for it too.
      await Promise.all(pages.map((page) => page.waitForFunction(idle)));
      return;
    }
    const roller = pages.find((_, j) => states[j]?.ready && states[j].turn < 1 && states[j].waiting === states[j].you);
    const button = roller ? throwButton(roller) : null;
    // Eden's starting item comes before the roll (R-SETUP-09): pick the first.
    const eden = roller ? roller.locator('.choice:not(.look) .choice-option button').first() : null;
    if (roller && button && (await enabled(button))) {
      await button.click();
      // Off the cards, so no zoom covers the table.
      await roller.mouse.move(1, 1);
    } else if (roller && eden && (await enabled(eden))) {
      await eden.click();
    } else {
      await pages[0]?.waitForTimeout(100);
    }
  }
  throw new Error('the first-player roll never ended');
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
