// Stack and responses (7.8): players always see why the game waits.
import { expect, test, type Page } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, endTurnButton, passButton } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

test('a played loot card shows on the stack; others may skip all', async ({ browser }) => {
  const pages = await startGame(browser, ['Sam', 'Tia'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  const name = await loot.getAttribute('aria-label');
  await loot.click();
  await active.getByRole('menu').getByRole('menuitem', { name: 'Play' }).click();
  // Some loot cards ask for a target before they go on the stack.
  const stack = active.getByRole('complementary', { name: 'Stack' });
  for (let i = 0; i < 10 && (await stack.locator('.stack-row').count()) === 0; i++) {
    const option = active.locator('.choice .choice-option button').first();
    if (await option.isVisible()) {
      await option.click();
    }
    await active.waitForTimeout(200);
  }

  for (const page of [active, other]) {
    await expect(page.getByRole('complementary', { name: 'Stack' })).toContainText(`plays ${name ?? ''}`);
  }
  // Both always see why the game waits.
  await expect(active.locator('.action-status')).toHaveText(/Respond, or pass|is choosing|Choose/);
  await expect(other.locator('.action-status')).toHaveText(/(Sam|Tia) (may respond|is choosing)/);
  await shot(other, 'stack');

  const bar = other.getByRole('region', { name: 'Actions' });
  await bar.getByRole('button', { name: 'Skip all' }).click();
  await expect(bar.getByRole('button', { name: 'Cancel skip' })).toBeVisible();
  await bar.getByRole('button', { name: 'Cancel skip' }).click();
  await expect(bar.getByRole('button', { name: 'Skip all' })).toBeVisible();
});

test('a player may wait for their turn; the end of the turn before asks them (N-16)', async ({ browser }) => {
  const pages = await startGame(browser, ['Wan', 'Oti'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await shot(other, 'wait-turn-offer');
  const bar = other.getByRole('region', { name: 'Actions' });
  await bar.getByRole('button', { name: 'Wait for my turn' }).click();
  await expect(bar.getByRole('button', { name: 'Cancel wait' })).toBeVisible();
  await shot(other, 'wait-turn-on');
  await bar.getByRole('button', { name: 'Cancel wait' }).click();
  await expect(bar.getByRole('button', { name: 'Wait for my turn' })).toBeVisible();
  await bar.getByRole('button', { name: 'Wait for my turn' }).click();
  await expect(bar.getByRole('button', { name: 'Cancel wait' })).toBeVisible();

  // The turn before ends: waiting is over, and nothing more is offered.
  await endTurnButton(active).click();
  await expect(bar.getByRole('button', { name: /Cancel wait|Wait for my turn/ })).toHaveCount(0);
  if (await passButton(other).isEnabled()) {
    await expect(other.locator('.action-status')).toHaveText(/You may act|Respond, or pass/);
    await shot(other, 'wait-turn-last-chance');
  }
});

test('the response timer counts down for the waiting player', async ({ browser }) => {
  const pages = await startGame(browser, ['Uma', 'Vic'], {
    url: fast(),
    setup: async (host) => {
      await host.getByLabel('Detailed setup').check();
      await host.getByLabel('Response timer').selectOption('60');
    },
  });
  const active = await activePage(pages, click);
  const timer = active.getByLabel('Time left');
  await expect(timer).toHaveText(/· \d+ s/);
  const first = Number((await timer.textContent())?.match(/\d+/)?.[0]);
  expect(first).toBeGreaterThan(0);
  await expect.poll(async () => Number((await timer.textContent())?.match(/\d+/)?.[0])).toBeLessThan(first);
});

// waitLine is a stack panel line with "You" turned into the page's own
// name, so every page's line can be compared.
async function waitLine(page: Page, me: string, label: string): Promise<string> {
  const line = page.getByRole('complementary', { name: 'Stack' }).getByLabel(label, { exact: true });
  if ((await line.count()) === 0) {
    return '';
  }
  return (await line.first().innerText()).replace(/\bYou\b/g, me);
}

// expectSameLines: every page names the same waiting and passed players.
async function expectSameLines(pages: Page[], names: string[], want: string): Promise<void> {
  for (const [i, page] of pages.entries()) {
    const me = names[i] ?? '';
    await expect.poll(async () => (await waitLine(page, me, 'Waiting for')) + (await waitLine(page, me, 'Passed'))).toBe(want);
  }
}

test('every player sees whom the stack waits for (V-11)', async ({ browser }) => {
  const names = ['Abe', 'Bea', 'Cal', 'Dot'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const me = names[pages.indexOf(active)] ?? '';
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  await loot.click();
  await active.getByRole('menu').getByRole('menuitem', { name: 'Play' }).click();
  const stack = active.getByRole('complementary', { name: 'Stack' });
  for (let i = 0; i < 10 && (await stack.locator('.stack-row').count()) === 0; i++) {
    const option = active.locator('.choice .choice-option button').first();
    if (await option.isVisible()) {
      await option.click();
    }
    await active.waitForTimeout(200);
  }
  await expect(stack).toContainText(/Waiting for|: /);
  const want = (await waitLine(active, me, 'Waiting for')) + (await waitLine(active, me, 'Passed'));
  expect(want).not.toBe('');
  await expectSameLines(pages, names, want);
  await shot(active, 'stack-waiting');

  // A pass moves the player to "Passed" on every table.
  const pass = passButton(active);
  test.skip(!(await pass.isEnabled()), 'the loot card asked something first');
  await pass.click();
  // Unless everyone else was skipped and the item resolved.
  await expect
    .poll(async () => (await stack.locator('.stack-row').count()) === 0 || (await waitLine(active, me, 'Passed')).includes(me))
    .toBe(true);
  await expectSameLines(pages, names, (await waitLine(active, me, 'Waiting for')) + (await waitLine(active, me, 'Passed')));
});

test('the stack panel: icons, colored names, and a hover with card and owner (15.21)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ike', 'Jud'], { url: fast() });
  const active = await activePage(pages, click);
  const stack = active.getByRole('complementary', { name: 'Stack' });
  // An empty stack says so, and Pass stays in the panel (S-08, S-09).
  await expect(stack).toContainText('Nothing on the stack');
  await expect(passButton(active)).toBeVisible();

  // Three items, as the server would send them: a loot card, a roll,
  // and damage. Only the screen changes; nobody acts on it.
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
      target: { player: number; is_player?: boolean };
      ability: { card: string; index: number };
      roll_for: { card: string; index: number };
    }
    interface S {
      game: { view: { stack: Item[] | null; turn: { active: number } } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const game = store.getState().game;
    if (!game) {
      throw new Error('no game');
    }
    const me = game.view.turn.active;
    const other = 1 - me;
    const none = { card: '', index: 0 };
    const stack: Item[] = [
      {
        id: 1,
        kind: 0,
        controller: me,
        source: 0,
        card: 'bomb',
        text: 'Deal 1 damage to a monster or player.',
        target: { player: other, is_player: true },
        ability: { card: 'bomb', index: 0 },
        roll_for: none,
      },
      { id: 2, kind: 3, controller: other, source: 0, roll: 5, target: { player: -1 }, ability: none, roll_for: none },
      { id: 3, kind: 4, controller: me, source: 0, amount: 1, target: { player: other, is_player: true }, ability: none, roll_for: none },
    ];
    store.setState({ game: { ...game, view: { ...game.view, stack } } });
  });

  const rows = stack.locator('.stack-row');
  await expect(rows).toHaveCount(3);
  // Newest on top; the target in red, the player in their color (S-01, S-03).
  await expect(rows.nth(0)).toHaveAttribute('aria-label', /^(Ike|Jud) deals 1 damage to (Ike|Jud)$/);
  await expect(rows.nth(0).locator('.stack-target')).toBeVisible();
  await expect(rows.nth(1).locator('.die')).toHaveAttribute('aria-label', 'Rolled 5');
  await expect(rows.nth(2)).toHaveAttribute('aria-label', /^(Ike|Jud) plays Bomb!?$/);

  // Hover the loot row: its card large, the used ability, the owner (S-06, S-07).
  const owner = (await rows.nth(2).locator('.stack-player').textContent()) ?? '';
  await rows.nth(2).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', /^Bomb/);
  await expect(zoom.locator('.card-zoom-ability')).toHaveText('Deal 1 damage to a monster or player.');
  await expect(zoom.locator('.card-zoom-owner')).toHaveText(owner);
  await shot(active, 'stack-panel-hover');
  await active.mouse.move(1200, 420);
  await expect(zoom).toBeHidden();
  await shot(active, 'stack-panel');
});
