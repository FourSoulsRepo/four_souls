// The replay player (10.2, ADR 031, UC-45): a fake record of the fake
// bridge plays on the table (fakeReplay.ts). A turn every ten steps; the
// active player gains a cent at step 4 of each turn (index 3); the
// finished record's last step is the game's end.
import type { Page } from '@playwright/test';
import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { constants, gunzipSync } from 'node:zlib';
import { expect, test } from './fixtures';
import { gameId, startGame } from './game';
import { appURL, serverPort } from './server';
import { shot } from './shot';
import { cardName, tr, trRe } from './locale';
import { activePage, click, endTurn, expectNextTurn } from './turn';

// watch opens the Records screen and watches the record of exactly
// these players.
async function watch(page: Page, players: string): Promise<void> {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.records'), exact: true }).click();
  await page
    .getByRole('list', { name: tr('records.list_label') })
    .getByRole('listitem')
    .filter({ hasText: new RegExp(`${players}(?!,)`) })
    .getByRole('button', { name: trRe('records.watch_label') })
    .click();
  await expect(bar(page)).toBeVisible();
}

function bar(page: Page) {
  return page.getByRole('region', { name: tr('replay.label') });
}

// where is the bar's "Step i / n · Turn t".
function where(page: Page) {
  return bar(page).getByRole('status');
}

// cents is a seat's cents counter showing n.
function cents(page: Page, player: string, n: number) {
  return page.getByRole('region', { name: player, exact: true }).getByLabel(tr('table.cents', { n }), { exact: true });
}

test('Watch puts the record on the table; forward plays, back shows at once', async ({ page }) => {
  await watch(page, 'Ann, Bo');
  await expect(where(page)).toHaveText(tr('replay.where', { i: 1, n: 12, turn: 1 }));
  // Nobody sits at a replay: no "(you)", no turn actions, no End turn, no menu.
  await expect(page.getByRole('region', { name: 'Ann', exact: true })).toBeVisible();
  await expect(page.getByRole('toolbar', { name: tr('turn.actions') })).toHaveCount(0);
  await expect(page.getByRole('button', { name: tr('turn.end') })).toHaveCount(0);
  await expect(page.getByRole('button', { name: tr('game_menu.button') })).toHaveCount(0);
  await expect(bar(page).getByRole('button', { name: tr('replay.back') })).toBeDisabled();
  await expect(cents(page, 'Ann', 3)).toBeVisible();

  const forward = bar(page).getByRole('button', { name: tr('replay.forward') });
  await forward.click();
  await forward.click();
  await expect(where(page)).toHaveText(tr('replay.where', { i: 3, n: 12, turn: 1 }));
  // Step 4 plays its effect: Ann's cent pops up, then the view shows it.
  await forward.click();
  await expect(page.getByText('+1¢')).toBeVisible();
  await expect(cents(page, 'Ann', 4)).toBeVisible();
  await shot(page, 'replay');

  // Back shows step 3 at once.
  await bar(page).getByRole('button', { name: tr('replay.back') }).click();
  await expect(where(page)).toHaveText(tr('replay.where', { i: 3, n: 12, turn: 1 }));
  await expect(cents(page, 'Ann', 3)).toBeVisible();
});

test('a turn and the slider jump; Esc goes back to the records', async ({ page }) => {
  await watch(page, 'Ann, Bo');
  await bar(page).getByRole('combobox', { name: tr('replay.turn_pick') }).selectOption({ label: tr('replay.turn', { n: 2 }) });
  await expect(where(page)).toHaveText(tr('replay.where', { i: 11, n: 12, turn: 2 }));
  await expect(page.getByRole('main', { name: tr('game.turn', { n: 2 }) })).toBeVisible();
  await expect(cents(page, 'Ann', 4)).toBeVisible();
  await expect(bar(page).getByRole('button', { name: tr('replay.forward') })).toBeEnabled();

  await bar(page).getByRole('slider', { name: tr('replay.slider') }).fill('5');
  await expect(where(page)).toHaveText(tr('replay.where', { i: 5, n: 12, turn: 1 }));
  await expect(bar(page).getByRole('combobox', { name: tr('replay.turn_pick') })).toHaveValue('0');

  await page.keyboard.press('Escape');
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toBeVisible();
});

test('autoplay at 4× plays the last turn to the end and stops', async ({ page }) => {
  await watch(page, 'Ann, Bo, Cy');
  await bar(page).getByRole('combobox', { name: tr('replay.turn_pick') }).selectOption({ label: tr('replay.turn', { n: 84 }) });
  await expect(where(page)).toHaveText(tr('replay.where', { i: 831, n: 840, turn: 84 }));
  await bar(page).getByRole('button', { name: tr('replay.speed_n', { n: 4 }) }).click();
  await expect(bar(page).getByRole('button', { name: tr('replay.speed_n', { n: 4 }) })).toHaveAttribute('aria-pressed', 'true');
  await bar(page).getByRole('button', { name: tr('replay.play') }).click();
  await expect(bar(page).getByRole('button', { name: tr('replay.pause') })).toBeVisible();
  await shot(page, 'replay-playing');
  // The player stops at the end, on the last step's view: the game is
  // over. Its short winner banner is not waited for (18.10).
  await expect(where(page)).toHaveText(tr('replay.where', { i: 840, n: 840, turn: 84 }), { timeout: 15_000 });
  await expect
    .poll(() =>
      page.evaluate(() => {
        interface S {
          game: { view: { over?: boolean; winners?: number[] } } | null;
        }
        const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
        return { over: v?.over, winners: v?.winners };
      }),
    )
    .toEqual({ over: true, winners: [1] });
  await expect(bar(page).getByRole('button', { name: tr('replay.play') })).toBeDisabled();
  await expect(bar(page).getByRole('button', { name: tr('replay.forward') })).toBeDisabled();
  await shot(page, 'replay-end');

  await bar(page).getByRole('button', { name: tr('replay.leave') }).click();
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toBeVisible();
});

test('pause stops the autoplay where it is', async ({ page }) => {
  await watch(page, 'Ann, Bo');
  await bar(page).getByRole('button', { name: tr('replay.play') }).click();
  // Step 4 has an effect and the longer pause after it.
  await expect(where(page)).toHaveText(tr('replay.where', { i: 4, n: 12, turn: 1 }));
  await bar(page).getByRole('button', { name: tr('replay.pause') }).click();
  const at = await where(page).textContent();
  await page.waitForTimeout(1500);
  await expect(where(page)).toHaveText(at ?? '');
});

test('the bar hides and shows; autoplay goes on meanwhile; Esc still leaves (10.7)', async ({ page }) => {
  await watch(page, 'Ann, Bo');
  const hide = bar(page).getByRole('button', { name: tr('replay.hide') });
  const show = bar(page).getByRole('button', { name: tr('replay.show') });
  await expect(show).toHaveCount(0);
  await bar(page).getByRole('button', { name: tr('replay.speed_n', { n: 2 }) }).click();
  await bar(page).getByRole('button', { name: tr('replay.play') }).click();
  await hide.click();
  // Only the small button stays, in the same corner.
  await expect(show).toBeVisible();
  await expect(where(page)).toHaveCount(0);
  await expect(bar(page).getByRole('button', { name: tr('replay.pause') })).toHaveCount(0);
  await shot(page, 'replay-bar-hidden');
  // The record plays on into turn 2 while the bar is hidden.
  await expect(page.getByRole('main', { name: tr('game.turn', { n: 2 }) })).toBeVisible({ timeout: 15_000 });
  await show.click();
  await expect(hide).toBeVisible();
  await expect(where(page)).toHaveText(trRe('replay.where', { turn: 2 }));
  await expect(where(page)).toHaveText(tr('replay.where', { i: 12, n: 12, turn: 2 }), { timeout: 15_000 });

  // Esc leaves the replay with the bar hidden too.
  await hide.click();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toBeVisible();
  // The next Watch shows the bar again.
  await watch(page, 'Ann, Bo');
  await expect(hide).toBeVisible();
});

test('Show hands hides every hand as backs and shows them again (10.3)', async ({ page }) => {
  await watch(page, 'Ann, Bo');
  const toggle = bar(page).getByRole('button', { name: tr('replay.hands') });
  const ann = page.getByRole('region', { name: 'Ann', exact: true });
  const bo = page.getByRole('region', { name: 'Bo', exact: true });
  // Hands show at first (ADR 031, 4.7): Ann's at the bottom, Bo's revealed.
  await expect(toggle).toHaveAttribute('aria-pressed', 'true');
  await expect(ann.getByRole('img', { name: cardName('a_penny'), exact: true })).toBeVisible();
  await expect(bo.getByRole('group', { name: tr('table.revealed_hand', { n: 3 }) })).toBeVisible();

  await toggle.click();
  await expect(toggle).toHaveAttribute('aria-pressed', 'false');
  await expect(ann.getByRole('img', { name: cardName('a_penny'), exact: true })).toHaveCount(0);
  await expect(ann.getByRole('group', { name: tr('table.your_hand') }).getByRole('img', { name: tr('deck.face_down.loot') })).toHaveCount(2);
  await expect(bo.getByRole('group', { name: tr('table.hand', { n: 3 }) })).toBeVisible();
  await expect(bo.getByRole('img', { name: cardName('bomb'), exact: true })).toHaveCount(0);
  // A step forward keeps them hidden.
  await bar(page).getByRole('button', { name: tr('replay.forward') }).click();
  await expect(where(page)).toHaveText(tr('replay.where', { i: 2, n: 12, turn: 1 }));
  await expect(ann.getByRole('group', { name: tr('table.your_hand') }).getByRole('img', { name: tr('deck.face_down.loot') })).toHaveCount(2);
  await shot(page, 'replay-hands-hidden');

  await toggle.click();
  await expect(toggle).toHaveAttribute('aria-pressed', 'true');
  await expect(ann.getByRole('img', { name: cardName('a_penny'), exact: true })).toBeVisible();
  await expect(bo.getByRole('group', { name: tr('table.revealed_hand', { n: 3 }) })).toBeVisible();
});

// futureCard is a card of a set this app does not know; only the
// record's header has its text (RP-10).
const futureCard = { ref: 'zz_future_card', name: 'Future Penny', type: 'Loot Card', effects: ['Gain 9¢. Then loot 1.'] };

// futureRecord is a two-step record: the future card lies in the shop,
// and Ann loots another copy of it in step 2.
function futureRecord(): { path: string; summary: object; steps: object[] } {
  const view = (looted: boolean) => {
    const card = (id: number, ref: string, more: object = {}) => ({ id, card: ref, charged: true, ...more });
    const hand = [card(11, 'a_penny'), ...(looted ? [card(12, futureCard.ref)] : [])];
    return {
      players: [
        { id: 0, character: card(1, 'isaac', { role: 1 }), hp: 2, atk: 1, cents: 3, souls: 0, in_play: [], hand_size: hand.length, hand, loot_plays: 1, attacks: 1, purchases: 1 },
        { id: 1, character: card(2, 'cain', { role: 1 }), hp: 2, atk: 1, cents: 3, souls: 0, in_play: [], hand_size: 1, hand: [card(21, 'bomb')], loot_plays: 1, attacks: 1, purchases: 1 },
      ],
      shop: [{ top: card(31, futureCard.ref) }, { top: card(32, 'breakfast') }],
      monsters: [{ top: card(41, 'big_spider', { stats: { hp: 2, dc: 3, atk: 1, printed_hp: 2, printed_dc: 3, printed_atk: 1 } }) }],
      deck_sizes: [80, 60, 90, 0],
      discards: [null, null, null, null],
      discard_sizes: [0, 0, 0, 0],
      stack: null,
      turn: { active: 0, number: 1, step: 4, loot_plays: 1, attacks: 1, purchases: 1 },
      waiting: { kind: 1, player: 0 },
      win_souls: 4,
    };
  };
  const path = '/fake/records/future.fsrec';
  return {
    path,
    summary: {
      path,
      game: 'future',
      server: '',
      started: Date.parse('2026-11-16T10:00:00Z'),
      seats: [
        { seat: 0, name: 'Ann' },
        { seat: 1, name: 'Bo' },
      ],
      sets: ['b2', 'zz'],
      challenge: '',
      difficulty: '',
      cards: [futureCard],
      steps: 2,
      turns: [{ number: 1, active: 0, step: 0 }],
    },
    steps: [
      { index: 0, n: 1, events: [], view: view(false) },
      { index: 1, n: 2, events: [{ kind: 'looted', player: 0, object: 12, card: futureCard.ref, private: true }], view: view(true) },
    ],
  };
}

test("a card the app lacks plays from the record's text (10.3)", async ({ page }) => {
  await page.goto('/?fake&nosplash');
  const record = futureRecord();
  await page.evaluate((r) => {
    const w = window as unknown as { fakeFiles: Record<string, unknown>; fakePick: string };
    w.fakeFiles = { [r.path]: { summary: r.summary, steps: r.steps } };
    w.fakePick = r.path;
  }, record);
  await page.getByRole('button', { name: tr('menu.records'), exact: true }).click();
  await page.getByRole('button', { name: tr('records.open_file') }).click();
  await expect(where(page)).toHaveText(tr('replay.where', { i: 1, n: 2, turn: 1 }));
  // The shop's card is a text card with the record's name and type.
  const shopCard = page.getByRole('group', { name: tr('board.shop') }).getByRole('img', { name: futureCard.name, exact: true });
  await expect(shopCard).toHaveCount(1);
  await expect(shopCard).toContainText(futureCard.type);
  // Ann loots one: the record plays on and her hand shows it.
  await bar(page).getByRole('button', { name: tr('replay.forward') }).click();
  await expect(where(page)).toHaveText(tr('replay.where', { i: 2, n: 2, turn: 1 }));
  const ann = page.getByRole('region', { name: 'Ann', exact: true });
  await expect(ann.getByRole('group', { name: tr('table.your_hand') }).getByRole('img', { name: futureCard.name, exact: true })).toBeVisible();
  // Its zoom shows the record's words.
  await shopCard.hover();
  await expect(page.getByTestId('card-zoom')).toContainText(futureCard.effects[0] ?? '');
  await shot(page, 'replay-unknown-card');

  // Back on the records list, the app forgets the record's cards.
  await bar(page).getByRole('button', { name: tr('replay.leave') }).click();
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toBeVisible();
  const known = await page.evaluate((ref) => {
    interface S {
      cards: Map<string, unknown>;
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().cards.has(ref);
  }, futureCard.ref);
  expect(known).toBe(false);
});

// Line is one line of a record file (pkg/record).
interface Line {
  kind: string;
  n?: number;
  to?: number;
  events?: unknown;
  view?: { turn: { number: number; active: number } };
  started?: string;
  game?: string;
  seats?: { seat: number; name: string }[];
  setup?: { sets: string[] };
}

// readRecord reads the test server's record of game id as Go's Open and
// Step give it: its live steps and a summary (bindings/records). The
// server flushes it every second, so the stream may be unfinished.
function readRecord(id: string): { path: string; summary: object; steps: object[] } | null {
  const dir = join(import.meta.dirname, '..', '..', '.cache', 'e2e', String(serverPort()), 'records');
  const name = readdirSync(dir).find((f) => f.includes(`-${id}_`));
  if (!name) {
    return null;
  }
  const path = join(dir, name);
  const text = gunzipSync(readFileSync(path), { finishFlush: constants.Z_SYNC_FLUSH }).toString('utf8');
  const lines = text
    .split('\n')
    .filter((l) => l.trim() !== '')
    .map((l) => JSON.parse(l) as Line);
  const header = lines[0];
  let live: Line[] = [];
  for (const l of lines) {
    if (l.kind === 'step') {
      live.push(l);
    } else if (l.kind === 'rollback') {
      live = live.filter((st) => (st.n ?? 0) < (l.to ?? 0));
    }
  }
  const steps = live.map((l, index) => ({ index, n: l.n, events: l.events ?? null, view: l.view }));
  const turns: { number: number; active: number; step: number }[] = [];
  live.forEach((l, step) => {
    const turn = l.view?.turn;
    if (turn && turn.number !== turns.at(-1)?.number) {
      turns.push({ number: turn.number, active: turn.active, step });
    }
  });
  const summary = {
    path,
    game: header?.game ?? '',
    server: '',
    started: Date.parse(header?.started ?? ''),
    seats: header?.seats ?? [],
    sets: header?.setup?.sets ?? [],
    challenge: '',
    difficulty: '',
    cards: [],
    steps: steps.length,
    turns,
  };
  return { path, summary, steps };
}

test('a real record from the test server plays to its end', async ({ browser, page }) => {
  test.setTimeout(120_000);
  // A short real game: the first-player roll and one whole turn.
  const players = await startGame(browser, ['Rea', 'Lee'], { url: `${appURL()}&anim=off` });
  const first = players[0];
  if (!first) {
    throw new Error('no pages');
  }
  const id = await gameId(first);
  await endTurn(await activePage(players, click));
  await expectNextTurn(players, click, 1);
  const step = await first.evaluate(() => {
    interface S {
      game: { step: number } | null;
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.step ?? 0;
  });
  // Wait for the server's flush of every step played so far.
  await expect.poll(() => readRecord(id)?.steps.at(-1), { timeout: 10_000 }).toMatchObject({ n: step });
  const record = readRecord(id);
  if (!record) {
    throw new Error('no record');
  }
  const n = record.steps.length;

  // The record plays through the fake bridge, as Go would serve it.
  await page.goto('/?fake&nosplash');
  await page.evaluate((r) => {
    const w = window as unknown as { fakeFiles: Record<string, unknown>; fakePick: string };
    w.fakeFiles = { [r.path]: { summary: r.summary, steps: r.steps } };
    w.fakePick = r.path;
  }, record);
  await page.getByRole('button', { name: tr('menu.records'), exact: true }).click();
  await page.getByRole('button', { name: tr('records.open_file') }).click();
  await expect(where(page)).toHaveText(trRe('replay.where', { i: 1, n }));
  await expect(page.getByRole('region', { name: 'Rea', exact: true })).toBeVisible();
  await expect(page.getByRole('region', { name: 'Lee', exact: true })).toBeVisible();
  await bar(page).getByRole('button', { name: tr('replay.speed_n', { n: 4 }) }).click();
  await bar(page).getByRole('button', { name: tr('replay.play') }).click();
  await expect(where(page)).toHaveText(trRe('replay.where', { i: n, n }), { timeout: 60_000 });
  await expect(bar(page).getByRole('button', { name: tr('replay.play') })).toBeDisabled();
  await expect(page.getByRole('main', { name: tr('game.turn', { n: 2 }) })).toBeVisible();
});
