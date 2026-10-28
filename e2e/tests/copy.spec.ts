// The copy mark (19.106): an item that uses another card's abilities
// shows a mark, and its tooltip names the card, for both players.
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
  // The other player sees the same mark.
  await expect.soft(other.locator('.card-copy')).toHaveCount(1);
});
