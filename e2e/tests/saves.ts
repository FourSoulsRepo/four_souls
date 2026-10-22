// Saved games in tests: save a running game, change it on disk, and
// continue it, so a test gets the card it needs (15.25).
import { expect, type Page } from '@playwright/test';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { serverPort } from './server';

// gameID is the id of the page's table, which names its save file.
export async function gameID(page: Page): Promise<string> {
  return page.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: { getState: () => { lobby: { table: { game: string } | null } } } } }).fourSouls.store;
    return store.getState().lobby.table?.game ?? '';
  });
}

// giveCard rewrites a saved game: the active player's first hand card
// becomes the card, and their first card in play becomes item, if given
// (15.92). The test server keeps saves under .cache/e2e.
export function giveCard(game: string, card: string, item?: string) {
  const file = join(import.meta.dirname, '..', '..', '.cache', 'e2e', String(serverPort()), 'saves', `${game}.fssave`);
  const save = JSON.parse(readFileSync(file, 'utf8')) as {
    state: { turn: { active: number }; players: { hand: number[]; in_play: number[] }[]; objects: { card: string }[] };
  };
  const { state } = save;
  const player = state.players[state.turn.active];
  const id = player?.hand[0];
  const obj = id === undefined ? undefined : state.objects[id];
  if (!obj) {
    throw new Error('the active player has no hand card');
  }
  obj.card = card;
  if (item !== undefined) {
    const at = player?.in_play[0];
    const held = at === undefined ? undefined : state.objects[at];
    if (!held) {
      throw new Error('the active player has no card in play');
    }
    held.card = item;
  }
  writeFileSync(file, JSON.stringify(save));
}

// saveAndContinue: the host saves, the active player gets the card (and
// the item in play), then everyone comes back to it.
export async function saveAndContinue(pages: Page[], names: string, card = 'bomb', item?: string) {
  const [host, ...guests] = pages;
  if (!host) {
    throw new Error('no pages');
  }
  const game = await gameID(host);
  await host.getByRole('button', { name: 'Menu' }).click();
  await host.getByRole('button', { name: 'Save and leave' }).click();
  for (const page of pages) {
    await expect(page.getByText('The game was saved.')).toBeVisible();
  }
  giveCard(game, card, item);
  await host.getByRole('button', { name: 'Back' }).click();
  await host.getByRole('button', { name: 'Host a game' }).click();
  await host.getByRole('button', { name: 'Continue a saved game' }).click();
  await host.getByRole('listitem').filter({ hasText: names }).getByRole('button', { name: 'Continue' }).click();
  for (const page of guests) {
    await page.getByRole('listitem').filter({ hasText: names }).getByRole('button', { name: 'Continue' }).click();
  }
  for (const page of pages) {
    await page.getByRole('button', { name: 'Ready', exact: true }).click();
  }
  for (const page of pages) {
    await expect(page.getByRole('main', { name: /^Turn / })).toBeVisible();
  }
}
