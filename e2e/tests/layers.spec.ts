// Layers (19.26, ADR 015): a window that must be on top is on top. The
// order is docs/wiki/Layers.md; these tests check the pairs a player
// meets.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL, serverPort } from './server';
import { shot } from './shot';
import { tr } from './locale';

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
      game: { step: number; view: { players: { hand?: { id: number }[] }[]; waiting: unknown }; allowed: unknown[] };
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
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
  });
  const dialog = lyra.getByRole('dialog', { name: tr('discard.label'), exact: true });
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
      game: { step: number; view: { waiting: unknown }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    game.view.waiting = { kind: 4, player: you, purpose: 9, text: { key: 'card.sack_of_pennies.recharge_sack_of_pennies' }, options: [{ key: 'ask.yes' }, { key: 'ask.no' }] };
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
  });
  const backdrop = tamsin.locator('.choice-backdrop');
  await expect(backdrop).toBeVisible();
  const menu = tamsin.getByRole('button', { name: tr('game_menu.button') });
  expect(await onTop(menu, backdrop)).toBe(true);
  await menu.click();
  await expect(tamsin.getByRole('dialog', { name: tr('game_menu.label') })).toBeVisible();
});

test('the rollback vote covers the table; the game menu opens over it (17.4)', async ({ browser }) => {
  const [uma] = await startGame(browser, ['Uma', 'Vik'], { url: fast() });
  if (!uma) {
    throw new Error('no pages');
  }
  // A Soft vote as the server sends it: Vik asks to undo a step.
  await uma.evaluate(() => {
    interface S {
      game: { step: number };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const you = s.lobby.table.you;
    const vote = { step: s.game.step, seat: 1 - you, author: 1 - you, undo: 1, voters: [you], votes: [], need: 1 };
    s.receive({ type: 'update', data: { ...structuredClone(s.game), events: [], rollback_mode: 'soft', rollback_vote: vote } });
  });
  const panel = uma.getByRole('dialog', { name: tr('rollback.label') });
  await expect(panel).toBeVisible();
  expect(await onTop(panel, uma.getByRole('main'))).toBe(true);
  await uma.getByRole('button', { name: tr('game_menu.button') }).click();
  const menu = uma.getByRole('dialog', { name: tr('game_menu.label') });
  await expect(menu).toBeVisible();
  expect(await onTop(menu, panel)).toBe(true);
});

// noPresence drops the server's presence frames on the way to the page:
// the other seat's real pointer would replace the injected one (19.48).
function noPresence(page: Page): Promise<void> {
  return page.routeWebSocket(`ws://127.0.0.1:${String(serverPort())}/ws`, (ws) => {
    const server = ws.connectToServer();
    server.onMessage((msg) => {
      if ((JSON.parse(String(msg)) as { type: string }).type !== 'presence') {
        ws.send(msg);
      }
    });
  });
}

test("another player's arrow stays under a choice (15.36)", async ({ browser }) => {
  const [wren] = await startGame(browser, ['Wren', 'Abel'], { url: fast(), open: noPresence });
  if (!wren) {
    throw new Error('no pages');
  }
  // A choice for this player, while the other one aims from their
  // character at the board's middle, as the server relays it.
  await wren.evaluate(() => {
    interface S {
      game: { step: number; view: { waiting: unknown; players: { character: { id: number } }[] }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const other = 1 - you;
    game.view.waiting = { kind: 4, player: you, purpose: 9, text: { key: 'card.sack_of_pennies.recharge_sack_of_pennies' }, options: [{ key: 'ask.yes' }, { key: 'ask.no' }] };
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
    const from = game.view.players[other]?.character.id ?? 0;
    s.receive({ type: 'presence', data: { seat: other, from: { object: from }, point: { x: 0.5, y: 0.5 } } });
  });
  const backdrop = wren.locator('.choice-backdrop');
  const head = wren.locator('.aim-arrow.others .aim-head');
  await expect(backdrop).toBeVisible();
  await expect(head).toBeAttached();
  expect(await onTop(backdrop, head)).toBe(true);
});

test('the board dims under a hovered history strip; the strip and the zoom stay above (20.21)', async ({ browser }) => {
  const [page] = await startGame(browser, ['Tess', 'Ugo'], { url: fast() });
  if (!page) {
    throw new Error('no pages');
  }
  // A headline as the store keeps it: a buy of Bomb!.
  await page.evaluate(() => {
    interface S {
      history: { key: number }[];
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } }).fourSouls;
    const { history } = store.getState();
    const key = (history.at(-1)?.key ?? 0) + 1;
    store.setState({ history: [...history, { key, turn: 1, icon: '¢', text: 'a buy', card: 'bomb', headline: true, badge: 'buy', owner: 0 } as { key: number }] });
  });
  const strip = page.getByRole('complementary', { name: tr('history.strip') });
  const tile = strip.locator('.strip-item').first();
  await expect(tile).toBeVisible();
  const dim = page.getByTestId('strip-dim');
  // The decks at the board's right edge: the zoom opens far from them.
  const decks = page.locator('.board-decks .deck-pair').last();
  await tile.hover();
  await expect(dim).toHaveClass(/on/);
  // drawn is what is drawn at the middle of l. The dim and the zoom let
  // the pointer through: the probe turns hit testing on for the dim and
  // the zoom's card, and off for the zoom's frame, which covers the window.
  const drawn = async (l: Locator) => {
    const b = await box(l);
    return page.evaluate(
      (p) => {
        const probe = document.createElement('style');
        probe.textContent = '.strip-dim, .card-zoom .card { pointer-events: auto !important; } .card-zoom, .card-zoom-body { pointer-events: none !important; }';
        document.head.append(probe);
        const top = document.elementFromPoint(p.x, p.y);
        probe.remove();
        return top?.closest('.strip-dim') ? 'dim' : top?.closest('.strip-item') ? 'tile' : top?.closest('.card-zoom') ? 'zoom' : (top?.className ?? '');
      },
      { x: b.x + b.width / 2, y: b.y + b.height / 2 },
    );
  };
  expect(await drawn(decks)).toBe('dim');
  expect(await drawn(tile)).toBe('tile');
  const zoom = page.getByTestId('card-zoom').locator('.card.size-zoom');
  await expect(zoom).toBeVisible();
  expect(await drawn(zoom)).toBe('zoom');
  await shot(page, 'layers-strip-dim');
});
