// Stack and responses (7.8): players always see why the game waits.
import { expect, test, type Page } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

test('a played loot card shows on the stack; others may skip all', async ({ browser }) => {
  const pages = await startGame(browser, ['Sam', 'Tia'], { url: fast() });
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
    url: fast(),
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

// waitLine is a stack panel line with "You" turned into the page's own
// name, so every page's line can be compared.
async function waitLine(page: Page, me: string, label: string): Promise<string> {
  const line = page.getByRole('complementary', { name: 'Stack' }).getByLabel(label, { exact: true });
  if ((await line.count()) === 0) {
    return '';
  }
  return (await line.first().innerText()).replace(/\bYou\b/g, me);
}

// expectSameLines: every page names the same waiting and passed players.
async function expectSameLines(pages: Page[], names: string[], want: string): Promise<void> {
  for (const [i, page] of pages.entries()) {
    const me = names[i] ?? '';
    await expect.poll(async () => (await waitLine(page, me, 'Waiting for')) + (await waitLine(page, me, 'Passed'))).toBe(want);
  }
}

test('every player sees whom the stack waits for (V-11)', async ({ browser }) => {
  const names = ['Abe', 'Bea', 'Cal', 'Dot'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const me = names[pages.indexOf(active)] ?? '';
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  await loot.click();
  await active.getByRole('menu').getByRole('menuitem', { name: 'Play' }).click();
  const stack = active.getByRole('complementary', { name: 'Stack' });
  for (let i = 0; i < 10 && !(await stack.isVisible()); i++) {
    const option = active.locator('.choice .choice-option button').first();
    if (await option.isVisible()) {
      await option.click();
    }
    await active.waitForTimeout(200);
  }
  await expect(stack).toContainText(/Waiting for|: /);
  const want = (await waitLine(active, me, 'Waiting for')) + (await waitLine(active, me, 'Passed'));
  expect(want).not.toBe('');
  await expectSameLines(pages, names, want);
  await shot(active, 'stack-waiting');

  // A pass moves the player to "Passed" on every table.
  const pass = active.getByRole('region', { name: 'Actions' }).getByRole('button', { name: 'Pass' });
  test.skip(!(await pass.isEnabled()), 'the loot card asked something first');
  await pass.click();
  // Unless everyone else was skipped and the item resolved.
  await expect.poll(async () => !(await stack.isVisible()) || (await waitLine(active, me, 'Passed')).includes(me)).toBe(true);
  await expectSameLines(pages, names, (await waitLine(active, me, 'Waiting for')) + (await waitLine(active, me, 'Passed')));
});
