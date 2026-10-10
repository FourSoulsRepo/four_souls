// Stack and responses (7.8): players always see why the game waits.
import { expect, test } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';

const fast = `${appURL}&anim=off`;

test.setTimeout(90_000);

test('a played loot card shows on the stack; others may skip all', async ({ browser }) => {
  const pages = await startGame(browser, ['Sam', 'Tia'], { url: fast });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  const name = await loot.getAttribute('aria-label');
  await loot.click();
  await active.getByRole('menu').getByRole('menuitem', { name: 'Play' }).click();
  // Some loot cards ask for a target before they go on the stack.
  const stack = active.getByRole('complementary', { name: 'Stack' });
  for (let i = 0; i < 10 && !(await stack.isVisible()); i++) {
    const option = active.locator('.choice .choice-option button').first();
    if (await option.isVisible()) {
      await option.click();
    }
    await active.waitForTimeout(200);
  }

  for (const page of [active, other]) {
    await expect(page.getByRole('complementary', { name: 'Stack' })).toContainText(`plays ${name ?? ''}`);
  }
  // Both always see why the game waits.
  await expect(active.locator('.action-status')).toHaveText(/Respond, or pass|is choosing|Choose/);
  await expect(other.locator('.action-status')).toHaveText(/(Sam|Tia) (may respond|is choosing)/);
  await shot(other, 'stack');

  const bar = other.getByRole('region', { name: 'Actions' });
  await bar.getByRole('button', { name: 'Skip all' }).click();
  await expect(bar.getByRole('button', { name: 'Cancel skip' })).toBeVisible();
  await bar.getByRole('button', { name: 'Cancel skip' }).click();
  await expect(bar.getByRole('button', { name: 'Skip all' })).toBeVisible();
});

test('the response timer counts down for the waiting player', async ({ browser }) => {
  const pages = await startGame(browser, ['Uma', 'Vic'], {
    url: fast,
    setup: async (host) => {
      await host.getByLabel('Detailed setup').check();
      await host.getByLabel('Response timer').selectOption('60');
    },
  });
  const active = await activePage(pages, click);
  const timer = active.getByLabel('Time left');
  await expect(timer).toHaveText(/· \d+ s/);
  const first = Number((await timer.textContent())?.match(/\d+/)?.[0]);
  expect(first).toBeGreaterThan(0);
  await expect.poll(async () => Number((await timer.textContent())?.match(/\d+/)?.[0])).toBeLessThan(first);
});
