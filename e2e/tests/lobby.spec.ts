// Host and join screens (7.3), against the real Go server.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { gameId, lobbyGame, startGame, throwFirstRoll } from './game';
import { appURL, serverPort } from './server';
import { realAssets, shot } from './shot';
import { openHost, openLan, openOnline } from './menu';

// address is the project's test server.
const address = () => `127.0.0.1:${String(serverPort())}`;

test('two apps reach the same lobby and start the game', async ({ browser }) => {
  // Separate contexts: two apps with their own settings.
  const ann = await (await browser.newContext()).newPage();
  const bo = await (await browser.newContext()).newPage();

  await realAssets(ann);
  await ann.goto(appURL());
  await openHost(ann);
  await ann.getByLabel('Nickname').fill('Ann');
  await ann.getByLabel('Detailed setup').check();
  await shot(ann, 'host');
  await ann.getByRole('button', { name: 'Host', exact: true }).click();
  const annSeats = ann.getByRole('list', { name: 'Seats' });
  await expect(annSeats).toContainText('Ann (you)');
  await expect(ann.getByText(address())).toBeVisible();
  const id = await gameId(ann);

  await realAssets(bo);
  await bo.goto(appURL());
  await openOnline(bo);
  await bo.getByLabel('Nickname').fill('Bo');
  await bo.getByLabel('Server address').fill(address());
  await bo.getByRole('button', { name: 'Connect' }).click();
  const game = lobbyGame(bo, id);
  await expect(game).toContainText("Ann's game · 1/2 players");
  await shot(bo, 'games');
  await game.getByRole('button', { name: 'Join' }).click();
  await expect(bo.getByRole('list', { name: 'Seats' })).toContainText('Bo (you)');
  await expect(annSeats).toContainText('Bo');
  await shot(ann, 'table');

  await ann.getByRole('button', { name: 'Ready', exact: true }).click();
  await expect(bo.getByRole('list', { name: 'Seats' })).toContainText('Ready');
  await bo.getByRole('button', { name: 'Ready', exact: true }).click();
  // The table opens before turn 1, for the first-player roll (15.43).
  await expect(ann.getByRole('main', { name: /^Turn / })).toBeVisible();
  await throwFirstRoll([ann, bo]);
  await expect(ann.getByRole('main', { name: 'Turn 1' })).toBeVisible();
  await expect(bo.getByRole('main', { name: 'Turn 1' })).toBeVisible();
});

test('a joined server is remembered by its name, until removed (N-18)', async ({ page }) => {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel('Nickname').fill('Cy');
  await page.getByLabel('Server address').fill(address());
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await expect(page.getByText('Test server', { exact: true })).toBeVisible();
  // Back goes one level up: the Online screen.
  await page.getByRole('button', { name: 'Back' }).click();
  const recent = page.getByRole('list', { name: 'Recent' });
  const server = recent.getByRole('button', { name: 'Test server', exact: true });
  await expect(server).toHaveAttribute('title', address());
  await shot(page, 'join-recent');

  // The name joins the address it stands for.
  await page.getByLabel('Server address').fill('');
  await server.click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await page.getByRole('button', { name: 'Back' }).click();

  await page.getByRole('button', { name: 'Remove Test server from Recent' }).click();
  await expect(recent).toBeHidden();
  await page.reload();
  await openOnline(page);
  await expect(page.getByRole('list', { name: 'Recent' })).toBeHidden();
});

test('a wrong address says why', async ({ page }) => {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel('Nickname').fill('Di');
  await page.getByLabel('Server address').fill('not an address');
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('status')).toHaveText('Type an address like 192.168.1.5, 192.168.1.5:4774 or games.example.com.');
  // Nobody listens on the port next to the test server.
  await page.getByLabel('Server address').fill(`127.0.0.1:${String(serverPort() - 1)}`);
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('status')).toHaveText('Could not connect to the server.', { timeout: 8000 });
});

test('a typed ws:// address joins the server (8.18)', async ({ page }) => {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel('Nickname').fill('Wes');
  await page.getByLabel('Server address').fill(`ws://${address()}/ws`);
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
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
  await page.getByLabel('Nickname').fill('Dom');
  await page.getByLabel('Server address').fill('games.test');
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect.poll(() => hello, { timeout: 8000 }).toContain('"type":"hello"');
  expect(opened.filter((u) => u.includes('games.test'))).toEqual(['wss://games.test/ws']);
});

test('a game found on the LAN is joined without typing (N-13)', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  await host.goto(appURL());
  await openHost(host);
  await host.getByLabel('Nickname').fill('Lex');
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toContainText('Lex (you)');
  const id = await gameId(host);

  // The fake LAN holds the test server; its real games list follows.
  const guest = await (await browser.newContext()).newPage();
  await realAssets(guest);
  await guest.goto(appURL());
  await openLan(guest);
  await guest.getByLabel('Nickname').fill('Mio');
  const lan = guest.getByRole('region', { name: 'On your network' });
  const found = lan.getByRole('listitem').filter({ hasText: address() });
  await expect(found).toContainText('1/2 players');
  await shot(guest, 'join-lan');
  await found.getByRole('button', { name: 'Join' }).click();
  await lobbyGame(guest, id).getByRole('button', { name: 'Join' }).click();
  await expect(guest.getByRole('list', { name: 'Seats' })).toContainText('Mio (you)');
});

test('with no game on the LAN, the address field is the way', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await openLan(page);
  const lan = page.getByRole('region', { name: 'On your network' });
  await expect(lan).toContainText('No games found yet. Still looking…');
  await expect(lan.getByLabel('Host address')).toBeVisible();
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
  const lan = page.getByRole('region', { name: 'On your network' });
  await expect(lan).toContainText('No games found yet');
  // A host opens a game: the next search finds it.
  await page.evaluate((f) => {
    (window as unknown as { fakeLan: unknown[] }).fakeLan = f;
  }, lanGame({ id: 'g-new', host: 'Neo', taken: 1, started: false }));
  await expect(lan.getByRole('listitem').filter({ hasText: "Neo's game · 1/2 players" })).toBeVisible();
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
  const row = again.getByRole('region', { name: 'On your network' }).getByRole('listitem').filter({ hasText: 'In progress' });
  await expect(row).toContainText("Ann's game · 2/2 players · In progress");
  await shot(again, 'lan-rejoin');

  // Somebody else sees it too, but cannot join it.
  const cy = await (await browser.newContext()).newPage();
  await setLan(cy, found);
  await cy.goto(appURL());
  await openLan(cy);
  const theirs = cy.getByRole('listitem').filter({ hasText: 'In progress' });
  await expect(theirs.getByRole('button', { name: 'Join' })).toBeDisabled();
  await expect(theirs.getByRole('button', { name: 'Rejoin' })).toHaveCount(0);

  await row.getByRole('button', { name: 'Rejoin' }).click();
  await expect(again.locator('.player-bottom')).toHaveAccessibleName('Bo (you)');
  await expect(ann.getByRole('dialog', { name: 'Paused' })).toHaveCount(0);
});

test('the LAN screen keeps the typed address, not a Recent list (15.97)', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await openLan(page);
  await page.getByLabel('Nickname').fill('Tia');
  await page.getByLabel('Host address').fill(address());
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  // Back goes one level up: the LAN screen; Esc then the menu.
  await page.getByRole('button', { name: 'Back' }).click();
  await expect(page.getByRole('heading', { name: 'LAN' })).toBeVisible();
  await expect(page.getByLabel('Host address')).toHaveValue(address());
  await page.keyboard.press('Escape');
  await expect(page.getByRole('group', { name: 'Play' })).toBeVisible();

  await page.reload();
  await openLan(page);
  await expect(page.getByLabel('Host address')).toHaveValue(address());
  await expect(page.getByRole('list', { name: 'Recent' })).toHaveCount(0);
  await shot(page, 'lan');
});

test('the host setup goes back to the LAN screen (15.97)', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await openHost(page);
  await expect(page.getByRole('heading', { name: 'Host a game' })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('heading', { name: 'LAN' })).toBeVisible();
  await page.getByRole('button', { name: 'Host a game' }).click();
  await page.getByRole('button', { name: 'Back' }).click();
  await expect(page.getByRole('heading', { name: 'LAN' })).toBeVisible();
});

test('a game password keeps others out (N-14)', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  await host.goto(appURL());
  await openHost(host);
  await host.getByLabel('Nickname').fill('Nox');
  await host.getByLabel('Password (optional)').fill('pw-123');
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toContainText('Nox (you)');
  const id = await gameId(host);

  const guest = await (await browser.newContext()).newPage();
  await realAssets(guest);
  await guest.goto(appURL());
  await openOnline(guest);
  await guest.getByLabel('Nickname').fill('Pip');
  await guest.getByLabel('Server address').fill(address());
  await guest.getByRole('button', { name: 'Connect' }).click();
  const game = lobbyGame(guest, id);
  await expect(game.getByLabel('Locked')).toBeVisible();
  const join = game.getByRole('button', { name: 'Join' });
  const status = guest.getByRole('status');

  await join.click();
  await expect(status).toHaveText('this game needs a password');
  await game.getByLabel("Password for Nox's game").fill('guess');
  await join.click();
  await expect(status).toHaveText('wrong password');
  await shot(guest, 'join-password');
  await game.getByLabel("Password for Nox's game").fill('pw-123');
  await join.click();
  await expect(guest.getByRole('list', { name: 'Seats' })).toContainText('Pip (you)');
});

test('the host picks card sets; the choice is remembered', async ({ browser }) => {
  const ctx = await browser.newContext();
  const host = await ctx.newPage();
  const guest = await (await browser.newContext()).newPage();

  await realAssets(host);
  await host.goto(appURL());
  await openHost(host);
  await host.getByLabel('Nickname').fill('Tess');
  const sets = host.getByRole('group', { name: 'Card sets' });
  // The core sets are on by default; the Base Game can't be turned off.
  await expect(sets.getByLabel('Base Game')).toBeChecked();
  await expect(sets.getByLabel('Base Game')).toBeDisabled();
  await expect(sets.getByLabel('Gold Box')).toBeChecked();
  await expect(sets.getByLabel('Four Souls+')).toBeChecked();
  await sets.getByLabel('Four Souls+').uncheck();
  await host.getByLabel('Players').selectOption('3');
  await host.getByLabel('Detailed setup').check();
  await host.getByLabel('Response timer').selectOption('90');
  await host.getByLabel('Password (optional)').fill('secret');
  await shot(host, 'host-sets');
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toContainText('Tess (you)');
  const id = await gameId(host);

  await guest.goto(appURL());
  await openOnline(guest);
  await guest.getByLabel('Nickname').fill('Lu');
  await guest.getByLabel('Server address').fill(address());
  await guest.getByRole('button', { name: 'Connect' }).click();
  const game = lobbyGame(guest, id);
  await expect(game).toContainText("Tess's game");
  await expect(game).toContainText('Base Game, Gold Box');
  await expect(game).not.toContainText('Four Souls+');

  // The next host screen in the same app starts from this setup.
  const again = await ctx.newPage();
  await again.goto(appURL());
  await openHost(again);
  const next = again.getByRole('group', { name: 'Card sets' });
  await expect(next.getByLabel('Gold Box')).toBeChecked();
  await expect(next.getByLabel('Four Souls+')).not.toBeChecked();
  await expect(again.getByLabel('Players')).toHaveValue('3');
  await expect(again.getByLabel('Detailed setup')).toBeChecked();
  await expect(again.getByLabel('Response timer')).toHaveValue('90');
  await expect(again.getByLabel('Password (optional)')).toHaveValue('');
});

test('a player creates a game on a dedicated server; another joins (N-17)', async ({ browser }) => {
  // The test server is a dedicated one: it lets anyone create.
  const ivy = await (await browser.newContext()).newPage();
  await realAssets(ivy);
  await ivy.goto(appURL());
  await openOnline(ivy);
  await ivy.getByLabel('Nickname').fill('Ivy');
  await ivy.getByLabel('Server address').fill(address());
  await ivy.getByRole('button', { name: 'Connect' }).click();
  const create = ivy.getByRole('button', { name: 'Create a game' });
  await expect(create).toBeVisible();
  await shot(ivy, 'games-create');
  await create.click();
  await expect(ivy.getByRole('heading', { name: 'Create a game' })).toBeVisible();
  await expect(ivy.getByText(`On the server ${address()}`)).toBeVisible();
  // The nickname went with hello; nothing is hosted here.
  await expect(ivy.getByLabel('Nickname')).toHaveCount(0);
  await expect(ivy.getByRole('button', { name: 'Continue a saved game' })).toHaveCount(0);
  // Back keeps the connection: the Games list again.
  await ivy.getByRole('button', { name: 'Back' }).click();
  await create.click();
  await ivy.getByLabel('Players').selectOption('2');
  await shot(ivy, 'create');
  await ivy.getByRole('button', { name: 'Create', exact: true }).click();
  await expect(ivy.getByRole('list', { name: 'Seats' })).toContainText('Ivy (you)');
  const id = await gameId(ivy);

  const jon = await (await browser.newContext()).newPage();
  await jon.goto(appURL());
  await openOnline(jon);
  await jon.getByLabel('Nickname').fill('Jon');
  await jon.getByLabel('Server address').fill(address());
  await jon.getByRole('button', { name: 'Connect' }).click();
  await lobbyGame(jon, id).getByRole('button', { name: 'Join' }).click();
  await expect(jon.getByRole('list', { name: 'Seats' })).toContainText('Jon (you)');
  await expect(ivy.getByRole('list', { name: 'Seats' })).toContainText('Jon');
});

test('an app that hosts offers no "Create a game" in its Games list (N-17)', async ({ page }) => {
  await page.goto(appURL());
  await openHost(page);
  await page.getByLabel('Nickname').fill('Kim');
  await page.getByRole('button', { name: 'Continue a saved game' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Create a game' })).toHaveCount(0);
});

test('deck ratios are on by default and locked when the sets are too small (GS-07)', async ({ page }) => {
  await realAssets(page);
  await page.goto(`${appURL()}&requiem`);
  await openHost(page);
  await page.getByLabel('Nickname').fill('Ratia');
  await page.getByLabel('Detailed setup').check();
  const sets = page.getByRole('group', { name: 'Card sets' });
  const ratios = page.getByLabel('Deck ratios');
  const label = page.locator('label', { has: ratios });
  // Every set together has the cards for the ratios.
  await expect(ratios).toBeChecked();
  await expect(ratios).toBeEnabled();

  // Without Requiem the groups are short: the box is off and locked,
  // and the tooltip says why.
  await sets.getByLabel('Requiem').uncheck();
  await expect(ratios).toBeDisabled();
  await expect(ratios).not.toBeChecked();
  await expect(label).toHaveAttribute('title', /^Too few cards/);
  await shot(page, 'host-deck-ratios-locked');

  // Requiem without Four Souls+ is enough for 3 players, not for 2.
  await sets.getByLabel('Requiem').check();
  await sets.getByLabel('Four Souls+').uncheck();
  await expect(ratios).toBeDisabled();
  await page.getByLabel('Players').selectOption('3');
  await expect(ratios).toBeEnabled();
  await expect(ratios).toBeChecked();
  await expect(label).not.toHaveAttribute('title', /./);

  await page.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(page.getByRole('list', { name: 'Seats' })).toContainText('Ratia (you)');
  await expect(page.getByText(/deck ratios/)).toBeVisible();
});


// Online is for servers: an address and Recent, no LAN games (15.98).
// Back and Esc go one level up from each screen.
test('Online: a server, its games, and one level back each time (15.98)', async ({ page }) => {
  await realAssets(page);
  await page.goto(appURL());
  await openOnline(page);
  await expect(page.getByRole('heading', { name: 'Online' })).toBeVisible();
  await expect(page.getByRole('region', { name: 'On your network' })).toHaveCount(0);
  await page.getByLabel('Nickname').fill('Ola');
  await page.getByLabel('Server address').fill(address());
  await page.getByRole('button', { name: 'Connect' }).click();
  await page.getByRole('button', { name: 'Create a game' }).click();
  await expect(page.getByRole('heading', { name: 'Create a game' })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('heading', { name: 'Online' })).toBeVisible();
  await expect(page.getByRole('list', { name: 'Recent' })).toContainText('Test server');
  await shot(page, 'online');
  await page.keyboard.press('Escape');
  await expect(page.getByRole('group', { name: 'Play' })).toBeVisible();
});

// Recent rows show at once, then whether each server is up (15.100).
test('Recent shows which servers are up and their open games (15.100)', async ({ page }) => {
  // Nobody listens on the port next to the test server.
  const stopped = `127.0.0.1:${String(serverPort() - 1)}`;
  await page.addInitScript(
    ([up, down]) => {
      const key = 'fake-settings';
      if (!localStorage.getItem(key)) {
        const recent = [
          { address: up, name: 'Test server' },
          { address: down, name: 'Gone server' },
        ];
        const rejoin = { address: '', token: '', game: '', via: '' };
        localStorage.setItem(key, JSON.stringify({ nickname: 'Rex', recent, animation: 'off', fullscreen: false, mat: '', dice: 'tear', rejoin }));
      }
    },
    [address(), stopped],
  );
  await realAssets(page);
  await page.goto(appURL());
  await openOnline(page);
  const recent = page.getByRole('list', { name: 'Recent' });
  const row = (name: string) => recent.getByRole('listitem').filter({ has: page.getByRole('button', { name, exact: true }) });
  // Other tests share the server, so its count is any number.
  await expect(row('Test server')).toContainText(/\d+ open games?/);
  await expect(row('Gone server')).toContainText('Down', { timeout: 8000 });
  await shot(page, 'online-status');

  // Back from its games asks again; the count follows the server.
  await row('Test server').getByRole('button', { name: 'Test server', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await page.getByRole('button', { name: 'Back' }).click();
  await expect(row('Test server')).toContainText(/\d+ open games?/);
});
