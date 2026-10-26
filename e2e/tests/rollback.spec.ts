// Rollback in Godmode (17.3, ADR 020): any player undoes a step from
// the history, the author too; the table jumps back to before it.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle, throwButton, turnAction } from './turn';

test.setTimeout(120_000);

const god = { rollback: 'god' };

function history(page: Page) {
  return page.getByRole('complementary', { name: 'History', exact: true });
}

async function openHistory(page: Page) {
  await page.getByRole('button', { name: 'History', exact: true }).click();
  await expect(history(page)).toBeVisible();
}

// attackToThrow declares an attack and picks its target, until the
// attacker's Throw button waits (as in fight.spec.ts).
async function attackToThrow(pages: Page[], active: Page) {
  await expect(turnAction(active, 'Attack')).toBeEnabled();
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
}

// expectRolledBack: the attack is undone on every page. The throw is
// gone, the attack is there to use again, the declare row is greyed,
// and a row says who rolled back.
async function expectRolledBack(pages: Page[], active: Page, who: string[]) {
  await expect(throwButton(active)).toHaveCount(0);
  await expect(turnAction(active, 'Attack')).toBeEnabled();
  for (const [i, page] of pages.entries()) {
    await expect(history(page).getByRole('listitem').filter({ hasText: /declares? an attack/ }).first()).toHaveClass(/undone/);
    await expect(history(page).getByRole('listitem').filter({ hasText: new RegExp(`^↶${who[i] ?? ''} back \\d+ steps?$`) })).toBeVisible();
  }
}

test('in Godmode the author rolls back their own step (17.3)', async ({ browser }) => {
  // Animations on: the restored table must jump, with no clash.
  const names = ['Gia', 'Hal'];
  const pages = await startGame(browser, names, { options: god });
  const active = await activePage(pages, click);
  const name = names[pages.indexOf(active)] ?? '';
  await attackToThrow(pages, active);
  for (const page of pages) {
    await openHistory(page);
  }
  await expect(history(active).getByText('Godmode')).toBeVisible();
  const button = history(active).getByRole('button', { name: 'Roll back: You declare an attack' });
  await expect(button).toBeEnabled();
  await button.hover();
  await shot(active, 'rollback-button');
  await button.click();
  await expectRolledBack(pages, active, pages.map((p) => (p === active ? 'You roll' : `${name} rolls`)));
  await active.mouse.move(1, 1);
  await shot(active, 'rollback-undone');
  // Undone rows offer no rollback any more.
  await expect(history(active).getByRole('button', { name: 'Roll back: You declare an attack' })).toHaveCount(0);
});

test('in Godmode another player rolls back (17.3)', async ({ browser }) => {
  const names = ['Ivo', 'Jen'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: god });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const name = names[pages.indexOf(active)] ?? '';
  const otherName = names[pages.indexOf(other)] ?? '';
  await attackToThrow(pages, active);
  for (const page of pages) {
    await openHistory(page);
  }
  await history(other).getByRole('button', { name: `Roll back: ${name} declares an attack` }).click();
  await expectRolledBack(pages, active, pages.map((p) => (p === other ? 'You roll' : `${otherName} rolls`)));
  await shot(other, 'rollback-other');
  // The game goes on: the attack can be made again.
  await attackToThrow(pages, active);
});

test('a game without a rollback mode has no rollback buttons (ADR 020)', async ({ browser }) => {
  const pages = await startGame(browser, ['Kai', 'Lux'], { url: `${appURL()}&anim=off` });
  const active = await activePage(pages, click);
  await attackToThrow(pages, active);
  await openHistory(active);
  await expect(history(active).getByRole('listitem').filter({ hasText: 'You declare an attack' })).toBeVisible();
  await expect(history(active).getByRole('button', { name: /^Roll back/ })).toHaveCount(0);
  await expect(history(active).getByText('Godmode')).toHaveCount(0);
});
