// What the turn owes and who joins an attack (22.5): the turn fields the
// client read nowhere before (19.180 – 19.188). Each test feeds both
// pages one update made from the real one, so the deal does not matter.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click } from './turn';
import { cardName, tr } from './locale';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(90_000);

// Table is what a test needs of the real update: the active seat, the
// other seat and its name, and a monster in a slot.
interface Table {
  active: number;
  other: number;
  otherName: string;
  monster: { id: number; card: string };
}

async function table(page: Page): Promise<Table> {
  return page.evaluate(() => {
    interface Store {
      getState: () => {
        game: {
          seats: { seat: number; name: string }[] | null;
          view: { turn: { active: number }; players: { id: number }[]; monsters: { top?: { id: number; card: string } }[] };
        } | null;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const game = store.getState().game;
    const view = game?.view;
    const top = view?.monsters.find((m) => m.top)?.top;
    if (!view || !top) {
      throw new Error('no game');
    }
    const active = view.turn.active;
    const other = view.players.find((p) => p.id !== active)?.id ?? -1;
    const otherName = game.seats?.find((seat) => seat.seat === other)?.name ?? '';
    return { active, other, otherName, monster: { id: top.id, card: top.card } };
  });
}

// Patch is what a fed update changes: turn fields, the open attack, the
// stack, the seats, and fields of the active player.
interface Patch {
  turn?: Record<string, unknown>;
  attack?: Record<string, unknown>;
  stack?: Record<string, unknown>[];
  seats?: Record<string, Record<string, unknown>>;
  activePlayer?: Record<string, unknown>;
}

// feed sends each page the next update: its own view with the patch.
async function feed(pages: Page[], patch: Patch): Promise<void> {
  for (const page of pages) {
    await page.evaluate((p: Patch) => {
      interface Game {
        step: number;
        seats: { seat: number }[] | null;
        view: { turn: { active: number }; players: { id: number }[]; attack?: unknown; stack?: unknown[] };
      }
      interface Store {
        getState: () => { game: Game | null; receive: (env: unknown) => void };
      }
      const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
      const s = store.getState();
      const game = structuredClone(s.game);
      if (!game) {
        throw new Error('no game');
      }
      const { view } = game;
      view.turn = { ...view.turn, ...p.turn };
      if (p.attack) {
        view.attack = p.attack;
      }
      if (p.stack) {
        view.stack = p.stack;
      }
      if (p.seats) {
        game.seats = (game.seats ?? []).map((seat) => ({ ...seat, ...p.seats?.[String(seat.seat)] }));
      }
      if (p.activePlayer) {
        view.players = view.players.map((pl) => (pl.id === view.turn.active ? { ...pl, ...p.activePlayer } : pl));
      }
      s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [] } });
    }, patch);
  }
}

async function twoPages(browser: Parameters<typeof startGame>[0], names: string[]): Promise<{ active: Page; other: Page; t: Table }> {
  const pages = await startGame(browser, names, { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  return { active, other, t: await table(active) };
}

const banner = (page: Page) => page.getByRole('region', { name: tr('banner.label') });
const monsterCard = (page: Page, id: number) => page.locator(`.card[data-zone="obj-${String(id)}"]`);

test('the monster a player must attack is marked, and the banner names it (19.180, 19.181)', async ({ browser }) => {
  const { active, other, t } = await twoPages(browser, ['Mia', 'Ned']);
  await feed([active, other], { turn: { must_attack: t.monster.id, must_attacks: 2 } });
  await shot(active, 'must-attack-monster');
  await expect(banner(active).locator('.prompt-hint')).toHaveText(tr('banner.must_attack_left', { target: cardName(t.monster.card), n: 2 }));
  // Everyone sees the mark on the monster.
  for (const page of [active, other]) {
    await expect(monsterCard(page, t.monster.id)).toHaveClass(/\bmust-attack\b/);
  }
  // A forced attack with no monster named: only the count.
  await feed([active, other], { turn: { must_attack: 0, must_attacks: 1 } });
  await expect(banner(active).locator('.prompt-hint')).toHaveText(tr('banner.must_attack_any', { n: 1 }));
  await expect(monsterCard(active, t.monster.id)).not.toHaveClass(/\bmust-attack\b/);
});

test('forced attacks on the monster deck mark the deck (19.182)', async ({ browser }) => {
  const { active, other, t } = await twoPages(browser, ['Ola', 'Pim']);
  await feed([active, other], { turn: { must_attacks: 2, must_attack_deck: 2 } });
  await shot(active, 'must-attack-deck');
  await expect(banner(active).locator('.prompt-hint')).toHaveText(tr('banner.must_attack_left', { target: tr('banner.monster_deck'), n: 2 }));
  for (const page of [active, other]) {
    await expect(page.locator('[data-zone="deck-monster"]')).toHaveClass(/\bmust-attack\b/);
    await expect(monsterCard(page, t.monster.id)).not.toHaveClass(/\bmust-attack\b/);
  }
});

test('the attack count says which attacks are deck only (19.183)', async ({ browser }) => {
  const { active, other } = await twoPages(browser, ['Quin', 'Rex']);
  await feed([active, other], { turn: { deck_attacks: 1 }, activePlayer: { attacks: 2 } });
  await shot(other, 'deck-attacks-other');
  await shot(active, 'deck-attacks-own');
  const line = tr('table.attacks_left_deck', { n: 2, deck: 1 });
  for (const page of [active, other]) {
    await expect(page.locator('.player-allowance').getByText(line, { exact: true })).toBeVisible();
  }
});

test('everyone sees who helps the next attack and shares its rewards (19.184, 19.185)', async ({ browser }) => {
  const { active, other, t } = await twoPages(browser, ['Sol', 'Tam']);
  const helper = t.otherName;
  await feed([active, other], { turn: { next_helpers: [t.other], next_share: [t.other] } });
  await shot(active, 'next-attack-helpers');
  const help = `${tr('attack.helped_by', { players: helper })} · ${tr('attack.shared_with', { players: helper })}`;
  for (const page of [active, other]) {
    await expect(page.locator('.player-next-attack')).toHaveText(tr('table.next_attack', { line: help }));
  }
});

test('a running attack shows its helpers on the stack (19.188)', async ({ browser }) => {
  const { active, other, t } = await twoPages(browser, ['Uma', 'Vic']);
  const helper = t.otherName;
  await feed([active, other], { attack: { target: t.monster.id, helpers: [t.other], share: [t.other] } });
  await shot(other, 'attack-helpers');
  const help = `${tr('attack.helped_by', { players: helper })} · ${tr('attack.shared_with', { players: helper })}`;
  for (const page of [active, other]) {
    const stack = page.getByRole('complementary', { name: tr('stack.title') });
    await expect(stack.locator('.stack-help')).toHaveText(help);
  }
});
