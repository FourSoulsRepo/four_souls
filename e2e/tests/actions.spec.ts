// Player actions (7.7): a full turn with the mouse, and with touch, on
// real games.
import { expect, test, type Page } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import {
  activePage,
  canEndTurn,
  click,
  endTurnButton,
  expectNextTurn,
  gameOver,
  settle,
  tap,
  turnNumber,
  type Tap,
} from './turn';

const fast = `${appURL}&anim=off`;

// A whole turn can take many moves (combat rolls, responses).
test.setTimeout(120_000);

// backToMe: the active player may end the turn again, or the turn or
// game moved on (a death ends it early).
function backToMe(active: Page, turn: number) {
  return async () => (await canEndTurn(active)) || (await turnNumber(active)) !== turn || (await gameOver(active));
}

// attack declares an attack; settling picks the target on the table.
async function attack(pages: Page[], active: Page, press: Tap, turn: number) {
  const button = active.getByRole('region', { name: 'Actions' }).getByRole('button', { name: 'Attack' });
  if (!(await canEndTurn(active)) || !(await button.isVisible())) {
    return;
  }
  await press(active, button);
  await settle(pages, press, backToMe(active, turn));
}

// endTurn ends the turn if it is still the active player's.
async function endTurn(active: Page, press: Tap) {
  if (await canEndTurn(active)) {
    await press(active, endTurnButton(active));
  }
}

test('a full turn with the mouse: drag a loot card, attack, end the turn', async ({ browser }) => {
  const pages = await startGame(browser, ['Max', 'Nia'], { url: fast });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);

  const loot = active.locator('.hand .card.hl-act').first();
  if (await loot.isVisible()) {
    await shot(active, 'actions');
    const from = await loot.boundingBox();
    const to = await active.getByRole('region', { name: 'Board' }).boundingBox();
    if (!from || !to) {
      throw new Error('no boxes');
    }
    await active.mouse.move(from.x + from.width / 2, from.y + from.height / 2);
    await active.mouse.down();
    await active.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
    await active.mouse.up();
    await settle(pages, click, backToMe(active, turn));
  }
  await attack(pages, active, click, turn);
  await endTurn(active, click);
  await expectNextTurn(pages, click, turn);
});

test('a full turn with touch: tap a loot card to play it, attack, end the turn', async ({ browser }) => {
  const pages = await startGame(browser, ['Oz', 'Pia'], { url: fast, touch: true });
  const active = await activePage(pages, tap);
  const turn = await turnNumber(active);

  const loot = active.locator('.hand .card.hl-act').first();
  if (await loot.isVisible()) {
    await loot.tap();
    const menu = active.getByRole('menu');
    await expect(menu).toBeVisible();
    await shot(active, 'action-menu');
    await menu.getByRole('menuitem', { name: 'Play' }).tap();
    await settle(pages, tap, backToMe(active, turn));
  }
  await attack(pages, active, tap, turn);
  await endTurn(active, tap);
  await expectNextTurn(pages, tap, turn);
});

test('a refused action shows the reason', async ({ browser }) => {
  const pages = await startGame(browser, ['Quin', 'Ray'], { url: fast });
  const active = await activePage(pages, click);
  // Send what the server must refuse: the other player cannot end this turn.
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await other.evaluate(() => {
    interface S {
      act: (i: { player: number; kind: number }) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store;
    store.getState().act({ player: 0, kind: 1 });
  });
  await expect(other.getByRole('region', { name: 'Actions' }).getByRole('status')).toContainText(/\(R-[A-Z]+-\d+\)/);
});
