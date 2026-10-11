// Settings screen (7.10): ST-01, ST-02, ST-03, ST-06.
import { expect, test } from '@playwright/test';
import { shot } from './shot';

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
