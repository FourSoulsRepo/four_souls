// Two-sided cards (14.15, ADR 023): a card shows the side that is up,
// each side by its own name, turns over when it flips, and its zoom
// shows the other side.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle } from './turn';

test.setTimeout(120_000);

// countFlips counts the card-flip animations the page starts from now.
async function countFlips(page: Page): Promise<void> {
  await page.evaluate(() => {
    const w = window as unknown as { flips: number };
    w.flips = 0;
    document.addEventListener('animationstart', (e) => {
      if (e.animationName === 'card-flip') {
        w.flips++;
      }
    });
  });
}

function flips(page: Page): Promise<number> {
  return page.evaluate(() => (window as unknown as { flips: number }).flips);
}

function history(page: Page) {
  return page.getByRole('complementary', { name: tr('history.title'), exact: true });
}

test('Anima Sola flips to The Revenant and back (14.15)', async ({ browser }) => {
  const names = ['Ivo', 'Jun'];
  // Godmode lets a test undo the flip: The Revenant itself only flips
  // back when attacked, and attacking an item is 14.14.
  // The pages offer Requiem (?requiem), so they can continue a game
  // with it; the deal stays the Base Game's.
  const pages = await startGame(browser, names, {
    url: `${appURL()}&requiem`,
    setup: (host) => host.getByRole('group', { name: tr('host.sets') }).getByLabel(tr('set.requiem')).uncheck(),
    options: { rollback: 'god' },
  });
  await activePage(pages, click);
  // The starting item becomes Anima Sola; its set joins the game.
  await saveAndContinue(pages, names.join(', '), 'a_penny', 'anima_sola', false, { r: 'Requiem' });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const name = names[pages.indexOf(active)] ?? '';
  const mine = active.getByRole('group', { name: tr('table.your_items') });
  const theirs = other.getByRole('group', { name: tr('table.items_of', { player: name }) });
  // Each side has its own name, never "Anima Sola/The Revenant".
  await expect(mine.getByLabel(cardName('anima_sola'), { exact: true })).toBeVisible();
  await expect(theirs.getByLabel(cardName('anima_sola'), { exact: true })).toBeVisible();
  for (const page of pages) {
    await countFlips(page);
  }

  // ↷: Flip this item.
  await click(active, mine.getByRole('button', { name: cardName('anima_sola'), exact: true }));
  await click(active, active.getByRole('menu').getByRole('menuitem').first());
  const revenant = mine.getByLabel(cardName('the_revenant'), { exact: true });
  await settle(pages, click, async () => (await revenant.count()) > 0);
  await expect(theirs.getByLabel(cardName('the_revenant'), { exact: true })).toBeVisible();
  await expect(mine.getByLabel(cardName('anima_sola'), { exact: true })).toHaveCount(0);
  // Both tables turned the card over.
  for (const page of pages) {
    await expect.poll(() => flips(page)).toBeGreaterThan(0);
  }

  // The zoom shows the side that is up, and the other one beside it.
  await revenant.hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.getByRole('img', { name: cardName('the_revenant'), exact: true })).toBeVisible();
  const side = zoom.locator('.card-zoom-side');
  await expect(side).toContainText(tr('zoom.other_side'));
  await expect(side.getByRole('img', { name: cardName('anima_sola'), exact: true })).toBeVisible();
  await shot(active, 'two-sided-zoom');
  await active.mouse.move(1, 1);

  // The history names the side that came up.
  await active.getByRole('button', { name: tr('history.title'), exact: true }).click();
  await expect(history(active).getByText(tr('history.flipped', { card: cardName('the_revenant') }))).toBeVisible();

  // Undo turns it back: Anima Sola again, on both tables.
  await history(active).getByRole('button', { name: tr('history.undo.undo') }).click();
  await expect(mine.getByLabel(cardName('anima_sola'), { exact: true })).toBeVisible();
  await expect(theirs.getByLabel(cardName('anima_sola'), { exact: true })).toBeVisible();
  await expect(mine.getByLabel(cardName('the_revenant'), { exact: true })).toHaveCount(0);
});
