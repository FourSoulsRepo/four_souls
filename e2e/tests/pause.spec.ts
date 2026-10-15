// Pause, vote, end of match (7.11).
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';

const fast = () => `${appURL()}&anim=off`;

test('a dropped player reconnects to the same seat; the others wait', async ({ browser }) => {
  const [amy, ben] = await startGame(browser, ['Amy', 'Ben'], { url: fast() });
  if (!amy || !ben) {
    throw new Error('no pages');
  }
  await ben.evaluate(() => {
    (window as unknown as { fourSouls: { drop: () => void } }).fourSouls.drop();
  });
  await expect(ben.getByRole('dialog', { name: 'Connection lost' })).toBeVisible();
  await expect(amy.getByRole('dialog', { name: 'Paused' })).toContainText('Waiting for Ben');
  await shot(amy, 'paused');
  await shot(ben, 'connection-lost');

  await ben.getByRole('button', { name: 'Reconnect' }).click();
  await expect(ben.getByRole('dialog', { name: 'Connection lost' })).toHaveCount(0);
  await expect(amy.getByRole('dialog', { name: 'Paused' })).toHaveCount(0);
  await expect(ben.locator('.player-bottom')).toHaveAccessibleName('Ben (you)');
});

test('a player who closed the app comes back to the same seat', async ({ browser }) => {
  const [kit, lou] = await startGame(browser, ['Kit', 'Lou'], { url: fast() });
  if (!kit || !lou) {
    throw new Error('no pages');
  }
  // Lou closes the app; Kit chooses to wait.
  const context = lou.context();
  await lou.close();
  const paused = kit.getByRole('dialog', { name: 'Paused' });
  await expect(paused).toContainText('Waiting for Lou');
  await paused.getByRole('button', { name: 'Wait' }).click();

  // Lou starts the app again (same settings) and connects to the server.
  const again = await context.newPage();
  await again.goto(fast());
  await again.getByRole('button', { name: 'Join a game' }).click();
  await again.getByRole('button', { name: 'Connect' }).click();
  await expect(again.locator('.player-bottom')).toHaveAccessibleName('Lou (you)');
  await expect(paused).toHaveCount(0);
});

test('players vote to go on without one who left', async ({ browser }) => {
  const [cal, dee] = await startGame(browser, ['Cal', 'Dee'], { url: fast() });
  if (!cal || !dee) {
    throw new Error('no pages');
  }
  await dee.close();
  const paused = cal.getByRole('dialog', { name: 'Paused' });
  await expect(paused).toContainText('Waiting for Dee');
  await paused.getByRole('button', { name: 'Go on without them' }).click();
  await expect(paused).toHaveCount(0);
});

test('the end screen names the winner and saves the record', async ({ browser }) => {
  const [eli] = await startGame(browser, ['Eli', 'Fox'], { url: fast() });
  if (!eli) {
    throw new Error('no pages');
  }
  // Finishing a real game takes long; play its last update instead.
  await eli.evaluate(() => {
    interface S {
      game: { step: number; events: unknown; view: { over?: boolean; winners?: number[] } };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const g = structuredClone(s.game);
    g.step += 1;
    g.events = [{ kind: 'game_won', player: s.lobby.table.you }];
    g.view.over = true;
    g.view.winners = [s.lobby.table.you];
    s.receive({ type: 'update', data: g });
  });
  const end = eli.getByRole('dialog', { name: 'Game over' });
  await expect(end).toContainText('You win!');
  await end.getByRole('button', { name: 'Save record' }).click();
  await expect(end.getByRole('status')).toHaveText(/^Saved to \/fake\/.+\.fsrec$/);
  await shot(eli, 'game-over');
  await end.getByRole('button', { name: 'Back to menu' }).click();
  await expect(eli.getByRole('button', { name: 'Host a game' })).toBeVisible();
});
