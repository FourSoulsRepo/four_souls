// Reasons and choice dialogs (15.29): a card or turn action the player
// can't use says why, in its zoom and where it is clicked; a choice
// shows big cards over the dimmed board, and "View board" hides it.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, turnAction } from './turn';
import { tr } from './locale';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

const card = (page: Page, id: number) => page.locator(`.card[data-zone="obj-${String(id)}"]`);
const note = (page: Page) => page.locator('.aim-why');

// hand lists the ids of the player's hand cards, from the store.
async function hand(page: Page): Promise<number[]> {
  return page.evaluate(() => {
    interface S {
      game: { view: { players: { hand?: { id: number }[] }[] } };
      lobby: { table: { you: number } };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return (s.game.view.players[s.lobby.table.you]?.hand ?? []).map((c) => c.id);
  });
}

// explained is the reason key the server gave for card id, after a click
// asked why (15.29); '' until the answer comes.
async function explained(page: Page, id: number): Promise<string> {
  return page.evaluate((object) => {
    interface S {
      explained: { object: number; reason: { key: string } }[];
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return s.explained.find((b) => b.object === object)?.reason.key ?? '';
  }, id);
}

test("an unusable card's zoom and a click on it say why (15.29)", async ({ browser }) => {
  const pages = await startGame(browser, ['Basil', 'Juno'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }

  // Not their turn: the other player's loot waits for it. Its reason
  // comes when asked: a click on the card asks, the zoom says it after.
  const [id] = await hand(other);
  if (id === undefined) {
    throw new Error('no hand card');
  }
  const loot = card(other, id);
  await loot.click();
  await expect.poll(() => explained(other, id)).toBe('refuse.loot_your_turn');
  // The answer shows at the card, where it was clicked, then fades.
  await expect(note(other)).toHaveText(tr('refuse.loot_your_turn'));
  await shot(other, 'why-click');
  await expect(note(other)).toHaveCount(0, { timeout: 5000 });

  // Its zoom says it from now on, for as long as the answer holds.
  await other.mouse.move(1, 1);
  await loot.hover();
  const why = other.getByTestId('card-zoom').locator('.card-zoom-why');
  await expect(why).toHaveText(tr('refuse.loot_your_turn'));
  await shot(other, 'zoom-why');

  // A dim turn icon says why too, and stays dim.
  const attack = turnAction(other, 'Attack');
  await expect(attack).toBeDisabled();
  await attack.click({ force: true });
  await expect(note(other)).toHaveText(tr('refuse.not_your_turn'));
});

test('"View board" hides a choice and brings it back (15.29)', async ({ browser }) => {
  const [corin] = await startGame(browser, ['Corin', 'Meda'], { url: fast() });
  if (!corin) {
    throw new Error('no pages');
  }
  // A card's target prompt on a monster, as the server sends it. What
  // to attack is picked on the table, with no dialog (18.11).
  const monster = await corin.evaluate(() => {
    interface Card {
      id: number;
      card: string;
    }
    interface S {
      game: { step: number; view: { monsters: { top?: Card }[]; waiting: unknown; options?: unknown }; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const m = game.view.monsters.find((x) => x.top)?.top ?? { id: 0, card: '' };
    game.view.waiting = { kind: 4, player: you, purpose: 6, options: [{ key: 'ask.card', args: { card: { card: m.card } } }] };
    game.view.options = [{ object: m.id, player: 0 }];
    game.allowed = [{ player: you, kind: 4, choice: 0 }];
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
    return m.id;
  });
  const dialog = corin.getByRole('dialog', { name: tr('choice.title.target') });
  await expect(dialog).toBeVisible();
  // The board under it is dimmed and takes no clicks.
  await expect(corin.locator('.choice-backdrop')).toBeVisible();
  await expect(card(corin, monster).click({ trial: true, timeout: 500 })).rejects.toThrow();
  await shot(corin, 'choice-big');

  await dialog.getByRole('button', { name: tr('choice.view_board') }).click();
  await expect(dialog).toHaveCount(0);
  const peek = corin.getByRole('group', { name: tr('choice.viewing') });
  await expect(peek).toContainText(tr('choice.title.target'));
  // The target glows on the table and answers there.
  await expect(card(corin, monster)).toHaveClass(/hl-pick/);
  await card(corin, monster).click({ trial: true });
  await shot(corin, 'choice-view-board');

  await peek.getByRole('button', { name: tr('choice.back') }).click();
  await expect(dialog).toBeVisible();
  await expect(peek).toHaveCount(0);
});
