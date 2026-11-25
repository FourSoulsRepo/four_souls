// A card under an item (19.244): Friendly Ball's monster lies behind
// the ball, face up, a little higher, so its top band shows above it,
// for every player, as on a real table. The ball keeps its place in the
// row. Hover on the band zooms the card; the ball's own zoom shows it
// beside the ball. Letting the monster go takes the band away.
import type { Locator, Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { expectFits, startGame } from './game';
import { cardName, tr } from './locale';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle } from './turn';

test.setTimeout(120_000);

interface BallSave {
  state: { objects: { card: string; charged?: boolean }[] };
}

interface Box {
  x: number;
  y: number;
  width: number;
  height: number;
}

// underBall: the cards the page's view has under Friendly Ball.
function underBall(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    interface S {
      game: { view: { players: { in_play: { card: string; under?: string[] }[] | null }[] } } | null;
    }
    const s = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState();
    const ball = s.game?.view.players.flatMap((p) => p.in_play ?? []).find((c) => c.card === 'friendly_ball');
    return ball?.under ?? [];
  });
}

async function boxOf(l: Locator): Promise<Box> {
  const box = await l.boundingBox();
  if (!box) {
    throw new Error('no box');
  }
  return box;
}

// bandShows: the top band of a holder's first card under it shows above
// the card on top, and takes the pointer there: nothing covers it, and
// the row's scroll does not clip it. It returns the band's middle.
async function bandShows(holder: Locator): Promise<{ x: number; y: number }> {
  const peek = holder.locator('.card-under-peek').first();
  const top = await boxOf(holder.locator(':scope > .card'));
  const band = await boxOf(peek);
  // The band lies above the card, in line with it.
  expect(band.y).toBeLessThan(top.y - 8);
  expect(Math.abs(band.x - top.x)).toBeLessThan(1);
  // Several under one card: each lies a band higher than the next.
  const peeks = holder.locator('.card-under-peek');
  const below = (await peeks.count()) > 1 ? await boxOf(peeks.nth(1)) : top;
  const at = { x: band.x + band.width / 2, y: (band.y + below.y) / 2 };
  const hit = await peek.evaluate((el, p) => {
    const found = document.elementFromPoint(p.x, p.y);
    // What is there instead names itself in the failure.
    return found !== null && el.contains(found) ? 'band' : `${found?.tagName ?? ''}.${String(found?.className ?? '')}`;
  }, at);
  expect(hit, 'the band takes the pointer at its middle').toBe('band');
  return at;
}

test('a monster under Friendly Ball shows its top above the ball for everyone (19.244)', async ({ browser }) => {
  // Names no other spec saves under: the server's saved games list is
  // shared by parallel workers, and Continue picks a game by its names.
  const names = ['Nim', 'Pel'];
  // The pages offer Requiem (?requiem), so they can continue a game
  // with Friendly Ball; the deal stays the Base Game's.
  const pages = await startGame(browser, names, {
    url: `${appURL()}&requiem`,
    setup: (host) => host.getByRole('group', { name: tr('host.sets') }).getByLabel(tr('set.requiem')).uncheck(),
  });
  await activePage(pages, click);
  // The starting item becomes Friendly Ball; the first slot holds Fly.
  await saveAndContinue(pages, names.join(', '), 'a_penny', 'friendly_ball', false, { r: 'Requiem' }, undefined, 'fly');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const name = names[pages.indexOf(active)] ?? '';
  const mine = active.getByRole('group', { name: tr('table.your_items') });
  const theirs = other.getByRole('group', { name: tr('table.items_of', { player: name }) });
  const ball = cardName('friendly_ball');
  const fly = cardName('fly');
  const under = (items: Locator) => items.getByRole('group', { name: tr('table.under', { card: ball }) });
  const ballOf = (items: Locator) => items.locator(`.card-holder > .card[aria-label="${ball}"]`);
  const holderOf = (items: Locator) => items.locator('.card-holder').filter({ has: items.page().locator(`:scope > .card[aria-label="${ball}"]`) });
  // Nothing under the ball yet: no band.
  await expect(under(mine)).toHaveCount(0);
  const mineBefore = await boxOf(ballOf(mine));
  const theirsBefore = await boxOf(ballOf(theirs));
  await shot(active, 'under-before');

  // ↷: put Fly under the ball. The other ability needs a card under
  // it, so the use goes straight to its target.
  await click(active, mine.getByRole('button', { name: ball, exact: true }));
  await click(active, active.getByRole('img', { name: fly, exact: true }));
  await settle(pages, click, async () => (await underBall(active)).includes('fly') && (await underBall(other)).includes('fly'));

  // Both players see Fly's top band above the ball; the ball has not
  // moved: it stays in line with the row, as before the catch.
  await expect(under(mine).getByRole('img', { name: fly, exact: true })).toBeVisible();
  await expect(under(theirs).getByRole('img', { name: fly, exact: true })).toBeVisible();
  expect(await boxOf(ballOf(mine))).toEqual(mineBefore);
  expect(await boxOf(ballOf(theirs))).toEqual(theirsBefore);
  const band = await bandShows(holderOf(mine));
  await bandShows(holderOf(theirs));
  await shot(active, 'under-band');
  await shot(other, 'under-band-other');

  // Hover on the band zooms Fly.
  const zoom = active.getByTestId('card-zoom');
  await active.mouse.move(band.x, band.y);
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', fly);
  await shot(active, 'under-band-zoom');
  await active.mouse.move(1, 1);
  await expect(zoom).toHaveCount(0);

  // The ball's zoom shows Fly beside it, as the card under it.
  await ballOf(mine).hover({ position: { x: 10, y: 40 } });
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', ball);
  const side = zoom.locator('.card-zoom-under');
  await expect(side).toContainText(tr('zoom.under'));
  await expect(side.getByRole('img', { name: fly, exact: true })).toBeVisible();
  await shot(active, 'under-ball-zoom');
  await active.mouse.move(1, 1);

  // The ball is spent: ready it again in the save, then let Fly go.
  await saveEditContinue(pages, names.join(', '), (game) => {
    editSave(game, (raw) => {
      for (const o of (raw as unknown as BallSave).state.objects) {
        if (o.card === 'friendly_ball') {
          o.charged = true;
        }
      }
    });
  });
  const again = await activePage(pages, click);
  const rest = pages.find((p) => p !== again);
  if (!rest) {
    throw new Error('no other page');
  }
  const mineAgain = again.getByRole('group', { name: tr('table.your_items') });
  const theirsAgain = rest.getByRole('group', { name: tr('table.items_of', { player: name }) });
  await expect(under(mineAgain).getByRole('img', { name: fly, exact: true })).toBeVisible();
  // Its one usable ability now is the release: a menu, then the slot.
  await click(again, mineAgain.getByRole('button', { name: ball, exact: true }));
  await click(again, again.getByRole('menu').getByRole('menuitem').first());
  await settle(pages, click, async () => (await underBall(again)).length === 0 && (await underBall(rest)).length === 0);

  // Fly is back in a slot; no band above the ball, which stays put.
  await expect(under(mineAgain)).toHaveCount(0);
  await expect(under(theirsAgain)).toHaveCount(0);
  expect(await boxOf(ballOf(mineAgain))).toEqual(again === active ? mineBefore : theirsBefore);
  await expect(again.locator('.board').getByRole('img', { name: fly, exact: true }).or(again.locator('.board').getByRole('button', { name: fly, exact: true })).first()).toBeVisible();
  await shot(again, 'under-released');
});

test('the bands show in every seat, a full row and any window size (19.244)', async ({ browser }) => {
  const pages = await startGame(browser, ['Kai', 'Lea', 'Max', 'Noa']);
  const page = pages[0];
  if (!page) {
    throw new Error('no pages');
  }
  const rows = {
    bottom: page.getByRole('group', { name: tr('table.your_items') }),
    left: page.locator('.player-left .opp-items'),
    top: page.locator('.player-top .opp-items'),
    right: page.locator('.player-right .opp-items'),
  };
  const firstBefore = new Map<string, Box>();
  for (const [seat, row] of Object.entries(rows)) {
    firstBefore.set(seat, await boxOf(row.locator('.card-holder > .card').first()));
  }
  // Every player gets items with cards under one of them, as the server
  // would send them; the viewer gets a full row, scrolling, with cards
  // under its second and last items. Only the screen changes; nobody
  // acts on it. A long name (Conjoined Fatty) wraps on a text card.
  await page.evaluate(() => {
    interface Card {
      id: number;
      card: string;
      role?: number;
      charged: boolean;
      under?: string[];
    }
    interface S {
      game: { view: { players: { id: number; in_play: Card[] | null }[] } } | null;
    }
    const { store } = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (s: Partial<S>) => void } } })
      .fourSouls;
    const { game } = store.getState();
    if (!game) {
      throw new Error('no game');
    }
    const items = ['box', 'chaos', 'boomerang', 'breakfast', 'brimstone', 'bum_friend', 'blood_lust', 'belly_button', 'baby_haunt', 'champion_belt'];
    const you = Number(document.querySelector('.player-bottom')?.getAttribute('data-zone')?.replace('player-', ''));
    const unders = [['conjoined_fatty'], ['fly'], ['fly', 'big_spider']];
    let other = 0;
    const players = game.view.players.map((p) => {
      const own: Card[] = [];
      if (p.id === you) {
        items.forEach((card, i) => {
          const under = i === 1 ? ['big_spider'] : i === items.length - 1 ? ['fly'] : undefined;
          own.push({ id: 9300 + i, card, role: 2, charged: true, ...(under ? { under } : {}) });
        });
      } else {
        own.push({ id: 9310 + other, card: items[other] ?? 'box', role: 2, charged: true, under: unders[other] ?? ['fly'] });
        other++;
      }
      return { ...p, in_play: [...(p.in_play ?? []), ...own] };
    });
    store.setState({ game: { ...game, view: { ...game.view, players } } });
  });
  // The new cards pop in (card-pop) before the checks.
  await page.waitForTimeout(400);

  const check = async (size: string) => {
    await expectFits(page);
    for (const [seat, row] of Object.entries(rows)) {
      const side = seat === 'left' || seat === 'right';
      // A top or bottom row's first card keeps its height: the bands
      // take no room. A side column makes room under the character, so
      // a band never hides its stats; the row moves down.
      if (size === 'default' && !side) {
        expect((await boxOf(row.locator('.card-holder > .card').first())).y, seat).toBe(firstBefore.get(seat)?.y);
      }
      if (side) {
        const character = await boxOf(page.locator(`.player-${seat} .opp-character > .card`));
        const band = await boxOf(row.locator('.card-under-peek').first());
        expect(band.y, seat).toBeGreaterThanOrEqual(character.y + character.height);
      }
      // Every card on top stays in line with the row's first card.
      const first = await boxOf(row.locator('.card-holder > .card').first());
      const holder = row.locator('.card-holder').filter({ has: page.locator('.card-under') }).first();
      await holder.scrollIntoViewIfNeeded();
      expect(Math.abs((await boxOf(holder.locator(':scope > .card'))).y - first.y), seat).toBeLessThan(1);
      await bandShows(holder);
    }
    // The full row scrolls sideways; its last card's band shows too.
    const own = rows.bottom;
    expect(await own.evaluate((el) => el.scrollWidth > el.clientWidth)).toBe(true);
    const last = own.locator('.card-holder').last();
    await last.scrollIntoViewIfNeeded();
    await bandShows(last);
    await own.evaluate((el) => {
      el.scrollLeft = 0;
    });
  };
  await check('default');
  await shot(page, 'under-seats');

  // Hover on the top player's band zooms the card under it.
  const topHolder = rows.top.locator('.card-holder').filter({ has: page.locator('.card-under') }).first();
  const topName = (await topHolder.locator('.card-under-peek .card').first().getAttribute('aria-label')) ?? '';
  expect(topName).not.toBe('');
  const top = await bandShows(topHolder);
  await page.mouse.move(top.x, top.y);
  await expect(page.getByTestId('card-zoom').locator('.card.size-zoom')).toHaveAttribute('aria-label', topName);
  await shot(page, 'under-seats-zoom');
  await page.mouse.move(1, 1);

  await page.setViewportSize({ width: 1024, height: 600 });
  await check('small');
  await shot(page, 'under-seats-small');
  await page.setViewportSize({ width: 1920, height: 1080 });
  await check('large');
  await shot(page, 'under-seats-large');
});
