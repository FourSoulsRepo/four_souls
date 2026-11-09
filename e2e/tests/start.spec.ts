// The start-of-game screen (20.22, ADR 027 "Start of game", UC-26 6a):
// Eden picks her starting item, then confirms; the others see who is
// choosing; several Edens choose at the same time.
import type { Page } from '@playwright/test';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { tr } from './locale';
import { gameID } from './saves';
import { appURL, serverPort } from './server';
import { shot } from './shot';
import { seen, throwButton } from './turn';

const root = join(import.meta.dirname, '..', '..');

// characters are a set's character cards, reprints left out.
function characters(set: string): string[] {
  const data = JSON.parse(readFileSync(join(root, 'pkg', 'card_db', 'data', `${set}.json`), 'utf8')) as { cards: { id: string; type?: string }[] };
  return data.cards.filter((c) => c.type === 'Character Card').map((c) => c.id);
}

// start is the page's start-of-game window, if it shows.
function start(page: Page) {
  return page.getByRole('dialog', { name: tr('choice.title.starting_item') });
}

// edenOf waits until one page shows the start-of-game screen; it
// returns that page and the others.
async function edenOf(pages: Page[]): Promise<{ eden: Page; others: Page[] }> {
  let eden: Page | undefined;
  await expect
    .poll(async () => {
      for (const p of pages) {
        if (await start(p).isVisible()) {
          eden = p;
        }
      }
      return eden !== undefined;
    })
    .toBe(true);
  if (!eden) {
    throw new Error('no Eden');
  }
  const found = eden;
  return { eden: found, others: pages.filter((p) => p !== found) };
}

test('Eden picks her starting item on the start-of-game screen, then confirms (20.22)', async ({ browser }) => {
  // The bans leave Eden and Isaac only: one of the two players is
  // always Eden, whose start-of-game screen this test checks.
  const bans = characters('b2').filter((c) => c !== 'eden' && c !== 'isaac');
  const pages = await startGame(browser, ['Ann', 'Bo'], { url: appURL(), roll: false, bans });
  const { eden, others } = await edenOf(pages);
  const [other] = others;
  if (!other) {
    throw new Error('no other page');
  }
  const screen = start(eden);
  const edenName = (await seen(eden)).you === 0 ? 'Ann' : 'Bo';
  await expect(screen.getByRole('heading', { name: tr('start.title') })).toBeVisible();
  await expect(screen).toContainText(tr('start.caption', { card: 'Eden' }));
  const confirm = screen.getByRole('button', { name: tr('start.confirm') });
  await expect(confirm).toBeDisabled();
  // The cards are dealt from the treasure deck; then they stay still.
  await expect(eden.locator('.choice.dealing')).toHaveCount(0, { timeout: 5000 });
  await eden.mouse.move(1, 1);
  await shot(eden, '20.22-start-own');

  // The other player sees who is choosing, no window, and no roll yet
  // (19.142); the badge names the start of the game.
  const wait = other.getByRole('status', { name: tr('start.waiting_label') });
  await expect(wait).toHaveText(tr('start.choosing', { player: edenName }));
  await expect(other.getByRole('dialog')).toHaveCount(0);
  await expect(other.getByRole('complementary', { name: tr('stack.title') }).getByRole('group', { name: tr('first_roll.label') })).toHaveCount(0);
  await expect(other.getByRole('status', { name: tr('turn_badge.label') })).toHaveText(tr('turn_badge.start'));
  await shot(other, '20.22-start-other');

  // A click picks a card, a second unpicks it; nothing is sent.
  const cards = screen.locator('.choice-card');
  await expect(cards).toHaveCount(3);
  await cards.nth(1).click();
  await expect(cards.nth(1)).toHaveAttribute('aria-pressed', 'true');
  await expect(confirm).toBeEnabled();
  await cards.nth(1).click();
  await expect(cards.nth(1)).toHaveAttribute('aria-pressed', 'false');
  await expect(confirm).toBeDisabled();
  // Another card: the pick moves to it.
  await cards.nth(1).click();
  await cards.nth(2).click();
  await expect(cards.nth(1)).toHaveAttribute('aria-pressed', 'false');
  await expect(cards.nth(2)).toHaveAttribute('aria-pressed', 'true');
  const kept = (await cards.nth(2).getAttribute('aria-label')) ?? '';
  await eden.mouse.move(1, 1);
  await shot(eden, '20.22-start-picked');
  await expect(start(eden)).toBeVisible();
  expect((await seen(eden)).waiting).toBe((await seen(eden)).you);

  // Confirm: the window goes; the kept card flies to Eden's items.
  await confirm.click();
  await expect(start(eden)).toHaveCount(0);
  const seat = (await seen(eden)).you;
  for (const page of pages) {
    await expect(page.locator(`[data-zone="player-${String(seat)}"]`).getByLabel(kept, { exact: true }).first()).toBeVisible();
    await expect(page.getByTestId('choice-kept')).toHaveCount(0, { timeout: 5000 });
  }
  await expect(wait).toHaveCount(0);
  // Then the first-player roll.
  const first = (await seen(eden)).you === 0 ? eden : other;
  await expect(throwButton(first)).toBeVisible({ timeout: 20_000 });
});

// twoEdens rewrites the saved game of a deal with one Eden: the other
// player becomes Eden too, with her own three cards, the next ones of
// the treasure deck, as the engine opens several start-of-game choices
// (20.22). The server deals each character once, so a deal never has two.
function twoEdens(game: string): void {
  const file = join(root, '.cache', 'e2e', String(serverPort()), 'saves', `${game}.fssave`);
  interface Choice {
    purpose: number;
    player: number;
    rule: string;
    objects: number[];
    look: { count: number };
    owner: number;
    to: number;
  }
  const save = JSON.parse(readFileSync(file, 'utf8')) as {
    state: {
      players: { character: number }[];
      objects: { card: string }[];
      decks: number[][];
      choice?: Choice;
      start_choices?: Choice[];
      waiting: Record<string, unknown>;
    };
  };
  const { state } = save;
  const first = state.start_choices?.[0] ?? state.choice;
  const deck = state.decks[0] ?? [];
  if (!first) {
    throw new Error('no start-of-game choice in the save');
  }
  const other = first.player === 0 ? 1 : 0;
  const character = state.objects[state.players[other]?.character ?? -1];
  if (!character) {
    throw new Error('no character');
  }
  character.card = 'eden';
  const objects = deck.slice(-6, -3).reverse();
  const choices = [first, { ...first, player: other, objects, owner: 0, to: 0 }].sort((a, b) => a.player - b.player);
  const [top] = choices;
  if (!top) {
    throw new Error('no choices');
  }
  state.start_choices = choices;
  state.choice = top;
  const options = top.objects.map((id) => ({ key: 'ask.card', args: { card: { card: state.objects[id]?.card ?? '' } } }));
  state.waiting = { kind: 4, player: top.player, options, purpose: 10, look: { count: 3 }, players: [0, 1] };
  writeFileSync(file, JSON.stringify(save));
}

test('two Edens choose at the same time; the roll waits for both (20.22)', async ({ browser }) => {
  // The bans leave Eden and Isaac only: one player is Eden; the saved
  // game then makes the other Eden too (twoEdens), so both choose.
  const bans = characters('b2').filter((c) => c !== 'eden' && c !== 'isaac');
  const names = ['Cyd', 'Dot'];
  const pages = await startGame(browser, names, { url: appURL(), roll: false, bans });
  await edenOf(pages);
  const [cyd, dot] = pages;
  if (!cyd || !dot) {
    throw new Error('no pages');
  }
  const game = await gameID(cyd);
  await cyd.getByRole('button', { name: tr('game_menu.button') }).click();
  await cyd.getByRole('button', { name: tr('game_menu.save') }).click();
  for (const page of pages) {
    await expect(page.getByText(tr('lobby.games.saved'))).toBeVisible();
  }
  twoEdens(game);
  await cyd.getByRole('button', { name: tr('lobby.back') }).click();
  await cyd.getByRole('button', { name: tr('host.title_host') }).click();
  await cyd.getByRole('button', { name: tr('host.continue_saved') }).click();
  await cyd.getByRole('listitem').filter({ hasText: 'Cyd' }).getByRole('button', { name: tr('lobby.games.continue') }).click();
  await dot.getByRole('listitem').filter({ hasText: 'Cyd' }).getByRole('button', { name: tr('lobby.games.continue') }).click();
  for (const page of pages) {
    await page.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  }
  for (const page of pages) {
    await expect(start(page)).toBeVisible();
    await expect(start(page).locator('.choice-card')).toHaveCount(3);
    await expect(page.locator('.choice.dealing')).toHaveCount(0, { timeout: 5000 });
    await page.mouse.move(1, 1);
  }
  // Each sees their own cards: no card twice.
  const labels = async (page: Page) => start(page).locator('.choice-card').evaluateAll((els) => els.map((el) => el.getAttribute('data-object')));
  const own = [...(await labels(cyd)), ...(await labels(dot))];
  expect(new Set(own).size).toBe(6);
  await shot(cyd, '20.22-start-two-edens');
  // The second seat confirms first; the first still chooses, unchanged.
  await start(dot).locator('.choice-card').first().click();
  await start(dot).getByRole('button', { name: tr('start.confirm') }).click();
  await expect(start(dot)).toHaveCount(0);
  await expect(dot.getByRole('status', { name: tr('start.waiting_label') })).toHaveText(tr('start.choosing', { player: 'Cyd' }));
  await expect(start(cyd)).toBeVisible();
  await expect(start(cyd).locator('.choice-card')).toHaveCount(3);
  await expect(throwButton(cyd)).toHaveCount(0);
  await start(cyd).locator('.choice-card').nth(1).click();
  await start(cyd).getByRole('button', { name: tr('start.confirm') }).click();
  // Both chose: the roll starts with the first seat.
  await expect(throwButton(cyd)).toBeVisible({ timeout: 20_000 });
  for (const page of pages) {
    await expect(page.getByRole('status', { name: tr('start.waiting_label') })).toHaveCount(0);
  }
});
