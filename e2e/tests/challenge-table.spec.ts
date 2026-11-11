// The challenge table (14.25, ADR 028 Decision 11, UC-42): both players
// see the final boss left of the monster slots, minions in a row under
// them, and the challenge panel by the decks: the rule card, the timer
// and the 4 presents. Saved games reach the exact states: a present at
// 0, enough souls, a known monster, the last round.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { repeatNames, startGame, statMonsters } from './game';
import { cardName, tr } from './locale';
import { editSave, saveEditContinue } from './saves';
import { shot } from './shot';
import { activePage, canEndTurn, click, endTurn, gameOver, settle, turnAction } from './turn';

test.setTimeout(180_000);

// The engine's numbers (state.go, statics.go): the monster deck, the
// in-play zone, the soul role, the next attack roll's boost.
const NEXT_ATTACK_ROLL = 18;
const MONSTER_DECK = 2;
const ZONE_IN_PLAY = 4;
const ROLE_SOUL = 5;

interface SaveObject {
  card: string;
  zone: { kind: number };
  role?: number;
  controller: number;
  damage?: number;
  counters?: { name?: string; count: number }[];
}

interface ChallengeSave {
  state: {
    turn: { active: number };
    boosts?: { stat: number; player: number; amount: number }[];
    players: { in_play: number[] }[];
    objects: SaveObject[];
    decks: number[][];
    monsters: { cards: number[] }[];
    bosses: { cards: number[] }[];
    challenge: { timer?: number; presents: { kind: string; object: number; destroyed?: boolean }[] };
  };
}

// edit changes the saved challenge game on disk.
function edit(game: string, change: (save: ChallengeSave) => void) {
  editSave(game, (save) => {
    change(save as unknown as ChallengeSave);
  });
}

// setCounters puts n counters on every present still there.
function setCounters(save: ChallengeSave, n: number) {
  for (const p of save.state.challenge.presents) {
    const obj = save.state.objects[p.object];
    if (!p.destroyed && obj) {
      obj.counters = [{ count: n }];
    }
  }
}

// giveSouls makes the top n cards of the monster deck souls of the
// active player, so the players hold n souls between them.
function giveSouls(save: ChallengeSave, n: number) {
  const { state } = save;
  const p = state.turn.active;
  for (let i = 0; i < n; i++) {
    const id = state.decks[MONSTER_DECK]?.pop();
    const obj = id === undefined ? undefined : state.objects[id];
    const player = state.players[p];
    if (id === undefined || !obj || !player) {
      throw new Error('no monster card left to make a soul');
    }
    Object.assign(obj, { zone: { kind: ZONE_IN_PLAY }, role: ROLE_SOUL, controller: p });
    player.in_play.push(id);
  }
}

// swapTop makes the first monster slot's card the given one; the card
// is one of a kind, so another object with it takes the old card.
function swapTop(save: ChallengeSave, card: string) {
  const { objects } = save.state;
  const id = save.state.monsters[0]?.cards.at(-1);
  const obj = id === undefined ? undefined : objects[id];
  if (!obj) {
    throw new Error('no monster in the first slot');
  }
  for (const other of objects) {
    if (other.card === card) {
      other.card = obj.card;
    }
  }
  obj.card = card;
  obj.damage = 0;
}

// Table is what the page's latest view says of the challenge: each
// present's counters and whether it is destroyed, and the minions.
interface Table {
  counters: number[];
  destroyed: string[];
  minions: number;
}

async function table(page: Page): Promise<Table> {
  return page.evaluate(() => {
    interface S {
      game: {
        view: {
          minions?: { top?: unknown }[];
          challenge?: { presents?: { kind: string; counters: number; destroyed?: boolean }[] };
        };
      } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    const presents = v?.challenge?.presents ?? [];
    return {
      counters: presents.map((p) => p.counters),
      destroyed: presents.filter((p) => p.destroyed).map((p) => p.kind),
      minions: (v?.minions ?? []).filter((s) => s.top).length,
    };
  });
}

function sum(n: number[]): number {
  return n.reduce((a, b) => a + b, 0);
}

const kinds = ['cents', 'loot', 'treasure', 'soul'] as const;

function caption(kind: (typeof kinds)[number]): string {
  return tr(`challenge.present.${kind}`);
}

function board(page: Page) {
  return page.getByRole('region', { name: tr('board.label') });
}

function boss(page: Page) {
  return board(page).getByRole('group', { name: tr('board.final_boss') });
}

function minionRow(page: Page) {
  return board(page).getByRole('group', { name: tr('board.minions') });
}

function presentsOf(page: Page) {
  return board(page).getByRole('group', { name: tr('challenge.presents') });
}

// hostKrampus hosts the Krampus challenge at Normal.
async function hostKrampus(host: Page) {
  await host.getByLabel(tr('host.mode')).selectOption('challenge');
  await expect(host.getByLabel(tr('host.difficulty'))).toHaveValue('normal');
}

// hostKrampusHard hosts it at Hard: no random minion attack is owed,
// so the player picks the attack's target on the table.
async function hostKrampusHard(host: Page) {
  await hostKrampus(host);
  await host.getByLabel(tr('host.difficulty')).selectOption('hard');
}

// expectPresents: every page shows each face-down present with the
// counters of the view, by its kind.
async function expectPresents(pages: Page[], counters: number[]) {
  for (const page of pages) {
    for (const [i, kind] of kinds.entries()) {
      await expect(presentsOf(page).getByRole('img', { name: tr('challenge.present_down', { present: caption(kind), n: counters[i] ?? 0 }) })).toBeVisible();
    }
  }
}

// turnOrOver: some page may end its turn, or the game is over.
async function turnOrOver(pages: Page[]): Promise<boolean> {
  for (const page of pages) {
    if ((await gameOver(page)) || (await canEndTurn(page))) {
      return true;
    }
  }
  return false;
}

// untilTurn plays until a page may end its turn, and returns it; none
// once the game is over.
async function untilTurn(pages: Page[]): Promise<Page | undefined> {
  await settle(pages, click, () => turnOrOver(pages));
  for (const page of pages) {
    if (await gameOver(page)) {
      return undefined;
    }
  }
  for (const page of pages) {
    if (await canEndTurn(page)) {
      return page;
    }
  }
  return undefined;
}

// first is the host's page.
function first(pages: Page[]): Page {
  const [page] = pages;
  if (!page) {
    throw new Error('no pages');
  }
  return page;
}

test('the challenge table: boss, panel, presents lose counters (14.25)', async ({ browser }) => {
  const pages = await startGame(browser, repeatNames(['Pia', 'Pom']), { setup: hostKrampus });
  const krampus = cardName('pilfering_krampus_norm');
  for (const page of pages) {
    // The final boss, landscape, left of the monster slots.
    const card = boss(page).getByRole('img', { name: krampus, exact: true });
    await expect(card).toBeVisible();
    const box = await card.boundingBox();
    expect(box && box.width > box.height).toBe(true);
    const monsters = await board(page).getByRole('group', { name: tr('board.monsters') }).boundingBox();
    expect(box && monsters && box.x + box.width <= monsters.x).toBe(true);
    // No minions yet: no row.
    await expect(minionRow(page)).toHaveCount(0);
    // The panel: the rule card, landscape, and the timer at 8 rounds.
    const panel = board(page).getByRole('group', { name: tr('challenge.panel') });
    const rule = panel.getByRole('img', { name: cardName('how_the_krampus_stole_christmas_norm'), exact: true });
    await expect(rule).toBeVisible();
    const ruleBox = await rule.boundingBox();
    expect(ruleBox && ruleBox.width > ruleBox.height).toBe(true);
    await expect(panel.getByRole('img', { name: tr('challenge.timer', { n: 8 }) })).toBeVisible();
    await expect(panel.getByText(tr('challenge.timer_label'))).toBeVisible();
  }
  // The presents start with 8, 2, 3 and 3 counters (R-KRAMPUS-05).
  expect((await table(first(pages))).counters).toEqual([8, 2, 3, 3]);
  await expectPresents(pages, [8, 2, 3, 3]);

  // As the turn starts, the active player rolls: counters come off
  // (R-KRAMPUS-12).
  await settle(pages, click, async () => sum((await table(first(pages))).counters) < 16);
  const { counters } = await table(first(pages));
  await expectPresents(pages, counters);
  const active = await activePage(pages, click);
  await shot(active, 'challenge-table');

  // The rule card zooms like any card, landscape.
  await board(active).getByRole('group', { name: tr('challenge.panel') }).getByRole('img', { name: cardName('how_the_krampus_stole_christmas_norm'), exact: true }).hover();
  const zoom = active.locator('.card-zoom .card.size-zoom.landscape');
  await expect(zoom).toBeVisible();
  const zoomBox = await zoom.boundingBox();
  expect(zoomBox && zoomBox.width > zoomBox.height).toBe(true);
  await shot(active, 'challenge-rule-zoom');
});

test('a present at 0 brings minions; the guarded boss is no target (14.25)', async ({ browser }) => {
  const names = repeatNames(['Mia', 'Max']);
  const pages = await startGame(browser, names, { setup: hostKrampusHard });
  await activePage(pages, click);
  // Every present at 1: the next start-of-turn roll takes the last
  // counter of one (R-KRAMPUS-08).
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      setCounters(save, 1);
    });
  });
  const active = await activePage(pages, click);
  await endTurn(active);
  await settle(pages, click, async () => (await table(first(pages))).minions > 0);

  // The present is destroyed; its minions are in the row under the
  // monster slots, on every screen.
  const now = await table(first(pages));
  expect(now.destroyed).toHaveLength(1);
  const gone = kinds.find((k) => now.destroyed.includes(k)) ?? 'cents';
  for (const page of pages) {
    await expect(presentsOf(page).getByRole('img', { name: tr('challenge.present_destroyed', { present: caption(gone) }) })).toBeVisible();
    await expect(minionRow(page).locator('.card')).toHaveCount(now.minions);
    const row = await minionRow(page).boundingBox();
    const monsters = await board(page).getByRole('group', { name: tr('board.monsters') }).boundingBox();
    expect(row && monsters && row.y >= monsters.y + monsters.height).toBe(true);
  }
  await shot(first(pages), 'challenge-minions');

  // The boss is guarded while a minion is in play (R-KRAMPUS-11): the
  // attack's targets glow on the table, the minions among them; the
  // boss is no target and dims.
  const attacker = await untilTurn(pages);
  if (!attacker) {
    throw new Error('the game ended');
  }
  await click(attacker, turnAction(attacker, 'Attack'));
  // The declared attack goes round for responses; then its target.
  const minion = minionRow(attacker).locator('.card.hl-pick').first();
  await settle(pages, click, async () => (await minion.count()) > 0);
  await expect(minion).toBeVisible();
  await expect(boss(attacker).locator('.card')).toBeVisible();
  await expect(boss(attacker).locator('.card.hl-pick')).toHaveCount(0);
  await shot(attacker, 'challenge-boss-guarded');
  await click(attacker, minion);
  await expect(minionRow(attacker).locator('.card.hl-pick')).toHaveCount(0);
});

test('the boss can be attacked once souls are enough; the challenge is won (14.25)', async ({ browser }) => {
  const names = repeatNames(['Bea', 'Bob']);
  // Delirium or Stoney in a slot gives Krampus +1 evasion (14.10,
  // R-CHAL-03): the boosted roll of 1 then misses, and 1 damage a hit
  // can't kill a healed Krampus in 6 turns (19.154).
  const pages = await startGame(browser, names, { setup: hostKrampus, bans: statMonsters });
  await activePage(pages, click);
  // 5 souls and no minion: Krampus may be attacked (R-KRAMPUS-11). It
  // has 1 HP left, and the attacker's next roll gets +3, so the first
  // roll hits (4+) and kills it; damage heals at the end of a turn. The
  // presents keep enough counters that no minion comes meanwhile.
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      giveSouls(save, 5);
      setCounters(save, 9);
      save.state.boosts = [...(save.state.boosts ?? []), { stat: NEXT_ATTACK_ROLL, player: save.state.turn.active, amount: 3 }];
      const id = save.state.bosses[0]?.cards.at(-1);
      const obj = id === undefined ? undefined : save.state.objects[id];
      if (!obj) {
        throw new Error('no final boss');
      }
      obj.damage = 4;
    });
  });
  const krampus = cardName('pilfering_krampus_norm');
  // The boosted roll hits; should it not, the next turns attack again.
  for (let i = 0; i < 6; i++) {
    const active = await untilTurn(pages);
    if (!active) {
      break;
    }
    const attack = turnAction(active, 'Attack');
    if (await attack.isEnabled()) {
      await click(active, attack);
      // The declared attack goes round for responses; then its target.
      const target = boss(active).getByRole('button', { name: krampus, exact: true });
      await settle(pages, click, async () => (await target.count()) > 0);
      await expect(target).toHaveClass(/hl-pick/);
      if (i === 0) {
        await shot(active, 'challenge-boss-target');
      }
      await click(active, target);
      await settle(pages, click, () => turnOrOver(pages));
    }
    if (!(await gameOver(active)) && (await canEndTurn(active))) {
      await endTurn(active);
    }
  }
  // Both players see the win (R-KRAMPUS-01, R-CHAL-14).
  for (const page of pages) {
    const over = page.getByRole('dialog', { name: tr('over.title') });
    await expect(over.getByRole('heading', { name: tr('over.challenge_won') })).toBeVisible();
  }
  await shot(first(pages), 'challenge-won');
});

test('All I Want for Christmas asks; a lost challenge says why (14.25)', async ({ browser }) => {
  const names = repeatNames(['Fay', 'Fox']);
  const pages = await startGame(browser, names, { setup: hostKrampus });
  await activePage(pages, click);
  // A Fly in the first slot: 1 HP and a 2 to hit, its reward 1¢.
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      swapTop(save, 'fly');
      setCounters(save, 9);
    });
  });
  const question = tr('card.how_the_krampus_stole_christmas.forgo_cents', { n: 1 });
  const asked = async () => {
    for (const page of pages) {
      if (await page.getByRole('dialog', { name: question }).isVisible()) {
        return page;
      }
    }
    return undefined;
  };
  let forgo: Page | undefined;
  for (let i = 0; i < 12 && !forgo; i++) {
    const active = await untilTurn(pages);
    if (!active) {
      throw new Error('the game ended before the Fly died');
    }
    const fly = board(active).getByRole('group', { name: tr('board.monsters') }).getByRole('button', { name: cardName('fly'), exact: true });
    const attack = turnAction(active, 'Attack');
    if (await attack.isEnabled()) {
      await click(active, attack);
      await settle(pages, click, async () => (await fly.count()) > 0);
      await click(active, fly);
      await settle(pages, click, async () => (await asked()) !== undefined || (await turnOrOver(pages)));
      forgo = await asked();
    }
    if (!forgo && (await canEndTurn(active))) {
      await endTurn(active);
    }
  }
  if (!forgo) {
    throw new Error('no forgo question');
  }
  // Yes: the 1¢ becomes a counter on the ¢ present (R-KRAMPUS-07).
  const before = (await table(forgo)).counters[0] ?? 0;
  await shot(forgo, 'challenge-forgo');
  await click(forgo, forgo.getByRole('dialog', { name: question }).getByRole('button', { name: tr('ask.yes'), exact: true }));
  await expect.poll(async () => (await table(forgo)).counters[0]).toBe(before + 1);

  // The last round: the timer runs out at its end (R-COOP-09).
  await activePage(pages, click);
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      save.state.challenge.timer = 1;
    });
  });
  for (const page of pages) {
    await expect(board(page).getByRole('img', { name: tr('challenge.timer', { n: 1 }) })).toBeVisible();
  }
  for (let i = 0; i < 4; i++) {
    const active = await untilTurn(pages);
    if (!active) {
      break;
    }
    await endTurn(active);
  }
  for (const page of pages) {
    const over = page.getByRole('dialog', { name: tr('over.title') });
    await expect(over.getByRole('heading', { name: tr('over.challenge_lost') })).toBeVisible();
    await expect(over.getByText(tr('event.challenge_lost.timer'))).toBeVisible();
  }
  await shot(first(pages), 'challenge-lost');
});
