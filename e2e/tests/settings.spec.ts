// Settings screen (7.10): ST-01, ST-02, ST-03, ST-06.
import { expect, test } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { realAssets, shot } from './shot';

test('the first start asks for a nickname, once', async ({ page }) => {
  await page.goto('/?fake&nosplash&firstrun');
  await expect(page.getByRole('heading', { name: 'Welcome!' })).toBeVisible();
  await shot(page, 'welcome');
  await page.getByLabel('Nickname').fill('Ada');
  await page.getByRole('button', { name: 'Continue' }).click();
  await expect(page.getByRole('button', { name: 'Host a game' })).toBeVisible();
  await page.reload();
  await expect(page.getByRole('button', { name: 'Host a game' })).toBeVisible();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Nickname')).toHaveValue('Ada');
});

test('fullscreen, theme and mat apply at once and survive a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Fullscreen').check();
  await page.getByLabel('Theme').selectOption('light');
  await page.getByRole('combobox', { name: 'Game mat' }).selectOption('basement');
  await expect(page.getByLabel('Game mat preview')).toHaveClass(/mat-basement/);
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');
  expect(await page.evaluate(() => (window as unknown as { fakeFullscreen: boolean }).fakeFullscreen)).toBe(true);
  await shot(page, 'settings-light');

  await page.reload();
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Fullscreen')).toBeChecked();
  await expect(page.getByRole('combobox', { name: 'Game mat' })).toHaveValue('basement');
  await page.getByLabel('Fullscreen').uncheck();
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  expect(await page.evaluate(() => (window as unknown as { fakeFullscreen: boolean }).fakeFullscreen)).toBe(false);
});

// Mat images (15.10, V-24): the setting puts the image under the table;
// without images, or with a missing one, the table shows colors.
test('a mat image lies under the table', async ({ browser }) => {
  const [mat] = await startGame(browser, ['Mat', 'Ned'], { url: `${appURL()}&mat=basement` });
  if (!mat) {
    throw new Error('no pages');
  }
  const table = mat.getByRole('main', { name: /^Turn / });
  await expect(table).toHaveClass(/mat-basement mat-image/);
  expect(await table.evaluate((el) => getComputedStyle(el).backgroundImage)).toContain('/cardimg/assets/mats/basement.webp');
  await shot(mat, 'table-mat');
});

test('the mat menu lists the mat images', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  const mats = page.getByRole('combobox', { name: 'Game mat' });
  await expect(mats.getByRole('option')).toHaveText(['Green felt', 'Basement', 'Cathedral']);
  await mats.selectOption('basement');
  await expect(page.getByLabel('Game mat preview')).toHaveClass(/mat-image/);
  await shot(page, 'settings-mat');
});

test('the first mat image is the default', async ({ page }) => {
  await page.goto(appURL());
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByRole('combobox', { name: 'Game mat' })).toHaveValue('basement');
  await expect(page.getByLabel('Game mat preview')).toHaveClass(/mat-basement mat-image/);
  // The image fails to load here; the basement colors stay under it.
  expect(await page.getByLabel('Game mat preview').evaluate((el) => getComputedStyle(el).backgroundImage)).toContain('rgb(90, 67, 48)');
});

test('without mat images the table shows colors', async ({ browser }) => {
  for (const [query, colors] of [
    ['&nomats', 'mat-basement'],
    ['&nomats&mat=felt', 'mat-felt'],
    ['&mat=dark_room', 'mat-basement'],
  ]) {
    const page = await (await browser.newContext()).newPage();
    await page.goto(`${appURL()}${query ?? ''}`);
    await page.getByRole('button', { name: 'Settings' }).click();
    await expect(page.getByLabel('Game mat preview')).toHaveClass(new RegExp(`${colors ?? ''}$`));
    await expect(page.getByLabel('Game mat preview')).not.toHaveClass(/mat-image/);
  }
});
