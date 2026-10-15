// Table layout (7.5), with real games on the Go server.
import type { Browser, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { expectFits, startGame, statMonsters } from './game';
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

// expectOwnArea checks the viewer's place (15.17): the character sits
// under the board's middle with live health and attack on it, coins are
// shown, and every hand card hangs partly past the bottom edge.
async function expectOwnArea(page: Page): Promise<void> {
  const own = page.locator('.player-bottom');
  const character = own.locator('.own-character > .card');
  await expect(character.getByLabel(/^\d+ health$/)).toBeVisible();
  await expect(character.getByLabel(/^\d+ attack$/)).toBeVisible();
  await expect(own.getByLabel(/^\d+ cents$/)).toBeVisible();
  const board = await boardBox(page);
  const card = await character.boundingBox();
  if (!board || !card) {
    throw new Error('no board or character');
  }
  expect(Math.abs(card.x + card.width / 2 - (board.x + board.width / 2))).toBeLessThan(2);
  const height = page.viewportSize()?.height ?? 0;
  const hand = await own.getByRole('group', { name: 'Your hand' }).locator('.hand-card').all();
  expect(hand.length).toBeGreaterThan(0);
  for (const c of hand) {
    const box = await c.boundingBox();
    expect(box?.y ?? height).toBeLessThan(height);
    expect((box?.y ?? 0) + (box?.height ?? 0)).toBeGreaterThan(height);
  }
}

// expectOpponents checks every other player's place (15.18): coins
// shown, live health and attack on the character, and the hand as one
// back per card, with its count. At the top, the character sits over
// the board's middle and the hand hangs past the top edge.
async function expectOpponents(page: Page): Promise<void> {
  const others = await page.locator('.player.opponent').all();
  expect(others.length).toBeGreaterThan(0);
  for (const other of others) {
    await expect(other.getByLabel(/^\d+ cents$/)).toBeVisible();
    const character = other.locator('.opp-character > .card');
    await expect(character.getByLabel(/^\d+ health$/)).toBeVisible();
    await expect(character.getByLabel(/^\d+ attack$/)).toBeVisible();
    const hand = other.getByRole('group', { name: /^Hand: \d+ cards$/ });
    const count = Number((await hand.locator('.pile-count').textContent()) ?? '');
    expect(count).toBeGreaterThan(0);
    await expect(hand.getByRole('img', { name: 'Face-down loot card' })).toHaveCount(Math.min(count, 8));
  }
  const top = page.locator('.player-top');
  if ((await top.count()) === 0) {
    return;
  }
  const board = await boardBox(page);
  const card = await top.locator('.opp-character > .card').boundingBox();
  if (!board || !card) {
    throw new Error('no board or character');
  }
  expect(Math.abs(card.x + card.width / 2 - (board.x + board.width / 2))).toBeLessThan(2);
  for (const back of await top.locator('.hidden-back').all()) {
    expect((await back.boundingBox())?.y ?? 0).toBeLessThan(0);
  }
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
    await expect(page.getByRole('group', { name: 'Shop' }).getByRole('img')).not.toHaveCount(0);
    await expect(page.getByRole('group', { name: 'Monsters' }).getByRole('img')).not.toHaveCount(0);
    await expectFits(page);
    await expectDecksClear(page);
    await expectOwnArea(page);
    await expectOpponents(page);
  }
  await shot(eve, 'table-2p');
});

test('coins show the coin image, or a drawn coin without it (15.38)', async ({ browser }) => {
  // Coy gets a coin image; Dot's fails to load, as without the images.
  const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8"><circle cx="4" cy="4" r="4" fill="#d90"/></svg>';
  let opened = 0;
  const [coy, dot] = await startGame(browser, ['Coy', 'Dot'], {
    open: (page) => {
      const served = opened++ === 0;
      return page.route('**/cardimg/assets/icons/coin.webp', (route) =>
        served ? route.fulfill({ body: svg, contentType: 'image/svg+xml' }) : route.abort(),
      );
    },
  });
  if (!coy || !dot) {
    throw new Error('no pages');
  }
  // Own coins and the opponent's (15.18) use the same icon.
  for (const coins of await coy.getByLabel(/^\d+ cents$/).all()) {
    const image = coins.locator('img.coin');
    await expect(image).toHaveAttribute('src', '/cardimg/assets/icons/coin.webp');
    await expect(image).toBeVisible();
    await expect(coins).not.toContainText('¢');
  }
  for (const coins of await dot.getByLabel(/^\d+ cents$/).all()) {
    await expect(coins.locator('.coin')).toHaveText('¢');
  }
});

test('the Attack icon shows the sword image, or a drawn sword without it (15.24)', async ({ browser }) => {
  // Sid gets a sword image; Uli's fails to load, as without the images.
  const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8"><rect width="8" height="8" fill="#555"/></svg>';
  let opened = 0;
  const [sid, uli] = await startGame(browser, ['Sid', 'Uli'], {
    open: (page) => {
      const served = opened++ === 0;
      return page.route('**/cardimg/assets/icons/sword.webp', (route) =>
        served ? route.fulfill({ body: svg, contentType: 'image/svg+xml' }) : route.abort(),
      );
    },
  });
  if (!sid || !uli) {
    throw new Error('no pages');
  }
  const attack = (page: Page) => page.getByRole('toolbar', { name: 'Turn actions' }).getByRole('button', { name: 'Attack' });
  const image = attack(sid).locator('img');
  await expect(image).toHaveAttribute('src', '/cardimg/assets/icons/sword.webp');
  await expect(image).toBeVisible();
  await expect(attack(uli).locator('img')).toHaveCount(0);
  await expect(attack(uli).locator('svg')).toBeVisible();
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
    await expectOwnArea(page);
    await expectOpponents(page);
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
    await expectOpponents(page);
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
        // Half a pixel either side for sub-pixel rounding (WebKit).
        expect(r.x).toBeGreaterThanOrEqual((area?.x ?? 0) - 0.5);
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

test('counters and changed monster stats show on the card, not in its zoom (15.23)', async ({ browser }) => {
  // Delirium or Stoney dealt into a slot would change the other monster (14.10).
  const [page] = await startGame(browser, ['Flo', 'Hux'], { bans: statMonsters });
  if (!page) {
    throw new Error('no pages');
  }
  const monster = page.getByRole('group', { name: 'Monsters' }).locator('.card').first();
  // The server sends each monster's stats; an untouched one shows the print only.
  await expect(monster).toBeVisible();
  await expect(monster.locator('.card-badge')).toHaveCount(0);

  // The first monster takes 1 damage and gets 2 counters; the second
  // gets +1 attack and is easier to hit, as the server would send it.
  // Only the screen changes; nobody acts on it.
  const hp = await page.evaluate(() => {
    interface Card {
      damage?: number;
      counters?: { count: number }[];
      stats?: { hp: number; dc: number; atk: number; printed_hp: number; printed_dc: number; printed_atk: number };
    }
    interface S {
      game: { view: { monsters: { top?: Card }[] } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const game = store.getState().game;
    const top = game?.view.monsters[0]?.top;
    const next = game?.view.monsters[1]?.top;
    if (!game || !top?.stats || !next?.stats) {
      throw new Error('no monster stats');
    }
    const hit: Card = { ...top, damage: 1, counters: [{ count: 2 }], stats: { ...top.stats, hp: top.stats.printed_hp - 1 } };
    const { printed_dc: dc, printed_atk: atk } = next.stats;
    const boosted: Card = { ...next, stats: { ...next.stats, dc: dc > 1 ? dc - 1 : dc + 1, atk: atk + 1 } };
    const monsters = game.view.monsters.map((s, i) => (i === 0 ? { ...s, top: hit } : i === 1 ? { ...s, top: boosted } : s));
    store.setState({ game: { ...game, view: { ...game.view, monsters } } });
    return top.stats.printed_hp - 1;
  });

  // The new health replaces the damage mark; dice and attack are as printed (C-03).
  await expect(monster.getByLabel(`${String(hp)} health`)).toBeVisible();
  await expect(monster.getByLabel('1 damage')).toHaveCount(0);
  await expect(monster.getByLabel(/ (dice|attack)$/)).toHaveCount(0);
  await expect(monster.getByLabel('2 counters')).toBeVisible();
  const second = page.getByRole('group', { name: 'Monsters' }).locator('.card').nth(1);
  await expect(second.getByLabel(/^\d dice$/)).toBeVisible();
  await expect(second.getByLabel(/^\d+ attack$/)).toBeVisible();
  await expect(second.getByLabel(/ health$/)).toHaveCount(0);
  await shot(page, 'monster-stats');

  // The zoom shows no badges, so none of its words hide (owner, 15.23).
  await monster.hover();
  const zoom = page.getByTestId('card-zoom').locator('.card.size-zoom');
  await expect(zoom).toBeVisible();
  await expect(zoom.locator('.card-badge')).toHaveCount(0);
  await shot(page, 'monster-stats-zoom');
});

test('a card with many abilities opens a dialog that uses the picked line (15.24)', async ({ browser }) => {
  const [page] = await startGame(browser, ['Lux', 'Ivy']);
  if (!page) {
    throw new Error('no pages');
  }
  // The viewer gets The Bone with 3 counters: its first three lines are
  // usable, the last needs 5 counters. Its lines come from the server
  // (C-05); a recorder stands in for sending.
  await page.evaluate(() => {
    interface Intent {
      player: number;
      kind: number;
      objects?: number[];
      choice?: number;
    }
    interface S {
      game: { step: number; view: { players: { in_play: unknown[] | null }[] }; allowed: Intent[] | null } | null;
      lobby: { table: { you: number } | null };
      select: (id: number | null) => void;
    }
    const w = window as unknown as {
      fourSouls: { store: { getState: () => S; setState: (s: object) => void } };
      picked?: Intent;
    };
    const { store } = w.fourSouls;
    const { game, lobby, select } = store.getState();
    const you = lobby.table?.you ?? 0;
    if (!game) {
      throw new Error('no game');
    }
    const bone = {
      id: 9001,
      card: 'the_bone',
      role: 2,
      charged: true,
      counters: [{ count: 3 }],
      abilities: [
        '↷: Put a counter on this.',
        'Remove 1 counter from this: Add +1 to a dice roll.',
        'Remove 2 counters from this: Deal 1 damage to a monster or player.',
        'Remove 5 counters from this: This becomes a soul and loses all abilities.',
      ],
    };
    const players = game.view.players.map((p, i) => (i === you ? { ...p, in_play: [...(p.in_play ?? []), bone] } : p));
    const allowed = [0, 1, 2].map((choice) => ({ player: you, kind: 7, objects: [9001], choice }));
    store.setState({
      game: { ...game, view: { ...game.view, players }, allowed },
      act: (intent: Intent) => {
        w.picked = intent;
        select(null);
      },
    });
  });

  await page.getByRole('group', { name: 'Your items' }).getByRole('button', { name: 'The Bone' }).click();
  const dialog = page.getByRole('dialog', { name: 'Use The Bone' });
  await expect(dialog).toBeVisible();
  await expect(dialog.locator('.card.size-zoom')).toBeVisible();
  const last = dialog.getByRole('button', { name: /Remove 5 counters/ });
  await expect(last).toBeDisabled();
  await shot(page, 'ability-dialog');

  await dialog.getByRole('button', { name: /Remove 2 counters/ }).click();
  await expect(dialog).toBeHidden();
  const picked = await page.evaluate(() => (window as unknown as { picked?: { objects?: number[]; choice?: number } }).picked);
  expect(picked).toMatchObject({ kind: 7, objects: [9001], choice: 2 });
});

test('souls show as crystals left of the character; a gained soul fills one more (15.39)', async ({ browser }) => {
  const [page] = await startGame(browser, ['Rho', 'Tau']);
  if (!page) {
    throw new Error('no pages');
  }
  const own = page.locator('.player-bottom');
  const other = page.locator('.player-top');
  // A slot per soul the game needs; none filled at the start (V-30).
  for (const area of [own, other]) {
    await expect(area.getByRole('group', { name: 'Souls: 0 of 4' })).toBeVisible();
    await expect(area.locator('.soul-slot')).toHaveCount(4);
    await expect(area.locator('.soul-slot.filled')).toHaveCount(0);
  }

  // setSouls gives each player souls, as the server would send them.
  // Only the screen changes; nobody acts on it.
  const setSouls = (souls: number[]) =>
    page.evaluate((next) => {
      interface S {
        game: { view: { players: { souls: number }[] } } | null;
      }
      const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
        .fourSouls;
      const game = store.getState().game;
      if (!game) {
        throw new Error('no game');
      }
      const players = game.view.players.map((p, i) => ({ ...p, souls: next[i] ?? p.souls }));
      store.setState({ game: { ...game, view: { ...game.view, players } } });
    }, souls);
  const you = (await own.getAttribute('data-zone')) === 'player-0' ? 0 : 1;
  const souls = (mine: number, theirs: number) => (you === 0 ? [mine, theirs] : [theirs, mine]);

  await setSouls(souls(1, 2));
  await expect(own.locator('.soul-slot.filled')).toHaveCount(1);
  await expect(other.locator('.soul-slot.filled')).toHaveCount(2);
  // A gained soul fills one more slot; the row keeps its four.
  await setSouls(souls(2, 2));
  await expect(own.getByRole('group', { name: 'Souls: 2 of 4' })).toBeVisible();
  await expect(own.locator('.soul-slot.filled')).toHaveCount(2);
  await expect(own.locator('.soul-slot')).toHaveCount(4);

  // The row sits left of the character (V-31); opponents get it smaller.
  for (const area of [own, other]) {
    const row = await area.locator('.souls').boundingBox();
    const card = await area.locator('.own-character .card, .opp-character .card').first().boundingBox();
    expect(row && card && row.x + row.width <= card.x).toBe(true);
  }
  const big = await own.locator('.soul-slot').first().boundingBox();
  const small = await other.locator('.soul-slot').first().boundingBox();
  expect(small?.width ?? 0).toBeLessThan(big?.width ?? 0);
  await shot(page, 'souls');
});

test('soul cards hide behind the crystals; hovering one zooms its card (15.56)', async ({ browser }) => {
  const [page] = await startGame(browser, ['Sigrid', 'Soren']);
  if (!page) {
    throw new Error('no pages');
  }
  const own = page.locator('.player-bottom');
  // Give the viewer two soul cards, as the server would send them: Mom
  // is worth two souls, Soul of Greed (a bonus soul) one. Only the
  // screen changes; nobody acts on it.
  await page.evaluate(() => {
    interface Card {
      id: number;
      card: string;
      role?: number;
      charged: boolean;
      souls?: number;
    }
    interface S {
      game: { view: { players: { id: number; souls: number; in_play: Card[] | null }[] } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const { game } = store.getState();
    if (!game) {
      throw new Error('no game');
    }
    const souls: Card[] = [
      { id: 9101, card: 'mom', role: 5, charged: false, souls: 2 },
      { id: 9102, card: 'soul_of_greed', role: 5, charged: false, souls: 1 },
    ];
    const you = Number(document.querySelector('.player-bottom')?.getAttribute('data-zone')?.replace('player-', ''));
    const players = game.view.players.map((p) =>
      p.id === you ? { ...p, souls: p.souls + 3, in_play: [...(p.in_play ?? []), ...souls] } : p,
    );
    store.setState({ game: { ...game, view: { ...game.view, players } } });
  });
  await expect(own.locator('.soul-slot.filled')).toHaveCount(3);
  await shot(page, 'soul-crystals');

  // No soul card in the items row (V-43).
  const items = own.getByRole('group', { name: 'Your items' });
  await expect(items.getByRole('img', { name: 'Mom!' })).toHaveCount(0);
  await expect(items.getByRole('img', { name: 'Soul Of Greed' })).toHaveCount(0);

  // Mom fills two crystals, Soul of Greed the third.
  const crystals = own.locator('.soul-slot.filled');
  await expect(crystals.nth(0)).toHaveAttribute('data-zone', 'obj-9101');
  await expect(crystals.nth(1)).toHaveAttribute('data-zone', 'obj-9101');
  await expect(crystals.nth(2)).toHaveAttribute('data-zone', 'obj-9102');

  // Hovering a crystal zooms its card beside it.
  const zoom = page.getByTestId('card-zoom').locator('.card.size-zoom');
  // Ready was clicked where the crystals now stand: move off first, so
  // the hover enters the crystal.
  await page.mouse.move(0, 0);
  await crystals.nth(1).hover();
  await expect(zoom).toHaveAttribute('aria-label', 'Mom!');
  await shot(page, 'soul-zoom');
  await crystals.nth(2).hover();
  await expect(zoom).toHaveAttribute('aria-label', 'Soul Of Greed');
  // An empty crystal zooms nothing.
  await own.locator('.soul-slot:not(.filled)').first().hover();
  await expect(zoom).toBeHidden();

  // An effect that aims at a soul picks it by its crystal: the other
  // player gains Pin, and a target prompt names it, as the server sends
  // it to the chooser. act is caught. (The choice dialog stands over
  // the own area, so the other player's soul is the one clicked.)
  await page.evaluate(() => {
    interface Player {
      id: number;
      souls: number;
      in_play: unknown[] | null;
    }
    interface S {
      game: { view: { players: Player[] } & Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: object) => void } } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    game.view.players = game.view.players.map((p) =>
      p.id === you
        ? p
        : { ...p, souls: p.souls + 1, in_play: [...(p.in_play ?? []), { id: 9103, card: 'pin', role: 5, charged: false, souls: 1 }] },
    );
    game.view.waiting = { kind: 4, player: you, purpose: 6, options: ['mom', 'pin'] };
    game.view.options = [
      { object: 9101, player: -1 },
      { object: 9103, player: -1 },
    ];
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    store.setState({
      act: (intent: unknown) => {
        (window as unknown as { sentChoice: unknown }).sentChoice = intent;
      },
    });
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  });
  const pin = page.locator('.player-top .soul-slot.filled');
  await expect(pin).toHaveAttribute('data-zone', 'obj-9103');
  await expect(pin).toHaveClass(/hl-pick/);
  await expect(page.getByRole('dialog', { name: 'Choose a target' })).toBeVisible();
  await shot(page, 'soul-pick');
  await pin.click();
  const sent = await page.evaluate(() => (window as unknown as { sentChoice: unknown }).sentChoice);
  expect(sent).toMatchObject({ kind: 4, choice: 1 });
});

test('a revealed hand shows face up; a hidden one stays backs (19.5)', async ({ browser }) => {
  const pages = await startGame(browser, ['Tam', 'Una', 'Vera', 'Wynn']);
  const una = pages[1];
  if (!una) {
    throw new Error('no pages');
  }
  // Reveal the hands at the top and right, as Guppy's Eye does: the
  // view sends their cards. Copies of the own hand stand in for them.
  // Only the screen changes; nobody acts on it.
  await una.evaluate(() => {
    interface Card {
      id: number;
      card: string;
    }
    interface S {
      game: { view: { players: { id: number; hand?: Card[]; hand_size: number }[] } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const { game } = store.getState();
    if (!game) {
      throw new Error('no game');
    }
    const own = game.view.players[1]?.hand ?? [];
    const revealed = (p: number) => [0, 1, 2].map((i) => ({ ...own[i % own.length], id: 9000 + p * 10 + i }) as Card);
    // Una is player 1: Wynn (3) sits at the top, Tam (0) at the right.
    const players = game.view.players.map((p) => (p.id === 0 || p.id === 3 ? { ...p, hand: revealed(p.id), hand_size: 3 } : p));
    store.setState({ game: { ...game, view: { ...game.view, players } } });
  });
  for (const name of ['Wynn', 'Tam']) {
    const area = una.getByRole('region', { name, exact: true });
    const hand = area.getByRole('group', { name: 'Revealed hand: 3 cards' });
    await expect(hand.locator('.card')).toHaveCount(3);
    await expect(hand.getByRole('img', { name: 'Face-down loot card' })).toHaveCount(0);
    await expect(hand).toHaveAttribute('data-zone', /^hand-\d$/);
  }
  // The hidden hand on the left keeps its backs and count.
  const vera = una.getByRole('region', { name: 'Vera', exact: true });
  await expect(vera.getByRole('group', { name: /^Hand: \d+ cards$/ })).toBeVisible();
  await expect(vera.getByRole('group', { name: /^Revealed hand/ })).toHaveCount(0);
  // A revealed card zooms on hover, like any table card.
  const zoom = una.getByTestId('card-zoom');
  await una.getByRole('region', { name: 'Wynn', exact: true }).getByRole('group', { name: /^Revealed hand/ }).locator('.card').last().hover();
  await expect(zoom).toBeVisible();
  await shot(una, 'revealed-hand-zoom');
  await una.mouse.move(0, 0);
  await expectFits(una);
  await shot(una, 'revealed-hand');
});
