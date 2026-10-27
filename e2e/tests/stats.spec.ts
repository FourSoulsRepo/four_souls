// Stat sources (15.92, V-41): a changed stat shows in color, its sum in
// the zoom, and each effect of this turn as an icon on the card.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, endTurn, expectNextTurn, settle, throwButton, turnAction, turnNumber } from './turn';
import { cardName, lineRe, tr, trRe } from './locale';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

// strengthOnSelf plays XI. Strength on the active player's own character
// and resolves it.
async function strengthOnSelf(pages: Page[], active: Page) {
  await click(active, active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('xi_strength') }).first());
  await click(active, active.locator('.own-character .card.hl-aim'));
  await expect(active.getByRole('complementary', { name: tr('stack.title') })).toContainText(trRe('stack.loot', { card: cardName('xi_strength') }));
  await settle(pages, click, stackDone(active));
}

// stackDone: the stack is empty for the active player.
const stackDone = (page: Page) => async () => (await page.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-row').count()) === 0;

test('a lasting item changes the badge; XI. Strength adds an icon (15.92)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ada', 'Ben'], { url: fast() });
  await activePage(pages, click);
  // The starting item becomes Brimstone: a lasting +1 ATK.
  await saveAndContinue(pages, 'Ada, Ben', 'xi_strength', 'brimstone');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const own = active.locator('.own-character .card');
  // A lasting item: a green ATK badge, no icon (owner, 2026-10-22).
  await expect.soft(own.locator('.card-atk.up')).toBeVisible();
  await expect.soft(own.locator('.card-source')).toHaveCount(0);
  await own.hover();
  await expect.soft(active.getByTestId('card-zoom').locator('.card-zoom-stats')).toContainText(`+ 1 ${cardName('brimstone')}`);
  await shot(active, 'stats-lasting');

  // XI. Strength on the player themself: +1 ATK this turn.
  await strengthOnSelf(pages, active);
  await active.mouse.move(0, 0);
  // One icon, on the character for both players, and the sum names it.
  await expect.soft(own.locator('.card-source')).toHaveCount(1);
  await expect.soft(own.locator('.card-source')).toHaveAttribute('aria-label', `${cardName('xi_strength')}: ${tr('stats.this_turn_effect', { effect: tr('stats.effect.atk', { amount: '+1' }) })}`);
  await expect.soft(other.locator('.opp-character .card .card-source')).toHaveCount(1);
  await shot(active, 'stats-temporary');
  await own.hover();
  await expect.soft(active.getByTestId('card-zoom').locator('.card-zoom-stats')).toContainText(`+ 1 ${cardName('brimstone')} + 1 ${cardName('xi_strength')} (${tr('stats.this_turn')})`);
  await shot(active, 'stats-sum');
  // Hovering the icon zooms XI. Strength, saying how long it lasts.
  await own.locator('.card-source').first().hover();
  await expect.soft(active.getByTestId('card-zoom')).toContainText(tr('stats.this_turn_effect', { effect: tr('stats.effect.atk', { amount: '+1' }) }));
  await shot(active, 'stats-icon');
});

test('a monster shows its changed dice and why (15.92)', async ({ browser }) => {
  const pages = await startGame(browser, ['Cy', 'Dee'], { url: fast() });
  await activePage(pages, click);
  // Baby Haunt: monsters have +1 DC on its owner's turn.
  await saveAndContinue(pages, 'Cy, Dee', 'bomb', 'baby_haunt');
  const active = await activePage(pages, click);
  const monster = active.locator('.board-slots .card').filter({ has: active.locator('.card-dc') }).first();
  await expect.soft(monster.locator('.card-dc.up')).toBeVisible();
  await monster.hover();
  await expect.soft(active.getByTestId('card-zoom').locator('.card-zoom-stats')).toContainText(`+ 1 ${cardName('baby_haunt')}`);
  await shot(active, 'stats-monster');
});

for (const anim of [true, false]) {
  test(`effect icons go at the end of the turn, animations ${anim ? 'on' : 'off'} (15.93)`, async ({ browser }) => {
    // Own names per variant: the saved games list finds a save by them,
    // and both variants may run at once on one server (19.80).
    const names = anim ? ['Eli', 'Fox'] : ['Kip', 'Lux'];
    const pages = await startGame(browser, names, { url: anim ? appURL() : fast() });
    await activePage(pages, click);
    await saveAndContinue(pages, names.join(', '), 'xi_strength');
    const active = await activePage(pages, click);
    await strengthOnSelf(pages, active);
    const icons = active.locator('.own-character .card .card-source');
    await expect(icons).toHaveCount(1);
    const turn = await turnNumber(active);
    await endTurn(active, click);
    await expectNextTurn(pages, click, turn);
    if (anim) {
      // The icon shrinks away; then it is gone.
      await expect(icons).toHaveCount(0);
    } else {
      // No fade: gone with the turn, at once.
      expect(await icons.count()).toBe(0);
    }
  });
}

test('a next-roll icon shows till the attack roll uses it (15.93)', async ({ browser }) => {
  const pages = await startGame(browser, ['Gil', 'Hap'], { url: fast() });
  await activePage(pages, click);
  // +1 to the next attack roll from the first monster, as Round Worm gives.
  await saveAndContinue(pages, 'Gil, Hap', 'bomb', undefined, true);
  const active = await activePage(pages, click);
  const icon = active.locator('.own-character .card .card-source.next-roll');
  await expect(icon).toHaveCount(1);
  await expect(icon).toHaveAttribute('aria-label', new RegExp(`: ${lineRe(tr('stats.effect.next_attack_roll', { amount: '+1' })).source}$`));
  await active.mouse.move(0, 0);
  await shot(active, 'stats-next-roll');
  await icon.hover();
  await expect(active.getByTestId('card-zoom')).toContainText(tr('stats.effect.next_attack_roll', { amount: '+1' }));
  await active.mouse.move(0, 0);
  // Attack; the roll uses the boost, and the icon goes.
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
  await settle(pages, click, async () => (await icon.count()) === 0);
  await expect(icon).toHaveCount(0);
});

test('the open attack says what the roll needs; Meat! makes it one lower (15.94)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ivo', 'Jo'], { url: fast() });
  await activePage(pages, click);
  // The starting item becomes Meat!: a lasting +1 to attack rolls.
  await saveAndContinue(pages, 'Ivo, Jo', 'bomb', 'meat');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
  for (const page of [active, other]) {
    const line = page.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-odds');
    await expect(line).toContainText(`− 1 ${cardName('meat')}`);
    const text = (await line.textContent()) ?? '';
    // "Hits on 3+ (DC 4 − 1 Meat!)", in the tests' language.
    // The tokens stand for patterns, put in after escaping.
    const odds = lineRe(tr('stats.odds_hit', { hits: 'QHITSQ', dc: 'QDCQ', terms: ` − 1 ${cardName('meat')}` })).source;
    const from = lineRe(tr('stats.hits_from', { n: 'QNQ' })).source.replace('QNQ', '(\\d)');
    const m = new RegExp(`^${odds.replace('QHITSQ', from).replace('QDCQ', '(\\d)')}`).exec(text);
    const always = new RegExp(`^${odds.replace('QHITSQ', lineRe(tr('stats.always_hits')).source).replace('QDCQ', '[12]')}`).test(text);
    expect(m !== null || always, text).toBe(true);
    if (m) {
      expect(Number(m[1])).toBe(Number(m[2]) - 1);
    }
  }
  await shot(active, 'stats-odds');
  // The thrown roll names Meat! and its result.
  await click(active, throwButton(active));
  await expect(active.getByRole('complementary', { name: tr('stack.title') })).toContainText(
    new RegExp(`${trRe('stack.roll_against_n').source}.*? ${lineRe(tr('stack.roll_result', { changes: `+1 ${cardName('meat')}`, n: 'QNQ' })).source.replace('QNQ', '\\d')}`),
  );
  await shot(active, 'stats-roll');
});
