// Table layout (7.5), with real games on the Go server.
import { expect, test, type Browser, type Page } from '@playwright/test';
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

type View = Record<string, unknown>;

// rewriteView changes the page's game updates before the app sees them;
// the real game goes on unchanged.
async function rewriteView(page: Page, change: (view: View) => void): Promise<void> {
  await page.routeWebSocket(new RegExp(`:${String(serverPort())}/`), (ws) => {
    const server = ws.connectToServer();
    server.onMessage((message) => {
      const env = JSON.parse(String(message)) as { type: string; data?: { view?: View } };
      const view = env.data?.view;
      if (env.type === 'update' && view && Array.isArray(view.shop)) {
        change(view);
      }
      ws.send(JSON.stringify(env));
    });
  });
}

// slots makes n slots of the given cards, with ids no real card has.
function slots(cards: string[], n: number, firstId: number): { top: { id: number; card: string; charged: boolean } }[] {
  return Array.from({ length: n }, (_, i) => ({ top: { id: firstId + i, card: cards[i % cards.length] ?? '', charged: true } }));
}

const monsterRefs = ['big_spider', 'black_bony', 'boom_fly', 'clotty', 'cod_worm', 'dinga'];
const treasureRefs = ['baby_haunt', 'belly_button', 'bobs_brain', 'breakfast', 'brimstone', 'bum_bo'];

// startPair starts a two-player game; the second page's updates go
// through change, the first page's do not.
async function startPair(browser: Browser, names: [string, string], change: (view: View) => void): Promise<[Page, Page]> {
  let opened = 0;
  const [plain, changed] = await startGame(browser, names, {
    open: (page) => (opened++ === 1 ? rewriteView(page, change) : Promise.resolve()),
  });
  if (!plain || !changed) {
    throw new Error('no pages');
  }
  return [plain, changed];
}

// boardBox is where the board is on screen.
async function boardBox(page: Page) {
  return page.getByRole('region', { name: 'Board' }).boundingBox();
}

// slotBoxes are the slots of a row, left to right. A slot, not its card:
// a new card pops in with a scale, which would change its box.
async function slotBoxes(page: Page, row: string) {
  const cards = page.getByRole('group', { name: row }).locator('> *');
  return Promise.all((await cards.all()).map(async (c) => (await c.boundingBox()) ?? { x: 0, y: 0, width: 0, height: 0 }));
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
  const [ned, ola] = await startPair(browser, ['Ned', 'Ola'], (view) => {
    const shop = view.shop as unknown[];
    view.shop = [{}, ...shop.slice(1), {}, {}, {}];
    view.monsters = [];
    view.bonus_souls = [];
  });
  await expect(ola.getByRole('group', { name: 'Shop' }).locator('.pile-empty')).toHaveCount(4);
  await expect(ola.getByRole('group', { name: 'Monsters' }).getByRole('img')).toHaveCount(0);
  await expect(ned.getByRole('group', { name: 'Shop' }).locator('.pile-empty')).toHaveCount(0);
  expect(await boardBox(ola)).toEqual(await boardBox(ned));
  await expectDecksClear(ola);
  await shot(ola, 'table-empty-slots');
});

// Some cards add monster or shop slots. Up to four slots show at full
// size; more share the row. The board never changes (V-23).
for (const [n, names] of [
  [3, ['Pam', 'Quy']],
  [4, ['Rex', 'Sol']],
  [6, ['Tom', 'Uma']],
] as const) {
  test(`${String(n)} monster and shop slots fit the same board`, async ({ browser }) => {
    const [plain, many] = await startPair(browser, [...names], (view) => {
      view.monsters = slots(monsterRefs, n, 900000);
      view.shop = slots(treasureRefs, n, 900100);
    });
    for (const row of ['Monsters', 'Shop']) {
      await expect(many.getByRole('group', { name: row }).locator('.card')).toHaveCount(n);
    }
    const board = await boardBox(plain);
    expect(await boardBox(many)).toEqual(board);
    const full = (await slotBoxes(plain, 'Shop'))[0]?.width ?? 0;
    expect(full).toBeGreaterThan(0);
    for (const row of ['Monsters', 'Shop']) {
      const area = await many.getByRole('group', { name: row }).boundingBox();
      const boxes = await slotBoxes(many, row);
      for (const [i, r] of boxes.entries()) {
        // Full size up to four slots; smaller, never overlapping, past that.
        if (n <= 4) {
          expect(r.width).toBeCloseTo(full, 0);
        } else {
          expect(r.width).toBeLessThan(full);
        }
        expect(r.x).toBeGreaterThanOrEqual(area?.x ?? 0);
        expect(r.x + r.width).toBeLessThanOrEqual((area?.x ?? 0) + (area?.width ?? 0) + 0.5);
        const next = boxes[i + 1];
        if (next) {
          expect(next.x).toBeGreaterThanOrEqual(r.x + r.width);
        }
      }
    }
    await expectDecksClear(many);
    await expectFits(many);
    await shot(many, `table-slots-${String(n)}`);
  });
}
