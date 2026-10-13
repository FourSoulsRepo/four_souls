// Host and join screens (7.3), against the real Go server.
import { expect, test } from '@playwright/test';
import { appURL, serverPort } from './server';
import { shot } from './shot';

// address is the project's test server.
const address = () => `127.0.0.1:${String(serverPort())}`;

test('two apps reach the same lobby and start the game', async ({ browser }) => {
  // Separate contexts: two apps with their own settings.
  const ann = await (await browser.newContext()).newPage();
  const bo = await (await browser.newContext()).newPage();

  await ann.goto(appURL());
  await ann.getByRole('button', { name: 'Host a game' }).click();
  await ann.getByLabel('Nickname').fill('Ann');
  await ann.getByLabel('Detailed setup').check();
  await shot(ann, 'host');
  await ann.getByRole('button', { name: 'Host', exact: true }).click();
  const annSeats = ann.getByRole('list', { name: 'Seats' });
  await expect(annSeats).toContainText('Ann (you)');
  await expect(ann.getByText(address())).toBeVisible();

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

test('the host picks card sets; the Base Game is always in', async ({ browser }) => {
  const host = await (await browser.newContext()).newPage();
  const guest = await (await browser.newContext()).newPage();

  await host.goto(appURL());
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByLabel('Nickname').fill('Tess');
  const sets = host.getByRole('group', { name: 'Card sets' });
  await expect(sets.getByLabel('Base Game')).toBeChecked();
  await expect(sets.getByLabel('Base Game')).toBeDisabled();
  await sets.getByLabel('Gold Box').check();
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
});
