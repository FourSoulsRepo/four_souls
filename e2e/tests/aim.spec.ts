// Aiming (15.25): one click on a card that needs a target starts aiming,
// the next click picks the target. A wrong target says why; the source
// does nothing; empty space or a right click cancels (15.33). A card
// with no target is dragged onto the board; one that asks its mode
// first aims after it, and a cancelled aim puts it back (15.34).
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
    game: { step: number; events: { kind: string; object?: number; card?: string }[] | null; view: { stack: StackItem[] | null } } | null;
    lobby: { table: { game: string } | null };
    sent: number | null;
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

// giveCard rewrites a saved game: the active player's first hand card
// becomes the card. The test server keeps saves under .cache/e2e.
function giveCard(game: string, card: string) {
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
  obj.card = card;
  writeFileSync(file, JSON.stringify(save));
}

// saveAndContinue: the host saves, the active player gets the card,
// then everyone comes back to it.
async function saveAndContinue(pages: Page[], names: string, card = 'bomb') {
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
  giveCard(game, card);
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

  // What was sent and at which step, to see that nothing goes out.
  const sent = () =>
    active.evaluate(() => {
      const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
      return { sent: s.sent, step: s.game?.step };
    });
  const before = await sent();

  // A right click cancels aiming; nothing is sent.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await active.mouse.click(640, 360, { button: 'right' });
  await expect(lifted).toHaveCount(0);
  await expect(monster).not.toHaveClass(/hl-aim/);
  await expect(active.getByRole('menu')).toHaveCount(0);
  expect(await sent()).toEqual(before);

  // A shop item is no target for a bomb: the reason shows beside the
  // click, and aiming goes on. Clicking the bomb again does nothing;
  // empty space cancels.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  const item = active.getByRole('group', { name: 'Shop' }).locator('.card').first();
  const at = await item.boundingBox();
  if (!at) {
    throw new Error('no shop item box');
  }
  await active.mouse.click(at.x + at.width / 2, at.y + at.height / 2);
  const why = active.locator('.aim-why');
  await expect(why).toHaveText('Pick a monster or a player.');
  await shot(active, 'aim-wrong');
  const note = await why.locator('div').boundingBox();
  if (!note) {
    throw new Error('no reason box');
  }
  // The note sits just above the click, centered on it.
  expect(Math.abs(note.x + note.width / 2 - (at.x + at.width / 2)) < 3 && note.y + note.height < at.y + at.height / 2).toBe(true);
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await active.getByText('Nothing on the stack').click();
  await expect(lifted).toHaveCount(0);
  await expect(why).toHaveCount(0);
  await expect(active.getByRole('menu')).toHaveCount(0);
  expect(await sent()).toEqual(before);

  // One click starts aiming: the card stays lifted, the targets glow, and
  // an arrow runs from the card's middle to the pointer.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  // A player glows by the character card, not by the area (15.33).
  const opponent = active.locator('.player.opponent');
  await expect(opponent.locator('.opp-character .card')).toHaveClass(/hl-aim/);
  await expect(opponent).not.toHaveClass(/hl-aim/);
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

// stackCards lists the cards on the page's stack.
function stackCards(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    return (store.getState().game?.view.stack ?? []).map((it) => it.card ?? '');
  });
}

// dragToBoard drags a hand card onto the board, by its top left: the
// next card in the fan covers its middle, and a tilted card misses
// its box's very corner. The card is measured once it stands still: a
// card just put back slides down into the fan. The first move stays
// on the card, so the card sees it and takes the pointer.
async function dragToBoard(page: Page, card: string) {
  const loc = page.getByRole('group', { name: 'Your hand' }).getByRole('button', { name: card }).first();
  let from = await loc.boundingBox();
  for (let still = false; !still; ) {
    await page.waitForTimeout(100);
    const now = await loc.boundingBox();
    still = !!from && !!now && from.x === now.x && from.y === now.y;
    from = now;
  }
  const to = await page.getByRole('region', { name: 'Board' }).boundingBox();
  if (!from || !to) {
    throw new Error('no boxes');
  }
  const x = from.x + 40;
  const y = from.y + 40;
  await page.mouse.move(x, y);
  await page.mouse.down();
  await page.mouse.move(x, y - 12, { steps: 2 });
  await page.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
  await page.mouse.up();
}

test('a loot card with no target dragged onto the board plays', async ({ browser }) => {
  const pages = await startGame(browser, ['Dot', 'Eve'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Dot, Eve', 'a_penny');
  const active = await activePage(pages, click);
  const hand = active.getByRole('group', { name: 'Your hand' });
  // Any role: once played, other pennies are no longer buttons.
  const penny = hand.getByLabel('A Penny!', { exact: true });
  const before = await penny.count();

  // Played loot is kept as it shows: the stack may resolve at once.
  await active.evaluate(() => {
    const w = window as unknown as { fourSouls: { store: Store }; played: string[] };
    w.played = [];
    let step = -1;
    w.fourSouls.store.subscribe((s) => {
      if (s.game && s.game.step !== step) {
        step = s.game.step;
        w.played.push(...(s.game.events ?? []).filter((e) => e.kind === 'loot_played').map((e) => e.card ?? ''));
      }
    });
  });

  await dragToBoard(active, 'A Penny!');
  await expect.poll(() => active.evaluate(() => (window as unknown as { played: string[] }).played)).toContain('a_penny');
  await expect(penny).toHaveCount(before - 1);
  // No arrow, no dialog: the card went straight to the stack.
  await expect(active.locator('.aim-arrow')).toHaveCount(0);
  await expect(active.getByRole('dialog')).toHaveCount(0);
});

test('a card placed first aims after its mode; a cancelled aim leaves it in hand', async ({ browser }) => {
  const pages = await startGame(browser, ['Fay', 'Gus'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Fay, Gus', 'dagaz');
  const active = await activePage(pages, click);
  const hand = active.getByRole('group', { name: 'Your hand' });
  const dagaz = hand.getByRole('button', { name: 'Dagaz' });
  const lifted = active.locator('.hand-card.aiming');
  const opponent = active.locator('.player.opponent .opp-character .card');
  const mode = 'Choose a player. Prevent the next 1 damage they would take this turn.';
  // place drags Dagaz onto the board and picks the mode that aims.
  const place = async () => {
    await dragToBoard(active, 'Dagaz');
    const ask = active.getByRole('dialog', { name: 'Choose one' });
    await expect(ask).toBeVisible();
    await shot(active, 'aim-mode');
    await ask.getByRole('button', { name: mode }).click();
    // The target prompt is an aim from the card, not a dialog.
    await expect(lifted).toHaveCount(1);
    await expect(opponent).toHaveClass(/hl-aim/);
    await expect(active.getByRole('dialog')).toHaveCount(0);
    // Players are aimed at by the character card; areas stay plain.
    await expect(active.locator('.player.hl-pick')).toHaveCount(0);
    await expect(active.getByRole('button', { name: /^Choose / })).toHaveCount(0);
  };

  await place();
  const to = await opponent.boundingBox();
  if (!to) {
    throw new Error('no character box');
  }
  await active.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
  await expect(active.locator('.aim-arrow')).toHaveClass(/\bon\b/);
  await shot(active, 'aim-placed');

  // A right click cancels the play: the card stays in hand, the stack
  // stays empty, and nothing aims any more.
  await active.mouse.click(640, 360, { button: 'right' });
  await expect(lifted).toHaveCount(0);
  await expect(opponent).not.toHaveClass(/hl-aim/);
  await expect(dagaz).toHaveCount(1);
  await expect(dagaz).toHaveClass(/hl-act/);
  expect(await stackCards(active)).not.toContain('dagaz');

  // Placed again, a click on the opponent's character plays it there.
  await place();
  await active.mouse.click(to.x + to.width / 2, to.y + to.height / 2);
  await expect.poll(() => stackCards(active)).toContain('dagaz');
  await expect(dagaz).toHaveCount(0);
});
