// Drives prompts through the UI the way a careful player would: answer
// choices with the first option, discard the first cards, pass when
// only responding. Used to get through a turn in tests.
import { expect, type Page } from '@playwright/test';

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
    const choice = page.locator('.choice');
    if (await choice.isVisible()) {
      // Prefer picking on the table, as a player would (7.7).
      const onTable = page.locator('.board .card.hl-pick, .player-pick').first();
      await press(page, (await onTable.isVisible()) ? onTable : choice.locator('.choice-option button').first());
      return true;
    }
    const pass = page.getByRole('region', { name: 'Actions' }).getByRole('button', { name: 'Pass' });
    if (await enabled(pass)) {
      await press(page, pass);
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

// endTurnButton is the active player's "End turn", when it is offered.
export function endTurnButton(page: Page) {
  return page.getByRole('region', { name: 'Actions' }).getByRole('button', { name: 'End turn' });
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
