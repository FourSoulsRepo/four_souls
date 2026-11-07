// Fight marks and played cards (15.31): every player can follow a fight
// and another player's loot card.
import { expect, test } from './fixtures';
import { startGame } from './game';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, playLoot, throwButton, turnAction, untilAttackRoll } from './turn';
import { cardName, tr, trRe } from './locale';

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
  await untilAttackRoll(pages, click, active);

  // Attacker and monster both have the red frame, on every screen.
  for (const page of [active, other]) {
    await expect(page.locator('.player.active .card.fight')).toHaveCount(1);
    await expect(page.locator('.board-monsters .card.fight')).toHaveCount(1);
  }
  // The target is named on the other's stack.
  const monster = await other.locator('.board-monsters .card.fight').getAttribute('aria-label');
  await expect(other.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-target').first()).toHaveText(monster ?? '');
  await shot(other, 'fight-mark');

  // The roll is seen by the other player too.
  await click(active, throwButton(active));
  await expect(other.getByRole('group', { name: tr('stack.last_roll') })).toBeVisible();
});

// The card is A Penny!, by a saved game: a deal may hold no playable loot
// card. It asks nothing and waits on the stack for its player, whose
// character was charged as the turn started (19.126). The play goes
// through the real server, which the 20.9 tests in playback.spec.ts
// fake (19.128).
test("another player's loot card shows large, then by the stack (15.31)", async ({ browser }) => {
  // A repeat (--repeat-each) gets its own names: the saved games list
  // keeps the earlier repeats' games.
  const r = test.info().repeatEachIndex;
  const names = ['Xan', 'Yol'].map((n) => (r > 0 ? `${n}${String(r)}` : n));
  // Animations on: the big card is an effect.
  const pages = await startGame(browser, names);
  await activePage(pages, click);
  await saveAndContinue(pages, names.join(', '), 'a_penny');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const player = names[pages.indexOf(active)] ?? '';
  const name = cardName('a_penny');
  // The player's own play is never shown back to them large (20.9): the
  // page notes any big card from now on.
  await active.evaluate(() => {
    const w = window as unknown as { bigSeen: boolean };
    w.bigSeen = false;
    new MutationObserver(() => {
      if (document.querySelector('[data-testid="played-pop"], .played')) {
        w.bigSeen = true;
      }
    }).observe(document.body, { childList: true, subtree: true });
  });
  await playLoot(active, active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name }).first(), click);

  // The other sees it large in the middle, named as the player's play.
  const pop = other.getByTestId('played-pop');
  await expect(pop).toBeVisible({ timeout: 10_000 });
  await expect(pop.getByRole('img', { name })).toBeVisible();
  await expect(pop).toContainText(tr('effect.plays', { player }));
  await other.waitForTimeout(400); // let it grow for the picture
  await shot(other, 'played-pop');
  // Then it shrinks into the stack, and stays large by it while it waits.
  const kept = other.getByRole('figure', { name: tr('stack.played_label', { player, card: name }) });
  await expect(pop).toHaveCount(0, { timeout: 8000 });
  await expect(kept).toBeVisible();
  await expect(kept.getByRole('img', { name })).toBeVisible();
  await expect(other.getByRole('complementary', { name: tr('stack.title') })).toContainText(trRe('stack.loot', { card: name }));
  await shot(other, 'played-card');
  expect(await active.evaluate(() => (window as unknown as { bigSeen: boolean }).bigSeen)).toBe(false);
});
