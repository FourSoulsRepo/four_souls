// Info panels (7.9): history, turn badge, what is left this turn.
import { expect, test } from '@playwright/test';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, canEndTurn, click, settle } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

test('the history shows only what the viewer may see', async ({ browser }) => {
  const names = ['Wes', 'Xan'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const history = (p: typeof active) => p.getByRole('complementary', { name: 'History', exact: true });
  // The full history opens from the strip (15.22).
  for (const page of pages) {
    await page.getByRole('button', { name: 'History', exact: true }).click();
  }
  for (const [i, page] of pages.entries()) {
    const other = names.at(1 - i) ?? '';
    const text = (await history(page).textContent()) ?? '';
    // Others' loots never name the card.
    expect(text).not.toMatch(new RegExp(`${other} loots (?!a card)`));
  }
  // The active player looted: their own history names the card.
  await expect(history(active)).toContainText(/You loot \S/);
  const other = pages.find((p) => p !== active);
  const activeName = names[pages.indexOf(active)] ?? '';
  await expect(history(other ?? active)).toContainText(`${activeName} loots a card`);
  await shot(active, 'history');
});

test('the turn badge, what is left, and hiding the history', async ({ browser }) => {
  const names = ['Yul', 'Zed'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const activeName = names[pages.indexOf(active)] ?? '';
  for (const page of pages) {
    await expect(page.getByRole('status', { name: 'Turn' })).toHaveText(new RegExp(`^Turn \\d+ · ${activeName} · \\w`));
  }
  await expect(active.getByLabel('Left this turn')).toHaveText(/Loot \d+ · Attacks \d+ · Buys \d+/);

  // The full history is closed; the strip opens it, Hide closes it.
  const panel = active.getByRole('complementary', { name: 'History', exact: true });
  const toggle = active.getByRole('button', { name: 'History', exact: true });
  await expect(panel).toHaveCount(0);
  await toggle.click();
  await expect(panel).toBeVisible();
  await expect(toggle).toHaveAttribute('aria-pressed', 'true');
  await panel.getByRole('button', { name: 'Hide' }).click();
  await expect(panel).toHaveCount(0);
  await toggle.click();
  await toggle.click();
  await expect(panel).toHaveCount(0);
});

test('a resolved loot card goes to the history strip (15.22)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ada', 'Ben'], { url: fast() });
  const active = await activePage(pages, click);
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  const name = (await loot.getAttribute('aria-label')) ?? '';
  await loot.click();
  await active.getByRole('menu').getByRole('menuitem', { name: 'Play' }).click();
  // Targets, passes and anything the card asks, until the stack is empty.
  const stackRows = active.getByRole('complementary', { name: 'Stack' }).locator('.stack-row');
  await settle(pages, click, async () => (await stackRows.count()) === 0 && (await canEndTurn(active)));
  for (const page of pages) {
    const strip = page.getByRole('complementary', { name: 'History strip' });
    // The loot is a played spell: a bolt badge (H-08).
    const played = strip.getByLabel(`plays ${name}`).first();
    await expect(played).toBeVisible();
    await expect(played.locator('.strip-spell')).toBeAttached();
  }
});

test('the history strip: icons, badges, and a hover with card and owner (15.22)', async ({ browser }) => {
  const pages = await startGame(browser, ['Cid', 'Dee'], { url: fast() });
  const active = await activePage(pages, click);
  // Resolved items as the server would send them: a loot card, an attack
  // roll, and damage, the newest last. They replace what the deal
  // resolved so far. Only the screen changes.
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
    interface Entry {
      key: number;
      turn: number;
      icon: string;
      text: string;
      card?: string;
      item?: Item;
    }
    interface S {
      game: { view: { turn: { active: number; number: number } } } | null;
      history: Entry[];
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const s = store.getState();
    if (!s.game) {
      throw new Error('no game');
    }
    const me = s.game.view.turn.active;
    const turn = s.game.view.turn.number;
    const other = 1 - me;
    const none = { card: '', index: 0 };
    const items: Item[] = [
      {
        id: 101,
        kind: 0,
        controller: me,
        source: 0,
        card: 'bomb',
        text: 'Deal 1 damage to a monster or player.',
        target: { player: other, is_player: true },
        ability: { card: 'bomb', index: 0 },
        roll_for: none,
      },
      { id: 102, kind: 3, controller: me, source: 0, roll: 4, attack: true, target: { player: -1 }, ability: none, roll_for: none },
      { id: 103, kind: 4, controller: me, source: 0, amount: 2, attack: true, target: { player: other, is_player: true }, ability: none, roll_for: none },
    ];
    const history = items.map((item, i) => ({ key: i + 1, turn, icon: '✦', text: '', ...(item.card ? { card: item.card } : {}), item }));
    store.setState({ history });
  });

  const strip = active.getByRole('complementary', { name: 'History strip' });
  const icons = strip.locator('.strip-item');
  await expect(icons).toHaveCount(3);
  // Newest on top; a sword for the attack, a bolt for the loot (H-08).
  await expect(icons.nth(0)).toHaveAttribute('aria-label', /^Attack: 2 damage to (Cid|Dee)$/);
  await expect(icons.nth(0).locator('.strip-attack')).toBeAttached();
  await expect(icons.nth(1).locator('.die')).toHaveAttribute('aria-label', 'Rolled 4');
  await expect(icons.nth(1).locator('.strip-attack')).toBeAttached();
  await expect(icons.nth(2)).toHaveAttribute('aria-label', /^(Cid|Dee) plays Bomb!?$/);
  await expect(icons.nth(2).locator('.strip-spell')).toBeAttached();

  // Hover the loot icon: its card large, the used ability, the owner (S-11, H-09).
  await icons.nth(2).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', /^Bomb/);
  await expect(zoom.locator('.card-zoom-ability')).toHaveText('Deal 1 damage to a monster or player.');
  await expect(zoom.locator('.card-zoom-owner')).toHaveText(/^(Cid|Dee)$/);
  await shot(active, 'history-strip-hover');
  await active.mouse.move(1000, 420);
  await expect(zoom).toBeHidden();
  await shot(active, 'history-strip');
});
