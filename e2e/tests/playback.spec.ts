// Event playback and animations (7.6).
import { expect, test } from '@playwright/test';
import { startGame } from './game';
import { shot } from './shot';

test('the animation speed setting survives a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Animations').selectOption('fast');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Animations')).toHaveValue('fast');
});

test('a fast burst of updates plays in order and ends on the last one', async ({ browser }) => {
  const [kai] = await startGame(browser, ['Kai', 'Lu']);
  if (!kai) {
    throw new Error('no pages');
  }
  // Mark one card: if updates redrew the table from scratch, the mark
  // would be gone and every card would pop in again.
  const character = kai.locator('.player-bottom .player-cards > .card').first();
  await character.evaluate((el) => el.setAttribute('data-mark', 'kept'));
  // Replay the current view 20 times as new turns, each with a dice roll
  // and a looted card, as fast as the page takes them.
  await kai.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { number: number; active: number } } } | null;
        lobby: { table: { you: number } | null };
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const you = s.lobby.table?.you ?? 0;
    for (let n = 2; n <= 21; n++) {
      const game = structuredClone(s.game);
      if (!game) {
        throw new Error('no game');
      }
      game.view.turn.number = n;
      s.receive({
        type: 'update',
        data: {
          ...game,
          step: n,
          events: [
            { kind: 'dice_rolled', player: you, amount: (n % 6) + 1, text: 'roll' },
            { kind: 'looted', player: you, object: 1000 + n },
          ],
        },
      });
    }
  });
  await expect(kai.getByTestId('effect-dice')).toBeVisible();
  await kai.waitForTimeout(300); // let the dice settle for the picture
  await shot(kai, 'effect-dice');
  await expect(kai.getByTestId('effect-fly')).toBeVisible();
  await expect(kai.getByRole('main', { name: 'Turn 21' })).toBeVisible({ timeout: 20_000 });
  await expect(kai.locator('.effects')).toHaveCount(0);
  await expect(character).toHaveAttribute('data-mark', 'kept');
  await expect(character).toHaveCSS('opacity', '1');
});
