// Card component (7.4) in the dev-only card gallery.
import { expect, test } from '@playwright/test';
import { shot } from './shot';

const gallery = '/?fake&nosplash&screen=cards';

test('cards show their image, or a readable blank card', async ({ page }) => {
  await page.goto(gallery);
  const hand = page.getByLabel('Hand size');
  // Isaac has an image.
  await expect(hand.getByRole('img', { name: 'Isaac' }).locator('img')).toBeVisible();
  // The D6's image fails to load: it falls back to its text.
  await expect(hand.getByRole('img', { name: 'The D6' })).toContainText('Eternal Treasure Card');
  // Big Spider has no image: name, stats and reward.
  const spider = hand.getByRole('img', { name: 'Big Spider' });
  await expect(spider).toContainText('HP 3');
  await expect(spider).toContainText('1x Loot');
  await expect(hand.getByRole('img', { name: 'Face-down treasure card' })).toBeVisible();
  await shot(page, 'cards');
});

test('hovering a card zooms it with its full text', async ({ page }) => {
  await page.goto(gallery);
  const zoom = page.getByTestId('card-zoom');
  await page.getByLabel('Table size').getByRole('img', { name: 'Big Spider' }).hover();
  await expect(zoom).toContainText('When this dies, the active player may attack');
  await shot(page, 'card-zoom');
  await page.getByRole('heading', { name: 'Cards' }).hover();
  await expect(zoom).toBeHidden();
});

test('a long press zooms a card on touch; a tap closes it', async ({ page }) => {
  await page.goto(gallery);
  const zoom = page.getByTestId('card-zoom');
  const d6 = page.getByLabel('Table size').getByRole('img', { name: 'The D6' });
  await d6.dispatchEvent('pointerdown', { pointerType: 'touch', clientX: 100, clientY: 100 });
  await expect(zoom).toContainText("This can't be destroyed");
  await zoom.dispatchEvent('pointerdown', { pointerType: 'touch' });
  await expect(zoom).toBeHidden();
  // A short tap does not zoom.
  await d6.dispatchEvent('pointerdown', { pointerType: 'touch', clientX: 100, clientY: 100 });
  await d6.dispatchEvent('pointerup', { pointerType: 'touch' });
  await page.waitForTimeout(600);
  await expect(zoom).toBeHidden();
});
