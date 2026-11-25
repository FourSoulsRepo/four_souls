// Effect icons (18.31, 18.34, ADR 037 item 3, UC-12): only what the table
// hides. A room's lasting effect gets no icon, its number shows in the
// badge sum; a played card's effect stays after the card is gone; all
// shields share one icon; above 4 icons a "⋯ K" chip opens the
// character's zoom. Saved games reach the exact states.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { repeatNames, startGame } from './game';
import { cardName, tr, trRe } from './locale';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle, stackHolds, holds } from './turn';

test.setTimeout(150_000);

// The engine's numbers (state.go, statics.go): the monster deck, the
// in-play zone, the soul role, and the stats the tests give.
const MONSTER_DECK = 2;
const ZONE_IN_PLAY = 4;
const ROLE_SOUL = 5;
const STAT_PLAYER_ATK = 0;
const STAT_ROLL = 6;
const STAT_LOOT_PLAYS = 8;
const STAT_PURCHASES = 10;
const STAT_DAMAGE_TAKEN = 36;

interface IconSave {
  state: {
    turn: { active: number };
    players: { character: number; in_play: number[] }[];
    objects: { card: string; zone: { kind: number }; role?: number; controller: number }[];
    rooms: { cards: number[] }[];
    decks: number[][];
    monsters: { cards: number[] }[];
    boosts?: { stat: number; player: number; amount: number; source?: number }[];
    shields?: { target: { player: number; is_player: boolean }; amount?: number; cap?: number; source?: number }[];
  };
}

function edit(game: string, change: (save: IconSave) => void) {
  editSave(game, (save) => {
    change(save as unknown as IconSave);
  });
}

// Seen is the page's latest view: each player's character object and
// card, and the loot discard's top card.
interface Seen {
  characters: { id: number; card: string }[];
  lootDiscard: string;
}

async function seenOf(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface S {
      game: { view: { players: { character: { id: number; card: string } }[]; discards: ({ card: string } | null)[] } } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    return {
      characters: (v?.players ?? []).map((p) => ({ id: p.character.id, card: p.character.card })),
      lootDiscard: v?.discards[1]?.card ?? '', // Deck.loot
    };
  });
}

// characterOf is a character card on the table by its object ID.
function characterOf(page: Page, id: number) {
  return page.locator(`[data-zone="obj-${String(id)}"]`).first();
}

// activeSeat is the seat of the page that may end its turn.
function activeSeat(pages: Page[], active: Page): number {
  return pages.indexOf(active);
}

test('a room gets no icon; its number shows in the badge sum (18.34)', async ({ browser }) => {
  const names = repeatNames(['Ros', 'Tam']);
  // The pages offer Requiem (?requiem), with its room deck on.
  const pages = await startGame(browser, names, {
    url: `${appURL()}&requiem`,
    // Bans: Blood Lust and Blind Rage as the first room force attacks
    // before a turn may end; activePage then plays fights for turns
    // (as rooms.spec.ts, 19.152).
    bans: ['blood_lust_r', 'blind_rage'],
    setup: async (host) => {
      await host.getByLabel(tr('host.detailed')).check();
      await host.getByLabel(tr('host.room_deck')).check();
    },
  });
  await activePage(pages, click);
  // Might for the Meek is the room: "Players who control the fewest
  // souls or tied for fewest have +1 ATK and may attack an additional
  // time on their turn." The active player gets the only soul: it hits
  // the other player only.
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, ({ state }) => {
      const room = state.rooms[0]?.cards.at(-1);
      const obj = room === undefined ? undefined : state.objects[room];
      if (!obj) {
        throw new Error('no room in play');
      }
      for (const other of state.objects) {
        if (other.card === 'might_for_the_meek') {
          other.card = obj.card; // one of a kind
        }
      }
      obj.card = 'might_for_the_meek';
      const p = state.turn.active;
      const id = state.decks[MONSTER_DECK]?.pop();
      const soul = id === undefined ? undefined : state.objects[id];
      const player = state.players[p];
      if (id === undefined || !soul || !player) {
        throw new Error('no monster card left to make a soul');
      }
      Object.assign(soul, { zone: { kind: ZONE_IN_PLAY }, role: ROLE_SOUL, controller: p });
      player.in_play.push(id);
    });
  });
  const active = await activePage(pages, click);
  const seen = await seenOf(active);
  const seat = activeSeat(pages, active);
  const hit = seen.characters.find((_, i) => i !== seat)?.id ?? -1;
  const term = `+ 1 ${cardName('might_for_the_meek')}`;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-room');
  for (const page of pages) {
    // The room lies on the table: no icon on anyone (ADR 037 item 3.1).
    for (const c of seen.characters) {
      await expect(characterOf(page, c.id).locator('.card-source')).toHaveCount(0);
    }
    // Its +1 ATK shows in the badge of the player it hits.
    await expect(characterOf(page, hit).locator('.card-atk.up')).toBeVisible();
  }
  // The zoom's ATK sum names the room (item 3.1, Context 8.3); its extra
  // attack is lasting too: no line in the effect list (item 4.1).
  await characterOf(active, hit).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card-zoom-stats')).toContainText(term);
  await expect(zoom.locator('.card-zoom-effects')).toHaveCount(0);
  await shot(active, 'effects-room-zoom');
  // The player with the soul: no change, no sum naming the room.
  await active.mouse.move(1, 1);
  await characterOf(active, seen.characters[seat]?.id ?? -1).hover();
  await expect(zoom).toBeVisible();
  await expect(zoom).not.toContainText(cardName('might_for_the_meek'));
});

test("a played card's effect keeps its icon after the card is gone (18.31, 18.34)", async ({ browser }) => {
  const names = repeatNames(['Jen', 'Jay']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Soul Heart: "Choose a player. Prevent the next 1 damage they would
  // take this turn." The active player picks themself.
  await saveAndContinue(pages, names.join(', '), 'soul_heart');
  const active = await activePage(pages, click);
  await click(active, active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('soul_heart') }).first());
  await click(active, active.locator('.own-character .card.hl-aim'));
  // Wait for the exact state (5.16.10): Soul Heart resolved and in the
  // loot discard, not an empty stack, which also holds before the play
  // reaches it. Its shield lasts the turn.
  await settle(pages, click, async () => (await seenOf(active)).lootDiscard === 'soul_heart' && (await holds(active, stackHolds(0))));
  const seen = await seenOf(active);
  const own = seen.characters[activeSeat(pages, active)]?.id ?? -1;
  const line = tr('stats.effect_from', { effect: tr('stats.this_turn_effect', { effect: tr('stats.effect.shield', { n: 1 }) }), card: cardName('soul_heart') });
  const label = `${tr('stats.shields', { n: 1 })}: ${line}`;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-gone-card');
  for (const page of pages) {
    const icon = characterOf(page, own).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', label);
    // The shields' icon counts them (ADR 037 item 3.15).
    await expect(icon.locator('.card-source-number')).toHaveText('1');
  }
});

test('all shields share one icon, listed in the order they block (18.34)', async ({ browser }) => {
  const names = repeatNames(['Ned', 'Nia']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Three shields on the active player, from three cards on the table:
  // 1 damage (as Soul Heart's), the whole next damage (as The Poop's),
  // and a cap of 1 (as The Jar's). Each blocks one damage instance,
  // oldest first (R-MECH-46, combat.go shield): the icon says 3, not a
  // sum.
  const sources: string[] = [];
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, ({ state }) => {
      const p = state.turn.active;
      const other = state.players.find((_, i) => i !== p);
      const slots = state.monsters.map((s) => s.cards.at(-1));
      const ids = [slots[0], slots[1], other?.character];
      const cards = ids.map((id) => (id === undefined ? undefined : state.objects[id]?.card));
      if (ids.some((id) => id === undefined) || cards.some((c) => c === undefined)) {
        throw new Error('missing cards for the shields');
      }
      sources.push(...cards.map((c) => c ?? ''));
      const target = { player: p, is_player: true };
      state.shields = [
        ...(state.shields ?? []),
        { target, amount: 1, source: ids[0] },
        { target, source: ids[1] },
        { target, cap: 1, source: ids[2] },
      ];
    });
  });
  const active = await activePage(pages, click);
  const seen = await seenOf(active);
  const own = seen.characters[activeSeat(pages, active)]?.id ?? -1;
  const thisTurn = (effect: string) => tr('stats.this_turn_effect', { effect });
  const lines = [
    tr('stats.effect_from', { effect: thisTurn(tr('stats.effect.shield', { n: 1 })), card: cardName(sources[0] ?? '') }),
    tr('stats.effect_from', { effect: thisTurn(tr('stats.effect.shield_all')), card: cardName(sources[1] ?? '') }),
    tr('stats.effect_from', { effect: thisTurn(tr('stats.effect.shield_cap', { n: 1 })), card: cardName(sources[2] ?? '') }),
  ];
  await active.mouse.move(1, 1);
  await shot(active, 'effects-shields');
  for (const page of pages) {
    const icon = characterOf(page, own).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', `${tr('stats.shields', { n: 3 })}: ${lines.join(', ')}`);
    await expect(icon.locator('.card-source-number')).toHaveText('3');
  }
  // Hover: each shield, in the order they block (item 3.15.3).
  await characterOf(active, own).locator('.card-source').hover();
  await expect(active.getByTestId('card-zoom').locator('.card-zoom-stats li')).toHaveText(lines);
  await shot(active, 'effects-shields-zoom');
});

test('above 4 icons a "⋯ K" chip opens the character\'s zoom (18.31, 18.34)', async ({ browser }) => {
  const names = repeatNames(['Max', 'Mia']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Six effects of this turn on the active player, each from its own
  // card, and one from a hidden card: 4 icons, then "⋯ 2".
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, ({ state }) => {
      const p = state.turn.active;
      const me = state.players[p];
      const other = state.players.find((_, i) => i !== p);
      const slots = state.monsters.map((s) => s.cards.at(-1));
      const [own, theirs] = [me?.in_play[0], other?.in_play[0]];
      if (!me || !other || own === undefined || theirs === undefined || slots[0] === undefined || slots[1] === undefined) {
        throw new Error('missing cards for the effects');
      }
      const boost = (stat: number, amount: number, source?: number) => ({ stat, player: p, amount, ...(source === undefined ? {} : { source }) });
      state.boosts = [
        ...(state.boosts ?? []),
        boost(STAT_PLAYER_ATK, 1, own),
        boost(STAT_DAMAGE_TAKEN, 1, slots[0]),
        boost(STAT_ROLL, 1, slots[1]),
        boost(STAT_LOOT_PLAYS, 1, other.character),
        boost(STAT_PURCHASES, 1),
      ];
      state.shields = [...(state.shields ?? []), { target: { player: p, is_player: true }, amount: 1, source: theirs }];
    });
  });
  const active = await activePage(pages, click);
  const seen = await seenOf(active);
  const character = seen.characters[activeSeat(pages, active)] ?? { id: -1, card: '' };
  const chipName = tr('badge.more_effects', { n: 2 });
  await active.mouse.move(1, 1);
  await shot(active, 'effects-many');
  for (const page of pages) {
    const card = characterOf(page, character.id);
    await expect(card.locator('.card-source')).toHaveCount(4);
    // A pill with no card art, unlike an icon's "+N" (ADR 037 item 3.16).
    const chip = card.getByRole('img', { name: chipName, exact: true });
    await expect(chip).toHaveText('⋯ 2');
    await expect(chip.locator('img')).toHaveCount(0);
  }
  // The chip zooms the character itself, as hovering the card does
  // (ADR 037 items 3.7, 3.9), with the list of every effect (18.32).
  await characterOf(active, character.id).getByRole('img', { name: chipName, exact: true }).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.getByRole('img', { name: cardName(character.card), exact: true })).toBeVisible();
  await shot(active, 'effects-chip-zoom');
  const list = zoom.getByRole('list', { name: tr('zoom.effects_you') });
  // This turn's effects, in the server's order; the source card's name
  // after what it does; a boost with no card is from "a hidden card"
  // (ADR 037 items 4.1–4.3, 1.6).
  const thisTurn = list.getByRole('list', { name: tr('zoom.effects_turn') }).getByRole('listitem');
  await expect(thisTurn).toHaveText([
    tr('stats.effect_from', { effect: tr('stats.effect.loot_plays', { amount: '+1', n: 1 }), card: cardName(seen.characters.find((c) => c.id !== character.id)?.card ?? '') }),
    tr('stats.effect_from', { effect: tr('stats.effect.purchases', { amount: '+1', n: 1 }), card: tr('stats.hidden_card') }),
    trRe('stats.effect_from', { effect: tr('stats.effect.shield', { n: 1 }) }),
  ]);
  // No repeats (ADR 037 item 4.4): ATK, rolls and damage taken stay in
  // their sum lines only.
  await expect(list).not.toContainText(tr('stats.effect.atk', { amount: '+1' }));
  await expect(list).not.toContainText(tr('stats.effect.roll', { amount: '+1' }));
  await expect(list).not.toContainText(tr('stats.effect.damage_taken', { amount: '+1' }));
});

test('a lasting effect of the own area has no icon and no zoom line (18.34)', async ({ browser }) => {
  const names = repeatNames(['Kim', 'Kai']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Steamy Sale!: "Shop items you purchase cost 5¢ less." It becomes the
  // active player's starting item: a card in their own area (group A).
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, ({ state }) => {
      const item = state.players[state.turn.active]?.in_play[0];
      const obj = item === undefined ? undefined : state.objects[item];
      if (!obj) {
        throw new Error('the active player has no item');
      }
      for (const other of state.objects) {
        if (other.card === 'steamy_sale') {
          other.card = obj.card; // one of a kind
        }
      }
      obj.card = 'steamy_sale';
    });
  });
  const active = await activePage(pages, click);
  const seen = await seenOf(active);
  const seat = activeSeat(pages, active);
  const character = seen.characters[seat] ?? { id: -1, card: '' };
  for (const [i, page] of pages.entries()) {
    const card = characterOf(page, character.id);
    // The card lies beside the character: no icon (ADR 037 item 3.1).
    await expect(card.locator('.card-source')).toHaveCount(0);
    await page.mouse.move(1, 1);
    await card.hover();
    const zoom = page.getByTestId('card-zoom');
    await expect(zoom).toBeVisible();
    await shot(page, i === seat ? 'effects-own-zoom' : 'effects-opponent-zoom');
    // Nor a line in the zoom's list, for its own player or another
    // (item 4.1): with nothing hidden, no list at all.
    await expect(zoom.locator('.card-zoom-effects')).toHaveCount(0);
    await expect(zoom).not.toContainText(tr('stats.effect.shop_price', { amount: '−5' }));
  }
});
