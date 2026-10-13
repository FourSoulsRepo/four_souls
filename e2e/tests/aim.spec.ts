// Aiming (15.25): one click on a card that needs a target starts aiming,
// the next click picks the target.
import { expect, test, type Page } from '@playwright/test';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { startGame } from './game';
import { appURL, serverPort } from './server';
import { shot } from './shot';
import { activePage, click, settle } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

interface Store {
  getState: () => {
    game: { step: number; events: { kind: string; object?: number }[] | null; view: { stack: StackItem[] | null } } | null;
    lobby: { table: { game: string } | null };
  };
  subscribe: (f: (s: ReturnType<Store['getState']>) => void) => void;
}

interface StackItem {
  card?: string;
  targets?: { object?: number }[];
}

// gameID is the id of the page's table, which names its save file.
async function gameID(page: Page): Promise<string> {
  return page.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    return store.getState().lobby.table?.game ?? '';
  });
}

// giveBomb rewrites a saved game: the active player's first hand card
// becomes a Bomb!. The test server keeps saves under .cache/e2e.
function giveBomb(game: string) {
  const file = join(import.meta.dirname, '..', '..', '.cache', 'e2e', String(serverPort()), 'saves', `${game}.fssave`);
  const save = JSON.parse(readFileSync(file, 'utf8')) as {
    state: { turn: { active: number }; players: { hand: number[] }[]; objects: { card: string }[] };
  };
  const { state } = save;
  const id = state.players[state.turn.active]?.hand[0];
  const obj = id === undefined ? undefined : state.objects[id];
  if (!obj) {
    throw new Error('the active player has no hand card');
  }
  obj.card = 'bomb';
  writeFileSync(file, JSON.stringify(save));
}

// saveAndContinue: the host saves, then everyone comes back to it.
async function saveAndContinue(pages: Page[], names: string) {
  const [host, ...guests] = pages;
  if (!host) {
    throw new Error('no pages');
  }
  const game = await gameID(host);
  await host.getByRole('button', { name: 'Menu' }).click();
  await host.getByRole('button', { name: 'Save and leave' }).click();
  for (const page of pages) {
    await expect(page.getByText('The game was saved.')).toBeVisible();
  }
  giveBomb(game);
  await host.getByRole('button', { name: 'Back' }).click();
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByRole('button', { name: 'Continue a saved game' }).click();
  await host.getByRole('listitem').filter({ hasText: names }).getByRole('button', { name: 'Continue' }).click();
  for (const page of guests) {
    await page.getByRole('listitem').filter({ hasText: names }).getByRole('button', { name: 'Continue' }).click();
  }
  for (const page of pages) {
    await page.getByRole('button', { name: 'Ready', exact: true }).click();
  }
  for (const page of pages) {
    await expect(page.getByRole('main', { name: /^Turn / })).toBeVisible();
  }
}

test('a bomb aimed at a monster by two clicks damages it', async ({ browser }) => {
  const pages = await startGame(browser, ['Bom', 'Cid'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Bom, Cid');
  const active = await activePage(pages, click);
  // Damage events are kept as they show, update by update.
  await active.evaluate(() => {
    const w = window as unknown as { fourSouls: { store: Store }; damaged: number[] };
    w.damaged = [];
    let step = -1;
    w.fourSouls.store.subscribe((s) => {
      if (s.game && s.game.step !== step) {
        step = s.game.step;
        w.damaged.push(...(s.game.events ?? []).filter((e) => e.kind === 'damaged').map((e) => e.object ?? -1));
      }
    });
  });

  const bomb = active.getByRole('group', { name: 'Your hand' }).getByRole('button', { name: 'Bomb!' }).first();
  const monster = active.getByRole('group', { name: 'Monsters' }).locator('.card').first();
  const lifted = active.locator('.hand-card.aiming');
  // Click the card's visible top: a locator click would scroll the hand,
  // which hangs past the window's edge, into view.
  const box = await bomb.boundingBox();
  if (!box) {
    throw new Error('no bomb box');
  }
  const clickBomb = () => active.mouse.click(box.x + box.width / 2, box.y + 30);

  // A right click cancels aiming; nothing is sent.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await active.mouse.click(640, 360, { button: 'right' });
  await expect(lifted).toHaveCount(0);
  await expect(monster).not.toHaveClass(/hl-aim/);
  await expect(active.getByRole('menu')).toHaveCount(0);

  // One click starts aiming: the card stays lifted, the targets glow, and
  // an arrow runs from the card's middle to the pointer.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  const id = Number((await monster.getAttribute('data-zone'))?.replace('obj-', ''));
  const to = await monster.boundingBox();
  if (!to) {
    throw new Error('no monster box');
  }
  const tip = { x: to.x + to.width / 2, y: to.y + to.height / 2 };
  // Off a target the head is plain; on one, a ring marks it.
  const arrow = active.locator('.aim-arrow');
  await active.mouse.move(to.x - 200, tip.y, { steps: 12 });
  await expect(arrow).toBeVisible();
  await expect(arrow).not.toHaveClass(/\bon\b/);
  await expect(active.locator('.aim-block').first()).toBeAttached();
  await shot(active, 'aim-off');
  await active.mouse.move(tip.x, tip.y, { steps: 12 });
  await expect(arrow).toHaveClass(/\bon\b/);
  await shot(active, 'aim');
  const geometry = await active.evaluate(() => {
    const box = (s: string) => document.querySelector(s)?.getBoundingClientRect();
    return { head: box('.aim-head'), ring: box('.aim-ring > div'), chain: box('.aim-chain'), card: box('.hand-card.aiming .card') };
  });
  const { head, ring, chain, card } = geometry;
  if (!head || !ring || !chain || !card) {
    throw new Error('no arrow');
  }
  // The head ends at the pointer, where the ring is centered; the chain
  // starts in the card's middle.
  const near = (a: number, b: number) => Math.abs(a - b) < 3;
  expect(tip.x >= head.left - 1 && tip.x <= head.right + 1 && tip.y >= head.top - 1 && tip.y <= head.bottom + 1).toBe(true);
  expect(near(ring.left + ring.width / 2, tip.x) && near(ring.top + ring.height / 2, tip.y)).toBe(true);
  const middle = { x: card.left + card.width / 2, y: card.top + card.height / 2 };
  expect(near(chain.left, middle.x) && near(chain.top, middle.y)).toBe(true);

  // The second click picks the target.
  await active.mouse.click(tip.x, tip.y);
  await expect(active.locator('.aim-arrow')).toHaveCount(0);

  // The click picked the monster: no dialog asked, the bomb aims at it.
  await expect
    .poll(() =>
      active.evaluate((obj) => {
        const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
        return (store.getState().game?.view.stack ?? []).some((it) => it.card === 'bomb' && it.targets?.[0]?.object === obj);
      }, id),
    )
    .toBe(true);
  const damaged = () => active.evaluate((obj) => (window as unknown as { damaged: number[] }).damaged.includes(obj), id);
  await settle(pages, click, damaged);
  expect(await damaged()).toBe(true);
});
