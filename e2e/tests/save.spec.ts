// Save and continue (7.13, N-11).
import { expect, test } from './fixtures';
import { gameId, lobbyGame, startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { tr } from './locale';

test.setTimeout(60_000);

test('the host saves the game; everyone comes back and plays on', async ({ browser }) => {
  const [gil, hob] = await startGame(browser, ['Gil', 'Hob'], { url: `${appURL()}&anim=off` });
  if (!gil || !hob) {
    throw new Error('no pages');
  }
  const turn = await gil.getByRole('main').getAttribute('aria-label');

  // Only the host is offered saving.
  await hob.getByRole('button', { name: tr('game_menu.button') }).click();
  await expect(hob.getByRole('dialog', { name: tr('game_menu.label') }).getByRole('button', { name: tr('game_menu.save') })).toHaveCount(0);
  await hob.getByRole('button', { name: tr('actions.cancel') }).click();

  await gil.getByRole('button', { name: tr('game_menu.button') }).click();
  await gil.getByRole('button', { name: tr('game_menu.save') }).click();
  for (const page of [gil, hob]) {
    await expect(page.getByText(tr('lobby.games.saved'))).toBeVisible();
  }

  // The host continues it later, from the host screen.
  // Back from the Games list is the LAN screen, with its Host button.
  await gil.getByRole('button', { name: tr('lobby.back') }).click();
  await gil.getByRole('button', { name: tr('host.title_host') }).click();
  await gil.getByRole('button', { name: tr('host.continue_saved') }).click();
  const save = gil.getByRole('listitem').filter({ hasText: 'Gil, Hob' });
  await save.getByRole('button', { name: tr('lobby.games.continue') }).click();
  await expect(gil.getByRole('heading', { name: tr('lobby.seats.saved_title') })).toBeVisible();
  const seats = gil.getByRole('list', { name: tr('lobby.seats.seats') });
  await expect(seats).toContainText(`Hob · ${tr('lobby.seats.not_back')}`);
  await shot(gil, 'saved-game');

  await hob.getByRole('listitem').filter({ hasText: 'Gil, Hob' }).getByRole('button', { name: tr('lobby.games.continue') }).click();
  await expect(seats).toContainText(`Hob · ${tr('lobby.seats.back')}`);
  for (const page of [gil, hob]) {
    await page.getByRole('button', { name: tr('lobby.seats.ready'), exact: true }).click();
  }
  for (const page of [gil, hob]) {
    await expect(page.getByRole('main', { name: turn ?? '' })).toBeVisible();
  }
});

test('the host ends the game; everyone is back on the game list (19.84)', async ({ browser }) => {
  const [gil, hob] = await startGame(browser, ['Gil', 'Hob'], { url: `${appURL()}&anim=off` });
  if (!gil || !hob) {
    throw new Error('no pages');
  }
  const id = await gameId(gil);

  // Only the host is offered ending it.
  await hob.getByRole('button', { name: tr('game_menu.button') }).click();
  await expect(hob.getByRole('dialog', { name: tr('game_menu.label') }).getByRole('button', { name: tr('game_menu.end') })).toHaveCount(0);
  await hob.getByRole('button', { name: tr('game_menu.cancel') }).click();

  // Back on the confirm step keeps the game.
  await gil.getByRole('button', { name: tr('game_menu.button') }).click();
  await gil.getByRole('button', { name: tr('game_menu.end') }).click();
  await expect(gil.getByRole('heading', { name: tr('game_menu.end_title') })).toBeVisible();
  await shot(gil, 'end-game');
  await gil.getByRole('button', { name: tr('game_menu.back') }).click();
  await gil.getByRole('button', { name: tr('game_menu.end') }).click();
  await gil.getByRole('dialog', { name: tr('game_menu.label') }).getByRole('button', { name: tr('game_menu.end') }).click();

  const ended = hob.getByRole('dialog', { name: tr('notice.game_idle.title') });
  await expect(ended).toContainText(tr('notice.game_ended.text'));
  await shot(hob, 'game-ended');
  await ended.getByRole('button', { name: tr('notice.ok') }).click();
  for (const page of [gil, hob]) {
    await expect(page.getByRole('heading', { name: tr('lobby.games.title'), exact: true })).toBeVisible();
    await expect(lobbyGame(page, id)).toHaveCount(0);
  }
  await expect(gil.getByRole('dialog', { name: tr('notice.game_idle.title') })).toHaveCount(0);
});
