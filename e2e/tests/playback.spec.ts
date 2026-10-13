// Event playback and animations (7.6).
import { expect, test } from '@playwright/test';
import { startGame } from './game';
import { shot } from './shot';

test('the animation speed setting survives a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Animations').selectOption('fast');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Animations')).toHaveValue('fast');
});

test('a fast burst of updates plays in order and ends on the last one', async ({ browser }) => {
  const [kai] = await startGame(browser, ['Kai', 'Lu']);
  if (!kai) {
    throw new Error('no pages');
  }
  // Mark one card: if updates redrew the table from scratch, the mark
  // would be gone and every card would pop in again.
  const character = kai.locator('.player-bottom .player-cards > .card').first();
  await character.evaluate((el) => el.setAttribute('data-mark', 'kept'));
  // Replay the current view 20 times as new turns, each with a dice roll
  // and a looted card, as fast as the page takes them.
  await kai.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { number: number; active: number } } } | null;
        lobby: { table: { you: number } | null };
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const you = s.lobby.table?.you ?? 0;
    for (let n = 2; n <= 21; n++) {
      const game = structuredClone(s.game);
      if (!game) {
        throw new Error('no game');
      }
      game.view.turn.number = n;
      s.receive({
        type: 'update',
        data: {
          ...game,
          step: n,
          events: [
            { kind: 'dice_rolled', player: you, amount: (n % 6) + 1, text: 'roll' },
            { kind: 'looted', player: you, object: 1000 + n },
          ],
        },
      });
    }
  });
  await expect(kai.getByTestId('effect-dice')).toBeVisible();
  await kai.waitForTimeout(300); // let the dice settle for the picture
  await shot(kai, 'effect-dice');
  await expect(kai.getByTestId('effect-fly')).toBeVisible();
  await expect(kai.getByRole('main', { name: 'Turn 21' })).toBeVisible({ timeout: 20_000 });
  await expect(kai.locator('.effects')).toHaveCount(0);
  await expect(character).toHaveAttribute('data-mark', 'kept');
  await expect(character).toHaveCSS('opacity', '1');
});

test('an attack lunges at the monster and drops its HP at the hit', async ({ browser }) => {
  const [vee] = await startGame(browser, ['Vee', 'Wyn']);
  if (!vee) {
    throw new Error('no pages');
  }
  // Replay the view with one hit of the own character on the first
  // monster, whose HP is one lower after it.
  const ids = await vee.evaluate(() => {
    interface Card {
      id: number;
      stats?: { hp: number };
    }
    interface Store {
      getState: () => {
        game: { view: { players: { id: number; character: Card }[]; monsters: { top?: Card }[] } } | null;
        lobby: { table: { you: number } | null };
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const character = game?.view.players.find((p) => p.id === you)?.character;
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !character || !monster?.stats) {
      throw new Error('no game');
    }
    const hp = monster.stats.hp;
    monster.stats.hp = hp - 1;
    s.receive({
      type: 'update',
      data: {
        ...game,
        step: 2,
        events: [{ kind: 'damaged', player: -1, object: monster.id, amount: 1, source: character.id, text: 'combat' }],
      },
    });
    return { character: character.id, monster: monster.id, hp };
  });
  const character = vee.locator(`[data-zone="obj-${String(ids.character)}"]`);
  const monster = vee.locator(`[data-zone="obj-${String(ids.monster)}"]`);
  // The card's copy flies; the card itself waits hidden at home.
  await expect(vee.getByTestId('lunge')).toBeVisible();
  await expect(character).toBeHidden();
  if (process.env.SCREENSHOTS) {
    // Hold the copy on its way into the monster for the picture.
    const lunge = vee.getByTestId('lunge');
    await lunge.evaluate((el) => {
      for (const a of el.getAnimations()) {
        a.pause();
        a.currentTime = 0.42 * Number(a.effect?.getTiming().duration ?? 0);
      }
    });
    await shot(vee, 'attack-lunge');
    await lunge.evaluate((el) => {
      for (const a of el.getAnimations()) {
        a.play();
      }
    });
  }
  await expect(vee.getByTestId('burst')).toHaveText('−1');
  await vee.waitForTimeout(200); // the burst pops
  await shot(vee, 'attack-hit');
  await expect(monster.getByLabel(`${String(ids.hp - 1)} health`)).toBeVisible();
  await expect(vee.locator('.effects-trail > *')).toHaveCount(0);
  await expect(character).toBeVisible();
});

// The view parts the replays below change.
interface Card {
  id: number;
  card: string;
  damage?: number;
  stats?: { hp: number };
}
interface View {
  players: { id: number; hp: number; character: Card }[];
  monsters: { top?: Card }[];
}
interface Store {
  getState: () => {
    game: { view: View } | null;
    lobby: { table: { you: number } | null };
    receive: (env: unknown) => void;
  };
}

test('damage from a card hits every target at once and reddens their HP', async ({ browser }) => {
  const [mo] = await startGame(browser, ['Mo', 'Nia']);
  if (!mo) {
    throw new Error('no pages');
  }
  // The own character's ability deals 1 damage to two monsters and to
  // its own player.
  const ids = await mo.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const me = game?.view.players.find((p) => p.id === you);
    const tops = (game?.view.monsters ?? []).flatMap((m) => (m.top?.stats ? [m.top] : [])).slice(0, 2);
    if (!game || !me || tops.length < 2) {
      throw new Error('no game');
    }
    const hps = tops.map((t) => t.stats?.hp ?? 0);
    const events: unknown[] = tops.map((t) => {
      if (t.stats) {
        t.stats.hp -= 1;
      }
      t.damage = 1;
      return { kind: 'damaged', player: -1, object: t.id, amount: 1, source: me.character.id };
    });
    me.hp -= 1;
    me.character.damage = 1;
    events.push({ kind: 'damaged', player: you, amount: 1, source: me.character.id });
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
    return { character: me.character.id, monsters: tops.map((t) => t.id), hps };
  });
  // A shot flies from the source to each target.
  await expect(mo.getByTestId('shot')).toHaveCount(3);
  if (process.env.SCREENSHOTS) {
    // Hold the shots on their way for the picture.
    await mo.evaluate(() => {
      for (const a of document.getAnimations()) {
        a.pause();
        a.currentTime = 0.6 * Number(a.effect?.getTiming().duration ?? 0);
      }
    });
    await shot(mo, 'damage-shots');
    await mo.evaluate(() => {
      for (const a of document.getAnimations()) {
        a.play();
      }
    });
  }
  // Every target bursts at once.
  await expect(mo.getByTestId('burst')).toHaveCount(3);
  await mo.waitForTimeout(200); // the bursts pop
  await shot(mo, 'damage-bursts');
  for (const [i, id] of ids.monsters.entries()) {
    const health = mo.locator(`[data-zone="obj-${String(id)}"]`).getByLabel(`${String((ids.hps[i] ?? 0) - 1)} health`);
    await expect(health).toHaveClass(/hurt/);
  }
  await expect(mo.locator(`[data-zone="obj-${String(ids.character)}"] .card-hp`)).toHaveClass(/hurt/);
  await expect(mo.locator('.effects-trail > *')).toHaveCount(0);
  await shot(mo, 'damage-red-hp');
});

test('a dying monster gets a skull before it leaves', async ({ browser }) => {
  const [ola] = await startGame(browser, ['Ola', 'Pim']);
  if (!ola) {
    throw new Error('no pages');
  }
  const id = await ola.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const slot = game?.view.monsters.find((m) => m.top?.stats);
    const top = slot?.top;
    if (!game || !slot || !top) {
      throw new Error('no game');
    }
    slot.top = undefined;
    const events = [{ kind: 'died', player: -1, object: 999, card: top.card, prev: top.id }];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
    return top.id;
  });
  const monster = ola.locator(`[data-zone="obj-${String(id)}"]`);
  await expect(ola.getByTestId('skull')).toBeVisible();
  await expect(monster).toBeVisible();
  await ola.waitForTimeout(250); // the skull pops
  await shot(ola, 'death-skull');
  await expect(monster).toHaveCount(0);
  await expect(ola.locator('.effects-trail > *')).toHaveCount(0);
});
