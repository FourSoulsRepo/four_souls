// Player actions (7.7): a full turn with the mouse, and with touch, on
// real games.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL, serverPort } from './server';
import { shot } from './shot';
import {
  activePage,
  canEndTurn,
  click,
  endTurnButton,
  expectNextTurn,
  gameOver,
  passButton,
  playLoot,
  settle,
  tap,
  throwButton,
  turnAction,
  turnNumber,
  type Tap,
} from './turn';

const fast = () => `${appURL()}&anim=off`;

// A whole turn can take many moves (combat rolls, responses).
test.setTimeout(120_000);

// backToMe: the active player may end the turn again, or the turn or
// game moved on (a death ends it early).
function backToMe(active: Page, turn: number) {
  return async () => (await canEndTurn(active)) || (await turnNumber(active)) !== turn || (await gameOver(active));
}

// attack declares an attack; settling picks the target on the table.
async function attack(pages: Page[], active: Page, press: Tap, turn: number) {
  const button = turnAction(active, 'Attack');
  if (!(await canEndTurn(active)) || !(await button.isEnabled())) {
    return;
  }
  await press(active, button);
  await settle(pages, press, backToMe(active, turn));
}

// endTurn ends the turn if it is still the active player's.
async function endTurn(active: Page, press: Tap) {
  if (await canEndTurn(active)) {
    await press(active, endTurnButton(active));
  }
}

test('a full turn with the mouse: drag a loot card, attack, end the turn', async ({ browser }) => {
  const pages = await startGame(browser, ['Max', 'Nia'], { url: fast() });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);

  const loot = active.locator('.hand .card.hl-act').first();
  if (await loot.isVisible()) {
    await shot(active, 'actions');
    const from = await loot.boundingBox();
    const to = await active.getByRole('region', { name: 'Board' }).boundingBox();
    if (!from || !to) {
      throw new Error('no boxes');
    }
    await active.mouse.move(from.x + from.width / 2, from.y + from.height / 2);
    await active.mouse.down();
    await active.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
    await active.mouse.up();
    await settle(pages, click, backToMe(active, turn));
  }
  await attack(pages, active, click, turn);
  await endTurn(active, click);
  await expectNextTurn(pages, click, turn);
});

test('a full turn with touch: tap a loot card to play it, attack, end the turn', async ({ browser }) => {
  const pages = await startGame(browser, ['Oz', 'Pia'], { url: fast(), touch: true });
  const active = await activePage(pages, tap);
  const turn = await turnNumber(active);

  const loot = active.locator('.hand .card.hl-act').first();
  if (await loot.isVisible()) {
    await playLoot(active, loot, tap);
    await settle(pages, tap, backToMe(active, turn));
  }
  await attack(pages, active, tap, turn);
  await endTurn(active, tap);
  await expectNextTurn(pages, tap, turn);
});

test('a refused action shows the reason', async ({ browser }) => {
  const pages = await startGame(browser, ['Quin', 'Ray'], { url: fast() });
  const active = await activePage(pages, click);
  // Send what the server must refuse: the other player cannot end this turn.
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await other.evaluate(() => {
    interface S {
      act: (i: { player: number; kind: number }) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store;
    store.getState().act({ player: 0, kind: 1 });
  });
  await expect(other.getByRole('region', { name: 'Actions' }).getByRole('status')).toContainText(/\(R-[A-Z]+-\d+\)/);
});

// allowedKinds reads the intent kinds the engine allows this page now,
// or none while the page may not act.
async function allowedKinds(page: Page): Promise<number[]> {
  return page.evaluate(() => {
    interface S {
      game: { step: number; allowed?: { kind: number }[] | null } | null;
      behind: boolean;
      sent: number | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const ready = s.game !== null && !s.behind && s.sent !== s.game.step;
    return ready ? (s.game?.allowed ?? []).map((a) => a.kind) : [];
  });
}

test('buy and attack icons are lit only when the engine allows them; B buys', async ({ browser }) => {
  const pages = await startGame(browser, ['Sol', 'Tia'], { url: fast() });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);
  for (const page of pages) {
    const kinds = await allowedKinds(page);
    const buy = turnAction(page, 'Buy');
    const attack = turnAction(page, 'Attack');
    await expect(buy).toBeVisible();
    await expect(attack).toBeVisible();
    await (kinds.includes(6) ? expect(buy).toBeEnabled() : expect(buy).toBeDisabled());
    await (kinds.includes(5) ? expect(attack).toBeEnabled() : expect(attack).toBeDisabled());
  }
  await expect(turnAction(active, 'Buy')).toBeEnabled();
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(turnAction(other, 'Buy')).toBeDisabled();
  await expect(turnAction(other, 'Attack')).toBeDisabled();
  await shot(active, 'turn-actions');

  // B declares the purchase; once the others pass, the Buy icon stays
  // pressed while the player picks what to buy (F-13).
  await active.keyboard.press('b');
  await settle(pages, click, async () => active.locator('.choice').isVisible());
  const choosing = active
    .getByRole('toolbar', { name: 'Turn actions' })
    .getByRole('button', { name: 'Choosing what to buy' });
  await expect(choosing).toHaveAttribute('aria-pressed', 'true');
  await shot(active, 'turn-actions-buy');
  await settle(pages, click, backToMe(active, turn));
});

// 19.7: an update that lands after the done check must not be answered.
// Updates come late; the slow check sees no dialog, then waits while the
// last pass opens it.
test('settle never answers a dialog its done check has not seen (19.7)', async ({ browser }) => {
  let lag = 0;
  const late = (page: Page) =>
    page.routeWebSocket(`ws://127.0.0.1:${String(serverPort())}/ws`, (ws) => {
      const server = ws.connectToServer();
      let last = Promise.resolve();
      ws.onMessage((msg) => {
        server.send(msg);
      });
      server.onMessage((msg) => {
        const at = Date.now() + lag;
        last = last
          .then(() => new Promise((r) => setTimeout(r, Math.max(at - Date.now(), 0))))
          .then(() => {
            ws.send(msg);
          });
      });
    });
  const pages = await startGame(browser, ['Ari', 'Ozy'], { url: fast(), open: late });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);
  lag = 150;
  await active.keyboard.press('b');
  // Back to the action phase without a dialog: settle answered it.
  await settle(pages, click, async () => {
    const done = (await active.locator('.choice').isVisible()) || (await canEndTurn(active));
    await active.waitForTimeout(300);
    return done;
  });
  await expect(active.locator('.choice')).toBeVisible();
  lag = 0;
  await settle(pages, click, backToMe(active, turn));
});

// recordActs keeps every intent the page sends in window.acts.
async function recordActs(page: Page): Promise<void> {
  await page.evaluate(() => {
    interface I {
      kind: number;
    }
    interface S {
      act: (i: I) => void;
    }
    const w = window as unknown as {
      acts: I[];
      fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } };
    };
    const { store } = w.fourSouls;
    const act = store.getState().act;
    w.acts = [];
    store.setState({
      act: (i) => {
        w.acts.push(i);
        act(i);
      },
    });
  });
}

async function lastAct(page: Page): Promise<number | undefined> {
  return page.evaluate(() => (window as unknown as { acts: { kind: number }[] }).acts.at(-1)?.kind);
}

test('the big End turn button: E ends the turn, Space passes', async ({ browser }) => {
  const pages = await startGame(browser, ['Yan', 'Zoe'], { url: fast() });
  const active = await activePage(pages, click);
  const turn = await turnNumber(active);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(endTurnButton(active)).toBeEnabled();
  await expect(endTurnButton(other, 'Their turn')).toBeDisabled();
  // The game's fonts are registered from their URLs; without the files
  // they fail to load and the system font shows.
  const families = await active.evaluate(() => [...document.fonts].map((f) => f.family.replace(/"/g, '')));
  expect(families).toEqual(expect.arrayContaining(['Four Souls Title', 'Four Souls Cyrillic']));
  await shot(active, 'end-turn-own');
  // Off the hand, so no card zoom covers the board.
  await other.mouse.move(1200, 420);
  await shot(other, 'end-turn-theirs');

  // Space passes priority, like the Pass button (F-15).
  for (const page of pages) {
    await recordActs(page);
  }
  await active.keyboard.press('b');
  const passing = async () => {
    for (const page of pages) {
      const pass = passButton(page);
      if ((await pass.count()) > 0 && (await pass.isEnabled())) {
        return page;
      }
    }
    return undefined;
  };
  await expect.poll(passing).toBeDefined();
  const page = await passing();
  if (!page) {
    throw new Error('nobody may pass');
  }
  await page.keyboard.press(' ');
  expect(await lastAct(page)).toBe(0);
  await settle(pages, click, backToMe(active, turn));

  // E ends the turn, like the button.
  if (await canEndTurn(active)) {
    await active.keyboard.press('e');
    expect(await lastAct(active)).toBe(1);
    await expectNextTurn(pages, click, turn);
  }
});

// feed plays the current view again as the next update, with events;
// player 'you' is the page's seat, 'other' the next one.
async function feed(page: Page, events: Record<string, unknown>[]) {
  await page.evaluate((events) => {
    interface S {
      game: { step: number; events: unknown };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const you = s.lobby.table.you;
    const g = structuredClone(s.game);
    g.step += 1;
    g.events = events.map((e) => ({ ...e, player: e.player === 'you' ? you : 1 - you }));
    s.receive({ type: 'update', data: g });
  }, events);
}

test('a "look at the top card" effect shows the card until it is closed (V-07)', async ({ browser }) => {
  const [pip, rue] = await startGame(browser, ['Pip', 'Rue'], { url: fast() });
  if (!pip || !rue) {
    throw new Error('no pages');
  }
  // Forever Alone's "Look at the top card of a deck", as the server
  // sends it: the looker gets the card, the others only the event.
  await feed(pip, [
    { kind: 'stack_resolved', player: 'you', card: 'forever_alone' },
    { kind: 'looked_at', player: 'you', object: 900, card: 'a_dime', private: true },
  ]);
  await feed(rue, [{ kind: 'looked_at', player: 'other', private: true }]);
  const look = pip.getByRole('dialog', { name: 'You look at' });
  await expect(look.getByRole('img', { name: 'A Dime!!' })).toBeVisible();
  await shot(pip, 'look-at');
  // It stays through the next update; only OK closes it.
  await feed(pip, []);
  await expect(look).toBeVisible();
  await look.getByRole('button', { name: 'OK' }).click();
  await expect(look).toHaveCount(0);
  await expect(rue.getByRole('dialog', { name: 'You look at' })).toHaveCount(0);
});

test('a "look, you may put it on the bottom" question shows the card (19.11)', async ({ browser }) => {
  const [opal] = await startGame(browser, ['Opal', 'Reed'], { url: fast() });
  if (!opal) {
    throw new Error('no pages');
  }
  // Cain's Eye at the start of the turn, as the server sends it to the
  // chooser: the question names the deck top and the view holds it.
  await opal.evaluate(() => {
    interface S {
      game: { view: Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const options = ['keep a_dime on top', 'put a_dime on the bottom'];
    game.view.subject = { id: 901, card: 'a_dime', charged: false };
    game.view.waiting = { kind: 4, player: you, purpose: 9, text: 'Put the top card on the bottom?', options };
    game.view.options = options.map(() => ({ player: -1 }));
    game.allowed = options.map((_, i) => ({ player: you, kind: 4, choice: i }));
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  });
  const dialog = opal.getByRole('dialog', { name: 'Put the top card on the bottom?' });
  await expect(dialog).toBeVisible();
  await shot(opal, 'look-may-bottom');
  await expect(dialog.getByRole('img', { name: 'A Dime!!' })).toBeVisible();
  // The buttons name the card by its title, not its ref.
  await expect(dialog.getByRole('button', { name: 'keep A Dime!! on top' })).toBeVisible();
});

test('"put them back in any order" shows each card in its place (19.16)', async ({ browser }) => {
  const [olaf] = await startGame(browser, ['Olaf', 'Rhea'], { url: fast() });
  if (!olaf) {
    throw new Error('no pages');
  }
  // Sleight of Hand's order prompt, as the server sends it to the
  // chooser. act is caught, so the test reads the intent it would send.
  const title = 'Put the cards back in any order';
  await olaf.evaluate((text) => {
    interface S {
      game: { view: Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: object) => void } } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    const refs = ['a_dime', 'a_penny', 'bomb'];
    game.view.ordering = refs.map((card, i) => ({ id: 901 + i, card, charged: false }));
    game.view.waiting = { kind: 5, player: you, purpose: 17, count: 3, text };
    game.allowed = [{ player: you, kind: 8, objects: [901, 902, 903] }];
    store.setState({
      act: (intent: unknown) => {
        (window as unknown as { sentOrder: unknown }).sentOrder = intent;
      },
    });
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  }, title);
  const dialog = olaf.getByRole('dialog', { name: title });
  await expect(dialog).toBeVisible();
  const place = (n: number) => dialog.getByRole('group', { name: `Place ${String(n)}` });
  const left = dialog.getByRole('group', { name: 'Not placed' });
  const putBack = dialog.getByRole('button', { name: 'Put back' });
  await expect(place(1)).toContainText('1 · top');
  await expect(putBack).toBeDisabled();
  // Bomb on top, then A Dime!!; A Penny! is still to place.
  await left.getByRole('button', { name: 'Bomb!' }).click();
  await left.getByRole('button', { name: 'A Dime!!' }).click();
  await expect(place(1).getByRole('button', { name: 'Bomb!' })).toBeVisible();
  await expect(place(2).getByRole('button', { name: 'A Dime!!' })).toBeVisible();
  await shot(olaf, 'order-two-placed');
  // A click on a placed card takes it back; the later ones move up.
  await place(1).getByRole('button', { name: 'Bomb!' }).click();
  await expect(place(1).getByRole('button', { name: 'A Dime!!' })).toBeVisible();
  await expect(left.getByRole('button', { name: 'Bomb!' })).toBeVisible();
  await left.getByRole('button', { name: 'Bomb!' }).click();
  await left.getByRole('button', { name: 'A Penny!' }).click();
  await expect(left).toHaveCount(0);
  await expect(place(2).getByRole('button', { name: 'Bomb!' })).toBeVisible();
  await expect(place(3).getByRole('button', { name: 'A Penny!' })).toBeVisible();
  await shot(olaf, 'order-all-placed');
  await putBack.click();
  const sent = await olaf.evaluate(() => (window as unknown as { sentOrder: unknown }).sentOrder);
  expect(sent).toMatchObject({ kind: 8, objects: [901, 903, 902] });
});

test('two "instead" effects: each card says what it gives first (19.17)', async ({ browser }) => {
  const [ugo] = await startGame(browser, ['Ugo', 'Wren'], { url: fast() });
  if (!ugo) {
    throw new Error('no pages');
  }
  // Bum-Bo! and Counterfeit Penny replace a gain of 10¢; the order
  // prompt as the server sends it to the chooser. act is caught.
  const title = 'Gaining 10¢: which card applies first?';
  await ugo.evaluate((text) => {
    interface S {
      game: { view: Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: object) => void } } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    game.view.waiting = { kind: 4, player: you, purpose: 0, options: ['counterfeit_penny', 'bum_bo'], text };
    game.view.options = [
      { object: 901, player: -1, outcome: { actions: [{ kind: 14, player: you, amount: 11, card: 'bum_bo' }] } },
      {
        object: 902,
        player: -1,
        outcome: { actions: [{ kind: 14, player: you, amount: 10, card: 'bum_bo' }], unused: ['counterfeit_penny'] },
      },
    ];
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    store.setState({
      act: (intent: unknown) => {
        (window as unknown as { sentChoice: unknown }).sentChoice = intent;
      },
    });
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  }, title);
  const dialog = ugo.getByRole('dialog', { name: title });
  await expect(dialog).toBeVisible();
  const option = (card: string) => dialog.locator('.choice-option').filter({ has: ugo.getByRole('button', { name: card, exact: true }) });
  await expect(option('Counterfeit Penny').getByRole('list', { name: 'Result' })).toHaveText('→ Bum-Bo! levels up by 11');
  await expect(option('Bum-Bo!').getByRole('list', { name: 'Result' })).toHaveText(
    '→ Bum-Bo! levels up by 10→ Counterfeit Penny does nothing',
  );
  await shot(ugo, 'instead-order');
  await option('Counterfeit Penny').locator('button').click();
  const sent = await ugo.evaluate(() => (window as unknown as { sentChoice: unknown }).sentChoice);
  expect(sent).toMatchObject({ kind: 4, choice: 0 });
});

test('trigger order: the title says the first one resolves last (19.19)', async ({ browser }) => {
  const [quill] = await startGame(browser, ['Quill', 'Zinnia'], { url: fast() });
  if (!quill) {
    throw new Error('no pages');
  }
  // Cain's Eye and Golden Horseshoe trigger at the start of the turn;
  // the order prompt as the server sends it to the chooser. act is caught.
  await quill.evaluate(() => {
    interface S {
      game: { view: Record<string, unknown>; allowed: unknown[] };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: object) => void } } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table.you;
    game.view.waiting = {
      kind: 4,
      player: you,
      purpose: 7,
      options: [
        'cains_eye: At the start of your turn, look at the top card of the loot deck. You may put it on the bottom.',
        'golden_horseshoe: At the start of your turn, look at the top card of the treasure deck. You may put it on the bottom.',
      ],
    };
    game.view.options = [{ player: -1 }, { player: -1 }];
    game.allowed = [
      { player: you, kind: 4, choice: 0 },
      { player: you, kind: 4, choice: 1 },
    ];
    store.setState({
      act: (intent: unknown) => {
        (window as unknown as { sentChoice: unknown }).sentChoice = intent;
      },
    });
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
  });
  const dialog = quill.getByRole('dialog', { name: /resolves last/ });
  await shot(quill, 'trigger-order');
  await expect(dialog).toBeVisible();
  const option = (card: string) => dialog.locator('.choice-option').filter({ has: quill.getByRole('button', { name: card, exact: true }) });
  await expect(option('Cain’s Eye')).toContainText('look at the top card of the loot deck');
  await expect(option('Golden Horseshoe')).toContainText('look at the top card of the treasure deck');
  await expect(dialog).not.toContainText('cains_eye');
  await option('Golden Horseshoe').locator('button').click();
  const sent = await quill.evaluate(() => (window as unknown as { sentChoice: unknown }).sentChoice);
  expect(sent).toMatchObject({ kind: 4, choice: 1 });
});

// recordRolls keeps every dice_rolled the page's updates bring, in
// window.rolls: the values the engine drew.
async function recordRolls(page: Page): Promise<void> {
  await page.evaluate(() => {
    interface S {
      game: { step: number; events: { kind: string; amount?: number }[] | null } | null;
    }
    const w = window as unknown as {
      rolls: number[];
      finals: number[];
      fourSouls: { store: { subscribe: (f: (s: S, old: S) => void) => void } };
    };
    w.rolls = [];
    w.finals = [];
    w.fourSouls.store.subscribe((s, old) => {
      if (s.game && s.game.step !== old.game?.step) {
        for (const e of s.game.events ?? []) {
          if (e.kind === 'dice_rolled' && e.amount !== undefined) {
            w.rolls.push(e.amount);
          }
          // Every value a roll takes; the last one is the final (15.63).
          if (['dice_rolled', 'roll_changed', 'roll_resolved'].includes(e.kind) && e.amount) {
            w.finals.push(e.amount);
          }
        }
      }
    });
  });
}

async function rolls(page: Page): Promise<number[]> {
  return page.evaluate(() => (window as unknown as { rolls: number[] }).rolls);
}

test('an attack waits for the roller to throw; the die shows the engine\'s value; the last roll stays (15.61, 15.63)', async ({ browser }) => {
  const pages = await startGame(browser, ['Thea', 'Ruan']);
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await expect(turnAction(active, 'Attack')).toBeEnabled();
  await click(active, turnAction(active, 'Attack'));
  // Settle stops before the throw: the attack's target is picked first.
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);

  // The roller sees the glowing button and the stack says it waits for them.
  const stack = (page: Page) => page.getByRole('complementary', { name: 'Stack' });
  await expect(throwButton(active)).toBeEnabled();
  await expect(stack(active).getByLabel('Waiting for')).toHaveText('Throw the die');
  await expect(active.locator('.prompt-banner.mine .prompt-title')).toHaveText(/^Throw the die/);
  await shot(active, 'throw-button');
  // The others see who throws, and the die shaking in their hand.
  const roller = await nameOf(active);
  const shaking = other.getByRole('region', { name: 'Throw' });
  await expect(shaking.locator('.throw-shake .die')).toHaveAttribute('aria-label', `${String(roller)}'s die`);
  await expect(shaking).toContainText(`${String(roller)} is throwing…`);
  await expect(throwButton(other)).toHaveCount(0);
  await expect(stack(other).getByLabel('Waiting for')).toHaveText(new RegExp(`^${String(roller)}: Throw the die`));
  await expect(other.locator('.prompt-title')).toHaveText(`Waiting for ${String(roller)}`);
  await shot(other, 'throw-others');

  // Space throws. The die that lands shows the engine's value, the same
  // on both screens.
  for (const page of pages) {
    await recordRolls(page);
  }
  await active.keyboard.press(' ');
  const effect = active.getByTestId('effect-dice').locator('.die');
  await expect(effect).toBeVisible();
  const label = await effect.getAttribute('aria-label');
  await expect(active.getByRole('region', { name: 'Throw' })).toHaveCount(0);
  await expect(other.getByRole('region', { name: 'Throw' })).toHaveCount(0, { timeout: 10_000 });
  for (const page of pages) {
    await expect.poll(async () => (await rolls(page)).map((v) => `Rolled ${String(v)}`)).toEqual([label]);
  }

  // Once the attack is over, the last roll stays as a small die by the
  // stack with its final value (15.63); hover names the roller.
  await settle(pages, click, backToMe(active, await turnNumber(active)));
  for (const page of pages) {
    const last = stack(page).getByRole('group', { name: 'Last roll' });
    const final = await page.evaluate(() => (window as unknown as { finals: number[] }).finals.at(-1));
    await expect(last.locator('.die')).toHaveAttribute('aria-label', `Rolled ${String(final)}`);
    await last.hover();
    await expect(last.getByRole('tooltip')).toBeVisible();
    await expect(last.getByRole('tooltip')).toHaveText(new RegExp(`^(Thea|Ruan) rolled ${String(final)}( — .+)?$`));
  }
  await shot(active, 'last-roll');
});

// giveLooks gives the page two dice looks of plain colors, so a test
// tells them apart without the private textures, and picks one as the
// viewer's own. The game's seats keep the looks sent in hello.
async function giveLooks(page: Page, own: string): Promise<Record<string, string>> {
  return page.evaluate((own) => {
    interface S {
      settings: { dice: string };
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (p: object) => void } } }).fourSouls.store;
    const paint = (color: string) => {
      const c = document.createElement('canvas');
      c.width = c.height = 96;
      const g = c.getContext('2d');
      if (g) {
        g.fillStyle = color;
        g.fillRect(0, 0, 96, 96);
      }
      return c.toDataURL();
    };
    const images = { tear: paint('#b33'), gold: paint('#db3') };
    store.setState({
      dice: Object.entries(images).map(([id, image]) => ({ id, name: id, image })),
      settings: { ...store.getState().settings, dice: own },
    });
    return images;
  }, own);
}

test("every screen draws a roll with the roller's dice look (15.6)", async ({ browser }) => {
  const pages = await startGame(browser, ['Odile', 'Pavo']);
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  // Both said hello with the default look; each seat lists it.
  for (const page of pages) {
    const looks = await page.evaluate(() => {
      interface S {
        game: { seats: { dice?: string }[] | null } | null;
      }
      const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
      return s.game?.seats?.map((x) => x.dice);
    });
    expect(looks).toEqual(['tear', 'tear']);
  }
  // The other player now likes gold dice; the roller's look still wins.
  const images = await giveLooks(active, 'tear');
  await giveLooks(other, 'gold');
  const image = (die: ReturnType<Page['locator']>) => die.locator('img').getAttribute('src');

  await expect(turnAction(active, 'Attack')).toBeEnabled();
  await click(active, turnAction(active, 'Attack'));
  await settle(pages, click, async () => (await throwButton(active).count()) > 0);
  const shaking = other.getByRole('region', { name: 'Throw' }).locator('.throw-shake .die');
  await expect(shaking).toBeVisible();
  expect(await image(shaking)).toBe(images.tear);
  await shot(other, 'dice-roller-look');

  await active.keyboard.press(' ');
  const landed = active.getByTestId('effect-dice').locator('.die');
  await expect(landed).toBeVisible();
  expect(await image(landed)).toBe(images.tear);
  await settle(pages, click, backToMe(active, await turnNumber(active)));
  for (const page of pages) {
    const last = page.getByRole('complementary', { name: 'Stack' }).getByRole('group', { name: 'Last roll' }).locator('.die');
    expect(await image(last)).toBe(images.tear);
  }
});

// nameOf is the page's own nickname, from the seats of its update.
async function nameOf(page: Page): Promise<string | undefined> {
  return page.evaluate(() => {
    interface S {
      game: { seats: { seat: number; name: string }[] | null } | null;
      lobby: { table: { you: number } | null };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    return s.game?.seats?.find((x) => x.seat === s.lobby.table?.you)?.name;
  });
}
