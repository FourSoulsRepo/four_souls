// Drives prompts through the UI the way a careful player would: throw
// a roll's die, answer choices with the first option, discard the first
// cards, pass when only responding, attack when the turn cannot end
// otherwise. Used to get through a turn in tests.
import { expect, type Locator, type Page } from '@playwright/test';
import { tr, trRe } from './locale';

export interface Seen {
  step: number;
  waiting: number;
  you: number;
  ready: boolean;
  turn: number;
  active: number;
  over: boolean;
}

// seen reads who must act from the store (?fake exposes it).
export async function seen(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface S {
      game: { step: number; view: { waiting: { player: number }; turn: { number: number; active: number }; over?: boolean } } | null;
      lobby: { table: { you: number } | null };
      behind: boolean;
      sent: number | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return {
      step: s.game?.step ?? -1,
      waiting: s.game?.view.waiting.player ?? -1,
      you: s.lobby.table?.you ?? -1,
      ready: s.game !== null && !s.behind && s.sent !== s.game.step,
      turn: s.game?.view.turn.number ?? 0,
      active: s.game?.view.turn.active ?? -1,
      over: s.game?.view.over ?? false,
    };
  });
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

// uncovered is the first of the targets a click would reach, if any.
async function uncovered(targets: Locator): Promise<Locator | undefined> {
  for (const t of await targets.all()) {
    if (await t.click({ trial: true, timeout: 300 }).then(() => true, () => false)) {
      return t;
    }
  }
  return undefined;
}

// reachable are all the targets a click would reach.
export async function reachable(targets: Locator): Promise<Locator[]> {
  const out: Locator[] = [];
  for (const t of await targets.all()) {
    if (await t.click({ trial: true, timeout: 300 }).then(() => true, () => false)) {
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
    const choice = page.locator('.choice:not(.look)');
    if (await choice.isVisible()) {
      // Prefer picking on the table, as a player would (7.7), where the
      // dialog leaves the target uncovered.
      await press(page, (await uncovered(page.locator('.board .card.hl-pick, .player-pick'))) ?? choice.locator('.choice-option button').first());
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
// glows, as nothing is left to do then (H-20, V-46).
export async function endTurn(page: Page, press: Tap = click) {
  const button = endTurnButton(page);
  const asks = !(await button.evaluate((b) => b.classList.contains('idle')));
  await press(page, button);
  if (asks) {
    await press(page, endTurnQuestion(page).getByRole('button', { name: tr('turn.end') }));
  }
}

// canEndTurn: the page may end the turn now (offered and enabled).
export async function canEndTurn(page: Page): Promise<boolean> {
  return (await seen(page)).ready && (await enabled(endTurnButton(page)));
}

// activePage settles the game until a player may end the turn.
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

// playLoot plays a hand card the way a player would: its menu's Play,
// or, for a card that needs a target, the first glowing target (15.25):
// any card, a character for its player (15.33), or a stack item.
export async function playLoot(page: Page, card: Locator, press: Tap): Promise<void> {
  await press(page, card);
  const menu = page.getByRole('menu');
  const target = page.locator('.card.hl-aim, .stack-row.hl-aim').first();
  await expect(menu.or(target)).toBeVisible();
  await press(page, (await menu.isVisible()) ? menu.getByRole('menuitem', { name: tr('actions.play') }) : target);
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
