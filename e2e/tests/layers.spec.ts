// Layers (19.26, ADR 015): a window that must be on top is on top. The
// order is docs/wiki/Layers.md; these tests check the pairs a player
// meets.
import type { Locator } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

async function box(l: Locator) {
  const b = await l.boundingBox();
  if (!b) {
    throw new Error('no box');
  }
  return b;
}

// onTop is true when the element drawn at the middle of where a and b
// overlap belongs to a. The zoom and the aiming arrows let the pointer
// through, so the probe turns hit testing back on for them first, and
// off again in the same task: a hover check in between would close the
// zoom.
async function onTop(a: Locator, b: Locator): Promise<boolean> {
  const [ba, bb] = [await box(a), await box(b)];
  const left = Math.max(ba.x, bb.x);
  const right = Math.min(ba.x + ba.width, bb.x + bb.width);
  const top = Math.max(ba.y, bb.y);
  const bottom = Math.min(ba.y + ba.height, bb.y + bb.height);
  expect(right - left, 'the two overlap').toBeGreaterThan(0);
  expect(bottom - top, 'the two overlap').toBeGreaterThan(0);
  return a.evaluate(
    (el, p) => {
      const probe = document.createElement('style');
      probe.textContent = '.card-zoom, .card-zoom *, .aim-arrow * { pointer-events: auto !important; }';
      document.head.append(probe);
      const at = document.elementFromPoint(p.x, p.y);
      probe.remove();
      return at !== null && el.contains(at);
    },
    { x: (left + right) / 2, y: (top + bottom) / 2 },
  );
}

test('a card zoomed from the discard window shows above it (19.26)', async ({ browser }) => {
  const [lyra] = await startGame(browser, ['Lyra', 'Orrin'], { url: fast() });
  if (!lyra) {
    throw new Error('no pages');
  }
  // The discard prompt as the server sends it: one card from the hand.
  await lyra.evaluate(() => {
    interface S {
      game: { view: { players: { hand?: { id: number }[] }[]; waiting: unknown }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const hand = (game.view.players[you]?.hand ?? []).map((c) => c.id);
    game.view.waiting = { kind: 2, player: you, count: 1 };
    game.allowed = [{ player: you, kind: 2, objects: hand }];
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  });
  const dialog = lyra.getByRole('dialog', { name: 'Discard', exact: true });
  const card = dialog.locator('.card').first();
  await expect(card).toBeVisible();
  await card.hover();
  const zoom = lyra.getByTestId('card-zoom').locator('.card-zoom-body');
  await expect(zoom).toBeVisible();
  await shot(lyra, 'discard-zoom');
  expect(await onTop(zoom, dialog)).toBe(true);
});

test('the Menu button stays above a choice and opens over it (15.29)', async ({ browser }) => {
  const [tamsin] = await startGame(browser, ['Tamsin', 'Ivor'], { url: fast() });
  if (!tamsin) {
    throw new Error('no pages');
  }
  // A choice as the server sends it: its backdrop dims the whole table.
  await tamsin.evaluate(() => {
    interface S {
      game: { view: { waiting: unknown }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    game.view.waiting = { kind: 4, player: you, purpose: 9, text: 'Keep going?', options: ['yes', 'No'] };
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  });
  const backdrop = tamsin.locator('.choice-backdrop');
  await expect(backdrop).toBeVisible();
  const menu = tamsin.getByRole('button', { name: 'Menu' });
  expect(await onTop(menu, backdrop)).toBe(true);
  await menu.click();
  await expect(tamsin.getByRole('dialog', { name: 'Game menu' })).toBeVisible();
});

test("another player's arrow stays under a choice (15.36)", async ({ browser }) => {
  const [wren] = await startGame(browser, ['Wren', 'Abel'], { url: fast() });
  if (!wren) {
    throw new Error('no pages');
  }
  // A choice for this player, while the other one aims from their
  // character at the board's middle, as the server relays it.
  await wren.evaluate(() => {
    interface S {
      game: { view: { waiting: unknown; players: { character: { id: number } }[] }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const other = 1 - you;
    game.view.waiting = { kind: 4, player: you, purpose: 9, text: 'Keep going?', options: ['yes', 'No'] };
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
    const from = game.view.players[other]?.character.id ?? 0;
    s.receive({ type: 'presence', data: { seat: other, from: { object: from }, point: { x: 0.5, y: 0.5 } } });
  });
  const backdrop = wren.locator('.choice-backdrop');
  const head = wren.locator('.aim-arrow.others .aim-head');
  await expect(backdrop).toBeVisible();
  await expect(head).toBeAttached();
  expect(await onTop(backdrop, head)).toBe(true);
});
