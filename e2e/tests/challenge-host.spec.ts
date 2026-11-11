// The host screen's game mode (14.24, ADR 028, UC-41): a regular game
// by default, or a co-op challenge with its difficulty. A challenge has
// exactly 2 seats and no bonus souls; the game list and the table name
// it, and the game starts with its final boss.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { tr } from './locale';
import { shot } from './shot';

// Played is what the page's game view says about the mode.
interface Played {
  challenge: string;
  difficulty: string;
  bosses: number;
}

async function played(page: Page): Promise<Played> {
  return page.evaluate(() => {
    interface S {
      game: { view: { challenge?: { id: string; difficulty: string }; bosses?: unknown[] } } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    return { challenge: v?.challenge?.id ?? '', difficulty: v?.challenge?.difficulty ?? '', bosses: v?.bosses?.length ?? 0 };
  });
}

// krampusHard is the game list's and the table's line for the game.
function krampusHard(): string {
  return tr('lobby.options.challenge', { challenge: tr('challenge.name.krampus'), difficulty: tr('challenge.difficulty.hard') });
}

test('a regular game stays the default', async ({ browser }) => {
  const pages = await startGame(browser, ['Rey', 'Reg', 'Rio'], {
    setup: async (host) => {
      await expect(host.getByLabel(tr('host.mode'))).toHaveValue('regular');
      await expect(host.getByLabel(tr('host.challenge'), { exact: true })).toHaveCount(0);
      await expect(host.getByLabel(tr('host.difficulty'))).toHaveCount(0);
      // The players stay 2–4; the test picked 3.
      await expect(host.getByLabel(tr('host.players'))).toBeEnabled();
      await expect(host.getByLabel(tr('host.players'))).toHaveValue('3');
      await shot(host, 'host-regular');
      await host.getByLabel(tr('host.detailed')).check();
      await expect(host.getByLabel(tr('host.bonus_souls'))).toBeVisible();
      await host.getByLabel(tr('host.bonus_souls')).scrollIntoViewIfNeeded();
      await shot(host, 'host-regular-detailed');
      await host.getByLabel(tr('host.detailed')).uncheck();
    },
    listed: async (row) => {
      await expect(row).not.toContainText(tr('challenge.name.krampus'));
    },
    roll: false,
  });
  for (const page of pages) {
    expect(await played(page)).toEqual({ challenge: '', difficulty: '', bosses: 0 });
  }
});

test('the host hosts the Krampus challenge at Hard for 2', async ({ browser }) => {
  const pages = await startGame(browser, ['Kim', 'Kai'], {
    setup: async (host) => {
      await host.getByLabel(tr('host.mode')).selectOption('challenge');
      // Krampus is the only challenge; Normal is its first difficulty.
      await expect(host.getByLabel(tr('host.challenge'), { exact: true })).toHaveValue('krampus');
      await expect(host.getByLabel(tr('host.difficulty'))).toHaveValue('normal');
      await host.getByLabel(tr('host.difficulty')).selectOption('hard');
      const seats = host.getByLabel(tr('host.players'));
      await expect(seats).toBeDisabled();
      await expect(seats).toHaveValue('2');
      await expect(host.getByText(tr('host.challenge_seats', { n: 2 }))).toBeVisible();
      await shot(host, 'host-challenge');
      // Detailed keeps every option but bonus souls.
      await host.getByLabel(tr('host.detailed')).check();
      await expect(host.getByLabel(tr('host.picking'))).toBeVisible();
      await expect(host.getByLabel(tr('host.rollback'))).toBeVisible();
      await expect(host.getByLabel(tr('host.random_seats'))).toBeVisible();
      await expect(host.getByLabel(tr('host.bonus_souls'))).toHaveCount(0);
      await host.getByLabel(tr('host.random_seats')).scrollIntoViewIfNeeded();
      await shot(host, 'host-challenge-detailed');
      // The simple setup sends the challenge too (GS-10).
      await host.getByLabel(tr('host.detailed')).uncheck();
    },
    listed: async (row) => {
      await expect(row).toContainText(krampusHard());
    },
    // The joiner's table names it before the start.
    seated: async (guest) => {
      await expect(guest.getByText(krampusHard())).toBeVisible();
    },
    roll: false,
  });
  for (const page of pages) {
    expect(await played(page)).toEqual({ challenge: 'krampus', difficulty: 'hard', bosses: 1 });
  }
});
