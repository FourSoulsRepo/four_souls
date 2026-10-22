// Save and continue (7.13, N-11).
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { openHost } from './menu';

test.setTimeout(60_000);

test('the host saves the game; everyone comes back and plays on', async ({ browser }) => {
  const [gil, hob] = await startGame(browser, ['Gil', 'Hob'], { url: `${appURL()}&anim=off` });
  if (!gil || !hob) {
    throw new Error('no pages');
  }
  const turn = await gil.getByRole('main').getAttribute('aria-label');

  // Only the host is offered saving.
  await hob.getByRole('button', { name: 'Menu' }).click();
  await expect(hob.getByRole('dialog', { name: 'Game menu' }).getByRole('button', { name: 'Save and leave' })).toHaveCount(0);
  await hob.getByRole('button', { name: 'Cancel' }).click();

  await gil.getByRole('button', { name: 'Menu' }).click();
  await gil.getByRole('button', { name: 'Save and leave' }).click();
  for (const page of [gil, hob]) {
    await expect(page.getByText('The game was saved.')).toBeVisible();
  }

  // The host continues it later, from the host screen.
  await gil.getByRole('button', { name: 'Back' }).click();
  await openHost(gil);
  await gil.getByRole('button', { name: 'Continue a saved game' }).click();
  const save = gil.getByRole('listitem').filter({ hasText: 'Gil, Hob' });
  await save.getByRole('button', { name: 'Continue' }).click();
  await expect(gil.getByRole('heading', { name: 'Saved game' })).toBeVisible();
  const seats = gil.getByRole('list', { name: 'Seats' });
  await expect(seats).toContainText('Hob · not back yet');
  await shot(gil, 'saved-game');

  await hob.getByRole('listitem').filter({ hasText: 'Gil, Hob' }).getByRole('button', { name: 'Continue' }).click();
  await expect(seats).toContainText('Hob · back');
  for (const page of [gil, hob]) {
    await page.getByRole('button', { name: 'Ready', exact: true }).click();
  }
  for (const page of [gil, hob]) {
    await expect(page.getByRole('main', { name: turn ?? '' })).toBeVisible();
  }
});
