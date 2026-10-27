// Stack and responses (7.8): players always see why the game waits.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { saveAndContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import {
  activePage,
  answerPlay,
  canEndTurn,
  click,
  endTurn,
  gameOver,
  passButton,
  playLoot,
  seen,
  settle,
  throwButton,
  turnAction,
  turnNumber,
} from './turn';
import { anyOf as any, cardName, lineRe, tr, trRe, whole } from './locale';

// choosing matches a choice's title: the app's, or the server's question
// (in English until 16.4).
const choosing = any(
  /Choose/,
  ...['replacement', 'attack', 'monster_slot', 'penalty_item', 'penalty_loot', 'purchase', 'target', 'trigger_order', 'mode', 'answer', 'starting_item', 'curse', 'bottom', 'choose'].map((k) =>
    trRe(`choice.title.${k}`),
  ),
);

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

test('a played loot card shows on the stack; nobody is offered skip all', async ({ browser }) => {
  const pages = await startGame(browser, ['Sam', 'Tia'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  const name = await loot.getAttribute('aria-label');
  await playLoot(active, loot, click);
  // Some loot cards ask for a mode or a target before they go on the stack.
  await answerPlay(active, click);

  for (const page of [active, other]) {
    await expect(page.getByRole('complementary', { name: tr('stack.title') })).toContainText(trRe('stack.loot', { card: name ?? '' }));
  }
  // Both always see why the game waits.
  await expect(active.locator('.prompt-title')).toHaveText(any(trRe('banner.respond'), choosing, trRe('banner.waiting_for')));
  await expect(other.locator('.prompt-title')).toHaveText(trRe('banner.waiting_for'));
  await shot(other, 'stack');

  // "Skip all" is gone (15.84); "Wait for my turn" stays.
  for (const page of [active, other]) {
    await expect(page.getByRole('button', { name: /Skip all|Cancel skip/ })).toHaveCount(0);
  }
});

// A loot card that aims after its mode still reaches the stack: the
// test picks its player on the table, not in a dialog (19.18).
test('a played Dagaz goes on the stack once its player is picked (19.18)', async ({ browser }) => {
  const pages = await startGame(browser, ['Kit', 'Lev'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Kit, Lev', 'dagaz');
  const active = await activePage(pages, click);
  await playLoot(active, active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('dagaz') }).first(), click);
  await answerPlay(active, click);
  await expect(active.getByRole('complementary', { name: tr('stack.title') })).toContainText(trRe('stack.loot', { card: cardName('dagaz') }));
});

test('a player may wait for their turn; the end of the turn before asks them (N-16)', async ({ browser }) => {
  const pages = await startGame(browser, ['Wan', 'Oti'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await shot(other, 'wait-turn-offer');
  const bar = other.getByRole('region', { name: tr('actions.label') });
  await bar.getByRole('button', { name: tr('actions.wait') }).click();
  await expect(bar.getByRole('button', { name: tr('actions.cancel_wait') })).toBeVisible();
  await shot(other, 'wait-turn-on');
  await bar.getByRole('button', { name: tr('actions.cancel_wait') }).click();
  await expect(bar.getByRole('button', { name: tr('actions.wait') })).toBeVisible();
  await bar.getByRole('button', { name: tr('actions.wait') }).click();
  await expect(bar.getByRole('button', { name: tr('actions.cancel_wait') })).toBeVisible();

  // The turn before ends: waiting is over, and nothing more is offered.
  await endTurn(active);
  await expect(bar.getByRole('button', { name: any(trRe('actions.cancel_wait'), trRe('actions.wait')) })).toHaveCount(0);
  if (await passButton(other).isEnabled()) {
    await expect(other.locator('.prompt-title')).toHaveText(any(trRe('banner.act'), trRe('banner.respond')));
    await shot(other, 'wait-turn-last-chance');
  }
});

test('the response timer counts down for the waiting player', async ({ browser }) => {
  const pages = await startGame(browser, ['Uma', 'Vic'], {
    url: fast(),
    setup: async (host) => {
      await host.getByLabel(tr('host.detailed')).check();
      await host.getByLabel(tr('host.timer')).selectOption('60');
    },
  });
  const active = await activePage(pages, click);
  const timer = active.getByLabel(tr('actions.time_left'));
  await expect(timer).toHaveText(new RegExp(`· ${trRe('actions.seconds').source}`));
  const first = Number((await timer.textContent())?.match(/\d+/)?.[0]);
  expect(first).toBeGreaterThan(0);
  await expect.poll(async () => Number((await timer.textContent())?.match(/\d+/)?.[0])).toBeLessThan(first);
});

// waitLine is a stack panel line with "You" turned into the page's own
// name, so every page's line can be compared.
async function waitLine(page: Page, me: string, label: string): Promise<string> {
  const line = page.getByRole('complementary', { name: tr('stack.title') }).getByLabel(label === 'Passed' ? tr('stack.passed_label') : tr('stack.waiting_label'), { exact: true });
  if ((await line.count()) === 0) {
    return '';
  }
  return (await line.first().innerText()).replaceAll(tr('stack.you'), me);
}

// expectSameLines: every page names the same waiting and passed players.
async function expectSameLines(pages: Page[], names: string[], want: string): Promise<void> {
  for (const [i, page] of pages.entries()) {
    const me = names[i] ?? '';
    await expect.poll(async () => (await waitLine(page, me, 'Waiting for')) + (await waitLine(page, me, 'Passed'))).toBe(want);
  }
}

test('every player sees whom the stack waits for (V-11)', async ({ browser }) => {
  const names = ['Abe', 'Bea', 'Cal', 'Dot'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const me = names[pages.indexOf(active)] ?? '';
  const loot = active.locator('.hand .card.hl-act').first();
  test.skip(!(await loot.isVisible()), 'no playable loot card in this deal');
  await playLoot(active, loot, click);
  const stack = active.getByRole('complementary', { name: tr('stack.title') });
  await answerPlay(active, click);
  await expect(stack).toContainText(any(trRe('stack.waiting'), /: /));
  const want = (await waitLine(active, me, 'Waiting for')) + (await waitLine(active, me, 'Passed'));
  expect(want).not.toBe('');
  await expectSameLines(pages, names, want);
  await shot(active, 'stack-waiting');

  // A pass moves the player to "Passed" on every table.
  const pass = passButton(active);
  test.skip(!(await pass.isEnabled()), 'the loot card asked something first');
  await pass.click();
  // Unless everyone else was skipped and the item resolved.
  await expect
    .poll(async () => (await stack.locator('.stack-row').count()) === 0 || (await waitLine(active, me, 'Passed')).includes(me))
    .toBe(true);
  await expectSameLines(pages, names, (await waitLine(active, me, 'Waiting for')) + (await waitLine(active, me, 'Passed')));
});

// fakeStack shows stack items as the server would send them: a loot
// card (1), a roll (2), and damage (3), the ones in ids. Only the
// screen changes; nobody acts on it.
async function fakeStack(page: Page, ids: number[]): Promise<void> {
  await page.evaluate((ids) => {
    interface Item {
      id: number;
      kind: number;
      controller: number;
      source: number;
      card?: string;
      roll?: number;
      amount?: number;
      attack?: boolean;
      text?: string;
      target: { player: number; is_player?: boolean };
      ability: { card: string; index: number };
      roll_for: { card: string; index: number };
    }
    interface S {
      game: { view: { stack: Item[] | null; turn: { active: number } } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const game = store.getState().game;
    if (!game) {
      throw new Error('no game');
    }
    const me = game.view.turn.active;
    const other = 1 - me;
    const none = { card: '', index: 0 };
    const stack: Item[] = [
      {
        id: 1,
        kind: 0,
        controller: me,
        source: 0,
        card: 'bomb',
        text: 'Deal 1 damage to a monster or player.',
        target: { player: other, is_player: true },
        ability: { card: 'bomb', index: 0 },
        roll_for: none,
      },
      { id: 2, kind: 3, controller: other, source: 0, roll: 5, target: { player: -1 }, ability: none, roll_for: none },
      { id: 3, kind: 4, controller: me, source: 0, amount: 1, target: { player: other, is_player: true }, ability: none, roll_for: none },
    ];
    store.setState({ game: { ...game, view: { ...game.view, stack: stack.filter((it) => ids.includes(it.id)) } } });
  }, ids);
}

test('the stack panel: icons, colored names, and a hover with card and owner (15.21)', async ({ browser }) => {
  const pages = await startGame(browser, ['Ike', 'Jud'], { url: fast() });
  const active = await activePage(pages, click);
  const stack = active.getByRole('complementary', { name: tr('stack.title') });
  // An empty stack says so, and Pass stays in the panel (S-08, S-09).
  await expect(stack).toContainText(tr('stack.empty'));
  await expect(passButton(active)).toBeVisible();

  // Three items: a loot card, a roll, and damage.
  await fakeStack(active, [1, 2, 3]);

  const rows = stack.locator('.stack-row');
  await expect(rows).toHaveCount(3);
  // Newest on top; the target in red, the player in their color (S-01, S-03).
  await expect(rows.nth(0)).toHaveAttribute('aria-label', whole(trRe('stack.damage', { n: 1 })));
  await expect(rows.nth(0).locator('.stack-target')).toBeVisible();
  await expect(rows.nth(1).locator('.die')).toHaveAttribute('aria-label', tr('dice.rolled', { n: 5 }));
  await expect(rows.nth(2)).toHaveAttribute('aria-label', whole(trRe('stack.loot', { card: cardName('bomb') })));

  // Hover the loot row: its card large, the used ability, the owner (S-06, S-07).
  const owner = (await rows.nth(2).locator('.stack-player').textContent()) ?? '';
  await rows.nth(2).hover();
  const zoom = active.getByTestId('card-zoom');
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', new RegExp(`^${lineRe(cardName('bomb')).source}`));
  await expect(zoom.locator('.card-zoom-ability')).toHaveText('Deal 1 damage to a monster or player.');
  await expect(zoom.locator('.card-zoom-owner')).toHaveText(owner);
  await shot(active, 'stack-panel-hover');
  await active.mouse.move(1200, 420);
  await expect(zoom).toBeHidden();
  await shot(active, 'stack-panel');
});

// An open attack or buy is no stack item (R-ATK-02, R-SHOP-02); the
// stack panel shows it on top for every player while it lasts (15.91).
const declared = (page: Page) => page.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-declared');

test('a declared attack shows on top of the stack, then its target (15.91)', async ({ browser }) => {
  const names = ['Ada', 'Bex'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const who = names[pages.indexOf(active)] ?? '';
  await click(active, turnAction(active, 'Attack'));
  // Until the target is picked, the attack has none.
  for (const page of pages) {
    await expect(declared(page)).toHaveText(tr('stack.declared.attack', { player: who }));
  }
  await shot(pages.find((p) => p !== active) ?? active, 'stack-declared-attack');
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
  const monster = await active.locator('.board-monsters .card.fight').getAttribute('aria-label');
  for (const page of pages) {
    // The declaration's own line; the fight line under it is 15.94's.
    await expect(declared(page).locator('.stack-line')).toHaveText(tr('stack.declared.attack_target', { player: who, target: monster ?? '' }));
    await expect(declared(page).locator('.stack-target')).toHaveText(monster ?? '');
    await expect(declared(page).locator('.stack-odds')).toHaveText(any(...['stats.hits_from', 'stats.hits_on', 'stats.cant_hit', 'stats.always_hits'].map((k) => trRe(k))));
  }
  await shot(active, 'stack-declared-attack-aimed');
});

test('a declared buy shows on top of the stack until it ends (15.91)', async ({ browser }) => {
  const names = ['Cid', 'Dot'];
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const who = names[pages.indexOf(active)] ?? '';
  const turn = await turnNumber(active);
  await active.keyboard.press('b');
  for (const page of pages) {
    await expect(declared(page)).toHaveText(tr('stack.declared.buy', { player: who }));
  }
  await shot(pages.find((p) => p !== active) ?? active, 'stack-declared-buy');
  await settle(pages, click, async () => (await canEndTurn(active)) || (await turnNumber(active)) !== turn || (await gameOver(active)));
  for (const page of pages) {
    await expect(declared(page)).toHaveCount(0);
  }
});

// endOpen: the page's view has the turn's end declared, still in the
// action step (R-TURN-08), so the others may respond.
async function endOpen(page: Page): Promise<boolean> {
  return page.evaluate(() => {
    interface S {
      game: { view: { turn: { step: number; end_declared?: boolean } } } | null;
    }
    const turn = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view.turn;
    return (turn?.end_declared ?? false) && turn?.step === 4;
  });
}

test('a declared end of turn shows on top of the stack (19.74)', async ({ browser }) => {
  const names = ['Eli', 'Fay'];
  const pages = await startGame(browser, names, { url: fast() });
  // The server passes for a player who cannot respond (N-05); then the
  // end window closes at once. A later turn's deal may let them respond.
  for (let tries = 0; tries < 6; tries++) {
    const active = await activePage(pages, click);
    const other = pages.find((p) => p !== active);
    if (!other) {
      throw new Error('no other page');
    }
    const who = names[pages.indexOf(active)] ?? '';
    const before = (await seen(other)).step;
    await endTurn(active, click);
    await expect.poll(async () => (await seen(other)).step).toBeGreaterThan(before);
    if (!(await endOpen(other))) {
      continue;
    }
    await expect(passButton(other)).toBeEnabled();
    await other.mouse.move(1200, 420);
    await shot(other, 'stack-declared-end');
    for (const page of pages) {
      await expect(declared(page)).toHaveText(tr('stack.declared.end_turn', { player: who }));
      await expect(declared(page).locator('.stack-player')).toHaveText(who);
    }
    // The row goes once the end phase starts (R-TURN-10).
    await passButton(other).click();
    for (const page of pages) {
      await expect(declared(page)).toHaveCount(0);
    }
    return;
  }
  test.skip(true, 'nobody could respond to an end of turn in six turns');
});

// Frame is one animation frame after the top stack item resolved.
interface Frame {
  // there: the resolved row is still in the list; opacity is its own.
  there: boolean;
  opacity: number;
  // next is the top of the row below it.
  next: number;
}

// resolveTop takes the top item (3) off the stack and samples the
// rows on every frame for a second, from the next one (15.51). top and next are where the two top rows
// were before.
async function resolveTop(page: Page): Promise<{ top: number; next: number; frames: Frame[] }> {
  return page.evaluate(async () => {
    interface S {
      game: { view: { stack: { id: number }[] | null } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const rows = [...document.querySelectorAll('.stack ol > li')];
    const [gone, below] = rows;
    const game = store.getState().game;
    if (!gone || !below || !game) {
      throw new Error('no rows to resolve');
    }
    const top = gone.getBoundingClientRect().top;
    const next = below.getBoundingClientRect().top;
    store.setState({ game: { ...game, view: { ...game.view, stack: (game.view.stack ?? []).filter((it) => it.id !== 3) } } });
    const frames: Frame[] = [];
    const start = performance.now();
    while (performance.now() - start < 1000) {
      await new Promise((r) => requestAnimationFrame(r));
      frames.push({ there: gone.isConnected, opacity: Number(getComputedStyle(gone).opacity), next: below.getBoundingClientRect().top });
    }
    return { top, next, frames };
  });
}

// quiet waits until the game's own effects are over, so its next view
// does not replace the faked stack.
async function quiet(page: Page): Promise<void> {
  await expect
    .poll(() =>
      page.evaluate(() => {
        const s = (window as unknown as { fourSouls: { store: { getState: () => { behind: boolean; playing: unknown } } } }).fourSouls.store.getState();
        return !s.behind && s.playing === null;
      }),
    )
    .toBe(true);
}

test('a resolved stack row fades out, then the next slides up (15.51)', async ({ browser }) => {
  const pages = await startGame(browser, ['Kai', 'Lou'], { url: appURL() });
  const active = await activePage(pages, click);
  await quiet(active);
  await fakeStack(active, [1, 2, 3]);
  const rows = active.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-row');
  await expect(rows).toHaveCount(3);
  // The new rows have faded in.
  await expect.poll(() => rows.evaluateAll((els) => els.every((el) => getComputedStyle(el).opacity === '1'))).toBe(true);

  await shot(active, 'stack-rows-before-resolve');
  const sampled = resolveTop(active);
  // Mid-fade, then mid-slide; screenshots only.
  if (process.env.SCREENSHOTS) {
    await active.waitForTimeout(100);
    await shot(active, 'stack-row-fading');
    await active.waitForTimeout(200);
    await shot(active, 'stack-rows-sliding');
  }
  const { top, next, frames } = await sampled;
  await shot(active, 'stack-rows-after-resolve');
  // It fades out while it still holds its place.
  expect(frames.some((f) => f.there && f.opacity > 0 && f.opacity < 1 && f.next === next)).toBe(true);
  // Then it is gone, and the row below slides up into its place, not at once.
  expect(frames.at(-1)?.there).toBe(false);
  expect(frames.some((f) => !f.there && f.next < next - 1 && f.next > top + 1)).toBe(true);
  expect(Math.abs((frames.at(-1)?.next ?? 0) - top)).toBeLessThan(1);
  await expect(rows).toHaveCount(2);
});

test('with animations off, stack rows change at once (15.51)', async ({ browser }) => {
  const pages = await startGame(browser, ['Mia', 'Ned'], { url: fast() });
  const active = await activePage(pages, click);
  await fakeStack(active, [1, 2, 3]);
  await expect(active.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-row')).toHaveCount(3);

  const { top, frames } = await resolveTop(active);
  // From the first frame: the row is gone and the next one is in its place.
  expect(frames.every((f) => !f.there && Math.abs(f.next - top) < 1)).toBe(true);
});
