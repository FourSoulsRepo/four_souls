// The first-player roll (15.43, GS-06): before turn 1 every player
// throws a die, ties throw again, and everyone sees each die and who
// goes first.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { shot } from './shot';
import { seen, throwButton } from './turn';
import { tr, trRe } from './locale';

test.setTimeout(120_000);

function rollList(page: Page) {
  return page.getByRole('complementary', { name: tr('stack.title') }).getByRole('group', { name: tr('first_roll.label') });
}

function rollRow(page: Page, name: string) {
  // A row's name is "Ann: 4", or "Ann: not rolled yet" in the tests' language.
  return rollList(page).getByRole('listitem', { name: new RegExp(`^(${name}: |${trRe('first_roll.not_rolled', { player: name }).source})`) });
}

// recordBanners keeps the text of every banner the page plays.
async function recordBanners(page: Page): Promise<void> {
  await page.evaluate(() => {
    interface S {
      playing: { key: number; effect: { kind: string; text?: string } } | null;
    }
    const w = window as unknown as { fourSouls: { store: { subscribe: (f: (s: S) => void) => void } }; banners: string[] };
    w.banners = [];
    let last = 0;
    w.fourSouls.store.subscribe((s) => {
      if (s.playing && s.playing.key !== last) {
        last = s.playing.key;
        if (s.playing.effect.kind === 'banner') {
          w.banners.push(s.playing.effect.text ?? '');
        }
      }
    });
  });
}

function banners(page: Page): Promise<string[]> {
  return page.evaluate(() => (window as unknown as { banners: string[] }).banners);
}

// bannerShot saves the banner once it has faded in, if one with this
// text shows before stop; a tie may never come.
async function bannerShot(page: Page, text: RegExp, name: string, stop: Promise<void>): Promise<void> {
  const banner = page.locator('.effect-banner', { hasText: text });
  const shown = async () => {
    await banner.waitFor({ timeout: 60_000 });
    await page.waitForFunction(() => Number(getComputedStyle(document.querySelector('.effect-banner') ?? document.body).opacity) > 0.95);
    return true;
  };
  if (await Promise.race([shown().catch(() => false), stop.then(() => false)])) {
    await shot(page, name);
  }
}

test("every player's first-player roll shows before the first turn", async ({ browser }) => {
  const names = ['Ann', 'Bo', 'Cy'];
  // Cain makes his player go first without a roll; this test needs the
  // roll. Eden's starting-item dialog comes right after the roll and
  // covers the History button this test opens.
  const pages = await startGame(browser, names, { roll: false, bans: ['cain', 'eden'] });
  const seats = new Map<number, Page>();
  for (const page of pages) {
    await recordBanners(page);
    seats.set((await seen(page)).you, page);
    await expect(page.getByRole('status', { name: tr('turn_badge.label') })).toHaveText(tr('turn_badge.first_roll'));
    await expect(rollList(page).getByRole('listitem')).toHaveCount(3);
    // Nobody has a turn yet: no End turn, no active seat.
    await expect(page.getByRole('button', { name: new RegExp(`${trRe('turn.end').source}|${trRe('turn.theirs').source}`) })).toHaveCount(0);
  }
  const [ann, bo] = pages;
  if (!ann || !bo) {
    throw new Error('no pages');
  }
  await expect(throwButton(ann)).toBeVisible();
  await expect(bo.getByRole('region', { name: tr('throw.label') })).toContainText(tr('throw.throwing', { player: 'Ann' }));
  await shot(bo, 'first-roll-waiting');

  // Ann throws first; every page shows her die, the same on each.
  await throwButton(ann).click();
  await ann.mouse.move(1, 1);
  await expect(rollRow(bo, 'Ann')).toHaveAttribute('aria-label', /^Ann: [1-6]$/);
  const label = (await rollRow(bo, 'Ann').getAttribute('aria-label')) ?? '';
  for (const page of pages) {
    await expect(rollRow(page, 'Ann')).toHaveAttribute('aria-label', label);
  }
  await shot(bo, 'first-roll-one-die');

  // Shots of the tie and winner banners, taken as they show.
  const winner = trRe('first_roll.by_roll');
  const tie = trRe('first_roll.tie');
  let started = () => undefined as void;
  const stop = new Promise<void>((resolve) => {
    started = resolve;
  });
  const shots = process.env.SCREENSHOTS
    ? Promise.all([bannerShot(bo, tie, 'first-roll-tie', stop), bannerShot(bo, winner, 'first-roll-winner', stop)])
    : Promise.resolve();

  // The others throw, round after round when the lowest ties. Up to a
  // minute: each tie adds a long banner and a round of throws.
  for (let i = 0; i < 600 && (await seen(ann)).turn < 1; i++) {
    const s = await seen(ann);
    const roller = seats.get(s.waiting);
    const button = roller ? throwButton(roller) : null;
    if (!roller || !button || !(await button.isVisible())) {
      await ann.waitForTimeout(100);
      continue;
    }
    await button.click();
    // Off the cards, so no zoom covers the table.
    await roller.mouse.move(1, 1);
  }
  for (const page of pages) {
    await expect(page.getByRole('status', { name: tr('turn_badge.label') })).toHaveText(new RegExp(`^${trRe('turn_badge.turn', { n: 1 }).source}$`), { timeout: 20_000 });
  }
  started();
  await shots;

  // Every page named who goes first, and why; a tie, if any, too.
  const ties: number[] = [];
  for (const page of pages) {
    const shown = await banners(page);
    expect(shown.filter((b) => winner.test(b))).toHaveLength(1);
    ties.push(shown.filter((b) => tie.test(b)).length);
  }
  expect(new Set(ties).size).toBe(1);
  for (const page of pages) {
    await expect(page.getByRole('status', { name: tr('turn_badge.label') })).toHaveText(new RegExp(`^${trRe('turn_badge.turn', { n: 1 }).source.replace('.*?', '(Ann|Bo|Cy)')}$`));
    await expect(rollList(page)).toHaveCount(0);
  }

  // The history keeps every die and who goes first.
  const first = (await bo.getByRole('status', { name: tr('turn_badge.label') }).textContent())?.split(' · ')[1] ?? '';
  for (const page of pages) {
    const you = names[(await seen(page)).you] ?? '';
    await page.getByRole('complementary', { name: tr('history.strip') }).getByRole('button', { name: tr('history.title') }).click();
    const history = page.getByRole('complementary', { name: tr('history.title'), exact: true });
    for (const name of names) {
      const roll = name === you ? new RegExp(`${trRe('history.first_roll.you').source}$`) : new RegExp(`${trRe('history.first_roll.other', { player: name }).source}$`);
      await expect(history.getByRole('listitem').filter({ hasText: roll }).first()).toBeVisible();
    }
    const goes = trRe('first_roll.by_roll', { line: first === you ? tr('first_roll.you_first') : tr('first_roll.first', { player: first }) });
    await expect(history.getByRole('listitem').filter({ hasText: goes })).toHaveCount(1);
    if ((ties[0] ?? 0) > 0) {
      // A tie seen as a banner is in the history too.
      await expect(history.getByRole('listitem').filter({ hasText: new RegExp(`${trRe('first_roll.tie').source}$`) }).first()).toBeVisible();
    }
  }
  await shot(bo, 'first-roll-history');
});
