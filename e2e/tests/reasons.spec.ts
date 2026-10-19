// Reasons and choice dialogs (15.29): a card or turn action the player
// can't use says why, in its zoom and where it is clicked; a choice
// shows big cards over the dimmed board, and "View board" hides it.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, turnAction } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

const card = (page: Page, id: number) => page.locator(`.card[data-zone="obj-${String(id)}"]`);
const note = (page: Page) => page.locator('.aim-why');

// blocked is the first hand card the server says this player can't
// play now, with its reason.
async function blocked(page: Page): Promise<{ id: number; reason: string }> {
  return page.evaluate(() => {
    interface S {
      game: { view: { players: { hand?: { id: number }[] }[] }; blocked?: { kind: number; object: number; reason: string }[] };
      lobby: { table: { you: number } };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const hand = (s.game.view.players[s.lobby.table.you]?.hand ?? []).map((c) => c.id);
    const b = (s.game.blocked ?? []).find((x) => x.kind === 3 && hand.includes(x.object));
    return { id: b?.object ?? -1, reason: b?.reason ?? '' };
  });
}

test("an unusable card's zoom and a click on it say why (15.29)", async ({ browser }) => {
  const pages = await startGame(browser, ['Basil', 'Juno'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }

  // Not their turn: the other player's loot waits for it.
  const b = await blocked(other);
  expect(b.reason).toBe('Loot is played on your turn');
  const loot = card(other, b.id);
  await loot.hover();
  const why = other.getByTestId('card-zoom').locator('.card-zoom-why');
  await expect(why).toHaveText(b.reason);
  await shot(other, 'zoom-why');

  // A click on it says the same where it was clicked, then fades.
  await loot.click();
  await expect(note(other)).toHaveText(b.reason);
  await shot(other, 'why-click');
  await expect(note(other)).toHaveCount(0, { timeout: 5000 });

  // A dim turn icon says why too, and stays dim.
  const attack = turnAction(other, 'Attack');
  await expect(attack).toBeDisabled();
  await attack.click({ force: true });
  await expect(note(other)).toHaveText('Not your turn');
});

test('"View board" hides a choice and brings it back (15.29)', async ({ browser }) => {
  const [corin] = await startGame(browser, ['Corin', 'Meda'], { url: fast() });
  if (!corin) {
    throw new Error('no pages');
  }
  // A choice of a monster on the table, as the server sends it.
  const monster = await corin.evaluate(() => {
    interface Card {
      id: number;
      card: string;
    }
    interface S {
      game: { view: { monsters: { top?: Card }[]; waiting: unknown; options?: unknown }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const m = game.view.monsters.find((x) => x.top)?.top ?? { id: 0, card: '' };
    game.view.waiting = { kind: 4, player: you, purpose: 1, options: [m.card] };
    game.view.options = [{ object: m.id, player: 0 }];
    game.allowed = [{ player: you, kind: 4, choice: 0 }];
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
    return m.id;
  });
  const dialog = corin.getByRole('dialog', { name: 'Choose what to attack' });
  await expect(dialog).toBeVisible();
  // The board under it is dimmed and takes no clicks.
  await expect(corin.locator('.choice-backdrop')).toBeVisible();
  await expect(card(corin, monster).click({ trial: true, timeout: 500 })).rejects.toThrow();
  await shot(corin, 'choice-big');

  await dialog.getByRole('button', { name: 'View board' }).click();
  await expect(dialog).toHaveCount(0);
  const peek = corin.getByRole('group', { name: 'Viewing the board' });
  await expect(peek).toContainText('Choose what to attack');
  // The target glows on the table and answers there.
  await expect(card(corin, monster)).toHaveClass(/hl-pick/);
  await card(corin, monster).click({ trial: true });
  await shot(corin, 'choice-view-board');

  await peek.getByRole('button', { name: 'Back to the choice' }).click();
  await expect(dialog).toBeVisible();
  await expect(peek).toHaveCount(0);
});
