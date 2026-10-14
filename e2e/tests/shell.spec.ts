// App shell (7.2): splash, main menu, settings.
import { expect, test, type Locator } from '@playwright/test';
import { realAssets, shot } from './shot';

test('splash shows the fan-game notice, then the menu', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake');
  await expect(page.getByText('Four Souls is an unofficial')).toBeVisible();
  await shot(page, 'splash');
  await expect(page.getByRole('button', { name: 'Host a game' })).toBeVisible({ timeout: 10_000 });
  await expect(page.getByText('app demo · rules engine demo')).toBeVisible();
  await shot(page, 'menu');
});

test('titles and menu buttons use the game fonts (15.40)', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  // WebKit drops the quotes around family names; Chromium keeps them.
  const family = (el: Locator) => el.evaluate((e) => getComputedStyle(e).fontFamily.replace(/"/g, ''));
  // Without the font files the faces fail to load, and the system font
  // at the end of each stack shows instead.
  const title = page.getByRole('heading', { name: 'Four Souls' });
  await expect.poll(() => family(title)).toMatch(/^Four Souls Cyrillic, -apple-system/);
  const button = page.getByRole('button', { name: 'Host a game' });
  await expect.poll(() => family(button)).toMatch(/^Four Souls Title, Four Souls Cyrillic, -apple-system/);
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect.poll(() => family(page.getByRole('heading', { name: 'Settings' }))).toMatch(/^Four Souls Cyrillic,/);
});

test('nickname survives a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Nickname').fill('  Ann  ');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await expect(page.getByLabel('Nickname')).toHaveValue('Ann');
  await shot(page, 'settings');

  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Nickname')).toHaveValue('Ann');
});

test('a bad nickname shows why', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Nickname').fill('x'.repeat(21));
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('nickname must be 1–20 characters');
  await page.getByRole('button', { name: 'Back' }).click();
  await expect(page.getByRole('button', { name: 'Host a game' })).toBeVisible();
});
