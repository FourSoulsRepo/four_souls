// Choosing cards (20.18, ADR 027 "Choosing cards"): a look at a deck's
// top cards, a pick, and what follows it. The card is given by a saved
// game, so no test depends on the deal (CLAUDE.md 5.16).
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { activePage, click, passButton, playLoot } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

// names are the players' names; a repeat (--repeat-each) gets its own,
// as the saved games list keeps the earlier repeats' games.
function names(base: string[]): string[] {
  const r = test.info().repeatEachIndex;
  return base.map((n) => (r > 0 ? `${n}${String(r)}` : n));
}

// playMoon gives the active player XVIII The Moon by a saved game and
// plays it; the pages pass until its question shows. It returns the
// active page and the other one.
async function playMoon(pages: Page[], who: string[]): Promise<{ active: Page; other: Page }> {
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'xviii_the_moon');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const moon = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('xviii_the_moon') });
  await playLoot(active, moon.first(), click);
  const dialog = active.getByRole('dialog', { name: tr('ask.one_on_top') });
  await expect
    .poll(
      async () => {
        if (await dialog.isVisible()) {
          return true;
        }
        for (const page of pages) {
          const pass = passButton(page);
          if ((await pass.isVisible()) && (await pass.isEnabled())) {
            await pass.click().catch(() => undefined);
          }
        }
        return false;
      },
      { timeout: 20_000 },
    )
    .toBe(true);
  return { active, other };
}

// history is the page's history lines.
function history(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    interface S {
      history: { text: string }[];
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().history.map((e) => e.text);
  });
}

test('after The Moon\'s pick no look window repeats the cards; the history keeps the look (19.130)', async ({ browser }) => {
  const who = names(['Ama', 'Bex']);
  const pages = await startGame(browser, who, { url: fast() });
  const { active } = await playMoon(pages, who);
  const dialog = active.getByRole('dialog', { name: tr('ask.one_on_top') });
  const options = dialog.locator('.choice-option');
  await expect(options).toHaveCount(5);
  const shown = await options.locator('.card').evaluateAll((els) => els.map((el) => el.getAttribute('aria-label') ?? ''));
  await options.first().getByRole('button').first().click();
  await expect(dialog).toHaveCount(0);
  // The look's events came: the history tells each card looked at.
  await expect.poll(async () => (await history(active)).filter((l) => shown.some((c) => l === tr('history.look.you', { card: c }))).length).toBeGreaterThan(0);
  await expect(active.getByRole('dialog', { name: tr('look.title') })).toHaveCount(0);
});
