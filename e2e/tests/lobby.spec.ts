// Host and join screens (7.3), against the real Go server.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { gameId, lobbyGame, startGame, throwFirstRoll } from './game';
import { appURL, serverPort } from './server';
import { realAssets, shot } from './shot';
import { openHost, openLan, openOnline } from './menu';
import { locale, tr, trRe, whole } from './locale';

// address is the project's test server.
const address = () => `127.0.0.1:${String(serverPort())}`;

test('two apps reach the same lobby and start the game', async ({ browser }) => {
  // Separate contexts: two apps with their own settings.
  const ann = await (await browser.newContext()).newPage();
  const bo = await (await browser.newContext()).newPage();

  await realAssets(ann);
  await ann.goto(appURL());
  await openHost(ann);
  await ann.getByLabel(tr('lobby.nickname')).fill('Ann');
  await ann.getByLabel(tr('host.detailed')).check();
  await shot(ann, 'host');
  await ann.getByRole('button', { name: tr('host.host'), exact: true }).click();
  const annSeats = ann.getByRole('list', { name: tr('lobby.seats.seats') });
  await expect(annSeats).toContainText(`Ann ${tr('lobby.seats.you')}`);
  await expect(ann.getByText(address())).toBeVisible();
  const id = await gameId(ann);

  await realAssets(bo);
  await bo.goto(appURL());
  await openOnline(bo);
  await bo.getByLabel(tr('lobby.nickname')).fill('Bo');
  await bo.getByLabel(tr('lobby.online.address')).fill(address());
  await bo.getByRole('button', { name: tr('lobby.connect') }).click();
  const game = lobbyGame(bo, id);
  await expect(game).toContainText(tr('lobby.game_row', { host: 'Ann', taken: 1, n: 2 }));
  await shot(bo, 'games');
  await game.getByRole('button', { name: tr('lobby.join') }).click();
  await expect(bo.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Bo ${tr('lobby.seats.you')}`);
  await expect(annSeats).toContainText('Bo');
  await shot(ann, 'table');

  await ann.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  await expect(bo.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(tr('lobby.seats.is_ready'));
  await bo.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  // The table opens before turn 1, for the first-player roll (15.43).
  await expect(ann.getByRole('main', { name: new RegExp(`^${trRe('game.turn').source}`) })).toBeVisible();
  await throwFirstRoll([ann, bo]);
  await expect(ann.getByRole('main', { name: tr('game.turn', { n: 1 }) })).toBeVisible();
  await expect(bo.getByRole('main', { name: tr('game.turn', { n: 1 }) })).toBeVisible();
});

test('a joined server is remembered by its name, until removed (N-18)', async ({ page }) => {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Cy');
  await page.getByLabel(tr('lobby.online.address')).fill(address());
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  await expect(page.getByText('Test server', { exact: true })).toBeVisible();
  // Back goes one level up: the Online screen.
  await page.getByRole('button', { name: tr('lobby.back') }).click();
  const recent = page.getByRole('list', { name: tr('lobby.online.recent') });
  const server = recent.getByRole('button', { name: 'Test server', exact: true });
  await expect(server).toHaveAttribute('title', address());
  await shot(page, 'join-recent');

  // The name joins the address it stands for.
  await page.getByLabel(tr('lobby.online.address')).fill('');
  await server.click();
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  await page.getByRole('button', { name: tr('lobby.back') }).click();

  await page.getByRole('button', { name: tr('lobby.online.forget_label', { server: 'Test server' }) }).click();
  await expect(recent).toBeHidden();
  await page.reload();
  await openOnline(page);
  await expect(page.getByRole('list', { name: tr('lobby.online.recent') })).toBeHidden();
});

test('a wrong address says why', async ({ page }) => {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Di');
  await page.getByLabel(tr('lobby.online.address')).fill('not an address');
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('error.bad_address'));
  // Nobody listens on the port next to the test server.
  await page.getByLabel(tr('lobby.online.address')).fill(`127.0.0.1:${String(serverPort() - 1)}`);
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('error.no_connection'), { timeout: 8000 });
});

test('a typed ws:// address joins the server (8.18)', async ({ page }) => {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Wes');
  await page.getByLabel(tr('lobby.online.address')).fill(`ws://${address()}/ws`);
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  await expect(page.getByText('Test server', { exact: true })).toBeVisible();
});

// games.test never resolves, so wss:// fails for real; the plain socket
// on 4774 is caught here, so no server is needed (N-19). A caught socket
// fires no "websocket" event: `opened` lists only the real attempts.
test('a bare domain tries wss:// first, then ws:// on 4774 (8.18)', async ({ page }) => {
  const opened: string[] = [];
  page.on('websocket', (ws) => opened.push(ws.url()));
  let hello = '';
  await page.routeWebSocket('ws://games.test:4774/ws', (ws) => {
    ws.onMessage((m) => {
      hello ||= String(m);
    });
  });
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Dom');
  await page.getByLabel(tr('lobby.online.address')).fill('games.test');
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect.poll(() => hello, { timeout: 8000 }).toContain('"type":"hello"');
  expect(opened.filter((u) => u.includes('games.test'))).toEqual(['wss://games.test/ws']);
});

test('a game found on the LAN is joined without typing (N-13)', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  await host.goto(appURL());
  await openHost(host);
  await host.getByLabel(tr('lobby.nickname')).fill('Lex');
  await host.getByRole('button', { name: tr('host.host'), exact: true }).click();
  await expect(host.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Lex ${tr('lobby.seats.you')}`);
  const id = await gameId(host);

  // The fake LAN holds the test server; its real games list follows.
  const guest = await (await browser.newContext()).newPage();
  await realAssets(guest);
  await guest.goto(appURL());
  await openLan(guest);
  await guest.getByLabel(tr('lobby.nickname')).fill('Mio');
  const lan = guest.getByRole('region', { name: tr('lobby.lan.on_network') });
  const found = lan.getByRole('listitem').filter({ hasText: address() });
  await expect(found).toContainText(trRe('lobby.game_row', { taken: 1, n: 2 }));
  await shot(guest, 'join-lan');
  await found.getByRole('button', { name: tr('lobby.join') }).click();
  await lobbyGame(guest, id).getByRole('button', { name: tr('lobby.join') }).click();
  await expect(guest.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Mio ${tr('lobby.seats.you')}`);
});

test('with no game on the LAN, the address field is the way', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await openLan(page);
  const lan = page.getByRole('region', { name: tr('lobby.lan.on_network') });
  await expect(lan).toContainText(tr('lobby.lan.none'));
  await expect(lan.getByLabel(tr('lobby.lan.address'))).toBeVisible();
  // The search runs by itself: there is nothing to press.
  await expect(lan.getByRole('button', { name: 'Search again' })).toHaveCount(0);
});

// lanGame is how the fake LAN lists a game on the test server.
function lanGame(game: { id: string; host: string; taken: number; started: boolean }) {
  return [{ address: address(), app: 'demo', games: [{ seats: 2, locked: false, ...game }] }];
}

// setLan makes page's LAN search find found from now on, also after a
// reload; the fake gives it the app's own protocol.
async function setLan(page: Page, found: ReturnType<typeof lanGame> | []) {
  const run = (f: unknown[]) => {
    (window as unknown as { fakeLan: unknown[] }).fakeLan = f;
  };
  await page.addInitScript(run, found);
  await page.evaluate(run, found).catch(() => undefined);
}

test('the LAN list fills by itself, without a click (15.97)', async ({ page }) => {
  await setLan(page, []);
  await page.goto(appURL());
  await openLan(page);
  const lan = page.getByRole('region', { name: tr('lobby.lan.on_network') });
  await expect(lan).toContainText(tr('lobby.lan.none'));
  // A host opens a game: the next search finds it.
  await page.evaluate((f) => {
    (window as unknown as { fakeLan: unknown[] }).fakeLan = f;
  }, lanGame({ id: 'g-new', host: 'Neo', taken: 1, started: false }));
  await expect(lan.getByRole('listitem').filter({ hasText: tr('lobby.game_row', { host: 'Neo', taken: 1, n: 2 }) })).toBeVisible();
});

// A started game keeps showing, "In progress"; only a player who held a
// seat can rejoin it (15.97, N-13).
test('a player whose app closed rejoins from the LAN list (15.97)', async ({ browser }) => {
  const [ann, bo] = await startGame(browser, ['Ann', 'Bo'], { url: `${appURL()}&anim=off` });
  if (!ann || !bo) {
    throw new Error('no pages');
  }
  const id = await gameId(bo);
  const found = lanGame({ id, host: 'Ann', taken: 2, started: true });
  // Bo's app closes; the same settings start it again.
  const context = bo.context();
  await bo.close();
  const again = await context.newPage();
  await realAssets(again);
  await setLan(again, found);
  await again.goto(`${appURL()}&anim=off`);
  await openLan(again);
  const row = again.getByRole('region', { name: tr('lobby.lan.on_network') }).getByRole('listitem').filter({ hasText: tr('lobby.lan.in_progress') });
  await expect(row).toContainText(`${tr('lobby.game_row', { host: 'Ann', taken: 2, n: 2 })} · ${tr('lobby.lan.in_progress')}`);
  await shot(again, 'lan-rejoin');

  // Somebody else sees it too, but cannot join it.
  const cy = await (await browser.newContext()).newPage();
  await setLan(cy, found);
  await cy.goto(appURL());
  await openLan(cy);
  const theirs = cy.getByRole('listitem').filter({ hasText: tr('lobby.lan.in_progress') });
  await expect(theirs.getByRole('button', { name: tr('lobby.join') })).toBeDisabled();
  await expect(theirs.getByRole('button', { name: tr('lobby.rejoin') })).toHaveCount(0);

  await row.getByRole('button', { name: tr('lobby.rejoin') }).click();
  await expect(again.locator('.player-bottom')).toHaveAccessibleName(tr('table.own_name', { player: 'Bo' }));
  await expect(ann.getByRole('dialog', { name: tr('pause.label') })).toHaveCount(0);
});

test('the LAN screen keeps the typed address, not a Recent list (15.97)', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await openLan(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Tia');
  await page.getByLabel(tr('lobby.lan.address')).fill(address());
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  // Back goes one level up: the LAN screen; Esc then the menu.
  await page.getByRole('button', { name: tr('lobby.back') }).click();
  await expect(page.getByRole('heading', { name: tr('menu.lan') })).toBeVisible();
  await expect(page.getByLabel(tr('lobby.lan.address'))).toHaveValue(address());
  await page.keyboard.press('Escape');
  await expect(page.getByRole('group', { name: tr('menu.play') })).toBeVisible();

  await page.reload();
  await openLan(page);
  await expect(page.getByLabel(tr('lobby.lan.address'))).toHaveValue(address());
  await expect(page.getByRole('list', { name: tr('lobby.online.recent') })).toHaveCount(0);
  await shot(page, 'lan');
});

test('the host setup goes back to the LAN screen (15.97)', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await openHost(page);
  await expect(page.getByRole('heading', { name: tr('host.title_host') })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('heading', { name: tr('menu.lan') })).toBeVisible();
  await page.getByRole('button', { name: tr('host.title_host') }).click();
  await page.getByRole('button', { name: tr('lobby.back') }).click();
  await expect(page.getByRole('heading', { name: tr('menu.lan') })).toBeVisible();
});

test('a game password keeps others out (N-14)', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  await host.goto(appURL());
  await openHost(host);
  await host.getByLabel(tr('lobby.nickname')).fill('Nox');
  await host.getByLabel(tr('host.password')).fill('pw-123');
  await host.getByRole('button', { name: tr('host.host'), exact: true }).click();
  await expect(host.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Nox ${tr('lobby.seats.you')}`);
  const id = await gameId(host);

  const guest = await (await browser.newContext()).newPage();
  await realAssets(guest);
  await guest.goto(appURL());
  await openOnline(guest);
  await guest.getByLabel(tr('lobby.nickname')).fill('Pip');
  await guest.getByLabel(tr('lobby.online.address')).fill(address());
  await guest.getByRole('button', { name: tr('lobby.connect') }).click();
  const game = lobbyGame(guest, id);
  await expect(game.getByLabel(tr('lobby.games.locked'))).toBeVisible();
  const join = game.getByRole('button', { name: tr('lobby.join') });
  const status = guest.getByRole('status');

  await join.click();
  await expect(status).toHaveText(tr('error.password_needed'));
  // The field keeps room to type next to a long set list (19.88).
  const field = await game.getByLabel(tr('lobby.games.password_label', { host: 'Nox' })).boundingBox();
  expect(field?.width).toBeGreaterThanOrEqual(120);
  await game.getByLabel(tr('lobby.games.password_label', { host: 'Nox' })).fill('guess');
  await join.click();
  await expect(status).toHaveText(tr('error.password_wrong'));
  await shot(guest, 'join-password');
  await game.getByLabel(tr('lobby.games.password_label', { host: 'Nox' })).fill('pw-123');
  await join.click();
  await expect(guest.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Pip ${tr('lobby.seats.you')}`);
});

test('the host picks card sets; the choice is remembered', async ({ browser }) => {
  const ctx = await browser.newContext();
  const host = await ctx.newPage();
  const guest = await (await browser.newContext()).newPage();

  await realAssets(host);
  await host.goto(appURL());
  await openHost(host);
  await host.getByLabel(tr('lobby.nickname')).fill('Tess');
  const sets = host.getByRole('group', { name: tr('host.sets') });
  // The core sets are on by default; the Base Game can't be turned off.
  await expect(sets.getByLabel(tr('set.base_game'))).toBeChecked();
  await expect(sets.getByLabel(tr('set.base_game'))).toBeDisabled();
  await expect(sets.getByLabel(tr('set.gold_box'))).toBeChecked();
  await expect(sets.getByLabel(tr('set.four_souls_plus'))).toBeChecked();
  await sets.getByLabel(tr('set.four_souls_plus')).uncheck();
  await host.getByLabel(tr('host.players')).selectOption('3');
  await host.getByLabel(tr('host.detailed')).check();
  await host.getByLabel(tr('host.timer')).selectOption('90');
  await host.getByLabel(tr('host.password')).fill('secret');
  await shot(host, 'host-sets');
  await host.getByRole('button', { name: tr('host.host'), exact: true }).click();
  await expect(host.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Tess ${tr('lobby.seats.you')}`);
  const id = await gameId(host);

  await guest.goto(appURL());
  await openOnline(guest);
  await guest.getByLabel(tr('lobby.nickname')).fill('Lu');
  await guest.getByLabel(tr('lobby.online.address')).fill(address());
  await guest.getByRole('button', { name: tr('lobby.connect') }).click();
  const game = lobbyGame(guest, id);
  await expect(game).toContainText(trRe('lobby.game_row', { host: 'Tess' }));
  await expect(game).toContainText(`${tr('set.base_game')}, ${tr('set.gold_box')}`);
  await expect(game).not.toContainText(tr('set.four_souls_plus'));

  // The next host screen in the same app starts from this setup.
  const again = await ctx.newPage();
  await again.goto(appURL());
  await openHost(again);
  const next = again.getByRole('group', { name: tr('host.sets') });
  await expect(next.getByLabel(tr('set.gold_box'))).toBeChecked();
  await expect(next.getByLabel(tr('set.four_souls_plus'))).not.toBeChecked();
  await expect(again.getByLabel(tr('host.players'))).toHaveValue('3');
  await expect(again.getByLabel(tr('host.detailed'))).toBeChecked();
  await expect(again.getByLabel(tr('host.timer'))).toHaveValue('90');
  await expect(again.getByLabel(tr('host.password'))).toHaveValue('');
});

test('a player creates a game on a dedicated server; another joins (N-17)', async ({ browser }) => {
  // The test server is a dedicated one: it lets anyone create.
  const ivy = await (await browser.newContext()).newPage();
  await realAssets(ivy);
  await ivy.goto(appURL());
  await openOnline(ivy);
  await ivy.getByLabel(tr('lobby.nickname')).fill('Ivy');
  await ivy.getByLabel(tr('lobby.online.address')).fill(address());
  await ivy.getByRole('button', { name: tr('lobby.connect') }).click();
  const create = ivy.getByRole('button', { name: tr('host.title_create') });
  await expect(create).toBeVisible();
  await shot(ivy, 'games-create');
  await create.click();
  await expect(ivy.getByRole('heading', { name: tr('host.title_create') })).toBeVisible();
  await expect(ivy.getByText(tr('host.on_server', { address: address() }))).toBeVisible();
  // The nickname went with hello; nothing is hosted here.
  await expect(ivy.getByLabel(tr('lobby.nickname'))).toHaveCount(0);
  await expect(ivy.getByRole('button', { name: tr('host.continue_saved') })).toHaveCount(0);
  // Back keeps the connection: the Games list again.
  await ivy.getByRole('button', { name: tr('lobby.back') }).click();
  await create.click();
  await ivy.getByLabel(tr('host.players')).selectOption('2');
  await shot(ivy, 'create');
  await ivy.getByRole('button', { name: tr('host.create'), exact: true }).click();
  await expect(ivy.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Ivy ${tr('lobby.seats.you')}`);
  const id = await gameId(ivy);

  const jon = await (await browser.newContext()).newPage();
  await jon.goto(appURL());
  await openOnline(jon);
  await jon.getByLabel(tr('lobby.nickname')).fill('Jon');
  await jon.getByLabel(tr('lobby.online.address')).fill(address());
  await jon.getByRole('button', { name: tr('lobby.connect') }).click();
  await lobbyGame(jon, id).getByRole('button', { name: tr('lobby.join') }).click();
  await expect(jon.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Jon ${tr('lobby.seats.you')}`);
  await expect(ivy.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText('Jon');
});

test('an app that hosts offers no "Create a game" in its Games list (N-17)', async ({ page }) => {
  await page.goto(appURL());
  await openHost(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Kim');
  await page.getByRole('button', { name: tr('host.continue_saved') }).click();
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  await expect(page.getByRole('button', { name: tr('host.title_create') })).toHaveCount(0);
});

test('deck ratios are on by default and locked when the sets are too small (GS-07)', async ({ page }) => {
  await realAssets(page);
  await page.goto(`${appURL()}&requiem`);
  await openHost(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Ratia');
  await page.getByLabel(tr('host.detailed')).check();
  const sets = page.getByRole('group', { name: tr('host.sets') });
  const ratios = page.getByLabel(tr('host.deck_ratios'));
  const label = page.locator('label', { has: ratios });
  // Every set together has the cards for the ratios.
  await expect(ratios).toBeChecked();
  await expect(ratios).toBeEnabled();

  // Without Requiem the groups are short: the box is off and locked,
  // and the tooltip says why.
  await sets.getByLabel(tr('set.requiem')).uncheck();
  await expect(ratios).toBeDisabled();
  await expect(ratios).not.toBeChecked();
  await expect(label).toHaveAttribute('title', whole(trRe('refuse.deck_ratios')));
  await shot(page, 'host-deck-ratios-locked');

  // Requiem without Four Souls+ is enough for 3 players, not for 2.
  await sets.getByLabel(tr('set.requiem')).check();
  await sets.getByLabel(tr('set.four_souls_plus')).uncheck();
  await expect(ratios).toBeDisabled();
  await page.getByLabel(tr('host.players')).selectOption('3');
  await expect(ratios).toBeEnabled();
  await expect(ratios).toBeChecked();
  await expect(label).not.toHaveAttribute('title', /./);

  await page.getByRole('button', { name: tr('host.host'), exact: true }).click();
  await expect(page.getByRole('list', { name: tr('lobby.seats.seats') })).toContainText(`Ratia ${tr('lobby.seats.you')}`);
  await expect(page.getByText(trRe('lobby.options.deck_ratios'))).toBeVisible();
});


// Online is for servers: an address and Recent, no LAN games (15.98).
// Back and Esc go one level up from each screen.
test('Online: a server, its games, and one level back each time (15.98)', async ({ page }) => {
  await realAssets(page);
  await page.goto(appURL());
  await openOnline(page);
  await expect(page.getByRole('heading', { name: tr('menu.online') })).toBeVisible();
  await expect(page.getByRole('region', { name: tr('lobby.lan.on_network') })).toHaveCount(0);
  await page.getByLabel(tr('lobby.nickname')).fill('Ola');
  await page.getByLabel(tr('lobby.online.address')).fill(address());
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  await page.getByRole('button', { name: tr('host.title_create') }).click();
  await expect(page.getByRole('heading', { name: tr('host.title_create') })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('heading', { name: tr('menu.online') })).toBeVisible();
  await expect(page.getByRole('list', { name: tr('lobby.online.recent') })).toContainText('Test server');
  await shot(page, 'online');
  await page.keyboard.press('Escape');
  await expect(page.getByRole('group', { name: tr('menu.play') })).toBeVisible();
});

// Recent rows show at once, then whether each server is up (15.100).
test('Recent shows which servers are up and their open games (15.100)', async ({ page }) => {
  // Nobody listens on the port next to the test server.
  const stopped = `127.0.0.1:${String(serverPort() - 1)}`;
  await page.addInitScript(
    ([up, down, language]) => {
      const key = 'fake-settings';
      if (!localStorage.getItem(key)) {
        const recent = [
          { address: up, name: 'Test server' },
          { address: down, name: 'Gone server' },
        ];
        const rejoin = { address: '', token: '', game: '', via: '' };
        localStorage.setItem(key, JSON.stringify({ nickname: 'Rex', recent, animation: 'off', fullscreen: false, mat: '', dice: 'tear', rejoin, language }));
      }
    },
    [address(), stopped, locale],
  );
  await realAssets(page);
  await page.goto(appURL());
  await openOnline(page);
  const recent = page.getByRole('list', { name: tr('lobby.online.recent') });
  const row = (name: string) => recent.getByRole('listitem').filter({ has: page.getByRole('button', { name, exact: true }) });
  // Other tests share the server, so its count is any number.
  await expect(row('Test server')).toContainText(trRe('lobby.online.open_games'));
  await expect(row('Gone server')).toContainText(tr('lobby.online.down'), { timeout: 8000 });
  await shot(page, 'online-status');

  // Back from its games asks again; the count follows the server.
  await row('Test server').getByRole('button', { name: 'Test server', exact: true }).click();
  await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
  await page.getByRole('button', { name: tr('lobby.back') }).click();
  await expect(row('Test server')).toContainText(trRe('lobby.online.open_games'));
});
