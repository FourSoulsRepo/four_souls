// The Update screen (18.24, ADR 029, UC-43): the fake bridge answers the
// check and the test sends the download's progress; no real network.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { shot } from './shot';
import { tr } from './locale';

const mb = 1 << 20;

interface FakeWindow {
  fakeUpdate?: object;
  fakeDownload?: { finish: (error?: string) => void };
  fakeOpened?: string;
  fakeQuit: boolean;
  fakeEmit: (name: string, ...data: unknown[]) => void;
}

// openUpdate opens the Update screen from the main menu; status is the
// check's answer, the fake's newer release when left out.
async function openUpdate(page: Page, status?: object): Promise<void> {
  await page.goto('/?fake&nosplash');
  if (status) {
    await page.evaluate((s) => {
      (window as unknown as FakeWindow).fakeUpdate = s;
    }, status);
  }
  await page.getByRole('button', { name: tr('menu.update') }).click();
  await expect(page.getByRole('heading', { name: tr('update.title') })).toBeVisible();
}

// progress sends one download progress event, as Go's Updater does.
async function progress(page: Page, done: number, total: number): Promise<void> {
  await page.evaluate(
    ([d, t]) => {
      (window as unknown as FakeWindow).fakeEmit('update:progress', { done: d, total: t });
    },
    [done, total],
  );
}

// noUpdate is a check whose release can't be installed, for reason.
function noUpdate(reason: string): object {
  const releases = 'https://github.com/FourSoulsRepo/four_souls/releases';
  return { current: '1.0.0', latest: 'v1.0.0', notes: '', url: `${releases}/tag/v1.0.0`, releases, asset: '', newer: false, canInstall: false, reason };
}

test('the menu opens the Update screen; it checks only on a press', async ({ page }) => {
  await openUpdate(page);
  await expect(page.getByText(tr('update.current', { version: 'demo' }))).toBeVisible();
  await expect(page.getByRole('button', { name: tr('update.check') })).toBeEnabled();
  await expect(page.getByText(tr('update.latest', { version: 'v9.9.9' }))).toHaveCount(0);
  await expect(page.getByRole('button', { name: tr('update.download') })).toHaveCount(0);
  await shot(page, 'update');
  await page.getByRole('button', { name: tr('update.back') }).click();
  await expect(page.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
});

test('a check shows the latest version and its notes', async ({ page }) => {
  await openUpdate(page);
  await page.getByRole('button', { name: tr('update.check') }).click();
  await expect(page.getByText(tr('update.latest', { version: 'v9.9.9' }))).toBeVisible();
  const notes = page.getByRole('region', { name: tr('update.notes') });
  await expect(notes).toContainText('A new card set.');
  await expect(page.getByRole('status')).toHaveText(tr('update.ready'));
  await expect(page.getByRole('button', { name: tr('update.download') })).toBeEnabled();
  await expect(page.getByRole('button', { name: tr('update.releases') })).toHaveCount(0);
  await shot(page, 'update-checked');
});

test('the download bar moves with the progress events, then installs', async ({ page }) => {
  await openUpdate(page);
  await page.getByRole('button', { name: tr('update.check') }).click();
  await page.getByRole('button', { name: tr('update.download') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('update.downloading'));
  // Busy: no second check or download while it runs.
  await expect(page.getByRole('button', { name: tr('update.check') })).toBeDisabled();
  await expect(page.getByRole('button', { name: tr('update.download') })).toBeDisabled();

  const bar = page.getByRole('progressbar', { name: tr('update.progress_label') });
  const fill = bar.locator('div');
  const width = () => fill.evaluate((e) => e.getBoundingClientRect().width);
  await progress(page, 0, 20 * mb);
  await expect(bar).toHaveAttribute('aria-valuenow', '0');
  await expect(page.getByText(tr('update.progress', { percent: 0, done: '0.0', total: '20.0' }))).toBeVisible();
  const empty = await width();
  await progress(page, 5 * mb, 20 * mb);
  await expect(bar).toHaveAttribute('aria-valuenow', '25');
  await expect(page.getByText(tr('update.progress', { percent: 25, done: '5.0', total: '20.0' }))).toBeVisible();
  const quarter = await width();
  expect(quarter).toBeGreaterThan(empty);
  await progress(page, 15 * mb, 20 * mb);
  await expect(bar).toHaveAttribute('aria-valuenow', '75');
  expect(await width()).toBeCloseTo(quarter * 3, 0);
  await shot(page, 'update-download');

  // A server without Content-Length: MB only, no percent.
  await progress(page, 3 * mb, 0);
  await expect(page.getByText(tr('update.progress_mb', { done: '3.0' }))).toBeVisible();
  await expect(bar).not.toHaveAttribute('aria-valuenow');

  await progress(page, 20 * mb, 20 * mb);
  await page.evaluate(() => {
    (window as unknown as FakeWindow).fakeDownload?.finish();
  });
  await expect(page.getByRole('status')).toHaveText(tr('update.downloaded'));
  await expect(bar).toHaveAttribute('aria-valuenow', '100');
  const quit = () => page.evaluate(() => (window as unknown as FakeWindow).fakeQuit);
  expect(await quit()).toBe(false);
  await page.getByRole('button', { name: tr('update.install') }).click();
  await expect.poll(quit).toBe(true);
});

test('a failed download says why in plain words', async ({ page }) => {
  await openUpdate(page);
  await page.getByRole('button', { name: tr('update.check') }).click();
  await page.getByRole('button', { name: tr('update.download') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('update.downloading'));
  await page.evaluate(() => {
    (window as unknown as FakeWindow).fakeDownload?.finish('{"key":"error.update_checksum"}');
  });
  await expect(page.getByRole('status')).toHaveText(tr('error.update_checksum'));
  // It may try again, or update by hand.
  await expect(page.getByRole('button', { name: tr('update.download') })).toBeEnabled();
  await expect(page.getByRole('button', { name: tr('update.releases') })).toBeVisible();
});

test('no update possible: the reason and the releases page', async ({ page }) => {
  const status = noUpdate('not_writable');
  await openUpdate(page, status);
  await page.getByRole('button', { name: tr('update.check') }).click();
  await expect(page.getByText(tr('update.reason.not_writable'))).toBeVisible();
  await expect(page.getByRole('button', { name: tr('update.download') })).toHaveCount(0);
  await shot(page, 'update-no-update');
  await page.getByRole('button', { name: tr('update.releases') }).click();
  await expect.poll(() => page.evaluate(() => (window as unknown as FakeWindow).fakeOpened)).toBe('https://github.com/FourSoulsRepo/four_souls/releases/tag/v1.0.0');
});

test('up to date: says so, with no releases link', async ({ page }) => {
  await openUpdate(page, noUpdate('not_newer'));
  await page.getByRole('button', { name: tr('update.check') }).click();
  await expect(page.getByText(tr('update.reason.not_newer'))).toBeVisible();
  await expect(page.getByRole('button', { name: tr('update.releases') })).toHaveCount(0);
});

test('a failed check says why and links the releases page', async ({ page }) => {
  await openUpdate(page, { error: '{"key":"error.update_check"}' });
  await page.getByRole('button', { name: tr('update.check') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('error.update_check'));
  await expect(page.getByRole('button', { name: tr('update.check') })).toBeEnabled();
  await page.getByRole('button', { name: tr('update.releases') }).click();
  await expect.poll(() => page.evaluate(() => (window as unknown as FakeWindow).fakeOpened)).toBe('https://github.com/FourSoulsRepo/four_souls/releases');
});

test('long notes scroll in their box; the page never scrolls', async ({ page }) => {
  await page.setViewportSize({ width: 1024, height: 600 });
  const notes = Array.from({ length: 80 }, (_, i) => `- Change ${String(i + 1)}`).join('\n');
  await openUpdate(page, { ...noUpdate(''), latest: 'v2.0.0', notes, newer: true, canInstall: true, asset: 'a.zip' });
  await page.getByRole('button', { name: tr('update.check') }).click();
  const box = page.getByRole('region', { name: tr('update.notes') });
  await expect(box).toContainText('Change 80');
  expect(await box.evaluate((e) => e.scrollHeight > e.clientHeight)).toBe(true);
  expect(await page.evaluate(() => document.documentElement.scrollHeight <= window.innerHeight)).toBe(true);
  await expect(page.getByRole('button', { name: tr('update.back') })).toBeInViewport();
});
