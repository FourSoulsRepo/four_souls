// Choosing cards (20.18, ADR 027 "Choosing cards"): a look at a deck's
// top cards, a pick, and what follows it. The card is given by a saved
// game, so no test depends on the deal (CLAUDE.md 5.16).
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
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
async function playMoon(pages: Page[], who: string[], before: (active: Page, other: Page) => Promise<void> = () => Promise.resolve()): Promise<{ active: Page; other: Page }> {
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'xviii_the_moon');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await before(active, other);
  const moon = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('xviii_the_moon') });
  await playLoot(active, moon.first(), click);
  await passUntil(pages, active.getByRole('dialog', { name: tr('ask.one_on_top') }));
  return { active, other };
}

// passUntil passes on every page until the locator shows.
async function passUntil(pages: Page[], shown: ReturnType<Page['locator']>): Promise<void> {
  await expect
    .poll(
      async () => {
        if (await shown.isVisible()) {
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
}

// listenSounds records the file of every sound played from now on, in
// window.played (15.50). Each sound still plays.
async function listenSounds(page: Page): Promise<void> {
  await page.evaluate(() => {
    const w = window as unknown as { played: string[] };
    w.played = [];
    const play = HTMLMediaElement.prototype.play;
    HTMLMediaElement.prototype.play = function (this: HTMLMediaElement) {
      w.played.push(new URL(this.src).pathname);
      return play.call(this);
    };
  });
}

function played(page: Page, sound: string): Promise<number> {
  return page.evaluate((s) => (window as unknown as { played: string[] }).played.filter((p) => p === `/cardimg/assets/sounds/${s}.mp3`).length, sound);
}

// watchDeal keeps where each choice card first shows, and the deck's
// centre then, in window.dealtAt: a dealt card starts on its deck.
async function watchDeal(page: Page, deck: string): Promise<void> {
  await page.evaluate((deck) => {
    const w = window as unknown as { dealtAt: { x: number; y: number; dx: number; dy: number }[] };
    w.dealtAt = [];
    const centre = (el: Element) => {
      const r = el.getBoundingClientRect();
      return { x: r.left + r.width / 2, y: r.top + r.height / 2 };
    };
    new MutationObserver(() => {
      const cards = [...document.querySelectorAll('.choice-card')];
      const pile = document.querySelector(`[data-zone="deck-${deck}"]`);
      if (w.dealtAt.length === 0 && cards.length > 0 && pile) {
        const d = centre(pile);
        w.dealtAt = cards.map((c) => ({ ...centre(c), dx: d.x, dy: d.y }));
      }
    }).observe(document.body, { childList: true, subtree: true });
  }, deck);
}

// dealtFar is how far the first dealt card was from its deck's centre
// when it showed, in px.
function dealtFar(page: Page): Promise<number> {
  return page.evaluate(() => {
    const at = (window as unknown as { dealtAt: { x: number; y: number; dx: number; dy: number }[] }).dealtAt[0];
    return at ? Math.hypot(at.x - at.dx, at.y - at.dy) : Infinity;
  });
}

// nameOf is the page's own player name.
function nameOf(page: Page): Promise<string> {
  return page.evaluate(() => {
    interface S {
      game: { seats?: { seat: number; name: string }[] } | null;
      lobby: { table: { you: number } | null };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return s.game?.seats?.find((x) => x.seat === s.lobby.table?.you)?.name ?? '';
  });
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

test('The Moon deals its cards from the loot deck; a pick sends them back; others see backs (20.18)', async ({ browser }) => {
  const who = names(['Cye', 'Dov']);
  // Animations on: the deal and the flight back are the test.
  const pages = await startGame(browser, who, { url: appURL() });
  const { active, other } = await playMoon(pages, who, async (a, o) => {
    await listenSounds(a);
    await listenSounds(o);
    await watchDeal(a, 'loot');
  });
  const dialog = active.getByRole('dialog', { name: tr('ask.one_on_top') });
  await active.waitForTimeout(150);
  await shot(active, '20.18-moon-deal-mid');
  // Each card starts on the deck and lands in one row; the card is the
  // button, named by the card, with no name button under it.
  expect(await dealtFar(active)).toBeLessThan(40);
  await expect(active.locator('.choice.dealing')).toHaveCount(0, { timeout: 5000 });
  const cards = dialog.locator('.choice-option > button.choice-card');
  await expect(cards).toHaveCount(5);
  await expect(dialog.locator('.choice-option > button:not(.choice-card)')).toHaveCount(0);
  const tops = await cards.evaluateAll((els) => els.map((el) => Math.round(el.getBoundingClientRect().top)));
  expect(new Set(tops).size).toBe(1);
  await expect.poll(() => played(active, 'choice-deal')).toBeGreaterThan(0);
  await shot(active, '20.18-moon-choice');

  // The other player sees five small backs over the loot deck.
  const player = await nameOf(active);
  const backs = other.locator('[data-zone="deck-loot"]').getByRole('img', { name: tr('deck.looking', { player, n: 5 }) });
  await expect(backs).toBeVisible();
  await expect(backs.locator('.deck-look-back')).toHaveCount(5);
  await other.waitForTimeout(800);
  await shot(other, '20.18-moon-choice-other');

  // Hover lifts the card; its zoom waits; a right click opens it at once.
  const first = cards.first();
  await first.hover();
  await expect.poll(() => first.evaluate((el) => getComputedStyle(el).translate)).toBe('0px -8px');
  await shot(active, '20.18-moon-choice-hover');
  await active.mouse.move(5, 5);
  await expect(active.getByTestId('card-zoom')).toHaveCount(0);
  await cards.nth(1).click({ button: 'right' });
  await expect(active.getByTestId('card-zoom')).toBeVisible({ timeout: 300 });
  await active.mouse.move(5, 5);

  // View board shrinks the cards away; back, they grow in with no deal.
  await dialog.getByRole('button', { name: tr('choice.view_board') }).click();
  await active.waitForTimeout(80);
  await shot(active, '20.18-moon-hide-mid');
  const peek = active.getByRole('group', { name: tr('choice.viewing') });
  await expect(peek).toBeVisible();
  await shot(active, '20.18-moon-view-board');
  await peek.getByRole('button', { name: tr('choice.back') }).click();
  await expect(dialog).toBeVisible();
  await expect(active.locator('.choice.dealing')).toHaveCount(0);

  // The pick sounds; the cards fly back to the deck, the window closes.
  await cards.nth(2).click();
  await expect.poll(() => played(active, 'choice-pick')).toBe(1);
  await expect(active.getByTestId('choice-back')).toBeVisible();
  await active.waitForTimeout(150);
  await shot(active, '20.18-moon-picked-150ms');
  await expect(dialog).toHaveCount(0);
  await expect(active.getByTestId('choice-back')).toHaveCount(0, { timeout: 5000 });
  await expect(backs).toHaveCount(0, { timeout: 5000 });
  await expect(active.getByRole('dialog', { name: tr('look.title') })).toHaveCount(0);
  await shot(active, '20.18-moon-after');
});

// feedOrder shows the page an order prompt of these cards on the monster
// deck, as the server sends it to the chooser (19.16, 20.18). act is
// caught: window.sentOrder is the intent it would send.
async function feedOrder(page: Page, refs: string[]): Promise<void> {
  await page.evaluate((refs) => {
    interface S {
      game: { step: number; view: Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: object) => void } } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    game.view.ordering = refs.map((card, i) => ({ id: 901 + i, card, charged: false }));
    // 2 is the monster deck (protocol Deck).
    game.view.waiting = { kind: 5, player: you, purpose: 17, count: refs.length, text: { key: 'ask.put_back' }, look: { deck: 2, count: refs.length } };
    game.allowed = [{ player: you, kind: 8, objects: refs.map((_, i) => 901 + i) }];
    store.setState({
      act: (intent: unknown) => {
        (window as unknown as { sentOrder: unknown }).sentOrder = intent;
      },
    });
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
  }, refs);
}

test('an order shows its cards at a choice\'s size when they fit one row; speed off moves nothing (20.18)', async ({ browser }) => {
  const [ola] = await startGame(browser, names(['Ola', 'Pim']), { url: fast() });
  if (!ola) {
    throw new Error('no pages');
  }
  const dialog = ola.getByRole('dialog', { name: tr('ask.put_back') });
  const left = dialog.getByRole('group', { name: tr('order.not_placed') });
  await feedOrder(ola, ['fly', 'leech', 'pooter']);
  await expect(left.locator('.choice-card')).toHaveCount(3);
  // Speed off: no deal, no flight.
  await expect(ola.locator('.choice.dealing')).toHaveCount(0);
  // Script motion only: a CSS hover transition may run under the pointer.
  expect(await dialog.evaluate((el) => el.getAnimations({ subtree: true }).filter((a) => !(a instanceof CSSTransition) && !(a instanceof CSSAnimation)).length)).toBe(0);
  const width = (await left.locator('.card').first().boundingBox())?.width ?? 0;
  expect(width).toBeGreaterThan(100);
  for (const ref of ['pooter', 'fly', 'leech']) {
    await left.getByRole('button', { name: cardName(ref) }).click();
  }
  await expect(dialog.getByRole('group', { name: tr('order.place', { n: 1 }) }).getByRole('button', { name: cardName('pooter') })).toBeVisible();
  await dialog.getByRole('button', { name: tr('order.put_back') }).click();
  await expect(ola.getByTestId('choice-back')).toHaveCount(0);
  expect(await ola.evaluate(() => (window as unknown as { sentOrder: unknown }).sentOrder)).toMatchObject({ kind: 8, objects: [903, 901, 902] });
});

test('an order of more cards than one row holds keeps table size (20.18)', async ({ browser }) => {
  const [rex] = await startGame(browser, names(['Rex', 'Sia']), { url: fast() });
  if (!rex) {
    throw new Error('no pages');
  }
  await feedOrder(rex, ['fly', 'leech', 'pooter', 'dip', 'clotty', 'hopper', 'horf']);
  const left = rex.getByRole('dialog', { name: tr('ask.put_back') }).getByRole('group', { name: tr('order.not_placed') });
  await expect(left.locator('.choice-card')).toHaveCount(7);
  expect(Math.round((await left.locator('.card').first().boundingBox())?.width ?? 0)).toBe(80);
});

test('Sleight of Hand deals the monster deck into the order and puts it back (20.18)', async ({ browser }) => {
  const who = names(['Tam', 'Uri']);
  const pages = await startGame(browser, who, { url: appURL() });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'bomb', 'sleight_of_hand');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await watchDeal(active, 'monster');
  const item = active.getByRole('button', { name: cardName('sleight_of_hand'), exact: true }).first();
  await click(active, item);
  await click(active, active.getByRole('menu').getByRole('menuitem', { name: tr('actions.use') }));
  const deck = active.getByRole('button', { name: tr('ask.deck.monster') });
  await passUntil(pages, deck);
  await shot(active, '20.18-sleight-deck');
  await deck.click();
  const dialog = active.getByRole('dialog', { name: tr('ask.put_back') });
  await expect(dialog).toBeVisible();
  await active.waitForTimeout(150);
  await shot(active, '20.18-sleight-deal-mid');
  expect(await dealtFar(active)).toBeLessThan(40);
  await expect(active.locator('.choice.dealing')).toHaveCount(0, { timeout: 5000 });
  await shot(active, '20.18-sleight-order');
  const backs = other.locator('[data-zone="deck-monster"]').getByRole('img', { name: tr('deck.looking', { player: await nameOf(active), n: 5 }) });
  await expect(backs.locator('.deck-look-back')).toHaveCount(5);
  await other.waitForTimeout(800);
  await shot(other, '20.18-sleight-order-other');
  const left = dialog.getByRole('group', { name: tr('order.not_placed') });
  for (let i = 0; i < 5; i++) {
    await left.locator('.choice-card').first().click();
    if (i === 1) {
      await shot(active, '20.18-sleight-two-placed');
    }
  }
  await dialog.getByRole('button', { name: tr('order.put_back') }).click();
  await expect(active.getByTestId('choice-back')).toBeVisible();
  await active.waitForTimeout(200);
  await shot(active, '20.18-sleight-back-mid');
  await expect(active.getByTestId('choice-back')).toHaveCount(0, { timeout: 5000 });
  await expect(backs).toHaveCount(0, { timeout: 5000 });
});

// feedCards shows the page a choice of these cards with no deck look,
// as Eden's start or Red Key ask (20.18). act is caught.
async function feedCards(page: Page, refs: string[]): Promise<void> {
  await page.evaluate((refs) => {
    interface S {
      game: { step: number; view: Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const w = window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: object) => void } }; grown: number };
    const store = w.fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const options = refs.map((card) => ({ key: 'ask.card', args: { card: { card } } }));
    game.view.waiting = { kind: 4, player: you, purpose: 9, options, text: { key: 'ask.one_on_top' } };
    game.view.options = options.map(() => ({ player: -1 }));
    game.allowed = options.map((_, i) => ({ player: you, kind: 4, choice: i }));
    store.setState({ act: () => undefined });
    // The cards' motion as they show: noted when they come in.
    w.grown = -1;
    const seen = new MutationObserver(() => {
      const card = document.querySelector('.choice-card');
      if (card && w.grown < 0) {
        w.grown = card.getAnimations().length;
        seen.disconnect();
      }
    });
    seen.observe(document.body, { childList: true, subtree: true });
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
  }, refs);
}

for (const anim of ['on', 'off'] as const) {
  test(`a card choice with no deck grows in; speed ${anim} (20.18)`, async ({ browser }) => {
    const [vic] = await startGame(browser, names(['Vic', 'Wes']), { url: anim === 'on' ? appURL() : fast() });
    if (!vic) {
      throw new Error('no pages');
    }
    await feedCards(vic, ['a_dime', 'bomb', 'the_d6']);
    const dialog = vic.getByRole('dialog', { name: tr('ask.one_on_top') });
    await expect(dialog.locator('.choice-card')).toHaveCount(3);
    const grown = await vic.evaluate(() => (window as unknown as { grown: number }).grown);
    // On: each card grows from 0; never a deal. Off: they are there at once.
    expect(grown).toBe(anim === 'on' ? 1 : 0);
    await expect(vic.locator('.choice.dealing')).toHaveCount(0);
    await expect(dialog.getByRole('button', { name: cardName('bomb') })).toBeVisible();
  });
}
