// App shell (7.2): splash, main menu, settings.
import type { Locator } from '@playwright/test';
import { expect, test } from './fixtures';
import { backToMenu, openHost } from './menu';
import { realAssets, shot } from './shot';
import { tr } from './locale';

test('splash shows the fan-game notice, then the menu', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake');
  await expect(page.getByText('Four Souls is an unofficial')).toBeVisible();
  await shot(page, 'splash');
  await expect(page.getByRole('button', { name: 'Play', exact: true })).toBeVisible({ timeout: 10_000 });
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
  await expect(page.getByRole('navigation').getByRole('button')).toHaveText(['Play', 'RecordsSoon', 'Settings', 'Exit']);
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
  const button = page.getByRole('button', { name: 'Play', exact: true });
  await expect.poll(() => family(button)).toMatch(/^Four Souls Title, Four Souls Cyrillic, -apple-system/);
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect.poll(() => family(page.getByRole('heading', { name: tr('settings.title') }))).toMatch(/^Four Souls Cyrillic,/);
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
  await page.getByLabel(tr('settings.nickname')).fill('  Ann  ');
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('settings.saved'));
  await expect(page.getByLabel(tr('settings.nickname'))).toHaveValue('Ann');
  await shot(page, 'settings');

  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel(tr('settings.nickname'))).toHaveValue('Ann');
});

test('a bad nickname shows why', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel(tr('settings.nickname')).fill('x'.repeat(21));
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(page.getByRole('status')).toHaveText('nickname must be 1–20 characters');
  await page.getByRole('button', { name: tr('settings.back') }).click();
  await expect(page.getByRole('button', { name: 'Play', exact: true })).toBeVisible();
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
  await openHost(page);
  await expect(page.getByRole('heading', { name: 'Host a game' })).toBeVisible();
  await expect(layers).toHaveCount(6);
  // Animations off: the sky stands still.
  await backToMenu(page);
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel(tr('settings.animations')).selectOption('off');
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(layers).toHaveCount(6);
  expect(await layers.last().evaluate((e) => getComputedStyle(e).animationName)).toBe('none');
});

// Play opens its group in place, under itself (M-06): the buttons above
// stay put, the ones below slide down. Still to come shows "Soon" (M-07).
test('Play opens LAN, Online and Learn to play in place (15.96)', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake&nosplash');
  const play = page.getByRole('button', { name: 'Play', exact: true });
  const settings = page.getByRole('button', { name: 'Settings' });
  const group = page.getByRole('group', { name: 'Play' });
  await expect(play).toHaveAttribute('aria-expanded', 'false');
  await expect(group).toHaveCount(0);
  const top = (l: typeof play) => l.evaluate((e) => e.getBoundingClientRect().top);
  // A hovered button lifts; the pointer leaves after each click. The
  // banner loads first in screenshots, so the menu has settled.
  await page.waitForLoadState('networkidle');
  const playTop = await top(play);
  const settingsTop = await top(settings);

  await play.click();
  await page.mouse.move(1, 1);
  await expect(play).toHaveAttribute('aria-expanded', 'true');
  await expect(group.getByRole('button')).toHaveText(['LAN', 'Online', 'Learn to playSoon']);
  await expect(group.getByRole('button', { name: 'Learn to play (soon)' })).toBeDisabled();
  await expect(page.getByRole('button', { name: 'Records (soon)' })).toBeDisabled();
  // Settings slides down by the group: three buttons and their gaps.
  const row = await group.getByRole('button', { name: 'LAN' }).evaluate((e) => (e as HTMLElement).offsetHeight + 12);
  await expect.poll(() => top(settings)).toBeCloseTo(settingsTop + 3 * row, 0);
  await expect.poll(() => top(play)).toBe(playTop);
  await shot(page, 'menu-play');

  // A second click closes it, and so does Esc.
  await play.click();
  await page.mouse.move(1, 1);
  await expect(group).toHaveCount(0);
  await expect.poll(() => top(settings)).toBe(settingsTop);
  await play.click();
  await page.mouse.move(1, 1);
  await expect(group).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(group).toHaveCount(0);
  await expect.poll(() => top(play)).toBe(playTop);
});

// Back from a Play screen finds the group still open.
test('the Play group stays open on the way back (15.96)', async ({ page }) => {
  await page.goto('/?fake&nosplash&anim=off');
  await page.getByRole('button', { name: 'Play', exact: true }).click();
  await page.getByRole('button', { name: 'Online', exact: true }).click();
  await page.getByRole('button', { name: 'Back' }).click();
  await expect(page.getByRole('group', { name: 'Play' })).toBeVisible();
});
