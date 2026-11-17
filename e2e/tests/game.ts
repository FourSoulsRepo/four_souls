// Starts a real game on the test server with one page per player.
import { expect, test, type Browser, type Locator, type Page } from '@playwright/test';
import { appURL, serverPort } from './server';
import { realAssets } from './shot';
import { enabled, seen, throwButton } from './turn';
import { openHost, openOnline } from './menu';
import { tr, trRe } from './locale';

// repeatNames makes the names unique to a --repeat-each run of the
// test. A saved game is found by its players' names, and repeats of one
// test run at the same time (19.152).
export function repeatNames(names: string[]): string[] {
  const r = test.info().repeatEachIndex;
  return r === 0 ? names : names.map((n) => `${n}${String(r)}`);
}

// statMonsters change other monsters' stats (14.10). A test that checks
// a fresh monster bans them, or a lucky deal puts a badge on it.
export const statMonsters = ['delirium', 'stoney'];

// createOptions adds options to every create message the page sends:
// the host screen has no ban list (14.9), rollback control (17.4) or
// seed. It wraps the page's own WebSocket send, so the socket is not
// routed through the test runner: every message keeps its timing.
function createOptions(page: Page, options: Record<string, unknown>): Promise<void> {
  return page.addInitScript((options) => {
    const send = WebSocket.prototype.send;
    WebSocket.prototype.send = function (this: WebSocket, data: Parameters<WebSocket['send']>[0]) {
      if (typeof data === 'string' && data.startsWith('{"type":"create"')) {
        const env = JSON.parse(data) as { data?: { options?: Record<string, unknown> } };
        if (env.data) {
          env.data.options = { ...env.data.options, ...options };
          data = JSON.stringify(env);
        }
      }
      send.call(this, data);
    };
  }, options);
}

// tableOptions are the options of the table the page sits at, as the
// server names them.
function tableOptions(page: Page): Promise<{ host_bans?: string[]; seed?: number }> {
  return page.evaluate(() => {
    interface S {
      lobby: { table?: { options?: { host_bans?: string[]; seed?: number } } };
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().lobby.table?.options ?? {};
  });
}

// testSeed is the test's own deal (ADR 032): a hash (FNV-1a) of its
// title path, so it is the same on every run and every repeat, and in
// both browsers. Never 0, which asks the server for a clock seed.
export function testSeed(): number {
  let h = 0x811c9dc5;
  for (const c of new TextEncoder().encode(test.info().titlePath.join(' › '))) {
    h = Math.imul(h ^ c, 0x01000193) >>> 0;
  }
  return h === 0 ? 1 : h;
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

// tableSeat is one row of the table screen's seat list.
export function tableSeat(page: Page, i: number): Locator {
  return page.getByRole('list', { name: tr('lobby.seats.seats') }).locator(`li[data-seat="${String(i)}"]`);
}

// addBot presses a free seat's Add bot button for the level, as the
// host (23.11, ADR 033 Decision 10); level is its translated name.
export async function addBot(page: Page, i: number, level: string): Promise<void> {
  await tableSeat(page, i).getByRole('button', { name: tr('lobby.bots.add_level', { level }) }).click();
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
    // seed fixes the deal (ADR 032): the test's own by default, the
    // same on every run; 0 asks for a random deal (a clock seed).
    seed = testSeed(),
    // roll throws every first-player die (15.43), so the test starts in
    // turn 1; false leaves the roll to the test.
    roll = true,
    // pick plays a draft's pick screen (14.11) once everyone is ready;
    // the game starts when it returns.
    pick = undefined as ((pages: Page[]) => Promise<void>) | undefined,
    // listed runs on each guest page once the game's row is in its list,
    // before Join: a test checks what the list shows.
    listed = (_: Locator) => Promise.resolve(),
    // seated runs on each guest page at the table, before Ready.
    seated = (_: Page) => Promise.resolve(),
    // bots are the levels of bots the host seats after the players
    // (23.12): 'standard' or 'hard', one more seat each.
    bots = [] as string[],
    // animations keeps each page's own speed: its saved setting, or the
    // url's anim flag. By default every page plays at "off" from Ready
    // on (TS-10, 22.6): a test that checks a motion or sound sets it.
    animations = false,
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
  // The seed shows when the test starts (CLAUDE.md 5.19.3).
  const seedText = seed === 0 ? 'random' : String(seed);
  console.log(`${test.info().titlePath.join(' › ')}: seed ${seedText}`);
  test.info().annotations.push({ type: 'seed', description: seedText });
  const create = { ...options, ...(bans.length > 0 ? { host_bans: bans } : {}), ...(seed === 0 ? {} : { seed }) };
  if (Object.keys(create).length > 0) {
    await createOptions(host, create);
  }
  await host.goto(url);
  await openHost(host);
  await host.getByLabel(tr('lobby.nickname')).fill(names[0] ?? '');
  await host.getByLabel(tr('host.players')).selectOption(String(names.length + bots.length));
  // Game tests play the Base Game only: they were written for its cards,
  // and a forced attack (Krampus) could keep a turn from ending.
  const sets = host.getByRole('group', { name: tr('host.sets') });
  for (const set of [tr('set.gold_box'), tr('set.four_souls_plus')]) {
    await sets.getByLabel(set).uncheck();
  }
  await setup(host);
  await host.getByRole('button', { name: tr('host.host'), exact: true }).click();
  await expect(host.getByRole('list', { name: tr('lobby.seats.seats') })).toBeVisible();
  const echoed = await tableOptions(host);
  expect(echoed.host_bans ?? []).toEqual(bans);
  expect(echoed.seed ?? 0).toBe(seed);
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
    await listed(lobbyGame(page, id));
    await join.dispatchEvent('click');
    await expect(page.getByRole('list', { name: tr('lobby.seats.seats') })).toBeVisible();
    await seated(page);
  }
  // The bots take the seats after the players, in order.
  for (const [i, level] of bots.entries()) {
    const at = names.length + i;
    await addBot(host, at, tr(`bot.level.${level}`));
    await expect(tableSeat(host, at)).toContainText(tr('bot.number', { n: i + 1 }));
  }
  if (!animations) {
    await Promise.all(pages.map((page) => setSpeed(page, 'off')));
  }
  for (const page of pages) {
    await page.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  }
  if (pick) {
    await pick(pages);
  }
  for (const page of pages) {
    await expect(page.getByRole('main', { name: new RegExp(`^${trRe('game.turn').source}`) })).toBeVisible();
  }
  if (roll) {
    // The roll is setup here; first-roll.spec.ts checks how it plays,
    // with roll: false. Animated, each throw waits for the dice before
    // it, and a tie throws again: 6 s, 12 s after a tie, 17 to 25 s with
    // ties under a full run, which used up a 30 s test (19.120). So it
    // plays with animations off; each page's speed comes back after:
    // "off" too, unless the test asked for animations.
    const speeds = await Promise.all(pages.map((page) => setSpeed(page, 'off')));
    await throwFirstRoll(pages);
    await Promise.all(pages.map((page, i) => setSpeed(page, speeds[i] ?? 'normal')));
  }
  return pages;
}

// setSpeed sets the page's animation speed for this session only, not
// in its saved settings, and returns the speed before.
export async function setSpeed(page: Page, speed: string): Promise<string> {
  return page.evaluate((speed) => {
    interface S {
      settings: { animation: string };
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } }).fourSouls.store;
    const { settings } = store.getState();
    store.setState({ settings: { ...settings, animation: speed } });
    return settings.animation;
  }, speed);
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
    // Eden's starting item comes before the roll (R-SETUP-09): pick the
    // first, then Confirm (20.22).
    const eden = roller ? roller.locator('.choice:not(.look) .choice-option button').first() : null;
    const confirm = roller ? roller.locator('.choice .choice-confirm') : null;
    if (roller && button && (await enabled(button))) {
      await button.click();
      // Off the cards, so no zoom covers the table.
      await roller.mouse.move(1, 1);
    } else if (roller && confirm && (await enabled(confirm))) {
      await confirm.click();
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
