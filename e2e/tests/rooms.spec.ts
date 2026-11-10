// Rooms on the table (14.13, R-SETUP-16): the host plays with the room
// deck, every player sees the room landscape, the active player uses
// it, and puts it into discard in the end phase for the next room.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr, trRe } from './locale';
import { editSave, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, endTurn, settle } from './turn';

test.setTimeout(120_000);

// ROOM_DECK is the engine's RoomDeck (state.go).
const ROOM_DECK = 3;

interface RoomSave {
  state: {
    turn: { active: number; monster_died?: boolean };
    players: { character: number }[];
    objects: { card: string }[];
    rooms: { cards: number[] }[];
    decks: number[][];
  };
}

// giveRoom makes the room in play Blood Donation and the room deck's
// top Samson's Blessing, as if a monster died this turn, so the end
// phase asks for the room step (R-TURN-12). The active character
// becomes Isaac: Blood Donation's 1 HP must not kill a 1 HP character.
function giveRoom(save: RoomSave) {
  const { state } = save;
  const swap = (id: number | undefined, card: string) => {
    const obj = id === undefined ? undefined : state.objects[id];
    if (!obj) {
      throw new Error(`no object for ${card}`);
    }
    // The card is one of a kind: another object with it takes the old one.
    for (const other of state.objects) {
      if (other.card === card) {
        other.card = obj.card;
      }
    }
    obj.card = card;
  };
  swap(state.rooms[0]?.cards.at(-1), 'blood_donation');
  swap(state.decks[ROOM_DECK]?.at(-1), 'samsons_blessing');
  const character = state.players[state.turn.active]?.character;
  const isaac = character === undefined ? undefined : state.objects[character];
  if (isaac) {
    isaac.card = 'isaac';
  }
  state.turn.monster_died = true;
}

// cents is the player's cents on the page's latest view.
function cents(page: Page, player: number): Promise<number> {
  return page.evaluate((player) => {
    interface S {
      game: { view: { players: { id: number; cents: number }[] } } | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return s.game?.view.players.find((p) => p.id === player)?.cents ?? -1;
  }, player);
}

function rooms(page: Page) {
  return page.getByRole('group', { name: tr('board.rooms') });
}

test('a room on the table: use it, change it in the end phase (14.13)', async ({ browser }) => {
  const names = ['Ivo', 'Jun'];
  // The pages offer Requiem (?requiem), ticked with the other sets; its
  // room deck is a detailed setup option, off by default.
  const pages = await startGame(browser, names, {
    url: `${appURL()}&requiem`,
    setup: async (host) => {
      await host.getByLabel(tr('host.detailed')).check();
      const box = host.getByLabel(tr('host.room_deck'));
      await expect(box).not.toBeChecked();
      await box.check();
    },
  });
  // Every player sees one room, landscape, over the room deck.
  for (const page of pages) {
    const room = rooms(page).locator('.card.kind-room').first();
    await expect(room).toBeVisible();
    const box = await room.boundingBox();
    expect(box && box.width > box.height).toBe(true);
    await expect(rooms(page).getByRole('group', { name: trRe('deck.room') })).toBeVisible();
  }
  await activePage(pages, click);
  await saveEditContinue(pages, names.join(', '), (game) => {
    editSave(game, (save) => {
      giveRoom(save as unknown as RoomSave);
    });
  });
  const active = await activePage(pages, click);
  const seat = pages.indexOf(active);
  const donation = rooms(active).getByRole('button', { name: cardName('blood_donation'), exact: true });
  for (const page of pages) {
    await expect(rooms(page).getByLabel(cardName('blood_donation'), { exact: true })).toBeVisible();
  }
  await shot(active, 'rooms-table');

  // The active player uses the room: pay 1 HP, gain 3¢ (R-CARD-36).
  const before = await cents(active, seat);
  await click(active, donation);
  await click(active, active.getByRole('menu').getByRole('menuitem').first());
  await settle(pages, click, async () => (await cents(active, seat)) === before + 3);

  // Ending the turn asks for the room step; the room goes, the next
  // room comes from the deck (R-TURN-12).
  await endTurn(active);
  const choice = active.locator('.choice');
  await settle(pages, click, async () => (await choice.getByRole('button', { name: cardName('blood_donation') }).count()) > 0);
  await shot(active, 'rooms-step');
  await click(active, choice.getByRole('button', { name: cardName('blood_donation') }).first());
  for (const page of pages) {
    await expect(rooms(page).getByLabel(cardName('samsons_blessing'), { exact: true })).toBeVisible();
    await expect(rooms(page).getByRole('group', { name: tr('deck.room_discard', { n: 1 }) })).toBeVisible();
  }
});
