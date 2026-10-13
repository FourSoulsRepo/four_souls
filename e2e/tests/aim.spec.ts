// Drag and aim (15.25): a card that needs a target is dragged onto it.
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

test('a bomb dragged onto a monster damages it', async ({ browser }) => {
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
  await expect(bomb).toBeVisible();
  const from = await bomb.boundingBox();
  if (!from) {
    throw new Error('no bomb box');
  }
  await active.mouse.move(from.x + from.width / 2, from.y + from.height / 2);
  await active.mouse.down();
  await active.mouse.move(from.x + from.width / 2, from.y - 60, { steps: 4 });

  // Aiming: the card stays lifted, an arrow points, the targets glow.
  const monster = active.getByRole('group', { name: 'Monsters' }).locator('.card.hl-aim').first();
  await expect(monster).toBeVisible();
  await expect(active.locator('.aim-arrow')).toBeVisible();
  await expect(active.locator('.hand-card.aiming')).toHaveCount(1);
  const zone = await monster.getAttribute('data-zone');
  const id = Number(zone?.replace('obj-', ''));
  const to = await monster.boundingBox();
  if (!to) {
    throw new Error('no monster box');
  }
  await active.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
  await shot(active, 'aim');
  await active.mouse.up();
  await expect(active.locator('.aim-arrow')).toHaveCount(0);

  // The drop picked the monster: no dialog asked, the bomb aims at it.
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
