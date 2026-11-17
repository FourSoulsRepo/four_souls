// Bots in the lobby (21.11, ADR 030 Decision 10): the host seats bots on
// free seats, changes their level or takes them away; the others see
// them with no buttons. Against the real Go server, which allows every
// level by default.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { gameId, lobbyGame, throwFirstRoll } from './game';
import { appURL, serverPort } from './server';
import { realAssets, shot } from './shot';
import { openHost, openOnline } from './menu';
import { tr, trRe } from './locale';

const standard = () => tr('bot.level.standard');
const hard = () => tr('bot.level.hard');

// hostTable opens a table of n seats from this app, with the Base Game
// only, as startGame does: a forced attack (Krampus) is not this test's.
async function hostTable(page: Page, name: string, n: number): Promise<Locator> {
  await realAssets(page);
  await page.goto(appURL());
  await openHost(page);
  await page.getByLabel(tr('lobby.nickname')).fill(name);
  await page.getByLabel(tr('host.players')).selectOption(String(n));
  const sets = page.getByRole('group', { name: tr('host.sets') });
  for (const set of [tr('set.gold_box'), tr('set.four_souls_plus')]) {
    await sets.getByLabel(set).uncheck();
  }
  await page.getByRole('button', { name: tr('host.host'), exact: true }).click();
  const seats = page.getByRole('list', { name: tr('lobby.seats.seats') });
  await expect(seats).toContainText(`${name} ${tr('lobby.seats.you')}`);
  return seats;
}

// seat is one row of the seat list.
function seat(seats: Locator, i: number): Locator {
  return seats.locator(`li[data-seat="${String(i)}"]`);
}

// addBot presses a free seat's Add bot button for the level.
async function addBot(seats: Locator, i: number, level: string): Promise<void> {
  await seat(seats, i).getByRole('button', { name: tr('lobby.bots.add_level', { level }) }).click();
}

// joinTable joins the game id from another app on the same server;
// listed checks the game's row first.
async function joinTable(page: Page, name: string, id: string, listed: (row: Locator) => Promise<void>): Promise<Locator> {
  await page.goto(appURL());
  await openOnline(page);
  await page.getByLabel(tr('lobby.nickname')).fill(name);
  await page.getByLabel(tr('lobby.online.address')).fill(`127.0.0.1:${String(serverPort())}`);
  await page.getByRole('button', { name: tr('lobby.connect') }).click();
  const join = lobbyGame(page, id).getByRole('button', { name: tr('lobby.join') });
  await expect(join).toBeEnabled();
  await listed(lobbyGame(page, id));
  await join.dispatchEvent('click');
  const seats = page.getByRole('list', { name: tr('lobby.seats.seats') });
  await expect(seats).toContainText(`${name} ${tr('lobby.seats.you')}`);
  return seats;
}

// Use case "Play with bots" (UC-44).
test('the host fills the table with 3 bots and starts alone', async ({ browser }) => {
  const ann = await (await browser.newContext()).newPage();
  const seats = await hostTable(ann, 'Ann', 4);
  await expect(seat(seats, 1).getByRole('group', { name: tr('lobby.bots.add') })).toBeVisible();
  await shot(ann, 'table-bots-free');

  await addBot(seats, 1, standard());
  await addBot(seats, 2, hard());
  await addBot(seats, 3, standard());
  for (const [i, level] of ['standard', 'hard', 'standard'].entries()) {
    const row = seat(seats, i + 1);
    const name = tr('bot.number', { n: i + 1 });
    await expect(row).toContainText(name);
    await expect(row.getByTestId('bot-icon')).toBeVisible();
    await expect(row.getByRole('combobox', { name: tr('lobby.bots.level', { name }) })).toHaveValue(level);
    await expect(row).toContainText(tr('lobby.seats.is_ready'));
    // A bot has no connection, yet it is never away (21.11 task 9).
    await expect(row).not.toContainText(tr('lobby.seats.away'));
  }
  await shot(ann, 'table-bots-host');

  // The host changes a level, takes a bot away and seats another: it
  // takes the smallest free number again.
  const bot3 = tr('bot.number', { n: 3 });
  await seat(seats, 3).getByRole('combobox', { name: tr('lobby.bots.level', { name: bot3 }) }).selectOption('hard');
  await expect(seat(seats, 3).getByRole('combobox', { name: tr('lobby.bots.level', { name: bot3 }) })).toHaveValue('hard');
  const bot2 = tr('bot.number', { n: 2 });
  await seat(seats, 2).getByRole('button', { name: tr('lobby.bots.remove', { name: bot2 }) }).click();
  await expect(seat(seats, 2)).toContainText(tr('lobby.seats.waiting'));
  await addBot(seats, 2, standard());
  await expect(seat(seats, 2).getByRole('combobox', { name: tr('lobby.bots.level', { name: bot2 }) })).toHaveValue('standard');
  await expect(ann.getByRole('status')).toHaveText('');

  // Every seat is taken; the game waits only for Ann.
  await ann.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  await expect(ann.getByRole('main', { name: new RegExp(`^${trRe('game.turn').source}`) })).toBeVisible();
  // The bots throw their own first-player dice.
  await throwFirstRoll([ann]);
  await expect(ann.getByRole('main', { name: trRe('game.turn') })).toBeVisible();
});

test('the other players see the bots and their levels, with no buttons', async ({ browser }) => {
  const ann = await (await browser.newContext()).newPage();
  const bo = await (await browser.newContext()).newPage();
  const annSeats = await hostTable(ann, 'Ann', 4);
  await addBot(annSeats, 2, hard());
  const id = await gameId(ann);

  // A bot's seat is taken in the game list.
  await realAssets(bo);
  const boSeats = await joinTable(bo, 'Bo', id, async (row) => {
    await expect(row).toContainText(tr('lobby.game_row', { host: 'Ann', taken: 2, n: 4 }));
  });
  await expect(seat(annSeats, 1)).toContainText('Bo');
  await addBot(annSeats, 3, standard());

  await expect(seat(boSeats, 2)).toContainText(tr('bot.name', { n: 1, level: hard() }));
  await expect(seat(boSeats, 2).getByTestId('bot-icon')).toBeVisible();
  await expect(seat(boSeats, 3)).toContainText(tr('bot.name', { n: 2, level: standard() }));
  await expect(seat(boSeats, 3)).not.toContainText(tr('lobby.seats.away'));
  await expect(boSeats.getByRole('button')).toHaveCount(0);
  await expect(boSeats.getByRole('combobox')).toHaveCount(0);
  await shot(bo, 'table-bots-guest');

  // A free seat shows no Add bot to a guest either.
  const bot2 = tr('bot.number', { n: 2 });
  await seat(annSeats, 3).getByRole('button', { name: tr('lobby.bots.remove', { name: bot2 }) }).click();
  await expect(seat(boSeats, 3)).toContainText(tr('lobby.seats.waiting'));
  await expect(boSeats.getByRole('button')).toHaveCount(0);
});

test('a server that allows no bots shows no bot buttons', async ({ page }) => {
  // The welcome of a server started with -bots=off lists no levels.
  await page.routeWebSocket(`ws://127.0.0.1:${String(serverPort())}/ws`, (ws) => {
    const server = ws.connectToServer();
    ws.onMessage((msg) => {
      server.send(msg);
    });
    server.onMessage((msg) => {
      const env = JSON.parse(String(msg)) as { type: string; data?: { bots?: string[] } };
      if (env.type === 'welcome' && env.data) {
        delete env.data.bots;
        ws.send(JSON.stringify(env));
        return;
      }
      ws.send(msg);
    });
  });
  const seats = await hostTable(page, 'Cy', 3);
  await expect(seat(seats, 1)).toContainText(tr('lobby.seats.waiting'));
  await expect(seats.getByRole('button')).toHaveCount(0);
});
