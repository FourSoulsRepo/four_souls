// Drives prompts through the UI the way a careful player would: answer
// choices with the first option, discard the first cards, pass when
// only responding, attack when the turn cannot end otherwise. Used to
// get through a turn in tests.
import { expect, type Locator, type Page } from '@playwright/test';

interface Seen {
  waiting: number;
  you: number;
  ready: boolean;
  turn: number;
  over: boolean;
}

// seen reads who must act from the store (?fake exposes it).
async function seen(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface S {
      game: { step: number; view: { waiting: { player: number }; turn: { number: number }; over?: boolean } } | null;
      lobby: { table: { you: number } | null };
      behind: boolean;
      sent: number | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return {
      waiting: s.game?.view.waiting.player ?? -1,
      you: s.lobby.table?.you ?? -1,
      ready: s.game !== null && !s.behind && s.sent !== s.game.step,
      turn: s.game?.view.turn.number ?? 0,
      over: s.game?.view.over ?? false,
    };
  });
}

// enabled: the element is there and enabled, without waiting for it.
async function enabled(l: ReturnType<Page['locator']>): Promise<boolean> {
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

// answer makes one move for whichever page must act, if any. It returns
// false when nobody can act right now.
async function answer(pages: Page[], press: Tap): Promise<boolean> {
  for (const page of pages) {
    const s = await seen(page);
    if (!s.ready || s.waiting !== s.you) {
      continue;
    }
    // exact: a question may mention "discard" too.
    const discard = page.getByRole('dialog', { name: 'Discard', exact: true });
    if (await discard.isVisible()) {
      const done = discard.getByRole('button', { name: /^Discard/ });
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
      await press(page, look.getByRole('button', { name: 'OK' }));
      return true;
    }
    const choice = page.locator('.choice:not(.look)');
    if (await choice.isVisible()) {
      // Prefer picking on the table, as a player would (7.7), where the
      // dialog leaves the target uncovered.
      await press(page, (await uncovered(page.locator('.board .card.hl-pick, .player-pick'))) ?? choice.locator('.choice-option button').first());
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
    if (await done()) {
      return;
    }
    if (!(await answer(pages, press))) {
      await pages[0]?.waitForTimeout(100);
    }
  }
  const states = await Promise.all(
    pages.map(async (p) => ({
      ...(await seen(p)),
      status: await p.locator('.action-status').textContent({ timeout: 1000 }).catch(() => null),
      buttons: await p.getByRole('region', { name: 'Actions' }).getByRole('button').allTextContents(),
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
  return page.getByRole('toolbar', { name: 'Turn actions' }).getByRole('button', { name, exact: true });
}

// passButton is Pass in the stack panel (15.21).
export function passButton(page: Page) {
  return page.getByRole('complementary', { name: 'Stack' }).getByRole('button', { name: 'Pass', exact: true });
}

// endTurnButton is the big button at the board's right edge (15.20):
// "End turn" on own turn, "Their turn" otherwise.
export function endTurnButton(page: Page, name: 'End turn' | 'Their turn' = 'End turn') {
  return page.getByRole('region', { name: 'Board' }).getByRole('button', { name, exact: true });
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
