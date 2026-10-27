// Rollback in Godmode (17.3, ADR 020): any player undoes a step from
// the history, the author too; the table jumps back to before it.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { saveAndContinue } from './saves';
import { activePage, answerPlay, canEndTurn, click, playLoot, settle, throwButton, turnAction } from './turn';
import { cardName, lineRe, tr, trRe } from './locale';

test.setTimeout(120_000);

const god = { rollback: 'god' };

function history(page: Page) {
  return page.getByRole('complementary', { name: tr('history.title'), exact: true });
}

async function openHistory(page: Page) {
  await page.getByRole('button', { name: tr('history.title'), exact: true }).click();
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
  return history(page).getByRole('button', { name: new RegExp(`^(${trRe('history.undo.undo').source}|${trRe('history.undo.ask').source})$`) });
}

// row is the history row of a line; who reads "You" on its own page.
function row(page: Page, text: RegExp) {
  return history(page).getByRole('listitem').filter({ hasText: text });
}

// The attack's two moves: its declaration, then its target.
const declared = new RegExp(`^⚔(${trRe('history.attack_declared.you').source}|${trRe('history.attack_declared.other').source})$`);
const aimed = new RegExp(`^⚔(${trRe('history.attack.you').source}|${trRe('history.attack.other').source})`);

// rolledBack is the row that says who rolled back: "You roll back 2
// steps" on their own page (you), else by name.
function rolledBack(you: boolean, name: string): RegExp {
  const line = you ? trRe('history.rolled_back.you') : trRe('history.rolled_back.other', { player: name });
  return new RegExp(`^↶(${line.source})$`);
}

// expectUndone: the target pick is undone on every page. The throw is
// gone, the target row is greyed, and a row says who rolled back.
async function expectUndone(pages: Page[], active: Page, who: RegExp[]) {
  await expect(throwButton(active)).toHaveCount(0);
  for (const [i, page] of pages.entries()) {
    await expect(row(page, aimed).first()).toHaveClass(/undone/);
    await expect(row(page, declared).first()).not.toHaveClass(/undone/);
    await expect(row(page, who[i] ?? /^$/)).toBeVisible();
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
  await expect(history(active).getByText(tr('host.rollback_god'))).toBeVisible();
  // The latest move is the target pick; the tooltip names it.
  await expect(undo(active)).toHaveText(new RegExp(`${trRe('history.undo.undo').source}$`));
  await expect(undo(active)).toHaveAttribute('title', new RegExp(`^${lineRe(tr('history.undo.title', { move: tr('history.attack.you') })).source}$`));
  await expect(undo(other)).toHaveAttribute('title', new RegExp(`^${lineRe(tr('history.undo.title', { move: tr('history.attack.other', { player: name }) })).source}$`));
  await undo(active).hover();
  await shot(active, 'rollback-undo');
  await undo(active).click();
  await expectUndone(pages, active, pages.map((p) => rolledBack(p === active, name)));
  await active.mouse.move(1, 1);
  await shot(active, 'rollback-undone');
  // A second press, by the other player, undoes the declaration too.
  await expect(undo(other)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.attack_declared.other', { player: name }) }));
  await undo(other).click();
  await expect(turnAction(active, 'Attack')).toBeEnabled();
  for (const page of pages) {
    await expect(row(page, declared).first()).toHaveClass(/undone/);
    await expect(row(page, new RegExp(`^↶(${trRe('history.rolled_back.you').source}|${trRe('history.rolled_back.other').source})$`))).toHaveCount(2);
  }
  await expect(row(active, rolledBack(false, otherName))).toBeVisible();
  await shot(other, 'rollback-other');
  // The game goes on: the attack can be made again.
  await attackToThrow(pages, active);
});

test('a game without a rollback mode has no Undo button (ADR 020)', async ({ browser }) => {
  const pages = await startGame(browser, ['Kai', 'Lux'], { url: `${appURL()}&anim=off` });
  const active = await activePage(pages, click);
  await attackToThrow(pages, active);
  await openHistory(active);
  await expect(history(active).getByRole('listitem').filter({ hasText: tr('history.attack_declared.you') })).toBeVisible();
  await expect(history(active).getByRole('button', { name: new RegExp(`${trRe('history.undo.undo').source}|${trRe('history.undo.ask').source}`, 'i') })).toHaveCount(0);
  await expect(history(active).getByText(tr('host.rollback_god'))).toHaveCount(0);
});

// pickSoft picks the Soft mode in the host's Detailed setup (17.4): the
// control is Disabled by default and says what each mode allows.
async function pickSoft(host: Page) {
  await host.getByLabel(tr('host.detailed')).check();
  const mode = host.getByLabel(tr('host.rollback'));
  await expect(mode).toHaveValue('');
  await expect(host.getByText(tr('host.rollback_hint.off'))).toBeVisible();
  await mode.selectOption('soft');
  const hint = host.getByText(tr('host.rollback_hint.soft'));
  await expect(hint).toBeVisible();
  await hint.scrollIntoViewIfNeeded();
  await shot(host, 'host-rollback');
}

function votePanel(page: Page) {
  return page.getByRole('dialog', { name: tr('rollback.label') });
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
  await expect(history(active).getByText(tr('history.mode_soft'))).toBeVisible();
  // The author has no Undo, only the ask; the others undo by a vote.
  await expect(undo(active)).toHaveText(new RegExp(`${trRe('history.undo.ask').source}$`));
  await expect(undo(active)).toHaveAttribute('title', new RegExp(`^${lineRe(tr('history.undo.ask_title', { move: tr('history.attack.you') })).source}$`));
  await expect(undo(opener)).toHaveText(new RegExp(`${trRe('history.undo.undo').source}$`));
  await expect(undo(opener)).toHaveAttribute('title', new RegExp(`^${lineRe(tr('history.undo.title', { move: tr('history.attack.other', { player: name }) })).source}$`));
  await shot(active, 'rollback-ask');
  await undo(opener).click();
  // The game waits: everyone sees the vote, only the voter answers.
  for (const page of pages) {
    await expect(votePanel(page)).toBeVisible();
    await expect(votePanel(page)).toContainText(tr('rollback.needed', { yes: 1, n: 2 }));
  }
  await expect(votePanel(active).getByRole('button')).toHaveCount(0);
  await expect(votePanel(opener).getByRole('button', { name: tr('rollback.cancel') })).toBeVisible();
  await expect(votePanel(voter)).toContainText(tr('rollback.want.other', { player: openerName, author: name }));
  await expect(votePanel(voter)).toContainText(trRe('rollback.undoes', { move: tr('history.attack.other', { player: name }) }));
  await expect(throwButton(active)).toHaveCount(0);
  // While the vote is open, Undo is greyed and says why.
  await expect(undo(voter)).toHaveAttribute('aria-disabled', 'true');
  await expect(undo(voter)).toHaveAttribute('title', tr('history.undo.voting'));
  await shot(voter, 'rollback-vote');
  await shot(opener, 'rollback-vote-opener');
  await votePanel(voter).getByRole('button', { name: tr('rollback.yes') }).click();
  await expectUndone(pages, active, pages.map((p) => rolledBack(p === opener, openerName)));
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
  await expect(ask).toHaveText(new RegExp(`${trRe('history.undo.ask').source}$`));
  await ask.click();
  await expect(votePanel(active)).toContainText(tr('rollback.ask.you'));
  await expect(votePanel(active).getByRole('button', { name: tr('rollback.yes') })).toHaveCount(0);
  await expect(votePanel(active).getByRole('button', { name: tr('rollback.cancel') })).toBeVisible();
  await expect(votePanel(other)).toContainText(tr('rollback.ask.other', { player: name }));
  await shot(active, 'rollback-ask-status');
  await votePanel(other).getByRole('button', { name: tr('rollback.no'), exact: true }).click();
  for (const page of pages) {
    await expect(votePanel(page)).toHaveCount(0);
    await expect(row(page, trRe('history.vote.declined'))).toBeVisible();
  }
  // Nothing was undone: the throw still waits.
  await expect(throwButton(active)).toBeVisible();
  await expect(row(active, aimed).first()).not.toHaveClass(/undone/);
  // Asked again, and cancelled by the asker.
  await ask.click();
  await expect(votePanel(other)).toBeVisible();
  await votePanel(active).getByRole('button', { name: tr('rollback.cancel') }).click();
  await expect(row(other, trRe('history.vote.cancelled.other', { player: name }))).toBeVisible();
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
  await expectUndone(pages, active, pages.map((p) => rolledBack(p === other, otherName)));
  // A roll is a barrier: once thrown, nothing before it is undone.
  await toThrow(pages, active);
  await click(active, throwButton(active));
  await expect(throwButton(active)).toHaveCount(0);
  const barriers = ['history.barrier.die', 'history.barrier.damage', 'history.barrier.died'].map((k) => trRe('history.undo.since', { barrier: tr(k) }).source);
  const greyed = new RegExp(`^(${barriers.join('|')})$`);
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
  const group = active.getByRole('group', { name: tr('table.your_hand') });
  const hand = group.locator(`[aria-label="${card}"]`);
  const copies = await hand.count();
  await playLoot(active, group.getByRole('button', { name: card }).first(), click);
  await answerPlay(active, click);
  await expect(stack(active)).toContainText(trRe('stack.loot', { card }));
  await expect(hand).toHaveCount(copies - 1);
  for (const page of pages) {
    await openHistory(page);
  }
  return { active, other, hand, copies, name: names[pages.indexOf(active)] ?? '', otherName: names[pages.indexOf(other)] ?? '' };
}

function stack(page: Page) {
  return page.getByRole('complementary', { name: tr('stack.title') });
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
  const { active, other, hand, copies, name, otherName } = await giveAndPlay(pages, names, 'a_penny', cardName('a_penny'));
  // The play has no history row (the history stays clean); Undo names it.
  await expect(undo(other)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.move.play.other', { player: name, card: cardName('a_penny') }) }));
  await expect(undo(active)).toHaveAttribute('title', tr('history.undo.ask_title', { move: tr('history.move.play.you', { card: cardName('a_penny') }) }));
  await expect(row(other, trRe('stack.loot', { card: cardName('a_penny') }))).toHaveCount(0);
  await shot(other, 'rollback-play-stack');
  await resolve(pages, active);
  const resolved = row(other, new RegExp(`^✦${lineRe(cardName('a_penny')).source}`));
  await expect(resolved).toBeVisible();
  // The result came with a pass, no target: Undo still names the play.
  await expect(undo(other)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.move.play.other', { player: name, card: cardName('a_penny') }) }));
  // Two players: the opponent's yes is the majority.
  await undo(other).click();
  await expectPlayUndone(pages, hand, copies);
  await expect(resolved).toHaveClass(/undone/);
  for (const page of pages) {
    await expect(row(page, rolledBack(page === other, otherName))).toBeVisible();
  }
  await other.mouse.move(1, 1);
  await shot(other, 'rollback-play-undone');
});

test('in Godmode one Undo after a resolved card undoes the play, not only its result (19.69)', async ({ browser }) => {
  const names = ['Vic', 'Wes'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: god });
  const { active, hand, copies } = await giveAndPlay(pages, names, 'a_penny', cardName('a_penny'));
  await resolve(pages, active);
  await expect(undo(active)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.move.play.you', { card: cardName('a_penny') }) }));
  await undo(active).click();
  // The card is back in the hand, not on the stack.
  await expectPlayUndone(pages, hand, copies);
});

test('a played card that asks a target is undone by one Undo (17.7)', async ({ browser }) => {
  const names = ['Xan', 'Yul'];
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off`, options: { rollback: 'soft' } });
  // Bomb! asks its target first; the play and the pick are one move.
  const { active, other, hand, copies, name } = await giveAndPlay(pages, names, 'bomb', cardName('bomb'));
  await expect(undo(other)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.move.play.other', { player: name, card: cardName('bomb') }) }));
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
  await click(active, active.getByRole('group', { name: tr('table.your_items') }).getByRole('button', { name: cardName('bum_friend') }));
  await click(active, active.getByRole('menu').getByRole('menuitem').first());
  for (const page of pages) {
    await openHistory(page);
  }
  await expect(undo(active)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.move.use.you', { card: cardName('bum_friend') }) }));
  // The pick puts a card on the loot deck, which gives no line; the
  // tooltip never borrows the next move's name.
  await settle(pages, click, async () => (await undo(active).getAttribute('title')) === tr('history.undo.title', { move: tr('history.choice.you') }));
  await expect(undo(other)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.choice.other', { player: name }) }));
  await undo(active).hover();
  await shot(active, 'rollback-choice');
  await undo(other).click();
  // Back before the pick: the item's use is the latest move again.
  await expect(undo(active)).toHaveAttribute('title', tr('history.undo.title', { move: tr('history.move.use.you', { card: cardName('bum_friend') }) }));
});
