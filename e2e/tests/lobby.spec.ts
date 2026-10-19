// Host and join screens (7.3), against the real Go server.
import { expect, test } from './fixtures';
import { appURL, serverPort } from './server';
import { realAssets, shot } from './shot';

// address is the project's test server.
const address = () => `127.0.0.1:${String(serverPort())}`;

test('two apps reach the same lobby and start the game', async ({ browser }) => {
  // Separate contexts: two apps with their own settings.
  const ann = await (await browser.newContext()).newPage();
  const bo = await (await browser.newContext()).newPage();

  await realAssets(ann);
  await ann.goto(appURL());
  await ann.getByRole('button', { name: 'Host a game' }).click();
  await ann.getByLabel('Nickname').fill('Ann');
  await ann.getByLabel('Detailed setup').check();
  await shot(ann, 'host');
  await ann.getByRole('button', { name: 'Host', exact: true }).click();
  const annSeats = ann.getByRole('list', { name: 'Seats' });
  await expect(annSeats).toContainText('Ann (you)');
  await expect(ann.getByText(address())).toBeVisible();

  await realAssets(bo);
  await bo.goto(appURL());
  await bo.getByRole('button', { name: 'Join a game' }).click();
  await bo.getByLabel('Nickname').fill('Bo');
  await bo.getByLabel('Host address').fill(address());
  await bo.getByRole('button', { name: 'Connect' }).click();
  const game = bo.getByRole('listitem').filter({ hasText: "Ann's game" });
  await expect(game).toContainText('1/2 players');
  await shot(bo, 'games');
  await game.getByRole('button', { name: 'Join' }).click();
  await expect(bo.getByRole('list', { name: 'Seats' })).toContainText('Bo (you)');
  await expect(annSeats).toContainText('Bo');
  await shot(ann, 'table');

  await ann.getByRole('button', { name: 'Ready', exact: true }).click();
  await expect(bo.getByRole('list', { name: 'Seats' })).toContainText('Ready');
  await bo.getByRole('button', { name: 'Ready', exact: true }).click();
  await expect(ann.getByRole('main', { name: 'Turn 1' })).toBeVisible();
  await expect(bo.getByRole('main', { name: 'Turn 1' })).toBeVisible();
});

test('a joined address is remembered', async ({ page }) => {
  await page.goto(appURL());
  await page.getByRole('button', { name: 'Join a game' }).click();
  await page.getByLabel('Nickname').fill('Cy');
  await page.getByLabel('Host address').fill(address());
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await page.getByRole('button', { name: 'Back' }).click();
  await page.getByRole('button', { name: 'Join a game' }).click();
  await expect(page.getByRole('button', { name: address() })).toBeVisible();
});

test('a wrong address says why', async ({ page }) => {
  await page.goto(appURL());
  await page.getByRole('button', { name: 'Join a game' }).click();
  await page.getByLabel('Nickname').fill('Di');
  await page.getByLabel('Host address').fill('not an address');
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('status')).toHaveText('Type an address like 192.168.1.5 or 192.168.1.5:4774.');
  // Nobody listens on the port next to the test server.
  await page.getByLabel('Host address').fill(`127.0.0.1:${String(serverPort() - 1)}`);
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('status')).toHaveText('Could not connect to the server.', { timeout: 8000 });
});

test('a game found on the LAN is joined without typing (N-13)', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  await host.goto(appURL());
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByLabel('Nickname').fill('Lex');
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toContainText('Lex (you)');

  // The fake LAN holds the test server; its real games list follows.
  const guest = await (await browser.newContext()).newPage();
  await realAssets(guest);
  await guest.goto(appURL());
  await guest.getByRole('button', { name: 'Join a game' }).click();
  await guest.getByLabel('Nickname').fill('Mio');
  const lan = guest.getByRole('region', { name: 'On your network' });
  const found = lan.getByRole('listitem').filter({ hasText: address() });
  await expect(found).toContainText('1/2 players');
  await shot(guest, 'join-lan');
  await found.getByRole('button', { name: 'Join' }).click();
  await guest.getByRole('listitem').filter({ hasText: "Lex's game" }).getByRole('button', { name: 'Join' }).click();
  await expect(guest.getByRole('list', { name: 'Seats' })).toContainText('Mio (you)');
});

test('with no game on the LAN, the address field is the way', async ({ page }) => {
  await page.goto(`${appURL()}&nolan`);
  await page.getByRole('button', { name: 'Join a game' }).click();
  const lan = page.getByRole('region', { name: 'On your network' });
  await expect(lan).toContainText('No games found');
  await expect(lan.getByRole('button', { name: 'Search again' })).toBeEnabled();
});

test('a game password keeps others out (N-14)', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  await host.goto(appURL());
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByLabel('Nickname').fill('Nox');
  await host.getByLabel('Password (optional)').fill('pw-123');
  await host.getByRole('button', { name: 'Host', exact: true }).click();
  await expect(host.getByRole('list', { name: 'Seats' })).toContainText('Nox (you)');

  const guest = await (await browser.newContext()).newPage();
  await realAssets(guest);
  await guest.goto(appURL());
  await guest.getByRole('button', { name: 'Join a game' }).click();
  await guest.getByLabel('Nickname').fill('Pip');
  await guest.getByLabel('Host address').fill(address());
  await guest.getByRole('button', { name: 'Connect' }).click();
  const game = guest.getByRole('listitem').filter({ hasText: "Nox's game" });
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
  await host.getByRole('button', { name: 'Host a game' }).click();
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

  await guest.goto(appURL());
  await guest.getByRole('button', { name: 'Join a game' }).click();
  await guest.getByLabel('Nickname').fill('Lu');
  await guest.getByLabel('Host address').fill(address());
  await guest.getByRole('button', { name: 'Connect' }).click();
  const game = guest.getByRole('listitem').filter({ hasText: "Tess's game" });
  await expect(game).toContainText('Base Game, Gold Box');
  await expect(game).not.toContainText('Four Souls+');

  // The next host screen in the same app starts from this setup.
  const again = await ctx.newPage();
  await again.goto(appURL());
  await again.getByRole('button', { name: 'Host a game' }).click();
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
  await ivy.getByRole('button', { name: 'Join a game' }).click();
  await ivy.getByLabel('Nickname').fill('Ivy');
  await ivy.getByLabel('Host address').fill(address());
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

  const jon = await (await browser.newContext()).newPage();
  await jon.goto(appURL());
  await jon.getByRole('button', { name: 'Join a game' }).click();
  await jon.getByLabel('Nickname').fill('Jon');
  await jon.getByLabel('Host address').fill(address());
  await jon.getByRole('button', { name: 'Connect' }).click();
  await jon.getByRole('listitem').filter({ hasText: "Ivy's game" }).getByRole('button', { name: 'Join' }).click();
  await expect(jon.getByRole('list', { name: 'Seats' })).toContainText('Jon (you)');
  await expect(ivy.getByRole('list', { name: 'Seats' })).toContainText('Jon');
});

test('an app that hosts offers no "Create a game" in its Games list (N-17)', async ({ page }) => {
  await page.goto(appURL());
  await page.getByRole('button', { name: 'Host a game' }).click();
  await page.getByLabel('Nickname').fill('Kim');
  await page.getByRole('button', { name: 'Continue a saved game' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Create a game' })).toHaveCount(0);
});

test('deck ratios are on by default and locked when the sets are too small (GS-07)', async ({ page }) => {
  await realAssets(page);
  await page.goto(`${appURL()}&requiem`);
  await page.getByRole('button', { name: 'Host a game' }).click();
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
