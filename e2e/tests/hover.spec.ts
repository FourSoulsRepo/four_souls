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
  const hand = (p: Page) => p.getByRole('group', { name: 'Your hand' }).locator('.card');
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
  await hoverTop(active, card);
  const zoom = active.getByTestId('card-zoom').locator('.card.size-zoom');
  await expect(zoom).toBeVisible();
  await shot(active, 'hand-hover');
  const c = await box(card);
  const z = await box(zoom);
  expect(z.y + z.height).toBeGreaterThan(720 - 40);
  expect(z.y + z.height).toBeLessThanOrEqual(720);
  expect(z.x).toBeLessThan(c.x + c.width / 2);
  expect(z.x + z.width).toBeGreaterThan(c.x + c.width / 2);

  // A monster zooms right beside itself, not over it.
  const monster = active.getByRole('group', { name: 'Monsters' }).locator('.card').first();
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

test('the zoom shows keyword tips beside the card (15.27)', async ({ page }) => {
  await page.goto(gallery);
  const zoom = page.getByTestId('card-zoom');
  const row = page.getByLabel('Table size');
  await row.getByRole('img', { name: 'The D6' }).hover();
  const tips = zoom.getByRole('note', { name: 'Keywords' });
  await expect(tips).toContainText('Eternal');
  await expect(tips).toContainText("This can't be destroyed or put into discard.");
  await shot(page, 'zoom-tips');
  // A card without keywords has no tips.
  await row.getByRole('img', { name: 'Isaac' }).hover();
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', 'Isaac');
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
  await expect(page.getByRole('group', { name: 'Your hand' }).locator('.card').first()).toBeVisible();
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
  const cards = page.getByRole('group', { name: 'Your hand' }).locator('.hand-card');
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
