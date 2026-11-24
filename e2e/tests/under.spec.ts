// A card under an item (19.244): Friendly Ball's monster peeks out from
// under the ball, face up, for every player, as on a real table. Hover
// on the peeking card zooms it; the ball's own zoom shows it beside the
// ball. Letting the monster go takes the peek away.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { cardName, tr } from './locale';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, click, settle } from './turn';

test.setTimeout(120_000);

interface BallSave {
  state: { objects: { card: string; charged?: boolean }[] };
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

test('a monster under Friendly Ball peeks out for everyone (19.244)', async ({ browser }) => {
  const names = ['Ivo', 'Jun'];
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
  const under = (items: typeof mine) => items.getByRole('group', { name: tr('table.under', { card: ball }) });
  // Nothing under the ball yet: no peek.
  await expect(under(mine)).toHaveCount(0);
  await shot(active, 'under-before');

  // ↷: put Fly under the ball. The other ability needs a card under
  // it, so the use goes straight to its target.
  await click(active, mine.getByRole('button', { name: ball, exact: true }));
  await click(active, active.getByRole('img', { name: fly, exact: true }));
  await settle(pages, click, async () => (await underBall(active)).includes('fly') && (await underBall(other)).includes('fly'));

  // Both players see Fly peeking out from under the ball.
  await expect(under(mine).getByRole('img', { name: fly, exact: true })).toBeVisible();
  await expect(under(theirs).getByRole('img', { name: fly, exact: true })).toBeVisible();
  await shot(active, 'under-peek');
  await shot(other, 'under-peek-other');

  // Hover on the peeking part zooms Fly.
  const zoom = active.getByTestId('card-zoom');
  const peek = under(mine).locator('.card');
  const box = await peek.boundingBox();
  if (!box) {
    throw new Error('the peeking card has no box');
  }
  await peek.hover({ position: { x: box.width - 6, y: box.height / 2 } });
  await expect(zoom.locator('.card.size-zoom')).toHaveAttribute('aria-label', fly);
  await shot(active, 'under-peek-zoom');
  await active.mouse.move(1, 1);
  await expect(zoom).toHaveCount(0);

  // The ball's zoom shows Fly beside it, as the card under it.
  await mine.getByRole('img', { name: ball, exact: true }).or(mine.getByRole('button', { name: ball, exact: true })).hover({ position: { x: 10, y: 40 } });
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

  // Fly is back in a slot; nothing peeks out from under the ball.
  await expect(under(mineAgain)).toHaveCount(0);
  await expect(under(theirsAgain)).toHaveCount(0);
  await expect(again.locator('.board').getByRole('img', { name: fly, exact: true }).or(again.locator('.board').getByRole('button', { name: fly, exact: true })).first()).toBeVisible();
  await shot(again, 'under-released');
});
