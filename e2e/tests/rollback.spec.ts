// Rollback in Godmode (17.3, ADR 020): any player undoes a step from
// the history, the author too; the table jumps back to before it.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { saveAndContinue } from './saves';
import { activePage, answerPlay, canEndTurn, click, playLoot, settle, throwButton, turnAction } from './turn';

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
  await toThrow(pages, active);
}

// toThrow plays on until the attacker's Throw button waits.
async function toThrow(pages: Page[], active: Page) {
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
}

// undo is the History's one Undo button (17.7): "Undo", or "Ask to
// undo" for the Soft author.
function undo(page: Page) {
  return history(page).getByRole('button', { name: /^(Undo|Ask to undo)$/ });
}

// row is the history row of a line; who reads "You" on its own page.
function row(page: Page, text: RegExp) {
  return history(page).getByRole('listitem').filter({ hasText: text });
}

// The attack's two moves: its declaration, then its target.
const declared = /^⚔(You declare|.+ declares) an attack$/;
const aimed = /^⚔(You attack|.+ attacks) /;

// expectUndone: the target pick is undone on every page. The throw is
// gone, the target row is greyed, and a row says who rolled back.
async function expectUndone(pages: Page[], active: Page, who: string[]) {
  await expect(throwButton(active)).toHaveCount(0);
  for (const [i, page] of pages.entries()) {
    await expect(row(page, aimed).first()).toHaveClass(/undone/);
    await expect(row(page, declared).first()).not.toHaveClass(/undone/);
    await expect(row(page, new RegExp(`^↶${who[i] ?? ''} back \\d+ steps?$`))).toBeVisible();
  }
}

test('in Godmode each Undo undoes one move, by the author or another player (17.3, 17.7)', async ({ browser }) => {
  // Animations on: the restored table must jump, with no clash.
  const names = ['Gia', 'Hal'];
  const pages = await startGame(browser, names, { options: god });
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
  await expect(history(active).getByText('Godmode')).toBeVisible();
  // The latest move is the target pick; the tooltip names it.
  await expect(undo(active)).toHaveText(/Undo$/);
  await expect(undo(active)).toHaveAttribute('title', /^Undo: You attack /);
  await expect(undo(other)).toHaveAttribute('title', new RegExp(`^Undo: ${name} attacks `));
  await undo(active).hover();
  await shot(active, 'rollback-undo');
  await undo(active).click();
  await expectUndone(pages, active, pages.map((p) => (p === active ? 'You roll' : `${name} rolls`)));
  await active.mouse.move(1, 1);
  await shot(active, 'rollback-undone');
  // A second press, by the other player, undoes the declaration too.
  await expect(undo(other)).toHaveAttribute('title', `Undo: ${name} declares an attack`);
  await undo(other).click();
  await expect(turnAction(active, 'Attack')).toBeEnabled();
  for (const page of pages) {
    await expect(row(page, declared).first()).toHaveClass(/undone/);
    await expect(row(page, /^↶(You roll|.+ rolls) back \d+ steps?$/)).toHaveCount(2);
  }
  await expect(row(active, new RegExp(`^↶${otherName} rolls back`))).toBeVisible();
  await shot(other, 'rollback-other');
  // The game goes on: the attack can be made again.
  await attackToThrow(pages, active);
});

test('a game without a rollback mode has no Undo button (ADR 020)', async ({ browser }) => {
  const pages = await startGame(browser, ['Kai', 'Lux'], { url: `${appURL()}&anim=off` });
  const active = await activePage(pages, click);
  await attackToThrow(pages, active);
  await openHistory(active);
  await expect(history(active).getByRole('listitem').filter({ hasText: 'You declare an attack' })).toBeVisible();
  await expect(history(active).getByRole('button', { name: /undo/i })).toHaveCount(0);
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
  const hint = host.getByText(/The other players vote to undo a move/);
  await expect(hint).toBeVisible();
  await hint.scrollIntoViewIfNeeded();
  await shot(host, 'host-rollback');
}

function votePanel(page: Page) {
  return page.getByRole('dialog', { name: 'Rollback vote' });
}

test('in Soft the author asks; another player opens a vote that passes (17.4, 17.7)', async ({ browser }) => {
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
  // The author has no Undo, only the ask; the others undo by a vote.
  await expect(undo(active)).toHaveText(/Ask to undo$/);
  await expect(undo(active)).toHaveAttribute('title', /^Ask to undo: You attack /);
  await expect(undo(opener)).toHaveText(/Undo$/);
  await expect(undo(opener)).toHaveAttribute('title', new RegExp(`^Undo: ${name} attacks `));
  await shot(active, 'rollback-ask');
  await undo(opener).click();
  // The game waits: everyone sees the vote, only the voter answers.
  for (const page of pages) {
    await expect(votePanel(page)).toBeVisible();
    await expect(votePanel(page)).toContainText('1 of 2 yes needed');
  }
  await expect(votePanel(active).getByRole('button')).toHaveCount(0);
  await expect(votePanel(opener).getByRole('button', { name: 'Cancel vote' })).toBeVisible();
  await expect(votePanel(voter)).toContainText(`${openerName} wants to undo ${name}'s move`);
  await expect(votePanel(voter)).toContainText(`“${name} attacks `);
  await expect(throwButton(active)).toHaveCount(0);
  // While the vote is open, Undo is greyed and says why.
  await expect(undo(voter)).toHaveAttribute('aria-disabled', 'true');
  await expect(undo(voter)).toHaveAttribute('title', 'A rollback vote is open');
  await shot(voter, 'rollback-vote');
  await shot(opener, 'rollback-vote-opener');
  await votePanel(voter).getByRole('button', { name: 'Yes, undo' }).click();
  await expectUndone(pages, active, pages.map((p) => (p === opener ? 'You roll' : `${openerName} rolls`)));
  for (const page of pages) {
    await expect(votePanel(page)).toHaveCount(0);
  }
  // The game goes on: the target is picked again.
  await toThrow(pages, active);
});

test('in Soft a declined ask leaves the move; the author may ask again (17.4)', async ({ browser }) => {
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
  const ask = undo(active);
  await expect(ask).toHaveText(/Ask to undo$/);
  await ask.click();
  await expect(votePanel(active)).toContainText('You ask to undo a move');
  await expect(votePanel(active).getByRole('button', { name: 'Yes, undo' })).toHaveCount(0);
  await expect(votePanel(active).getByRole('button', { name: 'Cancel vote' })).toBeVisible();
  await expect(votePanel(other)).toContainText(`${name} asks to undo a move`);
  await shot(active, 'rollback-ask-status');
  await votePanel(other).getByRole('button', { name: 'No', exact: true }).click();
  for (const page of pages) {
    await expect(votePanel(page)).toHaveCount(0);
    await expect(row(page, /The rollback was declined/)).toBeVisible();
  }
  // Nothing was undone: the throw still waits.
  await expect(throwButton(active)).toBeVisible();
  await expect(row(active, aimed).first()).not.toHaveClass(/undone/);
  // Asked again, and cancelled by the asker.
  await ask.click();
  await expect(votePanel(other)).toBeVisible();
  await votePanel(active).getByRole('button', { name: 'Cancel vote' }).click();
  await expect(row(other, new RegExp(`${name} cancelled the rollback vote`))).toBeVisible();
  await expect(throwButton(active)).toBeVisible();
});

test('in Soft with two players Undo rolls back at once; after a roll it is greyed (17.4, 17.7)', async ({ browser }) => {
  const names = ['Rex', 'Sia'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: { rollback: 'soft' } });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const otherName = names[pages.indexOf(other)] ?? '';
  await attackToThrow(pages, active);
  for (const page of pages) {
    await openHistory(page);
  }
  // The opponent is the only voter: their yes is the majority.
  await undo(other).click();
  await expectUndone(pages, active, pages.map((p) => (p === other ? 'You roll' : `${otherName} rolls`)));
  // A roll is a barrier: once thrown, nothing before it is undone.
  await toThrow(pages, active);
  await click(active, throwButton(active));
  await expect(throwButton(active)).toHaveCount(0);
  const greyed = /^Nothing to undo since (a die was rolled|damage was dealt|a card died)$/;
  for (const page of pages) {
    await expect(undo(page)).toHaveAttribute('aria-disabled', 'true');
    await expect(undo(page)).toHaveAttribute('title', greyed);
  }
  // A click on the greyed button says why, where it was clicked. It is
  // aria-disabled only, which Playwright counts as disabled.
  await undo(other).click({ force: true });
  await expect(other.getByRole('status').getByText(greyed)).toBeVisible();
  await shot(other, 'rollback-undo-greyed');
});

// giveAndPlay gives the active player a loot card by a saved game
// (15.25) and plays it, answering what it asks first. It returns the
// active page, the other, the card in the hand, and their names.
async function giveAndPlay(pages: Page[], names: string[], ref: string, card: string) {
  await activePage(pages, click);
  await saveAndContinue(pages, names.join(', '), ref);
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  // The deal may have given the card too: count the copies. A card is
  // a button only while it can be played, so count by name.
  const group = active.getByRole('group', { name: 'Your hand' });
  const hand = group.locator(`[aria-label="${card}"]`);
  const copies = await hand.count();
  await playLoot(active, group.getByRole('button', { name: card }).first(), click);
  await answerPlay(active, click);
  await expect(stack(active)).toContainText(`plays ${card}`);
  await expect(hand).toHaveCount(copies - 1);
  for (const page of pages) {
    await openHistory(page);
  }
  return { active, other, hand, copies, name: names[pages.indexOf(active)] ?? '', otherName: names[pages.indexOf(other)] ?? '' };
}

function stack(page: Page) {
  return page.getByRole('complementary', { name: 'Stack' });
}

// resolve passes until the stack is empty and the turn goes on.
async function resolve(pages: Page[], active: Page) {
  const rows = stack(active).locator('.stack-row');
  await settle(pages, click, async () => (await rows.count()) === 0 && (await canEndTurn(active)));
}

// expectPlayUndone: the card is back in the hand and the stack is empty
// on every page.
async function expectPlayUndone(pages: Page[], hand: Locator, copies: number) {
  await expect(hand).toHaveCount(copies);
  for (const page of pages) {
    await expect(stack(page).locator('.stack-row')).toHaveCount(0);
  }
}

test('in Soft a played loot card is undone by one Undo, on the stack and after (17.6, 17.7)', async ({ browser }) => {
  const names = ['Tom', 'Uma'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: { rollback: 'soft' } });
  const { active, other, hand, copies, name, otherName } = await giveAndPlay(pages, names, 'a_penny', 'A Penny!');
  // The play has no history row (the history stays clean); Undo names it.
  await expect(undo(other)).toHaveAttribute('title', `Undo: ${name} plays A Penny!`);
  await expect(undo(active)).toHaveAttribute('title', 'Ask to undo: You play A Penny!');
  await expect(row(other, /plays A Penny!/)).toHaveCount(0);
  await shot(other, 'rollback-play-stack');
  await resolve(pages, active);
  const resolved = row(other, /^✦A Penny!/);
  await expect(resolved).toBeVisible();
  // The result came with a pass, no target: Undo still names the play.
  await expect(undo(other)).toHaveAttribute('title', `Undo: ${name} plays A Penny!`);
  // Two players: the opponent's yes is the majority.
  await undo(other).click();
  await expectPlayUndone(pages, hand, copies);
  await expect(resolved).toHaveClass(/undone/);
  for (const page of pages) {
    const who = page === other ? 'You roll' : `${otherName} rolls`;
    await expect(row(page, new RegExp(`^↶${who} back \\d+ steps?$`))).toBeVisible();
  }
  await other.mouse.move(1, 1);
  await shot(other, 'rollback-play-undone');
});

test('in Godmode one Undo after a resolved card undoes the play, not only its result (19.69)', async ({ browser }) => {
  const names = ['Vic', 'Wes'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: god });
  const { active, hand, copies } = await giveAndPlay(pages, names, 'a_penny', 'A Penny!');
  await resolve(pages, active);
  await expect(undo(active)).toHaveAttribute('title', 'Undo: You play A Penny!');
  await undo(active).click();
  // The card is back in the hand, not on the stack.
  await expectPlayUndone(pages, hand, copies);
});

test('a played card that asks a target is undone by one Undo (17.7)', async ({ browser }) => {
  const names = ['Xan', 'Yul'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: { rollback: 'soft' } });
  // Bomb! asks its target first; the play and the pick are one move.
  const { active, other, hand, copies, name } = await giveAndPlay(pages, names, 'bomb', 'Bomb!');
  await expect(undo(other)).toHaveAttribute('title', `Undo: ${name} plays Bomb!`);
  await undo(other).click();
  // Back before the play, not at its target prompt.
  await expectPlayUndone(pages, hand, copies);
  await expect(active.locator('.card.hl-aim')).toHaveCount(0);
  await expect(turnAction(active, 'Attack')).toBeEnabled();
});

test("Bum Friend's put-back pick has no line of its own: Undo names it a choice (17.8)", async ({ browser }) => {
  const names = ['Zed', 'Ari'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: god });
  await activePage(pages, click);
  // The starting item becomes Bum Friend: loot 1, then put a hand card back.
  await saveAndContinue(pages, names.join(', '), 'a_penny', 'bum_friend');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const name = names[pages.indexOf(active)] ?? '';
  await click(active, active.getByRole('group', { name: 'Your items' }).getByRole('button', { name: 'Bum Friend' }));
  await click(active, active.getByRole('menu').getByRole('menuitem').first());
  for (const page of pages) {
    await openHistory(page);
  }
  await expect(undo(active)).toHaveAttribute('title', 'Undo: You use Bum Friend');
  // The pick puts a card on the loot deck, which gives no line; the
  // tooltip never borrows the next move's name.
  await settle(pages, click, async () => (await undo(active).getAttribute('title')) === 'Undo: your choice');
  await expect(undo(other)).toHaveAttribute('title', `Undo: ${name}'s choice`);
  await undo(active).hover();
  await shot(active, 'rollback-choice');
  await undo(other).click();
  // Back before the pick: the item's use is the latest move again.
  await expect(undo(active)).toHaveAttribute('title', 'Undo: You use Bum Friend');
});
