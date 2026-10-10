// Host and join screens (7.3), against the real Go server.
import { expect, test } from '@playwright/test';
import { appURL, serverPort } from './server';
import { shot } from './shot';

const address = `127.0.0.1:${String(serverPort)}`;

test('two apps reach the same lobby and start the game', async ({ browser }) => {
  // Separate contexts: two apps with their own settings.
  const ann = await (await browser.newContext()).newPage();
  const bo = await (await browser.newContext()).newPage();

  await ann.goto(appURL);
  await ann.getByRole('button', { name: 'Host a game' }).click();
  await ann.getByLabel('Nickname').fill('Ann');
  await ann.getByLabel('Detailed setup').check();
  await shot(ann, 'host');
  await ann.getByRole('button', { name: 'Host', exact: true }).click();
  const annSeats = ann.getByRole('list', { name: 'Seats' });
  await expect(annSeats).toContainText('Ann (you)');
  await expect(ann.getByText(address)).toBeVisible();

  await bo.goto(appURL);
  await bo.getByRole('button', { name: 'Join a game' }).click();
  await bo.getByLabel('Nickname').fill('Bo');
  await bo.getByLabel('Host address').fill(address);
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
  await expect(ann.getByRole('heading', { name: 'The game has started' })).toBeVisible();
  await expect(bo.getByRole('heading', { name: 'The game has started' })).toBeVisible();
});

test('a joined address is remembered', async ({ page }) => {
  await page.goto(appURL);
  await page.getByRole('button', { name: 'Join a game' }).click();
  await page.getByLabel('Nickname').fill('Cy');
  await page.getByLabel('Host address').fill(address);
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('heading', { name: 'Games' })).toBeVisible();
  await page.getByRole('button', { name: 'Back' }).click();
  await page.getByRole('button', { name: 'Join a game' }).click();
  await expect(page.getByRole('button', { name: address })).toBeVisible();
});

test('a wrong address says why', async ({ page }) => {
  await page.goto(appURL);
  await page.getByRole('button', { name: 'Join a game' }).click();
  await page.getByLabel('Nickname').fill('Di');
  await page.getByLabel('Host address').fill('not an address');
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('status')).toHaveText('Type an address like 192.168.1.5 or 192.168.1.5:4774.');
  await page.getByLabel('Host address').fill('127.0.0.1:1');
  await page.getByRole('button', { name: 'Connect' }).click();
  await expect(page.getByRole('status')).toHaveText('Could not connect to the server.');
});
