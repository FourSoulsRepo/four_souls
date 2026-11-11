// Trick or Treat on the table (14.33, ADR 028 Decision 11, UC-42):
// Mischievous The Haunt with its trick counters left of the monster
// slots, the Trick/Treat card in the challenge panel, fear counters by
// each player's name, and the monsters' evasion and rewards as the
// challenge changes them. Saved games reach the exact states: fear on
// a player, a changed reward, a sure hit or miss, enough souls, the
// last round.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { repeatNames, startGame, statMonsters } from './game';
import { cardName, tr, trRe } from './locale';
import { editSave, saveEditContinue } from './saves';
import { shot } from './shot';
import { activePage, canEndTurn, click, enabled, endTurn, gameOver, settle, throwButton, turnAction } from './turn';

test.setTimeout(180_000);

// The engine's numbers (state.go, statics.go, cards.go): the next
// attack roll's boost, the monster deck, the in-play zone, the soul
// role, the loot reward.
const NEXT_ATTACK_ROLL = 18;
const MONSTER_DECK = 2;
const ZONE_IN_PLAY = 4;
const ROLE_SOUL = 5;
const REWARD_LOOT = 1;

const haunt = 'mischievous_the_haunt_norm';

interface Counter {
  name?: string;
  count: number;
}

interface SaveObject {
  card: string;
  zone: { kind: number };
  role?: number;
  controller: number;
  damage?: number;
  counters?: Counter[];
  flipped?: boolean;
}

interface TrickSave {
  state: {
    turn: { active: number };
    boosts?: { stat: number; player: number; amount: number }[];
    players: { in_play: number[]; counters?: Counter[] }[];
    objects: SaveObject[];
    decks: number[][];
    monsters: { cards: number[] }[];
    bosses: { cards: number[] }[];
    reward_changes?: { object: number; kind?: number; amount?: number }[];
    challenge: { timer?: number; card?: number };
  };
}

// edit changes the saved challenge game on disk.
function edit(game: string, change: (save: TrickSave) => void) {
  editSave(game, (save) => {
    change(save as unknown as TrickSave);
  });
}

// bossOf is the final boss's object in the save.
function bossOf(save: TrickSave): SaveObject {
  const id = save.state.bosses[0]?.cards.at(-1);
  const obj = id === undefined ? undefined : save.state.objects[id];
  if (!obj) {
    throw new Error('no final boss');
  }
  return obj;
}

// slotTop is the first monster slot's object ID in the save.
function slotTop(save: TrickSave): number {
  const id = save.state.monsters[0]?.cards.at(-1);
  if (id === undefined) {
    throw new Error('no monster in the first slot');
  }
  return id;
}

// boost gives the active player's next attack roll n.
function boost(save: TrickSave, n: number) {
  save.state.boosts = [...(save.state.boosts ?? []), { stat: NEXT_ATTACK_ROLL, player: save.state.turn.active, amount: n }];
}

// giveSouls makes the top n cards of the monster deck souls of the
// active player, so the players hold n souls between them.
function giveSouls(save: TrickSave, n: number) {
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

// trickSide turns the challenge card back to Trick: the tests answer
// its flip question with yes (R-TRICK-11).
function trickSide(save: TrickSave) {
  const obj = save.state.objects[save.state.challenge.card ?? -1];
  if (!obj) {
    throw new Error('no challenge card');
  }
  obj.flipped = false;
}

// swapTop makes the first monster slot's card the given one; the card
// is one of a kind, so another object with it takes the old card.
function swapTop(save: TrickSave, card: string) {
  const { objects } = save.state;
  const obj = objects[slotTop(save)];
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

// Seen is what the page's latest view says of the challenge.
interface Seen {
  tricks: number;
  fear: number[];
  card: string;
  flipped: boolean;
  monsters: { id: number; card: string; dc: number; printed: number; rewards: number }[];
}

async function seenOf(page: Page): Promise<Seen> {
  return page.evaluate(() => {
    interface Card {
      id: number;
      card: string;
      flipped?: boolean;
      counters?: { name?: string; count: number }[];
      stats?: { dc: number; printed_dc: number; rewards?: unknown[] };
    }
    interface S {
      game: {
        view: {
          players: { counters?: { name?: string; count: number }[] }[];
          monsters?: { top?: Card }[];
          bosses?: { top?: Card }[];
          challenge?: { card: Card };
        };
      } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    const boss = v?.bosses?.[0]?.top;
    return {
      tricks: boss?.counters?.find((c) => c.name === 'trick')?.count ?? 0,
      fear: (v?.players ?? []).map((p) => p.counters?.find((c) => c.name === 'fear')?.count ?? 0),
      card: v?.challenge?.card.card ?? '',
      flipped: v?.challenge?.card.flipped ?? false,
      monsters: (v?.monsters ?? []).flatMap((s) =>
        s.top?.stats ? [{ id: s.top.id, card: s.top.card, dc: s.top.stats.dc, printed: s.top.stats.printed_dc, rewards: s.top.stats.rewards?.length ?? 0 }] : [],
      ),
    };
  });
}

function board(page: Page) {
  return page.getByRole('region', { name: tr('board.label') });
}

function boss(page: Page) {
  return board(page).getByRole('group', { name: tr('board.final_boss') });
}

function panel(page: Page) {
  return board(page).getByRole('group', { name: tr('challenge.panel') });
}

// tricks is the trick counter badge on the boss.
function tricks(page: Page, n: number) {
  return boss(page).getByLabel(tr('badge.counters_named', { n, name: 'trick' }), { exact: true });
}

// objectCard is a card on the table by its object ID.
function objectCard(page: Page, id: number) {
  return board(page).locator(`[data-zone="obj-${String(id)}"]`);
}

// trickNormal is the game list's and the table's line for the game.
function trickNormal(): string {
  return tr('lobby.options.challenge', { challenge: tr('challenge.name.trick_or_treat'), difficulty: tr('challenge.difficulty.normal') });
}

// hostTrick hosts Trick or Treat at Normal.
async function hostTrick(host: Page) {
  await host.getByLabel(tr('host.mode')).selectOption('challenge');
  const challenge = host.getByLabel(tr('host.challenge'), { exact: true });
  await expect(challenge.getByRole('option', { name: tr('challenge.name.trick_or_treat') })).toHaveCount(1);
  await challenge.selectOption('trick_or_treat');
  await expect(host.getByLabel(tr('host.difficulty'))).toHaveValue('normal');
  await shot(host, 'trick-host');
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

// asked is the page that shows the dialog, if any.
async function asked(pages: Page[], name: string): Promise<Page | undefined> {
  for (const page of pages) {
    if (await page.getByRole('dialog', { name }).isVisible()) {
      return page;
    }
  }
  return undefined;
}

function first(pages: Page[]): Page {
  const [page] = pages;
  if (!page) {
    throw new Error('no pages');
  }
  return page;
}

// attack declares an attack and picks the target card by its object ID
// once it glows.
async function attack(pages: Page[], active: Page, id: number) {
  await click(active, turnAction(active, 'Attack'));
  const target = board(active).locator(`.card.hl-pick[data-zone="obj-${String(id)}"]`);
  await settle(pages, click, async () => (await target.count()) > 0);
  await click(active, target);
}

test('Trick or Treat at Normal: the table, a flip at turn start, the timer runs out (14.33)', async ({ browser }) => {
  const names = repeatNames(['Tia', 'Tom']);
  const pages = await startGame(browser, names, {
    setup: hostTrick,
    listed: async (row) => {
      await expect(row).toContainText(trickNormal());
    },
  });
  for (const page of pages) {
    // The boss, landscape, with its 12 trick counters (R-TRICK-03).
    const card = boss(page).getByRole('img', { name: cardName(haunt), exact: true });
    await expect(card).toBeVisible();
    const box = await card.boundingBox();
    expect(box && box.width > box.height).toBe(true);
    await expect(tricks(page, 12)).toBeVisible();
    // The challenge card, Trick side up (R-TRICK-04); no presents.
    await expect(panel(page).getByRole('img', { name: cardName('trick_treat_normhardultra'), exact: true })).toBeVisible();
    await expect(panel(page).getByRole('img', { name: tr('challenge.timer', { n: 8 }) })).toBeVisible();
    await expect(board(page).getByRole('group', { name: tr('challenge.presents') })).toHaveCount(0);
  }

  // As the turn starts, the active player may flip it (R-TRICK-11).
  const question = tr('card.trick_treat.flip_to_treat');
  await settle(pages, click, async () => (await asked(pages, question)) !== undefined);
  const flipper = await asked(pages, question);
  if (!flipper) {
    throw new Error('no flip question');
  }
  await shot(flipper, 'trick-flip-ask');
  await click(flipper, flipper.getByRole('dialog', { name: question }).getByRole('button', { name: tr('ask.yes'), exact: true }));
  for (const page of pages) {
    await expect(panel(page).getByRole('img', { name: cardName('treat_normhardultra'), exact: true })).toBeVisible();
    await expect.poll(async () => (await seenOf(page)).flipped).toBe(true);
  }
  await flipper.mouse.move(1, 1);
  await shot(flipper, 'trick-table');

  // The last round: the timer runs out at its end (R-TRICK-16, R-COOP-09).
  await activePage(pages, click);
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      save.state.challenge.timer = 1;
    });
  });
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
  await shot(first(pages), 'trick-lost');
});

test('fear by the names raises the shown evasion; damage gives fear; changed rewards show (14.33)', async ({ browser }) => {
  const names = repeatNames(['Fen', 'Fiz']);
  // Delirium and Stoney change other monsters' evasion (14.10): the
  // test checks the +1 from fear alone.
  const pages = await startGame(browser, names, { setup: hostTrick, bans: statMonsters });
  await activePage(pages, click);
  // 2 fear on the active player: monsters +1 evasion (R-TRICK-09). The
  // first slot's monster gets loot 1 more till end of turn (R-TRICK-15),
  // and the next attack roll -6: it misses (R-TRICK-07).
  let target = 0;
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      const player = save.state.players[save.state.turn.active];
      if (!player) {
        throw new Error('no active player');
      }
      player.counters = [{ name: 'fear', count: 2 }];
      target = slotTop(save);
      save.state.reward_changes = [{ object: target, kind: REWARD_LOOT, amount: 1 }];
      boost(save, -6);
    });
  });
  const active = await activePage(pages, click);
  const now = await seenOf(active);
  const activeSeat = now.fear.indexOf(2);
  expect(activeSeat).toBeGreaterThanOrEqual(0);
  for (const page of pages) {
    // Both screens show the fear by the active player's name.
    await expect(page.getByRole('img', { name: tr('table.fear_title', { n: 2 }), exact: true })).toBeVisible();
    // Each monster shows its evasion one higher, up to 6.
    for (const m of now.monsters) {
      expect(m.dc).toBe(Math.min(m.printed + 1, 6));
      if (m.printed < 6) {
        await expect(objectCard(page, m.id).getByLabel(tr('badge.dc', { n: m.dc }), { exact: true })).toBeVisible();
      }
    }
    // The changed rewards: a ribbon on the card, in words.
    const ribbon = objectCard(page, target).getByRole('img', { name: trRe('stats.rewards_now') });
    await expect(ribbon).toBeVisible();
    await expect(ribbon).toContainText(trRe('stats.reward.loot'));
  }
  // The zoom names the boss as the evasion's source (V-41) and says the
  // rewards.
  const raised = now.monsters.find((m) => m.printed < 6) ?? now.monsters[0];
  if (raised) {
    await objectCard(active, raised.id).hover();
    const zoom = active.locator('.card-zoom');
    await expect(zoom.locator('.card-zoom-stats')).toContainText(cardName(haunt));
    await shot(active, 'trick-fear-evasion');
  }
  await objectCard(active, target).hover();
  await expect(active.locator('.card-zoom .card-zoom-stats')).toContainText(trRe('stats.rewards_now'));
  await shot(active, 'trick-rewards');
  await active.mouse.move(1, 1);

  // A miss: the monster's damage puts fear on the attacker (R-TRICK-07).
  await attack(pages, active, target);
  await settle(pages, click, async () => ((await seenOf(first(pages))).fear[activeSeat] ?? 0) > 2);
  const fear = (await seenOf(first(pages))).fear[activeSeat] ?? 0;
  for (const page of pages) {
    await expect(page.getByRole('img', { name: tr('table.fear_title', { n: fear }), exact: true })).toBeVisible();
  }
  await shot(active, 'trick-fear-gained');
});

test('the boss is guarded; a kill asks the trigger order and removes a trick counter (14.33)', async ({ browser }) => {
  const names = repeatNames(['Kit', 'Kay']);
  const pages = await startGame(browser, names, { setup: hostTrick });
  await activePage(pages, click);
  // A Fly in the first slot and a sure hit: it dies (R-TRICK-12).
  let fly = 0;
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      swapTop(save, 'fly');
      fly = slotTop(save);
      boost(save, 6);
      trickSide(save);
    });
  });
  const active = await activePage(pages, click);

  // 12 trick counters and no soul: the boss is no target; it dims while
  // the targets glow (R-TRICK-06).
  await click(active, turnAction(active, 'Attack'));
  const target = objectCard(active, fly);
  await settle(pages, click, async () => (await active.locator('.board .card.hl-pick').count()) > 0);
  await expect(target).toHaveClass(/hl-pick/);
  await expect(boss(active).locator('.card')).toBeVisible();
  await expect(boss(active).locator('.card.hl-pick')).toHaveCount(0);
  await shot(active, 'trick-boss-guarded');
  await click(active, target);

  // The attack-start roll (Ghostly Trickery) comes before the attack
  // roll: the active player throws it.
  await settle(pages, click, async () => enabled(throwButton(active)));
  await shot(active, 'trick-attack-start-roll');

  // The Fly dies: the Haunt's roll and the Trick side's counter come at
  // once; the active player orders them (R-ABIL-15).
  const order = tr('choice.title.trigger_order');
  await settle(pages, click, async () => (await active.locator('.choice').filter({ hasText: order }).count()) > 0 || (await seenOf(active)).tricks < 12);
  if ((await active.locator('.choice').filter({ hasText: order }).count()) > 0) {
    await shot(active, 'trick-trigger-order');
  }
  await settle(pages, click, async () => (await seenOf(first(pages))).tricks === 11);
  for (const page of pages) {
    await expect(tricks(page, 11)).toBeVisible();
  }
  await shot(active, 'trick-trick-removed');
});

test('the boss can be attacked once souls reach its trick counters; the challenge is won (14.33)', async ({ browser }) => {
  const names = repeatNames(['Wes', 'Wyn']);
  const pages = await startGame(browser, names, { setup: hostTrick });
  await activePage(pages, click);
  // 2 souls and 2 trick counters: the boss may be attacked (R-TRICK-06).
  // It has 1 HP left, and the next attack roll gets +6: it hits.
  await saveEditContinue(pages, names.join(', '), (game) => {
    edit(game, (save) => {
      giveSouls(save, 2);
      const b = bossOf(save);
      b.counters = [{ name: 'trick', count: 2 }];
      b.damage = 3;
      boost(save, 6);
    });
  });
  const name = cardName(haunt);
  for (let i = 0; i < 6; i++) {
    const active = await untilTurn(pages);
    if (!active) {
      break;
    }
    const attackIcon = turnAction(active, 'Attack');
    if (await attackIcon.isEnabled()) {
      await click(active, attackIcon);
      const target = boss(active).getByRole('button', { name, exact: true });
      await settle(pages, click, async () => (await target.count()) > 0);
      await expect(target).toHaveClass(/hl-pick/);
      await expect(tricks(active, 2)).toBeVisible();
      if (i === 0) {
        await shot(active, 'trick-boss-target');
      }
      await click(active, target);
      await settle(pages, click, () => turnOrOver(pages));
    }
    if (!(await gameOver(active)) && (await canEndTurn(active))) {
      await endTurn(active);
    }
  }
  // Both players see the win (R-TRICK-01, R-CHAL-14).
  for (const page of pages) {
    const over = page.getByRole('dialog', { name: tr('over.title') });
    await expect(over.getByRole('heading', { name: tr('over.challenge_won') })).toBeVisible();
  }
  await shot(first(pages), 'trick-won');
});
