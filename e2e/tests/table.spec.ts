// Table layout (7.5), with real games on the Go server.
import { expect, test, type Page } from '@playwright/test';
import { expectFits, startGame } from './game';
import { serverPort } from './server';
import { shot } from './shot';

// expectDecksClear fails when anything lies over a deck or a discard
// (F-02): their corners and middle must hit the pile itself.
async function expectDecksClear(page: Page): Promise<void> {
  const covered = await page.evaluate(() => {
    const out: string[] = [];
    for (const pile of document.querySelectorAll('[data-zone^="deck-"], [data-zone^="discard-"]')) {
      const r = pile.getBoundingClientRect();
      const points = [
        [r.left + 3, r.top + 3],
        [r.right - 3, r.top + 3],
        [r.left + 3, r.bottom - 3],
        [r.right - 3, r.bottom - 3],
        [(r.left + r.right) / 2, (r.top + r.bottom) / 2],
      ];
      for (const [x = 0, y = 0] of points) {
        const top = document.elementFromPoint(x, y);
        if (!top || !pile.contains(top)) {
          out.push(`${pile.getAttribute('data-zone') ?? ''} under ${top?.className ?? 'nothing'}`);
        }
      }
    }
    return out;
  });
  expect(covered).toEqual([]);
}

// emptySlots rewrites the page's game updates: the first shop slot is
// empty, three more empty slots follow, and no monsters or bonus souls
// are left. The real game goes on unchanged.
async function emptySlots(page: Page): Promise<void> {
  await page.routeWebSocket(new RegExp(`:${String(serverPort())}/`), (ws) => {
    const server = ws.connectToServer();
    server.onMessage((message) => {
      const env = JSON.parse(String(message)) as { type: string; data?: { view?: Record<string, unknown> } };
      const view = env.data?.view;
      if (env.type === 'update' && view && Array.isArray(view.shop)) {
        view.shop = [{}, ...view.shop.slice(1), {}, {}, {}];
        view.monsters = [];
        view.bonus_souls = [];
      }
      ws.send(JSON.stringify(env));
    });
  });
}

test('two players face each other; the table fits 1280x720', async ({ browser }) => {
  const [eve, fay] = await startGame(browser, ['Eve', 'Fay']);
  if (!eve || !fay) {
    throw new Error('no pages');
  }
  for (const [page, me, other] of [
    [eve, 'Eve', 'Fay'],
    [fay, 'Fay', 'Eve'],
  ] as const) {
    await expect(page.locator('.player-bottom')).toHaveAccessibleName(`${me} (you)`);
    await expect(page.locator('.player-top')).toHaveAccessibleName(other);
    await expect(page.getByRole('group', { name: 'Your hand' }).getByRole('img')).not.toHaveCount(0);
    await expect(page.getByRole('group', { name: /^Loot deck: \d+ cards$/ })).toBeVisible();
    // The opponent's hand is one loot back with its count (15.11).
    const theirHand = page.locator('.player-top').getByRole('group', { name: /^Hand: \d+ cards$/ });
    await expect(theirHand.getByRole('img', { name: 'Face-down loot card' })).toHaveCount(1);
    await expect(theirHand.locator('.pile-count')).toHaveText(/^\d+$/);
    await expect(page.getByRole('group', { name: 'Shop' }).getByRole('img')).not.toHaveCount(0);
    await expect(page.getByRole('group', { name: 'Monsters' }).getByRole('img')).not.toHaveCount(0);
    await expectFits(page);
    await expectDecksClear(page);
  }
  await shot(eve, 'table-2p');
});

test('four players sit left, top and right in turn order', async ({ browser }) => {
  const pages = await startGame(browser, ['Gus', 'Hal', 'Ivy', 'Jo']);
  const hal = pages[1];
  if (!hal) {
    throw new Error('no pages');
  }
  await expect(hal.locator('.player-bottom')).toHaveAccessibleName('Hal (you)');
  await expect(hal.locator('.player-left')).toHaveAccessibleName('Ivy');
  await expect(hal.locator('.player-top')).toHaveAccessibleName('Jo');
  await expect(hal.locator('.player-right')).toHaveAccessibleName('Gus');
  for (const page of pages) {
    await expectFits(page);
    await expectDecksClear(page);
  }
  await shot(hal, 'table-4p');
});

test('three players sit left and right; nothing covers the decks', async ({ browser }) => {
  const pages = await startGame(browser, ['Kip', 'Lea', 'Max']);
  const kip = pages[0];
  if (!kip) {
    throw new Error('no pages');
  }
  await expect(kip.locator('.player-left')).toHaveAccessibleName('Lea');
  await expect(kip.locator('.player-right')).toHaveAccessibleName('Max');
  for (const page of pages) {
    await expectFits(page);
    await expectDecksClear(page);
  }
  await shot(kip, 'table-3p');
});

test('the board keeps its size when shop slots empty (V-23)', async ({ browser }) => {
  let opened = 0;
  const [ned, ola] = await startGame(browser, ['Ned', 'Ola'], {
    open: (page) => (opened++ === 1 ? emptySlots(page) : Promise.resolve()),
  });
  if (!ned || !ola) {
    throw new Error('no pages');
  }
  const shop = ola.getByRole('group', { name: 'Shop' });
  await expect(shop.locator('.pile-empty')).toHaveCount(4);
  await expect(ola.getByRole('group', { name: 'Monsters' }).getByRole('img')).toHaveCount(0);
  await expect(ned.getByRole('group', { name: 'Shop' }).locator('.pile-empty')).toHaveCount(0);
  const full = await ned.getByRole('region', { name: 'Board' }).boundingBox();
  const empty = await ola.getByRole('region', { name: 'Board' }).boundingBox();
  expect(empty).toEqual(full);
  await expectDecksClear(ola);
  await shot(ola, 'table-empty-slots');
});
