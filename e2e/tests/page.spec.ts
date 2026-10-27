// Not a web page (8.8, V-10): no screen scrolls as a whole, and a mouse
// drag selects nothing except text worth copying.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { backToMenu, openHost, openLan, openOnline } from './menu';
import { appURL } from './server';
import { activePage, click } from './turn';
import { tr, trRe } from './locale';

// selection is the text the page has selected now.
function selection(page: Page): Promise<string> {
  return page.evaluate(() => window.getSelection()?.toString() ?? '');
}

// dragAcross drags the mouse from one corner of the box to the other,
// or from a given point to the far corner.
async function dragAcross(page: Page, box: Locator, from?: { x: number; y: number }): Promise<void> {
  const r = await box.boundingBox();
  if (!r) {
    throw new Error('no box');
  }
  await page.mouse.move(from?.x ?? r.x + 4, from?.y ?? r.y + 4);
  await page.mouse.down();
  await page.mouse.move(r.x + r.width / 2, r.y + r.height / 2, { steps: 5 });
  await page.mouse.move(r.x + r.width - 4, r.y + r.height - 4, { steps: 5 });
  await page.mouse.up();
}

// lineMiddle is the middle of a row's first line of text. Chromium
// starts no selection from the pixel between two lines (19.81).
async function lineMiddle(row: Locator): Promise<{ x: number; y: number }> {
  return row.evaluate((el) => {
    // The row's last span holds its words; the first is the icon.
    const range = document.createRange();
    range.selectNodeContents(el.lastElementChild ?? el);
    const r = range.getClientRects()[0] ?? el.getBoundingClientRect();
    return { x: r.left + 2, y: r.top + r.height / 2 };
  });
}

// expectNoPageScroll fails when the window scrolls as a whole, even after
// the wheel turns.
async function expectNoPageScroll(page: Page): Promise<void> {
  await page.mouse.move(10, 10);
  await page.mouse.wheel(0, 2000);
  const scroll = await page.evaluate(() => {
    const el = document.documentElement;
    return {
      x: window.scrollX,
      y: window.scrollY,
      overW: el.scrollWidth > window.innerWidth,
      overH: el.scrollHeight > window.innerHeight,
    };
  });
  expect(scroll).toEqual({ x: 0, y: 0, overW: false, overH: false });
}

test('menu screens never scroll, even in a small window', async ({ page }) => {
  await page.setViewportSize({ width: 560, height: 360 });
  await page.goto(appURL());
  await expectNoPageScroll(page);
  const settings = (p: Page) => p.getByRole('button', { name: tr('menu.settings') }).click();
  for (const [open, sign] of [
    [settings, tr('settings.nickname')],
    [openHost, tr('host.players')],
    [openLan, tr('lobby.lan.address')],
    [openOnline, tr('lobby.online.address')],
  ] as const) {
    await open(page);
    await expect(page.getByLabel(sign)).toBeVisible();
    await expectNoPageScroll(page);
    await backToMenu(page);
  }
});

test('a drag across the menu selects nothing', async ({ page }) => {
  await page.goto(appURL());
  await dragAcross(page, page.locator('body'));
  expect(await selection(page)).toBe('');
});

test('a drag across the table selects nothing; the history and address do', async ({ browser }) => {
  const pages = await startGame(browser, ['Kim', 'Lou'], { url: `${appURL()}&anim=off` });
  // The turn starts with a loot, so the history has a line.
  const kim = await activePage(pages, click);
  const history = kim.getByRole('complementary', { name: tr('history.title'), exact: true });
  // The full history sits on the board; it starts hidden (15.22).
  await dragAcross(kim, kim.getByRole('main', { name: new RegExp(`^${trRe('game.turn').source}`) }).locator('.board'));
  expect(await selection(kim)).toBe('');
  await expectNoPageScroll(kim);

  // Logs stay copyable. The newest row is the loot; the drag starts in
  // its text, not on a line edge (19.81).
  await kim.getByRole('button', { name: tr('history.title'), exact: true }).click();
  await expect(history).toContainText(trRe('history.loot.you'));
  await dragAcross(kim, history.getByRole('list'), await lineMiddle(history.getByRole('listitem').last()));
  expect(await selection(kim)).not.toBe('');
});

test('the host address can be copied', async ({ page }) => {
  await page.goto(appURL());
  await openHost(page);
  await page.getByLabel(tr('lobby.nickname')).fill('Max');
  await page.getByRole('button', { name: tr('host.host'), exact: true }).click();
  const address = page.getByRole('list', { name: tr('lobby.seats.addresses') }).getByRole('listitem').first();
  await address.click({ clickCount: 3 });
  expect(await selection(page)).toContain('127.0.0.1');
});
