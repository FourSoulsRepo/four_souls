// Stat sources (15.92, V-41): a changed stat shows in color, its sum in
// the zoom, and each effect of this turn as an icon on the card.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

// stackDone: the stack is empty for the active player.
const stackDone = (page: Page) => async () => (await page.getByRole('complementary', { name: 'Stack' }).locator('.stack-row').count()) === 0;

test('a lasting item changes the badge; XI. Strength adds an icon (15.92)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ada', 'Ben'], { url: fast() });
  await activePage(pages, click);
  // The starting item becomes Brimstone: a lasting +1 ATK.
  await saveAndContinue(pages, 'Ada, Ben', 'xi_strength', 'brimstone');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const own = active.locator('.own-character .card');
  // A lasting item: a green ATK badge, no icon (owner, 2026-10-22).
  await expect.soft(own.locator('.card-atk.up')).toBeVisible();
  await expect.soft(own.locator('.card-source')).toHaveCount(0);
  await own.hover();
  await expect.soft(active.getByTestId('card-zoom').locator('.card-zoom-stats')).toContainText('+ 1 Brimstone');
  await shot(active, 'stats-lasting');

  // XI. Strength on the player themself: +1 ATK this turn.
  await click(active, active.getByRole('group', { name: 'Your hand' }).getByRole('button', { name: 'XI. Strength' }).first());
  await click(active, active.locator('.own-character .card.hl-aim'));
  // Wait for the play to reach the stack, then resolve it.
  await expect(active.getByRole('complementary', { name: 'Stack' })).toContainText('plays XI. Strength');
  await settle(pages, click, stackDone(active));
  await active.mouse.move(0, 0);
  // One icon, on the character for both players, and the sum names it.
  await expect.soft(own.locator('.card-source')).toHaveCount(1);
  await expect.soft(own.locator('.card-source')).toHaveAttribute('aria-label', 'XI. Strength: +1 ATK this turn');
  await expect.soft(other.locator('.opp-character .card .card-source')).toHaveCount(1);
  await shot(active, 'stats-temporary');
  await own.hover();
  await expect.soft(active.getByTestId('card-zoom').locator('.card-zoom-stats')).toContainText('+ 1 Brimstone + 1 XI. Strength (this turn)');
  await shot(active, 'stats-sum');
  // Hovering the icon zooms XI. Strength, saying how long it lasts.
  await own.locator('.card-source').first().hover();
  await expect.soft(active.getByTestId('card-zoom')).toContainText('+1 ATK this turn');
  await shot(active, 'stats-icon');
});

test('a monster shows its changed dice and why (15.92)', async ({ browser }) => {
  const pages = await startGame(browser, ['Cy', 'Dee'], { url: fast() });
  await activePage(pages, click);
  // Baby Haunt: monsters have +1 DC on its owner's turn.
  await saveAndContinue(pages, 'Cy, Dee', 'bomb', 'baby_haunt');
  const active = await activePage(pages, click);
  const monster = active.locator('.board-slots .card').filter({ has: active.locator('.card-dc') }).first();
  await expect.soft(monster.locator('.card-dc.up')).toBeVisible();
  await monster.hover();
  await expect.soft(active.getByTestId('card-zoom').locator('.card-zoom-stats')).toContainText('+ 1 Baby Haunt');
  await shot(active, 'stats-monster');
});
