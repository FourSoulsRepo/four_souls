// App shell (7.2): splash, main menu, settings.
import type { Locator } from '@playwright/test';
import { expect, test } from './fixtures';
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

// The menu buttons stay put once they show: a title that changed size
// after the menu showed moved them, and a quick click hit the wrong
// button (19.49). Without the images the banner fails to load.
test('the menu buttons never move after they show (19.49)', async ({ page }) => {
  await page.addInitScript(() => {
    const w = window as unknown as { ys: number[] };
    w.ys = [];
    const t0 = performance.now();
    const frame = () => {
      const b = [...document.querySelectorAll('nav button')].find((x) => x.textContent === 'Settings');
      if (b) {
        const y = b.getBoundingClientRect().y;
        if (w.ys.at(-1) !== y) {
          w.ys.push(y);
        }
      }
      if (performance.now() - t0 < 2000) {
        requestAnimationFrame(frame);
      }
    };
    requestAnimationFrame(frame);
  });
  await page.goto('/?fake&nosplash');
  await expect(page.getByRole('button', { name: 'Settings' })).toBeVisible();
  await page.waitForTimeout(2000);
  expect(await page.evaluate(() => (window as unknown as { ys: number[] }).ys)).toHaveLength(1);
});

test('Exit is the last menu button and closes the app (15.79)', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await expect(page.getByRole('navigation').getByRole('button')).toHaveText(['Host a game', 'Join a game', 'Records', 'Settings', 'Exit']);
  const quit = () => page.evaluate(() => (window as unknown as { fakeQuit: boolean }).fakeQuit);
  expect(await quit()).toBe(false);
  await page.getByRole('button', { name: 'Exit' }).click();
  expect(await quit()).toBe(true);
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

test('the menu shows the banner, or the text title without it (15.41)', async ({ page }) => {
  const banner = page.getByRole('heading', { name: 'Four Souls' }).locator('img');
  // Served, the banner image is the title.
  const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="9"><rect width="16" height="9" fill="#a33"/></svg>';
  await page.route('**/cardimg/assets/banner.webp', (route) => route.fulfill({ body: svg, contentType: 'image/svg+xml' }));
  await page.goto('/?fake&nosplash');
  await expect(banner).toHaveAttribute('src', '/cardimg/assets/banner.webp');
  await expect(banner).toBeVisible();
  // Failing to load, as without the images, the text title shows.
  await page.unroute('**/cardimg/assets/banner.webp');
  await page.route('**/cardimg/assets/banner.webp', (route) => route.abort());
  await page.reload();
  await expect(page.getByRole('heading', { name: 'Four Souls' })).toHaveText('Four Souls');
  await expect(banner).toHaveCount(0);
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

test('the menus drift on the night sky, still with animations off (15.42)', async ({ page }) => {
  // The sky is a decoration with no role; its layers are found by class.
  const layers = page.locator('.sky .sky-layer');
  const transform = () => layers.last().evaluate((e) => getComputedStyle(e).transform);
  await page.goto('/?fake&nosplash');
  // The sky, then five star layers, behind the menu and moving.
  await expect(layers).toHaveCount(6);
  const start = await transform();
  await expect.poll(transform).not.toBe(start);
  // The menu draws over it, and clicks reach its buttons.
  await page.getByRole('button', { name: 'Host a game' }).click();
  await expect(page.getByRole('heading', { name: 'Host a game' })).toBeVisible();
  await expect(layers).toHaveCount(6);
  // Animations off: the sky stands still.
  await page.getByRole('button', { name: 'Back' }).click();
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Animations').selectOption('off');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(layers).toHaveCount(6);
  expect(await layers.last().evaluate((e) => getComputedStyle(e).animationName)).toBe('none');
});
