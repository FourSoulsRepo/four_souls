// The draft pick screen (14.11, GS-01, UC-39): the host picks a draft
// in detailed setup; each player sees only their own offers, each with
// its starting item; Eden shows none and picks one at the start.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { tr } from './locale';
import { shot } from './shot';

// Pick is what the page's pick screen offers, from its store.
interface Pick {
  offers: string[];
  items: string[];
  names: Record<string, string>;
}

async function offersOf(page: Page): Promise<Pick> {
  return page.evaluate(() => {
    interface S {
      lobby: { setup: { offers?: string[]; items?: string[] } | null };
      cards: Map<string, { name: string }>;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const offers = s.lobby.setup?.offers ?? [];
    const items = s.lobby.setup?.items ?? [];
    const names: Record<string, string> = {};
    for (const ref of [...offers, ...items]) {
      names[ref] = s.cards.get(ref)?.name ?? ref;
    }
    return { offers, items, names };
  });
}

// draft sets the host screen to a draft of n.
function draft(n: number) {
  return async (host: Page) => {
    await host.getByLabel(tr('host.detailed')).check();
    await expect(host.getByLabel(tr('host.draft_size'))).toHaveCount(0);
    await host.getByLabel(tr('host.picking')).selectOption('draft');
    await host.getByLabel(tr('host.draft_size')).selectOption(String(n));
  };
}

test('random stays the default picking', async ({ browser }) => {
  await startGame(browser, ['Ann', 'Bob'], {
    setup: async (host) => {
      await host.getByLabel(tr('host.detailed')).check();
      await expect(host.getByLabel(tr('host.picking'))).toHaveValue('random');
    },
  });
});

test('a draft offers each player characters with their starting items', async ({ browser }) => {
  const picks: string[] = [];
  // Two players are offered 5 of the Base Game's 11 characters each,
  // so one is left out. Banning The Forgotten, which the test does not
  // check, leaves exactly 10: Eden is always offered to someone.
  const pages = await startGame(browser, ['Ann', 'Bob'], {
    setup: draft(5),
    bans: ['the_forgotten'],
    roll: false,
    pick: async (pages) => {
      const seen: string[] = [];
      for (const [i, page] of pages.entries()) {
        await expect(page.getByRole('heading', { name: tr('pick.title') })).toBeVisible();
        await expect.poll(async () => (await offersOf(page)).offers.length).toBe(5);
        const { offers, items, names } = await offersOf(page);
        // Each sees only their own offers.
        expect(offers.filter((c) => seen.includes(c))).toEqual([]);
        seen.push(...offers);
        for (const [j, c] of offers.entries()) {
          const item = items[j] ?? '';
          const box = page.getByTestId(`item-${c}`);
          if (c === 'eden') {
            expect(item).toBe('');
            await expect(box).toHaveText(tr('pick.item_later'));
          } else {
            expect(item).not.toBe('');
            await expect(box.getByRole('img', { name: names[item] ?? item })).toBeVisible();
          }
        }
        if (i === 0) {
          await shot(page, 'draft-pick');
        }
      }
      expect(seen).toContain('eden');
      // Eden's player takes her; the other takes their first offer.
      for (const page of pages) {
        const { offers } = await offersOf(page);
        const card = offers.includes('eden') ? 'eden' : (offers[0] ?? '');
        picks.push(card);
        const confirm = page.getByRole('button', { name: tr('pick.confirm') });
        await expect(confirm).toBeDisabled();
        await page.getByTestId(`offer-${card}`).click();
        await confirm.click();
      }
    },
  });
  // The game starts with the picks; Eden chooses her item now.
  const eden = pages[picks.indexOf('eden')];
  await expect(eden!.getByRole('dialog', { name: tr('choice.title.starting_item') })).toBeVisible();
});

test('others see who has picked, not what', async ({ browser }) => {
  await startGame(browser, ['Cy', 'Dee'], {
    setup: draft(2),
    pick: async ([first, second]) => {
      const { offers } = await offersOf(first!);
      await first!.getByTestId(`offer-${offers[0] ?? ''}`).click();
      await first!.getByRole('button', { name: tr('pick.confirm') }).click();
      await expect(first!.getByText(tr('pick.waiting'))).toBeVisible();
      const seats = second!.getByRole('list', { name: tr('pick.players') });
      await expect(seats.getByText(`Cy: ${tr('pick.has_picked')}`)).toBeVisible();
      await expect(seats.getByText(`Dee: ${tr('pick.choosing')}`)).toBeVisible();
      const mine = await offersOf(second!);
      await second!.getByTestId(`offer-${mine.offers[1] ?? ''}`).click();
      await second!.getByRole('button', { name: tr('pick.confirm') }).click();
    },
  });
});
