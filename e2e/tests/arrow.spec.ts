// The target arrow (20.23, ADR 027 "Target arrow", UC-13 2e): a curve
// of marching dashes bowing up the screen, drawn at once, the system
// pointer hidden; a hovered target glows brighter; a card zooms only
// once the arrow rests on it. An attack's pick aims it too, from the
// attacker, and the other player sees it. Another player's arrow glides
// out from its source and lingers a moment after they stop. The hand
// card and the monsters are set by the save, so nothing depends on the
// deal.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { saveAndContinue } from './saves';
import { appURL, serverPort } from './server';
import { shot } from './shot';
import { activePage, click, settle, tablePicks, turnAction } from './turn';

const fast = () => `${appURL()}&anim=off`;
const monsters = ['fatty', 'pale_fatty'];

test.setTimeout(150_000);

interface Box {
  x: number;
  y: number;
}

async function middle(l: Locator): Promise<Box> {
  const b = await l.boundingBox();
  if (!b) {
    throw new Error('no box');
  }
  return { x: b.x + b.width / 2, y: b.y + b.height / 2 };
}

// ends reads where an arrow starts and where its tip is: the chain's
// and the ring's places, both drawn at a point.
function ends(page: Page, sel: string): Promise<{ from: Box; to: Box }> {
  return page.evaluate((s) => {
    const at = (q: string) => {
      const r = document.querySelector(`${s} ${q}`)?.getBoundingClientRect();
      return { x: r?.left ?? NaN, y: r?.top ?? NaN };
    };
    return { from: at('.aim-chain'), to: at('.aim-ring') };
  }, sel);
}

// characterOf finds a seat's character card on a page.
async function characterOf(page: Page, seat: 'you' | 'other'): Promise<Locator> {
  const id = await page.evaluate((who) => {
    interface S {
      game: { view: { players: { character: { id: number } }[] } } | null;
      lobby: { table: { you: number } | null };
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const you = s.lobby.table?.you ?? 0;
    return s.game?.view.players[who === 'you' ? you : 1 - you]?.character.id ?? 0;
  }, seat);
  return page.locator(`[data-zone="obj-${String(id)}"]`).first();
}

test('the own arrow bows up, marches, hides the pointer, and zooms a card only at rest (20.23)', async ({ browser }) => {
  // Animations on: the march and the zoom's wait are what this checks.
  const who = ['Jo', 'Kit'];
  const pages = await startGame(browser, who, { url: appURL() });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'bomb', undefined, false, {}, undefined, monsters);
  const active = await activePage(pages, click);
  const bomb = active.getByRole('group', { name: tr('table.your_hand') }).getByRole('button', { name: cardName('bomb') }).first();
  const monster = active.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  const arrow = active.locator('.aim-arrow');
  const zoom = active.getByTestId('card-zoom');
  await active.waitForTimeout(800);
  const box = await bomb.boundingBox();
  if (!box) {
    throw new Error('no bomb');
  }

  // A click aims: the arrow shows at once, with no move, and the
  // system pointer hides.
  await active.mouse.click(box.x + box.width / 2, box.y + 30);
  await expect(monster).toHaveClass(/hl-aim/);
  await expect(arrow).toHaveCount(1);
  await expect(active.locator('html')).toHaveClass(/\baim-cursor\b/);

  // Sideways onto empty board below the stack panel, it bows up: its
  // dashes stand well above the straight line from the card.
  await active.mouse.move(300, 420, { steps: 8 });
  await expect(arrow).not.toHaveClass(/\bon\b/);
  const { from, to } = await ends(active, '.aim-arrow');
  const lift = await active.evaluate(
    ({ a, b }) => {
      let most = 0;
      for (const el of document.querySelectorAll('.aim-arrow .aim-dash > div')) {
        const r = el.getBoundingClientRect();
        const p = { x: r.left + r.width / 2, y: r.top + r.height / 2 };
        // How far above the straight line the dash stands.
        const line = a.y + ((p.x - a.x) / (b.x - a.x)) * (b.y - a.y);
        most = Math.max(most, line - p.y);
      }
      return most;
    },
    { a: from, b: to },
  );
  expect(lift).toBeGreaterThan(40);
  // The dashes march: one running loop per dash, all in step.
  const march = await active.evaluate(() => {
    const runs = document.getAnimations().filter((a) => a.effect instanceof KeyframeEffect && a.effect.target?.classList.contains('aim-dash'));
    return { n: runs.length, starts: new Set(runs.map((a) => a.startTime)).size, playing: runs.every((a) => a.playState === 'running') };
  });
  expect(march.n).toBeGreaterThan(2);
  expect(march.starts).toBe(1);
  expect(march.playing).toBe(true);
  await shot(active, 'arrow-bow');

  // Over a valid target: the ring, and its glow brightens.
  const m = await middle(monster);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  await expect(arrow).toHaveClass(/\bon\b/);
  await expect(monster).toHaveAttribute('data-aim-over', '');
  // No zoom while it has just come, even past the table's usual wait...
  await active.waitForTimeout(500);
  await expect(zoom).toHaveCount(0);
  // ...but once it rests still, the card zooms; the head stays in view.
  await expect(zoom).toBeVisible();
  await shot(active, 'arrow-rest-zoom');
  const head = await active.locator('.aim-head').boundingBox();
  const body = await active.locator('.card-zoom-body').boundingBox();
  if (!head || !body) {
    throw new Error('no head or zoom');
  }
  expect(m.x < body.x || m.x > body.x + body.width).toBe(true);
  // Any move closes it.
  await active.mouse.move(m.x + 3, m.y + 2);
  await expect(zoom).toHaveCount(0);

  // Off the target its glow goes back; a right click cancels, and the
  // pointer comes back.
  await active.mouse.move(300, 420, { steps: 4 });
  await expect(monster).not.toHaveAttribute('data-aim-over');
  await active.mouse.click(300, 420, { button: 'right' });
  await expect(arrow).toHaveCount(0);
  await expect(active.locator('html')).not.toHaveClass(/\baim-cursor\b/);
  await expect(active.locator('[data-aim-over]')).toHaveCount(0);
});

test("an attack's pick aims an arrow from the character; the other player sees it (20.23)", async ({ browser }) => {
  const who = ['Lu', 'Max'];
  const pages = await startGame(browser, who, { url: fast() });
  await activePage(pages, click);
  await saveAndContinue(pages, who.join(', '), 'bomb', undefined, false, {}, undefined, monsters);
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  await turnAction(active, 'Attack').click();
  await settle(pages, click, async () => tablePicks(active).first().isVisible());
  const monster = active.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  await expect(monster).toHaveClass(/hl-pick/);
  const zone = (await monster.getAttribute('data-zone')) ?? '';

  // The arrow runs from the own character to the pointer; on a monster
  // it is on, and the monster's blue glow brightens.
  const m = await middle(monster);
  await active.mouse.move(m.x, m.y, { steps: 8 });
  const arrow = active.locator('.aim-arrow');
  await expect(arrow).toHaveCount(1);
  await expect(arrow).toHaveClass(/\bon\b/);
  await expect(monster).toHaveAttribute('data-aim-over', '');
  const mine = await middle(await characterOf(active, 'you'));
  const own = await ends(active, '.aim-arrow');
  expect(Math.hypot(own.from.x - mine.x, own.from.y - mine.y)).toBeLessThan(3);
  expect(Math.hypot(own.to.x - m.x, own.to.y - m.y)).toBeLessThan(3);
  await shot(active, 'arrow-attack');

  // The other player sees it from the attacker to the monster, marked.
  const theirs = other.locator('.aim-arrow.others');
  await expect(theirs).toHaveCount(1);
  await expect(other.locator(`[data-zone="${zone}"]`)).toHaveAttribute('data-presence', 'at');
  const attacker = await middle(await characterOf(other, 'other'));
  const there = await middle(other.locator(`[data-zone="${zone}"]`));
  await expect
    .poll(async () => {
      const e = await ends(other, '.aim-arrow.others');
      return Math.hypot(e.from.x - attacker.x, e.from.y - attacker.y) + Math.hypot(e.to.x - there.x, e.to.y - there.y);
    })
    .toBeLessThan(6);
  await shot(other, 'arrow-attack-others');

  // A click picks the monster as before; the arrows go.
  await active.mouse.click(m.x, m.y);
  await expect(arrow).toHaveCount(0);
  const id = Number(zone.replace('obj-', ''));
  await expect
    .poll(() =>
      active.evaluate(() => {
        interface S {
          game: { view: { attack?: { target?: number } } } | null;
        }
        return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view.attack?.target ?? 0;
      }),
    )
    .toBe(id);
  await expect(theirs).toHaveCount(0);
});

// noPresence drops the server's presence frames on the way to the page:
// the other seat's real pointer would replace the injected one (19.48).
function noPresence(page: Page): Promise<void> {
  return page.routeWebSocket(`ws://127.0.0.1:${String(serverPort())}/ws`, (ws) => {
    const server = ws.connectToServer();
    server.onMessage((msg) => {
      if ((JSON.parse(String(msg)) as { type: string }).type !== 'presence') {
        ws.send(msg);
      }
    });
  });
}

test("another player's arrow glides out from its source and lingers when they stop (20.23)", async ({ browser }) => {
  // Animations on: the glide is a motion.
  const [page] = await startGame(browser, ['Ned', 'Oda'], { url: appURL(), open: noPresence });
  if (!page) {
    throw new Error('no pages');
  }
  const arrow = page.locator('.aim-arrow.others');
  const monster = page.getByRole('group', { name: tr('board.monsters') }).locator('.card').first();
  const zone = (await monster.getAttribute('data-zone')) ?? '';
  const source = await middle(await characterOf(page, 'other'));
  const target = await middle(monster);
  // The other seat aims from its character at the monster, as the
  // server relays it; the head is read two frames later.
  const early = await page.evaluate(async (obj) => {
    interface S {
      game: { view: { players: { character: { id: number } }[] } };
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const other = 1 - s.lobby.table.you;
    const from = s.game.view.players[other]?.character.id ?? 0;
    s.receive({ type: 'presence', data: { seat: other, from: { object: from }, at: { object: obj } } });
    const frame = () => new Promise((r) => requestAnimationFrame(r));
    await frame();
    await frame();
    await frame();
    const r = document.querySelector('.aim-arrow.others .aim-ring')?.getBoundingClientRect();
    return { x: r?.left ?? NaN, y: r?.top ?? NaN };
  }, Number(zone.replace('obj-', '')));
  // It starts near its source, far from the monster...
  const way = Math.hypot(target.x - source.x, target.y - source.y);
  expect(Math.hypot(early.x - target.x, early.y - target.y)).toBeGreaterThan(way * 0.5);
  // ...and glides on to it.
  await expect
    .poll(async () => {
      const e = await ends(page, '.aim-arrow.others');
      return Math.hypot(e.to.x - target.x, e.to.y - target.y);
    })
    .toBeLessThan(3);
  await expect(page.locator('.aim-arrow.others .aim-ring > div')).toHaveCount(0);

  // They stop: the arrow stays a moment, then goes.
  const kept = await page.evaluate(async () => {
    interface S {
      lobby: { table: { you: number } };
      receive: (env: unknown) => void;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    s.receive({ type: 'presence', data: { seat: 1 - s.lobby.table.you } });
    const wait = (ms: number) => new Promise((r) => setTimeout(r, ms));
    await wait(100);
    const soon = document.querySelectorAll('.aim-arrow.others').length;
    await wait(400);
    return { soon, later: document.querySelectorAll('.aim-arrow.others').length };
  });
  expect(kept).toEqual({ soon: 1, later: 0 });
  await expect(arrow).toHaveCount(0);
});
