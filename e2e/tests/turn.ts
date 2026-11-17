// Drives prompts through the UI the way a careful player would: throw
// a roll's die, answer choices with the first option, discard the first
// cards, pass when only responding, attack when the turn cannot end
// otherwise. Used to get through a turn in tests.
import { expect, type Locator, type Page } from '@playwright/test';
import { tr, trRe } from './locale';

// Engine enums the waits read, in Go's order (frontend/src/actions.ts
// and stack.ts; pkg/rules_engine).
export const Intent = { pass: 0, endTurn: 1, discard: 2, playLoot: 3, choose: 4, attack: 5, purchase: 6, activate: 7, order: 8, throw: 9 } as const;
export const Ask = { none: 0, priority: 1, discard: 2, gameOver: 3, choose: 4, order: 5, throw: 6 } as const;
export const Purpose = {
  replacement: 0,
  attackTarget: 1,
  monsterSlot: 2,
  penaltyItem: 3,
  penaltyLoot: 4,
  purchase: 5,
  target: 6,
  triggerOrder: 7,
  mode: 8,
  answer: 9,
  startingItem: 10,
  cursed: 11,
  lowest: 12,
  recharge: 13,
  copyAbility: 14,
  die: 15,
  redirect: 16,
  order: 17,
  passCurse: 18,
  room: 19,
  mayReplace: 20,
} as const;
export const Item = { loot: 0, ability: 1, trigger: 2, roll: 3, damage: 4, death: 5, penalty: 6, deathStep: 7 } as const;

// SeenItem is one stack item, as a wait needs it.
export interface SeenItem {
  kind: number;
  controller: number;
  attack: boolean;
  // thrown: a roll's die is thrown (it has its number).
  thrown: boolean;
  card: string;
}

export interface Seen {
  step: number;
  waiting: number;
  you: number;
  // ready: the page shows the latest state and may act on it: nothing
  // behind, nothing sent for this step, no rollback vote open.
  ready: boolean;
  turn: number;
  active: number;
  over: boolean;
  // prompt is what the game waits for: its kind (Ask), purpose (a
  // choice's), and the text's key (a roll's label).
  prompt: { kind: number; purpose: number; key: string };
  // allowed are the intent kinds (Intent) the page may send now; none
  // while it is not ready.
  allowed: number[];
  stack: SeenItem[];
  // target is the running attack's target object; 0 for none.
  target: number;
}

// seen reads the game state from the store (?fake exposes it). Named
// waits read it, not the buttons that stand for it (TS-08, 22.8).
export async function seen(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface S {
      game: {
        step: number;
        allowed?: { kind: number }[] | null;
        rollback_vote?: unknown;
        view: {
          waiting: { player: number; kind: number; purpose?: number; text?: { key?: string } };
          turn: { number: number; active: number };
          over?: boolean;
          stack?: { kind: number; controller: number; attack?: boolean; roll?: number; card?: string }[] | null;
          attack?: { target?: number };
        };
      } | null;
      lobby: { table: { you: number } | null };
      behind: boolean;
      sent: number | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const g = s.game;
    const ready = g !== null && !s.behind && s.sent !== g.step && !(g.rollback_vote && !g.view.over);
    return {
      step: g?.step ?? -1,
      waiting: g?.view.waiting.player ?? -1,
      you: s.lobby.table?.you ?? -1,
      ready,
      turn: g?.view.turn.number ?? 0,
      active: g?.view.turn.active ?? -1,
      over: g?.view.over ?? false,
      prompt: { kind: g?.view.waiting.kind ?? 0, purpose: g?.view.waiting.purpose ?? 0, key: g?.view.waiting.text?.key ?? '' },
      allowed: ready ? (g?.allowed ?? []).map((a) => a.kind) : [],
      stack: (g?.view.stack ?? []).map((it) => ({
        kind: it.kind,
        controller: it.controller,
        attack: it.attack === true,
        thrown: !!it.roll,
        card: it.card ?? '',
      })),
      target: g?.view.attack?.target ?? 0,
    };
  });
}

// A Check is a named wait's condition on the seen state.
export type Check = (s: Seen) => boolean;

// may: the page may send an intent of this kind now (allowed by the
// engine, and the page ready): Attack, Buy, Pass, End turn, Throw.
export function may(kind: number): Check {
  return (s) => s.ready && s.allowed.includes(kind);
}

// asked: the game waits for the page's own player to answer a prompt
// of this kind (Ask), and of this purpose when given.
export function asked(kind: number, purpose?: number): Check {
  return (s) => s.ready && s.waiting === s.you && s.prompt.kind === kind && (purpose === undefined || s.prompt.purpose === purpose);
}

// question: the game asks the page's player this question (its text
// key), e.g. 'ask.one_on_top'.
export function question(key: string): Check {
  return (s) => asked(Ask.choose)(s) && s.prompt.key === key;
}

// stackHolds: the stack holds n items.
export function stackHolds(n: number): Check {
  return (s) => s.stack.length === n;
}

// turnFree: the stack is empty and the page's player may end the
// turn: the turn is back in its action step, nothing open.
export const turnFree: Check = (s) => s.stack.length === 0 && may(Intent.endTurn)(s);

// holds checks a condition on the page's state now, without waiting.
export async function holds(page: Page, check: Check): Promise<boolean> {
  return check(await seen(page));
}

// until waits until the page's state meets the condition. It answers
// nothing: settle does that.
export async function until(page: Page, check: Check, timeout?: number): Promise<void> {
  await expect.poll(async () => check(await seen(page)), timeout === undefined ? {} : { timeout }).toBe(true);
}

// untilPick answers prompts until the page's player is asked to pick,
// for this purpose (Purpose): an attack's target, what to buy, a slot.
export async function untilPick(pages: Page[], press: Tap, page: Page, purpose: number): Promise<void> {
  await settle(pages, press, () => holds(page, asked(Ask.choose, purpose)));
}

// enabled: the element is there and enabled, without waiting for it.
export async function enabled(l: ReturnType<Page['locator']>): Promise<boolean> {
  return (await l.count()) > 0 && (await l.isEnabled());
}

export type Tap = (page: Page, selector: ReturnType<Page['locator']>) => Promise<void>;

export const click: Tap = async (_, l) => {
  await l.click();
};

export const tap: Tap = async (_, l) => {
  await l.tap();
};

// reaches reports whether a click on the target's middle would reach
// it now: nothing covers that point. It reads the page once, with no
// wait. A trial click waited for the target to be stable, two animation
// frames alike, and under a full run's load that often took longer than
// any short timeout: settle then skipped the monster slots and picked
// the monster deck, a move more than the test expected (19.197).
async function reaches(target: Locator): Promise<boolean> {
  return target
    .evaluate((el) => {
      const inView = (r: DOMRect) => r.top >= 0 && r.left >= 0 && r.bottom <= window.innerHeight && r.right <= window.innerWidth;
      if (!inView(el.getBoundingClientRect())) {
        el.scrollIntoView({ block: 'nearest', inline: 'nearest' });
      }
      const r = el.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) {
        return false;
      }
      const hit = document.elementFromPoint(r.left + r.width / 2, r.top + r.height / 2);
      return hit !== null && el.contains(hit);
    })
    .catch(() => false);
}

// uncovered is the first of the targets a click would reach, if any.
async function uncovered(targets: Locator): Promise<Locator | undefined> {
  for (const t of await targets.all()) {
    if (await reaches(t)) {
      return t;
    }
  }
  return undefined;
}

// reachable are all the targets a click would reach.
export async function reachable(targets: Locator): Promise<Locator[]> {
  const out: Locator[] = [];
  for (const t of await targets.all()) {
    if (await reaches(t)) {
      out.push(t);
    }
  }
  return out;
}

// answer makes one move for whichever page must act, if any. It returns
// false when nobody can act right now. steps are the pages' steps when
// the done check began: when a page has moved on since, it makes no move
// and returns true, so done checks the new state first (19.7).
async function answer(pages: Page[], press: Tap, steps: number[]): Promise<boolean> {
  for (const [i, page] of pages.entries()) {
    const s = await seen(page);
    if (s.step !== steps[i]) {
      return true;
    }
    if (!s.ready || s.waiting !== s.you) {
      continue;
    }
    // A roll waits for its roller's throw (15.61).
    const roll = throwButton(page);
    if (await enabled(roll)) {
      await press(page, roll);
      return true;
    }
    // exact: a question may mention "discard" too.
    const discard = page.getByRole('dialog', { name: tr('discard.label'), exact: true });
    if (await discard.isVisible()) {
      const done = discard.getByRole('button', { name: new RegExp(`^${trRe('discard.button').source}`) });
      for (let i = 0; i < 12 && !(await enabled(done)); i++) {
        const next = discard.locator('.card.hl-act').first();
        if ((await next.count()) === 0) {
          break;
        }
        await press(page, next);
      }
      if (await enabled(done)) {
        await press(page, done);
      }
      return true;
    }
    // A "look at" dialog only shows cards; OK closes it.
    const look = page.locator('.choice.look');
    if (await look.isVisible()) {
      await press(page, look.getByRole('button', { name: tr('look.ok') }));
      return true;
    }
    // An order dialog (19.16): place the cards as they come, then confirm.
    const order = page.locator('.choice.order');
    if (await order.isVisible()) {
      const done = order.getByRole('button', { name: tr('order.put_back') });
      for (let i = 0; i < 12 && !(await enabled(done)); i++) {
        await press(page, order.locator('.card.hl-pick').first());
      }
      await press(page, done);
      return true;
    }
    // A "choose one-" card asks its mode with the card large (15.58).
    const mode = page.locator('.choice-backdrop .ability-dialog');
    if (await mode.isVisible()) {
      await press(page, mode.locator('.ability-line:enabled').first());
      return true;
    }
    // A start-of-game pick waits for Confirm (20.22).
    const confirm = page.locator('.choice .choice-confirm');
    if (await enabled(confirm)) {
      await press(page, confirm);
      return true;
    }
    const choice = page.locator('.choice:not(.look)');
    if (await choice.isVisible()) {
      // Prefer picking on the table, as a player would (7.7), where the
      // dialog leaves the target uncovered.
      await press(page, (await uncovered(page.locator('.board .card.hl-pick, .player-pick'))) ?? choice.locator('.choice-option button').first());
      return true;
    }
    // What to attack, what to buy, and a deck attack's slot are picked
    // on the table, with no dialog (18.11). The board first: a monster,
    // as the dialog's first option was, before a player.
    const pick = (await uncovered(page.locator('.board .card.hl-pick'))) ?? (await uncovered(tablePicks(page)));
    if (pick) {
      await press(page, pick);
      return true;
    }
    // A target prompt of the player's own card aims from it, with no
    // dialog (15.34): its targets glow on the table (19.10).
    const aim = await uncovered(page.locator('.card.hl-aim, .stack-row.hl-aim'));
    if (aim) {
      await press(page, aim);
      return true;
    }
    const pass = passButton(page);
    if (await enabled(pass)) {
      await press(page, pass);
      return true;
    }
    // A card may force attacks before the turn can end (Ambush).
    const attack = turnAction(page, 'Attack');
    if (!(await enabled(endTurnButton(page))) && (await enabled(attack))) {
      await press(page, attack);
      return true;
    }
  }
  return false;
}

// settle answers prompts until done() holds.
export async function settle(pages: Page[], press: Tap, done: () => Promise<boolean>): Promise<void> {
  for (let i = 0; i < 600; i++) {
    const steps = await Promise.all(pages.map(async (p) => (await seen(p)).step));
    if (await done()) {
      return;
    }
    if (!(await answer(pages, press, steps))) {
      await pages[0]?.waitForTimeout(100);
    }
  }
  const states = await Promise.all(
    pages.map(async (p) => ({
      ...(await seen(p)),
      status: await p.locator('.prompt-banner').textContent({ timeout: 1000 }).catch(() => null),
      buttons: await p.getByRole('region', { name: tr('actions.label') }).getByRole('button').allTextContents(),
      dialog: await p.locator('.choice').textContent({ timeout: 1000 }).catch(() => null),
      error: await p.locator('.action-error').textContent({ timeout: 1000 }).catch(() => null),
      detail: await p.evaluate(() => {
        interface D {
          game: { step: number; events: unknown; view: { stack: unknown } } | null;
          sent: number | null;
        }
        const s = (window as unknown as { fourSouls: { store: { getState: () => D } } }).fourSouls.store.getState();
        return { step: s.game?.step, sent: s.sent, stack: s.game?.view.stack, events: s.game?.events };
      }),
    })),
  );
  throw new Error(`the game never settled: ${JSON.stringify(states)}`);
}

// tablePicks are the options of a choice picked on the table, with no
// dialog (18.11): what to attack, what to buy, a deck attack's slot.
export function tablePicks(page: Page): Locator {
  return page.locator('.card.hl-pick, .player-pick');
}

// attackRoll: the page's player must throw their attack roll. The
// throw is for the attack roll, the roll is on the stack not yet thrown,
// the attack has its target, and the page may throw. Another roll
// does not count: Bob's Brain rolls when the attack is declared, before
// the target (19.85). A trigger of the target (Pride) resolves before
// the first roll goes on the stack (R-ATK-25, 19.151); `settle` answers
// its prompts on the way. The check does not need the roll on top.
export const attackRoll: Check = (s) =>
  asked(Ask.throw)(s) &&
  may(Intent.throw)(s) &&
  s.prompt.key === 'event.attack_roll' &&
  s.stack.some((it) => it.kind === Item.roll && it.attack && !it.thrown) &&
  s.target !== 0;

export function attackRollWaits(page: Page): Promise<boolean> {
  return holds(page, attackRoll);
}

// untilAttackRoll answers prompts after an Attack click until the
// attacker's attack roll waits for its throw (18.10). Prompts on the way
// get answered: the target, and any card the attack sets off.
export async function untilAttackRoll(pages: Page[], press: Tap, attacker: Page): Promise<void> {
  await settle(pages, press, () => attackRollWaits(attacker));
}

// stackWaits: the page shows the latest state, where the stack holds an
// item and waits for a response to it (V-11). A resolving card's
// question does not count: it shows as a prompt line instead (19.126).
export const stackWaits: Check = (s) => s.ready && s.prompt.kind === Ask.priority && s.stack.length > 0;

// untilStackWaits waits until every page shows the stack waiting for a
// response (19.126).
export async function untilStackWaits(pages: Page[]): Promise<void> {
  await expect.poll(async () => (await Promise.all(pages.map((p) => holds(p, stackWaits)))).every(Boolean)).toBe(true);
}

// markGame keeps the page's game state now, for untilResynced.
export function markGame(page: Page): Promise<void> {
  return page.evaluate(() => {
    interface S {
      game: unknown;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    (window as unknown as { gameMark: unknown }).gameMark = s.game;
  });
}

// untilResynced waits until the page has reconnected and shows the
// server's answer to the reconnect: a new game state, with nothing
// behind. Until then the page may still show the state from before the
// drop, and an intent sent on it is not answered (19.103).
export async function untilResynced(page: Page): Promise<void> {
  await page.waitForFunction(() => {
    interface S {
      game: unknown;
      connected: boolean;
      reconnecting: boolean;
      behind: boolean;
      sent: number | null;
    }
    const w = window as unknown as { fourSouls: { store: { getState: () => S } }; gameMark: unknown };
    const s = w.fourSouls.store.getState();
    return s.connected && !s.reconnecting && !s.behind && s.sent === null && s.game !== w.gameMark;
  });
}

// turnAction is an icon button by own name (15.19): "Buy" or "Attack".
export function turnAction(page: Page, name: 'Buy' | 'Attack') {
  const label = name === 'Buy' ? tr('turn.buy') : tr('turn.attack');
  return page.getByRole('toolbar', { name: tr('turn.actions') }).getByRole('button', { name: label, exact: true });
}

// throwButton is the roller's Throw button where the die lands (15.61).
export function throwButton(page: Page) {
  return page.getByRole('region', { name: tr('throw.label') }).getByRole('button', { name: trRe('throw.button') });
}

// passButton is Pass in the stack panel (15.21).
export function passButton(page: Page) {
  return page.getByRole('complementary', { name: tr('stack.title') }).getByRole('button', { name: tr('stack.pass'), exact: true });
}

// endTurnButton is the big button at the board's right edge (15.20):
// "End turn" on own turn, "Their turn" otherwise.
export function endTurnButton(page: Page, name: 'End turn' | 'Their turn' = 'End turn') {
  const label = name === 'End turn' ? tr('turn.end') : tr('turn.theirs');
  return page.getByRole('region', { name: tr('board.label') }).getByRole('button', { name: label, exact: true });
}

// endTurnQuestion asks before the turn ends while plays are left (V-46).
export function endTurnQuestion(page: Page) {
  return page.getByRole('dialog', { name: tr('turn.ask_title') });
}

// endTurn presses End turn and confirms when it asks: it asks unless it
// glows, as nothing is left to do then (H-20, V-46). Whether it glows
// comes from the store: nothing allowed but ending the turn or passing.
// The button must show the same before it is pressed.
export async function endTurn(page: Page, press: Tap = click) {
  const button = endTurnButton(page);
  const s = await seen(page);
  const idle = s.allowed.includes(Intent.endTurn) && s.allowed.every((k) => k === Intent.endTurn || k === Intent.pass);
  await (idle ? expect(button).toHaveClass(/\bidle\b/) : expect(button).not.toHaveClass(/\bidle\b/));
  await press(page, button);
  if (!idle) {
    await press(page, endTurnQuestion(page).getByRole('button', { name: tr('turn.end') }));
  }
}

// canEndTurn: the page may end the turn now (the engine allows it, and
// the page is ready).
export function canEndTurn(page: Page): Promise<boolean> {
  return holds(page, may(Intent.endTurn));
}

// activePage settles the game until a player may end the turn; then
// that page's End turn button shows it.
export async function activePage(pages: Page[], press: Tap): Promise<Page> {
  let active: Page | undefined;
  await settle(pages, press, async () => {
    for (const page of pages) {
      if (await canEndTurn(page)) {
        active = page;
        return true;
      }
    }
    return false;
  });
  if (!active) {
    throw new Error('no active player');
  }
  await expect(endTurnButton(active)).toBeEnabled();
  return active;
}

export async function turnNumber(page: Page): Promise<number> {
  return (await seen(page)).turn;
}

export async function gameOver(page: Page): Promise<boolean> {
  return (await seen(page)).over;
}

export async function expectNextTurn(pages: Page[], press: Tap, after: number): Promise<void> {
  const first = pages[0];
  if (!first) {
    throw new Error('no pages');
  }
  await settle(pages, press, async () => (await turnNumber(first)) > after || (await gameOver(first)));
  expect((await turnNumber(first)) > after || (await gameOver(first))).toBe(true);
}

// playLoot plays a hand card the way a player would (20.20): a click
// picks it up (a tap lifts it), and a click on the board plays it; a
// card that needs a target aims instead, and the first glowing target
// is pressed (15.25): any card, a character for its player (15.33), or a
// stack item. A new state can end an aim before it is used: the card is
// pressed again until it is held or a target shows (19.100).
export async function playLoot(page: Page, card: Locator, press: Tap): Promise<void> {
  const held = page.locator('.hand-card.held');
  const target = page.locator('.card.hl-aim, .stack-row.hl-aim').first();
  await expect
    .poll(
      async () => {
        if ((await held.count()) > 0 || (await target.isVisible())) {
          return true;
        }
        await press(page, card);
        return false;
      },
      { timeout: 10_000 },
    )
    .toBe(true);
  await press(page, (await held.count()) > 0 ? playSpot(page) : target);
}

// playSpot is where a click plays a held hand card: the board (20.20).
export function playSpot(page: Page): Locator {
  return page.getByRole('region', { name: tr('board.label') });
}

// answerPlay answers what a played loot card asks before it goes on the
// stack: its mode in a dialog, or its target, aimed at on the table
// (19.18): Dagaz protects a player, Lil Battery charges an item.
export async function answerPlay(page: Page, press: Tap): Promise<void> {
  const rows = page.getByRole('complementary', { name: tr('stack.title') }).locator('.stack-row');
  await expect
    .poll(
      async () => {
        if ((await rows.count()) > 0) {
          return true;
        }
        const option = page.locator('.choice .choice-option button, .ability-dialog .ability-line:enabled').first();
        const aim = (await option.isVisible()) ? option : await uncovered(page.locator('.card.hl-aim, .stack-row.hl-aim'));
        if (aim) {
          await press(page, aim);
        }
        return false;
      },
      { timeout: 10_000 },
    )
    .toBe(true);
}
