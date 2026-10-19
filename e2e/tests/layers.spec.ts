// Layers (19.26, ADR 015): a window that must be on top is on top. The
// order is docs/wiki/Layers.md; these tests check the pairs a player
// meets.
import type { Locator, Page } from '@playwright/test';
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
// overlap belongs to a. The zoom lets the pointer through, so the probe
// turns hit testing back on for it first.
async function onTop(page: Page, a: Locator, b: Locator): Promise<boolean> {
  const [ba, bb] = [await box(a), await box(b)];
  const left = Math.max(ba.x, bb.x);
  const right = Math.min(ba.x + ba.width, bb.x + bb.width);
  const top = Math.max(ba.y, bb.y);
  const bottom = Math.min(ba.y + ba.height, bb.y + bb.height);
  expect(right - left, 'the two overlap').toBeGreaterThan(0);
  expect(bottom - top, 'the two overlap').toBeGreaterThan(0);
  const probe = await page.addStyleTag({ content: '.card-zoom, .card-zoom * { pointer-events: auto !important; }' });
  const hit = await a.evaluate(
    (el, p) => {
      const at = document.elementFromPoint(p.x, p.y);
      return at !== null && el.contains(at);
    },
    { x: (left + right) / 2, y: (top + bottom) / 2 },
  );
  await probe.evaluate((el) => {
    el.remove();
  });
  return hit;
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
  expect(await onTop(lyra, zoom, dialog)).toBe(true);
});
