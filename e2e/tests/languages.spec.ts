// Each client shows the server's lines in its own language (16.4,
// LO-04): two players, one in the tests' language, one in another.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { locale, trIn } from './locale';
import { appURL } from './server';
import { shot } from './shot';
import { enabled } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

// other is the second page's language: the test language when the
// tests run in English, else English.
const other = locale === 'qps' ? 'en' : 'qps';

interface Seen {
  you: number;
  step: number;
  idle: boolean;
  turn: number;
  active: number;
  waiting: { kind: number; player: number };
  stack: number;
  hand: number[];
  blocked: { kind: number; object: number; reason: { key: string } }[];
  character: number;
}

// look reads what a page's store holds now.
function look(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface S {
      game: {
        step: number;
        blocked?: { kind: number; object: number; reason: { key: string } }[];
        view: {
          turn: { number: number; active: number };
          waiting: { kind: number; player: number };
          stack: unknown[] | null;
          players: { hand?: { id: number }[]; character: { id: number } }[];
        };
      };
      lobby: { table: { you: number } };
      playing: unknown;
      behind: boolean;
      sent: number | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const { view } = s.game;
    const you = s.lobby.table.you;
    return {
      you,
      step: s.game.step,
      idle: s.playing === null && !s.behind && s.sent !== s.game.step,
      turn: view.turn.number,
      active: view.turn.active,
      waiting: view.waiting,
      stack: (view.stack ?? []).length,
      hand: (view.players[you]?.hand ?? []).map((c) => c.id),
      blocked: s.game.blocked ?? [],
      character: view.players[you]?.character.id ?? 0,
    };
  });
}

// act sends an intent for the page's own seat, as its buttons would.
async function act(page: Page, intent: { kind: number; objects?: number[]; choice?: number }): Promise<void> {
  await page.evaluate((i) => {
    interface S {
      act: (i: unknown) => void;
      lobby: { table: { you: number } };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    s.act({ ...i, player: s.lobby.table.you });
  }, intent);
}

// answerUntil makes the page whose move it is throw, pass, or pick the
// first option, until done (a deal may ask, e.g. Eden's starting item).
async function answerUntil(pages: Page[], done: (seen: Seen[]) => boolean): Promise<void> {
  for (let i = 0; i < 300; i++) {
    const seen = await Promise.all(pages.map(look));
    if (done(seen)) {
      return;
    }
    // Throw a die, pick, or pass on a stack; never pass on an empty one.
    const j = seen.findIndex((s) => s.idle && s.waiting.player === s.you && (s.waiting.kind === 6 || s.waiting.kind === 4 || (s.waiting.kind === 1 && s.stack > 0)));
    const page = pages[j];
    const s = seen[j];
    if (page && s) {
      await act(page, s.waiting.kind === 6 ? { kind: 9 } : s.waiting.kind === 4 ? { kind: 4, choice: 0 } : { kind: 0 });
      await expect.poll(async () => (await look(page)).step).toBeGreaterThan(s.step);
    } else {
      await pages[0]?.waitForTimeout(100);
    }
  }
  throw new Error(`the game never got there: ${JSON.stringify(await Promise.all(pages.map(look)))}`);
}

// untilRollBanner answers Eden's starting-item choice when the deal gives
// her one: it comes before the first-player roll (R-SETUP-09, 19.94). It
// returns once every page shows the wanted banner text, so a deal with
// Eden and one without both end on the roll's prompt.
async function untilRollBanner(pages: Page[], want: (page: Page, s: Seen) => string): Promise<void> {
  for (let i = 0; i < 300; i++) {
    const seen = await Promise.all(pages.map(look));
    const shown = await Promise.all(pages.map((page) => page.locator('.prompt-banner .prompt-hint').textContent({ timeout: 300 }).catch(() => null)));
    const done = pages.every((page, j) => {
      const s = seen[j];
      return s !== undefined && shown[j] === want(page, s);
    });
    if (done) {
      return;
    }
    // The page that rolls may sit on Eden's choice; pick its first option.
    const j = seen.findIndex((s) => s.waiting.player === s.you && s.turn < 1);
    const eden = pages[j]?.locator('.choice:not(.look) .choice-option button').first();
    if (eden && (await enabled(eden))) {
      await eden.click();
    } else {
      await pages[0]?.waitForTimeout(100);
    }
  }
  throw new Error(`the roll banner never showed: ${JSON.stringify(await Promise.all(pages.map(look)))}`);
}

test('two players see the server in their own languages (16.4)', async ({ browser }) => {
  // Cain goes first with no first-player roll, whose prompt this test
  // reads (15.43).
  const pages = await startGame(browser, ['Lena', 'Mito'], { url: fast(), roll: false, bans: ['cain'] });
  const [first, second] = pages;
  if (!first || !second) {
    throw new Error('no pages');
  }
  // The second app picks the other language in its settings; it applies
  // at once (16.2).
  await second.evaluate(async (lang) => {
    interface S {
      settings: object;
      saveSettings: (s: object) => Promise<boolean>;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    await s.saveSettings({ ...s.settings, language: lang });
  }, other);
  const langOf = (page: Page) => (page === second ? other : locale);

  // A prompt: the first-player roll names its roll, a server line, in
  // each page's language, for the roller and for the one who waits.
  const label = (page: Page) => trIn(langOf(page), 'event.first_player_roll');
  // The wait ends on the roll's banner; Eden's choice, if dealt, is made first.
  await untilRollBanner(pages, (page, s) => trIn(langOf(page), s.waiting.player === s.you ? 'banner.throw_for' : 'banner.throw_for_other', { label: label(page) }));
  for (const page of pages) {
    const s = await look(page);
    const hint = page.locator('.prompt-banner .prompt-hint');
    const want = s.waiting.player === s.you ? 'banner.throw_for' : 'banner.throw_for_other';
    await expect(hint).toHaveText(trIn(langOf(page), want, { label: label(page) }));
  }
  await shot(second, `languages-${other}-throw`);
  // Turn 1 starts with the active player free to act.
  const free = (s: Seen) => s.idle && s.turn >= 1 && s.stack === 0 && s.waiting.kind === 1 && s.waiting.player === s.active;
  await answerUntil(pages, (seen) => seen.every(free));

  // A refusal: the player who waits for their turn can't play loot; the
  // card's zoom says why in that page's language.
  const seen = await Promise.all(pages.map(look));
  const activeAt = seen.findIndex((s) => s.active === s.you);
  const active = pages[activeAt];
  const waiter = pages[1 - activeAt];
  const w = seen[1 - activeAt];
  if (!active || !waiter || !w) {
    throw new Error('no active page');
  }
  // A click asks why; the answer shows at the card, in that page's
  // language, and the zoom says it after (15.29).
  const id = w.hand[0] ?? 0;
  const loot = waiter.locator(`.card[data-zone="obj-${String(id)}"]`);
  await loot.click();
  await expect(waiter.locator('.aim-why')).toHaveText(trIn(langOf(waiter), 'refuse.loot_your_turn'));
  await waiter.mouse.move(1, 1);
  await loot.hover();
  await expect(waiter.getByTestId('card-zoom').locator('.card-zoom-why')).toHaveText(trIn(langOf(waiter), 'refuse.loot_your_turn'));
  await shot(waiter, `languages-${langOf(waiter)}-why`);
  await waiter.mouse.move(1, 1);

  // A history line from a server event: the active player uses their
  // character's ↷; once it resolves, each History names its words with
  // the server's label line in that page's language.
  const a = seen[activeAt];
  await act(active, { kind: 7, objects: [a?.character ?? 0], choice: 0 });
  await answerUntil(pages, (now) => now.every((s) => free(s) && s.step > (a?.step ?? 0)));
  const words = '↷: Play an additional loot card this turn.';
  for (const page of pages) {
    const lang = langOf(page);
    await page.getByRole('button', { name: trIn(lang, 'history.title'), exact: true }).click();
    const history = page.getByRole('complementary', { name: trIn(lang, 'history.title'), exact: true });
    await expect(history.locator('li').last()).toContainText(trIn(lang, 'event.words', { text: words }));
    await shot(page, `languages-${lang}-history`);
  }
});
