// Info panels (7.9): history, turn badge, what is left this turn.
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, canEndTurn, click, gameOver, playLoot, settle, turnNumber } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

test('the history shows only what the viewer may see', async ({ browser }) => {
  const names = ['Wes', 'Xan'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const history = (p: typeof active) => p.getByRole('complementary', { name: 'History', exact: true });
  // The full history opens from the strip (15.22).
  for (const page of pages) {
    await page.getByRole('button', { name: 'History', exact: true }).click();
  }
  for (const [i, page] of pages.entries()) {
    const other = names.at(1 - i) ?? '';
    const text = (await history(page).textContent()) ?? '';
    // Others' loots never name the card.
    expect(text).not.toMatch(new RegExp(`${other} loots (?!a card)`));
  }
  // The active player looted: their own history names the card.
  await expect(history(active)).toContainText(/You loot \S/);
  const other = pages.find((p) => p !== active);
  const activeName = names[pages.indexOf(active)] ?? '';
  await expect(history(other ?? active)).toContainText(`${activeName} loots a card`);
  await shot(active, 'history');
});

test('the turn badge, what is left, and hiding the history', async ({ browser }) => {
  const names = ['Yul', 'Zed'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const activeName = names[pages.indexOf(active)] ?? '';
  for (const page of pages) {
    await expect(page.getByRole('status', { name: 'Turn' })).toHaveText(new RegExp(`^Turn \\d+ · ${activeName} · \\w`));
  }
  // Own: loot plays by the name, buys and attacks as badges on their icons (15.24).
  await expect(active.getByLabel('Loot plays left')).toHaveText(/^Loot \d+$/);
  const tools = active.getByRole('toolbar', { name: 'Turn actions' });
  await expect(tools.getByRole('button', { name: 'Buy', exact: true }).getByLabel(/^\d+ left$/)).toBeVisible();
  await expect(tools.getByRole('button', { name: 'Attack', exact: true }).getByLabel(/^\d+ left$/)).toBeVisible();
  await shot(active, 'turn-left');
  // A big count reads "99+" and stays clear of the character card.
  await active.evaluate(() => {
    interface S {
      game: { view: { players: { purchases: number }[]; turn: { active: number } } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const { game } = store.getState();
    if (!game) {
      throw new Error('no game');
    }
    const players = game.view.players.map((p, i) => (i === game.view.turn.active ? { ...p, purchases: 100 } : p));
    store.setState({ game: { ...game, view: { ...game.view, players } } });
  });
  const big = tools.getByLabel('100 left');
  await expect(big).toHaveText('99+');
  const badge = await big.boundingBox();
  const card = await active.locator('.own-character .card').boundingBox();
  expect((badge?.x ?? 0) + (badge?.width ?? 0)).toBeLessThanOrEqual(card?.x ?? 0);
  // Others see the active player's whole allowance by their name.
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(other.getByLabel('Left this turn')).toHaveText(/Loot \d+ · Attacks \d+ · Shop \d+/);
  await expect(other.getByLabel(/^\d+ left$/)).toHaveCount(0);

  // The full history is closed; the strip opens it, Hide closes it.
  const panel = active.getByRole('complementary', { name: 'History', exact: true });
  const toggle = active.getByRole('button', { name: 'History', exact: true });
  await expect(panel).toHaveCount(0);
  await toggle.click();
  await expect(panel).toBeVisible();
  await expect(toggle).toHaveAttribute('aria-pressed', 'true');
  await panel.getByRole('button', { name: 'Hide' }).click();
  await expect(panel).toHaveCount(0);
  await toggle.click();
  await toggle.click();
  await expect(panel).toHaveCount(0);
});

test('a resolved loot card goes to the history strip (15.22)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ada', 'Ben'], { url: fast() });
  const active = await activePage(pages, click);
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  const name = (await loot.getAttribute('aria-label')) ?? '';
  const turn = await turnNumber(active);
  await playLoot(active, loot, click);
  // Targets, passes and anything the card asks, until the stack is empty.
  // The card may end the turn (The Fool, a death): then the next one
  // counts too (19.8).
  const stackRows = active.getByRole('complementary', { name: 'Stack' }).locator('.stack-row');
  const moved = async () => (await canEndTurn(active)) || (await turnNumber(active)) !== turn || (await gameOver(active));
  await settle(pages, click, async () => (await stackRows.count()) === 0 && (await moved()));
  for (const page of pages) {
    const strip = page.getByRole('complementary', { name: 'History strip' });
    // The loot is a played spell: a bolt badge (H-08).
    const played = strip.getByLabel(`plays ${name}`).first();
    await expect(played).toBeVisible();
    await expect(played.locator('.strip-spell')).toBeAttached();
  }
});

test('the history strip: icons, badges, and a hover with card and owner (15.22)', async ({ browser }) => {
  const pages = await startGame(browser, ['Cid', 'Dee'], { url: fast() });
  const active = await activePage(pages, click);
  // Resolved items as the server would send them, the newest last: Blank
  // Rune, its roll and its result; a loot card, an attack roll, damage. They replace what the deal
  // resolved so far. Only the screen changes.
  await active.evaluate(() => {
    interface Item {
      id: number;
      kind: number;
      controller: number;
      source: number;
      card?: string;
      roll?: number;
      amount?: number;
      attack?: boolean;
      text?: string;
      target_card?: string;
      target: { player: number; object?: number; is_player?: boolean };
      ability: { card: string; index: number };
      roll_for: { card: string; index: number };
    }
    interface Entry {
      key: number;
      turn: number;
      icon: string;
      text: string;
      card?: string;
      item?: Item;
    }
    interface S {
      game: { view: { turn: { active: number; number: number } } } | null;
      history: Entry[];
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const s = store.getState();
    if (!s.game) {
      throw new Error('no game');
    }
    const me = s.game.view.turn.active;
    const turn = s.game.view.turn.number;
    const other = 1 - me;
    const none = { card: '', index: 0 };
    const rune = { card: 'blank_rune', index: 0 };
    const items: Item[] = [
      { id: 98, kind: 0, controller: me, source: 0, card: 'blank_rune', text: 'Roll.', target: { player: -1 }, ability: rune, roll_for: none },
      { id: 99, kind: 3, controller: me, source: 0, roll: 1, target: { player: -1 }, ability: none, roll_for: rune },
      {
        id: 100,
        kind: 2,
        controller: me,
        source: 0,
        card: 'blank_rune',
        text: '1: Each player gains 1¢.',
        target: { player: -1 },
        ability: rune,
        roll_for: none,
      },
      {
        id: 101,
        kind: 0,
        controller: me,
        source: 0,
        card: 'bomb',
        text: 'Deal 1 damage to a monster or player.',
        target: { player: other, is_player: true },
        ability: { card: 'bomb', index: 0 },
        roll_for: none,
      },
      {
        id: 102,
        kind: 3,
        controller: me,
        source: 0,
        roll: 4,
        attack: true,
        target_card: 'fly',
        target: { player: -1, object: 999 },
        ability: none,
        roll_for: none,
      },
      { id: 103, kind: 4, controller: me, source: 0, amount: 2, attack: true, target: { player: other, is_player: true }, ability: none, roll_for: none },
    ];
    const history = items.map((item, i) => ({ key: i + 1, turn, icon: '✦', text: '', ...(item.card ? { card: item.card } : {}), item, headline: true }));
    store.setState({ history });
  });

  const strip = active.getByRole('complementary', { name: 'History strip' });
  const icons = strip.locator('.strip-item');
  await expect(icons).toHaveCount(6);
  // Newest on top; a sword for the attack, a bolt for the loot (H-08).
  await expect(icons.nth(0)).toHaveAttribute('aria-label', /^Attack: (Cid|Dee) deals 2 damage to (Cid|Dee)$/);
  await expect(icons.nth(0).locator('.strip-attack')).toBeAttached();
  await expect(icons.nth(1).locator('.die')).toHaveAttribute('aria-label', 'Rolled 4');
  await expect(icons.nth(1).locator('.strip-attack')).toBeAttached();
  await expect(icons.nth(2)).toHaveAttribute('aria-label', /^(Cid|Dee) plays Bomb!?$/);
  await expect(icons.nth(2).locator('.strip-spell')).toBeAttached();

  // Hover the loot icon: its card large, the used ability, the owner (S-11, H-09).
  await icons.nth(2).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', /^Bomb/);
  await expect(zoom.locator('.card-zoom-ability')).toHaveText('Deal 1 damage to a monster or player.');
  await expect(zoom.locator('.card-zoom-owner')).toHaveText(/^(Cid|Dee)$/);
  await shot(active, 'history-strip-hover');

  // Hover the attack roll: the monster it was against, and the line.
  await icons.nth(1).hover();
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', /^Fly/);
  await expect(zoom.locator('.card-zoom-line')).toHaveText(/^Attack: (Cid|Dee) rolls 4 against Fly$/);
  await shot(active, 'history-strip-attack');

  // Blank Rune: the play says only "Roll."; its result says what happened.
  await icons.nth(5).hover();
  await expect(zoom.locator('.card-zoom-ability')).toHaveText('Roll.');
  await icons.nth(3).hover();
  await expect(zoom.locator('.card-zoom-ability')).toHaveText('1: Each player gains 1¢.');
  await shot(active, 'history-strip-result');
  await active.mouse.move(1000, 420);
  await expect(zoom).toBeHidden();
  await shot(active, 'history-strip');
});

test('damage a player takes names its card in the history strip (15.47)', async ({ browser }) => {
  const pages = await startGame(browser, ['Bob', 'Wes'], { url: fast() });
  const active = await activePage(pages, click);
  // Boom Fly died: each player controls the damage they take from it, and
  // the fly is no longer on the table; only source_card names it.
  await active.evaluate(() => {
    interface S {
      game: { view: { turn: { active: number; number: number } } } | null;
      history: unknown[];
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const s = store.getState();
    if (!s.game) {
      throw new Error('no game');
    }
    const turn = s.game.view.turn.number;
    const none = { card: '', index: 0 };
    const history = [1, 0].map((p, i) => ({
      key: i + 1,
      turn,
      icon: '✦',
      text: '',
      item: {
        id: 40 + i,
        kind: 4,
        controller: p,
        source: 307,
        source_card: 'boom_fly',
        amount: 1,
        label: 'damage',
        target: { player: p, is_player: true },
        ability: none,
        roll_for: none,
      },
      headline: true,
    }));
    store.setState({ history });
  });

  const strip = active.getByRole('complementary', { name: 'History strip' });
  const icons = strip.locator('.strip-item');
  await expect(icons).toHaveCount(2);
  await expect(icons.nth(0)).toHaveAttribute('aria-label', /^Boom Fly deals 1 damage to (Bob|Wes)$/);
  await expect(icons.nth(1)).toHaveAttribute('aria-label', /^Boom Fly deals 1 damage to (Bob|Wes)$/);

  // Hover: the fly large, and the line under it.
  await icons.nth(0).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', /^Boom Fly/);
  await expect(zoom.locator('.card-zoom-line')).toHaveText(/^Boom Fly deals 1 damage to (Bob|Wes)$/);
  await shot(active, 'history-strip-taken-damage');
});

test('a changed roll shows the die and the result in the history (19.1)', async ({ browser }) => {
  const pages = await startGame(browser, ['Jin', 'Kaz'], { url: fast() });
  const active = await activePage(pages, click);
  // Polycephalus took damage: the next attack roll has -2. The die says
  // 6, the game uses 4 (the play-test record of 19.1).
  await active.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { active: number }; monsters: { top?: { id: number; card: string; stats?: unknown } }[] } } | null;
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !monster) {
      throw new Error('no game');
    }
    const p = game.view.turn.active;
    const none = { card: '', index: 0 };
    const item = { id: 900, kind: 3, controller: p, source: 0, roll: 6, label: 'attack roll', attack: true, target: { object: monster.id, player: 0 }, ability: none, roll_for: none, target_card: monster.card };
    s.receive({
      type: 'update',
      data: {
        ...game,
        step: 900,
        events: [
          { kind: 'stack_resolved', player: p, amount: 6, text: 'attack roll', stack_id: 900, item },
          { kind: 'roll_resolved', player: p, amount: 4, text: 'attack' },
        ],
      },
    });
  });

  const strip = active.getByRole('complementary', { name: 'History strip' });
  const icon = strip.locator('.strip-item.stack-roll').first();
  await expect(icon).toHaveAttribute('aria-label', /rolls \d against /);
  await icon.hover();
  await shot(active, 'history-roll-changed-strip');
  await active.mouse.move(1000, 420);
  await active.getByRole('button', { name: 'History', exact: true }).click();
  const history = active.getByRole('complementary', { name: 'History', exact: true });
  await expect(history).toContainText(/Attack: \S+ rolls 6/);
  await shot(active, 'history-roll-changed');

  await expect(icon).toHaveAttribute('aria-label', /rolls 4 against /);
  await expect(icon.locator('.die')).toHaveAttribute('aria-label', 'Rolled 4');
  await expect(history).toContainText(/Attack: \S+ rolls 6 → 4 against /);
});

test('the history names every play: a buy, an attack and a kill (15.88)', async ({ browser }) => {
  const names = ['Lev', 'Mia'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const activeName = names[pages.indexOf(active)] ?? '';
  const other = pages.find((p) => p !== active) ?? active;
  const turn = await turnNumber(active);
  // A real purchase: with the starting cents it most likely fails, which
  // the history names too (15.46).
  await active.keyboard.press('b');
  await settle(pages, click, async () => (await canEndTurn(active)) || (await turnNumber(active)) !== turn || (await gameOver(active)));
  for (const page of pages) {
    await page.getByRole('button', { name: 'History', exact: true }).click();
  }
  const history = (p: typeof active) => p.getByRole('complementary', { name: 'History', exact: true });
  await expect(history(active)).toContainText('You declare a purchase');
  await expect(history(active)).toContainText(/You (buy \S|can't buy: the price is \d+¢)/);
  await expect(history(other)).toContainText(`${activeName} declares a purchase`);
  await expect(history(other)).toContainText(new RegExp(`${activeName} (buys \\S|can't buy)`));

  // A kill needs lucky dice; its events come as the engine sends them.
  await active.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { active: number }; monsters: { top?: { id: number; card: string; stats?: unknown } }[] } } | null;
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !monster) {
      throw new Error('no game');
    }
    const p = game.view.turn.active;
    s.receive({
      type: 'update',
      data: {
        ...game,
        step: 900,
        events: [
          { kind: 'attack_declared', player: p },
          { kind: 'attack_target', player: p, object: monster.id, card: monster.card },
          { kind: 'died', player: -1, object: 901, card: monster.card, prev: monster.id, by: p + 1 },
          { kind: 'gained_cents', player: p, amount: 3 },
        ],
      },
    });
  });
  await expect(history(active)).toContainText('You declare an attack');
  const kill = history(active).getByText(/^You kill \S/);
  await expect(kill).toBeVisible();
  await expect(history(active)).toContainText(/You attack \S/);
  await expect(history(active)).toContainText('You gain 3¢');
  await shot(active, 'history-every-play');
  await shot(other, 'history-every-play-other');
  // The line zooms the dead monster.
  await kill.hover();
  await expect(active.getByTestId('card-zoom').locator('.card.size-zoom')).toBeVisible();
  await shot(active, 'history-every-play-zoom');
});

test('a buy and a kill each add one tile to the history strip (15.89)', async ({ browser }) => {
  const pages = await startGame(browser, ['Nia', 'Oto'], { url: fast() });
  const active = await activePage(pages, click);
  const strip = active.getByRole('complementary', { name: 'History strip' });
  const tiles = strip.locator('.strip-item');
  const before = await tiles.count();
  // A kill resolves a death stack item first, then the monster dies; a
  // buy comes with its gain. Each is one tile.
  await active.evaluate(() => {
    interface Store {
      getState: () => {
        game: {
          view: {
            turn: { active: number };
            monsters: { top?: { id: number; card: string; stats?: unknown } }[];
            shop: { top?: { id: number; card: string } }[];
          };
        } | null;
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    const item = game?.view.shop.find((m) => m.top)?.top;
    if (!game || !monster || !item) {
      throw new Error('no game');
    }
    const p = game.view.turn.active;
    const none = { card: '', index: 0 };
    const death = { id: 950, kind: 5, controller: -1, source: 0, text: 'death', target: { player: -1, object: monster.id }, ability: none, roll_for: none };
    s.receive({
      type: 'update',
      data: {
        ...game,
        step: 900,
        events: [
          { kind: 'stack_resolved', player: -1, text: 'death', stack_id: 950, item: death },
          { kind: 'died', player: -1, object: 951, card: monster.card, prev: monster.id, by: p + 1 },
          { kind: 'gained_cents', player: p, amount: 3 },
          { kind: 'purchase_declared', player: p },
          { kind: 'purchased', player: p, object: 952, card: item.card },
          { kind: 'gained_item', player: p, object: 952, card: item.card },
        ],
      },
    });
  });
  await expect(tiles).toHaveCount(before + 2);
  const buy = tiles.nth(0);
  const kill = tiles.nth(1);
  await expect(buy).toHaveAttribute('aria-label', /^You buy \S/);
  await expect(buy.locator('.strip-buy')).toBeAttached();
  await expect(buy.locator('.card')).toBeAttached();
  await expect(kill).toHaveAttribute('aria-label', /^You kill \S/);
  await expect(kill.locator('.strip-kill')).toBeAttached();
  await expect(kill.locator('.card')).toBeAttached();
  await shot(active, 'history-strip-buy-kill');
  // Hover: the dead monster large, with the line under it.
  await kill.hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card.size-zoom')).toBeVisible();
  await expect(zoom.locator('.card-zoom-line')).toHaveText(/^You kill \S/);
  await shot(active, 'history-strip-kill-hover');
});

test('a damage line in the full history names who and whom (19.52)', async ({ browser }) => {
  const names = ['Pia', 'Quo'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const activeName = names[pages.indexOf(active)] ?? '';
  const otherName = names.find((x) => x !== activeName) ?? '';
  // Damage, a death and its soul step resolve as the engine sends them.
  await active.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { active: number }; monsters: { top?: { id: number; card: string; stats?: unknown } }[] } } | null;
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !monster) {
      throw new Error('no game');
    }
    const p = game.view.turn.active;
    const none = { card: '', index: 0 };
    const base = { controller: p, source: 0, ability: none, roll_for: none };
    const damage = { ...base, id: 960, kind: 4, amount: 1, label: 'damage', target: { player: 1 - p, is_player: true } };
    const hit = { ...base, id: 961, kind: 4, amount: 2, label: 'damage', target_card: monster.card, target: { player: -1, object: monster.id } };
    s.receive({
      type: 'update',
      data: {
        ...game,
        step: 900,
        events: [
          { kind: 'stack_resolved', player: p, text: 'damage', stack_id: 960, item: damage },
          { kind: 'stack_resolved', player: p, text: 'damage', stack_id: 961, item: hit },
        ],
      },
    });
  });
  await active.getByRole('button', { name: 'History', exact: true }).click();
  const history = active.getByRole('complementary', { name: 'History', exact: true });
  await expect(history).toContainText(`${activeName} deals 1 damage to ${otherName}`);
  await expect(history).toContainText(new RegExp(`${activeName} deals 2 damage to \\S`));
  await expect(history).not.toContainText('Roll — damage');
  await shot(active, 'history-damage-lines');
});

test('a strip tile shows what the play did: target, damage, kill and rewards (15.90)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ray', 'Sue'], { url: fast() });
  const active = await activePage(pages, click);
  // Bomb! at a monster, its damage, then the death and the rewards, as
  // the engine sends them: each resolution in its own update.
  await active.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { active: number }; monsters: { top?: { id: number; card: string; stats?: unknown } }[] } } | null;
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const game = structuredClone(store.getState().game);
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !monster) {
      throw new Error('no game');
    }
    const p = game.view.turn.active;
    const none = { card: '', index: 0 };
    const target = { player: -1, object: monster.id };
    const bomb = { id: 970, kind: 0, controller: p, source: 0, card: 'bomb', target, target_card: monster.card, ability: { card: 'bomb', index: 0 }, roll_for: none };
    const hit = { id: 971, kind: 4, controller: p, source: 0, source_card: 'bomb', amount: 1, label: 'damage', target, target_card: monster.card, ability: none, roll_for: none };
    const death = { id: 972, kind: 5, controller: -1, source: 0, label: 'death', target, ability: none, roll_for: none };
    const updates = [
      [{ kind: 'stack_resolved', player: p, card: 'bomb', stack_id: 970, item: bomb }],
      [{ kind: 'stack_resolved', player: p, text: 'damage', stack_id: 971, item: hit }],
      [
        { kind: 'stack_resolved', player: -1, text: 'death', stack_id: 972, item: death },
        { kind: 'died', player: -1, object: 973, card: monster.card, prev: monster.id, by: p + 1 },
        { kind: 'rewards_gained', player: p, object: 973, card: monster.card },
        { kind: 'gained_cents', player: p, amount: 3 },
      ],
    ];
    for (const [i, events] of updates.entries()) {
      store.getState().receive({ type: 'update', data: { ...game, step: 900 + i, events } });
    }
  });
  const strip = active.getByRole('complementary', { name: 'History strip' });
  const kill = strip.locator('.strip-item').nth(0);
  const damage = strip.locator('.strip-item').nth(1);
  const play = strip.locator('.strip-item').nth(2);
  await expect(kill).toHaveAttribute('aria-label', /^You kill \S/);
  const zoom = active.getByTestId('card-zoom');
  const results = zoom.getByRole('list', { name: 'What it did' }).getByRole('listitem');

  // The kill: the killer large, then the monster with its skull and the cents.
  await kill.hover();
  await expect(results).toHaveCount(2);
  await expect(results.nth(0)).toHaveAttribute('aria-label', /dies$/);
  await expect(results.nth(0).locator('.result-skull')).toBeAttached();
  await expect(results.nth(0).locator('.card')).toBeAttached();
  await expect(results.nth(1)).toHaveAttribute('aria-label', '+3¢');
  await expect(zoom.locator('.card-zoom-line')).toHaveText(/^You kill \S/);
  await shot(active, 'history-popout-kill');

  // The damage: Bomb! large, the monster with a starburst and the number.
  await damage.hover();
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', /^Bomb/);
  await expect(results).toHaveCount(1);
  await expect(results.nth(0)).toHaveAttribute('aria-label', /: 1 damage$/);
  await expect(results.nth(0).locator('.result-burst')).toHaveText('−1');
  await shot(active, 'history-popout-damage');

  // The play: Bomb! large with its ability, an arrow to its target.
  await play.hover();
  await expect(zoom.locator('.card-zoom-ability')).toBeVisible();
  await expect(results).toHaveCount(1);
  await expect(results.nth(0).locator('.card')).toBeAttached();
  await shot(active, 'history-popout-play');
});
