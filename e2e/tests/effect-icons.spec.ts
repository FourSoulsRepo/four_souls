// Effect icons for every effect (18.31, ADR 037 item 3, UC-12): a room's
// lasting effect shows on the players it hits only; a played card's
// effect stays after the card is gone; above 4 icons a "+K" chip opens
// the character's zoom. Saved games reach the exact states.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { repeatNames, startGame } from './game';
import { cardName, tr } from './locale';
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
    shields?: { target: { player: number; is_player: boolean }; amount?: number; source?: number }[];
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

test('a room shows its icon on the players it hits only (18.31)', async ({ browser }) => {
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
  // Tax for the Mighty is the room: "When any player controls a soul,
  // players who control the most souls ... must pay each other player
  // 2¢ to attack." The active player gets the only soul: it hits them.
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, ({ state }) => {
      const room = state.rooms[0]?.cards.at(-1);
      const obj = room === undefined ? undefined : state.objects[room];
      if (!obj) {
        throw new Error('no room in play');
      }
      for (const other of state.objects) {
        if (other.card === 'tax_for_the_mighty') {
          other.card = obj.card; // one of a kind
        }
      }
      obj.card = 'tax_for_the_mighty';
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
  const hit = seen.characters[seat]?.id ?? -1;
  const label = `${cardName('tax_for_the_mighty')}: ${tr('stats.effect.attack_tax', { n: 2 })}`;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-room');
  for (const page of pages) {
    // The room lies with the rooms, not in the player's area: an icon,
    // with its number (ADR 037 items 3.1, 3.5).
    const icon = characterOf(page, hit).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', label);
    await expect(icon.locator('.card-source-number')).toHaveText('+2');
    // The other player has no soul: not the most, no icon.
    for (const [other, c] of seen.characters.entries()) {
      if (other !== seat) {
        await expect(characterOf(page, c.id).locator('.card-source')).toHaveCount(0);
      }
    }
  }
  // The icon's zoom is the room's card, with what it does.
  await characterOf(active, hit).locator('.card-source').hover();
  await expect(active.getByTestId('card-zoom')).toContainText(tr('stats.effect.attack_tax', { n: 2 }));
  await shot(active, 'effects-room-zoom');
});

test("a played card's effect keeps its icon after the card is gone (18.31)", async ({ browser }) => {
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
  const label = `${cardName('soul_heart')}: ${tr('stats.this_turn_effect', { effect: tr('stats.effect.shield', { n: 1 }) })}`;
  await active.mouse.move(1, 1);
  await shot(active, 'effects-gone-card');
  for (const page of pages) {
    const icon = characterOf(page, own).locator('.card-source');
    await expect(icon).toHaveCount(1);
    await expect(icon).toHaveAttribute('aria-label', label);
    // A shield is no change of a number: no number (ADR 037 item 3.5).
    await expect(icon.locator('.card-source-number')).toHaveCount(0);
  }
});

test('above 4 icons a "+K" chip opens the character\'s zoom (18.31)', async ({ browser }) => {
  const names = repeatNames(['Max', 'Mia']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Six effects of this turn on the active player, each from its own
  // card, and one from a hidden card: 4 icons, then "+2".
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
    await expect(card.getByRole('img', { name: chipName, exact: true })).toHaveText('+2');
  }
  // The chip zooms the character itself, as hovering the card does
  // (ADR 037 items 3.7, 3.9); 18.32 adds the list of every effect.
  await characterOf(active, character.id).getByRole('img', { name: chipName, exact: true }).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.getByRole('img', { name: cardName(character.card), exact: true })).toBeVisible();
  await shot(active, 'effects-chip-zoom');
});
