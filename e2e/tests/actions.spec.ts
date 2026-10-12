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
  passButton,
  settle,
  tap,
  turnAction,
  turnNumber,
  type Tap,
} from './turn';

const fast = () => `${appURL()}&anim=off`;

// A whole turn can take many moves (combat rolls, responses).
test.setTimeout(120_000);

// backToMe: the active player may end the turn again, or the turn or
// game moved on (a death ends it early).
function backToMe(active: Page, turn: number) {
  return async () => (await canEndTurn(active)) || (await turnNumber(active)) !== turn || (await gameOver(active));
}

// attack declares an attack; settling picks the target on the table.
async function attack(pages: Page[], active: Page, press: Tap, turn: number) {
  const button = turnAction(active, 'Attack');
  if (!(await canEndTurn(active)) || !(await button.isEnabled())) {
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
  const pages = await startGame(browser, ['Max', 'Nia'], { url: fast() });
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
  const pages = await startGame(browser, ['Oz', 'Pia'], { url: fast(), touch: true });
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
  const pages = await startGame(browser, ['Quin', 'Ray'], { url: fast() });
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

// allowedKinds reads the intent kinds the engine allows this page now,
// or none while the page may not act.
async function allowedKinds(page: Page): Promise<number[]> {
  return page.evaluate(() => {
    interface S {
      game: { step: number; allowed?: { kind: number }[] | null } | null;
      behind: boolean;
      sent: number | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const ready = s.game !== null && !s.behind && s.sent !== s.game.step;
    return ready ? (s.game?.allowed ?? []).map((a) => a.kind) : [];
  });
}

test('buy and attack icons are lit only when the engine allows them; B buys', async ({ browser }) => {
  const pages = await startGame(browser, ['Sol', 'Tia'], { url: fast() });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);
  for (const page of pages) {
    const kinds = await allowedKinds(page);
    const buy = turnAction(page, 'Buy');
    const attack = turnAction(page, 'Attack');
    await expect(buy).toBeVisible();
    await expect(attack).toBeVisible();
    await (kinds.includes(6) ? expect(buy).toBeEnabled() : expect(buy).toBeDisabled());
    await (kinds.includes(5) ? expect(attack).toBeEnabled() : expect(attack).toBeDisabled());
  }
  await expect(turnAction(active, 'Buy')).toBeEnabled();
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(turnAction(other, 'Buy')).toBeDisabled();
  await expect(turnAction(other, 'Attack')).toBeDisabled();
  await shot(active, 'turn-actions');

  // B declares the purchase; once the others pass, the Buy icon stays
  // pressed while the player picks what to buy (F-13).
  await active.keyboard.press('b');
  await settle(pages, click, async () => active.locator('.choice').isVisible());
  const choosing = active
    .getByRole('toolbar', { name: 'Turn actions' })
    .getByRole('button', { name: 'Choosing what to buy' });
  await expect(choosing).toHaveAttribute('aria-pressed', 'true');
  await shot(active, 'turn-actions-buy');
  await settle(pages, click, backToMe(active, turn));
});

// recordActs keeps every intent the page sends in window.acts.
async function recordActs(page: Page): Promise<void> {
  await page.evaluate(() => {
    interface I {
      kind: number;
    }
    interface S {
      act: (i: I) => void;
    }
    const w = window as unknown as {
      acts: I[];
      fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } };
    };
    const { store } = w.fourSouls;
    const act = store.getState().act;
    w.acts = [];
    store.setState({
      act: (i) => {
        w.acts.push(i);
        act(i);
      },
    });
  });
}

async function lastAct(page: Page): Promise<number | undefined> {
  return page.evaluate(() => (window as unknown as { acts: { kind: number }[] }).acts.at(-1)?.kind);
}

test('the big End turn button: E ends the turn, Space passes', async ({ browser }) => {
  const pages = await startGame(browser, ['Yan', 'Zoe'], { url: fast() });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(endTurnButton(active)).toBeEnabled();
  await expect(endTurnButton(other, 'Their turn')).toBeDisabled();
  // The game's fonts are registered from their URLs; without the files
  // they fail to load and the system font shows.
  const families = await active.evaluate(() => [...document.fonts].map((f) => f.family.replace(/"/g, '')));
  expect(families).toEqual(expect.arrayContaining(['Four Souls Title', 'Four Souls Cyrillic']));
  await shot(active, 'end-turn-own');
  // Off the hand, so no card zoom covers the board.
  await other.mouse.move(1200, 420);
  await shot(other, 'end-turn-theirs');

  // Space passes priority, like the Pass button (F-15).
  for (const page of pages) {
    await recordActs(page);
  }
  await active.keyboard.press('b');
  const passing = async () => {
    for (const page of pages) {
      const pass = passButton(page);
      if ((await pass.count()) > 0 && (await pass.isEnabled())) {
        return page;
      }
    }
    return undefined;
  };
  await expect.poll(passing).toBeDefined();
  const page = await passing();
  if (!page) {
    throw new Error('nobody may pass');
  }
  await page.keyboard.press(' ');
  expect(await lastAct(page)).toBe(0);
  await settle(pages, click, backToMe(active, turn));

  // E ends the turn, like the button.
  if (await canEndTurn(active)) {
    await active.keyboard.press('e');
    expect(await lastAct(active)).toBe(1);
    await expectNextTurn(pages, click, turn);
  }
});
