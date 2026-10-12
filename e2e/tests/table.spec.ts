// Table layout (7.5), with real games on the Go server.
import { expect, test } from '@playwright/test';
import { expectFits, startGame } from './game';
import { shot } from './shot';

test('two players face each other; the table fits 1280x720', async ({ browser }) => {
  const [eve, fay] = await startGame(browser, ['Eve', 'Fay']);
  if (!eve || !fay) {
    throw new Error('no pages');
  }
  for (const [page, me, other] of [
    [eve, 'Eve', 'Fay'],
    [fay, 'Fay', 'Eve'],
  ] as const) {
    await expect(page.locator('.player-bottom')).toHaveAccessibleName(`${me} (you)`);
    await expect(page.locator('.player-top')).toHaveAccessibleName(other);
    await expect(page.getByRole('group', { name: 'Your hand' }).getByRole('img')).not.toHaveCount(0);
    await expect(page.getByRole('group', { name: /^Loot deck: \d+ cards$/ })).toBeVisible();
    // The opponent's hand is one loot back with its count (15.11).
    const theirHand = page.locator('.player-top').getByRole('group', { name: /^Hand: \d+ cards$/ });
    await expect(theirHand.getByRole('img', { name: 'Face-down loot card' })).toHaveCount(1);
    await expect(theirHand.locator('.pile-count')).toHaveText(/^\d+$/);
    await expect(page.getByRole('group', { name: 'Shop' }).getByRole('img')).not.toHaveCount(0);
    await expect(page.getByRole('group', { name: 'Monsters' }).getByRole('img')).not.toHaveCount(0);
    await expectFits(page);
  }
  await shot(eve, 'table-2p');
});

test('four players sit left, top and right in turn order', async ({ browser }) => {
  const pages = await startGame(browser, ['Gus', 'Hal', 'Ivy', 'Jo']);
  const hal = pages[1];
  if (!hal) {
    throw new Error('no pages');
  }
  await expect(hal.locator('.player-bottom')).toHaveAccessibleName('Hal (you)');
  await expect(hal.locator('.player-left')).toHaveAccessibleName('Ivy');
  await expect(hal.locator('.player-top')).toHaveAccessibleName('Jo');
  await expect(hal.locator('.player-right')).toHaveAccessibleName('Gus');
  for (const page of pages) {
    await expectFits(page);
  }
  await shot(hal, 'table-4p');
});
