// Aiming (15.25): one click on a card that needs a target starts aiming,
// the next click picks the target. A wrong target says why; the source
// does nothing; empty space or a right click cancels (15.33). A card
// with no target is dragged onto the board; one that asks its mode
// first aims after it, and a cancelled aim puts it back (15.34).
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { saveAndContinue } from './saves';
import { shot } from './shot';
import { activePage, click, settle } from './turn';
import { cardName, tr, trRe } from './locale';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

interface Store {
  getState: () => {
    game: { step: number; events: { kind: string; object?: number; card?: string }[] | null; view: { stack: StackItem[] | null } } | null;
    lobby: { table: { game: string } | null };
    sent: number | null;
  };
  subscribe: (f: (s: ReturnType<Store['getState']>) => void) => void;
}

interface StackItem {
  card?: string;
  targets?: { object?: number }[];
}

test('a bomb aimed at a monster by two clicks damages it', async ({ browser }) => {
  const pages = await startGame(browser, ['Bom', 'Cid'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Bom, Cid');
  const active = await activePage(pages, click);
  // Damage events are kept as they show, update by update.
  await active.evaluate(() => {
    const w = window as unknown as { fourSouls: { store: Store }; damaged: number[] };
    w.damaged = [];
    let step = -1;
    w.fourSouls.store.subscribe((s) => {
      if (s.game && s.game.step !== step) {
        step = s.game.step;
        w.damaged.push(...(s.game.events ?? []).filter((e) => e.kind === 'damaged').map((e) => e.object ?? -1));
      }
    });
  });

  const bomb = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('bomb') }).first();
  const monster = active.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  const lifted = active.locator('.hand-card.aiming');
  // Click the card's visible top: a locator click would scroll the hand,
  // which hangs past the window's edge, into view.
  const box = await bomb.boundingBox();
  if (!box) {
    throw new Error('no bomb box');
  }
  const clickBomb = () => active.mouse.click(box.x + box.width / 2, box.y + 30);

  // What was sent and at which step, to see that nothing goes out.
  const sent = () =>
    active.evaluate(() => {
      const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
      return { sent: s.sent, step: s.game?.step };
    });
  const before = await sent();

  // A right click cancels aiming; nothing is sent.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await active.mouse.click(640, 360, { button: 'right' });
  await expect(lifted).toHaveCount(0);
  await expect(monster).not.toHaveClass(/hl-aim/);
  await expect(active.getByRole('menu')).toHaveCount(0);
  expect(await sent()).toEqual(before);

  // A shop item is no target for a bomb: the reason shows beside the
  // click, and aiming goes on. Clicking the bomb again does nothing;
  // empty space cancels.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  const item = active.getByRole('group', { name: tr('board.shop') }).locator('.card').first();
  const at = await item.boundingBox();
  if (!at) {
    throw new Error('no shop item box');
  }
  await active.mouse.click(at.x + at.width / 2, at.y + at.height / 2);
  const why = active.locator('.aim-why');
  await expect(why).toHaveText(tr('refuse.pick.monster_or_player'));
  await shot(active, 'aim-wrong');
  const note = await why.locator('div').boundingBox();
  if (!note) {
    throw new Error('no reason box');
  }
  // The note sits just above the click, centered on it.
  expect(Math.abs(note.x + note.width / 2 - (at.x + at.width / 2)) < 3 && note.y + note.height < at.y + at.height / 2).toBe(true);
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await active.getByText(tr('stack.empty')).click();
  await expect(lifted).toHaveCount(0);
  await expect(why).toHaveCount(0);
  await expect(active.getByRole('menu')).toHaveCount(0);
  expect(await sent()).toEqual(before);

  // One click starts aiming: the card stays lifted, the targets glow, and
  // an arrow runs from the card's middle to the pointer.
  await clickBomb();
  await expect(lifted).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  // A player glows by the character card, not by the area (15.33).
  const opponent = active.locator('.player.opponent');
  await expect(opponent.locator('.opp-character .card')).toHaveClass(/hl-aim/);
  await expect(opponent).not.toHaveClass(/hl-aim/);
  const id = Number((await monster.getAttribute('data-zone'))?.replace('obj-', ''));
  const to = await monster.boundingBox();
  if (!to) {
    throw new Error('no monster box');
  }
  const tip = { x: to.x + to.width / 2, y: to.y + to.height / 2 };
  // Off a target the head is plain; on one, a ring marks it.
  const arrow = active.locator('.aim-arrow');
  await active.mouse.move(to.x - 200, tip.y, { steps: 12 });
  await expect(arrow).toBeVisible();
  await expect(arrow).not.toHaveClass(/\bon\b/);
  await expect(active.locator('.aim-dash').first()).toBeAttached();
  await shot(active, 'aim-off');
  await active.mouse.move(tip.x, tip.y, { steps: 12 });
  await expect(arrow).toHaveClass(/\bon\b/);
  await shot(active, 'aim');
  const geometry = await active.evaluate(() => {
    const box = (s: string) => document.querySelector(s)?.getBoundingClientRect();
    return { head: box('.aim-head'), ring: box('.aim-ring > div'), chain: box('.aim-chain'), card: box('.hand-card.aiming .card') };
  });
  const { head, ring, chain, card } = geometry;
  if (!head || !ring || !chain || !card) {
    throw new Error('no arrow');
  }
  // The head ends at the pointer, where the ring is centered; the chain
  // starts in the card's middle.
  const near = (a: number, b: number) => Math.abs(a - b) < 3;
  expect(tip.x >= head.left - 1 && tip.x <= head.right + 1 && tip.y >= head.top - 1 && tip.y <= head.bottom + 1).toBe(true);
  expect(near(ring.left + ring.width / 2, tip.x) && near(ring.top + ring.height / 2, tip.y)).toBe(true);
  const middle = { x: card.left + card.width / 2, y: card.top + card.height / 2 };
  expect(near(chain.left, middle.x) && near(chain.top, middle.y)).toBe(true);

  // The second click picks the target.
  await active.mouse.click(tip.x, tip.y);
  await expect(active.locator('.aim-arrow')).toHaveCount(0);

  // The click picked the monster: no dialog asked, the bomb aims at it.
  await expect
    .poll(() =>
      active.evaluate((obj) => {
        const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
        return (store.getState().game?.view.stack ?? []).some((it) => it.card === 'bomb' && it.targets?.[0]?.object === obj);
      }, id),
    )
    .toBe(true);
  const damaged = () => active.evaluate((obj) => (window as unknown as { damaged: number[] }).damaged.includes(obj), id);
  await settle(pages, click, damaged);
  expect(await damaged()).toBe(true);
});

// stackCards lists the cards on the page's stack.
function stackCards(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    return (store.getState().game?.view.stack ?? []).map((it) => it.card ?? '');
  });
}

// dragToBoard drags a hand card onto the board, by its top left: the
// next card in the fan covers its middle, and a tilted card misses
// its box's very corner. The card is measured once it stands still: a
// card just put back slides down into the fan. The first move stays
// on the card, so the card sees it and takes the pointer.
async function dragToBoard(page: Page, card: string) {
  const loc = page.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: card }).first();
  let from = await loc.boundingBox();
  for (let still = false; !still; ) {
    await page.waitForTimeout(100);
    const now = await loc.boundingBox();
    still = !!from && !!now && from.x === now.x && from.y === now.y;
    from = now;
  }
  const to = await page.getByRole('region', { name: tr('board.label') }).boundingBox();
  if (!from || !to) {
    throw new Error('no boxes');
  }
  const x = from.x + 40;
  const y = from.y + 40;
  await page.mouse.move(x, y);
  await page.mouse.down();
  await page.mouse.move(x, y - 12, { steps: 2 });
  await page.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
  await page.mouse.up();
}

test('a loot card with no target dragged onto the board plays', async ({ browser }) => {
  const pages = await startGame(browser, ['Dot', 'Eve'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Dot, Eve', 'a_penny');
  const active = await activePage(pages, click);
  const hand = active.getByRole('group', { name: tr('table.your_hand') });
  // Any role: once played, other pennies are no longer buttons.
  const penny = hand.getByLabel(cardName('a_penny'), { exact: true });
  const before = await penny.count();

  // Played loot is kept as it shows: the stack may resolve at once.
  await active.evaluate(() => {
    const w = window as unknown as { fourSouls: { store: Store }; played: string[] };
    w.played = [];
    let step = -1;
    w.fourSouls.store.subscribe((s) => {
      if (s.game && s.game.step !== step) {
        step = s.game.step;
        w.played.push(...(s.game.events ?? []).filter((e) => e.kind === 'loot_played').map((e) => e.card ?? ''));
      }
    });
  });

  await dragToBoard(active, cardName('a_penny'));
  await expect.poll(() => active.evaluate(() => (window as unknown as { played: string[] }).played)).toContain('a_penny');
  await expect(penny).toHaveCount(before - 1);
  // No arrow, no dialog: the card went straight to the stack.
  await expect(active.locator('.aim-arrow')).toHaveCount(0);
  await expect(active.getByRole('dialog')).toHaveCount(0);
});

// tiltOf is the held card's turn now, in degrees: its inner box's
// rotateX (pitch) and rotateY (roll); 0 where it has none.
async function tiltOf(page: Page): Promise<{ pitch: number; roll: number }> {
  const style = await page.locator('.hand-card.dragging .drag-tilt').evaluate((el) => (el as HTMLElement).style.transform);
  const deg = (axis: string) => Number(new RegExp(`rotate${axis}\\((-?[\\d.e-]+)deg\\)`).exec(style)?.[1] ?? 0);
  return { pitch: deg('X'), roll: deg('Y') };
}

test('a dragged card tilts by its motion and lies flat when still (20.6)', async ({ browser }) => {
  // Animations on: speed "off" keeps the card flat.
  const pages = await startGame(browser, ['Til', 'Ura'], { url: appURL() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Til, Ura', 'a_penny');
  const active = await activePage(pages, click);
  // The loops' files do not load here: they stay silent, with no
  // complaint from the audio context.
  const audio: string[] = [];
  active.on('console', (m) => {
    if (/audio/i.test(m.text())) {
      audio.push(m.text());
    }
  });
  const penny = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('a_penny') }).first();
  await expect(penny).toBeVisible();
  const from = await penny.boundingBox();
  if (!from) {
    throw new Error('no box');
  }
  const x = from.x + 40;
  const y = from.y + 40;
  await active.mouse.move(x, y);
  await active.mouse.down();
  // Up and to the right, fast: the top edge and the right edge lean away.
  await active.mouse.move(x + 10, y - 12, { steps: 2 });
  await active.mouse.move(x + 160, y - 200, { steps: 6 });
  await expect.poll(async () => (await tiltOf(active)).pitch).toBeGreaterThan(5);
  const tilt = await tiltOf(active);
  expect(tilt.roll).toBeGreaterThan(5);
  expect(tilt.pitch).toBeLessThanOrEqual(40);
  expect(tilt.roll).toBeLessThanOrEqual(45);
  await shot(active, 'drag-tilt');
  // Held still, it springs back to flat.
  await expect.poll(async () => Math.abs((await tiltOf(active)).roll)).toBeLessThan(1);
  // Dropped back on the hand, it is not played.
  await active.mouse.move(x, y, { steps: 4 });
  await active.mouse.up();
  await expect(active.locator('.hand-card.dragging')).toHaveCount(0);
  await expect(penny).toBeVisible();
  // Reduced motion: no tilt at all.
  await active.emulateMedia({ reducedMotion: 'reduce' });
  await active.mouse.move(x, y);
  await active.mouse.down();
  await active.mouse.move(x + 10, y - 12, { steps: 2 });
  await active.mouse.move(x + 160, y - 200, { steps: 6 });
  await active.waitForTimeout(100);
  expect(await tiltOf(active)).toEqual({ pitch: 0, roll: 0 });
  // Dropped back on the hand: nothing is played.
  await active.mouse.move(x, y, { steps: 4 });
  await active.mouse.up();
  await expect(active.locator('.hand-card.dragging')).toHaveCount(0);
  expect(audio).toEqual([]);
});

test('a card placed first aims after its mode; a cancelled aim leaves it in hand', async ({ browser }) => {
  // Dagaz is banned from the deal: the test gives one and counts it; a
  // dealt Dagaz made two in the hand (19.20).
  const pages = await startGame(browser, ['Fay', 'Gus'], { url: fast(), bans: ['dagaz'] });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Fay, Gus', 'dagaz');
  const active = await activePage(pages, click);
  const hand = active.getByRole('group', { name: tr('table.your_hand') });
  const dagaz = hand.getByRole('button', { name: cardName('dagaz') });
  const lifted = active.locator('.hand-card.aiming');
  const opponent = active.locator('.player.opponent .opp-character .card');
  const mode = 'Choose a player. Prevent the next 1 damage they would take this turn.';
  // place drags Dagaz onto the board and picks the mode that aims.
  const place = async () => {
    await dragToBoard(active, cardName('dagaz'));
    const ask = active.getByRole('dialog', { name: tr('choice.title.mode') });
    await expect(ask).toBeVisible();
    await shot(active, 'aim-mode');
    await ask.getByRole('button', { name: mode }).click();
    // The target prompt is an aim from the card, not a dialog.
    await expect(lifted).toHaveCount(1);
    await expect(opponent).toHaveClass(/hl-aim/);
    await expect(active.getByRole('dialog')).toHaveCount(0);
    // Players are aimed at by the character card; areas stay plain.
    await expect(active.locator('.player.hl-pick')).toHaveCount(0);
    await expect(active.getByRole('button', { name: new RegExp(`^${trRe('table.choose').source}`) })).toHaveCount(0);
  };

  await place();
  const to = await opponent.boundingBox();
  if (!to) {
    throw new Error('no character box');
  }
  await active.mouse.move(to.x + to.width / 2, to.y + to.height / 2, { steps: 12 });
  await expect(active.locator('.aim-arrow')).toHaveClass(/\bon\b/);
  await shot(active, 'aim-placed');

  // A right click cancels the play: the card stays in hand, the stack
  // stays empty, and nothing aims any more.
  await active.mouse.click(640, 360, { button: 'right' });
  await expect(lifted).toHaveCount(0);
  await expect(opponent).not.toHaveClass(/hl-aim/);
  await expect(dagaz).toHaveCount(1);
  await expect(dagaz).toHaveClass(/hl-act/);
  expect(await stackCards(active)).not.toContain('dagaz');

  // Placed again, a click on the opponent's character plays it there.
  await place();
  await active.mouse.click(to.x + to.width / 2, to.y + to.height / 2);
  await expect.poll(() => stackCards(active)).toContain('dagaz');
  await expect(dagaz).toHaveCount(0);
});

// A "choose one-" card asks its mode like The Bone asks its ability
// (15.58): the card large, a line per mode, the one with no target
// greyed. Cancel, a click outside, and "View board" work as on any
// choice; the picked mode then aims.
test('a "choose one-" card asks its mode with the card large (15.58)', async ({ browser }) => {
  // Dagaz is banned from the deal: the test gives one and counts it.
  const pages = await startGame(browser, ['Ida', 'Jon'], { url: fast(), bans: ['dagaz'] });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Ida, Jon', 'dagaz');
  const active = await activePage(pages, click);
  const dagaz = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('dagaz') });
  const ask = active.getByRole('dialog', { name: tr('choice.title.mode') });
  const open = async () => {
    await dragToBoard(active, cardName('dagaz'));
    await expect(ask).toBeVisible();
  };

  await open();
  await expect(ask.getByRole('heading', { name: cardName('dagaz') })).toBeVisible();
  await expect(ask.locator('.card.size-zoom')).toBeVisible();
  // No curse is in play at the start: that mode stays, greyed.
  const curse = ask.getByRole('button', { name: 'Destroy a curse.' });
  await expect(curse).toBeDisabled();
  await expect(curse).toHaveAttribute('title', tr('choice.no_target'));
  const protect = ask.getByRole('button', { name: /^Choose a player/ });
  await expect(protect).toBeEnabled();
  await shot(active, 'mode-dialog');

  // View board hides it; the bar brings it back.
  await ask.getByRole('button', { name: tr('choice.view_board') }).click();
  await expect(ask).toBeHidden();
  await active.getByRole('group', { name: tr('choice.viewing') }).getByRole('button', { name: tr('choice.back') }).click();
  await expect(ask).toBeVisible();

  // Cancel sends the prompt's cancel: Dagaz stays in hand.
  await ask.getByRole('button', { name: tr('actions.cancel') }).click();
  await expect(ask).toBeHidden();
  await expect(dagaz).toHaveClass(/hl-act/);
  // So does a click outside the dialog.
  await open();
  await active.mouse.click(8, 360);
  await expect(ask).toBeHidden();
  await expect(dagaz).toHaveClass(/hl-act/);

  // The picked mode is sent: its target prompt aims from the card.
  await open();
  await protect.click();
  await expect(ask).toBeHidden();
  await expect(active.locator('.hand-card.aiming')).toHaveCount(1);
});

// Presence (15.36): the other player sees the aim live, from a back
// lifted out of the aiming player's hand to the red-marked target, and
// the card the active player hovers. A waiting player shares no hover.
test('the other player sees the arrow and the red mark', async ({ browser }) => {
  const pages = await startGame(browser, ['Hal', 'Ivy'], { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Hal, Ivy');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no waiting player');
  }
  const bomb = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('bomb') }).first();
  const monster = active.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  const item = active.getByRole('group', { name: tr('board.shop') }).locator('.card').first();
  const arrow = other.locator('.aim-arrow.others');
  const lifted = other.locator('.hidden-back[data-presence~="from"]');
  // The same card on the other screen, by its zone.
  const there = async (l: Locator) => other.locator(`[data-zone="${(await l.getAttribute('data-zone')) ?? ''}"]`);
  const middle = async (l: Locator) => {
    const b = await l.boundingBox();
    if (!b) {
      throw new Error('no box');
    }
    return { x: b.x + b.width / 2, y: b.y + b.height / 2 };
  };
  // Empty board below the stack panel: no card, pile, or row.
  const space = { x: 300, y: 420 };
  const b = await bomb.boundingBox();
  if (!b) {
    throw new Error('no bomb box');
  }

  // Aiming: the arrow runs from a lifted back to the marked monster.
  await active.mouse.click(b.x + b.width / 2, b.y + 30);
  const m = await middle(monster);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await expect(arrow).toHaveCount(1);
  await expect(arrow).toHaveClass(/\bon\b/);
  await expect(lifted).toHaveCount(1);
  await expect(await there(monster)).toHaveAttribute('data-presence', 'at');
  // The head ends on the marked monster, not on a card passed on the way.
  // Its tip is where the ring would be; others' arrows draw no ring (20.23).
  await expect(other.locator('.aim-arrow.others .aim-ring > div')).toHaveCount(0);
  const ringAt = () =>
    other.evaluate(() => {
      const r = document.querySelector('.aim-arrow.others .aim-ring')?.getBoundingClientRect();
      return { x: r?.left ?? NaN, y: r?.top ?? NaN };
    });
  const near = async (to: { x: number; y: number }) => {
    const r = await ringAt();
    return Math.hypot(r.x - to.x, r.y - to.y);
  };
  const mm = await middle(await there(monster));
  await expect.poll(() => near(mm)).toBeLessThan(3);
  await shot(other, 'presence-aim');

  // Over empty space the arrow stays, with no mark.
  await active.mouse.move(space.x, space.y, { steps: 4 });
  await expect(arrow).not.toHaveClass(/\bon\b/);
  await expect(await there(monster)).not.toHaveAttribute('data-presence');
  // Both boards are drawn alike, so the head ends at the same point.
  await expect.poll(() => near(space)).toBeLessThan(3);
  await shot(other, 'presence-aim-off');

  // Cancelled, the arrow and the lifted back go.
  await active.mouse.click(space.x, space.y, { button: 'right' });
  await expect(arrow).toHaveCount(0);
  await expect(lifted).toHaveCount(0);

  // On the own turn, a hovered card glows for the other, in the
  // hovering player's color, not the red of an aim (20.13).
  const i = await middle(item);
  await active.mouse.move(i.x, i.y, { steps: 4 });
  await expect(await there(item)).toHaveAttribute('data-presence', 'hover');
  const color = await (await there(item)).evaluate((el) => (el as HTMLElement).style.getPropertyValue('--presence-color'));
  expect(color).toMatch(/^#[0-9a-f]{6}$/);
  await shot(other, 'presence-hover');
  await active.mouse.move(space.x, space.y, { steps: 4 });
  await expect(await there(item)).not.toHaveAttribute('data-presence');

  // The waiting player's hover reaches nobody.
  const o = await middle(other.getByRole('group', { name: tr('board.shop') }).locator('.card').first());
  await other.mouse.move(o.x, o.y, { steps: 4 });
  await active.waitForTimeout(500);
  const shown = await active.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: { getState: () => { presence: object } } } }).fourSouls.store.getState();
    return Object.keys(s.presence).length;
  });
  expect(shown).toBe(0);
  await expect(active.locator('[data-presence]')).toHaveCount(0);
});

// Others' hands live (20.13, ADR 027 "Other players live"): the hand
// card the active player hovers lifts a little on the other page, still
// a back; a card they drag floats out of the hand and sways until it is
// played or put back. Animations on, as players see it.
test('the other player sees a hovered hand card lift and a dragged one float', async ({ browser }) => {
  const pages = await startGame(browser, ['Lia', 'Max'], { url: appURL() });
  await activePage(pages, click);
  await saveAndContinue(pages, 'Lia, Max', 'a_penny');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no waiting player');
  }
  const penny = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('a_penny') }).first();
  await expect(penny).toBeVisible();
  const backs = other.locator('.hidden-back');
  // A hovered or held back is still a back: no card object in it.
  const hovered = other.locator('.hidden-back[data-presence~="hover"]');
  const held = other.locator('.hidden-back[data-presence~="held"]');
  // How far a back sits off its place in the fan, in px, either way.
  const lift = (l: Locator) => l.evaluate((el) => Math.abs(Number.parseFloat(getComputedStyle(el).translate.split(' ')[1] ?? '0')));
  await expect(backs.first()).toBeVisible();
  await expect(other.locator('.hidden-back[data-presence]')).toHaveCount(0);

  // Hovered: one back lifts a little toward the table; it stays a back.
  const from = await penny.boundingBox();
  if (!from) {
    throw new Error('no box');
  }
  const x = from.x + 40;
  const y = from.y + 40;
  await active.mouse.move(x, y, { steps: 4 });
  await expect(hovered).toHaveCount(1);
  await expect.poll(() => lift(hovered)).toBe(14);
  await expect(hovered.locator('[data-zone^="obj-"]')).toHaveCount(0);
  await shot(other, 'presence-hand-hover');

  // Dragged: the back floats further out, grows, and sways.
  await active.mouse.down();
  await active.mouse.move(x, y - 12, { steps: 2 });
  await active.mouse.move(x + 40, y - 200, { steps: 6 });
  await expect(held).toHaveCount(1);
  await expect(hovered).toHaveCount(0);
  await expect.poll(() => lift(held)).toBe(56);
  expect(await held.evaluate((el) => getComputedStyle(el).animationName)).toBe('others-sway');
  // It sways: its turn changes while the card is held.
  const turn = () => held.evaluate((el) => getComputedStyle(el).rotate);
  const first = await turn();
  await expect.poll(turn, { timeout: 5000 }).not.toBe(first);
  await shot(other, 'presence-held');

  // Put back on the hand: it goes back into the fan.
  await active.mouse.move(x, y, { steps: 4 });
  await active.mouse.up();
  await expect(held).toHaveCount(0);
});
