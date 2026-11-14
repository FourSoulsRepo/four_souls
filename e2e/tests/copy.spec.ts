// The copy mark (19.106): an item that uses another card's abilities
// shows a mark with the copied card's art, for both players. Hover on
// the mark zooms the copied card; the item's zoom names it (19.111).
import { expect, test } from './fixtures';
import { startGame } from './game';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';
import { cardName, tr } from './locale';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

test('an item that gains another item\'s abilities shows the copy mark (19.106)', async ({ browser }) => {
  const pages = await startGame(browser, ['Nia', 'Oz'], { url: fast() });
  await activePage(pages, click);
  // Brimstone gains The D6's abilities till end of turn.
  await saveAndContinue(pages, 'Nia, Oz', 'bomb', 'brimstone', false, {}, 'the_d6');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  // One mark, on the item, named for the card it copies.
  const mark = active.locator('.card-copy');
  await expect.soft(mark).toHaveCount(1);
  await expect.soft(mark).toHaveAttribute('aria-label', tr('badge.copy', { card: cardName('the_d6') }));
  await shot(active, 'copy-badge');
  // The mark takes the pointer: hover zooms the copied card (19.111).
  const zoom = active.getByTestId('card-zoom');
  await mark.hover();
  await expect.soft(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', cardName('the_d6'));
  await shot(active, 'copy-badge-zoom');
  // The item's own zoom names the copied card first.
  const item = active.locator('.card:has(.card-copy)');
  await item.hover({ position: { x: 10, y: 80 } });
  await expect.soft(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', cardName('brimstone'));
  await expect.soft(zoom.locator('.card-zoom-stats li').first()).toHaveText(tr('badge.copy', { card: cardName('the_d6') }));
  await shot(active, 'copy-card-zoom');
  // The other player sees the same mark.
  await expect.soft(other.locator('.card-copy')).toHaveCount(1);
});

// GAINED_END_OF_TURN is the engine's GainEndOfTurn (state.go).
const GAINED_END_OF_TURN = 1;

interface GainSave {
  state: {
    objects: { card: string; gained?: string; gained_till?: number; more_gained?: { card: string; till: number }[] }[];
    monsters: { cards: number[] }[];
  };
}

test('a monster that gains other monsters\' abilities shows a mark for each (19.157)', async ({ browser }) => {
  const pages = await startGame(browser, ['Vex', 'Kai'], { url: fast() });
  await activePage(pages, click);
  // Trick's 3–4 on Hard, twice: the first monster gains Pooter's and
  // Leaper's abilities till end of turn (R-TRICK-13, Q-46).
  await saveEditContinue(pages, 'Vex, Kai', (game) => {
    editSave(game, (raw) => {
      const { state } = raw as unknown as GainSave;
      const top = state.monsters[0]?.cards.at(-1);
      const monster = top === undefined ? undefined : state.objects[top];
      if (!monster) {
        throw new Error('monster slot 1 is empty');
      }
      monster.gained = 'pooter';
      monster.gained_till = GAINED_END_OF_TURN;
      monster.more_gained = [{ card: 'leaper', till: GAINED_END_OF_TURN }];
    });
  });
  const active = await activePage(pages, click);
  const slot = active.locator('[data-zone="monster-slot-0"]');
  // Two marks on the monster, in the order gained.
  const marks = slot.locator('.card-copy');
  await expect.soft(marks).toHaveCount(2);
  await expect.soft(marks.nth(0)).toHaveAttribute('aria-label', tr('badge.copy', { card: cardName('pooter') }));
  await expect.soft(marks.nth(1)).toHaveAttribute('aria-label', tr('badge.copy', { card: cardName('leaper') }));
  await shot(active, 'monster-copy-badges');
  // The monster's zoom names both copied cards first.
  const zoom = active.getByTestId('card-zoom');
  const card = slot.locator('.card').first();
  const box = await card.boundingBox();
  await card.hover({ position: { x: (box?.width ?? 60) / 2, y: (box?.height ?? 80) * 0.8 } });
  const lines = zoom.locator('.card-zoom-stats li');
  await expect.soft(lines.nth(0)).toHaveText(tr('badge.copy', { card: cardName('pooter') }));
  await expect.soft(lines.nth(1)).toHaveText(tr('badge.copy', { card: cardName('leaper') }));
  await shot(active, 'monster-copy-zoom');
  // The other player sees the same marks.
  const other = pages.find((p) => p !== active);
  await expect.soft(other?.locator('[data-zone="monster-slot-0"] .card-copy') ?? marks).toHaveCount(2);
});
