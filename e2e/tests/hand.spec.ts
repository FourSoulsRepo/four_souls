// Playing a card from the hand (20.20, ADR 027 "Playing from the hand"):
// a click picks a card up and a click on the board plays it; no menu.
// A held card grows, glows gold over the board, and flies home when put
// back. A card that needs a target stops at the board's edge and aims.
// A played card waits where it was let go until the server answers,
// then shrinks into its stack row; refused, it goes home. Each test
// gets its card from a saved game, so none depends on the deal
// (CLAUDE.md 5.16).
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { saveAndContinue } from './saves';
import { appURL, serverPort } from './server';
import { activePage, click, playSpot, tap } from './turn';

const fast = () => `${appURL()}&anim=off`;

test.setTimeout(120_000);

// names are the players' names; a repeat (--repeat-each) gets its own,
// as the saved games list keeps the earlier repeats' games.
function names(base: string[]): string[] {
  const r = test.info().repeatEachIndex;
  return base.map((n) => (r > 0 ? `${n}${String(r)}` : n));
}

interface Store {
  getState: () => {
    game: { step: number; view: { stack: { card?: string; source: number; targets?: { object?: number }[] }[] | null } } | null;
    sent: number | null;
    held: number | null;
  };
}

// handCard is the first hand card of that name.
function handCard(page: Page, card: string): Locator {
  return page.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName(card) }).first();
}

// topOf is a point on the card's visible top part: the hand hangs past
// the window's edge, and a locator click would scroll it into view.
async function topOf(card: Locator): Promise<{ x: number; y: number }> {
  let from = await card.boundingBox();
  // Measured once it stands still: a card just put back slides home.
  for (let still = false; !still; ) {
    await card.page().waitForTimeout(100);
    const now = await card.boundingBox();
    still = !!from && !!now && Math.abs(from.x - now.x) < 0.5 && Math.abs(from.y - now.y) < 0.5;
    from = now;
  }
  if (!from) {
    throw new Error('no card box');
  }
  return { x: from.x + from.width / 2, y: from.y + 30 };
}

async function middle(l: Locator): Promise<{ x: number; y: number }> {
  const b = await l.boundingBox();
  if (!b) {
    throw new Error('no box');
  }
  return { x: b.x + b.width / 2, y: b.y + b.height / 2 };
}

// stack lists the cards on the page's stack.
function stack(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    return (store.getState().game?.view.stack ?? []).map((it) => it.card ?? '');
  });
}

// sent is what the page sent last and its step now: nothing goes out
// while it stays the same.
function sent(page: Page): Promise<{ sent: number | null; step?: number }> {
  return page.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    return { sent: s.sent, step: s.game?.step };
  });
}

// spy records every sound played and every WAAPI animation started
// from now on, in window.played and window.anims; each still runs.
async function spy(page: Page): Promise<void> {
  await page.evaluate(() => {
    const w = window as unknown as { played: string[]; anims: { cls: string; ms: number; easing: string }[] };
    w.played = [];
    w.anims = [];
    const play = HTMLMediaElement.prototype.play;
    HTMLMediaElement.prototype.play = function (this: HTMLMediaElement) {
      w.played.push(new URL(this.src).pathname);
      return play.call(this);
    };
    const animate = Element.prototype.animate;
    Element.prototype.animate = function (this: Element, keys: Keyframe[] | PropertyIndexedKeyframes | null, o?: number | KeyframeAnimationOptions) {
      const t = typeof o === 'number' ? { duration: o } : (o ?? {});
      w.anims.push({ cls: typeof this.className === 'string' ? this.className : '', ms: Number(t.duration ?? 0), easing: t.easing ?? '' });
      return animate.call(this, keys, o);
    };
  });
}

function played(page: Page, sound: string): Promise<number> {
  return page.evaluate((s) => (window as unknown as { played: string[] }).played.filter((p) => p === `/cardimg/assets/sounds/${s}.mp3`).length, sound);
}

// anims are the WAAPI animations started on elements whose class has
// cls, by their time in ms.
function anims(page: Page, cls: string): Promise<number[]> {
  return page.evaluate((c) => (window as unknown as { anims: { cls: string; ms: number }[] }).anims.filter((a) => a.cls.split(' ').includes(c)).map((a) => a.ms), cls);
}

const expoOut = 'cubic-bezier(0.16, 1, 0.3, 1)';

test('a click picks a hand card up; the board plays it, the hand puts it back (20.20)', async ({ browser }) => {
  const who = names(['Ann', 'Bo']);
  const pages = await startGame(browser, who, { url: appURL() });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'a_penny');
  const active = await activePage(pages, click);
  await spy(active);
  const penny = handCard(active, 'a_penny');
  const held = active.locator('.hand-card.held');
  const board = playSpot(active);
  const before = await sent(active);
  const at = await topOf(penny);

  // A click picks it up: no menu; it grows to its held size, 0.2 s,
  // expo-out, and sounds.
  await active.mouse.click(at.x, at.y);
  await expect(held).toHaveCount(1);
  await expect(active.getByRole('menu')).toHaveCount(0);
  await expect.poll(() => played(active, 'card-grab')).toBe(1);
  expect(await anims(active, 'hand-card')).toContain(200);
  // It follows the pointer, button up; over the board its glow is gold.
  const m = await middle(board);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await expect(held).toHaveClass(/\bover\b/);
  const box = await held.boundingBox();
  if (!box) {
    throw new Error('no held box');
  }
  expect(Math.abs(box.x + box.width / 2 - m.x)).toBeLessThan(box.width);
  // Off the board it is green again; a click on the hand puts it back,
  // flying home in 0.5 s. Nothing is sent.
  await active.mouse.move(at.x, at.y, { steps: 8 });
  await expect(held).not.toHaveClass(/\bover\b/);
  await active.mouse.click(at.x, at.y);
  await expect(held).toHaveCount(0);
  await expect.poll(() => anims(active, 'hand-card')).toContain(500);
  expect(await sent(active)).toEqual(before);

  // A right click and Esc put it back too.
  await active.mouse.click(at.x, at.y);
  await expect(held).toHaveCount(1);
  await active.mouse.click(m.x, m.y, { button: 'right' });
  await expect(held).toHaveCount(0);
  await active.mouse.click(at.x, at.y);
  await expect(held).toHaveCount(1);
  await active.keyboard.press('Escape');
  await expect(held).toHaveCount(0);
  expect(await sent(active)).toEqual(before);

  // Picked up again, a click on the board plays it.
  await active.mouse.click(at.x, at.y);
  await expect(held).toHaveCount(1);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await active.evaluate(() => {
    (window as unknown as { anims: unknown[] }).anims.length = 0;
  });
  await active.mouse.click(m.x, m.y);
  await expect.poll(() => stack(active)).toContain('a_penny');
  // It shrinks into its stack row in 0.4 s, and the rest of the hand
  // closes up in 0.5 s.
  await expect.poll(() => anims(active, 'hand-drop')).toContain(400);
  await expect(active.locator('.hand-drop')).toHaveCount(0);
  if ((await active.locator('.own .hand-card').count()) > 0) {
    expect(await anims(active, 'hand-card')).toContain(500);
  }
  // Every motion of the hand eases out, expo.
  const easings = await active.evaluate(() =>
    (window as unknown as { anims: { cls: string; easing: string }[] }).anims.filter((a) => /\bhand-(card|drop)\b/.test(a.cls)).map((a) => a.easing),
  );
  expect(easings.every((e) => e === expoOut)).toBe(true);
});

test('speed "off": a held card neither grows nor flies (20.20)', async ({ browser }) => {
  const who = names(['Cy', 'Di']);
  const pages = await startGame(browser, who, { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'a_penny');
  const active = await activePage(pages, click);
  await spy(active);
  const at = await topOf(handCard(active, 'a_penny'));
  const held = active.locator('.hand-card.held');
  await active.mouse.click(at.x, at.y);
  await expect(held).toHaveCount(1);
  await active.keyboard.press('Escape');
  await expect(held).toHaveCount(0);
  await active.mouse.click(at.x, at.y);
  const m = await middle(playSpot(active));
  await active.mouse.click(m.x, m.y);
  await expect.poll(() => stack(active)).toContain('a_penny');
  expect(await anims(active, 'hand-card')).toEqual([]);
  expect(await anims(active, 'hand-drop')).toEqual([]);
});

test('a card that needs a target stops at the board edge and aims (20.20)', async ({ browser }) => {
  const who = names(['Ed', 'Flo']);
  const pages = await startGame(browser, who, { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'bomb');
  const active = await activePage(pages, click);
  await spy(active);
  const bomb = handCard(active, 'bomb');
  const monster = active.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  const stopped = active.locator('.hand-card.stopped');
  const dragging = active.locator('.hand-card.dragging');
  const arrow = active.locator('.aim-arrow');
  const board = await active.getByRole('region', { name: tr('board.label') }).boundingBox();
  if (!board) {
    throw new Error('no board');
  }
  const at = await topOf(bomb);
  const below = { x: at.x, y: board.y + board.height + 20 };
  const inside = { x: at.x, y: board.y + board.height - 40 };
  const m = await middle(monster);

  // Dragged up, it stops on reaching the board: the aim takes over.
  await active.mouse.move(at.x, at.y);
  await active.mouse.down();
  await active.mouse.move(at.x, at.y - 12, { steps: 2 });
  await active.mouse.move(below.x, below.y, { steps: 6 });
  await expect(dragging).toHaveCount(1);
  await active.mouse.move(inside.x, inside.y, { steps: 4 });
  await expect(stopped).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  const stop = await stopped.boundingBox();
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await expect(arrow).toHaveClass(/\bon\b/);
  // The card stays where it stopped; the arrow starts from it.
  expect(await stopped.boundingBox()).toEqual(stop);
  const chain = await active.locator('.aim-chain').boundingBox();
  if (!stop || !chain) {
    throw new Error('no boxes');
  }
  expect(Math.abs(chain.x - (stop.x + stop.width / 2)) < 3 && Math.abs(chain.y - (stop.y + stop.height / 2)) < 3).toBe(true);

  // Back below the board, the aim ends and the drag goes on.
  await active.mouse.move(below.x, below.y, { steps: 8 });
  await expect(stopped).toHaveCount(0);
  await expect(dragging).toHaveCount(1);
  await expect(monster).not.toHaveClass(/hl-aim/);
  expect(await played(active, 'aim-cancel')).toBe(0);

  // Up again and let go on empty board: the aim stays, the card waits
  // at the edge; a right click cancels it, with its sound, and the card
  // goes home. Nothing is sent.
  const before = await sent(active);
  await active.mouse.move(inside.x, inside.y, { steps: 4 });
  await expect(stopped).toHaveCount(1);
  await active.mouse.move(300, 420, { steps: 4 });
  await active.mouse.up();
  await expect(stopped).toHaveCount(1);
  await expect(monster).toHaveClass(/hl-aim/);
  await active.mouse.click(300, 420, { button: 'right' });
  await expect(stopped).toHaveCount(0);
  await expect(monster).not.toHaveClass(/hl-aim/);
  await expect.poll(() => played(active, 'aim-cancel')).toBe(1);
  expect(await sent(active)).toEqual(before);

  // Dragged in again and let go on the monster: the bomb plays at it.
  const id = Number((await monster.getAttribute('data-zone'))?.replace('obj-', ''));
  const again = await topOf(bomb);
  await active.mouse.move(again.x, again.y);
  await active.mouse.down();
  await active.mouse.move(again.x, again.y - 12, { steps: 2 });
  await active.mouse.move(below.x, below.y, { steps: 6 });
  await active.mouse.move(inside.x, inside.y, { steps: 4 });
  await expect(stopped).toHaveCount(1);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await active.mouse.up();
  await expect
    .poll(() =>
      active.evaluate((obj) => {
        const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
        return (store.getState().game?.view.stack ?? []).some((it) => it.card === 'bomb' && it.targets?.[0]?.object === obj);
      }, id),
    )
    .toBe(true);
  // No dialog asked; the click that ended the drag picked nothing more.
  await expect(active.getByRole('dialog')).toHaveCount(0);
  await expect(arrow).toHaveCount(0);
});

// holdAnswers keeps the server's messages from the page while hold()
// says so, then lets them through in order; spoil() makes the page's
// next intents point at no card, so the server refuses them.
function wire(hold: () => boolean, spoil: () => boolean) {
  return (page: Page) =>
    page.routeWebSocket(`ws://127.0.0.1:${String(serverPort())}/ws`, (ws) => {
      const server = ws.connectToServer();
      const kept: (string | Buffer)[] = [];
      ws.onMessage((msg) => {
        const env = JSON.parse(String(msg)) as { type: string; data?: { intent?: { objects?: number[] } } };
        if (spoil() && env.type === 'intent' && env.data?.intent) {
          env.data.intent.objects = [999_999];
          server.send(JSON.stringify(env));
          return;
        }
        server.send(msg);
      });
      const flush = () => {
        while (!hold() && kept.length > 0) {
          const m = kept.shift();
          if (m !== undefined) {
            ws.send(m);
          }
        }
      };
      const timer = setInterval(flush, 20);
      server.onClose(() => {
        clearInterval(timer);
      });
      server.onMessage((msg) => {
        if (hold() || kept.length > 0) {
          kept.push(msg);
          return;
        }
        ws.send(msg);
      });
    });
}

test('a played card waits where it was let go until the server answers (20.20)', async ({ browser }) => {
  let holding = false;
  let spoiling = false;
  const who = names(['Gil', 'Hu']);
  const pages = await startGame(browser, who, {
    url: appURL(),
    open: wire(
      () => holding,
      () => spoiling,
    ),
  });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'a_penny');
  const active = await activePage(pages, click);
  await spy(active);
  const penny = handCard(active, 'a_penny');
  const waiting = active.locator('.hand-card.waiting');
  const m = await middle(playSpot(active));
  const at = await topOf(penny);

  // A slow answer: the card stays at the drop point, still in the hand.
  holding = true;
  await active.mouse.move(at.x, at.y);
  await active.mouse.down();
  await active.mouse.move(at.x, at.y - 12, { steps: 2 });
  await active.mouse.move(m.x, m.y, { steps: 10 });
  await active.mouse.up();
  await expect(waiting).toHaveCount(1);
  const dropped = await waiting.boundingBox();
  await active.waitForTimeout(1500);
  const still = await waiting.boundingBox();
  // Its middle stays put; its size may still have been growing at the drop.
  const mid = (b: { x: number; y: number; width: number; height: number } | null) => (b ? { x: b.x + b.width / 2, y: b.y + b.height / 2 } : { x: NaN, y: NaN });
  expect(Math.hypot(mid(dropped).x - mid(still).x, mid(dropped).y - mid(still).y)).toBeLessThan(1);
  expect(await stack(active)).not.toContain('a_penny');
  // The answer: it shrinks into its stack row, from where it waited.
  holding = false;
  await expect.poll(() => stack(active)).toContain('a_penny');
  await expect.poll(() => anims(active, 'hand-drop')).toContain(400);
  await expect(waiting).toHaveCount(0);
  await expect(active.locator('.hand-drop')).toHaveCount(0);
});

test('a refused play goes home and says why (20.20)', async ({ browser }) => {
  let spoiling = false;
  const who = names(['Ike', 'Jo']);
  const pages = await startGame(browser, who, {
    url: appURL(),
    open: wire(
      () => false,
      () => spoiling,
    ),
  });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'a_penny');
  const active = await activePage(pages, click);
  await spy(active);
  const penny = handCard(active, 'a_penny');
  const at = await topOf(penny);
  const m = await middle(playSpot(active));
  spoiling = true;
  await active.mouse.click(at.x, at.y);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await active.mouse.click(m.x, m.y);
  // The server refuses: the card flies home, the reason shows, and the
  // refusal sounds.
  await expect.poll(() => played(active, 'play-refused')).toBe(1);
  await expect(active.locator('.hand-card.waiting')).toHaveCount(0);
  await expect.poll(() => anims(active, 'hand-card')).toContain(500);
  await expect(penny).toBeVisible();
  expect(await stack(active)).not.toContain('a_penny');
  spoiling = false;
});

test('touch: a tap lifts a hand card; a tap on the board plays it (20.20)', async ({ browser }) => {
  const who = names(['Kai', 'Lu']);
  const pages = await startGame(browser, who, { url: fast(), touch: true });
  await activePage(pages, tap);
  await saveAndContinue(pages, who.join(', '), 'a_penny');
  const active = await activePage(pages, tap);
  const penny = handCard(active, 'a_penny');
  const lifted = active.locator('.hand-card.held.lifted');
  const at = await topOf(penny);
  const before = await sent(active);
  // A tap lifts it in the fan; a tap on it again puts it back.
  await active.touchscreen.tap(at.x, at.y);
  await expect(lifted).toHaveCount(1);
  await expect(active.getByRole('menu')).toHaveCount(0);
  const up = await topOf(lifted);
  await active.touchscreen.tap(up.x, up.y);
  await expect(lifted).toHaveCount(0);
  expect(await sent(active)).toEqual(before);
  // Lifted again, a tap on the board plays it.
  await active.touchscreen.tap(at.x, at.y);
  await expect(lifted).toHaveCount(1);
  await tap(active, playSpot(active));
  await expect.poll(() => stack(active)).toContain('a_penny');
});

test('a card with modes is picked up and played by clicks; its mode is asked (20.20)', async ({ browser }) => {
  // Dagaz is banned from the deal: the test gives one (19.20).
  const who = names(['Mo', 'Ned']);
  const pages = await startGame(browser, who, { url: fast(), bans: ['dagaz'] });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'dagaz');
  const active = await activePage(pages, click);
  const at = await topOf(handCard(active, 'dagaz'));
  await active.mouse.click(at.x, at.y);
  await expect(active.locator('.hand-card.held')).toHaveCount(1);
  const m = await middle(playSpot(active));
  await active.mouse.click(m.x, m.y);
  const ask = active.getByRole('dialog', { name: tr('choice.title.mode') });
  await expect(ask).toBeVisible();
  // The card is back in the hand while the mode is asked, not over it.
  await expect(active.locator('.hand-card.waiting, .hand-card.held')).toHaveCount(0);
  await ask.getByRole('button', { name: tr('actions.cancel') }).click();
  await expect(ask).toBeHidden();
  await expect(handCard(active, 'dagaz')).toHaveClass(/hl-act/);
});
