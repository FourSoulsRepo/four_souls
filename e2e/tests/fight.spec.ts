// Fight marks and played cards (15.31): every player can follow a fight
// and another player's loot card.
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, answerPlay, click, playLoot, settle, throwButton, turnAction } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

test('a second player sees the fight mark and the attack roll (15.31)', async ({ browser }) => {
  const pages = await startGame(browser, ['Vea', 'Wim'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(other.locator('.card.fight')).toHaveCount(0);
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);

  // Attacker and monster both have the red frame, on every screen.
  for (const page of [active, other]) {
    await expect(page.locator('.player.active .card.fight')).toHaveCount(1);
    await expect(page.locator('.board-monsters .card.fight')).toHaveCount(1);
  }
  // The target is named on the other's stack.
  const monster = await other.locator('.board-monsters .card.fight').getAttribute('aria-label');
  await expect(other.getByRole('complementary', { name: 'Stack' }).locator('.stack-target').first()).toHaveText(monster ?? '');
  await shot(other, 'fight-mark');

  // The roll is seen by the other player too.
  await click(active, throwButton(active));
  await expect(other.getByRole('group', { name: 'Last roll' })).toBeVisible();
});

test("another player's loot card shows large, then by the stack (15.31)", async ({ browser }) => {
  // Animations on: the pop is an effect.
  const pages = await startGame(browser, ['Xan', 'Yol']);
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  const name = (await loot.getAttribute('aria-label')) ?? '';
  await playLoot(active, loot, click);

  // The other sees it pop large in the middle, then it stays by the stack.
  // A card may ask its mode or target first; the pop may come and go
  // while that answer waits for the stack row, so both run at once (19.58).
  const pop = other.getByTestId('played-pop');
  await Promise.all([answerPlay(active, click), expect(pop).toBeVisible({ timeout: 10_000 })]);
  await expect(pop.getByRole('img', { name })).toBeVisible();
  await other.waitForTimeout(400); // let it grow for the picture
  await shot(other, 'played-pop');
  const kept = other.getByRole('figure', { name: new RegExp(`plays ${name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}$`) });
  await expect(kept).toBeVisible();
  await expect(kept.getByRole('img', { name })).toBeVisible();
  await shot(other, 'played-card');
  // The player's own play is not shown back to them.
  await expect(active.locator('.played')).toHaveCount(0);
});
