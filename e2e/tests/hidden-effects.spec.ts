// Hidden effects that are no stats (18.37, ADR 037 items 3.17, 3.18, 8,
// UC-12): a skip, a copied loot, a game-wide effect on its user's
// character only, and Red Candle's chosen number as a die on its card.
// Saved games give the cards, so nothing depends on the deal.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { repeatNames, startGame } from './game';
import { cardName, tr } from './locale';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, playLoot, seen, settle } from './turn';

test.setTimeout(150_000);

// The engine's numbers (state.go): the monster deck and the discard zone.
const MONSTER_DECK = 2;
const ZONE_DISCARD = 2;

// Seen is the page's latest view: each player's character, the loot
// discard's top card, and the active player's items.
interface Seen {
  characters: { id: number; card: string }[];
  lootDiscard: string;
  items: { id: number; card: string; charged: boolean; dice?: number[] }[];
}

async function seenOf(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface Item {
      id: number;
      card: string;
      charged: boolean;
      dice?: number[];
    }
    interface S {
      game: {
        view: {
          turn: { active: number };
          players: { id: number; character: { id: number; card: string }; in_play: Item[] | null }[];
          discards: ({ card: string } | null)[];
        };
      } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    const active = v?.players.find((p) => p.id === v.turn.active);
    return {
      characters: (v?.players ?? []).map((p) => ({ id: p.character.id, card: p.character.card })),
      lootDiscard: v?.discards[1]?.card ?? '', // Deck.loot
      items: (active?.in_play ?? []).map((c) => ({ id: c.id, card: c.card, charged: c.charged, ...(c.dice ? { dice: c.dice } : {}) })),
    };
  });
}

// characterOf is a character card on the table by its object ID.
function characterOf(page: Page, id: number) {
  return page.locator(`[data-zone="obj-${String(id)}"]`).first();
}

// ownItem is a card in the page's own items.
function ownItem(page: Page, card: string) {
  return page.getByRole('group', { name: tr('table.your_items') }).getByRole('button', { name: cardName(card), exact: true });
}

// quiet: the stack is empty and the page may act again.
async function quiet(page: Page): Promise<boolean> {
  const s = await seen(page);
  return s.ready && s.stack.length === 0;
}

const thisTurn = (effect: string) => tr('stats.this_turn_effect', { effect });
const everyone = (effect: string) => tr('stats.for_everyone', { effect });

test("Famine's skip: an icon on the skipping player, its line under \"Until used\" (18.37)", async ({ browser }) => {
  const names = repeatNames(['Fay', 'Gus']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Famine: "When this dies, the active player skips their next turn."
  // It died: it lies in the monster discard, and the active player owes
  // one skipped turn from it (Player.SkipSources, 18.36).
  await saveEditContinue(pages, names.join(', '), (game) => {
    editSave(game, (save) => {
      const state = (
        save as unknown as {
          state: {
            turn: { active: number };
            players: { skip_turns?: number; skip_sources?: number[] }[];
            objects: { card: string; zone: { kind: number; deck?: number } }[];
            decks: number[][];
            discards: (number[] | null)[];
          };
        }
      ).state;
      const id = state.decks[MONSTER_DECK]?.pop();
      const obj = id === undefined ? undefined : state.objects[id];
      const player = state.players[state.turn.active];
      if (id === undefined || !obj || !player) {
        throw new Error('no monster card left for Famine');
      }
      for (const other of state.objects) {
        if (other.card === 'famine') {
          other.card = obj.card; // one of a kind
        }
      }
      Object.assign(obj, { card: 'famine', zone: { kind: ZONE_DISCARD, deck: MONSTER_DECK } });
      // An empty discard is null in the save.
      state.discards[MONSTER_DECK] = [...(state.discards[MONSTER_DECK] ?? []), id];
      player.skip_turns = 1;
      player.skip_sources = [id];
    });
  });
  const active = await activePage(pages, click);
  const seenNow = await seenOf(active);
  const seat = pages.indexOf(active);
  const own = seenNow.characters[seat]?.id ?? -1;
  const other = seenNow.characters.find((_, i) => i !== seat)?.id ?? -1;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-skip');
  // The skip outlasts this turn: no "this turn" in its line, and no
  // number on the icon for one turn (ADR 037 item 3.12).
  const line = tr('stats.effect.skip_turn', { n: 1 });
  for (const page of pages) {
    const icon = characterOf(page, own).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', `${cardName('famine')}: ${line}`);
    await expect(icon.locator('.card-source-number')).toHaveCount(0);
    await expect(characterOf(page, other).locator('.card-source')).toHaveCount(0);
  }
  // The zoom files it under "Until used", not "This turn" (18.37).
  await characterOf(active, own).hover();
  const zoom = active.getByTestId('card-zoom');
  const list = zoom.getByRole('list', { name: tr('zoom.effects_you') });
  await expect(list.getByRole('list', { name: tr('zoom.effects_later') }).getByRole('listitem')).toHaveText([
    tr('stats.effect_from', { effect: line, card: cardName('famine') }),
  ]);
  await expect(list.getByRole('list', { name: tr('zoom.effects_turn') })).toHaveCount(0);
  await shot(active, 'effects-skip-zoom');
});

test('A Lucky Penny!: rewards doubled for every player, on its user only (18.37)', async ({ browser }) => {
  const names = repeatNames(['Hal', 'Ivy']);
  // The pages offer Requiem (?requiem), so they can continue a game
  // with A Lucky Penny!'s set; the deal stays the Base Game's.
  const pages = await startGame(browser, names, {
    url: `${appURL()}&anim=off&requiem`,
    setup: (host) => host.getByRole('group', { name: tr('host.sets') }).getByLabel(tr('set.requiem')).uncheck(),
  });
  await activePage(pages, click);
  await saveAndContinue(pages, names.join(', '), 'a_lucky_penny', undefined, false, { r: 'Requiem' });
  const active = await activePage(pages, click);
  // A Lucky Penny!: "Gain 1¢. Rewards are doubled till end of turn."
  await playLoot(active, active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('a_lucky_penny') }).first(), click);
  // Wait for the exact state (5.16.10): the penny resolved and in the
  // loot discard, the stack empty.
  await settle(pages, click, async () => (await seenOf(active)).lootDiscard === 'a_lucky_penny' && (await quiet(active)));
  const seenNow = await seenOf(active);
  const seat = pages.indexOf(active);
  const own = seenNow.characters[seat]?.id ?? -1;
  const other = seenNow.characters.find((_, i) => i !== seat)?.id ?? -1;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-everyone');
  // Game-wide: on the character of the player who played it, and its
  // line says it is for every player (ADR 037 item 3.18).
  const line = thisTurn(everyone(tr('stats.effect.double_rewards')));
  for (const page of pages) {
    const icon = characterOf(page, own).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', `${cardName('a_lucky_penny')}: ${line}`);
    await expect(icon.locator('.card-source-number')).toHaveCount(0);
    await expect(characterOf(page, other).locator('.card-source')).toHaveCount(0);
  }
  await characterOf(active, own).hover();
  const list = active.getByTestId('card-zoom').getByRole('list', { name: tr('zoom.effects_you') });
  await expect(list.getByRole('list', { name: tr('zoom.effects_turn') }).getByRole('listitem')).toHaveText([
    tr('stats.effect_from', { effect: everyone(tr('stats.effect.double_rewards')), card: cardName('a_lucky_penny') }),
  ]);
  await shot(active, 'effects-everyone-zoom');
});

test('Blank Card: the next loot copied shows as an icon (18.37)', async ({ browser }) => {
  const names = repeatNames(['Jo', 'Kit']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  await saveAndContinue(pages, names.join(', '), 'bomb', 'blank_card');
  const active = await activePage(pages, click);
  // Blank Card: "↷: The next time you play a non-trinket, non-ambush
  // loot card this turn, copy it." A click opens its menu.
  await click(active, ownItem(active, 'blank_card'));
  await click(active, active.getByRole('menu').getByRole('menuitem').first());
  // Wait for the exact state (5.16.10): the card spent, its ability
  // resolved, the stack empty.
  await settle(pages, click, async () => (await seenOf(active)).items.some((c) => c.card === 'blank_card' && !c.charged) && (await quiet(active)));
  const seenNow = await seenOf(active);
  const own = seenNow.characters[pages.indexOf(active)]?.id ?? -1;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-blank-card');
  const line = thisTurn(tr('stats.effect.copy_next_loot'));
  for (const page of pages) {
    const icon = characterOf(page, own).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', `${cardName('blank_card')}: ${line}`);
    await expect(icon.locator('.card-source-number')).toHaveCount(0);
  }
});

test("Red Candle: the chosen number is a die on its card, on no character (18.37, ADR 037 item 8)", async ({ browser }) => {
  const names = repeatNames(['Lu', 'Mo']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  await saveAndContinue(pages, names.join(', '), 'bomb', 'red_candle', false, { fsp2: 'Four Souls+ V2' });
  const active = await activePage(pages, click);
  // Red Candle: "↷: Before a dice is rolled, choose a number. Till the
  // end of turn, each time that number is rolled, deal 1 damage to a
  // monster or player." The number is asked as the ability resolves.
  await click(active, ownItem(active, 'red_candle'));
  await click(active, active.getByRole('menu').getByRole('menuitem').first());
  const text = 'card.red_candle.choose_a_number';
  const dialog = active.getByRole('dialog', { name: tr(text) });
  await settle(pages, click, () => dialog.isVisible());
  await dialog.getByRole('button', { name: '4', exact: true }).click();
  await settle(pages, click, async () => (await seenOf(active)).items.some((c) => c.card === 'red_candle' && c.dice?.includes(4) === true) && (await quiet(active)));
  const candle = (await seenOf(active)).items.find((c) => c.card === 'red_candle')?.id ?? -1;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-red-candle');
  for (const page of pages) {
    // The die with the chosen 4 on the card (15.67), for everyone.
    const card = page.locator(`.card[data-zone="obj-${String(candle)}"]`);
    await expect(card.getByLabel(tr('badge.kept_die', { n: 4 }))).toBeVisible();
    // Not on any character (ADR 037 item 8.1).
    await expect(page.locator('.card-source')).toHaveCount(0);
  }
});

test("Sibling Rivalry: the chosen player shows it, the user does not (18.38)", async ({ browser }) => {
  const names = repeatNames(['Ann', 'Bo']);
  // The pages offer Requiem (?requiem), so they can continue a game
  // with Sibling Rivalry's set; the deal stays the Base Game's.
  const pages = await startGame(browser, names, {
    url: `${appURL()}&anim=off&requiem`,
    setup: (host) => host.getByRole('group', { name: tr('host.sets') }).getByLabel(tr('set.requiem')).uncheck(),
  });
  await activePage(pages, click);
  await saveAndContinue(pages, names.join(', '), 'bomb', 'sibling_rivalry', false, { r: 'Requiem' });
  const active = await activePage(pages, click);
  const seat = pages.indexOf(active);
  const seenNow = await seenOf(active);
  const own = seenNow.characters[seat]?.id ?? -1;
  const chosenSeat = seenNow.characters.findIndex((_, i) => i !== seat);
  const chosen = seenNow.characters[chosenSeat]?.id ?? -1;
  // Sibling Rivalry: "↷: Choose a non-active player. The next time the
  // active player declares an attack this turn, the chosen player must
  // make an attack roll after each attack roll the active player makes
  // for the attack. If that monster dies this attack, the chosen player
  // also gains the rewards." Its one ability asks for a target, so a
  // click aims from it (15.25).
  await click(active, ownItem(active, 'sibling_rivalry'));
  await expect(characterOf(active, chosen)).toHaveClass(/hl-aim/);
  await click(active, characterOf(active, chosen));
  // Wait for the exact state (5.16.10): the card spent, its ability
  // resolved, the stack empty.
  await settle(pages, click, async () => (await seenOf(active)).items.some((c) => c.card === 'sibling_rivalry' && !c.charged) && (await quiet(active)));
  await active.mouse.move(1, 1);
  await shot(active, 'effects-sibling-rivalry');
  // The line names the active player; it says when itself, so no
  // "this turn", and has no number (ADR 037 item 3.12).
  const line = tr('stats.effect.next_attack_helper', { player: names[seat] ?? '' });
  for (const page of pages) {
    const icon = characterOf(page, chosen).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', `${cardName('sibling_rivalry')}: ${line}`);
    await expect(icon.locator('.card-source-number')).toHaveCount(0);
    await expect(characterOf(page, own).locator('.card-source')).toHaveCount(0);
  }
  await characterOf(active, chosen).hover();
  const list = active.getByTestId('card-zoom').getByRole('list', { name: tr('zoom.effects_on', { player: names[chosenSeat] ?? '' }) });
  await expect(list.getByRole('list', { name: tr('zoom.effects_turn') }).getByRole('listitem')).toHaveText([
    tr('stats.effect_from', { effect: line, card: cardName('sibling_rivalry') }),
  ]);
  await shot(active, 'effects-sibling-rivalry-zoom');
});
