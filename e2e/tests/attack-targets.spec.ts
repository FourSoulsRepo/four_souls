// Requiem's attack targets (14.14): a player while The Siren is out
// (R-ATK-20), an attackable item (Punching Bag, R-MECH-04), and a
// helper's extra attack roll named by its roller (R-ATK-21). The save
// sets the cards, so nothing depends on the deal.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle, tablePicks, throwButton, turnAction, untilAttackRoll } from './turn';

test.setTimeout(150_000);

// The pages offer Requiem (?requiem), so they can continue a game with
// it; the deal stays the Base Game's, unticked as in two-sided.spec.ts.
const opts = () => ({
  url: `${appURL()}&anim=off&requiem`,
  setup: (host: Page) => host.getByRole('group', { name: tr('host.sets') }).getByLabel(tr('set.requiem')).uncheck(),
});

interface Seat {
  id: number;
  hp: number;
  character: { id: number; card: string };
  in_play: { id: number; card: string; damage?: number }[] | null;
}

interface View {
  players: Seat[];
  attack?: { target?: number; odds?: { dc: number } };
  stack?: { kind: number; attack?: boolean; controller: number; roll?: number }[] | null;
  waiting: { kind: number; player: number };
}

// view is the page's latest game view, and its own seat.
function view(page: Page): Promise<{ view: View; you: number }> {
  return page.evaluate(() => {
    interface S {
      game: { view: View } | null;
      lobby: { table: { you: number } | null };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    if (!s.game) {
      throw new Error('no game');
    }
    return { view: s.game.view, you: s.lobby.table?.you ?? 0 };
  });
}

async function seats(page: Page): Promise<{ me: Seat; them: Seat; all: Seat[] }> {
  const { view: v, you } = await view(page);
  const me = v.players.find((p) => p.id === you);
  const them = v.players.find((p) => p.id !== you);
  if (!me || !them) {
    throw new Error('no seats');
  }
  return { me, them, all: v.players };
}

async function pair(pages: Page[]): Promise<[Page, Page]> {
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  return [active, other];
}

// stack is the stack panel.
function stack(page: Page) {
  return page.getByRole('complementary', { name: tr('stack.title') });
}

test('The Siren lets the active player attack a player (14.14, R-ATK-20)', async ({ browser }) => {
  const names = ['Ada', 'Bo'];
  const pages = await startGame(browser, names, opts());
  await activePage(pages, click);
  // The Siren in the first slot; Requiem joins the game for her.
  await saveAndContinue(pages, names.join(', '), 'bomb', undefined, false, { r: 'Requiem' }, undefined, 'the_siren');
  const [active, other] = await pair(pages);
  const victim = names[pages.indexOf(other)] ?? '';
  const { me, them } = await seats(active);

  await turnAction(active, 'Attack').click();
  await settle(pages, click, async () => tablePicks(active).first().isVisible());
  // The other player's character glows as a target, beside The Siren.
  const character = active.locator(`[data-zone="obj-${String(them.character.id)}"]`).first();
  await expect(character).toHaveClass(/hl-pick/);
  await expect(active.getByRole('group', { name: tr('board.monsters') }).getByLabel(cardName('the_siren'), { exact: true })).toHaveClass(/hl-pick/);
  // The own character does not.
  await expect(active.locator(`[data-zone="obj-${String(me.character.id)}"]`).first()).not.toHaveClass(/hl-pick/);

  // The attack arrow aims at the player as at a monster (15.32, 20.23).
  const box = await character.boundingBox();
  if (!box) {
    throw new Error('no character box');
  }
  await active.mouse.move(box.x + box.width / 2, box.y + box.height / 2, { steps: 8 });
  await expect(active.locator('.aim-arrow')).toHaveClass(/\bon\b/);
  await shot(active, 'attack-player-aim');
  await active.mouse.click(box.x + box.width / 2, box.y + box.height / 2);

  // The attack is on the player; their evasion is The Siren's 4.
  await untilAttackRoll(pages, click, active);
  const { view: v } = await view(active);
  expect(v.attack?.target).toBe(them.character.id);
  expect(v.attack?.odds?.dc).toBe(4);
  // Both tables name the player, not their character: the attack and
  // its roll.
  for (const page of pages) {
    await expect(stack(page).locator('.stack-declared')).toContainText(victim);
    await expect(stack(page)).toContainText(tr('stack.roll_against', { player: names[pages.indexOf(active)] ?? '', target: victim }));
    await expect(stack(page)).not.toContainText(cardName(them.character.card));
  }
  await shot(active, 'attack-player-roll');

  // One roll: a hit hurts the player, a miss hurts the attacker.
  await click(active, throwButton(active));
  await settle(pages, click, async () => {
    const now = await seats(active);
    return now.them.hp < them.hp || now.me.hp < me.hp;
  });
});

test('Punching Bag can be attacked (14.14, R-MECH-04)', async ({ browser }) => {
  const names = ['Cy', 'Di'];
  const pages = await startGame(browser, names, opts());
  const first = await activePage(pages, click);
  const seat = pages.indexOf(first);
  // The other player's first item in play becomes Punching Bag.
  await saveEditContinue(pages, names.join(', '), (game) => {
    editSave(game, (save) => {
      const s = save as unknown as {
        sets: string[];
        state: { sets: string[]; turn: { active: number }; players: { in_play: number[] }[]; objects: { card: string }[] };
      };
      const owner = s.state.players.find((_, i) => i !== s.state.turn.active);
      const id = owner?.in_play[0];
      const obj = id === undefined ? undefined : s.state.objects[id];
      if (!obj) {
        throw new Error('the other player has no item');
      }
      obj.card = 'punching_bag';
      s.sets.push('r');
      s.state.sets.push('Requiem');
    });
  });
  const [active, other] = await pair(pages);
  expect(pages.indexOf(active)).toBe(seat);
  const owner = names[pages.indexOf(other)] ?? '';
  const bagView = (await seats(active)).them.in_play?.find((c) => c.card === 'punching_bag');
  if (!bagView) {
    throw new Error('no Punching Bag');
  }

  await turnAction(active, 'Attack').click();
  await settle(pages, click, async () => tablePicks(active).first().isVisible());
  const bag = active.getByRole('group', { name: tr('table.items_of', { player: owner }) }).getByLabel(cardName('punching_bag'), { exact: true });
  await expect(bag).toHaveClass(/hl-pick/);
  await bag.hover();
  await expect(active.locator('.aim-arrow')).toHaveClass(/\bon\b/);
  await shot(active, 'attack-item-aim');
  await click(active, bag);

  // The attack is on the bag, at its own evasion 5.
  await untilAttackRoll(pages, click, active);
  const { view: v } = await view(active);
  expect(v.attack?.target).toBe(bagView.id);
  expect(v.attack?.odds?.dc).toBe(5);
  for (const page of pages) {
    await expect(stack(page).locator('.stack-declared')).toContainText(cardName('punching_bag'));
  }
  await shot(active, 'attack-item-roll');

  // One roll: a hit damages the bag, a miss hurts the attacker.
  const before = (await seats(active)).me.hp;
  await click(active, throwButton(active));
  await settle(pages, click, async () => {
    const now = await seats(active);
    const hit = now.them.in_play?.find((c) => c.id === bagView.id);
    return (hit?.damage ?? 0) > 0 || hit === undefined || now.me.hp < before;
  });
});

test("a helper's extra attack roll shows whose roll it is (14.14, R-ATK-21)", async ({ browser }) => {
  const names = ['Ed', 'Fay'];
  const pages = await startGame(browser, names, opts());
  await activePage(pages, click);
  // Sibling Rivalry in play; a plain monster in the first slot.
  await saveAndContinue(pages, names.join(', '), 'bomb', 'sibling_rivalry', false, { r: 'Requiem' }, undefined, 'fatty');
  const [active, other] = await pair(pages);
  const helper = names[pages.indexOf(other)] ?? '';
  const { me, them } = await seats(active);

  // ↷: choose the other player; they roll after each attack roll.
  const rivalry = active.getByRole('group', { name: tr('table.your_items') }).getByRole('button', { name: cardName('sibling_rivalry'), exact: true });
  // Its one ability asks for a target, so a click aims from it (15.25).
  await click(active, rivalry);
  await expect(active.locator(`[data-zone="obj-${String(them.character.id)}"]`).first()).toHaveClass(/hl-aim/);
  await click(active, active.locator(`[data-zone="obj-${String(them.character.id)}"]`).first());
  await settle(pages, click, async () => {
    const { view: v } = await view(active);
    return (v.stack?.length ?? 0) === 0 && v.waiting.player === me.id;
  });

  await turnAction(active, 'Attack').click();
  await settle(pages, click, async () => tablePicks(active).first().isVisible());
  await click(active, active.getByRole('group', { name: tr('board.monsters') }).getByLabel(cardName('fatty'), { exact: true }));
  await untilAttackRoll(pages, click, active);
  await click(active, throwButton(active));

  // After the active player's roll, the helper's own roll waits for them.
  await settle(pages, click, async () => {
    const { view: v } = await view(other);
    return (v.stack ?? []).some((it) => it.kind === 3 && it.attack === true && it.controller === them.id && !it.roll) && (await throwButton(other).count()) > 0;
  });
  // Both tables say it is the helper's roll.
  for (const page of pages) {
    await expect(page.getByRole('region', { name: tr('throw.label') })).toContainText(page === other ? '' : helper);
    await expect(stack(page)).toContainText(helper);
  }
  await shot(other, 'attack-helper-roll-own');
  await shot(active, 'attack-helper-roll-others');
});
