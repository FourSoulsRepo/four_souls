// The history strip like the reference's (20.21, ADR 027 "History" 11):
// rolls and damage at half height, a frame in the owner's color, a new
// tile sliding in from the left, nothing moving under the pointer, the
// strip as one hit area with a sound per tile and a dimmed board,
// swords for an attack, targets in rows, a fading foot. Each test feeds
// the strip its own entries or updates, so none depends on the deal
// (CLAUDE.md 5.16).
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, lineRe, tr } from './locale';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';

const off = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

// names are the players' names; a repeat (--repeat-each) gets its own.
function names(base: string[]): string[] {
  const r = test.info().repeatEachIndex;
  return base.map((n) => (r > 0 ? `${n}${String(r)}` : n));
}

const expoOut = 'cubic-bezier(0.16, 1, 0.3, 1)';

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
  label?: { key: string };
  source_card?: string;
  target_card?: string;
  target: { player: number; object?: number; is_player?: boolean };
  ability: { card: string; index: number };
  roll_for: { card: string; index: number };
}

interface View {
  turn: { active: number; number: number };
  players: { id: number; character: { id: number; card: string } }[];
  monsters: { top?: { id: number; card: string; stats?: unknown } }[];
}

interface Store {
  getState: () => {
    game: { step: number; view: View } | null;
    history: { key: number; turn: number }[];
    receive: (env: unknown) => void;
  };
  setState: (s: { history: unknown[] }) => void;
}

// Scene is what a test needs to know about the table it feeds.
interface Scene {
  me: number;
  hero: string;
  monster: string;
}

// attackKill sends, as the engine does, a loot play of Bomb! at another
// player, then an attack on a monster: its roll, its damage and the
// kill, each resolution in its own update.
async function attackKill(page: Page): Promise<Scene> {
  return page.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const game = structuredClone(store.getState().game);
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !monster) {
      throw new Error('no game');
    }
    const p = game.view.turn.active;
    const hero = game.view.players.find((pl) => pl.id === p)?.character;
    if (!hero) {
      throw new Error('no character');
    }
    const none = { card: '', index: 0 };
    const target = { player: -1, object: monster.id };
    const bomb: Item = {
      id: 980,
      kind: 0,
      controller: p,
      source: 0,
      card: 'bomb',
      target: { player: 1 - p, is_player: true },
      ability: { card: 'bomb', index: 0 },
      roll_for: none,
    };
    const roll: Item = { id: 981, kind: 3, controller: p, source: 0, roll: 5, attack: true, label: { key: 'event.attack_roll' }, target, target_card: monster.card, ability: none, roll_for: none };
    const hit: Item = {
      id: 982,
      kind: 4,
      controller: p,
      source: hero.id,
      source_card: hero.card,
      amount: 1,
      attack: true,
      label: { key: 'event.combat_damage' },
      target,
      target_card: monster.card,
      ability: none,
      roll_for: none,
    };
    const death: Item = { id: 983, kind: 5, controller: -1, source: 0, label: { key: 'event.death' }, target, ability: none, roll_for: none };
    const updates = [
      [{ kind: 'stack_resolved', player: p, card: 'bomb', stack_id: 980, item: bomb }],
      [{ kind: 'stack_resolved', player: p, label: { key: 'event.attack_roll' }, stack_id: 981, item: roll, amount: 5 }],
      [{ kind: 'stack_resolved', player: p, label: { key: 'event.combat_damage' }, stack_id: 982, item: hit }],
      [
        { kind: 'stack_resolved', player: -1, label: { key: 'event.death' }, stack_id: 983, item: death },
        { kind: 'died', player: -1, object: 984, card: monster.card, prev: monster.id, by: p + 1 },
      ],
    ];
    for (const [i, events] of updates.entries()) {
      store.getState().receive({ type: 'update', data: { ...game, step: game.step + 1 + i, events } });
    }
    return { me: p, hero: hero.card, monster: monster.card };
  });
}

// addTiles puts n more plain tiles on top of the strip at once, as an
// update's headlines would: buys of Bomb! by the viewer. outcome, if
// set, is the popout's targets of the newest.
async function addTiles(page: Page, n: number, outcome: string[] = []): Promise<void> {
  await page.evaluate(
    ({ n, outcome }) => {
      const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
      const s = store.getState();
      const turn = s.game?.view.turn.number ?? 1;
      const me = s.game?.view.turn.active ?? 0;
      const key = (s.history.at(-1)?.key ?? 0) + 1;
      const tiles = Array.from({ length: n }, (_, i) => ({
        key: key + i,
        turn,
        icon: '¢',
        text: `tile ${String(key + i)}`,
        card: 'bomb',
        headline: true,
        badge: 'buy',
        owner: me,
        ...(i === n - 1 && outcome.length > 0 ? { from: 'bomb', outcome: outcome.map((c) => ({ text: c, card: c })) } : {}),
      }));
      store.setState({ history: [...s.history, ...tiles] });
    },
    { n, outcome },
  );
}

function strip(page: Page): Locator {
  return page.getByRole('complementary', { name: tr('history.strip') });
}

async function box(l: Locator) {
  const b = await l.boundingBox();
  if (!b) {
    throw new Error('no box');
  }
  return b;
}

// spy records every sound played and every WAAPI animation started
// from now on, in window.played and window.anims; each still runs.
async function spy(page: Page): Promise<void> {
  await page.evaluate(() => {
    const w = window as unknown as { played: string[]; anims: { cls: string; ms: number; easing: string }[] };
    w.played = [];
    w.anims = [];
    const play = HTMLMediaElement.prototype.play;
    HTMLMediaElement.prototype.play = function (this: HTMLMediaElement) {
      w.played.push(new URL(this.src).pathname);
      return play.call(this);
    };
    const animate = Element.prototype.animate;
    Element.prototype.animate = function (this: Element, keys: Keyframe[] | PropertyIndexedKeyframes | null, o?: number | KeyframeAnimationOptions) {
      const t = typeof o === 'number' ? { duration: o } : (o ?? {});
      w.anims.push({ cls: typeof this.className === 'string' ? this.className : '', ms: Number(t.duration ?? 0), easing: t.easing ?? '' });
      return animate.call(this, keys, o);
    };
  });
}

function played(page: Page, sound: string): Promise<number> {
  return page.evaluate((s) => (window as unknown as { played: string[] }).played.filter((p) => p === `/cardimg/assets/sounds/${s}.mp3`).length, sound);
}

// anims are the WAAPI animations started on elements whose class has
// cls, with their time in ms and their curve.
function anims(page: Page, cls: string): Promise<{ ms: number; easing: string }[]> {
  return page.evaluate(
    (c) =>
      (window as unknown as { anims: { cls: string; ms: number; easing: string }[] }).anims
        .filter((a) => a.cls.split(' ').includes(c))
        .map((a) => ({ ms: a.ms, easing: a.easing })),
    cls,
  );
}

test('rolls and damage are half height; a thin frame has the owner color (20.21)', async ({ browser }) => {
  const pages = await startGame(browser, names(['Ada', 'Ben']), { url: off() });
  const active = await activePage(pages, click);
  await attackKill(active);
  const tiles = strip(active).locator('.strip-item');
  const [kill, damage, roll, play] = [tiles.nth(0), tiles.nth(1), tiles.nth(2), tiles.nth(3)];
  await expect(kill.locator('.strip-kill')).toBeAttached();
  await expect(damage).toHaveClass(/strip-half/);
  await expect(roll).toHaveClass(/strip-half/);
  await expect(kill).toHaveClass(/strip-full/);
  await expect(play).toHaveClass(/strip-full/);
  expect((await box(roll)).height).toBeLessThan((await box(play)).height * 0.7);
  // The frame: the color the popout names its owner in.
  for (const tile of [kill, damage, roll, play]) {
    await expect(tile).toHaveClass(/owned/);
  }
  const frame = await roll.evaluate((el) => getComputedStyle(el, '::after').borderTopColor);
  await roll.hover();
  const owner = await active.getByTestId('card-zoom').locator('.card-zoom-owner').evaluate((el) => getComputedStyle(el).color);
  expect(frame).toBe(owner);
});

test("an attack's popout: the attacker large, swords, the monster with its marks (20.21)", async ({ browser }) => {
  const pages = await startGame(browser, names(['Cal', 'Dot']), { url: off() });
  const active = await activePage(pages, click);
  const scene = await attackKill(active);
  const tiles = strip(active).locator('.strip-item');
  const [kill, damage, roll, play] = [tiles.nth(0), tiles.nth(1), tiles.nth(2), tiles.nth(3)];
  await expect(kill.locator('.strip-kill')).toBeAttached();
  const zoom = active.getByTestId('card-zoom');
  const results = zoom.getByRole('list', { name: tr('zoom.outcome') }).getByRole('listitem');

  // The roll: the attacker's character large, swords, the monster with
  // the damage the attack dealt and its death; the die in the caption.
  await roll.hover();
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', new RegExp(`^${lineRe(cardName(scene.hero)).source}`));
  await expect(zoom.getByTestId('zoom-swords')).toBeVisible();
  await expect(results).toHaveCount(1);
  await expect(results.nth(0).locator('.card')).toHaveAttribute('aria-label', new RegExp(`^${lineRe(cardName(scene.monster)).source}`));
  await expect(results.nth(0).locator('.result-skull')).toBeAttached();
  await expect(zoom.locator('.card-zoom-line')).toContainText('5');
  await shot(active, 'history-popout-attack-roll');

  // The damage and the kill: swords too.
  await damage.hover();
  await expect(zoom.getByTestId('zoom-swords')).toBeVisible();
  await kill.hover();
  await expect(zoom.getByTestId('zoom-swords')).toBeVisible();
  // A play: an arrow.
  await play.hover();
  await expect(zoom.getByTestId('zoom-arrow')).toBeVisible();
  await expect(zoom.getByTestId('zoom-swords')).toHaveCount(0);
});

test('a new tile slides in from the left while the rest move down; speed scales it (20.21)', async ({ browser }) => {
  const pages = await startGame(browser, names(['Eli', 'Fay']), { url: appURL() });
  const page = pages[0];
  if (!page) {
    throw new Error('no page');
  }
  await addTiles(page, 3);
  const tiles = strip(page).locator('.strip-item');
  await expect(tiles.first()).toHaveAttribute('aria-label', /tile/);
  await page.mouse.move(900, 700);
  await page.waitForTimeout(1200);
  await spy(page);
  // Just after it comes, the new tile is still left of the one below it,
  // which keeps its x; the list moves down, so nothing goes up.
  const xs = await page.evaluate(async () => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const key = (s.history.at(-1)?.key ?? 0) + 1;
    store.setState({ history: [...s.history, { key, turn: s.history.at(-1)?.turn ?? 1, icon: '¢', text: 'new tile', card: 'bomb', headline: true, badge: 'buy' }] });
    await new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r)));
    const [a, b] = [...document.querySelectorAll('.history-strip .strip-item')].map((el) => el.getBoundingClientRect().left);
    return { a: a ?? 0, b: b ?? 0 };
  });
  expect(xs.a).toBeLessThan(xs.b - 5);
  await shot(page, 'history-strip-slide');
  expect(await anims(page, 'strip-item')).toEqual([{ ms: 1000, easing: expoOut }]);
  expect(await anims(page, 'strip-list')).toEqual([{ ms: 1000, easing: expoOut }]);
  // Settled: in line with the rest.
  await expect.poll(async () => Math.abs((await box(tiles.nth(0))).x - (await box(tiles.nth(1))).x)).toBeLessThan(1);
});

test('speed "fast" halves the slide and "off" places a tile at once (20.21)', async ({ browser }) => {
  const fast = await startGame(browser, names(['Gus', 'Hal']), { url: `${appURL()}&anim=fast` });
  const page = fast[0];
  if (!page) {
    throw new Error('no page');
  }
  await addTiles(page, 1);
  await page.waitForTimeout(800);
  await spy(page);
  await addTiles(page, 1);
  await expect.poll(() => anims(page, 'strip-item')).toEqual([{ ms: 500, easing: expoOut }]);

  const still = await startGame(browser, names(['Ida', 'Jon']), { url: off() });
  const other = still[0];
  if (!other) {
    throw new Error('no page');
  }
  await addTiles(other, 1);
  await spy(other);
  await addTiles(other, 2);
  await expect(strip(other).locator('.strip-item').first()).toHaveAttribute('aria-label', /tile/);
  await other.waitForTimeout(300);
  expect(await anims(other, 'strip-item')).toEqual([]);
  expect(await anims(other, 'strip-list')).toEqual([]);
});

test('nothing moves while the pointer is on the strip; new tiles come in when it leaves (20.21)', async ({ browser }) => {
  const pages = await startGame(browser, names(['Kai', 'Lea']), { url: off() });
  const page = pages[0];
  if (!page) {
    throw new Error('no page');
  }
  await addTiles(page, 3);
  const tiles = strip(page).locator('.strip-item');
  await expect(tiles.first()).toHaveAttribute('aria-label', /tile/);
  const count = await tiles.count();
  const top = await tiles.first().getAttribute('data-key');
  await tiles.first().hover();
  await expect(page.getByTestId('card-zoom')).toBeVisible();
  await addTiles(page, 2);
  await page.waitForTimeout(300);
  await expect(tiles).toHaveCount(count);
  await expect(tiles.first()).toHaveAttribute('data-key', top ?? '');
  await expect(tiles.first()).toHaveClass(/hovered/);
  await page.mouse.move(900, 700);
  await expect(tiles).toHaveCount(count + 2);
  await expect(tiles.first()).not.toHaveAttribute('data-key', top ?? '');
});

test('the strip is one hit area: the nearest tile shows, sounds once, and the board dims (20.21)', async ({ browser }) => {
  const pages = await startGame(browser, names(['Mia', 'Ned']), { url: appURL() });
  const page = pages[0];
  if (!page) {
    throw new Error('no page');
  }
  await addTiles(page, 5);
  await page.mouse.move(900, 700);
  await page.waitForTimeout(1200);
  await spy(page);
  const tiles = strip(page).locator('.strip-item');
  const zoom = page.getByTestId('card-zoom');
  const dim = page.getByTestId('strip-dim');
  const [a, b] = [await box(tiles.nth(0)), await box(tiles.nth(1))];
  const x = a.x + a.width / 2;
  // In the gap between two tiles: the nearer one shows.
  await page.mouse.move(x, a.y + a.height + 1);
  await expect(tiles.nth(0)).toHaveClass(/hovered/);
  await expect(zoom).toBeVisible();
  await expect(dim).toHaveClass(/on/);
  await expect.poll(() => played(page, 'history-hover')).toBe(1);
  // A move inside the same tile plays nothing more.
  await page.mouse.move(x, a.y + a.height / 2);
  // Sliding down the strip over the gaps: the card and the dim never go.
  const end = await box(tiles.nth(3));
  for (let y = a.y + a.height / 2; y <= end.y + end.height / 2; y += 2) {
    await page.mouse.move(x, y);
    const shown = await zoom.isVisible();
    const dimmed = await dim.evaluate((el) => el.classList.contains('on'));
    expect(shown && dimmed, `card and dim at y=${String(y)}`).toBe(true);
  }
  await expect(tiles.nth(3)).toHaveClass(/hovered/);
  expect(await played(page, 'history-hover')).toBe(4);
  expect(b.y).toBeGreaterThan(a.y);
  // The dim comes in over 0.5 s and is full by then.
  await expect.poll(() => dim.evaluate((el) => Number(getComputedStyle(el).opacity))).toBe(1);
  await shot(page, 'history-strip-dim');
  // Off the strip: the card goes and the dim clears.
  await page.mouse.move(900, 700);
  await expect(zoom).toBeHidden();
  await expect(dim).not.toHaveClass(/on/);
  await expect.poll(() => dim.evaluate((el) => Number(getComputedStyle(el).opacity))).toBe(0);
});

test('a popout lays its targets in one row up to four, then in two rows (20.21)', async ({ browser }) => {
  const pages = await startGame(browser, names(['Oli', 'Pam']), { url: off() });
  const page = pages[0];
  if (!page) {
    throw new Error('no page');
  }
  const zoom = page.getByTestId('card-zoom');
  const results = zoom.getByRole('list', { name: tr('zoom.outcome') }).getByRole('listitem');
  const tiles = strip(page).locator('.strip-item');
  const four = ['two_cents', 'bomb', 'a_penny', 'dice_shard'];
  await addTiles(page, 1, four);
  await tiles.first().hover();
  await expect(results).toHaveCount(4);
  const row = await Promise.all([0, 1, 2, 3].map(async (i) => (await box(results.nth(i))).y));
  expect(new Set(row.map(Math.round)).size).toBe(1);
  await page.mouse.move(900, 700);
  await addTiles(page, 1, [...four, 'pills', 'a_nickel']);
  await tiles.first().hover();
  await expect(results).toHaveCount(6);
  const ys = await Promise.all([0, 2, 3, 5].map(async (i) => (await box(results.nth(i))).y));
  expect(Math.round(ys[0] ?? 0)).toBe(Math.round(ys[1] ?? 1));
  expect(ys[2] ?? 0).toBeGreaterThan((ys[1] ?? 0) + 10);
  expect(Math.round(ys[2] ?? 0)).toBe(Math.round(ys[3] ?? 1));
  // Never above hand size.
  expect((await box(results.nth(0).locator('.card'))).width).toBeLessThanOrEqual(100.5);
  await shot(page, 'history-popout-two-rows');
});

test("the strip's foot fades when more tiles are below (20.21)", async ({ browser }) => {
  const pages = await startGame(browser, names(['Quin', 'Rae']), { url: off() });
  const page = pages[0];
  if (!page) {
    throw new Error('no page');
  }
  const area = page.getByTestId('strip-tiles');
  await addTiles(page, 1);
  await expect(strip(page).locator('.strip-item').first()).toHaveAttribute('aria-label', /tile/);
  const few = await area.evaluate((el) => el.scrollHeight > el.clientHeight);
  if (!few) {
    await expect(area).not.toHaveClass(/strip-more/);
  }
  await addTiles(page, 25);
  await expect(area).toHaveClass(/strip-more/);
  await shot(page, 'history-strip-fade');
  await area.evaluate((el) => {
    el.scrollTop = el.scrollHeight;
  });
  await expect(area).not.toHaveClass(/strip-more/);
});
