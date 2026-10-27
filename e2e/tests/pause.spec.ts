// Pause, vote, end of match (7.11).
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL, serverPort } from './server';
import { realAssets, shot } from './shot';
import { openOnline } from './menu';
import { locale, tr, trRe } from './locale';
import { activePage, click, settle, throwButton, turnAction } from './turn';

const fast = () => `${appURL()}&anim=off`;

test('a dropped player reconnects by itself to the same seat', async ({ browser }) => {
  const [amy, ben] = await startGame(browser, ['Amy', 'Ben'], { url: fast() });
  if (!amy || !ben) {
    throw new Error('no pages');
  }
  await ben.evaluate(() => {
    (window as unknown as { fourSouls: { drop: () => void } }).fourSouls.drop();
  });
  // The app reconnects at once (N-20): no dialog, the pause ends.
  await expect(ben.locator('.player-bottom')).toHaveAccessibleName(tr('table.own_name', { player: 'Ben' }));
  await expect(ben.getByRole('status', { name: tr('lost.reconnecting_label') })).toHaveCount(0);
  await expect(ben.getByRole('dialog', { name: tr('lost.title') })).toHaveCount(0);
  await expect(amy.getByRole('dialog', { name: tr('pause.label') })).toHaveCount(0);
});

test('an intent lost with a dropped connection can be sent again (19.86)', async ({ browser }) => {
  const pages = await startGame(browser, ['Gwen', 'Hugo'], { url: fast() });
  const active = await activePage(pages, click);
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
  // The throw never reaches the server: the connection drops with it,
  // as when the server closes a client that reads too slowly.
  await active.evaluate(() => {
    const w = window as unknown as { fourSouls: { drop: () => void } };
    const send = WebSocket.prototype.send;
    WebSocket.prototype.send = function (data) {
      if (String(data).includes('"type":"intent"')) {
        WebSocket.prototype.send = send;
        w.fourSouls.drop();
        return;
      }
      send.call(this, data);
    };
  });
  await click(active, throwButton(active));
  // The app reconnects to the same step; the Throw button is back.
  await expect(throwButton(active)).toBeVisible();
  await click(active, throwButton(active));
  await expect(throwButton(active)).toHaveCount(0);
});

test('a player who closed the app comes back to the same seat', async ({ browser }) => {
  const [kit, lou] = await startGame(browser, ['Kit', 'Lou'], { url: fast() });
  if (!kit || !lou) {
    throw new Error('no pages');
  }
  // Lou closes the app; Kit chooses to wait.
  const context = lou.context();
  await lou.close();
  const paused = kit.getByRole('dialog', { name: tr('pause.label') });
  await expect(paused).toContainText(tr('pause.title', { players: 'Lou' }));
  await paused.getByRole('button', { name: tr('pause.wait') }).click();

  // Lou starts the app again (same settings) and connects to the server.
  const again = await context.newPage();
  await again.goto(fast());
  await openOnline(again);
  await again.getByRole('button', { name: tr('lobby.connect') }).click();
  await expect(again.locator('.player-bottom')).toHaveAccessibleName(tr('table.own_name', { player: 'Lou' }));
  await expect(paused).toHaveCount(0);
});

test('players vote to go on without one who left', async ({ browser }) => {
  const [cal, dee] = await startGame(browser, ['Cal', 'Dee'], { url: fast() });
  if (!cal || !dee) {
    throw new Error('no pages');
  }
  await dee.close();
  const paused = cal.getByRole('dialog', { name: tr('pause.label') });
  await expect(paused).toContainText(tr('pause.title', { players: 'Dee' }));
  await shot(cal, 'paused');
  await paused.getByRole('button', { name: tr('pause.kick') }).click();
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
  const end = eli.getByRole('dialog', { name: tr('over.title') });
  await expect(end).toContainText(tr('over.you_win'));
  await end.getByRole('button', { name: tr('over.save') }).click();
  await expect(end.getByRole('status')).toHaveText(new RegExp(`^${trRe('over.saved').source.replace('.*?', '/fake/.+\\.fsrec')}$`));
  await shot(eli, 'game-over');
  await end.getByRole('button', { name: tr('over.menu') }).click();
  await expect(eli.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
});

// savedSeat is the seat token page's app keeps in its settings.
function savedSeat(page: Page): Promise<string> {
  return page.evaluate(() => {
    interface S {
      settings: { rejoin: { token: string } };
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().settings.rejoin.token;
  });
}

// restart closes page's app and starts it again with the same settings.
async function restart(page: Page): Promise<Page> {
  const context = page.context();
  await page.close();
  const again = await context.newPage();
  await realAssets(again);
  await again.goto(fast());
  return again;
}

// A crash or a closed window keeps the seat; the restarted menu offers
// it, labelled by the path it was joined by. Leaving on purpose or the
// end of the game forgets it (15.99, M-05).
test('a restarted app offers Rejoin in the menu, until the seat is left (15.99)', async ({ browser }) => {
  const [ann, bo] = await startGame(browser, ['Ann', 'Bo'], { url: fast() });
  if (!ann || !bo) {
    throw new Error('no pages');
  }
  // Bo joined by Online, Ann hosted on LAN.
  const bo2 = await restart(bo);
  const offer = bo2.getByRole('button', { name: tr('menu.rejoin', { game: `${tr('menu.rejoin_host', { host: 'Ann' })} · ${tr('menu.online')}` }) });
  await expect(offer).toBeVisible();
  // The menu itself did not wait for the answer, nor move for it.
  await expect(bo2.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
  await shot(bo2, 'menu-rejoin');
  await offer.click();
  await expect(bo2.locator('.player-bottom')).toHaveAccessibleName(tr('table.own_name', { player: 'Bo' }));

  const ann2 = await restart(ann);
  await ann2.getByRole('button', { name: tr('menu.rejoin', { game: `${tr('menu.rejoin_host', { host: 'Ann' })} · ${tr('menu.lan')}` }) }).click();
  await expect(ann2.locator('.player-bottom')).toHaveAccessibleName(tr('table.own_name', { player: 'Ann' }));
  // Ann leaves on purpose: no Rejoin, also after a restart.
  await ann2.getByRole('button', { name: tr('game_menu.button') }).click();
  await ann2.getByRole('button', { name: tr('game_menu.leave') }).click();
  await expect(ann2.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
  await expect.poll(() => savedSeat(ann2)).toBe('');
  const ann3 = await restart(ann2);
  await ann3.waitForLoadState('networkidle');
  await expect(ann3.getByRole('button', { name: new RegExp(`^${trRe('menu.rejoin').source}`) })).toHaveCount(0);

  // The game ends for Bo: the seat is forgotten.
  await bo2.evaluate(() => {
    interface S {
      game: { step: number; events: unknown; view: { over?: boolean; winners?: number[] } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const g = structuredClone(s.game);
    g.step += 1;
    g.events = [];
    g.view.over = true;
    g.view.winners = [0];
    s.receive({ type: 'update', data: g });
  });
  await expect(bo2.getByRole('dialog', { name: tr('over.title') })).toBeVisible();
  await expect.poll(() => savedSeat(bo2)).toBe('');
});

// A seat whose game is gone is forgotten when the menu asks (15.99).
test('a dead saved seat shows no Rejoin and is forgotten (15.99)', async ({ page }) => {
  await page.addInitScript(([port, language]) => {
    const key = 'fake-settings';
    if (!localStorage.getItem(key)) {
      const rejoin = { address: `127.0.0.1:${String(port)}`, token: 'ab'.repeat(16), game: 'gone', via: 'online' };
      localStorage.setItem(key, JSON.stringify({ nickname: 'Pam', recent: [], animation: 'off', fullscreen: false, mat: '', dice: 'tear', rejoin, language }));
    }
  }, [serverPort(), locale] as const);
  await page.goto(fast());
  await expect(page.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
  await expect.poll(() => savedSeat(page)).toBe('');
  await expect(page.getByRole('button', { name: new RegExp(`^${trRe('menu.rejoin').source}`) })).toHaveCount(0);
});
