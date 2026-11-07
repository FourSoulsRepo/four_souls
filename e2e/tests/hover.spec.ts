// Hand hover and glow (15.27): playable hand cards glow green, the rest
// stay plain (H-12). A hovered hand card rises large from its place
// (H-13); any other card zooms right beside itself (F-16, T-27). Keyword
// tips stand beside the zoom.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';
import { cardName, tr } from './locale';

const fast = () => `${appURL()}&anim=off`;
const gallery = '/?fake&nosplash&screen=cards';
// green is the halo behind a playable hand card.
const green = /rgba\(110, 255, 70, /;

test.setTimeout(120_000);

// glow is the halo behind a hand card, drawn by its slot in the fan.
function glow(card: Locator): Promise<string> {
  return card.evaluate((el) => (el.parentElement ? getComputedStyle(el.parentElement, '::before').boxShadow : ''));
}

async function box(l: Locator) {
  const b = await l.boundingBox();
  if (!b) {
    throw new Error('no box');
  }
  return b;
}

// hoverTop hovers a hand card by its visible top: it hangs past the
// window's edge, and a locator hover would scroll it into view.
async function hoverTop(page: Page, card: Locator) {
  const b = await box(card);
  await page.mouse.move(b.x + b.width / 2, b.y + 30);
}

test('playable hand cards glow green; a hovered one rises large (15.27)', async ({ browser }) => {
  const pages = await startGame(browser, ['Glo', 'Hov'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const hand = (p: Page) => p.getByRole('group', { name: tr('table.your_hand') }).locator('.card');
  const playable = active.locator('.hand .card.hl-act');

  await shot(active, 'hand-glow');
  // The active player may play loot: those cards glow green.
  await expect(playable.first()).toBeVisible();
  for (const card of await playable.all()) {
    expect(await glow(card)).toMatch(green);
  }
  for (const card of await active.locator('.hand .card:not(.hl-act)').all()) {
    expect(await glow(card)).not.toMatch(green);
  }
  // The other player holds no priority: nothing in their hand glows.
  await expect(hand(other).first()).toBeVisible();
  await expect(other.locator('.hand .card.hl-act')).toHaveCount(0);
  for (const card of await hand(other).all()) {
    expect(await glow(card)).not.toMatch(green);
  }

  // A hovered hand card rises: its zoom stands on the window's bottom
  // edge, over the card.
  const card = hand(active).first();
  // Its neighbours stay where they are (20.5): the same edges, give or
  // take the browser's float noise.
  const others = active.locator('.own .hand-card');
  const edges = () =>
    others.evaluateAll((els) =>
      els.flatMap((el) => {
        const r = el.getBoundingClientRect();
        return [r.left, r.top, r.right, r.bottom];
      }),
    );
  // The fan settles first (19.122): the turn's loot card can come just
  // before, and the other cards then turn to their new pose for
  // 150 ms. Only the cards' own transitions count, not their glow.
  await expect.poll(() => others.evaluateAll((els) => els.filter((el) => el.getAnimations().length > 0).length)).toBe(0);
  const before = await edges();
  await hoverTop(active, card);
  const zoom = active.getByTestId('card-zoom').locator('.card.size-zoom');
  // A hand card answers at once, with no wait (20.5, ADR 027).
  await expect(zoom).toBeVisible({ timeout: 250 });
  const after = await edges();
  expect(after).toHaveLength(before.length);
  after.forEach((v, i) => {
    expect(v).toBeCloseTo(before.at(i) ?? NaN, 1);
  });
  await shot(active, 'hand-hover');
  const c = await box(card);
  const z = await box(zoom);
  // The zoom's body: an unusable card's reason stands under it (15.29).
  const body = await box(active.getByTestId('card-zoom').locator('.card-zoom-body'));
  expect(body.y + body.height).toBeGreaterThan(720 - 40);
  expect(body.y + body.height).toBeLessThanOrEqual(720);
  expect(z.x).toBeLessThan(c.x + c.width / 2);
  expect(z.x + z.width).toBeGreaterThan(c.x + c.width / 2);

  // A monster zooms right beside itself, not over it.
  const monster = active.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  await monster.hover();
  await expect(zoom).toBeVisible();
  await shot(active, 'monster-hover');
  const m = await box(monster);
  const mz = await box(zoom);
  const gap = Math.max(mz.x - (m.x + m.width), m.x - (mz.x + mz.width));
  expect(gap).toBeGreaterThanOrEqual(0);
  expect(gap).toBeLessThan(40);
  expect(mz.y).toBeLessThan(m.y + m.height);
  expect(mz.y + mz.height).toBeGreaterThan(m.y);
});

// A table card's zoom waits for the pointer to rest (20.5, ADR 027
// principle 2): none at 0.2 s, there at 0.6 s. While it shows, the
// next card's comes at once.
test('a table card zooms after a pause; the next one at once (20.5)', async ({ page }) => {
  await page.goto(gallery);
  const row = page.getByLabel('Table size');
  const zoom = page.getByTestId('card-zoom').locator('.card.size-zoom');
  const isaac = await box(row.getByRole('img', { name: cardName('isaac') }));
  const start = Date.now();
  await page.mouse.move(isaac.x + isaac.width / 2, isaac.y + isaac.height / 2);
  await page.waitForTimeout(Math.max(0, 200 - (Date.now() - start)));
  expect(await zoom.count()).toBe(0);
  // It grows out of the card (0.15 s): caught while it runs.
  await page.waitForFunction(() => (document.querySelector('.card-zoom-body')?.getAnimations().length ?? 0) > 0, undefined, { polling: 'raf' });
  await page.waitForTimeout(Math.max(0, 600 - (Date.now() - start)));
  expect(await zoom.isVisible()).toBe(true);
  await shot(page, 'zoom-after-pause');
  const d6 = await box(row.getByRole('img', { name: cardName('the_d6') }));
  await page.mouse.move(d6.x + d6.width / 2, d6.y + d6.height / 2);
  await expect(zoom).toHaveAttribute('aria-label', cardName('the_d6'), { timeout: 250 });
  // Away from every card, it closes.
  await page.mouse.move(2, 2);
  await expect(zoom).toHaveCount(0);
});

// With animations off the zoom still waits, but shows at full size at
// once: no growth (20.5).
test('with animations off a table zoom does not grow (20.5)', async ({ page }) => {
  await page.goto(`${gallery}&anim=off`);
  const card = page.getByLabel('Table size').getByRole('img', { name: cardName('isaac') });
  await card.hover();
  const body = page.getByTestId('card-zoom').locator('.card-zoom-body');
  await expect(body).toBeVisible();
  expect(await body.evaluate((el) => el.getAnimations().length)).toBe(0);
});

test('the zoom shows keyword tips beside the card (15.27)', async ({ page }) => {
  await page.goto(gallery);
  const zoom = page.getByTestId('card-zoom');
  const row = page.getByLabel('Table size');
  await row.getByRole('img', { name: cardName('the_d6') }).hover();
  const tips = zoom.getByRole('note', { name: tr('zoom.keywords') });
  await expect(tips).toContainText('Eternal');
  await expect(tips).toContainText("This can't be destroyed or put into discard.");
  await shot(page, 'zoom-tips');
  // A card without keywords has no tips.
  await row.getByRole('img', { name: cardName('isaac') }).hover();
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', cardName('isaac'));
  await expect(tips).toHaveCount(0);
});

// reach is how much of each hand card the pointer can hit, in square
// pixels: the points inside the window where that card is on top.
function reach(page: Page): Promise<number[]> {
  return page.evaluate(() =>
    [...document.querySelectorAll('.own .hand-card')].map((el) => {
      const r = el.getBoundingClientRect();
      let hits = 0;
      for (let x = Math.ceil(r.left); x < r.right; x += 2) {
        for (let y = Math.ceil(r.top); y < Math.min(r.bottom, window.innerHeight); y += 2) {
          if (document.elementFromPoint(x, y)?.closest('.hand-card') === el) {
            hits++;
          }
        }
      }
      return hits * 4;
    }),
  );
}

test('every card of a 12-card hand can be pointed at (19.15)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ful', 'Lan'], { url: fast() });
  const page = await activePage(pages, click);
  await expect(page.getByRole('group', { name: tr('table.your_hand') }).locator('.card').first()).toBeVisible();
  // Copies of the own hand stand in for a hand of 12; only the screen
  // changes, nobody acts on it.
  await page.evaluate(() => {
    interface Card {
      id: number;
      card: string;
    }
    interface S {
      game: { view: { players: { id: number; hand?: Card[]; hand_size: number }[] } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const { game } = store.getState();
    if (!game) {
      throw new Error('no game');
    }
    const players = game.view.players.map((p) => {
      const own = p.hand;
      if (!own?.length) {
        return p;
      }
      const hand = Array.from({ length: 12 }, (_, i) => ({ ...own[i % own.length], id: 9000 + i }) as Card);
      return { ...p, hand, hand_size: hand.length };
    });
    store.setState({ game: { ...game, view: { ...game.view, players } } });
  });
  const cards = page.getByRole('group', { name: tr('table.your_hand') }).locator('.hand-card');
  await expect(cards).toHaveCount(12);
  await page.mouse.move(0, 0);
  await shot(page, 'hand-12');
  // Each card, the end ones too, shows at least a strip of 30 by 50
  // pixels to the pointer.
  const areas = await reach(page);
  for (const [i, a] of areas.entries()) {
    expect(a, `card ${String(i)}`).toBeGreaterThanOrEqual(30 * 50);
  }
});

// A zoom closes when the pointer leaves its card, even when the part
// under the pointer was redrawn meanwhile: a card whose image failed to
// load swaps it for a drawn face, and no leave event came (19.50).
test('a zoom closes when the pointer leaves a card whose face was redrawn (19.50)', async ({ page }) => {
  await page.goto(gallery);
  const card = page.getByLabel('Table size').getByRole('img', { name: cardName('big_spider') });
  await card.hover();
  const zoom = page.getByTestId('card-zoom');
  await expect(zoom).toBeVisible();
  // The deepest hovered node inside the card goes, as a redraw does.
  await card.evaluate((el) => {
    const hovered = [...el.querySelectorAll(':hover')].at(-1);
    if (!hovered) {
      throw new Error('nothing hovered in the card');
    }
    hovered.remove();
  });
  await page.mouse.move(2, 2);
  await expect(zoom).toHaveCount(0);
});
