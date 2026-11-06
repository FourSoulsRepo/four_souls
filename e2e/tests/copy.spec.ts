// The copy mark (19.106): an item that uses another card's abilities
// shows a mark with the copied card's art, for both players. Hover on
// the mark zooms the copied card; the item's zoom names it (19.111).
import { expect, test } from './fixtures';
import { startGame } from './game';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';
import { cardName, tr } from './locale';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

test('an item that gains another item\'s abilities shows the copy mark (19.106)', async ({ browser }) => {
  const pages = await startGame(browser, ['Nia', 'Oz'], { url: fast() });
  await activePage(pages, click);
  // Brimstone gains The D6's abilities till end of turn.
  await saveAndContinue(pages, 'Nia, Oz', 'bomb', 'brimstone', false, {}, 'the_d6');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  // One mark, on the item, named for the card it copies.
  const mark = active.locator('.card-copy');
  await expect.soft(mark).toHaveCount(1);
  await expect.soft(mark).toHaveAttribute('aria-label', tr('badge.copy', { card: cardName('the_d6') }));
  await shot(active, 'copy-badge');
  // The mark takes the pointer: hover zooms the copied card (19.111).
  const zoom = active.getByTestId('card-zoom');
  await mark.hover();
  await expect.soft(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', cardName('the_d6'));
  await shot(active, 'copy-badge-zoom');
  // The item's own zoom names the copied card first.
  const item = active.locator('.card:has(.card-copy)');
  await item.hover({ position: { x: 10, y: 80 } });
  await expect.soft(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', cardName('brimstone'));
  await expect.soft(zoom.locator('.card-zoom-stats li').first()).toHaveText(tr('badge.copy', { card: cardName('the_d6') }));
  await shot(active, 'copy-card-zoom');
  // The other player sees the same mark.
  await expect.soft(other.locator('.card-copy')).toHaveCount(1);
});
