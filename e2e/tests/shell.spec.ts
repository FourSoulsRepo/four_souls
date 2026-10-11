// App shell (7.2): splash, main menu, settings.
import { expect, test } from '@playwright/test';
import { shot } from './shot';

test('splash shows the fan-game notice, then the menu', async ({ page }) => {
  await page.goto('/?fake');
  await expect(page.getByText('Four Souls is an unofficial')).toBeVisible();
  await shot(page, 'splash');
  await expect(page.getByRole('button', { name: 'Host a game' })).toBeVisible({ timeout: 10_000 });
  await expect(page.getByText('app demo · rules engine demo')).toBeVisible();
  await shot(page, 'menu');
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
