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

// pickSoft picks the Soft mode in the host's Detailed setup (17.4): the
// control is Disabled by default and says what each mode allows.
async function pickSoft(host: Page) {
  await host.getByLabel('Detailed setup').check();
  const mode = host.getByLabel('Rollback');
  await expect(mode).toHaveValue('');
  await expect(host.getByText('No rollback: every step stands.')).toBeVisible();
  await mode.selectOption('soft');
  const hint = host.getByText(/The other players vote to undo a step/);
  await expect(hint).toBeVisible();
  await hint.scrollIntoViewIfNeeded();
  await shot(host, 'host-rollback');
}

function votePanel(page: Page) {
  return page.getByRole('dialog', { name: 'Rollback vote' });
}

test('in Soft the author asks; another player opens a vote that passes (17.4)', async ({ browser }) => {
  const names = ['Mia', 'Ned', 'Oda'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, setup: pickSoft });
  const active = await activePage(pages, click);
  const name = names[pages.indexOf(active)] ?? '';
  const [opener, voter] = pages.filter((p) => p !== active);
  if (!opener || !voter) {
    throw new Error('no other pages');
  }
  const openerName = names[pages.indexOf(opener)] ?? '';
  await attackToThrow(pages, active);
  for (const page of pages) {
    await openHistory(page);
  }
  await expect(history(active).getByText('Soft rollback')).toBeVisible();
  // The author has no direct button, only the ask.
  await expect(history(active).getByRole('button', { name: /^Roll back/ })).toHaveCount(0);
  await expect(history(active).getByRole('button', { name: 'Ask to undo: You declare an attack' })).toBeVisible();
  await history(opener).getByRole('button', { name: `Roll back: ${name} declares an attack` }).click();
  // The game waits: everyone sees the vote, only the voter answers.
  for (const page of pages) {
    await expect(votePanel(page)).toBeVisible();
    await expect(votePanel(page)).toContainText('1 of 2 yes needed');
  }
  await expect(votePanel(active).getByRole('button')).toHaveCount(0);
  await expect(votePanel(opener).getByRole('button', { name: 'Cancel vote' })).toBeVisible();
  await expect(votePanel(voter)).toContainText(`${openerName} wants to undo ${name}'s step`);
  await expect(throwButton(active)).toHaveCount(0);
  await shot(voter, 'rollback-vote');
  await shot(opener, 'rollback-vote-opener');
  await votePanel(voter).getByRole('button', { name: 'Yes, undo' }).click();
  await expectRolledBack(pages, active, pages.map((p) => (p === opener ? 'You roll' : `${openerName} rolls`)));
  for (const page of pages) {
    await expect(votePanel(page)).toHaveCount(0);
  }
  // The game goes on.
  await attackToThrow(pages, active);
});

test('in Soft a declined ask leaves the step; the author may ask again (17.4)', async ({ browser }) => {
  const names = ['Pia', 'Quin'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, setup: pickSoft });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const name = names[pages.indexOf(active)] ?? '';
  await attackToThrow(pages, active);
  for (const page of pages) {
    await openHistory(page);
  }
  const ask = history(active).getByRole('button', { name: 'Ask to undo: You declare an attack' });
  await ask.click();
  await expect(votePanel(active)).toContainText('You ask to undo a step');
  await expect(votePanel(active).getByRole('button', { name: 'Yes, undo' })).toHaveCount(0);
  await expect(votePanel(active).getByRole('button', { name: 'Cancel vote' })).toBeVisible();
  await expect(votePanel(other)).toContainText(`${name} asks to undo a step`);
  await shot(active, 'rollback-ask-status');
  await votePanel(other).getByRole('button', { name: 'No', exact: true }).click();
  for (const page of pages) {
    await expect(votePanel(page)).toHaveCount(0);
    await expect(history(page).getByRole('listitem').filter({ hasText: 'The rollback was declined' })).toBeVisible();
  }
  // Nothing was undone: the throw still waits.
  await expect(throwButton(active)).toBeVisible();
  await expect(history(active).getByRole('listitem').filter({ hasText: 'You declare an attack' })).not.toHaveClass(/undone/);
  // Asked again, and cancelled by the asker.
  await ask.click();
  await expect(votePanel(other)).toBeVisible();
  await votePanel(active).getByRole('button', { name: 'Cancel vote' }).click();
  await expect(history(other).getByRole('listitem').filter({ hasText: `${name} cancelled the rollback vote` })).toBeVisible();
  await expect(throwButton(active)).toBeVisible();
});

test('in Soft with two players the other player rolls back at once (17.4)', async ({ browser }) => {
  const names = ['Rex', 'Sia'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: { rollback: 'soft' } });
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
  await expect(history(active).getByRole('button', { name: /^Roll back/ })).toHaveCount(0);
  // The opponent is the only voter: their yes is the majority.
  await history(other).getByRole('button', { name: `Roll back: ${name} declares an attack` }).click();
  await expectRolledBack(pages, active, pages.map((p) => (p === other ? 'You roll' : `${otherName} rolls`)));
  // A roll is a barrier: once thrown, the attack is no target any more.
  await attackToThrow(pages, active);
  await click(active, throwButton(active));
  await expect(throwButton(active)).toHaveCount(0);
  await expect(history(other).getByRole('button', { name: `Roll back: ${name} declares an attack` })).toHaveCount(0);
});
