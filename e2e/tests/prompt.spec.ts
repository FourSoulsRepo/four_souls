// Prompt banners (15.28, V-20): every prompt says what to do to the
// player it asks, and whom the game waits for to everyone else. The
// stack panel repeats it; while the player picks a target, the cards
// that answer glow and the rest dim.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle, throwButton, turnAction } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

// Ask is a prompt fed to a page as the server would send it: waiting
// is the prompt; mine says whether it asks this page's player or the
// other one. The asked player gets the intents that answer it.
interface Ask {
  waiting: { kind: number; purpose?: number; count?: number; text?: string };
  answer: 'none' | 'pass' | 'choose' | 'discard' | 'order' | 'throw';
  mine: boolean;
  stack?: boolean;
}

// ask feeds the prompt to the page's store and returns the ids it used.
async function ask(page: Page, a: Ask, step: number): Promise<{ monster: number; character: number }> {
  return page.evaluate(
    ({ a, step }) => {
      interface Card {
        id: number;
        card: string;
      }
      interface S {
        game: {
          view: {
            players: { character: Card; hand?: Card[] }[];
            monsters: { top?: Card }[];
            waiting: unknown;
            stack: unknown;
            options?: unknown;
            ordering?: unknown;
            passed?: number[];
            responding?: number[];
          };
          allowed: unknown[];
        };
        lobby: { table: { you: number } };
        receive: (env: unknown) => void;
      }
      const store = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store;
      const s = store.getState();
      const game = structuredClone(s.game);
      const you = s.lobby.table.you;
      const v = game.view;
      const asker = a.mine ? you : 1 - you;
      const hand = v.players[you]?.hand ?? [];
      const monster = v.monsters.find((m) => m.top)?.top ?? { id: 0, card: '' };
      const character = v.players[you]?.character.id ?? 0;
      v.waiting = { ...a.waiting, player: asker, options: [monster.card] };
      v.stack = a.stack
        ? [{ id: 99, kind: 0, controller: 1 - asker, source: hand[0]?.id ?? 0, card: hand[0]?.card, target: { player: asker }, ability: { card: '', index: 0 }, roll_for: { card: '', index: 0 } }]
        : null;
      v.passed = [];
      v.responding = a.stack ? [asker] : [];
      v.options = [{ object: monster.id, player: 0 }];
      v.ordering = hand.slice(0, 2);
      const ids = hand.map((c) => c.id);
      const answers: Record<string, unknown[]> = {
        none: [],
        pass: [{ player: you, kind: 0 }],
        choose: [{ player: you, kind: 4, choice: 0 }],
        discard: [{ player: you, kind: 2, objects: ids }],
        order: [{ player: you, kind: 8, objects: ids.slice(0, 2) }],
        throw: [{ player: you, kind: 9 }],
      };
      game.allowed = a.mine ? (answers[a.answer] ?? []) : [];
      s.receive({ type: 'update', data: { ...game, step, events: [] } });
      return { monster: monster.id, character };
    },
    { a, step },
  );
}

const banner = (page: Page) => page.getByRole('region', { name: 'Prompt' });
const stackLine = (page: Page) => page.getByRole('complementary', { name: 'Stack' }).getByLabel('Waiting for');
const card = (page: Page, id: number) => page.locator(`.card[data-zone="obj-${String(id)}"]`);

// dimmed: the card is dimmed for a pick (15.28); a spent card's own
// darker look is another filter.
async function dimmed(l: Locator): Promise<boolean> {
  return l.evaluate((el) => getComputedStyle(el).filter.includes('brightness(0.45)'));
}

// nameOf is the page's own nickname, from the seats of its update.
async function nameOf(page: Page): Promise<string> {
  return page.evaluate(() => {
    interface S {
      game: { seats: { seat: number; name: string }[] };
      lobby: { table: { you: number } };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return s.game.seats.find((x) => x.seat === s.lobby.table.you)?.name ?? '';
  });
}

test('the banner says what to do, or whom the game waits for (15.28)', async ({ browser }) => {
  const pages = await startGame(browser, ['Pell', 'Brin'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const me = await nameOf(active);

  // The turn: the active player is told it is theirs and what they may do.
  await expect(banner(active)).toHaveClass(/mine/);
  await expect(banner(active).locator('.prompt-title')).toHaveText('Your turn');
  await expect(banner(active).locator('.prompt-hint')).toHaveText(/end the turn/);
  await expect(banner(other)).not.toHaveClass(/mine/);
  await expect(banner(other).locator('.prompt-title')).toHaveText(`Waiting for ${me}`);
  await expect(banner(other).locator('.prompt-hint')).toHaveText('Playing their turn');
  await shot(active, 'prompt-your-turn');
  await shot(other, 'prompt-waiting');

  // Every other prompt, to the player it asks and to the other one.
  // title is what the asked player reads; their stack panel says the
  // same, and the other's names the asked player first.
  const prompts: { name: string; a: Omit<Ask, 'mine'>; title: string; hint?: string }[] = [
    { name: 'respond', a: { waiting: { kind: 1 }, answer: 'pass', stack: true }, title: 'Respond, or pass', hint: 'On the stack: ' },
    { name: 'choose', a: { waiting: { kind: 4, purpose: 6 }, answer: 'choose' }, title: 'Choose a target' },
    { name: 'discard', a: { waiting: { kind: 2, count: 1 }, answer: 'discard' }, title: 'Discard 1 loot card' },
    { name: 'order', a: { waiting: { kind: 5, purpose: 12, count: 2 }, answer: 'order' }, title: 'Choose the card for the very bottom' },
    { name: 'throw', a: { waiting: { kind: 6, text: 'attack roll' }, answer: 'throw' }, title: 'Throw the die', hint: 'For the attack roll' },
  ];
  let step = 100;
  for (const p of prompts) {
    step++;
    await ask(active, { ...p.a, mine: true }, step);
    await ask(other, { ...p.a, answer: 'none', mine: false }, step);
    await expect(banner(active), p.name).toHaveClass(/mine/);
    await expect(banner(active).locator('.prompt-title'), p.name).toHaveText(p.title);
    if (p.hint) {
      await expect(banner(active).locator('.prompt-hint'), p.name).toContainText(p.hint);
    }
    await expect(banner(other).locator('.prompt-title'), p.name).toHaveText(`Waiting for ${me}`);
    await expect(banner(other).locator('.prompt-hint'), p.name).toContainText(p.title);
    if (p.name !== 'respond') {
      await expect(stackLine(active).first(), p.name).toHaveText(p.title);
      await expect(stackLine(other).first(), p.name).toHaveText(new RegExp(`^${me}: ${p.title}`));
    }
    await shot(active, `prompt-${p.name}`);
  }
  // The other player is never told it is their move.
  await expect(banner(other)).not.toHaveClass(/mine/);
});

test('answers glow and the rest dims while the player picks (15.28)', async ({ browser }) => {
  const pages = await startGame(browser, ['Yarrow', 'Kestrel'], { url: fast() });
  const active = await activePage(pages, click);

  // A response: Pass answers it, so it glows.
  await ask(active, { waiting: { kind: 1 }, answer: 'pass', mine: true, stack: true }, 100);
  await expect(active.locator('.stack-pass')).toHaveClass(/answer/);

  // A target on the table: it glows; the player's own cards dim.
  const { monster, character } = await ask(active, { waiting: { kind: 4, purpose: 6 }, answer: 'choose', mine: true }, 101);
  await expect(active.locator('main.game')).toHaveClass(/picking/);
  await expect(card(active, monster)).toHaveClass(/hl-pick/);
  await expect.poll(() => dimmed(card(active, monster))).toBe(false);
  await expect.poll(() => dimmed(card(active, character))).toBe(true);
  await expect(active.locator('.stack-pass')).not.toHaveClass(/answer/);
  await shot(active, 'prompt-pick');

  // The player's own aim: only its targets glow, and the banner says how
  // to cancel it.
  await ask(active, { waiting: { kind: 1 }, answer: 'pass', mine: true }, 102);
  await expect(active.locator('main.game')).not.toHaveClass(/picking/);
  await expect.poll(() => dimmed(card(active, character))).toBe(false);
  await active.evaluate(
    ({ monster, character }) => {
      interface S {
        setAiming: (a: unknown) => void;
      }
      const store = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store;
      store.getState().setAiming({ source: character, targets: [{ object: monster, player: 0 }] });
    },
    { monster, character },
  );
  await expect(banner(active).locator('.prompt-title')).toHaveText('Choose a target');
  await expect(banner(active).locator('.prompt-hint')).toHaveText(/Esc cancels$/);
  await expect(card(active, monster)).toHaveClass(/hl-aim/);
  await expect(active.locator('main.game')).toHaveClass(/picking/);
  await shot(active, 'prompt-aim');
  await active.keyboard.press('Escape');
  await expect(active.locator('main.game')).not.toHaveClass(/picking/);
});

test('the others see what the active player does (15.30)', async ({ browser }) => {
  const pages = await startGame(browser, ['Tarn', 'Ulla'], { url: fast() });
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const name = await nameOf(active);
  const doing = other.locator('.player.active .player-activity');
  // Holding priority on their turn: thinking. The own seat says nothing.
  await expect(doing).toHaveText(`${name} is thinking`);
  await expect(active.locator('.player-activity')).toHaveCount(0);
  await shot(other, 'activity-thinking');

  // Attacking: the line names the monster, as the stack's roll does.
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
  const monster = await other.getByRole('complementary', { name: 'Stack' }).locator('.stack-row .stack-target').first().textContent();
  await expect(doing).toHaveText(`${name} attacks ${monster ?? ''}`);
  await expect(doing.locator('.stack-target')).toHaveText(monster ?? '');
  await shot(other, 'activity-attacks');
});
