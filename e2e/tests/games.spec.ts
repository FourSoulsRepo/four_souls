// Whole games with every card set, played by random moves through the
// UI, for 2, 3 and 4 players (owner request, 2026-10-15). A game must
// end with a winner, with no server error, no page error and no console
// error or warning on any page. Each test logs how long its game took.
//
// Slow, so opt-in: FULL_GAMES=1, or `make e2e-games`.
// GAME_SEED repeats the moves' choices; the server's deal stays random.
// GAME_TURNS caps a game (default 300 turns).
import { copyFileSync, mkdirSync, readdirSync, statSync } from 'node:fs';
import { join } from 'node:path';
import type { CDPSession, Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { serverPort } from './server';
import { enabled, endTurnButton, endTurnQuestion, passButton, reachable, seen, throwButton, turnAction, type Seen } from './turn';
import { tr, trRe } from './locale';

const maxTurns = Number(process.env.GAME_TURNS ?? 300);
// stallMs: no step for this long fails the game.
const stallMs = 60_000;

// mulberry32 is a small seeded random source, so a failing run's choices
// can be repeated with its seed.
function mulberry32(seed: number): () => number {
  let a = seed >>> 0;
  return () => {
    a = (a + 0x6d2b79f5) >>> 0;
    let t = a;
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

type Rnd = () => number;

function one<T>(rnd: Rnd, list: T[]): T | undefined {
  return list[Math.floor(rnd() * list.length)];
}

// fakeFont: the fake bridge names the fonts at their real paths, which
// only screenshots serve (fakeBridge.ts); elsewhere Vite answers with
// its page, and Chromium warns that the font does not decode.
const fakeFont = /^(Failed to decode downloaded font: \S+\/cardimg\/assets\/fonts\/|OTS parsing error: invalid sfntVersion)/;

const reloaded = 'page reloaded';

// watch collects what a game must never show: console errors and
// warnings, uncaught page errors, error messages from the server, and a
// socket the server closes. It runs before the page loads.
function watch(page: Page, who: string, problems: string[]): void {
  page.on('console', (m) => {
    if (fakeFont.test(m.text())) {
      return;
    }
    if (m.type() === 'error' || m.type() === 'warning') {
      problems.push(`${who}: console ${m.type()}: ${m.text()}`);
    }
  });
  page.on('pageerror', (e) => {
    problems.push(`${who}: page error: ${e.message}`);
  });
  // The first load is the app; another one drops the game, e.g. Vite
  // reloading every page when a Wails build rewrites frontend/wailsjs.
  let loads = 0;
  page.on('load', () => {
    if (++loads > 1) {
      problems.push(`${who}: ${reloaded}`);
    }
  });
  page.on('websocket', (ws) => {
    ws.on('framereceived', (f) => {
      const text = String(f.payload);
      if (!text.includes('"error"')) {
        return;
      }
      const env = JSON.parse(text) as { type: string; data?: unknown };
      if (env.type === 'error') {
        problems.push(`${who}: server error: ${JSON.stringify(env.data)}`);
      }
    });
    ws.on('close', () => {
      if (!page.isClosed()) {
        problems.push(`${who}: socket closed`);
      }
    });
  });
}

// everySet checks every card set on the host form: startGame leaves
// only the Base Game.
async function everySet(host: Page): Promise<void> {
  const sets = host.getByRole('group', { name: tr('host.sets') }).getByRole('checkbox');
  for (const set of await sets.all()) {
    if (await set.isEnabled()) {
      await set.check();
    }
  }
}

// Weights of the free moves: attacking moves a game toward its end, so
// it is the likeliest; a player mostly lets others' stack resolve.
const weights = { attack: 6, buy: 3, card: 2, endTurn: 1, passOwn: 1, passOther: 6 };

// move makes one random move for a page that must act. It returns what
// it did, or null when nothing on the page could be pressed.
async function move(page: Page, s: Seen, rnd: Rnd): Promise<string | null> {
  const roll = throwButton(page);
  if (await enabled(roll)) {
    await roll.click();
    return 'throw';
  }
  const discard = page.getByRole('dialog', { name: tr('discard.label'), exact: true });
  if (await discard.isVisible()) {
    const done = discard.getByRole('button', { name: new RegExp(`^${trRe('discard.button').source}`) });
    for (let i = 0; i < 12 && !(await enabled(done)); i++) {
      const card = one(rnd, await discard.locator('.card.hl-act').all());
      if (!card) {
        break;
      }
      await card.click();
    }
    if (!(await enabled(done))) {
      return null;
    }
    await done.click();
    return 'discard';
  }
  const look = page.locator('.choice.look');
  if (await look.isVisible()) {
    await look.getByRole('button', { name: tr('look.ok') }).click();
    return 'look';
  }
  const order = page.locator('.choice.order');
  if (await order.isVisible()) {
    const done = order.getByRole('button', { name: tr('order.put_back') });
    for (let i = 0; i < 12 && !(await enabled(done)); i++) {
      const card = one(rnd, await order.locator('.card.hl-pick').all());
      if (!card) {
        break;
      }
      await card.click();
    }
    if (!(await enabled(done))) {
      return null;
    }
    await done.click();
    return 'order';
  }
  const menu = page.getByRole('menu');
  if (await menu.isVisible()) {
    const item = one(rnd, await menu.getByRole('menuitem').filter({ hasNotText: 'Cancel' }).all());
    await (item ?? menu.getByRole('menuitem', { name: tr('actions.cancel') })).click();
    return 'menu';
  }
  // End turn asks while plays are left (V-46): mostly yes, sometimes no.
  const question = endTurnQuestion(page);
  if (await question.isVisible()) {
    await question.getByRole('button', { name: rnd() < 0.75 ? 'End turn' : 'Keep playing' }).click();
    return 'end turn question';
  }
  const ability = page.locator('.ability-dialog');
  if (await ability.isVisible()) {
    const line = one(rnd, await ability.locator('.ability-line:enabled').all());
    await (line ?? ability.locator('.ability-cancel')).click();
    return 'ability';
  }
  const choice = page.locator('.choice:not(.look)');
  if (await choice.isVisible()) {
    const picks = [
      ...(await reachable(page.locator('.board .card.hl-pick, .player-pick'))),
      ...(await choice.locator('.choice-option button').all()),
    ];
    const pick = one(rnd, picks);
    if (!pick) {
      return null;
    }
    await pick.click();
    return 'choice';
  }
  const aim = one(rnd, await reachable(page.locator('.card.hl-aim, .stack-row.hl-aim')));
  if (aim) {
    await aim.click();
    return 'aim';
  }
  const own = s.active === s.you;
  const free: [string, number, Locator][] = [];
  const add = async (name: string, weight: number, l: Locator) => {
    if (await enabled(l)) {
      free.push([name, weight, l]);
    }
  };
  await add('attack', weights.attack, turnAction(page, 'Attack'));
  await add('buy', weights.buy, turnAction(page, 'Buy'));
  await add('end turn', weights.endTurn, endTurnButton(page));
  await add('pass', own ? weights.passOwn : weights.passOther, passButton(page));
  const card = one(rnd, await reachable(page.locator('.card.hl-act')));
  if (card) {
    free.push(['card', weights.card, card]);
  }
  let r = rnd() * free.reduce((sum, [, w]) => sum + w, 0);
  for (const [name, w, l] of free) {
    r -= w;
    if (r < 0) {
      await l.click();
      return name;
    }
  }
  return null;
}

// progressEvery: a progress line every this many turns, so a slow
// game shows whether it moves toward a winner.
const progressEvery = 10;

// souls reads each player's souls and the soul value that wins.
async function souls(page: Page): Promise<string> {
  return page.evaluate(() => {
    interface S {
      game: { view: { players: { souls: number }[]; win_souls?: number } } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    return `${(v?.players ?? []).map((p) => String(p.souls)).join('/')} of ${String(v?.win_souls ?? '?')}`;
  });
}

interface Result {
  seconds: number;
  turns: number;
  steps: number;
  moves: number;
}

// roundMs bounds one round of the loop: a page that stops taking input
// would block a key press or a click forever (19.12).
const roundMs = 30_000;

// bounded rejects when p takes longer than ms.
async function bounded<T>(p: Promise<T>, ms: number, what: string): Promise<T> {
  let timer: ReturnType<typeof setTimeout> | undefined;
  const late = new Promise<never>((_, reject) => {
    timer = setTimeout(() => {
      reject(new Error(`${what}: no answer in ${String(ms)} ms`));
    }, ms);
  });
  try {
    return await Promise.race([p, late]);
  } finally {
    clearTimeout(timer);
  }
}

// debuggers are each page's DevTools session, opened with the page: a
// page stuck in a script still pauses, but no longer enables one.
const debuggers = new WeakMap<Page, CDPSession>();

async function attachDebugger(page: Page): Promise<void> {
  const cdp = await page.context().newCDPSession(page);
  await cdp.send('Debugger.enable');
  debuggers.set(page, cdp);
}

// frozenStack pauses a page that does not answer and returns where its
// script runs.
async function frozenStack(page: Page): Promise<string[]> {
  const cdp = debuggers.get(page);
  if (!cdp) {
    return ['no debugger'];
  }
  const paused = new Promise<string[]>((resolve) => {
    cdp.on('Debugger.paused', (e) => {
      resolve(e.callFrames.slice(0, 25).map((f) => `${f.functionName || '(anonymous)'} ${f.url}:${String(f.location.lineNumber + 1)}:${String((f.location.columnNumber ?? 0) + 1)}`));
    });
  });
  await bounded(cdp.send('Debugger.pause'), 10_000, 'pause');
  return bounded(paused, 10_000, 'paused');
}

// diagnose saves what each page shows when a game fails: a page that
// answers gives its store and a screenshot; one that does not, the
// script it is stuck in.
async function diagnose(pages: Page[], names: string[]): Promise<void> {
  for (const [i, page] of pages.entries()) {
    const name = names[i] ?? String(i);
    const out: Record<string, unknown> = { name };
    const alive = await bounded(page.evaluate(() => true), 3_000, 'evaluate').catch(() => false);
    out.alive = alive;
    if (alive) {
      out.store = await page
        .evaluate(() => {
          interface S {
            game: { step: number; view: { waiting: unknown; turn: unknown; stack: unknown }; allowed: unknown; aims?: unknown; events: unknown } | null;
            behind: boolean;
            sent: number | null;
            selected: unknown;
            aiming: unknown;
            error: string;
          }
          const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
          const g = s.game;
          return {
            step: g?.step,
            sent: s.sent,
            behind: s.behind,
            selected: s.selected,
            aiming: s.aiming,
            error: s.error,
            waiting: g?.view.waiting,
            turn: g?.view.turn,
            stack: g?.view.stack,
            allowed: g?.allowed,
            aims: g?.aims,
            events: g?.events,
          };
        })
        .catch((e: unknown) => String(e));
      await bounded(page.screenshot({ path: test.info().outputPath(`${name}.png`) }), 10_000, 'screenshot').catch(() => undefined);
    } else {
      out.stack = await frozenStack(page).catch((e: unknown) => [String(e)]);
    }
    const text = JSON.stringify(out, null, 2);
    await test.info().attach(`${name}.json`, { body: text, contentType: 'application/json' });
    console.log(`diagnosis ${name}: alive ${String(alive)}${alive ? '' : `\n${(out.stack as string[]).join('\n')}`}`);
  }
}

// play makes random moves until the game ends.
async function play(pages: Page[], names: string[], problems: string[], rnd: Rnd): Promise<Result> {
  const start = Date.now();
  let moves = 0;
  let last = { step: -1, at: Date.now() };
  let reported = 0;
  // round looks once at every page and makes at most one move. It
  // returns the result when the game is over.
  const round = async (): Promise<Result | null> => {
    const gone = problems.filter((p) => p.endsWith(reloaded));
    if (gone.length > 0) {
      throw new Error(`the game is gone: ${gone.join(', ')}`);
    }
    const states = await Promise.all(pages.map(seen));
    const first = states[0];
    if (!first) {
      throw new Error('no pages');
    }
    if (first.over) {
      return { seconds: (Date.now() - start) / 1000, turns: first.turn, steps: first.step, moves };
    }
    if (first.turn >= reported + progressEvery) {
      reported = first.turn - (first.turn % progressEvery);
      const seconds = ((Date.now() - start) / 1000).toFixed(0);
      console.log(`  ${names[0] ?? ''}: turn ${String(first.turn)}, step ${String(first.step)}, ${seconds} s, souls ${await souls(pages[0] as Page)}`);
    }
    if (first.turn > maxTurns) {
      throw new Error(`no winner after ${String(maxTurns)} turns`);
    }
    if (first.step !== last.step) {
      last = { step: first.step, at: Date.now() };
    } else if (Date.now() - last.at > stallMs) {
      throw new Error(`the game stalled at step ${String(first.step)}: ${JSON.stringify(states)}`);
    }
    let acted = false;
    for (const [i, page] of pages.entries()) {
      const s = states[i];
      if (!s?.ready || s.waiting !== s.you) {
        continue;
      }
      // A click can miss when the state moves on under it; the next
      // round looks again.
      const did = await move(page, s, rnd).catch(() => null);
      if (did) {
        acted = true;
        moves++;
      } else {
        // An aim or a menu with nothing to press: back out of it.
        await page.keyboard.press('Escape');
      }
      break;
    }
    if (!acted) {
      await pages[0]?.waitForTimeout(50);
    }
    return null;
  };
  for (;;) {
    try {
      const r = await bounded(round(), roundMs, 'a round of moves');
      if (r) {
        return r;
      }
    } catch (e) {
      await diagnose(pages, names);
      throw e;
    }
  }
}

// keepRecord copies a failed game's match record out of the test
// server's folder, which the next run deletes (playwright.config.ts),
// into .cache/e2e-games/. The record holds the deal and every move.
function keepRecord(names: string[]): string {
  const root = join(import.meta.dirname, '..', '..', '.cache');
  const dir = join(root, 'e2e', String(serverPort()), 'records');
  const end = `_${names.join('_')}.fsrec`;
  const file = readdirSync(dir)
    .filter((f) => f.endsWith(end))
    .map((f) => join(dir, f))
    .sort((a, b) => statSync(b).mtimeMs - statSync(a).mtimeMs)[0];
  if (!file) {
    return `no record ends with ${end}`;
  }
  const kept = join(root, 'e2e-games');
  mkdirSync(kept, { recursive: true });
  const to = join(kept, file.slice(dir.length + 1));
  copyFileSync(file, to);
  return to;
}

test.describe('whole random games', () => {
  test.skip(process.env.FULL_GAMES !== '1', 'slow; run with FULL_GAMES=1 (make e2e-games)');
  test.describe.configure({ timeout: 60 * 60_000 });

  for (const players of [2, 3, 4]) {
    test(`${String(players)} players finish a game with every card set`, async ({ browser }) => {
      const seed = Number(process.env.GAME_SEED ?? Date.now() % 1_000_000);
      // A finished game stays in the lobby list, so each repeat
      // (--repeat-each) hosts under its own name.
      const run = `${String(players)}r${String(test.info().repeatEachIndex)}`;
      const names = ['Ann', 'Bob', 'Cid', 'Dee'].slice(0, players).map((n) => `Rnd${run}${n}`);
      const problems: string[] = [];
      let opened = 0;
      const pages = await startGame(browser, names, {
        setup: everySet,
        open: (page) => {
          // A click that never lands fails fast: the next round looks
          // again, and the stall check sees a game that stops moving.
          page.setDefaultTimeout(10_000);
          watch(page, names[opened++] ?? '?', problems);
          return attachDebugger(page);
        },
      });
      console.log(`${String(players)} players, seed ${String(seed)}: started`);
      test.info().annotations.push({ type: 'seed', description: String(seed) });
      let passed = false;
      try {
        const r = await play(pages, names, problems, mulberry32(seed));
        for (const page of pages) {
          await expect(page.getByRole('dialog', { name: tr('over.title') })).toBeVisible();
        }
        const line = `${String(players)} players, seed ${String(seed)}: ${r.seconds.toFixed(0)} s, ${String(r.turns)} turns, ${String(r.steps)} steps, ${String(r.moves)} moves`;
        console.log(line);
        test.info().annotations.push({ type: 'game', description: line });
        expect(problems).toEqual([]);
        passed = true;
      } finally {
        if (!passed) {
          const kept = keepRecord(names);
          console.log(`${String(players)} players, seed ${String(seed)}: failed; record ${kept}`);
          test.info().annotations.push({ type: 'record', description: kept });
        }
      }
    });
  }
});
