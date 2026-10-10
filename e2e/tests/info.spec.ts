// Info panels (7.9): history, turn badge, what is left this turn.
import { expect, test } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';

const fast = `${appURL}&anim=off`;

test.setTimeout(90_000);

test('the history shows only what the viewer may see', async ({ browser }) => {
  const names = ['Wes', 'Xan'];
  const pages = await startGame(browser, names, { url: fast });
  const active = await activePage(pages, click);
  const history = (p: typeof active) => p.getByRole('complementary', { name: 'History' });
  for (const [i, page] of pages.entries()) {
    const other = names.at(1 - i) ?? '';
    const text = (await history(page).textContent()) ?? '';
    // Others' loots never name the card.
    expect(text).not.toMatch(new RegExp(`${other} loots (?!a card)`));
  }
  // The active player looted: their own history names the card.
  await expect(history(active)).toContainText(/You loot \S/);
  const other = pages.find((p) => p !== active);
  const activeName = names[pages.indexOf(active)] ?? '';
  await expect(history(other ?? active)).toContainText(`${activeName} loots a card`);
  await shot(active, 'history');
});

test('the turn badge, what is left, and hiding the history', async ({ browser }) => {
  const names = ['Yul', 'Zed'];
  const pages = await startGame(browser, names, { url: fast });
  const active = await activePage(pages, click);
  const activeName = names[pages.indexOf(active)] ?? '';
  for (const page of pages) {
    await expect(page.getByRole('status', { name: 'Turn' })).toHaveText(new RegExp(`^Turn \\d+ · ${activeName} · \\w`));
  }
  await expect(active.getByLabel('Left this turn')).toHaveText(/Loot \d+ · Attacks \d+ · Buys \d+/);

  await active.getByRole('complementary', { name: 'History' }).getByRole('button', { name: 'Hide' }).click();
  await expect(active.getByRole('complementary', { name: 'History' })).toHaveCount(0);
  await active.getByRole('button', { name: 'History' }).click();
  await expect(active.getByRole('complementary', { name: 'History' })).toBeVisible();
});
